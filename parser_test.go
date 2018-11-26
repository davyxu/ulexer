package golexer2

import "testing"

func TestFull(t *testing.T) {
	lex := NewLexer([]rune("  struct 2 enum 3"))

	for !lex.EOF() {

		SkipWhiteSpace(lex)

		if TryString(lex, "enum ") {
			SkipWhiteSpace(lex)
			t.Log(ExpectInt32(lex))
		}

		if TryString(lex, "struct") {
			SkipWhiteSpace(lex)
			t.Log(ExpectInt32(lex))
		}

	}
}

func TestExpect(t *testing.T) {
	lex := NewLexer([]rune("b"))

	lex.Run(func(lex *Lexer) {
		ExpectString(lex, "a")
		ExpectLogResult("Expect 'a'", t)
	})
}

// 回车处理
func TestLineEnd(t *testing.T) {
	// Mac OS 9 以及之前的系统的换行符是 CR，从 Mac OS X （后来改名为“OS X”）开始的换行符是 LF即‘\n'，和Unix/Linux统一了。
	lex := NewLexer([]rune("\n \r \r\n"))

	lex.Run(func(lex *Lexer) {
		for !lex.EOF() {
			SkipString(lex, []rune{' ', '\t'})
			SkipLineEnd(lex)
		}

		if lex.Line() != 3 {
			t.FailNow()
		}
	})
}

func TestSvcID(t *testing.T) {
	lex := NewLexer([]rune("game#1@dev"))

	lex.Run(func(lex *Lexer) {

		svcName := ReadStringUtil(lex, '#')
		if svcName != "game" {
			t.FailNow()
		}

		lex.Consume(1)

		svcIndex := ExpectInt32(lex)

		if svcIndex != 1 {
			t.FailNow()
		}

		lex.Consume(1)

		group := ReadStringUtil(lex, 0)

		if group != "dev" {
			t.FailNow()
		}

	})
}
