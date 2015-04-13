package outobjs

import (
	"time"
)

type OutFeedback struct {
	Id       int        `json:"id"`
	Uid      int64      `json:"uid"`
	Member   *OutMember `json:"member"`
	Category string     `json:"category"`
	PostTime time.Time  `json:"post_time"`
	Title    string     `json:"title"`
	Content  string     `json:"content"`
	Img      int64      `json:"img"`
	ImgUrl   string     `json:"img_url"`
	Contact  string     `json:"contact"`
	Source   string     `json:"source"`
}

type OutFeedbackPagedList struct {
	Total       int            `json:"total"`
	TotalPage   int            `json:"pages"`
	CurrentPage int            `json:"current_page"`
	Size        int            `json:"size"`
	Lists       []*OutFeedback `json:"lists"`
}
