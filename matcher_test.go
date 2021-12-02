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

func (self *TestLexer) Try(text string, callback func(lex *Lexer)) *TestLexer {

	self.lex = NewLexer(text)

	self.runErr = Try(self.lex, callback)

	return self
}

func (self *TestLexer) RunUnit(ulist []*TestUnitInput, t *testing.T) *TestLexer {

	for _, u := range ulist {

		self.lex = NewLexer(u.Input)

		self.runErr = Try(self.lex, func(lex *Lexer) {

			var expect string
			if u.Expect == "" {
				expect = u.Input
			}

			if Expect(lex, u.Matcher).String() != expect {
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

func (self *TestLexer) ContainError(t *testing.T, str string) {

	if self.runErr == nil || !strings.Contains(self.runErr.Error(), str) {
		t.FailNow()
	}
}

func TestNumeral(t *testing.T) {

	new(TestLexer).Try("1 3.14 -2.5", func(lex *Lexer) {
		if Expect(lex, Numeral()).String() != "1" {
			t.FailNow()
		}

		if Expect(lex, WhiteSpace()).String() != " " {
			t.FailNow()
		}

		if Expect(lex, Numeral()).String() != "3.14" {
			t.FailNow()
		}

		if Expect(lex, WhiteSpace()).String() != " " {
			t.FailNow()
		}

		if Expect(lex, Numeral()).String() != "-2.5" {
			t.FailNow()
		}
	}).MustNoError(t)

}

// 回车处理
func TestLineEnd(t *testing.T) {

	new(TestLexer).Try("\n1\r2\r\n", func(lex *Lexer) {

		Expect(lex, LineEnd())

		if Expect(lex, Numeral()).String() != "1" {
			t.FailNow()
		}

		Expect(lex, LineEnd())

		if Expect(lex, Numeral()).String() != "2" {
			t.FailNow()
		}

		Expect(lex, LineEnd())

	}).MustNoError(t).MustEOF(t)

}

// 标识符识别
func TestExpectIdentifier(t *testing.T) {
	new(TestLexer).Try("b1 full _c", func(lex *Lexer) {
		if Expect(lex, Identifier()).String() != "b1" {
			t.FailNow()
		}
		Expect(lex, WhiteSpace())
		if Expect(lex, Identifier()).String() != "full" {
			t.FailNow()
		}
		Expect(lex, WhiteSpace())
		Expect(lex, Identifier())
	}).MustNoError(t)

	new(TestLexer).Try("a 1c", func(lex *Lexer) {
		if Expect(lex, Identifier()).String() != "a" {
			t.FailNow()
		}
		Expect(lex, WhiteSpace())
		Expect(lex, Identifier()) // 触发错误
	}).ContainError(t, "Expect Identifier")

}

// C行注释
func TestCLineComment(t *testing.T) {

	new(TestLexer).Try(`// abc
123`, func(lex *Lexer) {

		Expect(lex, CLineComment())
		Expect(lex, LineEnd())
		Expect(lex, Numeral())

	}).MustNoError(t).MustEOF(t)
}

// Unix行注释
func TestUnixLineComment(t *testing.T) {

	new(TestLexer).Try(`# abc
123`, func(lex *Lexer) {

		Expect(lex, UnixLineComment())
		Expect(lex, LineEnd())
		Expect(lex, Numeral())

	}).MustNoError(t).MustEOF(t)
}

func TestSvcID(t *testing.T) {

	new(TestLexer).Try("game#1@dev", func(lex *Lexer) {
		svcName := Expect(lex, Identifier()).String()
		if svcName != "game" {
			t.FailNow()
		}

		// 跳过#
		Expect(lex, Contain('#'))

		svcIndex := Expect(lex, Numeral()).Int32()

		if svcIndex != 1 {
			t.FailNow()
		}

		// 跳过@
		Expect(lex, Contain('@'))

		group := Expect(lex, Identifier()).String()

		if group != "dev" {
			t.FailNow()
		}
	}).MustNoError(t).MustEOF(t)
}

func TestString(t *testing.T) {

	new(TestLexer).Try(`a "hello" 'world'`, func(lex *Lexer) {
		if Expect(lex, Identifier()).String() != "a" {
			t.FailNow()
		}

		Expect(lex, WhiteSpace())

		if Expect(lex, String()).String() != "hello" {
			t.FailNow()
		}

		Expect(lex, WhiteSpace())

		if Expect(lex, String()).String() != "world" {
			t.FailNow()
		}

	}).MustNoError(t)

}

func TestIgnore(t *testing.T) {
	new(TestLexer).Try(` a b '`, func(lex *Lexer) {

		Ignore(lex, WhiteSpace())
		Ignore(lex, WhiteSpace())
		Expect(lex, Contain('a'))
		Ignore(lex, WhiteSpace())
		Expect(lex, Contain('b'))
		Ignore(lex, WhiteSpace())

	}).MustNoError(t)

}
