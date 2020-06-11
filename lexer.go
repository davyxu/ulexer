package ulexer

import (
	"errors"
	"fmt"
	"runtime"
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

type Matcher interface {
	MatchRune(index int, r rune) bool
	TokenType() string
}

var ErrEOF = errors.New("EOF")

func (self *Lexer) read(m Matcher, consume bool) (ret *Token, eof bool) {

	// 上一个步骤有EOF, 抛出后结束解析
	if self.Peek(0) == 0 {
		panic(ErrEOF)
	}

	var index int
	for {

		r := self.Peek(index)

		if r == 0 {
			eof = true
			break
		}

		if !m.MatchRune(index, r) {
			break
		}

		index++
	}

	if index > 0 {
		ret = new(Token)
		ret.lit = self.ToLiteral(index)
		ret.t = m.TokenType()
		ret.begin = self.Pos()
		ret.end = ret.begin + index
		ret.col = self.Col()
		ret.line = self.Line()

		if consume {
			self.Consume(index)
		}

	}

	return
}

func (self *Lexer) Skip(m Matcher) {
	self.read(m, true)
}

func (self *Lexer) Try(mlist ...Matcher) *Token {

	for _, m := range mlist {

		if tk, _ := self.read(m, true); tk != nil {
			return tk
		}
	}

	return nil
}

func (self *Lexer) Expect(m Matcher) *Token {

	tk, eof := self.read(m, true)
	if tk == nil {

		if !eof {
			var str string
			if s, ok := m.(fmt.Stringer); ok {
				str = s.String()
			}

			self.Error("Expect %s %s", m.TokenType(), str)
		}

		return nil
	}

	return tk
}

func (self *Lexer) Is(m Matcher) *Token {
	tk, _ := self.read(m, false)
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
