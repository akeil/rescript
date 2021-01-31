package rescript

// PipelineFunc is a function that can be chained to process a set of tokens.
//
// The output is the modified set of tokens. A PipelineFunc may change, remove
// or insert tokens.
type PipelineFunc func(n *Node) *Node

// BuildPipeline combines several pipeline functions into one.
func BuildPipeline(p ...PipelineFunc) PipelineFunc {
	return func(n *Node) *Node {
		for _, f := range p {
			n = f(n)
		}
		return n
	}
}

// Dehyphenate merges words that are separated by a hyphen.
func Dehyphenate(n *Node) *Node {
	count := 0
	state := 0
	var t *Token

	for node := n; node != nil; node = node.Next() {
		t = node.Token()
		switch state {
		case 0:
			if t.IsWord() {
				state = 1
				count = 1
			}
		case 1:
			if t.IsDash() {
				state = 2
				count++
			} else {
				state = 0
				count = 0
			}
		case 2:
			if t.IsWhitespace() {
				state = 2 // same
				count++
			} else if t.IsWord() {
				state = 3
			} else {
				state = 0
				count = 0
			}
		}

		// We have found a hyphenated word if we have reached state = 3
		//
		// Current node is the last part of the word
		// We need to merge `count` preceeding nodes
		//
		// We update the start node (not the current)
		// because the start node might be the return value of this function
		if state == 3 {
			// go back to the start of the hyphenated word
			start := node.Behind(count)
			// this will become the merged word
			s := start.Token().String()

			// drop `count` following nodes
			for i := 0; i < count; i++ {
				next := start.Next()
				if next.Token().IsWord() {
					s += next.Token().String()
				}
				next.Remove()
			}

			// make the merged word part of the list
			start.Update(NewToken(s))

			// "fix" the iterator - we have dropped the current node, reset it
			node = start
			// reset state
			count = 0
			state = 0
		}
	}

	return n
}
