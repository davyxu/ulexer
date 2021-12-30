package ulexer

import (
	"fmt"
	"runtime"
)

func Expect(lex *Lexer, m Matcher) *Token {

	tk := lex.Read(m)

	if tk == nil {

		var target string
		if s, ok := m.(fmt.Stringer); ok {
			target = s.String()
		} else {
			target = m.TokenType()
		}

		lex.Error("expect '%s'", target)
	}

	return tk
}

func Ignore(lex *Lexer, m Matcher) {
	state := lex.State
	tk := lex.Read(m)

	if tk == nil && !lex.EOF() {
		lex.State = state
	}
}

func Is(lex *Lexer, m Matcher, refToken **Token) bool {

	tk := lex.Read(m)
	if tk != nil {
		*refToken = tk
		return true
	}

	return false
}

func Try(lex *Lexer, callback func(lex *Lexer)) (retErr error) {

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

	callback(lex)

	return
}

func Select(lex *Lexer, mlist ...Matcher) *Token {

	for _, m := range mlist {

		tk := lex.Read(m)

		if tk != nil {
			return tk
		}
	}

	return nil
}
