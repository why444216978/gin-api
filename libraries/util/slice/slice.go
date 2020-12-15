package slice

import (
	"errors"
	"fmt"
	"reflect"
)

//删除切片指定位置元素
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

//在指定位置插入元素
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

//更新指定位置元素
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

func SliceContains(sl []interface{}, v interface{}) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func SliceContainsInt(sl []int, v int) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func SliceContainsInt64(sl []int64, v int64) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func SliceContainsString(sl []string, v string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

// SliceMerge merges interface slices to one slice.
func SliceMerge(slice1, slice2 []interface{}) (c []interface{}) {
	c = append(slice1, slice2...)
	return
}

func SliceMergeInt(slice1, slice2 []int) (c []int) {
	c = append(slice1, slice2...)
	return
}

func SliceMergeInt64(slice1, slice2 []int64) (c []int64) {
	c = append(slice1, slice2...)
	return
}

func SliceMergeString(slice1, slice2 []string) (c []string) {
	c = append(slice1, slice2...)
	return
}

func SliceUniqueInt64(s []int64) []int64 {
	size := len(s)
	if size == 0 {
		return []int64{}
	}

	m := make(map[int64]bool)
	for i := 0; i < size; i++ {
		m[s[i]] = true
	}

	realLen := len(m)
	ret := make([]int64, realLen)

	idx := 0
	for key := range m {
		ret[idx] = key
		idx++
	}

	return ret
}

func SliceUniqueInt(s []int) []int {
	size := len(s)
	if size == 0 {
		return []int{}
	}

	m := make(map[int]bool)
	for i := 0; i < size; i++ {
		m[s[i]] = true
	}

	realLen := len(m)
	ret := make([]int, realLen)

	idx := 0
	for key := range m {
		ret[idx] = key
		idx++
	}

	return ret
}

func SliceUniqueString(s []string) []string {
	size := len(s)
	if size == 0 {
		return []string{}
	}

	m := make(map[string]bool)
	for i := 0; i < size; i++ {
		m[s[i]] = true
	}

	realLen := len(m)
	ret := make([]string, realLen)

	idx := 0
	for key := range m {
		ret[idx] = key
		idx++
	}

	return ret
}

func SliceSumInt64(intslice []int64) (sum int64) {
	for _, v := range intslice {
		sum += v
	}
	return
}

func SliceSumInt(intslice []int) (sum int) {
	for _, v := range intslice {
		sum += v
	}
	return
}

func SliceSumFloat64(intslice []float64) (sum float64) {
	for _, v := range intslice {
		sum += v
	}
	return
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
	res, err := slice.InsertSliceByIndex(a , i, 9)
	if err != nil{
		panic(err)
	}
	data := res.([]int)
	fmt.Println(data)

	res, err = slice.DeleteSliceByPos(data, i)
	if err != nil{
		panic(err)
	}
	data = res.([]int)
	fmt.Println(data)

	res, err = slice.UpdateSliceByIndex(data, i , 6)
	if err != nil{
		panic(err)
	}
	data = res.([]int)
	fmt.Println(data)
*/
