package outobjs

import (
	"libs/credits/client"
	"libs/passport"
	"time"
)

func GetOutSimpleMember(uid int64) *OutSimpleMember {
	if uid <= 0 {
		return nil
	}
	pst := passport.NewMemberProvider()
	member := pst.Get(uid)
	if member == nil {
		return nil
	}
	m := &OutSimpleMember{
		Uid:         member.Uid,
		NickName:    member.NickName,
		Avatar:      member.Avatar,
		AvatarUrl:   file.GetFileUrl(member.Avatar),
		IsCertified: member.Certified,
	}
	return m
}

func GetOutMember(uid int64, sourceUid int64) *OutMember {
	if uid <= 0 {
		return nil
	}
	pst := passport.NewMemberProvider()
	member := pst.Get(uid)
	state := pst.GetState(uid)
	profile := pst.GetProfile(uid)
	if member == nil {
		return nil
	}
	//用户关系
	var fship *OutUserSTRelation = nil
	if sourceUid > 0 {
		sMember := pst.Get(sourceUid)
		if sMember != nil {
			//不能为同一个人
			if uid != sourceUid {
				rt := friendship.Relation(sourceUid, uid)
				fship = &OutUserSTRelation{
					Source: &OutUserRelation{
						Uid:        sourceUid,
						NickName:   sMember.NickName,
						FollowedBy: rt.Source.FollowedBy,
						Following:  rt.Source.Following,
					},
					Target: &OutUserRelation{
						Uid:        uid,
						NickName:   member.NickName,
						FollowedBy: rt.Target.FollowedBy,
						Following:  rt.Target.Following,
					},
				}
			}
		}
	}
	c_client, c_transport, c_err := client.NewClient("")
	j_client, j_transport, j_err := client.NewClient("jings")
	var credits int64 = 0
	var jings int64 = 0
	if c_err == nil {
		defer func() {
			if c_transport != nil {
				c_transport.Close()
			}
		}()
		_credits, err := c_client.GetCredit(member.Uid)
		if err == nil {
			credits = _credits
		}
	}
	if j_err == nil {
		defer func() {
			if j_transport != nil {
				j_transport.Close()
			}
		}()
		_jings, err := j_client.GetCredit(member.Uid)
		if err == nil {
			jings = _jings
		}
	}

	m := &OutMember{
		Uid:             member.Uid,
		UserName:        member.UserName,
		NickName:        member.NickName,
		Email:           member.Email,
		Avatar:          member.Avatar,
		AvatarUrl:       file.GetFileUrl(member.Avatar),
		CreateTime:      time.Unix(member.CreateTime, 0),
		IsCertified:     member.Certified,
		CertifiedReason: member.CertifiedReason,
		FriendShip:      fship,
		Credits:         credits,
		Jings:           jings,
	}
	if state != nil {
		m.Fans = state.Fans
		m.Vods = state.Vods
		m.Notes = state.Notes
		m.Friends = friendship.FriendCounts(uid)
	}
	if profile != nil {
		m.Gender = profile.Gender
	}
	return m
}

type OutMemberPageList struct {
	Total       int          `json:"total"`
	TotalPage   int          `json:"pages"`
	CurrentPage int          `json:"current_page"`
	Size        int          `json:"size"`
	Time        int64        `json:"t"`
	Lists       []*OutMember `json:"lists"`
}

type OutMember struct {
	Uid             int64              `json:"uid"`
	UserName        string             `json:"user_name"`
	NickName        string             `json:"nick_name"`
	Email           string             `json:"email"`
	Avatar          int64              `json:"avatar_id"`
	AvatarUrl       string             `json:"avatar_url"`
	CreateTime      time.Time          `json:"create_time"`
	Vods            int                `json:"vods"`
	Fans            int                `json:"fans"`
	Notes           int                `json:"shares"`
	Friends         int                `json:"friends"`
	IsCertified     bool               `json:"is_certified"`
	CertifiedReason string             `json:"certified_reason"`
	FriendShip      *OutUserSTRelation `json:"friendship"`
	Credits         int64              `json:"credits"`
	Jings           int64              `json:"jings"`
	Gender          string             `json:"gender"`
}

type OutSimpleMember struct {
	Uid         int64  `json:"uid"`
	NickName    string `json:"nick_name"`
	Avatar      int64  `json:"avatar_id"`
	AvatarUrl   string `json:"avatar_url"`
	IsCertified bool   `json:"is_certified"`
}

type OutMemberInfo struct {
	*OutMember
	//state
	Logins        int       `json:"logins"`
	LastLoginIP   string    `json:"last_login_ip"`
	LastLoginTime time.Time `json:"last_login_time"`
	LongiTude     float32   `json:"last_login_longitude"`
	LatiTude      float32   `json:"last_login_latitude"`
	Vods          int       `json:"vods"`
	Fans          int       `json:"fans"`
	Notes         int       `json:"notes"`
	Count1        int       `json:"count1"`
	Count2        int       `json:"count2"`
	Count3        int       `json:"count3"`
	Count4        int       `json:"count4"`
	Count5        int       `json:"count5"`
	//profile
	BackgroundImg    int64  `json:"bg_img_id"`
	BackgroundImgUrl string `json:"bg_img_url"`
	Description      string `json:"description"`
	Gender           string `json:"gender"`
	RealName         string `json:"real_name"`
	BirthYear        int    `json:"birth_year"`
	BirthMonth       int    `json:"birth_month"`
	BirthDay         int    `json:"birth_day"`
	Mobile           string `json:"mobile"` //打***
	IDCard           string `json:"idcard"`
	QQ               string `json:"qq"`
	Alipay           string `json:"alipay"`
	Youku            string `json:"youku"`
	Bio              string `json:"bio"`
	Field1           string `json:"field1"`
	Field2           string `json:"field2"`
	Field3           string `json:"field3"`
	Field4           string `json:"field4"`
	Field5           string `json:"field5"`
}

type OutMemberGame struct {
	AddTime time.Time `json:"add_time"`
	Game    *OutGame  `json:"game"`
}

type OutMemberRole struct {
	Id       int64            `json:"id"`
	Uid      int64            `json:"uid"`
	Member   *OutSimpleMember `json:"member"`
	RoleId   int64            `json:"role_id"`
	Role     *OutRole         `json:"role"`
	PostTime time.Time        `json:"post_time"`
	Expries  int              `json:"expries"`
}

type OutRole struct {
	Id       int64  `json:"id"`
	RoleName string `json:"role_name"`
	Icon     string `json:"icon"`
	Enabled  bool   `json:"enabled"`
}
