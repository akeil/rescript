package rescript

import (
	"math"
	"time"

	"github.com/akeil/rmtool/pkg/lines"
)

const (
	speedFactor int64 = 200
	strokeGap         = 500
	minSpeed          = 0.01
)

// ConvertLayer convert a Layer from a reMarkable drawing to a MyScript stroke group.
func ConvertLayer(l lines.Layer) StrokeGroup {
	t := time.Now()
	strokes := make([]Stroke, len(l.Strokes))

	i := 0
	for _, s := range l.Strokes {
		if isTextStroke(s.BrushType) {
			stroke, tx := convertStroke(t, s)
			strokes[i] = stroke
			// add some millis to t for each new stroke
			t = tx.Add(time.Millisecond * strokeGap)
			i++
		}
	}

	return StrokeGroup{
		Strokes:  strokes[:i],
		PenStyle: defaultPenStyle,
	}
}

func convertStroke(t time.Time, s lines.Stroke) (Stroke, time.Time) {
	size := len(s.Dots)
	x := make([]int, size)
	y := make([]int, size)
	ts := make([]int64, size)
	p := make([]float64, size)
	ms := toMillis(t)

	x0 := -1
	y0 := -1
	i := 0
	for _, dot := range s.Dots {
		x1 := int(math.Round(float64(dot.X)))
		y1 := int(math.Round(float64(dot.Y)))
		// avoid duplicate points
		if x0 != x1 || y0 != y1 {
			s := math.Max(minSpeed, float64(dot.Speed))
			offset := float64(speedFactor) / s
			ms += int64(math.Round(offset))

			x[i] = int(math.Round(float64(dot.X)))
			y[i] = int(math.Round(float64(dot.Y)))
			ts[i] = ms
			p[i] = coercePressure(dot.Pressure)

			i++
		}
		x0 = x1
		y0 = y1
	}

	return Stroke{
		PointerType: lookupPointer(s.BrushType),
		PointerID:   1,
		X:           x[:i],
		Y:           y[:i],
		Timestamp:   ts[:i],
		Pressure:    p[:i],
	}, fromMillis(ms)
}

func coercePressure(p float32) float64 {
	return math.Max(0.0, math.Min(1.0, float64(p)))
}

func toMillis(t time.Time) int64 {
	nanos := t.UnixNano()
	return nanos / 1000000
}

func fromMillis(n int64) time.Time {
	secs := int64(n / 1000)
	nanos := (int64(n) - (secs * 1000)) * 1000000
	return time.Unix(secs, nanos)
}

func isTextStroke(bt lines.BrushType) bool {
	switch bt {
	case lines.Eraser,
		lines.EraseArea,
		lines.Highlighter,
		lines.HighlighterV5:
		return false
	default:
		return true
	}
}

func lookupPointer(bt lines.BrushType) PointerType {
	switch bt {
	case lines.Eraser, lines.EraseArea:
		return Eraser
	default:
		return Pen
	}
}
