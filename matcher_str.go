package ulexer

import (
	"strings"
	"unicode"
)

// 匹配标识符
func Identifier() Matcher {
	return (*identifierMatcher)(nil)
}

type identifierMatcher int

func (*identifierMatcher) TokenType() string {
	return "Identifier"
}

func (self *identifierMatcher) Read(lex *Lexer) (tk *Token) {

	var count int
	for {
		c := lex.Peek(count)

		isBasic := unicode.IsLetter(c) || c == '_'

		switch {
		case count == 0 && isBasic:
		case count > 0 && (isBasic || unicode.IsDigit(c)):
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

// 包含字面量
func Contain(literal interface{}) Matcher {

	self := &containMatcher{}

	switch v := literal.(type) {
	case string:
		self.literal = []rune(v)
	case rune:
		self.literal = []rune{v}
	default:
		panic("invalid contain")
	}

	return self
}

type containMatcher struct {
	literal []rune
}

func (*containMatcher) TokenType() string {
	return "Contain"
}

func (c *containMatcher) String() string {
	return string(c.literal)
}

func (self *containMatcher) Read(lex *Lexer) (tk *Token) {

	var count int
	for {
		c := lex.Peek(count)

		if count >= len(self.literal) {
			break
		}

		if c != self.literal[count] {
			break
		}

		count++
	}

	if count == 0 {
		return EmptyToken
	}

	tk = lex.NewToken(count, self)

	lex.Consume(count)

	return
}

// 匹配字符串
func String() Matcher {
	return (*stringMatcher)(nil)
}

type stringMatcher int

func (*stringMatcher) TokenType() string {
	return "String"
}

func (self *stringMatcher) Read(lex *Lexer) (tk *Token) {

	beginChar := lex.Peek(0)
	if beginChar != '"' && beginChar != '\'' {
		return EmptyToken
	}

	state := lex.State

	lex.Consume(1)

	var (
		escaping bool
		closed   bool
		sb       strings.Builder
	)

	var count int
	for {
		c := lex.Peek(count)

		if escaping {
			switch c {
			case 'n':
				sb.WriteRune('\n')
			case 'r':
				sb.WriteRune('\r')
			case '"', '\'':
				sb.WriteRune(c)
			default:
				sb.WriteRune('\\')
				sb.WriteRune(c)
			}

			escaping = false
		} else if c != beginChar {
			if c == '\\' {
				escaping = true
			} else if c != 0 {
				sb.WriteRune(c)
			}
		} else {
			closed = true
			break
		}

		if c == '\n' || c == 0 {
			break
		}

		count++
	}

	if !closed {
		lex.State = state
		return EmptyToken
	}

	end := count + 1

	tk = lex.NewTokenLiteral(end, self, sb.String())

	lex.Consume(end)

	return
}
