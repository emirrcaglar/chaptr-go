// SPDX-License-Identifier: GPL-3.0-or-later
// Go port and modifications of chaptr 1.5.0, 2026-10-27, by emirrcaglar.
// Upstream: https://github.com/johnthreekay/chaptr. See NOTICE for attribution.

package chaptr

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type TokenKind uint8

const (
	Word TokenKind = iota
	Delimiter
	Bracketed
	Parenthesized
	Curly
)

// Token is a slice of the original input. Start and End are byte offsets;
// bracket tokens include their markers in the offsets but omit them in Text.
// Delimiters hold a single Unicode character in Text.
type Token struct {
	Kind  TokenKind `json:"kind"`
	Text  string    `json:"text"`
	Start int       `json:"start"`
	End   int       `json:"end"`
}

// Tokenize splits a filename into words, delimiters, and bracket spans.
// Unclosed brackets consume the rest of the input. Brackets do not nest.
func Tokenize(input string) []Token {
	tokens := make([]Token, 0, len(input)/3)
	for i := 0; i < len(input); {
		start := i
		r, size := utf8.DecodeRuneInString(input[i:])
		i += size
		kind := Delimiter
		closer := byte(0)
		switch r {
		case '[':
			kind, closer = Bracketed, ']'
		case '(':
			kind, closer = Parenthesized, ')'
		case '{':
			kind, closer = Curly, '}'
		}
		if closer != 0 {
			end := len(input)
			if p := strings.IndexByte(input[i:], closer); p >= 0 {
				end = i + p
			}
			text := input[i:end]
			i = end
			if i < len(input) {
				i++
			}
			tokens = append(tokens, Token{kind, text, start, i})
			continue
		}
		if isWordChar(r) {
			kind = Word
			for i < len(input) {
				c, n := utf8.DecodeRuneInString(input[i:])
				if !isWordChar(c) {
					break
				}
				i += n
			}
		}
		tokens = append(tokens, Token{kind, input[start:i], start, i})
	}
	return tokens
}

func isWordChar(r rune) bool {
	// Rust's is_alphanumeric includes Unicode's Other_Alphabetic property.
	return unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.Is(unicode.Other_Alphabetic, r) || r == '\'' || r == '!' || r == '?'
}
