package log

type LogConfig struct {
	Path                string `default:"Users/why/Desktop/go"` // default to ./
	File                string `default:"default.log"`          // default to log.log
	Rotate              bool   `default:"false"`                // whether to rotate file,default to false
	RotatingFileHandler string `default:"TIME"`                 // SIZE/CRON/TIME
	RotateSize          int64
	RotateInterval      int64 `default:"3600"`  // unit: second
	Mode                int   `default:"1"`     // 0 text 1 json
	Level               int   `default:"1"`     // DEBUG 0,INFO 1,WARN 2, ERROR 3
	Debug               bool  `default:"false"` // whether to output to stdout, default to false
	AsyncFormatter      bool  `default:"false"` // true: async formatter, false:sync formatter
}
