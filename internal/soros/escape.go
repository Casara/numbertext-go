package soros

import "strings"

// Soros special characters (spec 2.6) that must be quoted with a
// backslash to be used literally: \\, \", \$, \(, \), \[, \], \|, \#, \;
// and the \n newline notation. Each is mapped to a private-use-area rune
// while the template/pattern text is scanned for real syntax, then mapped
// back to the literal character once a piece of text is known to be plain
// output (see literal unescaping in template.go).
const (
	phBackslash = ''
	phQuote     = ''
	phDollar    = ''
	phLParen    = ''
	phRParen    = ''
	phLBracket  = ''
	phRBracket  = ''
	phPipe      = ''
	phHash      = ''
	phSemicolon = ''
	phNewline   = ''
)

var escapeToPlaceholder = map[byte]rune{
	'\\': phBackslash,
	'"':  phQuote,
	'$':  phDollar,
	'(':  phLParen,
	')':  phRParen,
	'[':  phLBracket,
	']':  phRBracket,
	'|':  phPipe,
	'#':  phHash,
	';':  phSemicolon,
	'n':  phNewline,
}

// placeholderToRegex maps a placeholder rune back to a RE2-safe literal
// representation, used when compiling a pattern (as opposed to rendering
// template output, see placeholderToLiteral).
var placeholderToRegex = map[rune]string{
	phBackslash: `\\`,
	phQuote:     `"`,
	phDollar:    `\$`,
	phLParen:    `\(`,
	phRParen:    `\)`,
	phLBracket:  `\[`,
	phRBracket:  `\]`,
	phPipe:      `\|`,
	phHash:      `#`,
	phSemicolon: `;`,
	phNewline:   "\n",
}

// revealEscapesForRegex maps placeholder runes back to their literal,
// regex-metacharacter-safe form so the result can be embedded in a
// pattern string passed to regexp.Compile.
func revealEscapesForRegex(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if lit, ok := placeholderToRegex[r]; ok {
			b.WriteString(lit)

			continue
		}
		b.WriteRune(r)
	}

	return b.String()
}

var placeholderToLiteral = map[rune]byte{
	phBackslash: '\\',
	phQuote:     '"',
	phDollar:    '$',
	phLParen:    '(',
	phRParen:    ')',
	phLBracket:  '[',
	phRBracket:  ']',
	phPipe:      '|',
	phHash:      '#',
	phSemicolon: ';',
	phNewline:   '\n',
}

// hideEscapes replaces every recognized "\X" escape sequence in s with a
// private-use placeholder rune, so the Soros syntax scanner (parseNodes,
// the pattern lexer) never mistakes an escaped character for real syntax.
// A backslash followed by an ASCII digit is left untouched: that's always
// a backreference ("\1".."\9" or "\0" for the whole match), never an
// escape sequence, and needs no special case here since a digit is never
// an escapeToPlaceholder key - it already falls through unmodified below.
func hideEscapes(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\\' && i+1 < len(s) {
			if ph, ok := escapeToPlaceholder[s[i+1]]; ok {
				b.WriteRune(ph)
				i++

				continue
			}
		}
		b.WriteByte(c)
	}

	return b.String()
}

// revealEscapes maps placeholder runes produced by hideEscapes back to
// their literal characters. It is applied only to text that is known to
// be plain output (a literalNode, or a compiled regex fragment) once
// syntax scanning is complete.
func revealEscapes(s string) string {
	if !strings.ContainsAny(s, string([]rune{
		phBackslash, phQuote, phDollar, phLParen, phRParen,
		phLBracket, phRBracket, phPipe, phHash, phSemicolon, phNewline,
	})) {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if c, ok := placeholderToLiteral[r]; ok {
			b.WriteByte(c)

			continue
		}
		b.WriteRune(r)
	}

	return b.String()
}
