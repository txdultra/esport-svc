package utils

import (
	"crypto/rand"
	"errors"
	"math/big"
	"time"
)

//异常处理
func Try(fun func(), handler func(interface{}), finally func()) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	defer finally()
	fun()
}

func WeekName(t time.Time, lanuage string) string {
	week := t.Weekday()
	switch week {
	case time.Monday:
		return "星期一"
	case time.Tuesday:
		return "星期二"
	case time.Wednesday:
		return "星期三"
	case time.Thursday:
		return "星期四"
	case time.Friday:
		return "星期五"
	case time.Saturday:
		return "星期六"
	case time.Sunday:
		return "星期天"
	default:
		return "none"
	}
}

func WeekName2(t time.Time, lanuage string) string {
	now := time.Now()
	if now.Year() == t.Year() && now.Month() == t.Month() && now.Day() == t.Day() {
		return "今天"
	}
	week := t.Weekday()
	switch week {
	case time.Monday:
		return "周一"
	case time.Tuesday:
		return "周二"
	case time.Wednesday:
		return "周三"
	case time.Thursday:
		return "周四"
	case time.Friday:
		return "周五"
	case time.Saturday:
		return "周六"
	case time.Sunday:
		return "周日"
	default:
		return "none"
	}
}

func TotalPages(total int, size int) int {
	if total == 0 || size == 0 {
		return 0
	}
	total_page := total / size
	if total%size > 0 {
		total_page++
	}
	if total%size == 0 {
		return total_page
	}
	return total_page
}

const (
	// Set of characters to use for generating random strings
	Alphabet     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	Numerals     = "1234567890"
	Alphanumeric = Alphabet + Numerals
	Ascii        = Alphanumeric + "~!@#$%^&*()-_+={}[]\\|<,>.?/\"';:`"
)

var MinMaxError = errors.New("Min cannot be greater than max.")

func IntRange(min, max int) (int, error) {
	var result int
	switch {
	case min > max:
		return result, MinMaxError
	case max == min:
		result = max
	case max > min:
		maxRand := max - min
		b, err := rand.Int(rand.Reader, big.NewInt(int64(maxRand)))
		if err != nil {
			return result, err
		}
		result = min + int(b.Int64())
	}
	return result, nil
}

// String returns a random string n characters long, composed of entities
// from charset.
func String(n int, charset string) (string, error) {
	randstr := make([]byte, n) // Random string to return
	charlen := big.NewInt(int64(len(charset)))
	for i := 0; i < n; i++ {
		b, err := rand.Int(rand.Reader, charlen)
		if err != nil {
			return "", err
		}
		r := int(b.Int64())
		randstr[i] = charset[r]
	}
	return string(randstr), nil
}

// StringRange returns a random string at least min and no more than max
// characters long, composed of entitites from charset.
func StringRange(min, max int, charset string) (string, error) {
	//
	// First determine the length of string to be generated
	//
	var err error      // Holds errors
	var strlen int     // Length of random string to generate
	var randstr string // Random string to return
	strlen, err = IntRange(min, max)
	if err != nil {
		return randstr, err
	}
	randstr, err = String(strlen, charset)
	if err != nil {
		return randstr, err
	}
	return randstr, nil
}

// AlphaRange returns a random alphanumeric string at least min and no more
// than max characters long.
func AlphaStringRange(min, max int) (string, error) {
	return StringRange(min, max, Alphanumeric)
}

// AlphaString returns a random alphanumeric string n characters long.
func AlphaString(n int) (string, error) {
	return String(n, Alphanumeric)
}

// ChoiceString returns a random selection from an array of strings.
func ChoiceString(choices []string) (string, error) {
	var winner string
	length := len(choices)
	i, err := IntRange(0, length)
	winner = choices[i]
	return winner, err
}

// ChoiceInt returns a random selection from an array of integers.
func ChoiceInt(choices []int) (int, error) {
	var winner int
	length := len(choices)
	i, err := IntRange(0, length)
	winner = choices[i]
	return winner, err
}

// A Choice contains a generic item and a weight controlling the frequency with
// which it will be selected.
type Choice struct {
	Weight int
	Item   interface{}
}

// WeightedChoice used weighted random selection to return one of the supplied
// choices. Weights of 0 are never selected. All other weight values are
// relative. E.g. if you have two choices both weighted 3, they will be
// returned equally often; and each will be returned 3 times as often as a
// choice weighted 1.
func WeightedChoice(choices []Choice) (Choice, error) {
	// Based on this algorithm:
	// http://eli.thegreenplace.net/2010/01/22/weighted-random-generation-in-python/
	var ret Choice
	sum := 0
	for _, c := range choices {
		sum += c.Weight
	}
	r, err := IntRange(0, sum)
	if err != nil {
		return ret, err
	}
	for _, c := range choices {
		r -= c.Weight
		if r < 0 {
			return c, nil
		}
	}
	err = errors.New("Internal error - code should not reach this point")
	return ret, err
}
