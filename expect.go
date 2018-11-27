package golexer2

import "strconv"

func ExpectString(lex *Lexer, str string) {

	if _, ok := lex.Visit(MatchCompleteString(str)); !ok {
		lex.Error("Expect '%s'", str)
	}
}

func ExpectInt32(lex *Lexer) int32 {

	raw, ok := lex.Visit(MatchNumber)

	if ok {

		v, err := strconv.ParseInt(raw.(string), 10, 32)

		if err != nil {
			lex.Error("Invalid integer '%s'", raw.(string))
		} else {
			return int32(v)
		}
	} else {
		lex.Error("Expect integer")
	}

	return 0
}
