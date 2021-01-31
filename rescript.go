package rescript

import (
	"io"
)

// Metadata holds information about a document.
type Metadata struct {
	Title   string
	PageIDs []string
}

// ComposeFunc is a function that generates an output document from the given
// set of tokens. THe result is written to the given writer.
type ComposeFunc func(w io.Writer, m Metadata, r map[string]*Node) error
