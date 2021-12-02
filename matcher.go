package ulexer

type Matcher interface {
	Read(lex *Lexer) *Token
	TokenType() string
}

// 匹配true/false
func Bool() Matcher {
	return (*boolMatcher)(nil)
}

type boolMatcher int

func (*boolMatcher) TokenType() string {
	return "Bool"
}

func (self *boolMatcher) Read(lex *Lexer) (tk *Token) {

	tk = Select(lex, Contain("true"), Contain("false"))

	if tk != EmptyToken {
		tk.t = self.TokenType()
	}

	return
}
