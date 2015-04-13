package vod

import (
	"libs"
	"libs/reptile"
	"strings"
	"time"
	"utils"

	"labix.org/v2/mgo/bson"
)

type Video struct {
	Id             int64
	Title          string
	Url            string
	Img            int64
	Dkey           string
	LastUpdateTime time.Time `orm:"column(lr_time)"`
	PostTime       time.Time
	AddTime        time.Time
	Seconds        float32
	Mt             bool
	Source         reptile.VOD_SOURCE
	Uid            int64
	GameId         int  `orm:"column(gid)"`
	NoIndex        bool `orm:"column(no_idx)"`
}

func (self *Video) TableName() string {
	return "vod_videos"
}

func (self *Video) TableEngine() string {
	return "INNODB"
}

func (self *Video) Key() string {
	if len(self.Url) == 0 {
		return ""
	}
	s := strings.ToLower(self.Url)
	return utils.Md5(s)
}

func (self *Video) Reptile() error {
	vs := &Vods{}
	return vs.Reptile(self.Id)
}

func (self *Video) ExpriedAndReptile() {
	if self.Id <= 0 || len(self.Url) == 0 {
		return
	}
	dur := time.Now().Sub(self.LastUpdateTime)
	if dur.Seconds() <= 0 { //lasttime 大于 now
		return
	}
	videoId := self.Id
	sync_key := vod_chan_sync_key(self.Id)

	if dur.Seconds() > VOD_RE_REPTILE_TIMEOUT.Seconds() { // || self.Seconds == float32(0) {
		libs.ChanSync(sync_key, func() {
			vs := &Vods{}
			vs.Reptile(videoId)
		})
		return
	}

	taskLocker.Lock()
	defer taskLocker.Unlock()
	if _, ok := reptileTasks[self.Id]; ok {
		return
	}
	timer := time.AfterFunc(VOD_RE_REPTILE_TIMEOUT-dur, func() {
		libs.ChanSync(sync_key, func() {
			vs := &Vods{}
			vs.Reptile(videoId)
		})
	})
	reptileTasks[videoId] = timer
}

type VideoCount struct {
	VideoId   int64 `orm:"column(vid);pk" json:"video_id"`
	Views     int   `json:"views"`
	Comments  int   `json:"comments"`
	Favorites int   `json:"favs"`
	Dings     int   `json:"dings"`
	Cais      int   `json:"cais"`
	Downloads int   `json:"downloads"`
	Ex1       int   `json:"ex1"`
	Ex2       int   `json:"ex2"`
	Ex3       int   `json:"ex3"`
	Ex4       int   `json:"ex4"`
	Ex5       int   `json:"ex5"`
}

func (self *VideoCount) TableName() string {
	return "vod_counts"
}

func (self *VideoCount) TableEngine() string {
	return "INNODB"
}

//视频播放类型
type VideoOpt struct {
	N       int                     `bson:"n"`
	Size    int                     `bson:"size"`
	Mode    reptile.VOD_STREAM_MODE `bson:"mode"`
	Seconds float32                 `bson:"seconds"`
	Flvs    []VideoFlv              `bson:"flvs"`
}

type VideoPlayFlvs struct {
	ID      bson.ObjectId `bson:"_id"`
	VideoId int64         `bson:"video_id"`
	OptFlvs []VideoOpt    `bson:"opts"`
}

//视频地址
type VideoFlv struct {
	Url     string `bson:"url"`
	No      int16  `bson:"no"`
	Size    int    `bson:"size"`
	Seconds int    `bson:"seconds"`
}

//type VideoOpt struct {
//	Id      int64
//	VideoId int64 `orm:"column(vid)"`
//	Flvs    int
//	Size    int
//	Mode    reptile.VOD_STREAM_MODE
//	Seconds float32
//}

//func (self *VideoOpt) TableName() string {
//	return "vod_opts"
//}

//func (self *VideoOpt) TableEngine() string {
//	return "INNODB"
//}

//type VideoFlv struct {
//	Id      int64  `json:"-"`
//	VideoId int64  `orm:"column(vid)" json:"-"`
//	OpId    int64  `orm:"column(opid)" json:"-"`
//	Url     string `json:"url"`
//	No      int16  `json:"no"`
//	Size    int    `json:"size"`
//	Seconds int    `json:"seconds"`
//}

//func (self *VideoFlv) TableName() string {
//	return "vod_flvs"
//}

//func (self *VideoFlv) TableEngine() string {
//	return "INNODB"
//}

//用户空间
type VodUcenter struct {
	Id         int64
	Uid        int64
	Source     reptile.VOD_SOURCE
	SiteUrl    string
	LastTime   time.Time
	ScanAll    bool
	CreateTime time.Time `orm:"auto_now;type(datetime)"`
}

func (self *VodUcenter) TableName() string {
	return "vod_us"
}

func (self *VodUcenter) TableEngine() string {
	return "INNODB"
}

// 多字段唯一键
func (self *VodUcenter) TableUnique() [][]string {
	return [][]string{
		[]string{"uid", "source"},
	}
}

//视频专辑
type VideoPlaylist struct {
	Id       int64
	Title    string
	Des      string
	PostTime time.Time
	Vods     int
	Img      int
	Uid      int
}

func (self *VideoPlaylist) TableName() string {
	return "vod_playlist"
}

func (self *VideoPlaylist) TableEngine() string {
	return "INNODB"
}

//专辑视频
type VideoPlaylistVod struct {
	Id         int64
	PlaylistId int64 `orm:"column(pid)"`
	VideoId    int64 `orm:"column(vid)"`
	No         int
}

func (self *VideoPlaylistVod) TableName() string {
	return "vod_playlist_vods"
}

func (self *VideoPlaylistVod) TableEngine() string {
	return "INNODB"
}
