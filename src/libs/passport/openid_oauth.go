package passport

import (
	"time"
)

type OPENID_MARK byte

const (
	OPENID_MARK_QQ     OPENID_MARK = 1
	OPENID_MARK_WEIXIN OPENID_MARK = 2
)

const (
	OPENID_QQ_NAME     = "QQ"
	OPENID_WEIXIN_NAME = "WEIXIN"
)

func OpenIDMarkName(mark OPENID_MARK) string {
	switch mark {
	case OPENID_MARK_QQ:
		return OPENID_QQ_NAME
	case OPENID_MARK_WEIXIN:
		return OPENID_WEIXIN_NAME
	default:
		return "NON"
	}
}

type OpenIDOAuth struct {
	AuthId                int64 `orm:"auto;pk"`
	Uid                   int64
	AuthType              byte
	AuthIdentifier        string `orm:"size(32)"`
	AuthEmail             string `orm:"size(64)"`
	AuthEmailVerified     string `orm:"size(32);column(auth_emailverified)"`
	AuthPreferredUserName string `orm:"size(64);column(auth_preferredusername)"`
	AuthProviderName      string `orm:"size(16);column(auth_providername)"`
	AuthToken             string
	AuthDate              time.Time
	OpenIDUid             string `orm:"size(32);column(open_uid)"`
}

func (self *OpenIDOAuth) TableName() string {
	return "common_openid_auth"
}

func (self *OpenIDOAuth) TableEngine() string {
	return "INNODB"
}

func (self *OpenIDOAuth) TableUnique() [][]string {
	return [][]string{
		[]string{"AuthIdentifier", "AuthType"},
	}
}
