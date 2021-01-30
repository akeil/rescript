package rescript

import (
	"fmt"
	"io"
	"strings"

	"github.com/akeil/rmtool"
)

// NewPlaintextComposer creates a new composer which creates plain text output
// for a regicnition result.
func NewPlaintextComposer() ComposeFunc {
	return composePlain
}

func composePlain(w io.Writer, doc *rmtool.Document, r map[string][]*Token) error {
	var err error
	sw := stringWriter{w}

	// Output the title if we have one
	title := doc.Name()
	if title != "" {
		_, err = sw.WriteString(strings.ToUpper(title) + "\n")
		if err != nil {
			return err
		}
	}

	// Output the text body from all pages
	for i, pageID := range doc.Pages() {
		res, ok := r[pageID]
		if ok {
			err = plainPage(sw, i, res)
			if err != nil {
				return err
			}
		}
		// TODO what should we do with pages w/o results?
	}

	// End the document with a newline
	_, err = sw.WriteString("\n")
	if err != nil {
		return err
	}

	return nil
}

func plainPage(sw io.StringWriter, idx int, tokens []*Token) error {
	var err error

	sw.WriteString(fmt.Sprintf("\n[Page %d]\n\n", idx+1))

	for _, t := range tokens {
		_, err = sw.WriteString(t.String())
		if err != nil {
			return err
		}
	}

	_, err = sw.WriteString("\n")
	if err != nil {
		return err
	}

	return nil
}
