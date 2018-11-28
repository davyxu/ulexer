package golexer2

import "unicode"

type whiteSpaceMatcher struct {
}

func (whiteSpaceMatcher) MatchRune(index int, r rune) bool {
	return unicode.In(r, unicode.White_Space)
}

func (whiteSpaceMatcher) TokenType() string {
	return "WhiteSpace"
}

var _whiteSpaceMatcher = new(whiteSpaceMatcher)

// 匹配空白符
func WhiteSpace() Matcher {
	return _whiteSpaceMatcher
}

type numberMatcher struct {
}

func (numberMatcher) MatchRune(index int, r rune) bool {
	return unicode.IsDigit(r)
}

func (numberMatcher) TokenType() string {
	return "Number"
}

var _numberMatcher = new(numberMatcher)

// 匹配数字
func Numbers() Matcher {
	return _numberMatcher
}

type letterMatcher struct {
}

func (letterMatcher) MatchRune(index int, r rune) bool {
	return unicode.IsLetter(r)
}

func (letterMatcher) TokenType() string {
	return "Letter"
}

var _letterMatcher = new(letterMatcher)

// 匹配字母
func Letters() Matcher {
	return _letterMatcher
}

type etcMatcher struct {
}

func (etcMatcher) MatchRune(index int, r rune) bool {

	if r == 0 {
		return false
	}

	return !unicode.IsLetter(r) && !unicode.IsDigit(r)
}

func (etcMatcher) TokenType() string {
	return "Etc"
}

var _etcMatcher = new(etcMatcher)

// 匹配除字母数字之外的字符
func Etc() Matcher {
	return _etcMatcher
}

type lineEndMatcher struct {
}

func (lineEndMatcher) MatchRune(index int, r rune) bool {
	return r == '\r' || r == '\n'
}

func (lineEndMatcher) TokenType() string {
	return "LineEnd"
}

var _lineEndMatcher = new(lineEndMatcher)

// 匹配行结尾
func LineEnd() Matcher {
	return _lineEndMatcher
}

type anyCharMatcher struct {
	list []rune
}

func (self *anyCharMatcher) String() string {
	return string(self.list)
}

func (self *anyCharMatcher) MatchRune(index int, r rune) bool {
	for _, libr := range self.list {
		if libr == r {
			return true
		}
	}

	return false
}

func (anyCharMatcher) TokenType() string {
	return "AnyChar"
}

// 匹配给定的任意字符
func AnyChar(list ...rune) Matcher {
	return &anyCharMatcher{
		list: list,
	}
}

type containCharMatcher struct {
	list []rune
}

func (self *containCharMatcher) String() string {
	return string(self.list)
}

func (self *containCharMatcher) MatchRune(index int, r rune) bool {
	for _, libr := range self.list {
		if libr != r {
			return false
		}
	}

	return true
}

func (containCharMatcher) TokenType() string {
	return "ContainChar"
}

// 匹配指定的所有字符
func ContainChar(list ...rune) Matcher {
	return &containCharMatcher{
		list: list,
	}
}

// 匹配指定的字符串
func ContainString(str string) Matcher {
	return &containCharMatcher{
		list: []rune(str),
	}
}

type identifierMatcher struct {
}

func (identifierMatcher) MatchRune(index int, r rune) bool {
	if index == 0 {
		if unicode.IsLetter(r) || r == '_' {
			return true
		}

	} else {

		if unicode.IsLetter(r) || r == '_' || unicode.IsDigit(r) {
			return true
		}
	}

	return false
}

func (identifierMatcher) TokenType() string {
	return "Identifier"
}

var _identifierMatcher = new(identifierMatcher)

// 匹配标识符
func Identifier() Matcher {
	return _identifierMatcher
}
