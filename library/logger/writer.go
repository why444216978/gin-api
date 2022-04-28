package logger

import (
	"io"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

func RotateWriter(infoFile, errFile string) (infoWriter io.Writer, errWriter io.Writer, err error) {
	if infoWriter, err = rotateWriter(infoFile); err != nil {
		return
	}

	if errWriter, err = rotateWriter(errFile); err != nil {
		return
	}

	return
}

func rotateWriter(filename string) (io.Writer, error) {
	// 保存7天内的日志，每1小时(整点)分割一次日志
	return rotatelogs.New(
		strings.Replace(filename, ".log", "", -1)+"-%Y%m%d%H.log", // 没有使用go风格反人类的format格式
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour),
	)
}
