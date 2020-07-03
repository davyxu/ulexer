package ulexer

// 匹配行结尾  Mac OS 9 以及之前的系统的换行符是 CR，从 Mac OS X （后来改名为“OS X”）开始的换行符是 LF即‘\n'，和Unix/Linux统一了。
func LineEnd() Matcher {
	return (*lineEndMatcher)(nil)
}

type lineEndMatcher int

func (*lineEndMatcher) TokenType() string {
	return "LineEnd"
}

func (*lineEndMatcher) MatchRune(index int, r rune) bool {
	return r == '\r' || r == '\n'
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
