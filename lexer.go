package golexer2

import (
	"fmt"
	"runtime"
)

type Lexer struct {
	src   []rune
	index int
}

func (self *Lexer) Current() rune {

	if self.EOF() {
		return 0
	}

	return self.src[self.index]
}

func (self *Lexer) Index() int {
	return self.index
}

func (self *Lexer) Count() int {
	return len(self.src)
}

func (self *Lexer) Peek(offset int) rune {

	if self.index+offset >= len(self.src) {
		return 0
	}

	return self.src[self.index+offset]
}

func (self *Lexer) Consume(n int) {
	self.index += n
}

func (self *Lexer) EOF() bool {
	return self.index >= len(self.src)
}

func (self *Lexer) StringRange(count int) string {

	end := self.index + count

	if end > len(self.src) {
		end = len(self.src)
	}

	return string(self.src[self.index:end])
}

func (self *Lexer) Error(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}

type MatchFunc func(lex *Lexer, index int, r rune) interface{}

func (self *Lexer) Visit(match MatchFunc) (ret interface{}, ok bool) {

	var count int
	for {

		r := self.Peek(count)

		raw := match(self, count, r)

		if op, isOpOK := raw.(MatchOp); isOpOK {
			switch op {
			case MatchOp_Stop:
				return
			case MatchOp_Next:
			default:
				panic("unknown op")
			}
		} else {
			ret = raw
			ok = true
			return
		}

		count++
	}

	return
}

func (self *Lexer) Run(callback func(lex *Lexer)) {

	defer func() {

		switch err := recover().(type) {
		case runtime.Error:
			panic(err)
		default:
			log.Errorln(err)
		}

	}()

	callback(self)

}

func NewLexer(s []rune) *Lexer {

	return &Lexer{
		src: s,
	}
}
