package ulexer

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
)

type State struct {
	pos  int
	line int
	col  int
}

type Lexer struct {
	src []rune
	State
	preHooker Hooker
}

func (self *Lexer) String() string {
	return fmt.Sprintf("%d,%d", self.line, self.col)
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

func (self *Lexer) Read(m Matcher) *Token {

	if self.preHooker != nil {

		// hooker里不能用hooker
		h := self.preHooker
		self.preHooker = nil

		tk := h(self)

		self.preHooker = h

		if tk != nil {
			return tk
		}
	}

	return m.Read(self)
}

type Hooker func(lex *Lexer) *Token

func (self *Lexer) SetPreHook(hook Hooker) (pre Hooker) {
	pre = self.preHooker
	self.preHooker = hook
	return
}

func NewLexer(s string) *Lexer {

	self := &Lexer{
		src: []rune(s),
	}

	self.State.line = 1
	self.State.col = 1

	return self
}

func (self *Lexer) Code() string {

	reader := bufio.NewReader(strings.NewReader(string(self.src)))

	var sb strings.Builder
	line := 0
	for {
		lineStr, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line++

		sb.WriteString(fmt.Sprintf("%d:	%s", line, lineStr))
	}

	return sb.String()
}
