package conversion

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"reflect"
	"strings"

	util_err "gin-api/libraries/util/error"
)

//深拷贝转换
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

func JsonEncode(v interface{}) string {
	b, err := json.Marshal(v)
	util_err.Must(err)
	return string(b)
}

func MapToJsonInt(data map[int]interface{}) string {
	jsonStr, err := json.Marshal(data)
	util_err.Must(err)
	return string(jsonStr)
}

func MapToJson(data map[string]interface{}) string {
	jsonStr, err := json.Marshal(data)
	util_err.Must(err)
	return string(jsonStr)
}

//json转map数组
func JsonToMapArray(data string) []map[string]interface{} {
	var res []map[string]interface{}
	if data == "" {
		return res
	}
	err := json.Unmarshal([]byte(data), &res)
	util_err.Must(err)

	return res
}

//json转map
func JsonToMap(data string) map[string]interface{} {
	var res map[string]interface{}
	if data == "" {
		return res
	}
	err := json.Unmarshal([]byte(data), &res)
	util_err.Must(err)
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
