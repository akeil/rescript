package rescript

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDehyphenate(t *testing.T) {
	assert := assert.New(t)
	var ta []*Token
	var tb []*Token

	// basic, unchanged
	ta = []*Token{
		NewToken("foo"),
		NewToken(" "),
		NewToken("bar"),
	}
	tb = Dehyphenate(ta)
	assert.Equal(len(ta), len(tb))

	// hyphenated word should be merged
	//
	// foo-
	// bar
	ta = []*Token{
		NewToken("foo"),
		NewToken("-"),
		NewToken("\n"),
		NewToken("bar"),
	}
	tb = Dehyphenate(ta)
	assert.Equal(1, len(tb), "Hyphenated words should be merged")

	// hyphenated word should be merged
	//
	// foo- bar
	ta = []*Token{
		NewToken("foo"),
		NewToken("-"),
		NewToken(" "),
		NewToken("bar"),
	}
	tb = Dehyphenate(ta)
	assert.Equal(1, len(tb), "Hyphenated words should be merged")
	assert.Equal("foobar", tb[0].String(), "Hyphenated words should be merged")

	// hyphenated word should be merged
	//
	// foo-bar
	ta = []*Token{
		NewToken("foo"),
		NewToken("-"),
		NewToken("bar"),
	}
	tb = Dehyphenate(ta)
	assert.Equal(1, len(tb), "Hyphenated words should be merged")
	assert.Equal("foobar", tb[0].String(), "Hyphenated words should be merged")

	// two words separated by dash should NOT be merged
	//
	// foo - bar
	ta = []*Token{
		NewToken("foo"),
		NewToken(" "),
		NewToken("-"),
		NewToken(" "),
		NewToken("bar"),
	}
	tb = Dehyphenate(ta)
	assert.Equal(len(ta), len(tb), "Words separated by dash should NOT be merged")

	// List elements are not hyphenated words
	//
	// foo
	// - item
	ta = []*Token{
		NewToken("foo"),
		NewToken("\n"),
		NewToken("-"),
		NewToken(" "),
		NewToken("item"),
	}
	tb = Dehyphenate(ta)
	assert.Equal(len(ta), len(tb), "List items should NOT be merged")

	// List elements are not hyphenated words (even if the list is not properly recognized)
	//
	// foo
	// -item
	ta = []*Token{
		NewToken("foo"),
		NewToken("\n"),
		NewToken("-"),
		NewToken("item"),
	}
	tb = Dehyphenate(ta)
	assert.Equal(len(ta), len(tb), "List items should NOT be merged")

}
