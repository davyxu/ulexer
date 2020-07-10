package ulexer

func NextToken(lex *Lexer, m Matcher) *Token {

	lex.Read(WhiteSpace())

	return lex.Expect(m)
}
