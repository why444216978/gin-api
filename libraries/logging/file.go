package logging

import (
	"gin-api/libraries/util/sys"
	"io"
	"os"
	"path/filepath"
)

type LogFile struct {
	path string //路径
	file string //文件名

	logFile *os.File
}

func (f *LogFile) FilePath() string {
	return filepath.Join(f.path, f.file) + "." + sys.HostName()
}

//open for write only
func (f *LogFile) OpenFile() error {
	if file, err := os.OpenFile(f.FilePath(), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644); err != nil {
		return err
	} else {
		f.logFile = file
		return nil
	}
}

func (f *LogFile) Close() error {
	if f.logFile == nil {
		return nil
	}
	return f.logFile.Close()
}

func (f *LogFile) Writer() io.Writer {
	return f.logFile
}
