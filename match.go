package golexer2

import (
	"unicode"
)

type MatchOp int

const (
	MatchOp_Next MatchOp = iota
	MatchOp_Stop
)

// 匹配空白符
func MatchWhiteSpace(lex *Lexer, index int, r rune) interface{} {
	if unicode.In(r, unicode.White_Space) {
		return MatchOp_Next
	}

	lex.Consume(index)

	return index
}

// 匹配数字,直到碰到非数字
func MatchNumber(lex *Lexer, index int, r rune) (ret interface{}) {
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

// 匹配所有字符串直到指定的字符为止
func MatchUtilChar(end rune) MatchFunc {

	return func(lex *Lexer, index int, r rune) (ret interface{}) {

		if end != r && r != 0 {
			return MatchOp_Next
		}

		ret = lex.StringRange(index)

		lex.Consume(index)

		return
	}
}

// 完整匹配指定字符串
func MatchCompleteString(str string) MatchFunc {

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

func MatchAnyChar(data []rune) MatchFunc {

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

func MatchLineEnd() MatchFunc {

	charsToConsume := 0
	return func(lex *Lexer, index int, r rune) interface{} {

		switch r {
		case '\n':
			lex.onNewLine()
			charsToConsume++
		case '\r':
			lex.onReturn()
			charsToConsume++
		default:
			lex.Consume(charsToConsume)
			return MatchOp_Stop
		}

		return MatchOp_Next
	}
}

func ReadStringUtil(lex *Lexer, r rune) string {

	if raw, ok := lex.Visit(MatchUtilChar(r)); ok {
		return raw.(string)
	}

	return ""
}

func SkipLineEnd(lex *Lexer) {
	lex.Visit(MatchLineEnd())
}

func SkipString(lex *Lexer, data []rune) {

	lex.Visit(MatchAnyChar(data))
}

func SkipWhiteSpace(lex *Lexer) {

	lex.Visit(MatchWhiteSpace)
}

func TryString(lex *Lexer, str string) bool {

	if _, ok := lex.Visit(MatchCompleteString(str)); ok {
		return true
	}

	return false
}
