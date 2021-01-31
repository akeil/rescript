package rescript

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlaintextPage(t *testing.T) {
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

	err := plaintextPage(&buf, 2, node)
	assert.Nil(err)

	s := string(buf.Bytes())
	expected := "\n[Page 3]\n\nfoo bar baz\nnewline\n"
	assert.Equal(expected, s)
}
