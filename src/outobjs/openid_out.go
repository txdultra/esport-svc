package outobjs

import (
	"libs/passport"
	"time"
)

type OutAccessToken struct {
	AccessToken string   `json:"access_token"`
	ExpiresIn   int      `json:"expires_in"`
	Uid         int64    `json:"uid"`
	UnProcs     []string `json:"un_procs"`
}

func GetAccessToken(access *passport.AccessToken) *OutAccessToken {
	ms := passport.NewMemberProvider()
	unprocs := ms.CheckMemberUnCompletedProcs(access.Uid, "v1.1") //非版本号
	//计算多久过期
	duration := time.Second * time.Duration(access.ExpiresIn)
	expireTime := access.LastTime.Add(duration)
	dur := expireTime.Sub(time.Now())

	return &OutAccessToken{
		AccessToken: access.AccessToken,
		ExpiresIn:   int(dur.Seconds()),
		Uid:         access.Uid,
		UnProcs:     unprocs,
	}
}
