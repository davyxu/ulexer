package ulexer

import (
	"strings"
	"testing"
)

type TestUnitInput struct {
	Input   string
	Matcher Matcher
	Expect  string
}

type TestLexer struct {
	runErr error
	lex    *Lexer
}

func (self *TestLexer) Run(text string, callback func(lex *Lexer)) *TestLexer {

	self.lex = NewLexer(text)

	self.runErr = self.lex.Run(callback)

	return self
}

func (self *TestLexer) RunUnit(ulist []*TestUnitInput, t *testing.T) *TestLexer {

	for _, u := range ulist {

		self.lex = NewLexer(u.Input)

		self.runErr = self.lex.Run(func(lex *Lexer) {

			var expect string
			if u.Expect == "" {
				expect = u.Input
			}

			if lex.Expect(u.Matcher).String() != expect {
				t.FailNow()
			}
		})
	}

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

	if self.runErr == nil || strings.TrimSpace(self.runErr.Error()) != str {
		t.FailNow()
	}
}

func TestNumeral(t *testing.T) {

	new(TestLexer).Run("1 3.14 -2.5", func(lex *Lexer) {
		if lex.Expect(Numeral()).String() != "1" {
			t.FailNow()
		}

		if lex.Expect(WhiteSpace()).String() != " " {
			t.FailNow()
		}

		if lex.Expect(Numeral()).String() != "3.14" {
			t.FailNow()
		}

		if lex.Expect(WhiteSpace()).String() != " " {
			t.FailNow()
		}

		if lex.Expect(Numeral()).String() != "-2.5" {
			t.FailNow()
		}
	}).MustNoError(t)

}

// 回车处理
func TestLineEnd(t *testing.T) {

	new(TestLexer).Run("\n1\r2\r\n", func(lex *Lexer) {

		lex.Expect(LineEnd())

		if lex.Expect(Numeral()).String() != "1" {
			t.FailNow()
		}

		lex.Expect(LineEnd())

		if lex.Expect(Numeral()).String() != "2" {
			t.FailNow()
		}

		lex.Expect(LineEnd())

	}).MustNoError(t).MustEOF(t)

}

// 标识符识别
func TestExpectIdentifier(t *testing.T) {
	new(TestLexer).Run("b1 full _c", func(lex *Lexer) {
		if lex.Expect(Identifier()).String() != "b1" {
			t.FailNow()
		}
		lex.Expect(WhiteSpace())
		if lex.Expect(Identifier()).String() != "full" {
			t.FailNow()
		}
		lex.Expect(WhiteSpace())
		lex.Expect(Identifier())
	}).MustNoError(t)

	new(TestLexer).Run("a 1c", func(lex *Lexer) {
		if lex.Expect(Identifier()).String() != "a" {
			t.FailNow()
		}
		lex.Expect(WhiteSpace())
		lex.Expect(Identifier()) // 触发错误
	}).ExpectError(t, "Expect Identifier")

}

// C行注释
func TestCLineComment(t *testing.T) {

	new(TestLexer).Run(`// abc
123`, func(lex *Lexer) {

		lex.Expect(CLineComment())
		lex.Expect(LineEnd())
		lex.Expect(Numeral())

	}).MustNoError(t).MustEOF(t)
}

// Unix行注释
func TestUnixLineComment(t *testing.T) {

	new(TestLexer).Run(`# abc
123`, func(lex *Lexer) {

		lex.Expect(UnixLineComment())
		lex.Expect(LineEnd())
		lex.Expect(Numeral())

	}).MustNoError(t).MustEOF(t)
}

func TestSvcID(t *testing.T) {

	new(TestLexer).Run("game#1@dev", func(lex *Lexer) {
		svcName := lex.Expect(Identifier()).String()
		if svcName != "game" {
			t.FailNow()
		}

		// 跳过#
		lex.Expect(Contain('#'))

		svcIndex := lex.Expect(Numeral()).Int32()

		if svcIndex != 1 {
			t.FailNow()
		}

		// 跳过@
		lex.Expect(Contain('@'))

		group := lex.Expect(Identifier()).String()

		if group != "dev" {
			t.FailNow()
		}
	}).MustNoError(t).MustEOF(t)
}

func TestString(t *testing.T) {

	new(TestLexer).Run(`a "hello" 'world'`, func(lex *Lexer) {
		if lex.Expect(Identifier()).String() != "a" {
			t.FailNow()
		}

		lex.Expect(WhiteSpace())

		if lex.Expect(String()).String() != "hello" {
			t.FailNow()
		}

		lex.Expect(WhiteSpace())

		if lex.Expect(String()).String() != "world" {
			t.FailNow()
		}

	}).MustNoError(t)

}
