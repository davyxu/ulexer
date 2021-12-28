package ulexer

import "unicode"

// 匹配空白符, 回车符属于空白符
func WhiteSpace() Matcher {
	return (*whiteSpaceMatcher)(nil)
}

type whiteSpaceMatcher struct{}

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
		return nil
	}

	tk = lex.NewToken(count, self)

	lex.Consume(count)

	return
}

// 匹配行结尾  Mac OS 9 以及之前的系统的换行符是 CR，从 Mac OS X （后来改名为“OS X”）开始的换行符是 LF即‘\n'，和Unix/Linux统一了。
func LineEnd() Matcher {
	return (*lineEndMatcher)(nil)
}

type lineEndMatcher struct{}

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
		return nil
	}

	tk = lex.NewToken(count, self)

	lex.Consume(count)

	return
}

// 流结束, TokenString: EOF
func FileEnd() Matcher {
	return (*fileEndMatcher)(nil)
}

type fileEndMatcher struct{}

func (*fileEndMatcher) TokenType() string {
	return "FilEnd"
}

func (self *fileEndMatcher) Read(lex *Lexer) (tk *Token) {
	if lex.EOF() {
		return lex.NewTokenLiteral(0, self, "@EOF")
	}
	return nil
}
