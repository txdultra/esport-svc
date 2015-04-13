package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

func FriendTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	nowTime := time.Now()
	sDur := nowTime.Sub(t)
	years := DurationYears(sDur)
	days := DurationDays(sDur)
	if years < 1 {
		if days < 30 {
			if days < 7 {
				if days < 1 {
					if sDur.Hours() < 1 {
						if sDur.Minutes() < 1 {
							if sDur.Seconds() < 1 {
								return fmt.Sprintf("刚刚")
							} else {
								return fmt.Sprintf("%d秒之前", int(sDur.Seconds()))
							}
						} else {
							return fmt.Sprintf("%d分钟之前", int(sDur.Minutes()))
						}
					} else {
						return fmt.Sprintf("%d小时之前", int(sDur.Hours()))
					}
				} else if days == 1 {
					return fmt.Sprintf("昨天")
				} else if days == 2 {
					return fmt.Sprintf("前天")
				} else {
					return fmt.Sprintf("%d天之前", days)
				}
			} else {
				return fmt.Sprintf("%d星期之前", days/7)
			}
		} else {
			return fmt.Sprintf("%d月之前", days/30)
		}
	} else {
		return fmt.Sprintf("%d年之前", years)
	}
}

func IntToCh(i int) string {
	str := strconv.FormatInt(int64(i), 10)
	result := ""
	for _, s := range str {
		ss := string(s)
		switch ss {
		case "0":
			result += "零"
		case "1":
			result += "一"
		case "2":
			result += "二"
		case "3":
			result += "三"
		case "4":
			result += "四"
		case "5":
			result += "五"
		case "6":
			result += "六"
		case "7":
			result += "七"
		case "8":
			result += "八"
		case "9":
			result += "九"
		}
	}
	return result
}

func MonthToCh(time time.Time) string {
	month := int(time.Month())
	switch month {
	case 1:
		return "1月"
	case 2:
		return "2月"
	case 3:
		return "3月"
	case 4:
		return "4月"
	case 5:
		return "5月"
	case 6:
		return "6月"
	case 7:
		return "7月"
	case 8:
		return "8月"
	case 9:
		return "9月"
	case 10:
		return "10月"
	case 11:
		return "11月"
	case 12:
		return "12月"
	}
	return ""
}

func DurationYears(dur time.Duration) int {
	hours := dur.Hours()
	hs := int(hours)
	if hs == 0 {
		return 0
	}
	return hs / 8760
}

func DurationDays(dur time.Duration) int {
	hours := dur.Hours()
	hs := int(hours)
	if hs == 0 {
		return 0
	}
	return hs / 24
}

func NanoToTime(nano int64) time.Time {
	_sec := nano / 1000000000
	_nano := nano % 1000000000
	return time.Unix(_sec, _nano)
}

func MsToTime(ms int64) time.Time {
	_sec := ms / 1000
	_nano := ms % 1000 * 1000000
	return time.Unix(_sec, _nano)
}

func IsMobile(mobile string) bool {
	if len(mobile) == 0 {
		return false
	}
	if matched, err := regexp.MatchString(`^1[3458][0-9]{9}$`, mobile); err != nil || !matched {
		return false
	}
	return true
}
