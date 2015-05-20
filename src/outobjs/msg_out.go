package outobjs

import (
	"libs/message"
	"time"
)

type OutAtMsg struct {
	FromUid    int64            `json:"from_uid"`
	ToUid      int64            `json:"to_uid"`
	FromMember *OutMember       `json:"from_member"`
	MsgType    message.MSG_TYPE `json:"msg_type"`
	Text       string           `json:"text"`
	RefId      string           `json:"ref_id"`
	PostTime   time.Time        `json:"post_time"`
	FriendTime string           `json:"friend_time"`
	Comment    *OutComment      `json:"comment"`
}

type OutAtMsgPageList struct {
	Total       int         `json:"total"`
	TotalPage   int         `json:"pages"`
	CurrentPage int         `json:"current_page"`
	Size        int         `json:"size"`
	Time        int64       `json:"t"`
	Lists       []*OutAtMsg `json:"lists"`
}
