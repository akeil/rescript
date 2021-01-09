package rescript

import (
	"encoding/binary"
	"hash"
)

type LanguageCode string

type PointerType string

const (
	LangEN LanguageCode = "en_US"
	LangDE LanguageCode = "de_DE"

	Pen    PointerType = "PEN"
	Touch  PointerType = "TOUCH"
	Eraser PointerType = "ERASER"

	defaultContentType = "Text"
	defaultConversion  = "DIGITAL_EDIT"
	defaultPenStyle    = "color: #000000; -myscript-pen-width: ;"
	defaultResolution  = 96
)

// Request is the root element for requests against the batch enpoint.
type Request struct {
	Width           int           `json:"width"`
	Height          int           `json:"height"`
	ContentType     string        `json:"contentType"`
	ConversionState string        `json:"conversionState"`
	XDpi            int           `json:"xDPI"`
	YDpi            int           `json:"yDPI"`
	Theme           string        `json:"theme"`
	StrokeGroups    []StrokeGroup `json:"strokeGroups"`
	Configuration   Configuration `json:"configuration"`
}

// NewRequest creates a request with default values.
func NewRequest() Request {
	return Request{
		ContentType:     defaultContentType,
		ConversionState: defaultConversion,
		XDpi:            defaultResolution,
		YDpi:            defaultResolution,
		StrokeGroups:    make([]StrokeGroup, 0),
	}
}

func (r Request) checksum(h hash.Hash) {
	r.Configuration.checksum(h)
	for _, sg := range r.StrokeGroups {
		sg.checksum(h)
	}
}

// A StrokeGroup contains the digital ink strokes, i.e. the actual handwriting.
type StrokeGroup struct {
	PenStyle string   `json:"penStyle"`
	Strokes  []Stroke `json:"strokes"`
}

// NewStrokeGroup creates and empty stroke group.
func NewStrokeGroup() StrokeGroup {
	return StrokeGroup{
		PenStyle: defaultPenStyle,
		Strokes:  make([]Stroke, 0),
	}
}

// checksum is used internally to determine the cache key.
// It calculates a checksum over all relevant request parameters.
func (s StrokeGroup) checksum(h hash.Hash) {
	for _, st := range s.Strokes {
		st.checksum(h)
	}
}

// A Stroke is a single stroke of digital ink.
//
// It consists of a series of X,Y coordinates and their related
// timestamps and optional pressure values.
type Stroke struct {
	ID          string      `json:"id,omitempty"`          // opt
	PointerType PointerType `json:"pointerType,omitempty"` // opt
	PointerID   int         `json:"pointerId,omitempty"`   // opt
	X           []int       `json:"x"`
	Y           []int       `json:"y"`
	Timestamp   []int64     `json:"t"`
	Pressure    []float64   `json:"p"` // opt
}

// NewStroke creates a stroke with default values and no points.
func NewStroke() Stroke {
	return Stroke{
		PointerType: Pen,
		PointerID:   1,
		X:           make([]int, 0),
		Y:           make([]int, 0),
		Timestamp:   make([]int64, 0),
		Pressure:    make([]float64, 0),
	}
}

func (s Stroke) checksum(h hash.Hash) {
	h.Write([]byte(s.PointerType))
	for i := 0; i < len(s.X); i++ {
		hashInt(h, s.X[i])
		hashInt(h, s.Y[i])
		hashFloat(h, s.Pressure[i])
		// Timestamp is based on Now() and thus different every time
		//binary.Write(h, binary.LittleEndian, s.Timestamp[i])
	}
}

// Configuration --------------------------------------------------------------

// Configuration is the root object for configuration options to a batch call.
type Configuration struct {
	Language LanguageCode         `json:"lang"`
	Text     *TextConfiguration   `json:"text,omitempty"`
	Export   *ExportConfiguration `json:"export,omitempty"`
}

// NewConfiguration creates a new configuration object with the given values.
func NewConfiguration(lang LanguageCode, bbox, chars, words bool) Configuration {
	return Configuration{
		Language: lang,
		Text:     NewTextConfiguration(),
		Export:   NewExportConfiguration(bbox, chars, words),
	}
}

func (c Configuration) checksum(h hash.Hash) {
	h.Write([]byte(c.Language))
	c.Export.checksum(h)
}

// TextConfiguration holds settings holds settings for text recognition.
type TextConfiguration struct {
	Guides GuidesConfiguration `json:"guides"`
	// Margin
	MimeTypes         []string             `json:"mimeTypes"`
	SmartGuide        bool                 `json:"smartGuide"`
	SmartGuideFadeout FadeoutConfiguration `json:"smartGuideFadeout"`
}

// NewTextConfiguration creates a default configuration.
func NewTextConfiguration() *TextConfiguration {
	return &TextConfiguration{
		Guides:            GuidesConfiguration{Enable: true},
		MimeTypes:         []string{"text/plain", "application/vnd.myscript.jiix"},
		SmartGuide:        true,
		SmartGuideFadeout: FadeoutConfiguration{Enable: false, Duration: 10000},
	}
}

type GuidesConfiguration struct {
	Enable bool `json:"enable"`
}

type FadeoutConfiguration struct {
	Enable   bool `json:"enable"`
	Duration int  `json:"duration"`
}

// ExportConfiguration holds settings for the API response.
type ExportConfiguration struct {
	Jiix            JiixConfiguration `json:"jiix"`
	ImageResolution int               `json:"image-resolution"`
}

// NewExportConfiguration creates an export config.
func NewExportConfiguration(bbox, chars, words bool) *ExportConfiguration {
	return &ExportConfiguration{
		ImageResolution: 300,
		Jiix:            NewJiixConfiguration(bbox, chars, words),
	}
}

func (e *ExportConfiguration) checksum(h hash.Hash) {
	e.Jiix.checksum(h)
}

// JiixConfiguration holds options for the JIIX format.
//
// JIIX is the "deep dive" format from the MyScript API.
type JiixConfiguration struct {
	Strokes     bool     `json:"strokes"`
	BoundingBox bool     `json:"bounding-box"`
	Text        JiixText `json:"text"`
}

// NewJiixConfiguration creates a JIIX configuration with the given flags.
//
// If bbox is set, each recognition result is additionally described by a
// bounding box which refers to the input drawing.
//
// If words is set, the result will contain additional entries for each word,
// if chars is also set, the character index for each word is included.
func NewJiixConfiguration(bbox, chars, words bool) JiixConfiguration {
	return JiixConfiguration{
		Strokes:     false, // not sure what this param does
		BoundingBox: bbox,
		Text: JiixText{
			Chars: chars,
			Words: words,
		},
	}
}

func (j JiixConfiguration) checksum(h hash.Hash) {
	hashBool(h, j.Strokes)
	hashBool(h, j.BoundingBox)
	hashBool(h, j.Text.Chars)
	hashBool(h, j.Text.Words)
}

type JiixText struct {
	Chars bool `json:"chars"`
	Words bool `json:"words"`
}

// cChecksum Helpers ----------------------------------------------------------

func hashBool(h hash.Hash, b bool) {
	binary.Write(h, binary.LittleEndian, b)
}

func hashInt(h hash.Hash, i int) {
	binary.Write(h, binary.LittleEndian, int64(i))
}

func hashFloat(h hash.Hash, f float64) {
	binary.Write(h, binary.LittleEndian, f)
}
