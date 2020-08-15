package error

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
