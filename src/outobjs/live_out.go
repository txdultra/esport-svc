package outobjs

import (
	"libs/lives"
	"libs/reptile"
	//"sort"
	"fmt"
	"time"
)

type OutPersonalLiveList struct {
	Total       int                `json:"total"`
	TotalPage   int                `json:"pages"`
	CurrentPage int                `json:"current_page"`
	Size        int                `json:"size"`
	Time        int64              `json:"t"`
	Lives       []*OutPersonalLive `json:"lives"`
}

type OutPersonalLive struct {
	Id        int64               `json:"live_id"`
	Title     string              `json:"title"`
	Des       string              `json:"description"`
	ImgId     int64               `json:"img_id"`
	ImgUrl    string              `json:"img_url"`
	Uid       int64               `json:"uid"`
	StreamUrl string              `json:"stream_url"`
	Status    reptile.LIVE_STATUS `json:"status"`
	SortField int                 `json:"sort_field"`
	Member    *OutMember          `json:"member"`
	Onlines   int                 `json:"onlines"`
	OfGames   []*OutGame          `json:"games"`
	RepMethod reptile.REP_METHOD  `json:"rep_method"`
	ShareUrl  string              `json:"share_url"`
}

func GetOutPersonalLive(per *lives.LivePerson) *OutPersonalLive {
	lps := &lives.LivePers{}
	lp_games := lps.GetOfGames(per.Id)
	ofGames := []*OutGame{}
	for _, ofgid := range lp_games {
		ofg := GetOutGameById(ofgid)
		if ofg != nil {
			ofGames = append(ofGames, ofg)
		}
	}

	oc := &lives.OnlineCounter{}
	return &OutPersonalLive{
		Id:        per.Id,
		Title:     per.Name,
		Des:       per.Des,
		ImgId:     per.Img,
		ImgUrl:    file.GetFileUrl(per.Img),
		Uid:       per.Uid,
		StreamUrl: per.StreamUrl,
		Status:    per.LiveStatus,
		SortField: int(per.LiveStatus),
		Member:    GetOutMember(per.Uid, 0),
		Onlines:   oc.GetChannelCounts(lives.LIVE_TYPE_PERSONAL, int(per.Id)),
		OfGames:   ofGames,
		RepMethod: reptile.LiveRepMethod(per.Rep), //per.RepMethod,
		ShareUrl:  getPLiveShareUrl(per.Id),
	}
}

func getPLiveShareUrl(id int64) string {
	return fmt.Sprintf("http://www.dianjingquan.cn/plive/%d", id)
}

type OutLiveChannel struct {
	Id       int64      `json:"id"`
	Title    string     `json:"title"`
	Uid      int64      `json:"uid"`
	ImgId    int64      `json:"img_id"`
	ImgUrl   string     `json:"img_url"`
	Member   *OutMember `json:"member"`
	Onlines  int        `json:"onlines"`
	ShareUrl string     `json:"share_url"`
}

type OutChannelStream struct {
	Id          int64              `json:"id"`
	Rep         string             `json:"rep"`
	StreamUrl   string             `json:"stream_url"`
	IsDefault   bool               `json:"is_def"`
	RepMethod   reptile.REP_METHOD `json:"rep_method"`
	LoadingInfo string             `json:"loading_info"`
}

type OutChannelV1 struct {
	Channel *OutLiveChannel     `json:"channel"`
	Streams []*OutChannelStream `json:"streams"`
}

/////////////////////////////////////////////////////////////////////////////////////

type OutProgramV1 struct {
	Date     time.Time        `json:"-"`
	Year     int              `json:"year"`
	YearStr  string           `json:"year_str"`
	Month    int              `json:"month"`
	MonthStr string           `json:"month_str"`
	Day      int              `json:"day"`
	DayStr   string           `json:"day_str"`
	Week     string           `json:"week"`
	Programs []*OutProgramObj `json:"programs"`
}

type OutProgramObj struct {
	Id               int64               `json:"id"`
	Title            string              `json:"title"`
	SubTitle         string              `json:"sub_title"`
	StartTime        time.Time           `json:"start_time"`
	EndTime          time.Time           `json:"end_time"`
	MatchId          int                 `json:"match_id"`
	Match            *OutMatch           `json:"match"`
	DefaultChannelId int64               `json:"def_channel_id"`
	Channels         []*OutLiveChannel   `json:"channels"`
	Game             *OutGame            `json:"game"`
	GameId           int                 `json:"game_id"`
	SubPrograms      []*OutSubProgramObj `json:"sub_programs"`
	IsLiving         bool                `json:"is_living"`
	IsExpired        bool                `json:"is_expired"`
	ImgId            int64               `json:"img_id"`
	ImgUrl           string              `json:"img_url"`
}

type OutSubProgramObj struct {
	Id              int64                           `json:"id"`
	ChannelId       int64                           `json:"channel_id"`
	GameId          int                             `json:"game_id"`
	ProgramId       int64                           `json:"program_id"`
	Game            *OutGame                        `json:"game"`
	Vs1             string                          `json:"vs1_name"`
	Vs1Uid          int64                           `json:"vs1_uid"`
	Vs1Img          int64                           `json:"vs1_img_id"`
	Vs1ImgUrl       string                          `json:"vs1_img_url"`
	Vs2             string                          `json:"vs2_name"`
	Vs2Uid          int64                           `json:"vs2_uid"`
	Vs2Img          int64                           `json:"vs2_img_id"`
	Vs2ImgUrl       string                          `json:"vs2_img_url"`
	Title           string                          `json:"title"`
	ImgId           int64                           `json:"img_id"`
	ImgUrl          string                          `json:"img_url"`
	Type            lives.LIVE_SUBPROGRAM_VIEW_TYPE `json:"t"`
	StartTime       time.Time                       `json:"start_time"`
	EndTime         time.Time                       `json:"end_time"`
	IsLiving        bool                            `json:"is_living"`
	SubScribeLocked bool                            `json:"subscribe_locked"`
	IsScribed       bool                            `json:"is_scribed"`
	IsExpired       bool                            `json:"is_expired"`
	BetCid          int64                           `json:"bet_cid"`
}

type OutUserProgramNotice struct {
	ProgramId   int64                     `json:"pid"`
	LastTime    time.Time                 `json:"last_time"`
	SubPrograms []OutUserSubProgramNotice `json:"subs"`
}

type OutUserSubProgramNotice struct {
	SubId           int64 `json:"sub_id"`
	IsLiving        bool  `json:"is_living"`
	SubScribeLocked bool  `json:"subscribe_locked"`
}

type OutClientProxyRepCallback struct {
	ClientUrl string `json:"client_url"`
	State     string `json:"state"`
	Error     string `json:"err"`
}

func GetOutLiveChannel(channel *lives.LiveChannel, uid int64) *OutLiveChannel {
	if channel == nil {
		return nil
	}
	oc := &lives.OnlineCounter{}
	return &OutLiveChannel{
		Id:       channel.Id,
		Title:    channel.Name,
		Uid:      channel.Uid,
		ImgId:    channel.Img,
		ImgUrl:   file.GetFileUrl(channel.Img),
		Member:   GetOutMember(channel.Uid, uid),
		Onlines:  oc.GetChannelCounts(lives.LIVE_TYPE_ORGANIZATION, int(channel.Id)),
		ShareUrl: getJLiveShareUrl(channel.Id),
	}
}

func getJLiveShareUrl(id int64) string {
	return fmt.Sprintf("http://www.dianjingquan.cn/jlive/%d", id)
}

func GetOutLiveChannelById(channelId int64, uid int64) *OutLiveChannel {
	live := &lives.LiveOrgs{}
	channel := live.GetChannel(channelId)
	if channel == nil {
		return nil
	}
	return GetOutLiveChannel(channel, uid)
}

//转换节目单输出对象
func ConvertOutProgramObj(program *lives.LiveProgram, uid int64) *OutProgramObj {
	if program == nil {
		return nil
	}
	spms := &lives.LiveSubPrograms{}
	subps := spms.Gets(program.Id)
	out_subpms := []*OutSubProgramObj{} //二级节目单s
	isPLiving := false
	for _, sub := range subps {
		out_subpm := ConvertOutSubProgramObj(sub, program, uid)
		if out_subpm == nil {
			continue
		}
		if !isPLiving && out_subpm.IsLiving { //判断一级菜单是否播放状态
			isPLiving = true
		}
		out_subpms = append(out_subpms, out_subpm)
	}
	out_channels := []*OutLiveChannel{}
	sps := &lives.LivePrograms{}
	channel_ids := sps.GetChannelIds(program.Id)
	for _, cid := range channel_ids {
		channel := GetOutLiveChannelById(cid, uid)
		if channel != nil {
			out_channels = append(out_channels, channel)
		}
	}
	out_pm := OutProgramObj{
		Id:               program.Id,
		Title:            program.Title,
		SubTitle:         program.SubTitle,
		StartTime:        program.StartTime,
		EndTime:          program.EndTime,
		MatchId:          program.MatchId,
		Match:            GetOutMatchById(program.MatchId),
		DefaultChannelId: program.DefaultChannelId,
		Channels:         out_channels,
		Game:             GetOutGameById(program.GameId),
		GameId:           program.GameId,
		SubPrograms:      out_subpms,
		IsLiving:         isPLiving,
		IsExpired:        time.Now().After(program.EndTime),
		ImgId:            program.Img,
		ImgUrl:           file.GetFileUrl(program.Img),
	}
	return &out_pm
}

func ConvertOutSubProgramObj(sub *lives.LiveSubProgram, program *lives.LiveProgram, uid int64) *OutSubProgramObj {
	if sub == nil || program == nil {
		return nil
	}
	spms := &lives.LiveSubPrograms{}
	isSLiving := spms.IsLiving(sub.Id) //判断二级菜单是否播放状态
	isSLocked := spms.IsLocked(sub.Id)
	isScribed := false
	if uid > 0 {
		pns := lives.NewProgramNoticeService()
		isScribed = pns.IsSubsuribed(uid, sub.Id)
	}
	out_subpm := &OutSubProgramObj{
		Id:              sub.Id,
		ChannelId:       program.DefaultChannelId,
		GameId:          sub.GameId,
		ProgramId:       program.Id,
		Game:            GetOutGameById(sub.GameId),
		Vs1:             sub.Vs1Name,
		Vs1Uid:          sub.Vs1Uid,
		Vs1Img:          sub.Vs1Img,
		Vs1ImgUrl:       file.GetFileUrl(sub.Vs1Img),
		Vs2:             sub.Vs2Name,
		Vs2Uid:          sub.Vs2Uid,
		Vs2Img:          sub.Vs2Img,
		Vs2ImgUrl:       file.GetFileUrl(sub.Vs2Img),
		Title:           sub.Title,
		ImgId:           sub.Img,
		ImgUrl:          file.GetFileUrl(sub.Img),
		Type:            sub.ViewType,
		StartTime:       sub.StartTime,
		EndTime:         sub.EndTime,
		IsLiving:        isSLiving,
		SubScribeLocked: isSLocked,
		IsExpired:       time.Now().After(sub.EndTime),
		IsScribed:       isScribed,
		BetCid:          sub.BetId,
	}
	return out_subpm
}

////////////////////////////////////////////////////////////////////////////////
// for admin
////////////////////////////////////////////////////////////////////////////////
type OutPersonalLiveListForAdmin struct {
	Total       int                        `json:"total"`
	TotalPage   int                        `json:"pages"`
	CurrentPage int                        `json:"current_page"`
	Size        int                        `json:"size"`
	Time        int64                      `json:"t"`
	Lives       []*OutPersonalLiveForAdmin `json:"lives"`
}

type OutPersonalLiveForAdmin struct {
	Id            int64               `json:"live_id"`
	Title         string              `json:"title"`
	Des           string              `json:"description"`
	ImgId         int64               `json:"img_id"`
	ImgUrl        string              `json:"img_url"`
	Uid           int64               `json:"uid"`
	StreamUrl     string              `json:"stream_url"`
	Status        reptile.LIVE_STATUS `json:"status"`
	Member        *OutMember          `json:"member"`
	Onlines       int                 `json:"onlines"`
	OfGames       []*OutGame          `json:"games"`
	RepMethod     reptile.REP_METHOD  `json:"rep_method"`
	ReptileUrl    string              `json:"rep_url"`
	ReptileDes    string              `json:"rep_desc"`
	PostTime      time.Time           `json:"post_time"`
	Enabled       bool                `json:"enabled"`
	ShowOnlineMin int                 `json:"show_online_min"`
	ShowOnlineMax int                 `json:"show_online_max"`
}

type OutLiveChannelForAdmin struct {
	Id     int64      `json:"id"`
	Title  string     `json:"title"`
	Uid    int64      `json:"uid"`
	Childs int        `json:"childs"`
	ImgId  int64      `json:"img_id"`
	ImgUrl string     `json:"img_url"`
	Member *OutMember `json:"member"`
}

func GetOutLiveChannelForAdmin(channel *lives.LiveChannel, uid int64) *OutLiveChannelForAdmin {
	return &OutLiveChannelForAdmin{
		Id:     channel.Id,
		Title:  channel.Name,
		Uid:    channel.Uid,
		Childs: channel.Childs,
		ImgId:  channel.Img,
		ImgUrl: file.GetFileUrl(channel.Img),
		Member: GetOutMember(channel.Uid, uid),
	}
}

type OutChannelStreamForAdmin struct {
	Id         int64               `json:"id"`
	ImgId      int64               `json:"img_id"`
	ImgUrl     string              `json:"img_url"`
	Rep        reptile.REP_SUPPORT `json:"rep"`
	ReptileUrl string              `json:"reptile_url"`
	ChannelId  int64               `json:"channel_id"`
	StreamUrl  string              `json:"stream_url"`
	IsDefault  bool                `json:"is_def"`
	RepMethod  reptile.REP_METHOD  `json:"rep_method"`
	AllowRep   bool                `json:"allow_rep"`
	Enabled    bool                `json:"enabled"`
}

func GetOutLiveChannelStreamForAdmin(stream *lives.LiveStream) *OutChannelStreamForAdmin {
	return &OutChannelStreamForAdmin{
		Id:         stream.Id,
		ImgId:      stream.Img,
		ImgUrl:     file.GetFileUrl(stream.Img),
		Rep:        stream.Rep,
		ReptileUrl: stream.ReptileUrl,
		ChannelId:  stream.ChannelId,
		StreamUrl:  stream.StreamUrl,
		IsDefault:  stream.Default,
		RepMethod:  reptile.LiveRepMethod(stream.Rep),
		AllowRep:   stream.AllowRep,
		Enabled:    stream.Enabled,
	}
}
