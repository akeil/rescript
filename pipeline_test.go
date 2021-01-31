package rescript

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDehyphenate(t *testing.T) {
	assert := assert.New(t)
	var ta *Node
	var tb *Node

	head := func(n *Node) *Node {
		nn := n
		for node := n; node != nil; node = node.Next() {
			nn = node
		}
		return nn
	}

	str := func(n *Node) string {
		return n.Token().String()
	}

	// basic, unchanged
	ta = buildSampleList("foo", " ", "bar")
	tb = Dehyphenate(ta)
	assert.Equal(str(ta), str(tb))
	assert.Equal(str(head(ta)), str(head(tb)))

	// hyphenated word should be merged
	//
	// foo-
	// bar
	ta = buildSampleList("foo", "-", "\n", "bar")
	tb = Dehyphenate(ta)
	assert.Equal("foobar", str(tb), "Hyphenated words should be merged")

	// hyphenated word should be merged
	//
	// foo- bar
	ta = buildSampleList("foo", "-", " ", "bar")
	tb = Dehyphenate(ta)
	assert.Equal("foobar", str(tb), "Hyphenated words should be merged")

	// subsequent hyphenations should also be merged
	//
	// foo- bar
	ta = buildSampleList("foo", "-", " ", "bar", " ", "abc", "-", " ", "def")
	tb = Dehyphenate(ta)
	assert.Equal("foobar", str(tb), "Hyphenated words should be merged")
	assert.Equal(" ", str(tb.Next()), "Hyphenated words should be merged")
	assert.Equal("abcdef", str(tb.Next().Next()), "Hyphenated words should be merged")

	// hyphenated word should be merged
	//
	// foo-bar
	ta = buildSampleList("foo", "-", "bar")
	tb = Dehyphenate(ta)
	assert.Equal("foobar", str(tb), "Hyphenated words should be merged")

	// two words separated by dash should NOT be merged
	//
	// foo - bar
	ta = buildSampleList("foo", " ", "-", " ", "bar")
	tb = Dehyphenate(ta)
	assert.Equal("foo", str(tb), "Words separated by dash should NOT be merged")
	assert.Equal("bar", str(head(tb)), "Words separated by dash should NOT be merged")

	// List elements are not hyphenated words
	//
	// foo
	// - item
	ta = buildSampleList("foo", "\n", "-", " ", "item")
	tb = Dehyphenate(ta)
	assert.Equal("foo", str(tb), "Words separated by dash should NOT be merged")
	assert.Equal("item", str(head(tb)), "Words separated by dash should NOT be merged")

	// List elements are not hyphenated words (even if the list is not properly recognized)
	//
	// foo
	// -item
	ta = buildSampleList("foo", "\n", "-", "item")
	tb = Dehyphenate(ta)
	assert.Equal("foo", str(tb), "Words separated by dash should NOT be merged")
	assert.Equal("item", str(head(tb)), "Words separated by dash should NOT be merged")
}

func buildSampleList(s ...string) *Node {
	var tail *Node
	var head *Node
	for _, w := range s {
		n := NewNode(NewToken(w))
		if head != nil {
			head.InsertAfter(n)
		} else {
			head = n
			tail = n
		}
	}

	return tail
}
