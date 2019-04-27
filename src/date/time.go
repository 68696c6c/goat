package date

import (
	"fmt"
	"time"
)

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

func TimeToPrettyString(t time.Time) string {
	suffix := "th"
	switch t.Day() {
	case 1, 21, 31:
		suffix = "st"
	case 2, 22:
		suffix = "nd"
	case 3, 23:
		suffix = "rd"
	}
	m := t.Month().String()
	d := t.Day()
	y := t.Year()
	return fmt.Sprintf("%s %d%s %d", m, d, suffix, y)
}
