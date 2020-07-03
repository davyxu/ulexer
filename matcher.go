package ulexer

import "unicode"

// 匹配空白符, 回车符属于空白符
func WhiteSpace() Matcher {
	return (*whiteSpaceMatcher)(nil)
}

type whiteSpaceMatcher int

func (*whiteSpaceMatcher) TokenType() string {
	return "WhiteSpace"
}

func (*whiteSpaceMatcher) MatchRune(index int, r rune) bool {
	return unicode.In(r, unicode.White_Space)
}

func (self *whiteSpaceMatcher) Read(lex *Lexer) (tk *Token) {

	var count int
	for {
		c := lex.Peek(count)

		if !unicode.In(c, unicode.White_Space) {
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

// 匹配正数, 负数, 浮点数
func Numeral() Matcher {
	return (*numeralMatcher)(nil)
}

type numeralMatcher int

func (*numeralMatcher) TokenType() string {
	return "Numeral"
}

func (*numeralMatcher) MatchRune(index int, r rune) bool {
	return unicode.IsDigit(r)
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

// 匹配标识符
func Identifier() Matcher {
	return (*identifierMatcher)(nil)
}

type identifierMatcher int

func (*identifierMatcher) TokenType() string {
	return "Identifier"
}

func (*identifierMatcher) MatchRune(index int, r rune) bool {
	if index == 0 {
		if unicode.IsLetter(r) || r == '_' {
			return true
		}

	} else {

		if unicode.IsLetter(r) || r == '_' || unicode.IsDigit(r) {
			return true
		}
	}

	return false
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

	self := &literalMatcher{}

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

type literalMatcher struct {
	literal []rune
}

func (*literalMatcher) TokenType() string {
	return "Literal"
}

func (self *literalMatcher) Read(lex *Lexer) (tk *Token) {

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
