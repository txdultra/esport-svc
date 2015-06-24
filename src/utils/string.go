package utils

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	//"time"
)

//随机字符串
func RandomStrings(length int) string {
	if length > 32 {
		panic("length not more than 32")
	}
	str := NewGuid()
	str = strings.Replace(str, "-", "", -1)
	str = strings.Replace(str, "_", "", -1)
	return SubstrByByte(str, length)
}

func NewGuid() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func Md5(str string) string {
	hx := md5.New()
	hx.Write([]byte(str))
	return hex.EncodeToString(hx.Sum(nil))
}

// 按字节截取字符串 utf-8不乱码
func SubstrByByte(str string, length int) string {
	bs := []byte(str)[:length]
	bl := 0
	for i := len(bs) - 1; i >= 0; i-- {
		switch {
		case bs[i] >= 0 && bs[i] <= 127:
			return string(bs[:i+1])
		case bs[i] >= 128 && bs[i] <= 191:
			bl++
		case bs[i] >= 192 && bs[i] <= 253:
			cl := 0
			switch {
			case bs[i]&252 == 252:
				cl = 6
			case bs[i]&248 == 248:
				cl = 5
			case bs[i]&240 == 240:
				cl = 4
			case bs[i]&224 == 224:
				cl = 3
			default:
				cl = 2
			}
			if bl+1 == cl {
				return string(bs[:i+cl])
			}
			return string(bs[:i])
		}
	}
	return ""
}

func IntToIp(ip int64) string {
	ox1 := ip / 16777216
	ip = ip % 16777216
	ox2 := ip / 65536
	ip = ip % 65536
	ox3 := ip / 256
	ip = ip % 256
	return fmt.Sprint(ox1, ".", ox2, ".", ox3, ".", ip)
}

func IpToInt(ip string) int64 {
	if len(ip) == 0 {
		return 0
	}
	exp := 3
	var iip float64 = 0
	for _, v := range strings.Split(ip, ".") {
		ix, _ := strconv.ParseFloat(v, 64)
		y := ix * math.Pow(256, float64(exp))
		iip = iip + y
		exp = exp - 1
	}
	return int64(iip)
}

func StripSQLInjection(sql string) string {
	reg1 := regexp.MustCompile("(\\%27)|(\\')|(\\-\\-)")
	reg2 := regexp.MustCompile("((\\%27)|(\\'))\\s*((\\%6F)|o|(\\%4F))((\\%72)|r|(\\%52))")
	reg3 := regexp.MustCompile("\\s+exec(\\s|\\+)+(s|x)p\\w+")
	rstr := reg1.ReplaceAllString(sql, "")
	rstr = reg2.ReplaceAllString(rstr, "")
	rstr = reg3.ReplaceAllString(rstr, "")
	return rstr
}

func UrlParamKV(url string) map[string]string {
	kvs := make(map[string]string)
	ps := strings.Split(url, "?")
	if len(ps) <= 1 {
		return kvs
	}
	return UrlParamMap(ps[1])
}

func UrlParamMap(parameter string) map[string]string {
	kvs := make(map[string]string)
	vss := strings.Split(parameter, "&")
	for _, v := range vss {
		kv := strings.Split(v, "=")
		if len(kv) == 2 {
			kvs[kv[0]] = kv[1]
		}
	}
	return kvs
}

func UrlEncode(urlstr string) string {
	if len(urlstr) == 0 {
		return ""
	}
	return url.QueryEscape(urlstr)
}

func UrlDecode(enstr string) (string, error) {
	if len(enstr) == 0 {
		return "", nil
	}
	return url.QueryUnescape(enstr)
}

func ExtractAts(text string) []string {
	pattern := `@[^@\s]*?[:：，,.。 ]`
	re := regexp.MustCompile(pattern)
	_text := text + " "
	tos := re.FindAllString(_text, -1)
	uss := []string{}
	for _, s := range tos {
		_s := s[1 : len(s)-1]
		uss = append(uss, StripSQLInjection(_s))
	}
	return uss
}

func StringReplace(s string, start int, end int, rep_char string) string {
	if len(s) == 0 || len(s) < start || len(s) < end {
		return s
	}
	fs := s[0:start]
	es := s[end:]
	ssf := func(n int, c string) string {
		_s := ""
		for i := 0; i < n; i++ {
			_s += c
		}
		return _s
	}
	newstr := fs + ssf(end-start, rep_char) + es
	return newstr
}

func ToBase64(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func ToWebBase64(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

func FromBase64(src string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(src)
}

func FromWebBase64(src string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(src)
}

func XmlEscape(input string) string {
	buf := bytes.NewBuffer(nil)
	xml.Escape(buf, []byte(input))
	return buf.String()
}

func CensorWords(source string) string {
	if len(source) == 0 {
		return source
	}
	_cp := source
	for c, _ := range censorWords {
		if strings.Contains(_cp, c) {
			rpl := strings.Repeat("*", len([]rune(c)))
			_cp = strings.Replace(_cp, c, rpl, -1)
		}
	}
	return _cp
}

func IsChineseChar(r rune) bool {
	if unicode.Is(unicode.Scripts["Han"], r) {
		return true
	}
	return false
}

func ReplaceRepeatString(sourceStr string, repStr string, n int, newStr string) string {
	patten := fmt.Sprintf("[%s]{%d,}", repStr, n)
	rep := regexp.MustCompile(patten)
	return rep.ReplaceAllString(sourceStr, newStr)
}
