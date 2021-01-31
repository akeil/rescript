package rescript

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComposePlaintext(t *testing.T) {
	assert := assert.New(t)

	var buf bytes.Buffer

	append := func(n *Node, s string) {
		head := n
		for node := n; node != nil; node = node.Next() {
			head = node
		}
		head.InsertAfter(NewNode(NewToken(s)))
	}

	node := NewNode(NewToken("foo"))
	append(node, " ")
	append(node, "bar")
	append(node, " ")
	append(node, "baz")
	append(node, "\n")
	append(node, "newline")

	node1 := NewNode(NewToken("second page"))

	m := Metadata{
		Title:   "My Title",
		PageIDs: []string{"page0", "page1"},
	}

	nodes := map[string]*Node{
		"page0": node,
		"page1": node1,
	}

	c := NewPlaintextComposer()
	err := c(&buf, m, nodes)
	assert.Nil(err)

	s := string(buf.Bytes())
	expected := "MY TITLE\n\n[Page 1]\n\nfoo bar baz\nnewline\n\n[Page 2]\n\nsecond page\n"
	assert.Equal(expected, s)
}

func TestPlaintextError(t *testing.T) {
	assert := assert.New(t)

	node := NewNode(NewToken("foo"))
	w := failWriter{}

	err := plaintextPage(w, 2, node)
	assert.Error(err)
}

type failWriter struct{}

func (f failWriter) WriteString(s string) (int, error) {
	return 0, errors.New("test failure")
}
