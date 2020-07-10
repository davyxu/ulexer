package ulexer

import "unicode"

// 数字, 无符号整形
func UInteger() Matcher {
	return (*uIntegerMatcher)(nil)
}

type uIntegerMatcher int

func (*uIntegerMatcher) TokenType() string {
	return "UInteger"
}

func (self *uIntegerMatcher) Read(lex *Lexer) (tk *Token) {
	var count int

	for {
		c := lex.Peek(count)

		switch {
		case unicode.IsDigit(c):
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

// 有符号整数=正数, 0, 负数
func Integer() Matcher {
	return (*integerMatcher)(nil)
}

type integerMatcher int

func (*integerMatcher) TokenType() string {
	return "Integer"
}

func (self *integerMatcher) Read(lex *Lexer) (tk *Token) {

	var count int

	for {
		c := lex.Peek(count)

		switch {
		case unicode.IsDigit(c):
		case count == 0 && c == '-':
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

	var isFloat bool

	for {
		c := lex.Peek(count)

		switch {
		case unicode.IsDigit(c):

		case c == '-' && count == 0:

		case c == '.' && count > 0 && !isFloat:
			isFloat = true
		case isFloat && (c == 'e' || c == '+'):

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
