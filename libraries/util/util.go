package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func MapToJson(data map[string]interface{}) string {
	jsonStr, err := json.Marshal(data)
	Must(err)
	return string(jsonStr)
}

//json转map数组
func JsonToMapArray(data string) []map[string]interface{} {
	var res []map[string]interface{}
	if data == "" {
		return res
	}
	err := json.Unmarshal([]byte(data), &res)
	Must(err)

	return res
}

//json转map
func JsonToMap(data string) map[string]interface{} {
	var res map[string]interface{}
	if data == "" {
		return res
	}
	err := json.Unmarshal([]byte(data), &res)
	Must(err)
	return res
}

//url的path转文件名
func UriToFilePathByDate(uriPath string, dir string) string {
	pathArr := strings.Split(uriPath, "/")
	fileName := strings.Join(pathArr, "-")
	writePath := CreateDateDir(dir, "") //根据时间检测是否存在目录，不存在创建
	writePath = RightAddPathPos(writePath)
	fileName = path.Join(writePath, fileName[1:len(fileName)]+".log")
	return fileName
}

//uri转log路径
func UriToFilePath(uri string, dir string) string {
	pathArr := strings.Split(uri, "/")
	fileName := strings.Join(pathArr, "-")
	fileName = path.Join(dir, fileName[1:len(fileName)]+".log")
	if fileName[len(fileName)-1:len(fileName)] != "/" {
		fileName = fileName + "/"
	}
	return fileName
}

//uri转log文件名
func UriToFileName(uri string) string {
	pathArr := strings.Split(uri, "/")
	fileName := strings.Join(pathArr, "-")
	fileName = fileName + ".log"
	return fileName
}

//url转log文件名
func LogByUrl(fullUrl string) string {
	u, err := url.Parse(fullUrl)
	if err != nil {
		panic(err)
	}

	pathArr := strings.Split(u.Path, "/")
	fileName := strings.Join(pathArr, "-")
	writePath := "/tmp/logs/2020-01-12"
	fileName = path.Join(writePath, fileName[1:len(fileName)]+".log")

	return fileName
}

//uri的query转map
func ParseUriQueryToMap(query string) map[string]interface{} {
	queryMap := strings.Split(query, "&")
	res := make(map[string]interface{}, len(queryMap))
	if query == "" {
		return res
	}
	for _, item := range queryMap {
		itemMap := strings.Split(item, "=")
		res[itemMap[0]] = itemMap[1]
	}
	return res
}

//根据最大值生成随机整数
func RandomN(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func Must2(i interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return i
}

//获得某一天0点的时间戳
func GetDaysAgoZeroTime(day int) int64 {
	date := time.Now().AddDate(0, 0, day).Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", date)
	return t.Unix()
}

//时间戳转人可读
func TimeToHuman(target int) string {
	var res = ""
	if target == 0 {
		return res
	}

	t := int(time.Now().Unix()) - target
	data := [7]map[string]interface{}{
		{"key": 31536000, "value": "年"},
		{"key": 2592000, "value": "个月"},
		{"key": 604800, "value": "星期"},
		{"key": 86400, "value": "天"},
		{"key": 3600, "value": "小时"},
		{"key": 60, "value": "分钟"},
		{"key": 1, "value": "秒"},
	}
	for _, v := range data {
		var c = t / v["key"].(int)
		if 0 != c {
			res = strconv.Itoa(c) + v["value"].(string) + "前"
			break
		}
	}

	return res
}

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
func CreateDateDir(Path string, prex string) string {
	folderName := time.Now().Format("20060102")
	if prex != "" {
		folderName = prex + folderName
	}
	folderPath := filepath.Join(Path, folderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		os.MkdirAll(folderPath, 0777) //0777也可以os.ModePerm
		os.Chmod(folderPath, 0777)
	}
	return folderPath
}

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

func DeleteSliceByPos(slice interface{}, index int) (interface{}, error) {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return slice, errors.New("not slice")
	}
	if v.Len() == 0 || index < 0 || index > v.Len()-1 {
		return slice, errors.New("index error")
	}
	return reflect.AppendSlice(v.Slice(0, index), v.Slice(index+1, v.Len())).Interface(), nil
}
func InsertSliceByIndex(slice interface{}, index int, value interface{}) (interface{}, error) {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return slice, errors.New("not slice")
	}
	if index < 0 || index > v.Len() || reflect.TypeOf(slice).Elem() != reflect.TypeOf(value) {
		return slice, errors.New("index error")
	}
	if index == v.Len() {
		return reflect.Append(v, reflect.ValueOf(value)).Interface(), nil
	}
	v = reflect.AppendSlice(v.Slice(0, index+1), v.Slice(index, v.Len()))
	v.Index(index).Set(reflect.ValueOf(value))
	return v.Interface(), nil
}
func UpdateSliceByIndex(slice interface{}, index int, value interface{}) (interface{}, error) {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return slice, errors.New("not slice")
	}
	if index > v.Len()-1 || reflect.TypeOf(slice).Elem() != reflect.TypeOf(value) {
		return slice, errors.New("index error")
	}
	v.Index(index).Set(reflect.ValueOf(value))

	return v.Interface(), nil
}

//备忘：切片指定位置插入和删除原理
func sliceInsertAndDelete() {
	//insert
	data := []int{1, 2, 3, 4, 5}
	left := data[:3]
	right := data[3:]
	tmp := append([]int{}, left...)
	tmp = append(tmp, 0)
	res := append(tmp, right...)
	fmt.Println(res)

	//delete
	data = []int{1, 2, 3, 4, 5}
	left = data[:3]
	right = data[3+1:]
	res = append(left, right...)
	fmt.Println(res)
}

/*
	slice test code:
	i := 1
	a := []int{1, 2, 3}
	fmt.Println(a)
	res, err := util.InsertSliceByIndex(a , i, 9)
	util.Must(err)
	data := res.([]int)
	fmt.Println(data)

	res, err = util.DeleteSliceByPos(data, i)
	util.Must(err)
	data = res.([]int)
	fmt.Println(data)

	res, err = util.UpdateSliceByIndex(data, i , 6)
	util.Must(err)
	data = res.([]int)
	fmt.Println(data)
*/

//结构体转map
func StructToMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[strings.ToLower(obj1.Field(i).Name)] = obj2.Field(i).Interface()
	}
	return data
}

/*func StructToByte(tmp struct{}){
	tmp := &Test{Name: "why", Age: 34, Id: 1}
	length := unsafe.Sizeof(tmp)
	data := &SliceMock{
		addr: uintptr(unsafe.Pointer(tmp)),
		cap : int(length),
		len : int(length),
	}
	ret := *(*[]byte)(unsafe.Pointer(data))
}*/

//断言
func Assertion(data interface{}) interface{} {
	switch data.(type) {
	case string:
		return data.(string)
	case int:
		return data.(int)
	case int8:
		return data.(int8)
	case int32:
		return data.(int32)
	case int64:
		return data.(int64)
	case float32:
		return data.(float32)
	case float64:
		return data.(float64)
	default:
		return data
	}
	return nil
}

func BytesToString(b *[]byte) *string {
	s := bytes.NewBuffer(*b)
	r := s.String()
	return &r
}

func ExternalIP() (string, error) {
	iFaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iFace := range iFaces {
		if iFace.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iFace.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iFace.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

//获得本机名
func HostName() string {
	hostNamePrefix := ""
	host, err := os.Hostname()
	Must(err)
	if err == nil {
		parts := strings.SplitN(host, ".", 2)
		if len(parts) > 0 {
			hostNamePrefix = parts[0]
		}
	}
	return hostNamePrefix
}
