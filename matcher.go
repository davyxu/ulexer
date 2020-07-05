package ulexer

import "unicode"

// 匹配正数, 负数, 浮点数
func Numeral() Matcher {
	return (*numeralMatcher)(nil)
}

type numeralMatcher int

func (*numeralMatcher) TokenType() string {
	return "Numeral"
}

func (self *numeralMatcher) Read(lex *Lexer) (tk *Token) {

	var count int

	var dot bool

	for {
		c := lex.Peek(count)

		switch {
		case unicode.IsDigit(c):

		case c == '-' && count == 0:

		case c == '.' && count > 0 && !dot:
			dot = true
		default:
			goto ExitFor
		}

		count++
	}

ExitFor:

	if count == 0 {
		return EmptyToken
	}

	tk = lex.NewToken(count, self)

	lex.Consume(count)

	return
}
