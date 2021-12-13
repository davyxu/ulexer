package ulexer

// C的行注释
func CLineComment() Matcher {
	return (*cLineCommentMatcher)(nil)
}

func (*cLineCommentMatcher) TokenType() string {
	return "CLineComment"
}

type cLineCommentMatcher struct{}

func (self *cLineCommentMatcher) Read(lex *Lexer) (tk *Token) {

	if lex.Peek(0) != '/' || lex.Peek(1) != '/' {
		return nil
	}

	lex.Consume(2)

	var count int
	for {
		c := lex.Peek(count)

		if c == '\r' || c == '\n' {
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

// Unix的行注释
func UnixLineComment() Matcher {
	return (*unixLineCommentMatcher)(nil)
}

func (*unixLineCommentMatcher) TokenType() string {
	return "UnixLineComment"
}

type unixLineCommentMatcher struct{}

func (self *unixLineCommentMatcher) Read(lex *Lexer) (tk *Token) {

	if lex.Peek(0) != '#' {
		return nil
	}

	lex.Consume(2)

	var count int
	for {
		c := lex.Peek(count)

		if c == '\r' || c == '\n' {
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
