# ulexer
A simple lexer better than parse combinator


# Example

```
func TestSvcID(t *testing.T) {

	new(TestLexer).Run("game#1@dev", func(lex *Lexer) {
		svcName := lex.Expect(Identifier()).ToString()
		if svcName != "game" {
			t.FailNow()
		}

		// 跳过#
		lex.Expect(Contain('#'))

		svcIndex := lex.Expect(Numeral()).ToInt32()

		if svcIndex != 1 {
			t.FailNow()
		}

		// 跳过@
		lex.Expect(Contain('@'))

		group := lex.Expect(Identifier()).ToString()

		if group != "dev" {
			t.FailNow()
		}
	}).MustNoError(t).MustEOF(t)
}
```