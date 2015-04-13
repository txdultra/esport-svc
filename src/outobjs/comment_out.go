package outobjs

import (
	"libs/comment"
	"time"
	"utils"
)

type OutCommentPageList struct {
	Total       int           `json:"total"`
	TotalPage   int           `json:"pages"`
	CurrentPage int           `json:"current_page"`
	Size        int           `json:"size"`
	Time        int64         `json:"t"`
	Lists       []*OutComment `json:"lists"`
	LastNano    int64         `json:"last_nano"`
}

type OutScreenNameUid struct {
	ScreenName string `json:"screen_name"`
	Uid        int64  `json:"uid"`
}

type OutComment struct {
	Id         string              `json:"id"`
	Pid        string              `json:"parent_id"`
	FromId     int64               `json:"from_id"`
	Position   int                 `json:"position"`
	Title      string              `json:"title"`
	Text       string              `json:"text"`
	Reply      *OutReplyComment    `json:"reply"`
	PostTime   time.Time           `json:"post_time"`
	FriendTime string              `json:"friend_time"`
	AllowReply bool                `json:"allow_reply"`
	InVisible  bool                `json:"invisible"`
	IP         string              `json:"ip"`
	Area       string              `json:"area"`
	Member     *OutMember          `json:"member"`
	AtUsers    []*OutScreenNameUid `json:"at_users"`
	RefId      int64               `json:"ref_id"`
}

type OutReplyComment struct {
	Id       string     `bson:"id" json:"id"`
	FromId   int64      `bson:"from_id" json:"from_id"`
	Position int        `bson:"position" json:"position"`
	Text     string     `bson:"text" json:"text" `
	PostTime time.Time  `bson:"post_time" json:"post_time"`
	IP       string     `bson:"ip" json:"ip"`
	Member   *OutMember `bson:"member" json:"member"`
}

func GetOutComment(cmt *comment.Comment) *OutComment {
	if cmt == nil {
		return nil
	}
	var out_r *OutReplyComment
	if cmt.Reply != nil {
		out_r = &OutReplyComment{
			Id:       cmt.Reply.ID,
			FromId:   cmt.Reply.FromId,
			Position: cmt.Reply.Position,
			Text:     cmt.Reply.Text,
			PostTime: cmt.Reply.PostTime,
			IP:       cmt.Reply.IP,
			Member:   GetOutMember(cmt.Reply.FromId, 0),
		}
	}
	ats := []*OutScreenNameUid{}
	for k, v := range cmt.Ats {
		ats = append(ats, &OutScreenNameUid{
			ScreenName: k,
			Uid:        v,
		})
	}
	out_c := &OutComment{
		Id:         cmt.ID.Hex(),
		Pid:        cmt.Pid,
		FromId:     cmt.FromId,
		Position:   cmt.Position,
		Title:      cmt.Title,
		Text:       cmt.Text,
		Reply:      out_r,
		PostTime:   cmt.PostTime,
		FriendTime: utils.FriendTime(cmt.PostTime),
		AllowReply: cmt.AllowReply,
		InVisible:  cmt.InVisible,
		IP:         cmt.IP,
		Area:       cmt.Area,
		Member:     GetOutMember(cmt.FromId, 0),
		AtUsers:    ats,
		RefId:      cmt.RefId,
	}
	return out_c
}
