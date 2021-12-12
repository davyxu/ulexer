package ulexer

import (
	"strings"
	"unicode"
)

// 数字, 无符号整形
func UInteger() Matcher {
	return (*uIntegerMatcher)(nil)
}

type uIntegerMatcher struct{}

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

type integerMatcher struct{}

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

type numeralMatcher struct{}

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

// 匹配十六进制, 带0x前缀的, 自动去掉, 使用Token.Numeral(32, 16, false)获取数值
func Hex() Matcher {
	return (*hexMatcher)(nil)
}

type hexMatcher struct{}

func (*hexMatcher) TokenType() string {
	return "Hex"
}

func isHexNumber(c rune) bool {
	if unicode.IsNumber(c) {
		return true
	}
	switch c {
	case 'a', 'b', 'c', 'd', 'e', 'f':
		return true
	}

	return false
}

func (self *hexMatcher) Read(lex *Lexer) (tk *Token) {

	var count int

	for {
		c := lex.Peek(count)

		switch {
		case isHexNumber(c):
		case c == 'x' && count == 1:
		default:
			goto ExitFor
		}

		count++
	}

ExitFor:

	if count == 0 {
		return EmptyToken
	}

	str := lex.ToLiteral(count)

	// 纯十六进制parse时, 在输出字符前加0x, 方便直接转换为16进制
	if strings.HasPrefix(str, "0x") {
		str = str[2:]
	}

	tk = lex.NewTokenLiteral(count, self, str)

	lex.Consume(count)

	return
}
