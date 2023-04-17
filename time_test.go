package goat

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Time_YMDHISStringToTime(t *testing.T) {
	date, err := time.Parse("2006-01-02 15:04:05", "2011-01-19 12:30:00")
	require.Nil(t, err, "failed to create test date")
	result, err := YMDHISStringToTime("2011-01-19 12:30:00")
	require.Nil(t, err, "failed to parse test date")
	assert.Equal(t, date, result, "date not parsed correctly")
}

func Test_Time_TimeToYMDHISString(t *testing.T) {
	date, err := time.Parse("2006-01-02 15:04:05", "2015-05-23 10:45:30")
	require.Nil(t, err, "failed to create test date")
	result := TimeToYMDHISString(date)
	assert.Equal(t, "2015-05-23 10:45:30", result, "date not formatted correctly")
}

func Test_Time_YMDStringToTime(t *testing.T) {
	date, err := time.Parse("2006-01-02", "1999-08-07")
	require.Nil(t, err, "failed to create test date")
	result, err := YMDStringToTime("1999-08-07")
	require.Nil(t, err, "failed to parse test date")
	assert.Equal(t, date, result, "date not parsed correctly")
}

func Test_Time_TimeToYMDString(t *testing.T) {
	date, err := time.Parse("2006-01-02", "1988-11-01")
	require.Nil(t, err, "failed to create test date")
	result := TimeToYMDString(date)
	assert.Equal(t, "1988-11-01", result, "date not formatted correctly")
}

func Test_Time_TimeToPrettyString(t *testing.T) {
	date, err := time.Parse("2006-01-02", "2006-08-04")
	require.Nil(t, err, "failed to create test date")
	result := TimeToPrettyString(date)
	assert.Equal(t, "August 4th 2006", result, "date not formatted correctly")
}

func Test_Time_TimeToPrettyString_EndsWithOne(t *testing.T) {
	date, err := time.Parse("2006-01-02", "2000-01-01")
	require.Nil(t, err, "failed to create test date")
	result := TimeToPrettyString(date)
	assert.Equal(t, "January 1st 2000", result, "date not formatted correctly")
}

func Test_Time_TimeToPrettyString_EndsWithTwo(t *testing.T) {
	date, err := time.Parse("2006-01-02", "1966-04-02")
	require.Nil(t, err, "failed to create test date")
	result := TimeToPrettyString(date)
	assert.Equal(t, "April 2nd 1966", result, "date not formatted correctly")
}

func Test_Time_TimeToPrettyString_EndsWithThree(t *testing.T) {
	date, err := time.Parse("2006-01-02", "1993-02-03")
	require.Nil(t, err, "failed to create test date")
	result := TimeToPrettyString(date)
	assert.Equal(t, "February 3rd 1993", result, "date not formatted correctly")
}
