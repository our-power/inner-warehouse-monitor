package util

import (
	"time"
)

func TransformDuration2Seconds(begin time.Time, end time.Time) (duration int) {
	if begin.After(end) {
		return 1
	}
	return int(end.Sub(begin).Seconds())
}

func GetTodayTime(now time.Time) (todayTime int) {
	year, month, day := now.Date()
	todayStartPoint := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return TransformDuration2Seconds(todayStartPoint, now)
}
