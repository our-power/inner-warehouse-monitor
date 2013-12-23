package util

import (
   "time" 
)

func FormatTime(date string, time_index int)(timestamp int64) {
    seconds := time_index*30
	hour := seconds/3600
	minute := seconds%3600/60
	second := seconds%60

	hour_str := string(hour)
	if hour < 10 {
		hour_str = "0" + hour_str
	}
	minute_str := string(minute)
	if minute < 10 {
		minute_str = "0" + minute_str
	}
	second_str := string(second)
	if second < 10 {
		second_str = "0" + second_str
	}

	time_str := bodyParts[0] + " " + hour_str + ":" + minute_str + ":" + second_str
	timestamp := time.Parse("20060102 15:04:05", time_str).Unix()
}
