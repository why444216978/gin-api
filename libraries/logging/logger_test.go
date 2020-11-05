package logging

import (
	"fmt"
	"testing"
	"time"
)

type TestData struct {
	A int
}

func Test_SyncFormatterLogger(t *testing.T) {
	tmp := &TestData{A: 1}

	logConfig := &LogConfig{
		File:           "logger_test.log",
		Mode:           1,
		Debug:          true,
		AsyncFormatter: false,
	}
	Init(logConfig)
	fmt.Printf("it should be: %d\n", tmp.A)
	Debug(defaultLogHeader, tmp)
	tmp.A = 2
	time.Sleep(time.Second * 2)
}

func Test_AsyncFormatterLogger(t *testing.T) {
	tmp := &TestData{A: 1}
	logConfig := &LogConfig{
		File:           "logger_test.log",
		Mode:           1,
		Debug:          true,
		AsyncFormatter: true,
	}
	Init(logConfig)
	fmt.Printf("it should be: %d\n", tmp.A)
	Debug(defaultLogHeader, tmp)
	tmp.A = 2
	time.Sleep(time.Second * 2)
}
