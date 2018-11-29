package golexer2

import "unicode"

type whiteSpaceMatcher int

func (*whiteSpaceMatcher) MatchRune(index int, r rune) bool {
	return unicode.In(r, unicode.White_Space)
}

func (*whiteSpaceMatcher) TokenType() string {
	return "WhiteSpace"
}

// 匹配空白符
func WhiteSpace() Matcher {
	return (*whiteSpaceMatcher)(nil)
}

type numberMatcher int

func (*numberMatcher) MatchRune(index int, r rune) bool {
	return unicode.IsDigit(r)
}

func (*numberMatcher) TokenType() string {
	return "Number"
}

// 匹配数字
func Numbers() Matcher {
	return (*numberMatcher)(nil)
}

type letterMatcher int

func (*letterMatcher) MatchRune(index int, r rune) bool {
	return unicode.IsLetter(r)
}

func (*letterMatcher) TokenType() string {
	return "Letter"
}

// 匹配字母
func Letters() Matcher {
	return (*letterMatcher)(nil)
}

type etcMatcher int

func (*etcMatcher) MatchRune(index int, r rune) bool {

	if r == 0 {
		return false
	}

	return !unicode.IsLetter(r) && !unicode.IsDigit(r)
}

func (*etcMatcher) TokenType() string {
	return "Etc"
}

// 匹配除字母数字之外的字符
func Etc() Matcher {
	return (*etcMatcher)(nil)
}

type lineEndMatcher int

func (*lineEndMatcher) MatchRune(index int, r rune) bool {
	return r == '\r' || r == '\n'
}

func (*lineEndMatcher) TokenType() string {
	return "LineEnd"
}

// 匹配行结尾
func LineEnd() Matcher {
	return (*lineEndMatcher)(nil)
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

func (*anyCharMatcher) TokenType() string {
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

func (*containCharMatcher) TokenType() string {
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

type identifierMatcher int

func (*identifierMatcher) MatchRune(index int, r rune) bool {
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

func (*identifierMatcher) TokenType() string {
	return "Identifier"
}

// 匹配标识符
func Identifier() Matcher {
	return (*identifierMatcher)(nil)
}
