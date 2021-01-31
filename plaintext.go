package rescript

import (
	"fmt"
	"io"
	"strings"
)

// NewPlaintextComposer creates a new composer which creates plain text output
// for a regicnition result.
func NewPlaintextComposer() ComposeFunc {
	return composePlain
}

func composePlain(w io.Writer, m Metadata, r map[string]*Node) error {
	var err error
	sw := stringWriter{w}

	// Output the title if we have one
	title := m.Title
	if title != "" {
		_, err = sw.WriteString(strings.ToUpper(title) + "\n")
		if err != nil {
			return err
		}
	}

	// Output the text body from all pages
	for i, pageID := range m.PageIDs {
		tail, ok := r[pageID]
		if ok {
			err = plaintextPage(sw, i, tail)
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

func plaintextPage(sw io.StringWriter, idx int, n *Node) error {
	var err error

	_, err = sw.WriteString(fmt.Sprintf("\n[Page %d]\n\n", idx+1))
	if err != nil {
		return err
	}

	for node := n; node != nil; node = node.Next() {
		_, err = sw.WriteString(node.Token().String())
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
