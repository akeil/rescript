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
	ta = BuildLinkedList([]*Token{
		NewToken("foo"),
		NewToken(" "),
		NewToken("bar"),
	})
	tb = Dehyphenate(ta)
	assert.Equal(str(ta), str(tb))
	assert.Equal(str(head(ta)), str(head(tb)))

	// hyphenated word should be merged
	//
	// foo-
	// bar
	ta = BuildLinkedList([]*Token{
		NewToken("foo"),
		NewToken("-"),
		NewToken("\n"),
		NewToken("bar"),
	})
	tb = Dehyphenate(ta)
	assert.Equal("foobar", str(tb), "Hyphenated words should be merged")

	// hyphenated word should be merged
	//
	// foo- bar
	ta = BuildLinkedList([]*Token{
		NewToken("foo"),
		NewToken("-"),
		NewToken(" "),
		NewToken("bar"),
	})
	tb = Dehyphenate(ta)
	assert.Equal("foobar", str(tb), "Hyphenated words should be merged")

	// subsequent hyphenations should also be merged
	//
	// foo- bar
	ta = BuildLinkedList([]*Token{
		NewToken("foo"),
		NewToken("-"),
		NewToken(" "),
		NewToken("bar"),
		NewToken(" "),
		NewToken("abc"),
		NewToken("-"),
		NewToken(" "),
		NewToken("def"),
	})
	tb = Dehyphenate(ta)
	assert.Equal("foobar", str(tb), "Hyphenated words should be merged")
	assert.Equal(" ", str(tb.Next()), "Hyphenated words should be merged")
	assert.Equal("abcdef", str(tb.Next().Next()), "Hyphenated words should be merged")

	// hyphenated word should be merged
	//
	// foo-bar
	ta = BuildLinkedList([]*Token{
		NewToken("foo"),
		NewToken("-"),
		NewToken("bar"),
	})
	tb = Dehyphenate(ta)
	assert.Equal("foobar", str(tb), "Hyphenated words should be merged")

	// two words separated by dash should NOT be merged
	//
	// foo - bar
	ta = BuildLinkedList([]*Token{
		NewToken("foo"),
		NewToken(" "),
		NewToken("-"),
		NewToken(" "),
		NewToken("bar"),
	})
	tb = Dehyphenate(ta)
	assert.Equal("foo", str(tb), "Words separated by dash should NOT be merged")
	assert.Equal("bar", str(head(tb)), "Words separated by dash should NOT be merged")

	// List elements are not hyphenated words
	//
	// foo
	// - item
	ta = BuildLinkedList([]*Token{
		NewToken("foo"),
		NewToken("\n"),
		NewToken("-"),
		NewToken(" "),
		NewToken("item"),
	})
	tb = Dehyphenate(ta)
	assert.Equal("foo", str(tb), "Words separated by dash should NOT be merged")
	assert.Equal("item", str(head(tb)), "Words separated by dash should NOT be merged")

	// List elements are not hyphenated words (even if the list is not properly recognized)
	//
	// foo
	// -item
	ta = BuildLinkedList([]*Token{
		NewToken("foo"),
		NewToken("\n"),
		NewToken("-"),
		NewToken("item"),
	})
	tb = Dehyphenate(ta)
	assert.Equal("foo", str(tb), "Words separated by dash should NOT be merged")
	assert.Equal("item", str(head(tb)), "Words separated by dash should NOT be merged")
}
