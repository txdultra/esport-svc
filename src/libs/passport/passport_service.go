package passport

import (
	"libs/passport/service"
	"utils"
)

var mp = NewMemberProvider()
var ats = NewAccessTokenService()

type PassportServiceImpl struct{}

func (p PassportServiceImpl) GetMember(uid int64) (r *service.Result_, err error) {
	if uid <= 0 {
		return &service.Result_{
			Code:    100,
			ErrorA1: "uid不能小于等于0",
			Member:  nil,
		}, nil
	}
	member := mp.Get(uid)
	if member == nil {
		return &service.Result_{
			Code:    101,
			ErrorA1: "用户不存在",
			Member:  nil,
		}, nil
	}
	return &service.Result_{
		Code:    0,
		ErrorA1: "",
		Member: &service.Member{
			Uid:           member.Uid,
			UserName:      member.UserName,
			NickName:      member.NickName,
			Email:         member.Email,
			CreateTime:    member.CreateTime,
			Avatar:        member.Avatar,
			PushId:        member.PushId,
			PushChannelId: member.PushChannelId,
			PushProxy:     int32(member.PushProxy),
			DeviceType:    string(member.DeviceType),
			Certified:     member.Certified,
		},
	}, nil
}

func (p PassportServiceImpl) GetMemberByAccessToken(accessToken string) (r *service.Result_, err error) {
	if len(accessToken) == 0 {
		return &service.Result_{
			Code:    103,
			ErrorA1: "access_token非法",
			Member:  nil,
		}, nil
	}
	accTkn := utils.StripSQLInjection(accessToken)
	if accessToken != accTkn {
		return &service.Result_{
			Code:    103,
			ErrorA1: "access_token非法",
			Member:  nil,
		}, nil
	}

	token, resultErr := ats.GetTokenObj(accessToken)
	if resultErr != nil {
		return &service.Result_{
			Code:    104,
			ErrorA1: resultErr.ErrorDescription,
			Member:  nil,
		}, nil
	}
	member := mp.Get(token.Uid)
	if member == nil {
		return &service.Result_{
			Code:    101,
			ErrorA1: "用户不存在",
			Member:  nil,
		}, nil
	}
	return &service.Result_{
		Code:    0,
		ErrorA1: "",
		Member: &service.Member{
			Uid:           member.Uid,
			UserName:      member.UserName,
			NickName:      member.NickName,
			Email:         member.Email,
			CreateTime:    member.CreateTime,
			Avatar:        member.Avatar,
			PushId:        member.PushId,
			PushChannelId: member.PushChannelId,
			PushProxy:     int32(member.PushProxy),
			DeviceType:    string(member.DeviceType),
			Certified:     member.Certified,
		},
	}, nil
}
