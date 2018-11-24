package golexer2

import (
	"unicode"
)

type MatchOp int

const (
	MatchOp_Next MatchOp = iota
	MatchOp_Stop
)

func ReadWhiteSpace(lex *Lexer, index int, r rune) interface{} {
	if unicode.In(r, unicode.White_Space) {
		return MatchOp_Next
	}

	lex.Consume(index)

	return index
}

func ReadNumber(lex *Lexer, index int, r rune) (ret interface{}) {
	isDigit := unicode.IsDigit(r)

	// 首字符必须是数字
	if index == 0 {
		if !isDigit {
			return MatchOp_Stop
		}
	}

	if isDigit {
		return MatchOp_Next
	}

	ret = lex.StringRange(index)

	lex.Consume(index)

	return
}

func ReadUtilChar(end rune) MatchFunc {

	return func(lex *Lexer, index int, r rune) (ret interface{}) {

		if end != r && r != 0 {
			return MatchOp_Next
		}

		ret = lex.StringRange(index)

		lex.Consume(index)

		return
	}
}

func MatchString(str string) MatchFunc {

	data := []rune(str)

	return func(lex *Lexer, index int, r rune) (ret interface{}) {

		if index >= len(data) {
			lex.Consume(index)
			return
		}

		if data[index] != r {
			return MatchOp_Stop
		}

		return MatchOp_Next
	}
}

func ContainString(str string) MatchFunc {

	data := []rune(str)

	return func(lex *Lexer, index int, r rune) (ret interface{}) {

		for _, dr := range data {
			if dr == r {
				return MatchOp_Next
			}
		}

		lex.Consume(index)

		return MatchOp_Stop
	}
}

func ReadStringUtil(lex *Lexer, r rune) string {

	if raw, ok := lex.Visit(ReadUtilChar('#')); ok {
		return raw.(string)
	}

	return ""
}

func SkipString(lex *Lexer, str string) {

	lex.Visit(ContainString(str))
}

func SkipWhiteSpace(lex *Lexer) {

	lex.Visit(ReadWhiteSpace)
}

func TryString(lex *Lexer, str string) bool {

	if _, ok := lex.Visit(MatchString(str)); ok {
		return true
	}

	return false
}
