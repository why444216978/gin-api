package file

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	util_err "gin-api/libraries/util/error"
)

//使用io.WriteString()函数进行数据的写入，不存在则创建
func WriteWithIo(filePath, content string) error {
	//os.O_WRONLY | os.O_CREATE | O_EXCL    【如果已经存在，则失败】
	//os.O_WRONLY | os.O_CREATE    【如果已经存在，会覆盖写，不会清空原来的文件，而是从头直接覆盖写】
	//os.O_WRONLY | os.O_CREATE | os.O_APPEND    【如果已经存在，则在尾部添加写】

	fileObj, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		fmt.Println("Failed to open the file", err.Error())
		return err
	}

	if content != "" {
		if _, err := io.WriteString(fileObj, content); err == nil {
			fmt.Println("Successful appending to the file with os.OpenFile and io.WriteString.", content)
			return nil
		}
		return err
	}

	return nil
}

//读取指定字节
func ReadLimit(str string, len int64) string {
	reader := strings.NewReader(str)
	limitReader := &io.LimitedReader{R: reader, N: len}

	var res string
	for limitReader.N > 0 {
		tmp := make([]byte, 1)
		limitReader.Read(tmp)
		res += string(tmp)
	}
	return res
}

//读取整个文件
func ReadFile(dir string) string{
	data, err := ioutil.ReadFile(dir)
	if err != nil {
		panic(err)
		return ""
	}
	return string(data)
}

//按行读取文件
func ReadFileLine(dir string) map[int]string {
	file, err := os.OpenFile(dir, os.O_RDWR, 0666)
	util_err.Must(err)
	defer file.Close()

	/* stat, err := file.Stat()
	util_err.Must(err)
	size := stat.Size */

	buf := bufio.NewReader(file)
	res := make(map[int]string)
	i := 0
	for {
		line, _, err := buf.ReadLine()
		context := string(line)
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		res[i] = context
		i++
	}
	return res
}

func ReadJsonFile(dir string) string {
	jsonFile, err := os.Open(dir)
	util_err.Must(err)

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	return string(byteValue)
}
