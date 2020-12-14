package string

import (
	"strings"
	"unicode/utf8"
)

//在字符串中查找指定字串，并返回left或right部分
func Substr(str string, target string, turn string, hasPos bool) string {
	pos := strings.Index(str, target)

	if pos == -1 {
		return ""
	}

	if turn == "left" {
		if hasPos == true {
			pos = pos + 1
		}
		return str[:pos]
	} else if turn == "right" {
		if hasPos == false {
			pos = pos + 1
		}
		return str[pos:]
	} else {
		panic("params 3 error")
	}
}

func GetStringUtf8Len(str string) int{
	return utf8.RuneCountInString(str)
}

func Utf8Index(str, substr string) int {
	index := strings.Index(str, substr)
	if index < 0{
		return -1
	}
	return utf8.RuneCountInString(str[:index])
}

