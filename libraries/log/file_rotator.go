package log

import (
	"fmt"
	"log"
	"os"
	"time"
)

// rotator 主要完成文件的自动切割
// 支持按大小，按整点时间，按定时切割
// 先Rename，再Reopen

const (
	SIZED_ROTATING_FILE_HANDLER = "SIZE" //Not Implement Yet
	CRON_ROTATING_FILE_HANDLER  = "CRON"
	TIMED_ROTATING_FILE_HANDLER = "TIME"
)

const (
	MIN_ROTATE_INTERVAL     = time.Minute          // one minute
	MAX_ROTATE_INTERVAL     = 365 * 24 * time.Hour // one year
	DEFAULT_ROTATE_INTERVAL = time.Hour            // one hour
)

type FileRotator interface {
	FilePath() string
	NextFileSuffix() string
	NextFilePath() string // NextFilePath = FilePath + NextFileSuffix

	Rotate() error
	NextRotate() //计算下一次Rotate的时机,仅用于TimedFileRotator
	Monitor() <-chan time.Time
}

func NewFileRotator(c *LogConfig) (logRotator FileRotator) {
	switch c.RotatingFileHandler {
	case SIZED_ROTATING_FILE_HANDLER:
		log.Fatal("unimplemented RotatingFileHandler:", c.RotatingFileHandler)
	case CRON_ROTATING_FILE_HANDLER:
		logRotator = NewCronFileRotator(c)
	case TIMED_ROTATING_FILE_HANDLER:
		logRotator = NewTimedFileRotator(c)
	default:
		log.Fatal("unexpected RotatingFileHandler:", c.RotatingFileHandler)
	}

	return
}

type CronFileRotator struct {
	LogFile

	date     time.Time     //下一次rotate时间
	interval time.Duration //rotate间隔,默认3600s
}

func NewCronFileRotator(c *LogConfig) *CronFileRotator {
	rotator := &CronFileRotator{
		LogFile: LogFile{
			path: c.Path,
			file: c.File,
		},
		interval: formatInterval(c.RotateInterval),
	}

	rotator.NextRotate()

	return rotator
}

func (cfr *CronFileRotator) NextFilePath() string {
	return cfr.FilePath() + "." + cfr.NextFileSuffix()
}

func (cfr *CronFileRotator) Rotate() (err error) {
	log.Println("Rotate")
	targetFilePath := cfr.NextFilePath()
	for i := 0; ; i++ {
		if _, err = os.Stat(targetFilePath); os.IsNotExist(err) {
			break
		} else {
			targetFilePath = fmt.Sprintf("%s.%d", targetFilePath, i)
			log.Println(i, targetFilePath)
		}
	}

	if err = os.Rename(cfr.FilePath(), targetFilePath); err != nil {
		log.Printf("[rename %s to %s] Rotator Rename error:%s\n", cfr.FilePath(), targetFilePath, err)
		return err
	}

	return nil
}

func (h *CronFileRotator) NextRotate() {
	var (
		previousDate time.Time = h.date
	)

	if previousDate.IsZero() {
		now := time.Now()
		year, month, day := now.Date()

		//first align to today's zero
		previousDate = time.Date(year, month, day, 0, 0, 0, 0, time.Local)
		previousDate = previousDate.Add(now.Sub(previousDate) / h.interval * h.interval)
	}

	h.date = previousDate.Add(h.interval)
}

func (h *CronFileRotator) Monitor() <-chan time.Time {
	return time.After(h.date.Sub(time.Now()))
}

func (h *CronFileRotator) NextFileSuffix() string {
	//machine-name.20160714
	//日志的后缀名为:切分动作发生时间减去一个切分时间单位
	var (
		format string
	)
	if h.interval >= 24*time.Hour {
		format = "20060102"
	} else if h.interval >= time.Hour {
		format = "2006010215"
	} else {
		format = "200601021504"
	}

	return h.date.Add(-1 * h.interval).Format(format)
}

//TimedFileRotator 是一种特殊的CronFileRotator
//TimedFileRotator 要求interval恰好为day/hour/minute
type TimedFileRotator struct {
	CronFileRotator
}

func NewTimedFileRotator(c *LogConfig) *TimedFileRotator {
	rotator := &TimedFileRotator{
		CronFileRotator: CronFileRotator{
			LogFile: LogFile{
				path: c.Path,
				file: c.File,
			},
			interval: alignInterval(c.RotateInterval),
		},
	}

	rotator.NextRotate()

	return rotator
}

func formatInterval(interval int64) time.Duration {
	if interval == 0 {
		return DEFAULT_ROTATE_INTERVAL
	}

	duration := time.Duration(interval) * time.Second

	if duration < MIN_ROTATE_INTERVAL {
		return MIN_ROTATE_INTERVAL
	}
	if duration > MAX_ROTATE_INTERVAL {
		return MAX_ROTATE_INTERVAL
	}

	//align to minute
	if duration%time.Minute != 0 {
		return duration / time.Minute * time.Minute
	}

	return duration
}

//alignInterval align interval to day/hour/minute
func alignInterval(interval int64) time.Duration {
	formatted := formatInterval(interval)
	if formatted < time.Hour {
		return time.Minute
	} else if formatted < 24*time.Hour {
		return time.Hour
	} else {
		return 24 * time.Hour
	}
}
