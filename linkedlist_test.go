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

func TestAhead(t *testing.T) {
	assert := assert.New(t)

	start, middle, end := sampleList()

	assert.Equal(start.Ahead(1).Token().String(), middle.Token().String())
	assert.Equal(start.Ahead(2).Token().String(), end.Token().String())
	assert.Nil(start.Ahead(3))
	assert.Nil(start.Ahead(4))
	assert.Nil(start.Behind(1))
}

func TestBehind(t *testing.T) {
	assert := assert.New(t)

	start, middle, end := sampleList()
	assert.Equal(end.Behind(1).Token().String(), middle.Token().String())
	assert.Equal(end.Behind(2).Token().String(), start.Token().String())
	assert.Nil(end.Behind(3))
	assert.Nil(end.Behind(4))
	assert.Nil(end.Ahead(1))
}

func sampleList() (*Node, *Node, *Node) {
	start := NewNode(NewToken("start"))
	middle := NewNode(NewToken("middle"))
	end := NewNode(NewToken("end"))

	start.InsertAfter(middle)
	middle.InsertAfter(end)

	return start, middle, end
}
