package rescript

// Result is the response returned by the MyScript batch enpoint.
//
// The Label field contains the complete recognized text.
// If the "words" option was enabled, the liust of Words contains
// the individual words and whitespace.
type Result struct {
	ID          string      `json:"id"`
	Version     string      `json:"version"`
	Type        string      `json:"type"`
	Label       string      `json:"label"`
	BoundingBox BoundingBox `json:"bounding-box"`
	Words       []Word      `json:"words"`
	Chars       []Char      `json:"chars"`
	Linebreaks  []Linebreak `json:"linebreaks"`
}

// Word is a single recognized "word", including whitespace or punctuation.
//
// The recognized content is held in the `Label`; concatenating all labesl
// gives the full text.
// See:
// https://developer.myscript.com/docs/interactive-ink/1.4/reference/web/jiix/#word-object
type Word struct {
	Label       string      `json:"label"`
	ReflowLabel string      `json:"reflow-label"`
	FirstChar   int         `json:"first-char,omitempty"`
	LastChar    int         `json:"last-char,omitempty"`
	BoundingBox BoundingBox `json:"bounding-box,omitempty"`
	Candidates  []string    `json:"candidates,omitempty"`
	Items       []Item      `json:"items,omitempty"`
}

type Item struct {
	ID              string      `json:"id"`
	Type            string      `json:"type"`
	Timestamp       string      `json:"timestamp"` // 2021-01-09 13:23:42.196250
	Label           string      `json:"label"`
	Baseline        float64     `json:"baseline"`
	XHeight         float64     `json:"x-height"`
	LeftSideBearing float64     `json:"left-side-bearing"`
	BoundingBox     BoundingBox `json:"bounding-box"`
}

type Char struct {
	Label       string      `json:"label"`
	Word        int         `json:"word"`
	Grid        []Point     `json:"grid"`
	BoundingBox BoundingBox `json:"bounding-box,omitempty"`
	Items       []Item      `json:"items,omitempty"`
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Linebreak struct {
	Line int `json:"line"`
}

// Coordinates are in **millimeters**
type BoundingBox struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

func (b BoundingBox) IsZero() bool {
	return b.X == 0 && b.Y == 0 && b.Width == 0 && b.Height == 0
}
