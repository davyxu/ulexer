package golexer2

import (
	"testing"
)

func TestExpect(t *testing.T) {
	testMode = true

	lex := NewLexer([]rune("1"))

	lex.Run(func(lex *Lexer) {
		lex.Expect(Letters())
	})

	ExpectLogResult(lex, "Expect Letter", t)
}

// 回车处理
func TestLineEnd(t *testing.T) {
	testMode = true

	// Mac OS 9 以及之前的系统的换行符是 CR，从 Mac OS X （后来改名为“OS X”）开始的换行符是 LF即‘\n'，和Unix/Linux统一了。
	lex := NewLexer([]rune("\n \r \r\n"))

	lex.Run(func(lex *Lexer) {
		for !lex.EOF() {

			lex.Skip(AnyChar(' ', '\t'))
			lex.Skip(LineEnd())
		}

		if lex.Line() != 3 {
			t.FailNow()
		}
	})
}

func TestSvcID(t *testing.T) {
	testMode = true

	lex := NewLexer([]rune("game#1@dev"))

	lex.Run(func(lex *Lexer) {

		svcName := lex.Expect(Letters()).ToString()
		if svcName != "game" {
			t.FailNow()
		}

		// 跳过#
		lex.Expect(ContainChar('#'))

		svcIndex := lex.Expect(Numbers()).ToInt32()

		if svcIndex != 1 {
			t.FailNow()
		}

		// 跳过@
		lex.Expect(ContainChar('@'))

		group := lex.Expect(Letters()).ToString()

		if group != "dev" {
			t.FailNow()
		}

	})
}

// 标识符识别
func TestExpectIdentifier(t *testing.T) {
	testMode = true

	lex := NewLexer([]rune("b1 full 1c"))

	lex.Run(func(lex *Lexer) {

		if lex.Expect(Identifier()).ToString() != "b1" {
			t.FailNow()
		}
		lex.Skip(WhiteSpace())
		if lex.Expect(Identifier()).ToString() != "full" {
			t.FailNow()
		}
		lex.Skip(WhiteSpace())
		lex.Expect(Identifier())

	})

	ExpectLogResult(lex, "Expect Identifier", t)
}

// parser尝试逻辑
func TestTryList(t *testing.T) {
	testMode = true
	lex := NewLexer([]rune("game#1@dev"))

	lex.Run(func(lex *Lexer) {

		for !lex.EOF() {

			tk := lex.Try(Numbers(), WhiteSpace(), Identifier(), Etc())
			if tk != nil {

				t.Log(tk.ToString(), tk.Type(), tk.Pos())

			} else {
				t.Log("unknown tk", lex.Pos())
				break
			}
		}

	})
}
