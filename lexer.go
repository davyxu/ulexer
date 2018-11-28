package golexer2

import (
	"fmt"
	"github.com/davyxu/golog"
	"runtime"
)

type Lexer struct {
	src    []rune
	pos    int
	line   int
	col    int
	logger *golog.Logger
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
	panic(fmt.Sprintf(format, args...))
}

type Matcher interface {
	MatchRune(index int, r rune) bool
	TokenType() string
}

func (self *Lexer) match(m Matcher) (ret *Token) {

	var index int
	for {

		r := self.Peek(index)

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
		self.Consume(index)
	}

	return
}

func (self *Lexer) Skip(m Matcher) {
	self.match(m)
}

func (self *Lexer) Try(mlist ...Matcher) *Token {

	for _, m := range mlist {

		if tk := self.match(m); tk != nil {
			return tk
		}
	}

	return nil
}

func (self *Lexer) Expect(m Matcher) *Token {

	tk := self.match(m)
	if tk == nil {
		var str string
		if s, ok := m.(fmt.Stringer); ok {
			str = s.String()
		}

		self.Error("Expect %s %s", m.TokenType(), str)
		return nil
	}

	return tk
}

func (self *Lexer) Run(callback func(lex *Lexer)) {

	defer func() {

		switch err := recover().(type) {
		case runtime.Error:
			panic(err)
		case nil:
		default:
			self.logger.Errorf("%s", err)
		}

	}()

	callback(self)

}

var (
	testMode bool
)

func NewLexer(s []rune) *Lexer {

	self := &Lexer{
		src:    s,
		line:   1,
		col:    1,
		logger: golog.New("golexer2"),
	}

	self.logger.SetParts()

	if testMode {
		golog.ClearAll() // 解决logger会重名
		self.logger.SetOutptut(new(outputCacher))
	}

	return self
}
