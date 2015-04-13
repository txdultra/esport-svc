package utils

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func extractDomain(urlString string) string {
	if hasScheme, _ := regexp.MatchString(`https?://.*`, urlString); !hasScheme {
		urlString = "http://" + urlString
	}
	u, err := url.Parse(urlString)
	if err != nil {
		fmt.Println("Provided argument is not a valid URL!")
		os.Exit(1)
	}
	d := strings.Split(strings.Split(u.Host, ":")[0], ".")
	return strings.Join(d[len(d)-2:], ".")
}

func IsUrl(url string) bool {
	_url := strings.ToLower(url)
	if strings.HasPrefix(_url, "http://") || strings.HasPrefix(_url, "https://") {
		return true
	}
	return false
}
