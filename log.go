package golexer2

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func ExpectLogResult(lex *Lexer, expect string, t *testing.T) {

	got := strings.TrimSpace(lex.logger.GetOutput().(fmt.Stringer).String())
	if got != expect {
		t.Logf("Expect log output '%s' got '%s'", expect, got)
		t.FailNow()
	}

}

type outputCacher struct {
	sb strings.Builder
}

func (self *outputCacher) String() string {
	return self.sb.String()
}

func (self *outputCacher) Write(p []byte) (n int, err error) {

	self.sb.Write(p)

	return os.Stdout.Write(p)
}
