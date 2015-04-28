package outobjs

import (
	"libs"
	"libs/version"
	"time"
)

func GetOutGame(game *libs.Game) *OutGame {
	if game == nil {
		return nil
	}
	return &OutGame{
		Id:      game.Id,
		Name:    game.Name,
		En:      game.En,
		ImgId:   game.Img,
		ImgUrl:  file.GetFileUrl(game.Img),
		Enabled: game.Enabled,
	}
}

func GetOutGameById(id int) *OutGame {
	bas := &libs.Bas{}
	game := bas.GetGame(id)
	if game == nil {
		return nil
	}
	return GetOutGame(game)
}

func GetOutMatchById(id int) *OutMatch {
	if id <= 0 {
		return nil
	}
	bas := &libs.Bas{}
	match := bas.Match(id)
	if match == nil {
		return nil
	}
	return &OutMatch{
		Id:       match.Id,
		Name:     match.Name,
		SubTitle: match.SubTitle,
		En:       match.En,
		ImgId:    match.Img,
		ImgUrl:   file.GetFileUrl(match.Img),
		Des1:     match.Des1,
		Des2:     match.Des2,
		Des3:     match.Des3,
	}
}

func GetOutMatch(match *libs.Match) *OutMatch {
	if match == nil {
		return nil
	}
	return &OutMatch{
		Id:       match.Id,
		Name:     match.Name,
		SubTitle: match.SubTitle,
		En:       match.En,
		ImgId:    match.Img,
		ImgUrl:   file.GetFileUrl(match.Img),
		Des1:     match.Des1,
		Des2:     match.Des2,
		Des3:     match.Des3,
	}
}

type OutGame struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	En      string `json:"en"`
	ImgId   int64  `json:"img_id"`
	ImgUrl  string `json:"img_url"`
	Enabled bool   `json:"enabled"`
}

type OutMatch struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	SubTitle string `json:"sub_title"`
	En       string `json:"en"`
	ImgId    int64  `json:"img_id"`
	ImgUrl   string `json:"img_url"`
	Des1     string `json:"des1"`
	Des2     string `json:"des2"`
	Des3     string `json:"des3"`
}

type OutVersion struct {
	Ver         string                  `json:"ver"`
	VerName     string                  `json:"ver_name"`
	Description string                  `json:"desc"`
	Platform    version.MOBILE_PLATFORM `json:"plat"`
	IsExpried   bool                    `json:"is_expries"`
	DownloadUrl string                  `json:"down_url"`
	NewVersion  *OutNewVersion          `json:"new_version"`
}

type OutVersionForAdmin struct {
	Id          int64                   `json:"id"`
	Version     float64                 `json:"version"`
	Ver         string                  `json:"ver"`
	VerName     string                  `json:"ver_name"`
	Description string                  `json:"desc"`
	PostTime    time.Time               `json:"post_time"`
	Platform    version.MOBILE_PLATFORM `json:"plat"`
	IsExpried   bool                    `json:"is_expries"`
	DownloadUrl string                  `json:"down_url"`
}

type OutNewVersion struct {
	Ver         string                  `json:"ver"`
	VerName     string                  `json:"ver_name"`
	Description string                  `json:"desc"`
	Platform    version.MOBILE_PLATFORM `json:"plat"`
	IsExpried   bool                    `json:"is_expries"`
	DownloadUrl string                  `json:"down_url"`
}

type OutApiModHost struct {
	ModName string `json:"mod"`
	BaseUrl string `json:"base_url"`
	Version string `json:"v"`
}

type OutSmiley struct {
	Ext          string `json:"ext"`
	Code         string `json:"code"`
	Bins         string `json:"bins"`
	Points       int    `json:"points"`
	Category     string `json:"category"`
	DisplayOrder int    `json:"display_order"`
}

type OutFileUrl struct {
	FileId  int64  `json:"file_id"`
	FileUrl string `json:"file_url"`
}

type OutPicture struct {
	Id           int64  `json:"img_id"`
	Title        string `json:"title"`
	ThumbnailPic string `json:"thumbnail_pic"`
	BmiddlePic   string `json:"bmiddle_pic"`
	OriginalPic  string `json:"original_pic"`
	Views        int    `json:"views"`
}

////////////////////////////////////////////////////////////////////////////////
// for admin
////////////////////////////////////////////////////////////////////////////////
type OutGameForAdmin struct {
	Id           int       `json:"id"`
	Name         string    `json:"name"`
	En           string    `json:"en"`
	Img          int64     `json:"img_id"`
	ImgUrl       string    `json:"img_url"`
	Enabled      bool      `json:"enabled"`
	PostTime     time.Time `json:"post_time"`
	DisplayOrder int       `json:"display_order"`
}

type OutMatchForAdmin struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	SubTitle string `json:"sub_title"`
	En       string `json:"en"`
	ImgId    int64  `json:"img_id"`
	ImgUrl   string `json:"img_url"`
	Des1     string `json:"des1"`
	Des2     string `json:"des2"`
	Des3     string `json:"des3"`
}

type OutRecommendForAdmin struct {
	Id           int64      `json:"id"`
	RefId        int64      `json:"ref_id"`
	RefType      string     `json:"ref_t"`
	Title        string     `json:"title"`
	ImgUrl       string     `json:"img_url"`
	ImgId        int64      `json:"img_id"`
	Categroy     string     `json:"category"`
	Enabled      bool       `json:"enabled"`
	DisplayOrder int        `json:"no"`
	PostTime     time.Time  `json:"post_time"`
	PostMember   *OutMember `json:"member"`
}
