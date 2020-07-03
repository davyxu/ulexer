package ulexer

// C的行注释
func CLineComment() Matcher {
	return (*cLineCommentMatcher)(nil)
}

func (*cLineCommentMatcher) TokenType() string {
	return "CLineComment"
}

type cLineCommentMatcher int

func (*cLineCommentMatcher) MatchRune(index int, r rune) bool {
	switch index {
	case 0:
		return r == '/'
	case 1:
		return r == '/'
	default:
		return r != '\r' && r != '\n'
	}
}

func (self *cLineCommentMatcher) Read(lex *Lexer) (tk *Token) {

	if lex.Peek(0) != '/' || lex.Peek(1) != '/' {
		return EmptyToken
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
		return EmptyToken
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

type unixLineCommentMatcher int

func (*unixLineCommentMatcher) MatchRune(index int, r rune) bool {
	switch index {
	case 0:
		return r == '#'
	default:
		return r != '\r' && r != '\n'
	}
}

func (self *unixLineCommentMatcher) Read(lex *Lexer) (tk *Token) {

	if lex.Peek(0) != '#' {
		return EmptyToken
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
		return EmptyToken
	}

	tk = lex.NewToken(count, self)

	lex.Consume(count)

	return
}
