package rescript

import (
	"fmt"
	"io"

	"github.com/akeil/rmtool"
)

// NewMarkdownComposer creates a new composer which generates output in markdown format.
func NewMarkdownComposer() ComposeFunc {
	return composeMarkdown
}

type stringWriter struct {
	io.Writer
}

func (sw stringWriter) WriteString(s string) (int, error) {
	return sw.Write([]byte(s))
}

func composeMarkdown(w io.Writer, doc *rmtool.Document, r map[string][]*Token) error {
	var err error
	sw := stringWriter{w}

	// TODO: we might write a yaml frontmatter here
	sw.WriteString(fmt.Sprintf("# %v\n\n", doc.Name()))

	for i, pageID := range doc.Pages() {
		res, ok := r[pageID]
		if ok {
			err = markdownPage(sw, i, res)
			if err != nil {
				return err
			}
			// TODO: not after the last page
			sw.WriteString("\n\n---\n\n") // thematic break after each page
		}
		// TODO what should we do with pages w/o results?
	}

	// end the document with a newline
	_, err = sw.WriteString("\n")
	if err != nil {
		return err
	}

	return nil
}

// "Improve" the result
// recognize lists:
// lines starting with "-" or "*"
// in some cases, add missing space
// add a newline before the *first* and after the *last* line
// of consecutive list entries

func markdownPage(w io.StringWriter, idx int, tokens []*Token) error {
	var err error

	w.WriteString(fmt.Sprintf("**Page %d**\n\n", idx+1))

	for _, t := range tokens {
		// TODO: we might attempt to "guess" markdown here,

		_, err = w.WriteString(t.String())
		if err != nil {
			return err
		}
	}

	return nil
}
