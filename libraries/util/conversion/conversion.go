package conversion

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"

	"gin-frame/libraries/util/error"
)

func MapToJsonInt(data map[int]interface{}) string {
	jsonStr, err := json.Marshal(data)
	error.Must(err)
	return string(jsonStr)
}

func MapToJson(data map[string]interface{}) string {
	jsonStr, err := json.Marshal(data)
	error.Must(err)
	return string(jsonStr)
}

//json转map数组
func JsonToMapArray(data string) []map[string]interface{} {
	var res []map[string]interface{}
	if data == "" {
		return res
	}
	err := json.Unmarshal([]byte(data), &res)
	error.Must(err)

	return res
}

//json转map
func JsonToMap(data string) map[string]interface{} {
	var res map[string]interface{}
	if data == "" {
		return res
	}
	err := json.Unmarshal([]byte(data), &res)
	error.Must(err)
	return res
}

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
