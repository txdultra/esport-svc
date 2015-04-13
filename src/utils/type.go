package utils

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"time"
)

func IntToBytes(i int) []byte {
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.BigEndian, i)
	return buffer.Bytes()
}

func BytesToInt(b []byte) int {
	buffer := bytes.NewBuffer(b)
	var i int
	binary.Read(buffer, binary.BigEndian, &i)
	return i
}

func StrToTime(str string) (time.Time, error) {
	const TimeFormat = "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t, err := time.ParseInLocation(TimeFormat, str, loc)
	return t, err
}

func TimeMillisecond(t time.Time) int64 {
	return t.UnixNano() / 1000000
}

func StrToDuration(s string) time.Duration {
	dur, err := time.ParseDuration(s)
	if err != nil {
		panic("time duration format fail")
	}
	return dur
}

func IntToDuration(i int, t string) time.Duration {
	dur, _ := time.ParseDuration(strconv.Itoa(i) + t)
	return dur
}

func IsInArrayInt(src []int, tar []int) bool {
	for _, s := range src {
		for _, t := range tar {
			if s == t {
				return true
			}
		}
	}
	return false
}

func IsInArrayInt64(src []int64, tar []int64) bool {
	for _, s := range src {
		for _, t := range tar {
			if s == t {
				return true
			}
		}
	}
	return false
}

func IsZero(t time.Time) bool {
	et := time.Unix(0, 0)
	if et.Equal(t) {
		return true
	}
	return false
}
