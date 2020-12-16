package string

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

//截取字符串，并返回实际截取的长度和子串
func SubStr(str string, start, end int64 )(int64, string){
	reader:= strings.NewReader(str)

	// Calling NewSectionReader method with its parameters
	r:= io.NewSectionReader(reader, start, end)

	// Calling Copy method with its parameters
	var buf bytes.Buffer
	n, err:= io.Copy(&buf, r)
	if err != nil {
		panic(err)
	}
	return n, buf.String()
}

//在字符串中查找指定子串，并返回left或right部分
func SubstrTarget(str string, target string, turn string, hasPos bool) string {
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

//连接字符串和其他类型
//fmt.Println(JoinStringAndInt("why", 123))
func JoinStringAndInt(val... interface{}) string {
	return fmt.Sprint(val...)
}


