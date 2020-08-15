package time

import (
	"strconv"
	"time"
)

//获得某一天0点的时间戳
func GetDaysAgoZeroTime(day int) int64 {
	date := time.Now().AddDate(0, 0, day).Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", date)
	return t.Unix()
}

//时间戳转人可读
func TimeToHuman(target int) string {
	var res = ""
	if target == 0 {
		return res
	}

	t := int(time.Now().Unix()) - target
	data := [7]map[string]interface{}{
		{"key": 31536000, "value": "年"},
		{"key": 2592000, "value": "个月"},
		{"key": 604800, "value": "星期"},
		{"key": 86400, "value": "天"},
		{"key": 3600, "value": "小时"},
		{"key": 60, "value": "分钟"},
		{"key": 1, "value": "秒"},
	}
	for _, v := range data {
		var c = t / v["key"].(int)
		if 0 != c {
			res = strconv.Itoa(c) + v["value"].(string) + "前"
			break
		}
	}

	return res
}
