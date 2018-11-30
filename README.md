# golexer2
A simple lexer better than parse combinator


# Example

```
func TestSvcID(t *testing.T) {

	new(TestLexer).Run([]rune("game#1@dev"), func(lex *Lexer) {
		svcName := lex.Expect(Letters()).ToString()
		if svcName != "game" {
			t.FailNow()
		}

		// 跳过#
		lex.Expect(ContainChar('#'))

		svcIndex := lex.Expect(Numbers()).ToInt32()

		if svcIndex != 1 {
			t.FailNow()
		}

		// 跳过@
		lex.Expect(ContainChar('@'))

		group := lex.Expect(Letters()).ToString()

		if group != "dev" {
			t.FailNow()
		}
	}).MustNoError(t).MustEOF(t)
}
```