package admincp

import (
	"controllers"
	"encoding/json"
	"fmt"
	"libs"
	"libs/reptile"
	"libs/vod"
	"net/url"
	"outobjs"
	"strconv"
	"strings"
	"time"
	"utils"
)

// 视频管理 API
type VodCPController struct {
	AdminController
}

func (c *VodCPController) Prepare() {
	c.AdminController.Prepare()
}

func (c *VodCPController) parseVodUrl(srcUrl string) string {
	_url, err := url.Parse(srcUrl)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s://%s%s", _url.Scheme, _url.Host, _url.Path)
}

// @Title 添加新视频
// @Description 添加新视频
// @Param   url   path	string true  "视频url"
// @Param   game_id   path	int true  "所属游戏"
// @Param   img_id   path	int true  "图片id"
// @Param   title   path	string true  "标题"
// @Param   of_uid   path	int true  "所属播主"
// @Param   reping   path	bool true  "立即抓取"
// @Param   nosearch   path	bool true  "不被搜索到"
// @Success 200 {object} libs.Error
// @router /add [post]
func (c *VodCPController) Add() {
	url, _ := utils.UrlDecode(c.GetString("url"))
	game_id, _ := c.GetInt64("game_id")
	img_id, _ := c.GetInt64("img_id")
	title, _ := utils.UrlDecode(c.GetString("title"))
	of_uid, _ := c.GetInt64("of_uid")
	reping, _ := c.GetBool("reping")
	nosearch, _ := c.GetBool("nosearch")
	if len(url) == 0 {
		c.Json(libs.NewError("admincp_vod_add_fail", "GM010_001", "url不能为空", ""))
		return
	}
	if game_id <= 0 {
		c.Json(libs.NewError("admincp_vod_add_fail", "GM010_002", "必须设置所属游戏", ""))
		return
	}
	if len(title) == 0 {
		c.Json(libs.NewError("admincp_vod_add_fail", "GM010_003", "必须设置标题", ""))
		return
	}
	if of_uid <= 0 {
		c.Json(libs.NewError("admincp_vod_add_fail", "GM010_004", "必须设置所属主播", ""))
		return
	}

	source := reptile.GetVodSource(url)
	if source == reptile.VOD_SOURCE_NONE {
		c.Json(libs.NewError("admincp_vod_add_fail", "GM010_005", "抓取地址未被支持", ""))
		return
	}

	_url := c.parseVodUrl(url)
	if len(_url) == 0 {
		c.Json(libs.NewError("admincp_vod_add_fail", "GM010_005", "抓取地址格式错误", ""))
		return
	}

	video := &vod.Video{}
	video.Title = title
	video.Url = _url
	video.Img = img_id
	video.PostTime = time.Now()
	video.Mt = true
	video.Source = source
	video.Uid = of_uid
	video.GameId = int(game_id)
	video.NoIndex = nosearch
	vods := &vod.Vods{}
	id, err := vods.Create(video, reping)
	if id > 0 {
		c.Json(libs.NewError("admincp_vod_add_succ", controllers.RESPONSE_SUCCESS, "新视频添加成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_vod_add_fail", "GM010_006", "添加失败:"+err.Error(), ""))
}

// @Title 修改视频
// @Description 修改视频
// @Param   id   path	int true  "视频编号"
// @Param   url   path	string true  "视频url(自动抓取视频不允许修改)"
// @Param   game_id   path	int true  "所属游戏"
// @Param   img_id   path	int true  "图片id"
// @Param   title   path	string true  "标题"
// @Param   of_uid   path	int true  "所属播主"
// @Param   nosearch   path	bool true  "不被搜索到"
// @Success 200 {object} libs.Error
// @router /update [post]
func (c *VodCPController) Update() {
	id, _ := c.GetInt64("id")
	url, _ := utils.UrlDecode(c.GetString("url"))
	game_id, _ := c.GetInt64("game_id")
	img_id, _ := c.GetInt64("img_id")
	title, _ := utils.UrlDecode(c.GetString("title"))
	of_uid, _ := c.GetInt64("of_uid")
	nosearch, _ := c.GetBool("nosearch")
	if id <= 0 {
		c.Json(libs.NewError("admincp_vod_update_fail", "GM010_018", "id不能小于0", ""))
		return
	}
	if len(url) == 0 {
		c.Json(libs.NewError("admincp_vod_update_fail", "GM010_010", "url不能为空", ""))
		return
	}
	//if game_id <= 0 {
	//	c.Json(libs.NewError("admincp_vod_update_fail", "GM010_011", "必须设置所属游戏", ""))
	//  return
	//}
	if len(title) == 0 {
		c.Json(libs.NewError("admincp_vod_update_fail", "GM010_012", "必须设置标题", ""))
		return
	}
	if of_uid <= 0 {
		c.Json(libs.NewError("admincp_vod_update_fail", "GM010_013", "必须设置所属主播", ""))
		return
	}
	source := reptile.GetVodSource(url)
	if source == reptile.VOD_SOURCE_NONE {
		c.Json(libs.NewError("admincp_vod_update_fail", "GM010_014", "抓取地址未被支持", ""))
		return
	}
	vods := &vod.Vods{}
	old_v := vods.Get(id, false)
	if old_v == nil {
		c.Json(libs.NewError("admincp_vod_update_fail", "GM010_015", "原视频不存在", ""))
		return
	}
	if !old_v.Mt && old_v.Url != url {
		c.Json(libs.NewError("admincp_vod_update_fail", "GM010_016", "不得修改自动抓取视频", ""))
		return
	}
	if old_v.Mt {
		_url := c.parseVodUrl(url)
		if len(_url) == 0 {
			c.Json(libs.NewError("admincp_vod_add_fail", "GM010_005", "抓取地址格式错误", ""))
			return
		}
		old_v.Url = _url
		old_v.Source = source
	}
	old_v.GameId = int(game_id)
	old_v.Img = img_id
	old_v.Title = title
	old_v.Uid = of_uid
	old_v.NoIndex = nosearch
	err := vods.Update(old_v)
	if err == nil {
		c.Json(libs.NewError("admincp_vod_update_succ", controllers.RESPONSE_SUCCESS, "视频更新成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_vod_update_fail", "GM010_017", "其他错误:"+err.Error(), ""))
}

// @Title 视频
// @Description 视频
// @Param   id   path	int true  "视频id"
// @Success 200 {object} outobjs.OutVodForAdmin
// @router /get [get]
func (c *VodCPController) Get() {
	id, _ := c.GetInt64("id")
	if id <= 0 {
		c.Json(libs.NewError("admincp_vod_get_fail", "GM010_071", "id不能小于0", ""))
		return
	}
	vods := &vod.Vods{}
	v := vods.Get(id, false)
	if v == nil {
		c.Json(libs.NewError("admincp_vod_get_fail", "GM010_072", "视频不存在", ""))
		return
	}
	c.Json(*outobjs.GetOutVodForAdmin(v))
}

// @Title 视频列表
// @Description 视频列表
// @Param   query   path	int true  "标题关键字"
// @Param   game_id  path	int true  "所属游戏(0未分游戏视频)"
// @Param   page   path	int true  "page"
// @Param   size   path	int true  "size"
// @Param   of_uid   path	int true  "所属播主"
// @Param   unarchived   path	bool true  "是否未归档"
// @Success 200 {object} outobjs.OutVodPageForAdmin
// @router /list [get]
func (c *VodCPController) List() {
	query, _ := utils.UrlDecode(c.GetString("query"))
	game_id, _ := c.GetInt64("game_id")
	page, _ := c.GetInt64("page")
	size, _ := c.GetInt64("size")
	of_uid, _ := c.GetInt64("of_uid")
	unarchived, _ := c.GetBool("unarchived")

	if size <= 0 {
		size = 20
	}

	params := make(map[string]string)
	if len(query) > 0 {
		params["title__icontains"] = query
	}
	if game_id > 0 {
		params["gid"] = strconv.Itoa(int(game_id))
	}
	if of_uid > 0 {
		params["uid"] = strconv.FormatInt(of_uid, 10)
	}
	if unarchived {
		params["gid"] = "0"
	}
	vods := &vod.Vods{}
	total, lists := vods.DbQuery(params, int(page), int(size))
	outvs := []*outobjs.OutVodForAdmin{}
	for _, v := range lists {
		outvs = append(outvs, outobjs.GetOutVodForAdmin(v))
	}
	out := outobjs.OutVodPageForAdmin{
		CurrentPage: int(page),
		Total:       total,
		Pages:       utils.TotalPages(total, int(size)),
		Size:        int(size),
		Lists:       outvs,
	}
	c.Json(out)
}

// @Title 批量更新视频游戏归类
// @Description 批量更新视频游戏归类
// @Param   input   path	int true  "json字符串"
// @Success 200
// @router /update/ofgame_batch [post]
func (c *VodCPController) UpdateOfGameBatch() {
	input, _ := utils.UrlDecode(c.GetString("input"))
	type InputParam struct {
		VideoId  int64 `json:"vid"`
		OfGameId int   `json:"game_id"`
	}
	vps := []*InputParam{}
	err := json.Unmarshal([]byte(input), &vps)
	if err != nil {
		c.Json(libs.NewError("admincp_vod_update_batch_fail", "GM010_030", "输入json格式错误", ""))
		return
	}
	vods := &vod.Vods{}
	for _, vp := range vps {
		video := vods.Get(vp.VideoId, false)
		if video != nil {
			video.GameId = vp.OfGameId
			vods.Update(video)
		}
	}
	c.Json(libs.NewError("admincp_vod_update_batch_succ", controllers.RESPONSE_SUCCESS, "视频集合更新成功", ""))
}

// @Title 添加自动抓取主播空间地址
// @Description 添加自动抓取主播空间地址
// @Param   uid  path	int true  "主播id"
// @Param   uc_url path	string true  "空间地址"
// @Success 200
// @router /uc/add [post]
func (c *VodCPController) UserVodCenterReptileAdd() {
	uid, _ := c.GetInt64("uid")
	ucurl, _ := utils.UrlDecode(c.GetString("uc_url"))

	if len(ucurl) == 0 {
		c.Json(libs.NewError("admincp_vod_uc_add_fail", "GM010_040", "url不能为空", ""))
		return
	}
	if uid <= 0 {
		c.Json(libs.NewError("admincp_vod_uc_add_fail", "GM010_041", "必须设置所属主播", ""))
		return
	}
	source := reptile.GetUcSource(ucurl)
	if source == reptile.VOD_SOURCE_NONE {
		c.Json(libs.NewError("admincp_vod_uc_add_fail", "GM010_042", "抓取地址未被支持", ""))
		return
	}

	uc := vod.VodUcenter{}
	uc.Uid = uid
	uc.SiteUrl = ucurl
	uc.Source = source
	ucp := &vod.VodUcenterReptile{}
	_, err := ucp.Create(uc)
	if err == nil {
		c.Json(libs.NewError("admincp_vod_uc_add_succ", controllers.RESPONSE_SUCCESS, "已成功添加主播抓取空间地址", ""))
		return
	}
	c.Json(libs.NewError("admincp_vod_uc_add_fail", "GM010_043", "其他错误:"+err.Error(), ""))
}

// @Title 删除自动抓取主播空间地址
// @Description 添加自动抓取主播空间地址
// @Param   id  path int true  "id"
// @Success 200
// @router /uc/del [delete]
func (c *VodCPController) UserVodCenterReptileDel() {
	id, _ := c.GetInt64("id")
	if id <= 0 {
		c.Json(libs.NewError("admincp_vod_uc_del_fail", "GM010_050", "必须提供删除id", ""))
		return
	}
	ucp := &vod.VodUcenterReptile{}
	err := ucp.Delete(id)
	if err == nil {
		c.Json(libs.NewError("admincp_vod_uc_del_succ", controllers.RESPONSE_SUCCESS, "成功删除", ""))
		return
	}
	c.Json(libs.NewError("admincp_vod_uc_del_fail", "GM010_051", "其他错误:"+err.Error(), ""))
}

// @Title 完全扫描主播空间地址内视频
// @Description 完全扫描主播空间地址内视频
// @Param   id  path int true  "id"
// @Success 200
// @router /uc/scan_all [post]
func (c *VodCPController) UserVodCenterReptileScanAll() {
	id, _ := c.GetInt64("id")
	if id <= 0 {
		c.Json(libs.NewError("admincp_vod_uc_scanall_fail", "GM010_060", "必须提供删除id", ""))
		return
	}
	ucp := &vod.VodUcenterReptile{}
	ucv := ucp.Get(id)
	if ucv == nil {
		c.Json(libs.NewError("admincp_vod_uc_scanall_fail", "GM010_061", "对象不存在", ""))
		return
	}
	ucv.ScanAll = true
	err := ucp.Update(*ucv)
	if err == nil {
		c.Json(libs.NewError("admincp_vod_uc_scanall_succ", controllers.RESPONSE_SUCCESS, "设置成功,下次抓取时将进行全部扫描", ""))
		return
	}
	c.Json(libs.NewError("admincp_vod_uc_del_fail", "GM010_062", "其他错误:"+err.Error(), ""))
}

// @Title 抓取主播空间地址列表
// @Description 抓取主播空间地址列表
// @Param   page  path int false  "page"
// @Param   size  path int false  "size"
// @Success 200 {object} outobjs.OutVodUcenterPageForAdmin
// @router /uc/list [get]
func (c *VodCPController) UserVodCenterReptileList() {
	page, _ := c.GetInt("page")
	size, _ := c.GetInt("size")
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	out_vuc := []*outobjs.OutVodUcenterForAdmin{}
	ucp := &vod.VodUcenterReptile{}
	total, list := ucp.Gets(int(page), int(size))
	for _, vuc := range list {
		out_vuc = append(out_vuc, &outobjs.OutVodUcenterForAdmin{
			Id:         vuc.Id,
			Uid:        vuc.Uid,
			Member:     outobjs.GetOutMember(vuc.Uid, 0),
			Source:     vuc.Source,
			SiteUrl:    vuc.SiteUrl,
			LastTime:   vuc.LastTime,
			ScanAll:    vuc.ScanAll,
			CreateTime: vuc.CreateTime,
		})
	}
	out := outobjs.OutVodUcenterPageForAdmin{
		CurrentPage: int(page),
		Total:       total,
		Pages:       utils.TotalPages(total, int(size)),
		Size:        int(size),
		Lists:       out_vuc,
	}
	c.Json(out)
}

// @Title 改变抓取主播空间的归属播主
// @Description 改变抓取主播空间的归属播主
// @Param   ucid  path int true  "播主抓取集id"
// @Param   to_uid  path int true  "切换到uid"
// @Success 200
// @router /uc/change_user [post]
func (c *VodCPController) UserVodCenterChangeUser() {
	ucid, _ := c.GetInt64("ucid")
	toUid, _ := c.GetInt64("to_uid")
	if ucid <= 0 || toUid <= 0 {
		c.Json(libs.NewError("admincp_vod_uc_change_user_fail", "GM010_070", "参数错误", ""))
		return
	}
	ucs := &vod.VodUcenterReptile{}
	err := ucs.ChangeToUser(ucid, toUid, true)
	if err == nil {
		c.Json(libs.NewError("admincp_vod_uc_change_user_succ", controllers.RESPONSE_SUCCESS, "设置成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_vod_uc_change_user_fail", "GM010_079", "参数错误:"+err.Error(), ""))
}

// @Title 合集列表
// @Description 合集列表
// @Param   page  path int false  "page"
// @Param   size  path int false  "size"
// @Success 200 {object} outobjs.OutVodPlaylistPagedListForAdmin
// @router /playlists [get]
func (c *VodCPController) VodPlayLists() {
	page, _ := c.GetInt("page")
	size, _ := c.GetInt("size")
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	ups := &vod.Vods{}
	total, list := ups.GetPlaylistsForAdmin(page, size)
	outp := []*outobjs.OutVodPlaylist{}
	for _, p := range list {
		outp = append(outp, outobjs.GetOutVodPlaylist(p))
	}
	outpl := &outobjs.OutVodPlaylistPagedListForAdmin{
		CurrentPage: page,
		Total:       total,
		Pages:       utils.TotalPages(total, size),
		Size:        size,
		Lists:       outp,
	}
	c.Json(outpl)
}

// @Title 合集
// @Description 合集
// @Param   plid  path int true  "合集id"
// @Success 200 {object} outobjs.OutVodPlaylist
// @router /playlist [get]
func (c *VodCPController) VodPlayListGet() {
	plid, _ := c.GetInt64("plid")
	if plid <= 0 {
		c.Json(libs.NewError("admincp_vod_playlist_get_fail", "GM010_080", "参数错误", ""))
		return
	}
	ups := &vod.Vods{}
	pl := ups.GetPlaylist(plid)
	if pl == nil {
		c.Json(libs.NewError("admincp_vod_playlist_get_fail", "GM010_081", "合集不存在", ""))
		return
	}
	outp := outobjs.GetOutVodPlaylist(pl)
	c.Json(outp)
}

// @Title 创建合集
// @Description 创建合集
// @Param   title  path string true  "标题"
// @Param   description  path string false  "描述"
// @Param   img  path int false  "图片"
// @Success 200  {object} libs.Error
// @router /playlist/create [post]
func (c *VodCPController) CreatePlaylist() {
	title, _ := utils.UrlDecode(c.GetString("title"))
	des, _ := utils.UrlDecode(c.GetString("description"))
	img, _ := c.GetInt64("img")
	current_uid := c.CurrentUid()
	if len(title) == 0 {
		c.Json(libs.NewError("admincp_vod_playlist_create_fail", "GM010_090", "标题不能为空", ""))
		return
	}

	ups := &vod.Vods{}
	pl := &vod.VideoPlaylist{
		Title:    title,
		Des:      des,
		PostTime: time.Now(),
		Vods:     0,
		Img:      img,
		Uid:      current_uid,
	}
	_, err := ups.CreatePlaylist(pl, nil)
	if err != nil {
		c.Json(libs.NewError("admincp_vod_playlist_create_fail", "GM010_092", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("admincp_vod_playlist_create_success", controllers.RESPONSE_SUCCESS, "创建成功", ""))
}

// @Title 更新合集中的视频
// @Description 更新合集中的视频
// @Param   plid  path int true  "合集id"
// @Param   vids  path string false  "视频ids(1|2|xx 1表示视频id,2表示排序no,xx表示标题)"
// @Success 200  {object} libs.Error
// @router /playlist/update_vods [post]
func (c *VodCPController) UpdatePlaylistVods() {
	plid, _ := c.GetInt64("plid")
	vidstrs := c.GetString("vids")
	vids := []*vod.VPVData{}
	vidss := strings.Split(vidstrs, ",")
	for _, vstr := range vidss {
		vd := strings.Split(vstr, "|")
		if len(vd) != 3 {
			continue
		}
		_id, _ := strconv.ParseInt(vd[0], 10, 64)
		_no, _ := strconv.Atoi(vd[1])
		_title, _ := utils.UrlDecode(vd[2])
		if _id > 0 {
			vids = append(vids, &vod.VPVData{
				VideoId: _id,
				No:      _no,
				Title:   _title,
			})
		}
	}
	if len(vids) == 0 {
		c.Json(libs.NewError("admincp_vod_playlist_append_fail", "GM010_1001", "添加视频编号数不能为0", ""))
		return
	}
	ups := &vod.Vods{}
	err := ups.UpdatePlaylistVods(plid, vids)
	if err != nil {
		c.Json(libs.NewError("admincp_vod_playlist_append_fail", "GM010_1002", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("admincp_vod_playlist_append_success", controllers.RESPONSE_SUCCESS, "添加成功", ""))
}

//// @Title 合集中添加视频
//// @Description 合集中添加视频
//// @Param   plid  path int true  "合集id"
//// @Param   vids  path string false  "视频ids(1|2 1表示视频id,2表示排序no)"
//// @Success 200  {object} libs.Error
//// @router /playlist/append_vods [post]
//func (c *VodCPController) AppendPlaylistVods() {
//	plid, _ := c.GetInt64("plid")
//	vidstrs := c.GetString("vids")
//	vids := make(map[int64]int)
//	vidss := strings.Split(vidstrs, ",")
//	for _, vstr := range vidss {
//		vd := strings.Split(vstr, "|")
//		if len(vd) != 2 {
//			continue
//		}
//		_id, _ := strconv.ParseInt(vd[0], 10, 64)
//		_no, _ := strconv.ParseInt(vd[1], 10, 64)
//		if _id > 0 {
//			vids[_id] = int(_no)
//		}
//	}
//	if len(vids) == 0 {
//		c.Json(libs.NewError("admincp_vod_playlist_append_fail", "GM010_1001", "添加视频编号数不能为0", ""))
//		return
//	}
//	ups := &vod.Vods{}
//	err := ups.AppenedPlsVods(plid, vids)
//	if err != nil {
//		c.Json(libs.NewError("admincp_vod_playlist_append_fail", "GM010_1002", err.Error(), ""))
//		return
//	}
//	c.Json(libs.NewError("admincp_vod_playlist_append_success", controllers.RESPONSE_SUCCESS, "添加成功", ""))
//}

//// @Title 删除合集中视频
//// @Description 删除合集中视频
//// @Param   plid  path int true  "合集id"
//// @Param   vid  path int true  "视频id"
//// @Success 200  {object} libs.Error
//// @router /playlist/remove_vod [post]
//func (c *VodCPController) RemovePlaylostVods() {
//	plid, _ := c.GetInt64("plid")
//	vid, _ := c.GetInt64("vid")
//	if plid <= 0 {
//		c.Json(libs.NewError("admincp_vod_playlist_remove_fail", "GM010_1010", "合集id错误", ""))
//		return
//	}
//	if vid <= 0 {
//		c.Json(libs.NewError("admincp_vod_playlist_remove_fail", "GM010_1011", "视频id错误", ""))
//		return
//	}
//	ups := &vod.Vods{}
//	err := ups.RemovePlsVod(plid, vid)
//	if err != nil {
//		c.Json(libs.NewError("admincp_vod_playlist_remove_fail", "GM010_1012", err.Error(), ""))
//		return
//	}
//	c.Json(libs.NewError("admincp_vod_playlist_remove_success", controllers.RESPONSE_SUCCESS, "删除成功", ""))
//}

// @Title 合集中视频
// @Description 合集中视频
// @Param   plid  path int true  "合集id"
// @Param   page  path int false  "页"
// @Param   size  path int false  "页数"
// @Success 200  {object} outobjs.OutVodPlaylistVodPageForAdmin
// @router /playlist/vods [get]
func (c *VodCPController) PlaylistVods() {
	plid, _ := c.GetInt64("plid")
	page, _ := c.GetInt("page")
	size, _ := c.GetInt("size")
	if plid <= 0 {
		c.Json(libs.NewError("admincp_vod_playlist_vods_fail", "GM010_1020", "合集id错误", ""))
		return
	}
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 1000
	}
	ups := &vod.Vods{}
	total, plvs := ups.GetPlaylistPlvsForAdmin(plid, page, size)
	outp := []*outobjs.OutVodPlaylistVod{}
	for _, plv := range plvs {
		vod := ups.Get(plv.VideoId, false)
		if vod == nil {
			continue
		}
		_outvod := outobjs.GetOutVodForAdmin(vod)
		if len(plv.Title) > 0 {
			_outvod.Title = plv.Title
		}
		outp = append(outp, &outobjs.OutVodPlaylistVod{
			No:    plv.No,
			VodId: plv.VideoId,
			Vod:   _outvod,
		})
	}
	outpl := &outobjs.OutVodPlaylistVodPageForAdmin{
		CurrentPage: page,
		Total:       total,
		Pages:       utils.TotalPages(total, size),
		Size:        size,
		Lists:       outp,
	}
	c.Json(outpl)
}
