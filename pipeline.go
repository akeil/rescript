package rescript

// PipelineFunc is a function that can be chained to process a set of tokens.
//
// The output is the modified set of tokens. A PipelineFunc may change, remove
// or insert tokens.
type PipelineFunc func(t []*Token) []*Token

// BuildPipeline combines several pipeline functions into one.
func BuildPipeline(p ...PipelineFunc) PipelineFunc {
	return func(t []*Token) []*Token {
		for _, f := range p {
			t = f(t)
		}
		return t
	}
}

// Dehyphenate merges words that are separated by a hyphen.
func Dehyphenate(t []*Token) []*Token {
	result := make([]*Token, len(t))
	count := 0
	drop := 0

	stage := 0
	for _, token := range t {
		result[count] = token
		count++
		switch stage {
		case 0:
			if token.IsWord() {
				stage = 1
				drop = 1
			}
		case 1:
			if token.IsDash() {
				stage = 2
				drop++
			} else {
				stage = 0
			}
		case 2:
			if token.IsWhitespace() {
				stage = 2 // same
				drop++
			} else if token.IsWord() {
				stage = 3
				drop++
			} else {
				stage = 0
			}
		}
		if stage == 3 {
			// lets go back and remove some items from the list
			s := ""
			for i := 0; i < drop; i++ {
				idx := count - drop + i
				pt := result[idx]
				if pt.IsWord() {
					s += pt.String()
				}
			}
			count -= drop
			result[count] = NewToken(s)
			count++
			stage = 0
			drop = 0
		}
	}
	return result[0:count]
}
