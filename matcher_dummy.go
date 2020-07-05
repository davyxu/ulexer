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

// 匹配行结尾  Mac OS 9 以及之前的系统的换行符是 CR，从 Mac OS X （后来改名为“OS X”）开始的换行符是 LF即‘\n'，和Unix/Linux统一了。
func LineEnd() Matcher {
	return (*lineEndMatcher)(nil)
}

type lineEndMatcher int

func (*lineEndMatcher) TokenType() string {
	return "LineEnd"
}

func (self *lineEndMatcher) Read(lex *Lexer) (tk *Token) {

	var count int
	for {
		c := lex.Peek(count)

		if c != '\n' && c != '\r' {
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
