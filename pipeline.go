package rescript

type PipelineFunc func(t []*Token) []*Token

func BuildPipeline(p ...PipelineFunc) PipelineFunc {
	return func(t []*Token) []*Token {
		for _, f := range p {
			t = f(t)
		}
		return t
	}
}

func Dehyphenate(t []*Token) []*Token {
	return t
}
