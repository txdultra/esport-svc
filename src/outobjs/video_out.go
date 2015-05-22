package outobjs

import (
	"fmt"
	"libs/reptile"
	"libs/vod"
	"time"
)

//输出对象格式
type OutFlvs struct {
	Opt  *OutVideoOpt   `json:"opt"`
	Flvs []*OutVideoFlv `json:"flvs"`
}

type OutVideoOpt struct {
	Flvs    int                     `json:"flvs"`
	Size    int                     `json:"size"`
	Mode    reptile.VOD_STREAM_MODE `json:"mode"`
	Seconds float32                 `json:"seconds"`
}

type OutVideoFlv struct {
	Url     string `json:"url"`
	No      int16  `json:"no"`
	Size    int    `json:"size"`
	Seconds int    `json:"seconds"`
}

type OutVideoDownClarity struct {
	Name string                  `json:"name"`
	Mode reptile.VOD_STREAM_MODE `json:"stream_mode"`
	Size int                     `json:"size"`
}

type OutVideoInfo struct {
	Id              int64           `json:"id"`
	Title           string          `json:"title"`
	ImgUrl          string          `json:"img_url"`
	ImgId           int64           `json:"img_id"`
	PostTime        time.Time       `json:"post_time"`
	Seconds         float32         `json:"seconds"`
	Uid             int64           `json:"uid"`
	GameId          int             `json:"game_id"`
	Game            *OutGame        `json:"game"`
	Counts          *vod.VideoCount `json:"counts"`
	Member          *OutMember      `json:"member"`
	ShareUrl        string          `json:"share_url"`
	UseCallbackPlay bool            `json:"use_callback_play"`
}

type OutVideoInfoByGame struct {
	Game *OutGame        `json:"game"`
	Vods []*OutVideoInfo `json:"vods"`
}

type OutVideoPageList struct {
	Total       int             `json:"total"`
	TotalPage   int             `json:"pages"`
	CurrentPage int             `json:"current_page"`
	Size        int             `json:"size"`
	Time        int64           `json:"t"`
	Vods        []*OutVideoInfo `json:"vods"`
}

type OutStreamMode struct {
	Mode reptile.VOD_STREAM_MODE `json:"mode"`
	Name string                  `json:"name"`
}

type OutVodRecommend struct {
	Id           int64  `json:"id"`
	RefId        int64  `json:"ref_id"`
	RefType      string `json:"ref_t"`
	Title        string `json:"title"`
	ImgUrl       string `json:"img_url"`
	ImgId        int64  `json:"img_id"`
	DisplayOrder int    `json:"no"`
}

type OutVodClientProxyRepCallback struct {
	ClientUrl string `json:"client_url"`
	State     string `json:"state"`
	Error     string `json:"err"`
}

var vs *vod.Vods = &vod.Vods{}

func GetOutVideoInfo(v *vod.Video) *OutVideoInfo {
	if v == nil {
		return nil
	}
	out := &OutVideoInfo{
		Id:              v.Id,
		ImgUrl:          file.GetFileUrl(v.Img),
		ImgId:           v.Img,
		Title:           v.Title,
		PostTime:        v.PostTime,
		Seconds:         v.Seconds,
		Uid:             v.Uid,
		GameId:          v.GameId,
		Game:            GetOutGameById(v.GameId),
		Counts:          vs.GetCount(v.Id),
		Member:          GetOutMember(v.Uid, 0),
		ShareUrl:        getVodShareUrl(v.Id),
		UseCallbackPlay: true,
	}
	return out
}

func getVodShareUrl(id int64) string {
	return fmt.Sprintf("http://www.dianjingquan.cn/vod/%d", id)
}

////////////////////////////////////////////////////////////////////////////////
// for admin
////////////////////////////////////////////////////////////////////////////////
type OutVodPageForAdmin struct {
	CurrentPage int               `json:"current_page"`
	Total       int               `json:"total"`
	Pages       int               `json:"pages"`
	Size        int               `json:"size"`
	Lists       []*OutVodForAdmin `json:"lists"`
}

type OutVodForAdmin struct {
	Id             int64              `json:"id"`
	Title          string             `json:"title"`
	Url            string             `json:"url"`
	Img            int64              `json:"img_id"`
	ImgUrl         string             `json:"img_url"`
	Dkey           string             `json:"key"`
	LastUpdateTime time.Time          `json:"last_time"`
	PostTime       time.Time          `json:"post_time"`
	AddTime        time.Time          `json:"add_time"`
	Seconds        float32            `json:"seconds"`
	Mt             bool               `json:"mt"`
	Source         reptile.VOD_SOURCE `json:"source"`
	Uid            int64              `json:"uid"`
	Member         *OutMember         `json:"member"`
	GameId         int                `json:"game_id"`
	Game           *OutGame           `json:"game"`
	NoIndex        bool               `json:"no_idx"`
}

type OutVodUcenterPageForAdmin struct {
	CurrentPage int                      `json:"current_page"`
	Total       int                      `json:"total"`
	Pages       int                      `json:"pages"`
	Size        int                      `json:"size"`
	Lists       []*OutVodUcenterForAdmin `json:"lists"`
}

type OutVodUcenterForAdmin struct {
	Id         int64              `json:"id"`
	Uid        int64              `json:"uid"`
	Member     *OutMember         `json:"member"`
	Source     reptile.VOD_SOURCE `json:"source"`
	SiteUrl    string             `json:"url"`
	LastTime   time.Time          `json:"last_rep_time"`
	ScanAll    bool               `json:"scan_all"`
	CreateTime time.Time          `json:"create_time"`
}

type OutVodPlaylistPagedListForAdmin struct {
	CurrentPage int               `json:"current_page"`
	Total       int               `json:"total"`
	Pages       int               `json:"pages"`
	Size        int               `json:"size"`
	Lists       []*OutVodPlaylist `json:"lists"`
}

type OutVodPlaylist struct {
	Id       int64            `json:"id"`
	Title    string           `json:"title"`
	Des      string           `json:"des"`
	PostTime time.Time        `json:"post_time"`
	Vods     int              `json:"vods"`
	Img      int64            `json:"img"`
	ImgUrl   string           `json:"img_url"`
	Uid      int64            `json:"uid"`
	Member   *OutSimpleMember `json:"member"`
}

func GetOutVodForAdmin(video *vod.Video) *OutVodForAdmin {
	return &OutVodForAdmin{
		Id:             video.Id,
		Title:          video.Title,
		Url:            video.Url,
		Img:            video.Img,
		ImgUrl:         file.GetFileUrl(video.Img),
		Dkey:           video.Dkey,
		LastUpdateTime: video.LastUpdateTime,
		PostTime:       video.PostTime,
		AddTime:        video.AddTime,
		Seconds:        video.Seconds,
		Mt:             video.Mt,
		Source:         video.Source,
		Uid:            video.Uid,
		Member:         GetOutMember(video.Uid, 0),
		GameId:         video.GameId,
		Game:           GetOutGameById(video.GameId),
		NoIndex:        video.NoIndex,
	}
}

func GetOutVodPlaylist(pl *vod.VideoPlaylist) *OutVodPlaylist {
	return &OutVodPlaylist{
		Id:       pl.Id,
		Title:    pl.Title,
		Des:      pl.Des,
		PostTime: pl.PostTime,
		Vods:     pl.Vods,
		Img:      pl.Img,
		ImgUrl:   file.GetFileUrl(pl.Img),
		Uid:      pl.Uid,
		Member:   GetOutSimpleMember(pl.Uid),
	}
}
