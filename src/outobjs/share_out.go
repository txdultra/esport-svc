package outobjs

import (
	"libs/share"
	"time"
)

type OutShare struct {
	Id               string              `json:"id"`
	ShareType        int                 `json:"share_type"`
	Source           string              `json:"source"`
	Geo              string              `json:"geo"`
	Text             string              `json:"text"`
	CreateTime       time.Time           `json:"create_time"`
	FriendTime       string              `json:"friend_time"`
	Ts               int64               `json:"ts"`
	RepostCount      int                 `json:"reposts_count"`
	CommentsCount    int                 `json:"comments_count"`
	NewCommentCounts int                 `json:"new_comments_count"`
	AttitudesCount   int                 `json:"attitudes_count"`
	Vods             []*OutShareVod      `json:"vods"`
	Pics             []*OutSharePic      `json:"pics"`
	Ats              []*OutScreenNameUid `json:"at_users"`
	Member           *OutMember          `json:"member"`
	Comments         []*OutShareComment  `json:"comments"`
	RefUids          []int64             `json:"ref_uids"`
}

type OutShareComment struct {
	Id         int64                    `json:"id"`
	Uid        int64                    `json:"uid"`
	UMember    *OutSimpleMember         `json:"u_member"`
	RUid       int64                    `json:"ruid"`
	RUMember   *OutSimpleMember         `json:"ru_member"`
	T          share.SHARE_COMMENT_TYPE `json:"t"`
	Content    string                   `json:"content"`
	FriendTime string                   `json:"friend_time"`
	Ts         int64                    `json:"ts"`
}

type OutShareNotice struct {
	Id         string                  `json:"id"`
	Sid        int64                   `json:"sid"`
	SContent   string                  `json:"s_content"`
	LUid       int64                   `json:"luid"`
	LUMember   *OutSimpleMember        `json:"lu_member"`
	RUid       int64                   `json:"ruid"`
	RUMember   *OutSimpleMember        `json:"ru_member"`
	Content    string                  `json:"content"`
	Pic        *OutSharePic            `json:"pic"`
	Ts         int64                   `json:"ts"`
	T          share.SHARE_NOTICE_TYPE `json:"t"`
	ST         int                     `json:"st"`
	FriendTime string                  `json:"friend_time"`
}

type OutShareVod struct {
	Id           int64      `json:"vod_id"`
	Title        string     `json:"title"`
	ThumbnailPic string     `json:"thumbnail_pic"`
	BmiddlePic   string     `json:"bmiddle_pic"`
	OriginalPic  string     `json:"original_pic"`
	Views        int        `json:"views"`
	Member       *OutMember `json:"member"`
}

type OutSharePic struct {
	Id           int64  `json:"img_id"`
	Title        string `json:"title"`
	ThumbnailPic string `json:"thumbnail_pic"`
	BmiddlePic   string `json:"bmiddle_pic"`
	OriginalPic  string `json:"original_pic"`
	Views        int    `json:"views"`
}

type OutSharePageList struct {
	Total       int         `json:"total"`
	TotalPage   int         `json:"pages"`
	CurrentPage int         `json:"current_page"`
	Size        int         `json:"size"`
	Time        int64       `json:"t"`
	Lists       []*OutShare `json:"lists"`
	LastNano    int64       `json:"last_nano"`
}

type OutShareCommentPageList struct {
	CurrentPage int                `json:"current_page"`
	Size        int                `json:"size"`
	Time        int64              `json:"t"`
	Lists       []*OutShareComment `json:"lists"`
}

type OutShareNoticePageList struct {
	Total       int               `json:"total"`
	TotalPage   int               `json:"pages"`
	CurrentPage int               `json:"current_page"`
	Size        int               `json:"size"`
	Time        int64             `json:"t"`
	Lists       []*OutShareNotice `json:"lists"`
}
