package rescript

import (
	"unicode"
)

// Token represents a single text element that was reconized from the
// handwriting input. The fully reconized text consists of a list of tokens.
//
// RULES for tokenization:
//
// - consecutive whitespace is split into multiple tokens
// - punctuation is a single token
type Token struct {
	text  string
	runes []rune
}

// NewToken creates a new token with the given content.
func NewToken(s string) *Token {
	return &Token{s, []rune(s)}
}

func (t *Token) String() string {
	return t.text
}

func (t *Token) isSingle() bool {
	return len(t.runes) == 1
}

// The various IsXxx functions from Go's unicode package refer to unicode
// categories. See:
// https://en.wikipedia.org/wiki/Unicode_character_property

// IsWhitespace tells whether this token is a single whiotespace character.
// For a typical recognition result, this will be either a single space or
// newline.
func (t *Token) IsWhitespace() bool {
	if !t.isSingle() {
		return false
	}
	// TODO: how is this different from
	// https://golang.org/src/unicode/graphic.go?s=3997:4022#L116

	// copied from
	// https://github.com/reiver/go-whitespace/blob/master/whitespace.go
	switch t.runes[0] {
	case
		'\u0009', // horizontal tab
		'\u000A', // line feed
		'\u000B', // vertical tab
		'\u000C', // form feed
		'\u000D', // carriage return
		'\u0020', // space
		'\u0085', // next line
		'\u00A0', // no-break space
		'\u1680', // ogham space mark
		'\u180E', // mongolian vowel separator
		'\u2000', // en quad
		'\u2001', // em quad
		'\u2002', // en space
		'\u2003', // em space
		'\u2004', // three-per-em space
		'\u2005', // four-per-em space
		'\u2006', // six-per-em space
		'\u2007', // figure space
		'\u2008', // punctuation space
		'\u2009', // thin space
		'\u200A', // hair space
		'\u2028', // line separator
		'\u2029', // paragraph separator
		'\u202F', // narrow no-break space
		'\u205F', // medium mathematical space
		'\u3000': // ideographic space
		return true
	default:
		return false
	}
}

// IsNewline tells if this token is a newline.
func (t *Token) IsNewline() bool {
	if !t.isSingle() {
		return false
	}
	// copied from
	// https://github.com/reiver/go-whitespace/blob/master/mandatorybreak.go
	switch t.runes[0] {
	case
		'\u000A', // line feed
		'\u000B', // vertical tab
		'\u000C', // form feed
		'\u000D', // carriage return
		'\u0085', // next line
		'\u2028', // line separator
		'\u2029': // paragraph separator
		return true
	default:
		return false
	}
}

// StartsUpper tells if this token starts with an uppercase letter
func (t *Token) StartsUpper() bool {
	if len(t.runes) == 0 {
		return false
	}
	return unicode.IsUpper(t.runes[0])
}

// IsPunctuation tells if this token belongs to the punctuation group.
func (t *Token) IsPunctuation() bool {
	return t.isSingle() && unicode.IsPunct(t.runes[0])
}

// IsWord tells if this token is a "word".
//
// This is determined by looking at the recongnized characters. This may still
// be a poorly recognized word.
func (t *Token) IsWord() bool {
	// TODO: not sure is this holds
	// "words" only consist of letters - right?
	if len(t.runes) == 0 {
		return false
	}
	for _, r := range t.runes {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// IsDash tells if this token is a dash ("-"), including several unicode variants.
func (t *Token) IsDash() bool {
	if !t.isSingle() {
		return false
	}

	// source:
	// https://www.fileformat.info/info/unicode/category/Pd/list.htm
	switch t.runes[0] {
	case
		'\u002D', // hyphen-minus
		'\u058A', // armenian hyphen
		'\u05BE', // hebrew punctuation maqaf
		'\u1400', // canadian syllabics hyphen
		'\u2010', // hyphen
		'\u2011', // non-breaking hyphen
		'\u2012', // figure dash
		'\u2013', // en dsh
		'\u2014', // em dash
		'\u2015', // horizontal bar
		'\u2E3A', // two em dash
		'\u2E3B', // three em dash
		'\uFE58', // small em dash
		'\uFE63', // small hyphen minus
		'\uFE0D': // fullwidth hyphen minus
		return true
	default:
		return false
	}
}
