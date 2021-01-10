package rescript

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/akeil/rmtool"
	"github.com/akeil/rmtool/pkg/lines"
)

// The Recognizer organizes calls to the MyScript API to convert notbooks
// from handwriting to a recognize Result.
//
// The recognizer also manages caching to avoid repeated calls to the API
// if a page has not changed.
type Recognizer struct {
	ms       *MyScript
	cacheDir string
	cacheMx  sync.RWMutex
}

// NewRecognizer creates a recognizer withthe given credentials for the
// MyScript API.
//
// If cacheDir is non-empty, it will be used to cache responses from the API.
// If it is empty, caching is disabled.
func NewRecognizer(appKey, hmacKey, cacheDir string) *Recognizer {
	return &Recognizer{
		ms:       NewMyScript(appKey, hmacKey),
		cacheDir: cacheDir,
	}
}

// Recognize performs handwriting recognition on all pages of the given document.
// It resturns a map of page-IDs and recognition results.
func (r *Recognizer) Recognize(doc *rmtool.Document, l LanguageCode) (map[string]Result, error) {
	var resultsMx sync.Mutex
	results := make(map[string]Result)

	var group errgroup.Group
	for _, p := range doc.Pages() {
		pageID := p
		group.Go(func() error {
			d, err := doc.Drawing(pageID)
			if err != nil {
				return err
			}
			res, err := r.recognizeDrawing(d, l)
			if err != nil {
				return err
			}
			resultsMx.Lock()
			results[pageID] = res
			resultsMx.Unlock()
			return nil
		})
	}

	err := group.Wait()
	if err != nil {
		return results, err
	}

	return results, nil
}

func (r *Recognizer) recognizeDrawing(d *lines.Drawing, l LanguageCode) (Result, error) {
	groups := make([]StrokeGroup, len(d.Layers))
	t := int64(0)
	for i, l := range d.Layers {
		g, tx := ConvertLayer(t, l)
		t = tx
		groups[i] = g
	}

	req := prepareRequest(l)
	req.StrokeGroups = groups

	k, err := cacheKey(req)
	if err == nil {
		cached, err := r.readCache(k)
		if err == nil {
			return cached, nil
		}
	}

	res, err := r.ms.Batch(req)
	if err != nil {
		return res, err
	}

	if k != "" {
		go r.writeCache(k, res)
	}

	return res, err
}

func (r *Recognizer) readCache(key string) (Result, error) {
	var res Result

	if r.cacheDir == "" {
		return res, fmt.Errorf("cache dir not set")
	}

	r.cacheMx.RLock()
	defer r.cacheMx.RUnlock()

	p := filepath.Join(r.cacheDir, key+".cache.json")
	f, err := os.Open(p)
	if err != nil {
		return res, err
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *Recognizer) writeCache(key string, res Result) error {
	if r.cacheDir == "" {
		return fmt.Errorf("cache dir not set")
	}

	r.cacheMx.Lock()
	defer r.cacheMx.Unlock()

	err := os.MkdirAll(r.cacheDir, 0755)
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
	}
	p := filepath.Join(r.cacheDir, key+".cache.json")
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(res)
}

func prepareRequest(l LanguageCode) Request {
	req := NewRequest()
	req.Width = lines.MaxWidth
	req.Height = lines.MaxHeight
	guides := false // recommended to turn off in Offscreen usage
	bbox := true
	chars := false
	words := true
	req.Configuration = NewConfiguration(l, guides, bbox, chars, words)

	return req
}

func cacheKey(req Request) (string, error) {
	cs := sha1.New()
	req.checksum(cs)
	return hex.EncodeToString(cs.Sum(nil)), nil
}
