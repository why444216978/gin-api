package slice

import (
	"errors"
	"fmt"
	"reflect"
)

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
