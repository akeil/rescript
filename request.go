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
	mmPerInch          = 25.4
	singlePointerID    = -1
)

// Request is the root element for requests against the batch enpoint.
type Request struct {
	Width  int64 `json:"width"`
	Height int64 `json:"height"`
	// ContentType controls the "recognition type" of the MyScript API.
	// It must be one of Text, Diagram, Math, Raw Content, Text Document
	ContentType string `json:"contentType"`
	// Must be DIGITAL_EDIT when the ReST API is used.
	ConversionState string        `json:"conversionState"`
	XDpi            int64         `json:"xDPI"`
	YDpi            int64         `json:"yDPI"`
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
	binary.Write(h, binary.LittleEndian, r.Width)
	binary.Write(h, binary.LittleEndian, r.Height)
	binary.Write(h, binary.LittleEndian, r.ConversionState)
	binary.Write(h, binary.LittleEndian, r.ContentType)
	binary.Write(h, binary.LittleEndian, r.XDpi)
	binary.Write(h, binary.LittleEndian, r.YDpi)
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
//
// The Timestamp can be based on "0" as long as it increases from point to point.
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
		binary.Write(h, binary.LittleEndian, int64(s.X[i]))
		binary.Write(h, binary.LittleEndian, int64(s.Y[i]))
		binary.Write(h, binary.LittleEndian, s.Pressure[i])
		binary.Write(h, binary.LittleEndian, s.Timestamp[i])
	}
}

// Configuration --------------------------------------------------------------

// Configuration is the root object for configuration options to a batch call.
type Configuration struct {
	Language   LanguageCode             `json:"lang"`
	Text       *TextConfiguration       `json:"text,omitempty"`
	Export     *ExportConfiguration     `json:"export,omitempty"`
	RawContent *RawContentConfiguration `json:"raw-content,omitempty"`
}

// NewConfiguration creates a new configuration object with the given values.
func NewConfiguration(lang LanguageCode, guides, bbox, chars, words bool) Configuration {
	return Configuration{
		Language:   lang,
		Text:       NewTextConfiguration(guides),
		Export:     NewExportConfiguration(bbox, chars, words),
		RawContent: NewRawContentConfiguration(),
	}
}

func (c Configuration) checksum(h hash.Hash) {
	h.Write([]byte(c.Language))
	c.Text.checksum(h)
	c.Export.checksum(h)
	c.RawContent.checksum(h)
}

// TextConfiguration holds settings holds settings for text recognition.
type TextConfiguration struct {
	Guides        GuidesConfiguration  `json:"guides"`
	Margin        MarginConfiguration  `json:"margin"`
	Configuration RawTextConfiguration `json:"configuration"`
}

// NewTextConfiguration creates a default configuration.
func NewTextConfiguration(guides bool) *TextConfiguration {
	return &TextConfiguration{
		Guides: GuidesConfiguration{Enable: guides},
		Margin: MarginConfiguration{
			Top:    0,
			Left:   0,
			Right:  0,
			Bottom: 0,
		},
		Configuration: RawTextConfiguration{
			AddLKText: true, // defult: true
		},
	}
}

func (t *TextConfiguration) checksum(h hash.Hash) {
	binary.Write(h, binary.LittleEndian, t.Guides.Enable)
	binary.Write(h, binary.LittleEndian, t.Configuration.AddLKText)
	binary.Write(h, binary.LittleEndian, t.Margin.Top)
	binary.Write(h, binary.LittleEndian, t.Margin.Left)
	binary.Write(h, binary.LittleEndian, t.Margin.Right)
	binary.Write(h, binary.LittleEndian, t.Margin.Bottom)
}

type MarginConfiguration struct {
	Top    int32 `json:"top"`
	Left   int32 `json:"left"`
	Right  int32 `json:"right"`
	Bottom int32 `json:"bottom"`
}

type GuidesConfiguration struct {
	Enable bool `json:"enable"`
}

// ExportConfiguration holds settings for the API response.
type ExportConfiguration struct {
	Jiix            JiixConfiguration `json:"jiix"`
	ImageResolution int64             `json:"image-resolution"`
}

// NewExportConfiguration creates an export config.
func NewExportConfiguration(bbox, chars, words bool) *ExportConfiguration {
	return &ExportConfiguration{
		ImageResolution: 300,
		Jiix:            NewJiixConfiguration(bbox, chars, words),
	}
}

func (e *ExportConfiguration) checksum(h hash.Hash) {
	binary.Write(h, binary.LittleEndian, e.ImageResolution)
	e.Jiix.checksum(h)
}

// JiixConfiguration holds options for the JIIX format.
//
// JIIX is the "deep dive" format from the MyScript API.
type JiixConfiguration struct {
	Strokes     bool     `json:"strokes"`
	BoundingBox bool     `json:"bounding-box"`
	Style       bool     `json:"style"`
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
		Strokes:     true, // not sure what this param does (no effect?)
		BoundingBox: bbox,
		Style:       true, // not sure what this does
		Text: JiixText{
			Chars: chars,
			Words: words,
		},
	}
}

func (j JiixConfiguration) checksum(h hash.Hash) {
	binary.Write(h, binary.LittleEndian, j.Strokes)
	binary.Write(h, binary.LittleEndian, j.BoundingBox)
	binary.Write(h, binary.LittleEndian, j.Text.Chars)
	binary.Write(h, binary.LittleEndian, j.Text.Words)
}

type JiixText struct {
	Chars bool `json:"chars"`
	Words bool `json:"words"`
}

type RawContentConfiguration struct {
	Recognition RawRecognitionConfiguration `json:"recognition"`
	Text        RawTextConfiguration        `json:"text"`
}

func NewRawContentConfiguration() *RawContentConfiguration {
	return &RawContentConfiguration{
		Recognition: RawRecognitionConfiguration{
			Text:  true,
			Shape: true,
		},
		Text: RawTextConfiguration{
			AddLKText: true, // defult: true
		},
	}
}

func (r *RawContentConfiguration) checksum(h hash.Hash) {
	binary.Write(h, binary.LittleEndian, r.Recognition.Text)
	binary.Write(h, binary.LittleEndian, r.Recognition.Shape)
	binary.Write(h, binary.LittleEndian, r.Text.AddLKText)
}

type RawRecognitionConfiguration struct {
	Text  bool `json:"text"`
	Shape bool `json:"shape"`
}

type RawTextConfiguration struct {
	CustomResources []string `json:"customResources,omitempty"`
	CustomLexicon   []string `json:"customLexicon,omitempty"`
	AddLKText       bool     `json:"addLKText"`
}
