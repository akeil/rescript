package rescript

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeadTail(t *testing.T) {
	assert := assert.New(t)

	start, middle, end := sampleList()

	assert.False(start.IsHead())
	assert.True(start.IsTail())

	assert.False(middle.IsHead())
	assert.False(middle.IsTail())

	assert.True(end.IsHead())
	assert.False(end.IsTail())
}

func TestRemove(t *testing.T) {
	assert := assert.New(t)

	start, middle, end := sampleList()

	middle.Remove()
	assert.Equal(start.Next(), end)
	assert.Equal(end.Prev(), start)
}

func TestUpdate(t *testing.T) {
	assert := assert.New(t)

	start, middle, end := sampleList()

	middle.Update(NewToken("replaced"))
	assert.Equal("replaced", start.Next().Token().String())
	assert.Equal("replaced", end.Prev().Token().String())
}

func TestInsert(t *testing.T) {
	assert := assert.New(t)

	start, middle, end := sampleList()

	a := NewNode(NewToken("a"))
	b := NewNode(NewToken("b"))

	assert.Equal("middle", start.Next().Token().String())

	middle.InsertBefore(a)
	assert.Equal("a", start.Next().Token().String())
	assert.Equal("start", a.Prev().Token().String())
	assert.Equal("middle", a.Next().Token().String())
	assert.Equal("a", middle.Prev().Token().String())

	middle.InsertAfter(b)
	assert.Equal("b", end.Prev().Token().String())
	assert.Equal("middle", b.Prev().Token().String())
	assert.Equal("end", b.Next().Token().String())
	assert.Equal("b", end.Prev().Token().String())
}

func TestBuildList(t *testing.T) {
	assert := assert.New(t)

	tokens := []*Token{
		NewToken("foo"),
		NewToken("bar"),
		NewToken("baz"),
	}
	start := BuildLinkedList(tokens)

	assert.Equal(start.Token().String(), "foo")
	assert.Equal(start.Next().Token().String(), "bar")
	assert.Equal(start.Next().Next().Token().String(), "baz")
}

func sampleList() (*Node, *Node, *Node) {
	start := NewNode(NewToken("start"))
	middle := NewNode(NewToken("middle"))
	end := NewNode(NewToken("end"))

	start.InsertAfter(middle)
	middle.InsertAfter(end)

	return start, middle, end
}
