package golexer2

import (
	"github.com/davyxu/golog"
	"io"
	"strings"
	"testing"
)

var (
	log        = golog.New("golexer2")
	logCatcher = new(outputCacher)
	preOutput  io.Writer
)

func ExpectLogResult(expect string, t *testing.T) {

	got := logCatcher.String()
	if got != expect {
		t.Logf("Expect log output '%s' go '%s'", expect, got)
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

	return preOutput.Write(p)
}

func init() {
	log.SetParts()
	preOutput = log.GetOutput()
	log.SetOutptut(logCatcher)
}
