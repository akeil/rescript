package rescript

// Node is an element in a doubly linked list of Tokens.
type Node struct {
	prev *Node
	next *Node
	data *Token
}

// NewNode creates a new Node wrapping the given Token.
func NewNode(t *Token) *Node {
	return &Node{
		data: t,
	}
}

// BuildLinkedList creates a doubly linked list from the given list of tokens.
// Returns the TAIL (first entry) of the list.
func BuildLinkedList(t []*Token) *Node {
	var head *Node
	var tail *Node
	for _, tt := range t {
		n := NewNode(tt)
		if head != nil {
			head.InsertAfter(n)
			head = head.Next()
		} else {
			head = n
			tail = n
		}
	}

	return tail
}

// Token returns the Token for this node.
func (n *Node) Token() *Token {
	return n.data
}

// Next returns the next node or nil if this is the HEAD.
func (n *Node) Next() *Node {
	return n.next
}

// Prev returns the revious node or nil if this is the TAIL.
func (n *Node) Prev() *Node {
	return n.prev
}

// IsHead tells if this is the last node in the list.
func (n *Node) IsHead() bool {
	return n.next == nil
}

// IsTail tells if this is the first node in the list.
func (n *Node) IsTail() bool {
	return n.prev == nil
}

// Remove drops this token from the list
// and directly links the previous and next nodes.
func (n *Node) Remove() {
	if n.prev != nil {
		n.prev.next = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	}

	n.prev = nil
	n.next = nil
}

// InsertAfter adds the given not after this one
func (n *Node) InsertAfter(o *Node) {
	next := n.Next()
	n.next = o
	o.prev = n
	o.next = next

	if next != nil {
		next.prev = o
	}
}

// InsertBefore inserts the given Node before this one
func (n *Node) InsertBefore(o *Node) {
	//
	// prev <-- n --> next
	//
	// prev <-- o <-- n --> next
	prev := n.Prev()
	n.prev = o
	o.next = n
	o.prev = prev

	if prev != nil {
		prev.next = o
	}
}

// Update replaces the Token payload for this node with another token.
func (n *Node) Update(t *Token) {
	n.data = t
}
