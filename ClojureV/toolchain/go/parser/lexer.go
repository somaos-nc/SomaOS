package parser

import (
	"strings"
	"unicode"
)

type TokenType string

const (
	TokenLParen TokenType = "LParen"
	TokenRParen TokenType = "RParen"
	TokenLBrack TokenType = "LBracket"
	TokenRBrack TokenType = "RBracket"
	TokenIdent  TokenType = "Identifier"
	TokenNum    TokenType = "Number"
	TokenStr    TokenType = "String"
	TokenEOF    TokenType = "EOF"
)

type Token struct {
	Type  TokenType
	Value string
	Line  int
	Col   int
}

type Lexer struct {
	input string
	pos   int
	line  int
	col   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: input, line: 1, col: 1}
}

func (l *Lexer) current() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) advance() {
	if l.current() == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	l.pos++
}

func (l *Lexer) skipWhitespace() {
	for {
		c := l.current()
		if c == ';' { // Skip comments
			for l.current() != '\n' && l.current() != 0 {
				l.advance()
			}
		} else if unicode.IsSpace(rune(c)) || c == ',' { // Treat comma as whitespace in Clojure
			l.advance()
		} else {
			break
		}
	}
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	c := l.current()
	if c == 0 {
		return Token{Type: TokenEOF, Value: "", Line: l.line, Col: l.col}
	}

	startLine, startCol := l.line, l.col

	switch c {
	case '(':
		l.advance()
		return Token{Type: TokenLParen, Value: "(", Line: startLine, Col: startCol}
	case ')':
		l.advance()
		return Token{Type: TokenRParen, Value: ")", Line: startLine, Col: startCol}
	case '[':
		l.advance()
		return Token{Type: TokenLBrack, Value: "[", Line: startLine, Col: startCol}
	case ']':
		l.advance()
		return Token{Type: TokenRBrack, Value: "]", Line: startLine, Col: startCol}
	case '"':
		return l.readString()
	default:
		if isDigit(c) || (c == '-' && isDigit(l.peek())) {
			return l.readNumber()
		}
		return l.readIdentifier()
	}
}

func (l *Lexer) peek() byte {
	if l.pos+1 >= len(l.input) {
		return 0
	}
	return l.input[l.pos+1]
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isIdentChar(c byte) bool {
	if unicode.IsSpace(rune(c)) || c == '(' || c == ')' || c == '[' || c == ']' || c == '"' || c == ';' || c == 0 {
		return false
	}
	return true
}

func (l *Lexer) readString() Token {
	startLine, startCol := l.line, l.col
	l.advance() // Skip opening quote
	var sb strings.Builder
	for l.current() != '"' && l.current() != 0 {
		sb.WriteByte(l.current())
		l.advance()
	}
	if l.current() == '"' {
		l.advance()
	}
	return Token{Type: TokenStr, Value: sb.String(), Line: startLine, Col: startCol}
}

func (l *Lexer) readNumber() Token {
	startLine, startCol := l.line, l.col
	var sb strings.Builder
	if l.current() == '-' {
		sb.WriteByte(l.current())
		l.advance()
	}
	if l.current() == '0' && (l.peek() == 'x' || l.peek() == 'X') {
		sb.WriteByte(l.current())
		l.advance()
		sb.WriteByte(l.current())
		l.advance()
		for (l.current() >= '0' && l.current() <= '9') || (l.current() >= 'a' && l.current() <= 'f') || (l.current() >= 'A' && l.current() <= 'F') {
			sb.WriteByte(l.current())
			l.advance()
		}
	} else {
		for isDigit(l.current()) || l.current() == '.' {
			sb.WriteByte(l.current())
			l.advance()
		}
	}
	return Token{Type: TokenNum, Value: sb.String(), Line: startLine, Col: startCol}
}

func (l *Lexer) readIdentifier() Token {
	startLine, startCol := l.line, l.col
	var sb strings.Builder
	for isIdentChar(l.current()) {
		sb.WriteByte(l.current())
		l.advance()
	}
	return Token{Type: TokenIdent, Value: sb.String(), Line: startLine, Col: startCol}
}
