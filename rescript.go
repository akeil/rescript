package rescript

import (
	"io"

	"github.com/akeil/rmtool"
)

type Composer interface {
	Compose(w io.Writer, doc *rmtool.Document, r map[string]Result) error
}
