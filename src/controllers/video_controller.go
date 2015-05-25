package controllers

import (
	"fmt"
	"libs"
	"libs/reptile"
	"libs/search"
	"libs/stat"
	"libs/vod"
	"outobjs"
	"strconv"
	"strings"
	"time"
	"utils"
)

// 视频 API
type VideoController struct {
	BaseController
}

func (c *VideoController) Prepare() {
	c.BaseController.Prepare()
}

func (c *VideoController) URLMapping() {
	c.Mapping("Modes", c.Modes)
	c.Mapping("Get", c.Get)
	c.Mapping("Flvs", c.Flvs)
	c.Mapping("Recommends", c.Recommends)
	c.Mapping("List", c.List)
	c.Mapping("ListByGames", c.ListByGames)
	c.Mapping("FlvsCallback", c.FlvsCallback)
	c.Mapping("Download", c.Download)
	c.Mapping("PlaylistVods", c.PlaylistVods)
}

// @Title 获取所有清晰度格式
// @Description 支持清晰度中英对应列表(返回数组)
// @Success 200 {object} outobjs.OutStreamMode
// @router /modes [get]
func (c *VideoController) Modes() {
	modes := []*outobjs.OutStreamMode{}
	for _, m := range reptile.ALL_VOD_STREAM_MODES {
		sm := &outobjs.OutStreamMode{
			Mode: m,
			Name: reptile.ConvertVodModeName(m),
		}
		modes = append(modes, sm)
	}
	c.Json(modes)
}

// @Title 按游戏分类视频列表
// @Description 视频列表(返回数组)
// @Param   access_token     path   string  false  "access_token"
// @Param   game_ids path    string  false  "游戏分类,(示例:1,2,3,全部则传空)"
// @Param   size     path    int  false  "分类列表大小,默认4"
// @Success 200 {object} outobjs.OutVideoInfoByGame
// @router /list/by_games [get]
func (c *VideoController) ListByGames() {
	query, _ := utils.UrlDecode(c.GetString("query"))
	game_ids := c.GetString("game_ids")
	size, _ := c.GetInt("size")
	//timestamp, _ := c.GetInt("t")
	match_mode := c.GetString("match_mode")

	if len(match_mode) == 0 {
		match_mode = "any"
	}

	if len(query) == 0 {
		match_mode = "all"
	}
	if size <= 0 {
		size = 4
	}

	//1分钟间隔缓存
	query_cache_key := fmt.Sprintf("front_fast_cache.vods_bygames.query:words_%s_gameids_%s_p_%d_s_%d_mode_%s_t_%s",
		query, game_ids, 1, size, match_mode)
	c_obj := utils.GetLocalFastExpriesTimePartCache(query_cache_key)
	if c_obj != nil {
		c.Json(c_obj)
		return
	}

	bas := libs.Bas{}
	allGames := bas.Games()
	fgameids := []uint64{}
	spgids := strings.Split(game_ids, ",")
	if len(game_ids) > 0 {
		for _, id := range spgids {
			_id, err := strconv.Atoi(id)
			if err == nil && _id > 0 {
				isadd := true
				for _, g := range allGames {
					if g.Id == _id && !g.Enabled {
						isadd = false
						break
					}
				}
				if isadd {
					fgameids = append(fgameids, uint64(_id))
				}
			}
		}
	} else {
		for _, g := range allGames {
			if g.Enabled {
				fgameids = append(fgameids, uint64(g.Id))
			}
		}
	}

	vod := &vod.Vods{}
	out_list := []*outobjs.OutVideoInfoByGame{}
	for _, game_id := range fgameids {
		filters := []search.SearchFilter{
			search.SearchFilter{
				Attr:    "gid",
				Values:  []uint64{game_id},
				Exclude: false,
			},
		}
		_, videos := vod.Query(query, 1, int(size), match_mode, filters, nil)
		vod_list := []*outobjs.OutVideoInfo{}
		for _, v := range videos {
			_v := outobjs.GetOutVideoInfo(v)
			vod_list = append(vod_list, _v)
		}
		out_game := outobjs.GetOutGameById(int(game_id))
		if out_game == nil {
			continue
		}
		out_list = append(out_list, &outobjs.OutVideoInfoByGame{
			Game: out_game,
			Vods: vod_list,
		})
	}
	utils.SetLocalFastExpriesTimePartCache(5*time.Minute, query_cache_key, out_list)
	c.Json(out_list)
}

// @Title 获取视频列表
// @Description 视频列表(返回数组)
// @Param   access_token     path   string  false  "access_token"
// @Param   query    path    string  false  "查询字符串"
// @Param   game_ids path    string  false  "游戏分类,(示例:1,2,3,全部则传空)"
// @Param   size     path    int  false  "一页行数(默认20)"
// @Param   page     path    int  false  "页数(默认1)"
// @Param   exclude_vids path   string  false  "忽略指定视频ids(示例:1,2,3)"
// @Param   t 		 path  	 int  false  "时间戳(防止服务器更新后客户端翻页重复),由服务器生成,首次调用会回传t参数时间戳,客户端每次提交时带上(首次或手动刷新时留空)"
// @Param   match_mode path  string   false  "搜索模式(all,any,phrase,boolean,extended,fullscan,extended2),默认为any"
// @Success 200 {object} outobjs.OutVideoPageList
// @router /list [get]
func (c *VideoController) List() {

	fmt.Println(c.Ctx.Input.IP())

	query, _ := utils.UrlDecode(c.GetString("query"))
	game_ids := c.GetString("game_ids")
	size, _ := c.GetInt("size")
	page, _ := c.GetInt("page")
	exclude_vid_str := c.GetString("exclude_vids")
	timestamp, _ := c.GetInt64("t")
	match_mode := c.GetString("match_mode")
	nocache, _ := c.GetBool("nocache")

	t := time.Now()
	if timestamp > 0 {
		t = time.Unix(timestamp, 0)
	}

	//1分钟间隔缓存 测试时关闭
	query_cache_key := fmt.Sprintf("front_fast_cache.vods.query:words_%s_gameids_%s_p_%d_s_%d_exv_%s_mode_%s",
		query, game_ids, page, size, exclude_vid_str, match_mode)
	if !nocache {
		c_obj := utils.GetLocalFastTimePartCache(t, query_cache_key, utils.CACHE_INTERVAL_TIME_TYPE_MINUTE)
		if c_obj != nil {
			c.Json(c_obj)
			return
		}
	}

	if len(match_mode) == 0 {
		match_mode = "any"
	}

	if len(query) == 0 {
		match_mode = "all"
	}

	fgameids := []uint64{}
	exclude_vids := []uint64{}
	spgids := strings.Split(game_ids, ",")
	exclude_vid_sps := strings.Split(exclude_vid_str, ",")
	for _, id := range spgids {
		_id, err := strconv.Atoi(id)
		if err == nil && _id > 0 {
			fgameids = append(fgameids, uint64(_id))
		}
	}
	for _, id := range exclude_vid_sps {
		_id, err := strconv.Atoi(id)
		if err == nil && _id > 0 {
			exclude_vids = append(exclude_vids, uint64(_id))
		}
	}

	var filters []search.SearchFilter
	if len(fgameids) > 0 {
		filters = append(filters, search.SearchFilter{
			Attr:    "gid",
			Values:  fgameids,
			Exclude: false,
		})
	}
	if len(exclude_vids) > 0 {
		filters = append(filters, search.SearchFilter{
			Attr:    "video_id",
			Values:  exclude_vids,
			Exclude: true,
		})
	}
	var filterRanges []search.FilterRangeInt
	filterRanges = append(filterRanges, search.FilterRangeInt{
		Attr:    "add_time",
		Min:     0,
		Max:     uint64(t.Unix()),
		Exclude: false,
	})

	if size <= 0 {
		size = 20
	}
	if page <= 0 {
		page = 1
	}
	vod := &vod.Vods{}
	total, videos := vod.Query(query, int(page), int(size), match_mode, filters, filterRanges)
	vod_infos := []*outobjs.OutVideoInfo{}
	for _, v := range videos {
		info := outobjs.GetOutVideoInfo(v)
		vod_infos = append(vod_infos, info)
	}

	out := outobjs.OutVideoPageList{
		Total:       total,
		Pages:       utils.TotalPages(int(total), int(size)),
		CurrentPage: int(page),
		Size:        int(size),
		Vods:        vod_infos,
		Time:        t.Unix(),
	}
	utils.SetLocalFastTimePartCache(t, query_cache_key, utils.CACHE_INTERVAL_TIME_TYPE_MINUTE, out)
	c.Json(out)
}

// @Title 获取视频信息
// @Description 单个视频信息
// @Param   vid     path    int  true  "视频id"
// @Success 200 {object} outobjs.OutVideoInfo
// @router /:id([0-9]+) [get]
func (c *VideoController) Get() {
	vid, err := c.GetInt64(":id")
	if err != nil {
		c.Json(libs.NewError("vod_parameter", "V1201", "参数错误", ""))
		return
	}
	vs := &vod.Vods{}
	video := vs.Get(vid, false)
	if video == nil {
		c.Json(libs.NewError("vod_not_exist", "V1202", "视频不存在", ""))
		return
	}
	out := *outobjs.GetOutVideoInfo(video)
	c.Json(out)
}

// @Title 获取视频播放真实地址
// @Description 真实播放地址flvs(返回数组),标识"m_"为手机用m3u8文件地址
// @Param   vid     path    int  true  "视频id"
// @Success 200 {object} outobjs.OutFlvs
// @router /flvs [get]
func (c *VideoController) Flvs() {
	vid, err := c.GetInt64("vid")
	if err != nil {
		c.Json(libs.NewError("vod_parameter", "V1101", "参数错误", ""))
		return
	}
	vs := &vod.Vods{}
	video := vs.Get(vid, true)
	if video == nil {
		c.Json(libs.NewError("vod_not_exist", "V1102", "视频不存在", ""))
		return
	}
	vpf := vs.GetPlayFlvs(vid, true)
	out_opts := c.getOutOpts(vpf)
	c.Json(out_opts)
}

func (c *VideoController) getOutOpts(vpf *vod.VideoPlayFlvs) []*outobjs.OutFlvs {
	out_opts := []*outobjs.OutFlvs{}
	if vpf != nil {
		for _, v := range vpf.OptFlvs {
			out_opt := &outobjs.OutVideoOpt{
				Flvs:    v.N,
				Size:    v.Size,
				Mode:    v.Mode,
				Seconds: v.Seconds,
			}
			out_flvs := []*outobjs.OutVideoFlv{}
			for _, fv := range v.Flvs {
				_f := &outobjs.OutVideoFlv{
					Url:     fv.Url,
					No:      fv.No,
					Size:    fv.Size,
					Seconds: fv.Seconds,
				}
				out_flvs = append(out_flvs, _f)
			}
			jf := &outobjs.OutFlvs{
				Opt:  out_opt,
				Flvs: out_flvs,
			}
			out_opts = append(out_opts, jf)
		}
	}
	return out_opts
}

// @Title 获取视频播放地址回调
// @Description 获取视频播放地址回调,最终返回flvs(返回数组),标识"m_"为手机用m3u8文件地址
// @Param   vid     path    int  true  "视频id"
// @Param   state     path    string  false  "状态码(首次提交时无须填写)"
// @Param   content     path    string  false  "内容"
// @Success 200 {object} outobjs.OutFlvs
// @router /flvs_callback [post]
func (c *VideoController) FlvsCallback() {
	vid, err := c.GetInt64("vid")
	if err != nil {
		c.Json(libs.NewError("vod_callback_parameter", "V1105", "参数错误", ""))
		return
	}
	vs := &vod.Vods{}
	video := vs.Get(vid, false)
	if video == nil {
		c.Json(libs.NewError("vod_callback_not_exist", "V1106", "视频不存在", ""))
		return
	}
	content, _ := utils.UrlDecode(c.GetString("content"))
	state := c.GetString("state")
	rep := reptile.CallbackReptileService(video.Source)
	if rep == nil {
		c.Json(libs.NewError("vod_callback_fail", "V1107", "没有播放流", ""))
		return
	}
	vss, url, cmd, err := rep.CallbackReptile(content, video.Url, state)
	if cmd == reptile.REP_CALLBACK_COMMAND_COMPLETED {
		vpf := vs.VodStreamToVideoPlayFlvs(video.Id, vss)
		if vpf == nil {
			c.Json(libs.NewError("vod_callback_fail", "V1109", "抓取失败2", ""))
			return
		}
		out_opts := c.getOutOpts(vpf)
		c.Json(out_opts)
	}
	out_err := ""
	if err != nil {
		out_err = err.Error()
	}
	c.Json(&outobjs.OutVodClientProxyRepCallback{
		ClientUrl: url,
		State:     cmd,
		Error:     out_err,
	})
}

// @Title 获取视频模块大眼睛推存列表
// @Description 推荐列表(返回数组,game_id不填为全局模式)
// @Param   game_id     path   int  true  "游戏id"
// @Success 200 {object} outobjs.OutVodRecommend
// @router /recommends [get]
func (c *VideoController) Recommends() {
	service := libs.NewRecommendService()
	gameid, _ := c.GetInt("game_id")
	category := vod.VOD_RECOMMEND_CATEGORTY_NAME
	if gameid > 0 {
		category = fmt.Sprintf(vod.VOD_RECOMMEND_CATEGORTY_GAME, gameid)
	}
	recommends := service.Gets(category)
	ens := []*outobjs.OutVodRecommend{}
	for _, recommend := range recommends {
		if recommend.Enabled {
			_o := &outobjs.OutVodRecommend{
				Id:           recommend.Id,
				RefId:        recommend.RefId,
				RefType:      recommend.RefType,
				Title:        recommend.Title,
				ImgUrl:       file_storage.GetFileUrl(recommend.Img),
				ImgId:        recommend.Img,
				DisplayOrder: recommend.DisplayOrder,
			}
			ens = append(ens, _o)
		}
	}
	c.Json(ens)
}

// @Title 下载视频
// @Description 下载时提供的最新地址(返回数组)
// @Param   vid     path   int  true  "视频id"
// @Param   stream_mode     path  string  true  "清晰度标识"
// @Success 200 {object} outobjs.OutVideoFlv
// @router /download [get]
func (c *VideoController) Download() {
	vid, err := c.GetInt64("vid")
	if err != nil {
		c.Json(libs.NewError("vod_parameter", "V1301", "参数错误", ""))
		return
	}
	stream_mode := c.GetString("stream_mode")
	mode := reptile.ConvertVodStreamMode(stream_mode)
	if mode == reptile.VOD_STREAM_MODE_UNDEFINED {
		c.Json(libs.NewError("vod_stream_mode_notexist", "V1303", "指定的清晰度格式不存在", ""))
		return
	}
	vs := &vod.Vods{}
	video := vs.Get(vid, true)
	if video == nil {
		c.Json(libs.NewError("vod_not_exist", "V1302", "视频不存在", ""))
		return
	}
	//opts := vs.GetOpts(vid, true)
	//for _, opt := range opts {
	//	if opt.Mode == mode {
	//		flvs := vs.GetFlvs(opt.Id)
	//		if len(flvs) == 0 {
	//			c.Json(libs.NewError("vod_flvs_not_exist", "V1304", "视频的文件不存在", ""))
	//		}
	//		go stat.GetCounter(vod.MOD_NAME).IncrCount(vid, 1, "download")
	//		c.Json(flvs)
	//	}
	//}
	vpf := vs.GetPlayFlvs(vid, true)
	if vpf == nil {
		c.Json(libs.NewError("vod_flvs_not_exist", "V1304", "视频的文件不存在", ""))
		return
	}
	for _, v := range vpf.OptFlvs {
		if v.Mode != mode {
			continue
		}
		out_flvs := []*outobjs.OutVideoFlv{}
		for _, fv := range v.Flvs {
			_f := &outobjs.OutVideoFlv{
				Url:     fv.Url,
				No:      fv.No,
				Size:    fv.Size,
				Seconds: fv.Seconds,
			}
			out_flvs = append(out_flvs, _f)
		}
		go stat.GetCounter(vod.MOD_NAME).DoC(vid, 1, "download")
		c.Json(out_flvs)
		return
	}

	c.Json(libs.NewError("vod_stream_mode_notexist", "V1305", "需要的清晰度格式文件不存在", ""))
}

// @Title 下载视频清晰度
// @Description 下载视频清晰度(返回数组)
// @Param   vid     path   int  true  "视频id"
// @Success 200 {object} outobjs.OutVideoDownClarity
// @router /download/clarities [get]
func (c *VideoController) DownloadClarities() {
	vid, err := c.GetInt64("vid")
	if err != nil {
		c.Json(libs.NewError("vod_parameter", "V1401", "参数错误", ""))
		return
	}

	//6小时隔缓存
	query_cache_key := fmt.Sprintf("front_fast_cache.down_clarities:vid_%d", vid)
	c_obj := utils.GetLocalFastExpriesTimePartCache(query_cache_key)
	if c_obj != nil {
		c.Json(c_obj)
		return
	}

	vs := &vod.Vods{}
	video := vs.Get(vid, true)
	if video == nil {
		c.Json(libs.NewError("vod_not_exist", "V1402", "视频不存在", ""))
		return
	}
	vpf := vs.GetPlayFlvs(vid, true)
	if vpf == nil {
		c.Json(libs.NewError("vod_flvs_not_exist", "V1403", "视频的文件不存在", ""))
		return
	}
	clars := []*outobjs.OutVideoDownClarity{}
	sps := []reptile.VOD_STREAM_MODE{reptile.VOD_STREAM_MODE_STANDARD_SP, reptile.VOD_STREAM_MODE_HIGH_SP, reptile.VOD_STREAM_MODE_SUPER_SP}
	for _, sp := range sps {
		for _, opt := range vpf.OptFlvs {
			if opt.Mode == sp {
				clars = append(clars, &outobjs.OutVideoDownClarity{reptile.ConvertVodModeName(opt.Mode), opt.Mode, opt.Size})
			}
		}
	}

	//	for _, opt := range vpf.OptFlvs {
	//		switch opt.Mode {
	//		case reptile.VOD_STREAM_MODE_STANDARD_SP:
	//			clars = append(clars, &outobjs.OutVideoDownClarity{reptile.ConvertVodModeName(reptile.VOD_STREAM_MODE_STANDARD_SP), reptile.VOD_STREAM_MODE_STANDARD_SP, opt.Size})
	//			break
	//		case reptile.VOD_STREAM_MODE_HIGH_SP:
	//			clars = append(clars, &outobjs.OutVideoDownClarity{reptile.ConvertVodModeName(reptile.VOD_STREAM_MODE_HIGH_SP), reptile.VOD_STREAM_MODE_HIGH_SP, opt.Size})
	//			break
	//		case reptile.VOD_STREAM_MODE_SUPER_SP:
	//			clars = append(clars, &outobjs.OutVideoDownClarity{reptile.ConvertVodModeName(reptile.VOD_STREAM_MODE_SUPER_SP), reptile.VOD_STREAM_MODE_SUPER_SP, opt.Size})
	//		}
	//	}
	//	sortutil.AscByField(clars, "Size")
	utils.SetLocalFastExpriesTimePartCache(6*time.Hour, query_cache_key, clars)
	c.Json(clars)
}

// @Title 专辑视频列表
// @Description 专辑视频列表
// @Param   pid     path   int  true  "专辑id"
// @Param   page     path   int  false  "页"
// @Param   size     path   int  false  "页数量"
// @Success 200 {object} outobjs.OutVideoPageList
// @router /playlist/vods [get]
func (c *VideoController) PlaylistVods() {
	pid, _ := c.GetInt64("pid")
	page, _ := c.GetInt("page")
	size, _ := c.GetInt("size")
	if pid <= 0 {
		c.Json(libs.NewError("vod_parameter", "V1501", "参数错误", ""))
		return
	}
	pls := &vod.Vods{}
	total, vods := pls.GetPlaylistVods(pid, page, size)
	outp := []*outobjs.OutVideoInfo{}
	for _, vod := range vods {
		outp = append(outp, outobjs.GetOutVideoInfo(vod))
	}
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	outpl := &outobjs.OutVideoPageList{
		CurrentPage: page,
		Total:       total,
		Pages:       utils.TotalPages(total, size),
		Size:        size,
		Vods:        outp,
	}
	c.Json(outpl)
}
