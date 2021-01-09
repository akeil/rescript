package rescript

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	batchEndpoint = "/api/v4.0/iink/batch"
)

// MyScript is the client for the MyScript ReST API.
type MyScript struct {
	appKey string
	host   string
	client *http.Client
	sign   func(data []byte) string
}

// NewMyScript sets up a new client.
//
// It requires the application key and the HMAC key from ypur MyScript account.
func NewMyScript(appKey, hmacKey string) *MyScript {
	return &MyScript{
		appKey: appKey,
		host:   "https://cloud.myscript.com",
		client: &http.Client{},
		sign: func(data []byte) string {
			// see:
			// https://developer.myscript.com/support/account/registering-myscript-cloud/#computing-the-hmac-value

			// our "user key"
			key := []byte(appKey + hmacKey)
			// create a SHA-512 Mac
			mac := hmac.New(sha512.New, key)

			// the hmac over the payload data
			mac.Write(data)

			return hex.EncodeToString(mac.Sum(nil))
		},
	}
}

// Batch is the single endpoint fif the ReST API.
// It performs handwriting recognition.
func (m *MyScript) Batch(r Request) (Result, error) {
	var result Result
	// We need the JSON body as []byte because we need to create a signature over it.
	payload, err := json.Marshal(r)
	if err != nil {
		return result, err
	}

	u, err := m.resolveEndpoint(batchEndpoint)
	if err != nil {
		return result, err
	}

	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(payload))
	if err != nil {
		return result, err
	}

	// MyScript custom headers
	req.Header.Add("applicationKey", m.appKey)
	req.Header.Add("hmac", m.sign(payload))
	// Standard headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept", "application/vnd.myscript.jiix")

	res, err := m.client.Do(req)
	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		// TODO: Error Model
		e := make(map[string]interface{})
		err = json.NewDecoder(res.Body).Decode(&e)
		if err != nil {
			return result, err
		}
		fmt.Println(e)
		return result, fmt.Errorf("bad status code %v", res.StatusCode)
	}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (m *MyScript) resolveEndpoint(ep string) (*url.URL, error) {
	base, err := url.Parse(m.host)
	if err != nil {
		return nil, err
	}

	ref, err := url.Parse(ep)
	if err != nil {
		return nil, err
	}
	return base.ResolveReference(ref), nil
}
