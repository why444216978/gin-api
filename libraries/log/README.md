## Golang Logging
  - async logging
  - support formater
  - configurable

## QuickStart
```
package main

import (
	"time"

	"git.ymt360.com/go/logging"
)

func main() {
	msg := map[string]string{"AAAA": "BBBBB"}
	c := logging.LogConfig{Path: "./", File: "test.log", Mode: 1, Rotate: true, Debug: true}
	logger := logging.NewLogger(&c)
	header := logging.LogHeader{LogId: "abc"}
	logger.Error(&header, "123131")
	logger.Error(&header, msg)
    // 等待异步写入完成
	time.Sleep(100 * time.Millisecond)
	return
}
```
