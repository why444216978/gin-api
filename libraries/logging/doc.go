// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package logging implements YMT logging standard. The design ideas are derived from Python's
logging library.


LogConfig


	type LogConfig struct {
		Path                string // default to ./
		File                string // default to log.log
		Rotate              bool   // whether to rotate file,default to false
		RotatingFileHandler string // SIZE/CRON/TIME
		RotateSize          int64
		RotateInterval      int  // unit: second
		Mode                int  // 0 text 1 json
		Level               int  // DEBUG 0,INFO 1,WARN 2, ERROR 3
		Debug               bool // whether to output to stdout, default to false
	}


Usage

    c := logging.LogConfig{Path: "./", File: "test.log", Mode: 1, Rotate: true, Debug: true}
    logger := logging.NewLogger(&c)
    header := logging.LogHeader{LogId: "abc"}
    logger.Error(&header, "123131")
    // 等待异步写入完成
    time.Sleep(100 * time.Millisecond)



*/
package logging
