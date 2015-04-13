package outobjs

import (
	"time"
)

type OutCollectiblePageList struct {
	Total       int               `json:"total"`
	TotalPage   int               `json:"pages"`
	CurrentPage int               `json:"current_page"`
	Size        int               `json:"size"`
	Time        int64             `json:"t"`
	Lists       []*OutCollectible `json:"lists"`
}

type OutCollectible struct {
	Id             string      `json:"id"`
	RelId          string      `json:"rel_id"`
	RelType        string      `json:"rel_type"`
	PreviewContent string      `json:"preview_content"`
	PreviewImgId   int64       `json:"preview_img_id"`
	PreviewImgUrl  string      `json:"preview_img_url"`
	CreateTime     time.Time   `json:"create_time"`
	Vod            interface{} `json:"vod"`
	Org            interface{} `json:"org"`
	Per            interface{} `json:"per"`
}
