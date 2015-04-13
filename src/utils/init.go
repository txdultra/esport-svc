package utils

import (
	//"code.google.com/p/mahonia"
	"io/ioutil"
	"path/filepath"
	"strings"
)

//敏感字
var censorWords map[string]bool = make(map[string]bool)

func init() {
	init_censor_words()
}

func init_censor_words() {
	path := filepath.Join("data", "censor_words.txt")
	data, err := ioutil.ReadFile(path)
	if err == nil {
		words := strings.Split(string(data), "\n")
		for _, word := range words {
			_w := strings.Trim(word, " ")
			if len(_w) == 0 {
				continue
			}
			_w = strings.Replace(_w, "\n", "", -1)
			_w = strings.Replace(_w, "\r", "", -1)
			if _, ok := censorWords[_w]; ok {
				continue
			}
			censorWords[_w] = true
		}
	}
}
