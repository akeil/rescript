package rescript

import (
	"io"

	"github.com/akeil/rmtool"
)

// ComposeFunc is a function that generates an output document from the given
// set of tokens. THe result is written to the given writer.
type ComposeFunc func(w io.Writer, doc *rmtool.Document, r map[string]*Node) error
