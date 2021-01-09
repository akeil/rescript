package rescript

// Result is the response returned by the MyScript batch enpoint.
//
// The Label field contains the complete recognized text.
// If the "words" option was enabled, the liust of Words contains
// the individual words and whitespace.
type Result struct {
	Type        string      `json:"type"`
	BoundingBox BoundingBox `json:"bounding-box"`
	Label       string      `json:"label"`
	Words       []Word      `json:"words"`
}

type Word struct {
	Label       string      `json:"label"`
	Candidates  []string    `json:"candidates,omitempty"`
	FirstChar   int         `json:"first-char,omitempty"`
	LastChar    int         `json:"last-char,omitempty"`
	BoundingBox BoundingBox `json:"bounding-box,omitempty"`
	Items       []Item      `json:"item,omitempty"`
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

type BoundingBox struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

func (b BoundingBox) IsZero() bool {
	return b.X == 0 && b.Y == 0 && b.Width == 0 && b.Height == 0
}
