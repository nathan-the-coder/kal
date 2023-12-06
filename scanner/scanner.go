package scanner

import (
	"kal/error"
	"kal/token"
)


type Scanner struct {
	source  string
	tokens  []token.Token
	start   int
	current int
	line    int
}

func NewScanner(src string) *Scanner {
	var sc = &Scanner{
		source: src,
	}

	return sc
}

func (sc *Scanner) isAtEnd() bool { return sc.current >= len(sc.source) }

func (sc *Scanner) peek() byte {
	if sc.isAtEnd() {
		return byte(0)
	}
	return sc.source[sc.current]
}

func (sc *Scanner) peekNext() byte {
	if sc.current+1 >= len(sc.source) {
		return byte(0)
	}
	return sc.source[sc.current+1]
}

func (sc *Scanner) advance() byte {
	character := sc.source[sc.current]
	sc.current++
	return character
}

func (sc *Scanner) addToken(t token.TokenType) {
	sc.addTokenWithLiteral(t, nil)
}

func (sc *Scanner) addTokenWithLiteral(t token.TokenType, literal interface{}) {
	text := sc.source[sc.start:sc.current]
	newToken := token.Token{
		Type:    t,
		Lexeme:  text,
		Literal: literal,
		Line:    sc.line,
	}
	sc.tokens = append(sc.tokens, newToken)
}

func (sc *Scanner) match(expected byte) bool {
	if sc.isAtEnd() {
		return false
	}
	if sc.source[sc.current] != expected {
		return false
	}

	sc.current += 1
	return true
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isAlpha(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

func isAlNum(ch byte) bool {
	return isDigit(ch) || isAlpha(ch)
}

func (sc *Scanner) string() {
	for sc.peek() != '"' && !sc.isAtEnd() {
		if sc.peek() == '\n' {
			sc.line++
		}
		sc.advance()
	}
	if sc.isAtEnd() {
		error.Error(sc.line, "Unterminated string.")
		return
	}

	sc.advance()

	value := sc.source[sc.start+1 : sc.current-1]
	sc.addTokenWithLiteral(token.STRING, value)
}

func (sc *Scanner) identifier() {
	for isAlNum(sc.peek()) { sc.advance() }

	text := sc.source[sc.start:sc.current]
	tokenType := token.Keywords[text]
	if tokenType == "" {
		tokenType = token.IDENTIFIER
	}
	sc.addToken(tokenType)
}

func (sc *Scanner) number() {
	for isDigit(sc.peek()) {
		sc.advance()
	}

	if sc.peek() == '.' && isDigit(sc.peekNext()) {
		sc.advance()

		for isDigit(sc.peek()) {
			sc.advance()
		}
	}

	sc.addTokenWithLiteral(token.NUMBER, sc.source[sc.start:sc.current])
}

func (sc *Scanner) ScanToken() {
	var c = sc.advance()
	switch c {
	case '(':
		sc.addToken(token.LEFT_PAREN)
	case ')':
		sc.addToken(token.RIGHT_PAREN)
	case '{':
		sc.addToken(token.LEFT_BRACE)
	case '}':
		sc.addToken(token.RIGHT_BRACE)
	case ',':
		sc.addToken(token.COMMA)
	case '.':
		sc.addToken(token.DOT)
	case '-':
		sc.addToken(token.MINUS)
	case '+':
		sc.addToken(token.PLUS)
	case ';':
		sc.addToken(token.SEMICOLON)
	case '*':
		sc.addToken(token.STAR)
	case '!':
		if sc.match('=') {
			sc.addToken(token.BANG_EQUAL)
		} else {
			sc.addToken(token.BANG)
		}
	case '=':
		if sc.match('=') {
			sc.addToken(token.EQUAL_EQUAL)
		} else {
			sc.addToken(token.EQUAL)
		}
	case '<':
		if sc.match('=') {
			sc.addToken(token.LESS_EQUAL)
		} else {
			sc.addToken(token.LESS)
		}
	case '>':
		if sc.match('=') {
			sc.addToken(token.GREATER_EQUAL)
		} else {
			sc.addToken(token.GREATER)
		}
	case '/':
		if sc.match('/') {
			for sc.peek() != '\n' && !sc.isAtEnd() {
				sc.advance()
			}
		} else if sc.match('*') {
			for sc.peek() != '/' {
				sc.advance()
			}
		}

	case ' ':
	case '\r':
	case '\t':
		break

	case '\n':
		sc.line++
		break

	case '"':
		sc.string()

	default:
		if isDigit(c) {
			sc.number()
		} else if isAlpha(c) {
			sc.identifier()
		} else {
			error.Error(sc.line, "unexpected character")
		}
	}
}

func (sc *Scanner) ScanTokens() []token.Token {
	for !sc.isAtEnd() {
		sc.start = sc.current
		sc.ScanToken()
	}

	if !error.HadError {
		sc.tokens = append(sc.tokens, *token.NewToken(token.EOF, "", nil, sc.line))
	}
	return sc.tokens
}
