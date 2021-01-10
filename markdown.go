package rescript

import (
	"fmt"
	"io"

	"github.com/akeil/rmtool"
)

type markdownComposer struct {}

func NewMarkdownComposer() Composer {
	return &markdownComposer{}
}

func (m *markdownComposer) Compose(w io.Writer, doc *rmtool.Document, r map[string]Result) error {
	sw := stringWriter{w}
	return composeMarkdown(sw, doc, r)
}

type stringWriter struct {
	io.Writer
}

func (sw stringWriter) WriteString(s string) (int, error) {
	return sw.Write([]byte(s))
}

func composeMarkdown(w io.StringWriter, doc *rmtool.Document, r map[string]Result) error {
	var err error

	// TODO: we might write a yaml frontmatter here
	w.WriteString(fmt.Sprintf("# %v\n\n", doc.Name()))

	for i, pageID := range doc.Pages() {
		res, ok := r[pageID]
		if ok {
			err = page(w, i, res)
			if err != nil {
				return err
			}
			// TODO: not after the last page
			w.WriteString("\n\n---\n\n") // thematic break after each page
		}
		// TODO what should we do with pages w/o results?
	}

	// end the document with a newline
	_, err = w.WriteString("\n")
	if err != nil {
		return err
	}

	return nil
}

func page(w io.StringWriter, idx int, r Result) error {
	var err error

	w.WriteString(fmt.Sprintf("**Page %d**\n\n", idx+1))

	for _, wd := range r.Words {
		err = word(w, wd)
		if err != nil {
			return err
		}
	}

	return nil
}

func word(w io.StringWriter, wd Word) error {
	// TODO: we might attempt to "guess" markdown here,
	// e.g. prepend '#' before larger texts

	// reconize lists:
	// lines starting with "-" or "*"
	// in some cases, add missing space
	// add a newline vefoire the *first* and after the *last* line
	// of consecutive list entries

	// recongnize paragraphs
	// ic `word` is a newline, look at the vertical pixel distance
	// towards the next word - if higher than threshold, insert another newline
	s := wd.Label

	_, err := w.WriteString(s)
	if err != nil {
		return err
	}

	return nil
}
