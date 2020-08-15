package file

import (
	"bufio"
	"fmt"
	"io"
	"os"

	util_err "gin-frame/libraries/util/error"
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

func ReadFile(dir string) map[int]string {
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
