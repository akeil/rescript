package rescript

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordsToTokens(t *testing.T) {
	assert := assert.New(t)

	r := Result{
		Words: []Word{
			Word{Label: "foo"},
			Word{Label: " "},
			Word{Label: "bar"},
			Word{Label: " "},
			Word{Label: "baz"},
		},
	}

	n := toTokens(r)

	assert.Equal("foo", n.Token().String())
	n = n.Next()
	assert.Equal(" ", n.Token().String())
	n = n.Next()
	assert.Equal("bar", n.Token().String())

	n = n.Next()
	assert.Equal(" ", n.Token().String())
	n = n.Next()
	assert.Equal("baz", n.Token().String())

	assert.True(n.IsHead())
}
