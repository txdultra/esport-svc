package passport

import (
	"time"
	"utils"
)

type AccessToken struct {
	Id          int64
	Uid         int64
	AccessToken string `orm:"unique;size(32)"`
	ExpiresIn   int
	LastTime    time.Time
	LoginIp     int64
	App         int
}

func (self *AccessToken) ExpireDurition() time.Duration {
	duration := time.Second * time.Duration(self.ExpiresIn)
	expireTime := self.LastTime.Add(duration)
	dur := expireTime.Sub(time.Now())
	return dur
}

func (self *AccessToken) IpString() string {
	return utils.IntToIp(self.LoginIp)
}

func (self *AccessToken) TableName() string {
	return "common_access_tokens"
}

func (self *AccessToken) TableEngine() string {
	return "INNODB"
}
