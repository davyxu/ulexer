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
			self.col = 0
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
	Read(lex *Lexer) *Token
	TokenType() string
}

var ErrEOF = errors.New("EOF")

func (self *Lexer) NewToken(count int, m Matcher) (ret *Token) {
	ret = new(Token)
	ret.lit = self.ToLiteral(count)
	ret.t = m.TokenType()
	ret.begin = self.Pos()
	ret.end = ret.begin + count
	ret.col = self.Col()
	ret.line = self.Line()
	return
}

func (self *Lexer) Expect(m Matcher) *Token {

	if self.EOF() {
		panic(ErrEOF)
	}

	tk := m.Read(self)

	if tk == EmptyToken {
		var str string
		if s, ok := m.(fmt.Stringer); ok {
			str = s.String()
		}

		self.Error("Expect %s %s", m.TokenType(), str)
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
		line: 0,
		col:  0,
	}

	return self
}
