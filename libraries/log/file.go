package log

import (
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
	return filepath.Join(f.path, f.file)
	//return filepath.Join(f.path, f.file) + "." + misc.HostNamePrefix()
}

//open for write only
func (f *LogFile) OpenFile() error {
	if file, err := os.OpenFile(f.file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777); err != nil {
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
