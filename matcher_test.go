package ulexer

import (
	"strings"
	"testing"
)

type TestLexer struct {
	runErr error
	lex    *Lexer
}

func (self *Lexer) DebugPrint(t *testing.T) {
	for !self.EOF() {

		tk := self.Try(
			Numbers(),
			WhiteSpace(),
			Identifier(),
			CBlockComment(),
			CLineComment(),
			UnixLineComment(),
			Etc(),
		)
		if tk != nil {

			t.Logf("%02d %10s '%s'", tk.Pos(), tk.Type(), tk.ToString())

		} else {
			t.Log("unknown tk", self.Pos())
			break
		}
	}
}

func (self *TestLexer) Run(text string, callback func(lex *Lexer)) *TestLexer {

	self.lex = NewLexer(text)

	self.runErr = self.lex.Run(callback)

	return self
}

func (self *TestLexer) MustNoError(t *testing.T) *TestLexer {
	if self.runErr != nil {
		t.FailNow()
	}
	return self
}

func (self *TestLexer) MustEOF(t *testing.T) *TestLexer {
	if !self.lex.EOF() {
		t.FailNow()
	}
	return self
}

func (self *TestLexer) ExpectError(t *testing.T, str string) {

	if strings.TrimSpace(self.runErr.Error()) != str {
		t.FailNow()
	}
}

func TestExpect(t *testing.T) {
	new(TestLexer).Run("1", func(lex *Lexer) {
		lex.Expect(Letters())
	}).ExpectError(t, "Expect Letter")
}

// 回车处理
func TestLineEnd(t *testing.T) {

	// Mac OS 9 以及之前的系统的换行符是 CR，从 Mac OS X （后来改名为“OS X”）开始的换行符是 LF即‘\n'，和Unix/Linux统一了。
	new(TestLexer).Run("\n \r \r\n", func(lex *Lexer) {
		for !lex.EOF() {
			lex.Skip(AnyChar(' ', '\t'))
			lex.Skip(LineEnd())
		}

		if lex.Line() != 3 {
			t.FailNow()
		}
	}).MustNoError(t).MustEOF(t)

}

func TestSvcID(t *testing.T) {

	new(TestLexer).Run("game#1@dev", func(lex *Lexer) {
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
	}).MustNoError(t).MustEOF(t)
}

// 标识符识别
func TestExpectIdentifier(t *testing.T) {
	new(TestLexer).Run("b1 full 1c", func(lex *Lexer) {
		if lex.Expect(Identifier()).ToString() != "b1" {
			t.FailNow()
		}
		lex.Skip(WhiteSpace())
		if lex.Expect(Identifier()).ToString() != "full" {
			t.FailNow()
		}
		lex.Skip(WhiteSpace())
		lex.Expect(Identifier())
	}).ExpectError(t, "Expect Identifier")

}

// parser尝试逻辑
func TestTryList(t *testing.T) {

	new(TestLexer).Run("b1 full 1c game#1@dev", func(lex *Lexer) {

		for !lex.EOF() {

			tk := lex.Try(
				Numbers(),
				WhiteSpace(),
				Identifier(),
				Etc(),
			)
			if tk != nil {

				t.Logf("%02d %10s '%s'", tk.Pos(), tk.Type(), tk.ToString())

			} else {
				t.Log("unknown tk", lex.Pos())
				break
			}
		}

	}).MustNoError(t).MustEOF(t)
}

// C行注释
func TestCLineComment(t *testing.T) {

	new(TestLexer).Run(`// abc
123`, func(lex *Lexer) {

		lex.Expect(CLineComment())
		lex.Expect(LineEnd())
		lex.Expect(Numbers())

	}).MustNoError(t).MustEOF(t)
}

// C块注释
func TestCBlockComment(t *testing.T) {

	new(TestLexer).Run(`/*a
1*2/3*/
456`, func(lex *Lexer) {

		lex.Expect(CBlockComment())
		lex.Expect(WhiteSpace())
		lex.Expect(Numbers())

	}).MustNoError(t).MustEOF(t)
}

// Unix行注释
func TestUnixLineComment(t *testing.T) {

	new(TestLexer).Run(`# abc
123`, func(lex *Lexer) {

		lex.Expect(UnixLineComment())
		lex.Expect(LineEnd())
		lex.Expect(Numbers())

	}).MustNoError(t).MustEOF(t)
}
