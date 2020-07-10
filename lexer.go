package ulexer

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type Lexer struct {
	src  []rune
	pos  int
	line int
	col  int
}

func (self *Lexer) Pos() int {
	return self.pos
}

func (self *Lexer) Count() int {
	return len(self.src)
}

func (self *Lexer) Line() int {
	return self.line
}

func (self *Lexer) Col() int {
	return self.col
}

func (self *Lexer) Peek(offset int) rune {

	if self.pos+offset >= len(self.src) {
		return 0
	}

	return self.src[self.pos+offset]
}

func (self *Lexer) Consume(n int) {

	for i := 0; i < n; i++ {
		r := self.Peek(i)
		switch r {
		case '\n':
			self.line++
		case '\r':
			self.col = 1
		}
	}

	self.pos += n
	self.col += n
}

func (self *Lexer) EOF() bool {
	return self.pos >= len(self.src)
}

func (self *Lexer) ToLiteral(count int) string {

	end := self.pos + count

	if end > len(self.src) {
		end = len(self.src)
	}

	return string(self.src[self.pos:end])
}

func (self *Lexer) Error(format string, args ...interface{}) {
	panic(fmt.Errorf(format, args...))
}

var ErrEOF = errors.New("EOF")

func newToken(self *Lexer, count int, m Matcher) *Token {
	ret := new(Token)
	ret.t = m.TokenType()
	ret.begin = self.Pos()
	ret.end = ret.begin + count
	ret.col = self.Col()
	ret.line = self.Line()

	return ret
}

func (self *Lexer) NewToken(count int, m Matcher) (ret *Token) {
	ret = newToken(self, count, m)
	ret.lit = self.ToLiteral(count)
	return
}

func (self *Lexer) NewTokenLiteral(count int, m Matcher, literal string) (ret *Token) {
	ret = newToken(self, count, m)
	ret.lit = literal
	return
}

func (self *Lexer) Select(mlist ...Matcher) *Token {

	for _, m := range mlist {

		tk := self.Read(m)

		if tk != EmptyToken {
			return tk
		}
	}

	return EmptyToken
}

type MatchAction func(tk *Token)

func (self *Lexer) SelectAction(mlist []Matcher, alist []MatchAction) {

	if len(mlist) != len(alist) {
		panic("Matcher list should equal to Action list length")
	}

	var hit bool
	for index, m := range mlist {
		tk := self.Read(m)

		if tk != EmptyToken {

			action := alist[index]
			if action != nil {
				action(tk)
			}
			hit = true
			break
		}
	}

	if !hit {

		var sb strings.Builder

		for index, m := range mlist {
			if index > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(m.TokenType())
		}

		self.Error("Expect %s", sb.String())
	}

}

func (self *Lexer) Read(m Matcher) *Token {
	if self.EOF() {
		panic(ErrEOF)
	}

	return m.Read(self)
}

func (self *Lexer) Expect(m Matcher) *Token {

	tk := self.Read(m)

	if tk == EmptyToken {

		self.Error("Expect %s, got: %s", m.TokenType(), self.ToLiteral(10))
	}

	return tk
}

func (self *Lexer) Run(callback func(lex *Lexer)) (retErr error) {

	defer func() {

		switch raw := recover().(type) {
		case runtime.Error:
			panic(raw)
		case nil:
		case error:
			if raw != ErrEOF {
				retErr = raw
			}

		default:
			panic(raw)
		}

	}()

	callback(self)

	return
}

func NewLexer(s string) *Lexer {

	self := &Lexer{
		src:  []rune(s),
		line: 1,
		col:  1,
	}

	return self
}
