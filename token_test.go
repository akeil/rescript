package rescript

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenIs(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("Foo", NewToken("Foo").String())
	assert.Equal("", NewToken("").String())

	assert.True(NewToken(" ").IsWhitespace())
	assert.True(NewToken("\n").IsWhitespace())
	assert.False(NewToken("").IsWhitespace())
	assert.False(NewToken("Foo").IsWhitespace())

	assert.True(NewToken("\n").IsWhitespace())
	assert.True(NewToken("\n").IsNewline())
	assert.False(NewToken(" ").IsNewline())

	assert.True(NewToken(".").IsPunctuation())
	assert.True(NewToken(":").IsPunctuation())
	assert.True(NewToken(",").IsPunctuation())
	assert.True(NewToken("-").IsPunctuation())
	assert.True(NewToken("_").IsPunctuation())
	assert.False(NewToken("").IsPunctuation())
	assert.False(NewToken("Foo").IsPunctuation())

	assert.True(NewToken("Foo").StartsUpper())
	assert.False(NewToken("bar").StartsUpper())

	assert.True(NewToken("Foo").IsWord())
	assert.True(NewToken("a").IsWord())
	assert.False(NewToken("@").IsWord())
	assert.False(NewToken(" ").IsWord())

	// two words separated by space should not occur (= three separate tpkens)
	assert.False(NewToken("foo bar").IsWord())

	assert.True(NewToken("-").IsDash())
	assert.False(NewToken("_").IsDash())
}
