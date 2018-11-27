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

// 回车符的识别
func TestLine(t *testing.T) {
	lex := NewLexer([]rune("\n"))

	lex.Run(func(lex *Lexer) {
		SkipString(lex, " \t")
		ExpectString(lex, "\n")
		ExpectLogResult("", t)
	})
}

func TestSvcID(t *testing.T) {
	lex := NewLexer([]rune("game#1@dev"))

	lex.Run(func(lex *Lexer) {

		svcName := ReadStringUtil(lex, '#')
		if svcName != "game" {
			t.FailNow()
		}

		// 跳过#
		lex.Consume(1)

		svcIndex := ExpectInt32(lex)

		if svcIndex != 1 {
			t.FailNow()
		}

		// 跳过@
		lex.Consume(1)

		group := ReadStringUtil(lex, 0)

		if group != "dev" {
			t.FailNow()
		}

	})
}
