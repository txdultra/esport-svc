package version

import (
	"time"
)

type MOBILE_PLATFORM string

const (
	MOBILE_PLATFORM_ANDROID MOBILE_PLATFORM = "android"
	MOBILE_PLATFORM_APPLE   MOBILE_PLATFORM = "ios"
	MOBILE_PLATFORM_WPHONE  MOBILE_PLATFORM = "wphone"
)

type ClientVersion struct {
	Id           int64
	Version      float64
	Ver          string `orm:"column(ver)"`
	VerName      string
	Description  string
	PostTime     time.Time
	Platform     MOBILE_PLATFORM
	IsExpried    bool
	DownloadUrl  string `orm:"column(down_url)"`
	AllowVodDown bool   `orm:"column(allow_voddown)"`
}

func (self *ClientVersion) TableName() string {
	return "common_version"
}

func (self *ClientVersion) TableEngine() string {
	return "INNODB"
}

// 多字段唯一键
func (self *ClientVersion) TableUnique() [][]string {
	return [][]string{
		[]string{"version", "platform"},
	}
}
