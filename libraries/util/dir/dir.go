package dir

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//获得当前绝对路径
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}

//检测并补全路径左边的反斜杠
func LeftAddPathPos(path string) string {
	if path[:0] != "/" {
		path = "/" + path
	}
	return path
}

//检测并补全路径右边的反斜杠
func RightAddPathPos(path string) string {
	if path[len(path)-1:len(path)] != "/" {
		path = path + "/"
	}
	return path
}

//根据当天日期和给定dir返回log文件名路径
func FileNameByDate(dir string) string {
	fileName := time.Now().Format("2006-01-02")
	dir = RightAddPathPos(dir)
	return dir + fileName + ".log"
}

//不存在则创建目录
func CreateDir(folderPath string) {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		os.MkdirAll(folderPath, 0777) //0777也可以os.ModePerm
		os.Chmod(folderPath, 0777)
	}
}

//根据当前日期，不存在则创建目录
func CreateDateDir(path string, prex string) string {
	folderName := time.Now().Format("20060102")
	if prex != "" {
		folderName = prex + folderName
	}
	folderPath := filepath.Join(path, folderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		os.MkdirAll(folderPath, 0777) //0777也可以os.ModePerm
		os.Chmod(folderPath, 0777)
	}
	return folderPath
}

func GetDateDir(path string) string {
	return path + time.Now().Format("20660102")
}

//根据当前小时创建目录和日志文件
func CreateHourLogFile(path string, prex string) string {
	folderName := time.Now().Format("20060102")
	if prex != "" {
		folderName = prex + folderName
	}
	hourDir := time.Now().Format("2006010215")
	folderPath := filepath.Join(path, folderName, hourDir)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		os.MkdirAll(folderPath, 0777) //0777也可以os.ModePerm
		os.Chmod(folderPath, 0777)
	}
	return folderPath
}
