package golexer2

import (
	"fmt"
	"strconv"
)

type Token struct {
	t          string
	lit        string
	begin, end int
	col        int
	line       int
}

func (self *Token) Pos() int {
	return self.begin
}

func (self *Token) Type() string {
	if self == nil {
		return ""
	}

	return self.t
}

func (self *Token) ToString() string {
	if self == nil {
		return ""
	}

	return self.lit
}

func (self *Token) ToInt32() int32 {

	if self == nil {
		return 0
	}

	v, err := strconv.ParseInt(self.lit, 10, 32)

	if err != nil {
		panic(fmt.Sprintf("strconv.ParseInt '%s', %s", self.lit, err))
		return 0
	}

	return int32(v)
}

var (
	EmptyToken = Token{}
)
