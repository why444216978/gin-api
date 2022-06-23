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
	return rotatelogs.New(
		strings.Replace(filename, ".log", "", -1)+"-%Y%m%d%H.log", // 2022年1月1日12点 => filename-2022010112.log
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*7),  // log max save seven days
		rotatelogs.WithRotationTime(time.Hour), // rotate once an hour
	)
}
