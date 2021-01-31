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

func composeMarkdown(w io.Writer, doc *rmtool.Document, r map[string]*Node) error {
	var err error
	sw := stringWriter{w}

	// TODO: we might write a yaml frontmatter here
	sw.WriteString(fmt.Sprintf("# %v\n\n", doc.Name()))

	for i, pageID := range doc.Pages() {
		tail, ok := r[pageID]
		if ok {
			err = markdownPage(sw, i, tail)
			if err != nil {
				return err
			}
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

func markdownPage(sw io.StringWriter, idx int, n *Node) error {
	var err error

	_, err = sw.WriteString(fmt.Sprintf("**Page %d**\n\n", idx+1))
	if err != nil {
		return err
	}

	for node := n; node != nil; node = node.Next() {
		// TODO: we might attempt to "guess" markdown here,
		_, err = sw.WriteString(node.Token().String())
		if err != nil {
			return err
		}
	}

	// thematic break after each page
	// TODO: not after the last page
	_, err = sw.WriteString("\n\n---\n\n")
	if err != nil {
		return err
	}

	return nil
}
