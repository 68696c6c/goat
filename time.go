package goat

import "time"

const (
	ymdTemplate    = "2006-01-02"
	ymdhisTemplate = "2006-01-02 15:04:05"
)

func YMDHISStringToTime(str string) (time.Time, error) {
	return time.Parse(ymdhisTemplate, str)
}

func TimeToYMDHISString(t time.Time) string {
	return t.Format(ymdhisTemplate)
}

func YMDStringToTime(str string) (time.Time, error) {
	return time.Parse(ymdTemplate, str)
}

func TimeToYMDString(t time.Time) string {
	return t.Format(ymdTemplate)
}
