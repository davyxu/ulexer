package ulexer

import (
	"strconv"
	"strings"
)

type Token struct {
	t          string
	lit        string
	begin, end int
	col        int
	line       int
}

func (self *Token) SetType(t string) {
	self.t = t
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

func (self *Token) String() string {
	if self == nil {
		return ""
	}

	return self.lit
}

func (self *Token) Bool() bool {
	if self == nil {
		return false
	}

	if self.lit == "true" {
		return true
	}

	return false
}

func (self *Token) UInt8() uint8 {

	if self == nil {
		return 0
	}

	v, _ := strconv.ParseUint(self.lit, 10, 8)

	return uint8(v)
}

func (self *Token) Int32() int32 {

	if self == nil {
		return 0
	}

	v, _ := strconv.ParseInt(self.lit, 10, 32)

	return int32(v)
}

func (self *Token) UInt32() uint32 {

	if self == nil {
		return 0
	}

	v, _ := strconv.ParseUint(self.lit, 10, 32)

	return uint32(v)
}

func (self *Token) Int64() int64 {

	if self == nil {
		return 0
	}

	v, _ := strconv.ParseInt(self.lit, 10, 64)

	return v
}

func (self *Token) UInt64() uint64 {

	if self == nil {
		return 0
	}

	v, _ := strconv.ParseUint(self.lit, 10, 64)

	return v
}

func (self *Token) Float32() float32 {

	if self == nil {
		return 0
	}

	v, _ := strconv.ParseFloat(self.lit, 32)

	return float32(v)
}

func (self *Token) Float64() float64 {

	if self == nil {
		return 0
	}

	v, _ := strconv.ParseFloat(self.lit, 64)

	return v
}

// bitSize: 32/64 位, base: 10/16进制, unsign是否带符号, 自动转换为指定数值
func (self *Token) Numeral(bitSize, base int, signed bool) interface{} {

	if strings.Contains(self.lit, ".") {

		v, _ := strconv.ParseFloat(self.lit, bitSize)

		if bitSize == 32 {
			return float32(v)
		}

		return v
	} else {

		if signed {
			v, _ := strconv.ParseInt(self.lit, base, bitSize)

			switch bitSize {
			case 16:
				return int16(v)
			case 32:
				return int32(v)
			case 64:
				return v
			default:
				panic("unknown numeral bitsize")
			}

		} else {
			v, _ := strconv.ParseUint(self.lit, base, bitSize)

			switch bitSize {
			case 16:
				return uint16(v)
			case 32:
				return uint32(v)
			case 64:
				return v
			default:
				panic("unknown numeral bitsize")
			}
		}
	}

}

func MakeRawToken(t, literal string) *Token {
	return &Token{t: t, lit: literal}
}
