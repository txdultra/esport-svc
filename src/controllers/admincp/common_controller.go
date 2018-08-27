package admincp

import (
	"controllers"
	"encoding/json"
	"libs"
	//"libs/passport"
	"libs/matchrace"
	"libs/version"
	"outobjs"
	"time"
	"utils"
)

// 公共管理模块 API
type CommonCPController struct {
	AdminController
	storage libs.IFileStorage
}

func (c *CommonCPController) Prepare() {
	c.AdminController.Prepare()
	c.storage = libs.NewFileStorage()
}

// @Title 游戏
// @Description 游戏列表
// @Success 200 {object} outobjs.OutGameForAdmin
// @router /game/list [get]
func (c *CommonCPController) Games() {
	bas := &libs.Bas{}
	games := bas.Games()
	out_games := []*outobjs.OutGameForAdmin{}
	for _, g := range games {
		out_games = append(out_games, &outobjs.OutGameForAdmin{
			Id:           g.Id,
			Name:         g.Name,
			En:           g.En,
			Img:          g.Img,
			ImgUrl:       c.storage.GetFileUrl(g.Img),
			Enabled:      g.Enabled,
			PostTime:     g.PostTime,
			DisplayOrder: g.DisplayOrder,
			ForcedzSel:   g.ForcedzSel,
		})
	}
	c.Json(out_games)
}

// @Title 添加新游戏
// @Description 添加新游戏
// @Param   name   path	string true  "游戏名称"
// @Param   en   path	string true  "游戏英文名"
// @Param   img_id   path	int true  "图片id"
// @Param   enabled   path	bool true  "是否启用"
// @Param   display_order   path	int false  "排序"
// @Param   forcedz_sel   path	bool true  "强制选择"
// @Success 200 {object} libs.Error
// @router /game/add [post]
func (c *CommonCPController) GameAdd() {
	name, _ := utils.UrlDecode(c.GetString("name"))
	en, _ := utils.UrlDecode(c.GetString("en"))
	img_id, _ := c.GetInt64("img_id")
	enabled, _ := c.GetBool("enabled")
	display_order, _ := c.GetInt64("display_order")
	forcedz, _ := c.GetBool("forcedz_sel")
	if len(name) == 0 {
		c.Json(libs.NewError("admincp_common_game_add_fail", "GM001_001", "游戏名称不能空", ""))
		return
	}
	if len(en) == 0 {
		c.Json(libs.NewError("admincp_common_game_add_fail", "GM001_002", "游戏英文名不能空", ""))
		return
	}
	game := &libs.Game{}
	game.Name = name
	game.En = en
	game.Img = img_id
	game.Enabled = enabled
	game.PostTime = time.Now()
	game.DisplayOrder = int(display_order)
	game.ForcedzSel = forcedz
	bas := &libs.Bas{}
	err := bas.AddGame(game)
	if err == nil {
		c.Json(libs.NewError("admincp_common_game_add_succ", controllers.RESPONSE_SUCCESS, "添加成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_game_add_fail", "GM001_003", err.Error(), ""))
}

// @Title 更新游戏
// @Description 更新游戏
// @Param   id   path	int true  "游戏id"
// @Param   name   path	string true  "游戏名称"
// @Param   en   path	string true  "游戏英文名"
// @Param   img_id   path	int true  "图片id"
// @Param   enabled   path	bool true  "是否启用(默认启用)"
// @Param   forcedz_sel   path	bool true  "强制选择"
// @Param   display_order   path	int false  "排序"
// @Success 200 {object} libs.Error
// @router /game/update [post]
func (c *CommonCPController) GameUpdate() {
	id, _ := c.GetInt64("id")
	name, _ := utils.UrlDecode(c.GetString("name"))
	en, _ := utils.UrlDecode(c.GetString("en"))
	img_id, _ := c.GetInt64("img_id")
	enabled, _ := c.GetBool("enabled")
	forcedz, _ := c.GetBool("forcedz_sel")
	display_order, _ := c.GetInt("display_order")
	if len(c.Ctx.Input.Query("enabled")) == 0 {
		enabled = true
	}
	if id <= 0 {
		c.Json(libs.NewError("admincp_common_game_update_fail", "GM001_010", "游戏不存在", ""))
		return
	}
	if len(name) == 0 {
		c.Json(libs.NewError("admincp_common_game_update_fail", "GM001_011", "游戏名称不能空", ""))
		return
	}
	if len(en) == 0 {
		c.Json(libs.NewError("admincp_common_game_update_fail", "GM001_012", "游戏英文名不能空", ""))
		return
	}
	game := &libs.Game{}
	game.Id = int(id)
	game.Name = name
	game.En = en
	game.Img = img_id
	game.Enabled = enabled
	game.PostTime = time.Now()
	game.DisplayOrder = int(display_order)
	game.ForcedzSel = forcedz
	bas := &libs.Bas{}
	err := bas.UpdateGame(game)
	if err == nil {
		c.Json(libs.NewError("admincp_common_game_update_succ", controllers.RESPONSE_SUCCESS, "更新成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_game_update_fail", "GM001_013", err.Error(), ""))
}

// @Title 赛事列表
// @Description 赛事列表
// @Param   isall   path	bool false  "isall"
// @Param   enabled   path	bool false  "enabled"
// @Success 200 {object} outobjs.OutMatchForAdmin
// @router /match/list [get]
func (c *CommonCPController) Matchs() {
	enabled, _ := c.GetBool("enabled")
	all, _ := c.GetBool("isall")
	bas := &libs.Bas{}
	matchs := bas.Matchs()
	out_matchs := []*outobjs.OutMatchForAdmin{}
	for _, m := range matchs {
		if !all {
			if m.Enabled != enabled {
				continue
			}
		}
		out_matchs = append(out_matchs, &outobjs.OutMatchForAdmin{
			Id:       m.Id,
			Name:     m.Name,
			SubTitle: m.SubTitle,
			En:       m.En,
			ImgId:    m.Img,
			ImgUrl:   c.storage.GetFileUrl(m.Img),
			Des1:     m.Des1,
			Des2:     m.Des2,
			Des3:     m.Des3,
			Enabled:  m.Enabled,
			IconId:   m.Icon,
			IconUrl:  c.storage.GetFileUrl(m.Icon),
		})
	}
	c.Json(out_matchs)
}

// @Title 新建赛事
// @Description 新建赛事
// @Param   name   path	string true  "赛事名称"
// @Param   sub_title   path	string true  "二级标题"
// @Param   img_id   path	int true  "图片id"
// @Param   en   path	string true  "赛事英文名"
// @Param   des1   path	string true  "赛事描述1"
// @Param   des2   path	string true  "赛事描述2"
// @Param   des3   path	string true  "赛事描述3"
// @Param   enabled   path	bool true  "开启"
// @Param   icon_id	path int true "图标id"
// @Success 200 {object} libs.Error
// @router /match/add [post]
func (c *CommonCPController) MatchAdd() {
	name, _ := utils.UrlDecode(c.GetString("name"))
	sub_title, _ := utils.UrlDecode(c.GetString("sub_title"))
	en, _ := utils.UrlDecode(c.GetString("en"))
	img_id, _ := c.GetInt64("img_id")
	des1, _ := utils.UrlDecode(c.GetString("des1"))
	des2, _ := utils.UrlDecode(c.GetString("des2"))
	des3, _ := utils.UrlDecode(c.GetString("des3"))
	enabled, _ := c.GetBool("enabled")
	icon_id, _ := c.GetInt64("icon_id")

	if len(name) == 0 {
		c.Json(libs.NewError("admincp_common_match_add_fail", "GM002_001", "赛事名称不能空", ""))
		return
	}
	if len(en) == 0 {
		c.Json(libs.NewError("admincp_common_match_add_fail", "GM002_002", "赛事英文名不能空", ""))
		return
	}
	match := &libs.Match{}
	match.Name = name
	match.SubTitle = sub_title
	match.En = en
	match.Img = img_id
	match.Des1 = des1
	match.Des2 = des2
	match.Des3 = des3
	match.Enabled = enabled
	match.Icon = icon_id
	bas := &libs.Bas{}
	err := bas.AddMatch(match)
	if err == nil {
		c.Json(libs.NewError("admincp_common_match_add_succ", controllers.RESPONSE_SUCCESS, "添加成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_match_add_fail", "GM002_003", err.Error(), ""))
}

// @Title 更新赛事
// @Description 更新赛事
// @Param   id   path	string true  "赛事id"
// @Param   name   path	string true  "赛事名称"
// @Param   sub_title   path	string true  "二级标题"
// @Param   img_id   path	int true  "图片id"
// @Param   en   path	string true  "赛事英文名"
// @Param   des1   path	string true  "赛事描述1"
// @Param   des2   path	string true  "赛事描述2"
// @Param   des3   path	string true  "赛事描述3"
// @Param   enabled   path	bool true  "开启"
// @Param   icon_id	path int true "图标id"
// @Success 200 {object} libs.Error
// @router /match/update [post]
func (c *CommonCPController) MatchUpdate() {
	id, _ := c.GetInt("id")
	name, _ := utils.UrlDecode(c.GetString("name"))
	sub_title, _ := utils.UrlDecode(c.GetString("sub_title"))
	en, _ := utils.UrlDecode(c.GetString("en"))
	img_id, _ := c.GetInt64("img_id")
	des1, _ := utils.UrlDecode(c.GetString("des1"))
	des2, _ := utils.UrlDecode(c.GetString("des2"))
	des3, _ := utils.UrlDecode(c.GetString("des3"))
	enabled, _ := c.GetBool("enabled")
	icon_id, _ := c.GetInt64("icon_id")

	if id <= 0 {
		c.Json(libs.NewError("admincp_common_match_update_fail", "GM002_010", "赛事不存在", ""))
		return
	}
	if len(name) == 0 {
		c.Json(libs.NewError("admincp_common_match_update_fail", "GM002_011", "赛事名称不能空", ""))
		return
	}
	if len(en) == 0 {
		c.Json(libs.NewError("admincp_common_match_update_fail", "GM002_012", "赛事英文名不能空", ""))
		return
	}
	bas := &libs.Bas{}
	match := bas.Match(id)
	if match == nil {
		c.Json(libs.NewError("admincp_common_match_update_fail", "GM002_014", "对象不存在", ""))
	}
	match.Name = name
	match.SubTitle = sub_title
	match.En = en
	match.Img = img_id
	match.Des1 = des1
	match.Des2 = des2
	match.Des3 = des3
	match.Enabled = enabled
	match.Icon = icon_id

	err := bas.UpdateMatch(match)
	if err == nil {
		c.Json(libs.NewError("admincp_common_match_update_succ", controllers.RESPONSE_SUCCESS, "添加成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_match_update_fail", "GM002_013", err.Error(), ""))
}

// @Title 推荐项目列表
// @Description 推荐项目列表
// @Param   category   path	string true  "组类型"
// @Success 200 {object} outobjs.OutRecommendForAdmin
// @router /recommend/list [get]
func (c *CommonCPController) RecommendList() {
	category := c.GetString("category")
	recs := libs.NewRecommendService()
	rems := recs.Gets(category)
	out_rems := []*outobjs.OutRecommendForAdmin{}
	for _, rem := range rems {
		out_rems = append(out_rems, c.getOutRecommendForAdmin(rem))
	}
	c.Json(out_rems)
}

// @Title 添加推荐项目
// @Description 添加推荐项目
// @Param   ref_id   path	int true  "关联id"
// @Param   ref_type   path	string true  "关联id的数据对象类型(video)"
// @Param   title   path	string true  "标题"
// @Param   img_id   path	int true  "图片id"
// @Param   category   path	string true  "推荐分类(vod_home)"
// @Param   display   path	int true  "显示顺序"
// @Param   enabled   path	bool true  "是否启用"
// @Success 200 {object} libs.Error
// @router /recommend/add [post]
func (c *CommonCPController) RecommendAdd() {
	ref_id, _ := c.GetInt64("ref_id")
	ref_type := c.GetString("ref_type")
	title, _ := utils.UrlDecode(c.GetString("title"))
	img_id, _ := c.GetInt64("img_id")
	category := c.GetString("category")
	display, _ := c.GetInt64("display")
	enabled, _ := c.GetBool("enabled")
	if ref_id <= 0 {
		c.Json(libs.NewError("admincp_common_recommend_add_fail", "GM003_001", "关联id不能为0", ""))
		return
	}
	if len(ref_type) == 0 {
		c.Json(libs.NewError("admincp_common_recommend_add_faill", "GM003_002", "关联id的数据对象类型不能空", ""))
		return
	}
	if len(category) == 0 {
		c.Json(libs.NewError("admincp_common_recommend_add_fail", "GM003_003", "推荐分类不能空", ""))
		return
	}
	recommend := libs.Recommend{}
	recommend.RefId = ref_id
	recommend.RefType = ref_type
	recommend.Title = title
	recommend.Img = img_id
	recommend.Category = category
	recommend.DisplayOrder = int(display)
	recommend.Enabled = enabled
	recommend.PostTime = time.Now()
	recommend.PostUid = c.CurrentUid()
	recs := libs.NewRecommendService()
	id, err := recs.Create(recommend)
	if id <= 0 {
		c.Json(libs.NewError("admincp_common_recommend_add_fail", "GM003_004", "提交失败:"+err.Error(), ""))
		return
	}
	c.Json(libs.NewError("admincp_common_recommend_add_succ", controllers.RESPONSE_SUCCESS, "添加成功", ""))
}

// @Title 更新推荐项目
// @Description 更新推荐项目
// @Param   id   path	int true  "推荐id"
// @Param   ref_id   path	int true  "关联id"
// @Param   title   path	string true  "标题"
// @Param   img_id   path	int true  "图片id"
// @Param   display   path	int true  "显示顺序"
// @Param   enabled   path	bool true  "是否启用"
// @Success 200 {object} libs.Error
// @router /recommend/update [post]
func (c *CommonCPController) RecommendUpdate() {
	id, _ := c.GetInt64("id")
	ref_id, _ := c.GetInt64("ref_id")
	title, _ := utils.UrlDecode(c.GetString("title"))
	img_id, _ := c.GetInt64("img_id")
	display, _ := c.GetInt64("display")
	enabled, _ := c.GetBool("enabled")
	if id <= 0 {
		c.Json(libs.NewError("admincp_common_recommend_update_fail", "GM003_010", "推荐id非法", ""))
		return
	}
	if ref_id <= 0 {
		c.Json(libs.NewError("admincp_common_recommend_update_fail", "GM003_011", "关联id不能为0", ""))
		return
	}
	recs := libs.NewRecommendService()
	recommend := recs.Get(id)
	if recommend == nil {
		c.Json(libs.NewError("admincp_common_recommend_update_fail", "GM003_012", "推荐对象不存在", ""))
		return
	}
	recommend.RefId = ref_id
	recommend.Title = title
	recommend.Img = img_id
	recommend.DisplayOrder = int(display)
	recommend.Enabled = enabled
	recommend.PostTime = time.Now()
	recommend.PostUid = c.CurrentUid()
	err := recs.Update(*recommend)
	if err != nil {
		c.Json(libs.NewError("admincp_common_recommend_update_fail", "GM003_013", "提交失败:"+err.Error(), ""))
		return
	}
	c.Json(libs.NewError("admincp_common_recommend_update_succ", controllers.RESPONSE_SUCCESS, "更新成功", ""))
}

// @Title 删除推荐项目
// @Description 删除推荐项目
// @Param   id   path	int true  "推荐id"
// @Success 200 {object} libs.Error
// @router /recommend/del [delete]
func (c *CommonCPController) RecommendDelete() {
	id, _ := c.GetInt64("id")
	if id <= 0 {
		c.Json(libs.NewError("admincp_common_recommend_delete_fail", "GM003_020", "推荐id非法", ""))
		return
	}
	recs := libs.NewRecommendService()
	recommend := recs.Get(id)
	if recommend == nil {
		c.Json(libs.NewError("admincp_common_recommend_delete_fail", "GM003_021", "推荐对象不存在", ""))
		return
	}
	err := recs.Delete(id)
	if err != nil {
		c.Json(libs.NewError("admincp_common_recommend_delete_fail", "GM003_013", "提交失败:"+err.Error(), ""))
		return
	}
	c.Json(libs.NewError("admincp_common_recommend_delete_succ", controllers.RESPONSE_SUCCESS, "删除成功", ""))
}

// @Title 推荐项目列表
// @Description 推荐项目列表
// @Param   id   path	int true  "id"
// @Success 200 {object} outobjs.OutRecommendForAdmin
// @router /recommend/:id([0-9]+) [get]
func (c *CommonCPController) RecommendGet() {
	id, _ := c.GetInt64(":id")
	if id <= 0 {
		c.Json(libs.NewError("admincp_common_recommend_get_fail", "GM003_020", "id不能为0", ""))
		return
	}
	recs := libs.NewRecommendService()
	rem := recs.Get(id)
	if rem == nil {
		c.Json(libs.NewError("admincp_common_recommend_get_fail", "GM003_021", "对象不存在", ""))
		return
	}
	out := c.getOutRecommendForAdmin(rem)
	c.Json(out)
}

func (c *CommonCPController) getOutRecommendForAdmin(rem *libs.Recommend) *outobjs.OutRecommendForAdmin {
	return &outobjs.OutRecommendForAdmin{
		Id:           rem.Id,
		RefId:        rem.RefId,
		RefType:      rem.RefType,
		Title:        rem.Title,
		ImgUrl:       c.storage.GetFileUrl(rem.Img),
		ImgId:        rem.Img,
		Categroy:     rem.Category,
		Enabled:      rem.Enabled,
		DisplayOrder: rem.DisplayOrder,
		PostTime:     rem.PostTime,
		PostMember:   outobjs.GetOutMember(rem.PostUid, 0),
	}
}

// @Title 版本管理
// @Description 版本管理(数组)
// @Param   plat   path	string true  "平台(android,ios,wphone)"
// @Success 200 {object} outobjs.OutVersionForAdmin
// @router /versions [get]
func (c *CommonCPController) Versions() {
	plt := c.GetString("plat")
	platform := version.MOBILE_PLATFORM(plt)
	vcs := &version.VCS{}
	vers := vcs.GetClientVersions(platform)
	outs := []*outobjs.OutVersionForAdmin{}
	for _, v := range vers {
		outs = append(outs, &outobjs.OutVersionForAdmin{
			Id:           v.Id,
			Version:      v.Version,
			Ver:          v.Ver,
			VerName:      v.VerName,
			Description:  v.Description,
			PostTime:     v.PostTime,
			Platform:     v.Platform,
			IsExpried:    v.IsExpried,
			DownloadUrl:  v.DownloadUrl,
			AllowVodDown: v.AllowVodDown,
		})
	}
	c.Json(outs)
}

// @Title 版本管理
// @Description 版本管理(单个)
// @Param   id   path	int true  "id"
// @Param   plat   path	string true  "平台(android,ios,wphone)"
// @Success 200 {object} outobjs.OutVersionForAdmin
// @router /version [get]
func (c *CommonCPController) Version() {
	plt := c.GetString("plat")
	id, _ := c.GetInt64("id")
	platform := version.MOBILE_PLATFORM(plt)
	vcs := &version.VCS{}
	v := vcs.GetClientVersionById(platform, id)
	if v == nil {
		c.Json(libs.NewError("admincp_common_version_add_fail", "GM004_031", "版本不存在", ""))
		return
	}
	c.Json(outobjs.OutVersionForAdmin{
		Id:           v.Id,
		Version:      v.Version,
		Ver:          v.Ver,
		VerName:      v.VerName,
		Description:  v.Description,
		PostTime:     v.PostTime,
		Platform:     v.Platform,
		IsExpried:    v.IsExpried,
		DownloadUrl:  v.DownloadUrl,
		AllowVodDown: v.AllowVodDown,
	})
}

// @Title 新建版本
// @Description 新建版本
// @Param   version   path	float true  "内部版本号"
// @Param   ver   path	string true  "版本标识"
// @Param   ver_name   path	string true  "版本名称"
// @Param   desc   path	string true  "版本描述"
// @Param   plat   path	string true  "平台(android,ios,wphone)"
// @Param   down_url   path	string true  "下载地址"
// @Param   allow_voddown   path	bool false  "允许视频下载"
// @Success 200 {object} libs.Error
// @router /version/add [post]
func (c *CommonCPController) VersionAdd() {
	_version, _ := c.GetFloat("version")
	ver, _ := utils.UrlDecode(c.GetString("ver"))
	ver_name, _ := utils.UrlDecode(c.GetString("ver_name"))
	desc, _ := utils.UrlDecode(c.GetString("desc"))
	plt := c.GetString("plat")
	down_url, _ := utils.UrlDecode(c.GetString("down_url"))
	allow_voddown, _ := c.GetBool("allow_voddown")

	if _version <= 0 {
		c.Json(libs.NewError("admincp_common_version_add_fail", "GM004_001", "version不能小于等于0", ""))
		return
	}
	if len(ver) == 0 {
		c.Json(libs.NewError("admincp_common_version_add_fail", "GM004_002", "ver不能为空", ""))
		return
	}
	if len(ver_name) == 0 {
		c.Json(libs.NewError("admincp_common_version_add_fail", "GM004_003", "ver_name不能为空", ""))
		return
	}
	if len(down_url) == 0 {
		c.Json(libs.NewError("admincp_common_version_add_fail", "GM004_004", "下载地址不能为空", ""))
		return
	}
	v := &version.ClientVersion{
		Version:      _version,
		Ver:          ver,
		VerName:      ver_name,
		Description:  desc,
		Platform:     version.MOBILE_PLATFORM(plt),
		IsExpried:    false,
		DownloadUrl:  down_url,
		AllowVodDown: allow_voddown,
	}
	vcs := &version.VCS{}
	_, err := vcs.Create(v)
	if err == nil {
		c.Json(libs.NewError("admincp_common_version_add_succ", controllers.RESPONSE_SUCCESS, "添加成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_version_add_fail", "GM004_005", "添加失败:"+err.Error(), ""))
}

// @Title 更新版本信息
// @Description 更新版本信息
// @Param   id   path	int true  "id"
// @Param   version   path	float true  "内部版本号"
// @Param   ver   path	string true  "版本标识"
// @Param   ver_name   path	string true  "版本名称"
// @Param   desc   path	string true  "版本描述"
// @Param   is_expried   path	bool true  "是否过期"
// @Param   down_url   path	string true  "下载地址"
// @Param   plat   path	string true  "平台(android,ios,wphone)"
// @Param   allow_voddown   path	bool false  "允许视频下载"
// @Success 200 {object} libs.Error
// @router /version/update [post]
func (c *CommonCPController) VersionUpdate() {
	id, _ := c.GetInt64("id")
	_version, _ := c.GetFloat("version")
	ver, _ := utils.UrlDecode(c.GetString("ver"))
	ver_name, _ := utils.UrlDecode(c.GetString("ver_name"))
	desc, _ := utils.UrlDecode(c.GetString("desc"))
	plat := c.GetString("plat")
	is_expried, _ := c.GetBool("is_expried")
	down_url, _ := utils.UrlDecode(c.GetString("down_url"))
	allow_voddown, _ := c.GetBool("allow_voddown")

	if id <= 0 {
		c.Json(libs.NewError("admincp_common_version_update_fail", "GM004_011", "id不能小于等于0", ""))
		return
	}
	if _version <= 0 {
		c.Json(libs.NewError("admincp_common_version_update_fail", "GM004_012", "version不能小于等于0", ""))
		return
	}
	if len(ver) == 0 {
		c.Json(libs.NewError("admincp_common_version_update_fail", "GM004_013", "ver不能为空", ""))
		return
	}
	if len(ver_name) == 0 {
		c.Json(libs.NewError("admincp_common_version_update_fail", "GM004_014", "ver_name不能为空", ""))
		return
	}
	if len(down_url) == 0 {
		c.Json(libs.NewError("admincp_common_version_update_fail", "GM004_015", "下载地址不能为空", ""))
		return
	}
	if len(plat) == 0 {
		c.Json(libs.NewError("admincp_common_version_update_fail", "GM004_016", "platform不能为空", ""))
		return
	}
	vcs := &version.VCS{}
	v := vcs.GetClientVersionById(version.MOBILE_PLATFORM(plat), id)
	if v == nil {
		c.Json(libs.NewError("admincp_common_version_add_fail", "GM004_016", "版本不存在", ""))
		return
	}
	v.Version = _version
	v.Ver = ver
	v.VerName = ver_name
	v.Description = desc
	v.IsExpried = is_expried
	v.DownloadUrl = down_url
	v.AllowVodDown = allow_voddown
	err := vcs.Update(v)
	if err == nil {
		c.Json(libs.NewError("admincp_common_version_update_succ", controllers.RESPONSE_SUCCESS, "更新成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_version_update_fail", "GM004_017", "更新失败:"+err.Error(), ""))
}

// @Title 删除版本
// @Description 删除版本
// @Param   id   path	int true  "id"
// @Param   plat path	string true  "平台(android,ios,wphone)"
// @Success 200 {object} libs.Error
// @router /version/del [delete]
func (c *CommonCPController) VersionDel() {
	id, _ := c.GetInt64("id")
	plat := c.GetString("plat")
	vcs := &version.VCS{}
	if id <= 0 {
		c.Json(libs.NewError("admincp_common_version_del_fail", "GM004_041", "版本不存在", ""))
		return
	}
	if len(plat) == 0 {
		c.Json(libs.NewError("admincp_common_version_del_fail", "GM004_042", "platform不能为空", ""))
		return
	}
	err := vcs.Del(id, version.MOBILE_PLATFORM(plat))
	if err == nil {
		c.Json(libs.NewError("admincp_common_version_update_succ", controllers.RESPONSE_SUCCESS, "删除成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_version_update_fail", "GM004_043", "删除失败:"+err.Error(), ""))
}

// @Title 新增首页广告
// @Description 新增首页广告
// @Param   title   path	string true  "标题"
// @Param   img   path	int true  "图片"
// @Param   action   path	string true  "事件"
// @Param   args   path	string true  "参数"
// @Param   end_time   path	int true  "结束时间"
// @Param   waits   path	int true  "等待时间"
// @Success 200 {object} libs.Error
// @router /homead/add [post]
func (c *CommonCPController) HomeAdAdd() {
	title, _ := utils.UrlDecode(c.GetString("title"))
	img, _ := c.GetInt64("img")
	action := c.GetString("action")
	args, _ := utils.UrlDecode(c.GetString("args"))
	end_time, _ := c.GetInt64("end_time")
	waits, _ := c.GetInt("waits")
	if len(title) == 0 {
		c.Json(libs.NewError("admincp_common_homead_add_fail", "GM004_051", "标题不能为空", ""))
		return
	}
	if img <= 0 {
		c.Json(libs.NewError("admincp_common_homead_add_fail", "GM004_052", "图片未设置", ""))
		return
	}
	if len(action) == 0 {
		c.Json(libs.NewError("admincp_common_homead_add_fail", "GM004_053", "事件不能为空", ""))
		return
	}
	if end_time <= 0 {
		c.Json(libs.NewError("admincp_common_homead_add_fail", "GM004_054", "结束时间未设置", ""))
		return
	}
	if waits <= 0 {
		waits = 3
	}
	ad := &libs.HomeAd{
		Title:    title,
		Img:      img,
		Action:   action,
		Args:     args,
		EndTime:  time.Unix(end_time, 0),
		Waits:    waits,
		PostTime: time.Now(),
		PostUid:  c.CurrentUid(),
	}
	bas := &libs.Bas{}
	err := bas.CreateHomeAd(ad)
	if err == nil {
		c.Json(libs.NewError("admincp_common_homead_add_succ", controllers.RESPONSE_SUCCESS, "添加成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_homead_add_fail", "GM004_055", err.Error(), ""))
}

// @Title 获取最新首页广告
// @Description 获取最新首页广告
// @Success 200 {object} outobjs.OutHomeAdForAdmin
// @router /homead/get [get]
func (c *CommonCPController) GetLastHomeAd() {
	bas := &libs.Bas{}
	ad := bas.LastNewHomeAd()
	if ad == nil {
		c.Json(libs.NewError("admincp_common_homead_get_fail", "GM004_060", "不存在", ""))
		return
	}
	out_c := &outobjs.OutHomeAdForAdmin{
		Id:         ad.Id,
		Title:      ad.Title,
		Img:        ad.Img,
		ImgUrl:     c.storage.GetFileUrl(ad.Img),
		Action:     ad.Action,
		Args:       ad.Args,
		Waits:      ad.Waits,
		EndTime:    ad.EndTime,
		PostTime:   ad.PostTime,
		PostUid:    ad.PostUid,
		PostMember: outobjs.GetOutSimpleMember(ad.PostUid),
	}
	c.Json(out_c)
}

// @Title 删除首页广告
// @Description 删除首页广告
// @Param   id   path	int true  "id"
// @Success 200 {object} libs.Error
// @router /homead/del [delete]
func (c *CommonCPController) DeleteHomeAd() {
	id, _ := c.GetInt64("id")
	if id <= 0 {
		c.Json(libs.NewError("admincp_common_homead_del_fail", "GM004_070", "id非法", ""))
		return
	}
	bas := &libs.Bas{}
	err := bas.DeleteHomeAd(id)
	if err == nil {
		c.Json(libs.NewError("admincp_common_homead_del_succ", controllers.RESPONSE_SUCCESS, "删除成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_homead_del_fail", "GM004_071", err.Error(), ""))
}

// @Title 新增战队
// @Description 新增战队
// @Param   title   path	string true  "标题"
// @Param   img1   path	int true  "图片1"
// @Param   img2   path	int true  "图片2"
// @Param   img3   path	int true  "图片2"
// @Param   description   path	string false  "描述"
// @Param   tt   path	int true  "战队类型"
// @Param   del   path	bool false  "删除"
// @Param   pid   path	int false  "父id"
// @Success 200 {object} libs.Error
// @router /team/add [post]
func (c *CommonCPController) TeamAdd() {
	title, _ := utils.UrlDecode(c.GetString("title"))
	img1, _ := c.GetInt64("img1")
	img2, _ := c.GetInt64("img2")
	img3, _ := c.GetInt64("img3")
	desc, _ := utils.UrlDecode(c.GetString("description"))
	tt, _ := c.GetInt("tt")
	del, _ := c.GetBool("del")
	pid, _ := c.GetInt64("pid")

	if len(title) == 0 {
		c.Json(libs.NewError("admincp_common_teamadd_fail", "GM005_010", "title不能为空", ""))
		return
	}
	team := &libs.Team{
		Title:       title,
		Img1:        img1,
		Img2:        img2,
		Img3:        img3,
		Description: desc,
		TeamType:    libs.TEAM_TYPE(tt),
		Del:         del,
		ParentId:    pid,
	}
	bas := &libs.Bas{}
	err := bas.CreateTeam(team)
	if err == nil {
		c.Json(libs.NewError("admincp_common_teamadd_succ", controllers.RESPONSE_SUCCESS, "添加成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_teamadd_fail", "GM005_011", err.Error(), ""))
}

// @Title 更新战队
// @Description 新增战队
// @Param   id   path	int true  "id"
// @Param   title   path	string true  "标题"
// @Param   img1   path	int true  "图片1"
// @Param   img2   path	int true  "图片2"
// @Param   img3   path	int true  "图片2"
// @Param   description   path	string false  "描述"
// @Param   tt   path	int true  "战队类型"
// @Param   del   path	bool false  "删除"
// @Param   pid   path	int false  "父id"
// @Success 200 {object} libs.Error
// @router /team/update [post]
func (c *CommonCPController) TeamUpdate() {
	id, _ := c.GetInt64("id")
	title, _ := utils.UrlDecode(c.GetString("title"))
	img1, _ := c.GetInt64("img1")
	img2, _ := c.GetInt64("img2")
	img3, _ := c.GetInt64("img3")
	desc, _ := utils.UrlDecode(c.GetString("description"))
	tt, _ := c.GetInt("tt")
	del, _ := c.GetBool("del")
	pid, _ := c.GetInt64("pid")

	if len(title) == 0 {
		c.Json(libs.NewError("admincp_common_teamupdate_fail", "GM005_020", "title不能为空", ""))
		return
	}
	bas := &libs.Bas{}
	team := bas.GetTeam(id)
	if team == nil {
		c.Json(libs.NewError("admincp_common_teamupdate_fail", "GM005_021", "对象不存在", ""))
		return
	}
	team.Title = title
	team.Img1 = img1
	team.Img2 = img2
	team.Img3 = img3
	team.Description = desc
	team.TeamType = libs.TEAM_TYPE(tt)
	team.Del = del
	team.ParentId = pid
	err := bas.UpdateTeam(team)
	if err == nil {
		c.Json(libs.NewError("admincp_common_teamupdate_succ", controllers.RESPONSE_SUCCESS, "更新成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_teamupdate_fail", "GM005_022", err.Error(), ""))
}

// @Title 删除战队
// @Description 删除战队
// @Param   id   path	int true  "id"
// @Success 200 {object} libs.Error
// @router /team/remove [delete]
func (c *CommonCPController) DelTeam() {
	id, _ := c.GetInt64("id")
	if id <= 0 {
		c.Json(libs.NewError("admincp_common_teamdel_fail", "GM005_040", "id不能为0", ""))
		return
	}
	bas := &libs.Bas{}
	team := bas.GetTeam(id)
	if team == nil {
		c.Json(libs.NewError("admincp_common_teamdel_fail", "GM005_041", "对象不存在", ""))
		return
	}
	err := bas.DelTeam(id)
	if err == nil {
		c.Json(libs.NewError("admincp_common_teamdel_succ", controllers.RESPONSE_SUCCESS, "删除成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_teamdel_fail", "GM005_042", err.Error(), ""))
}

// @Title 获取战队
// @Description 获取战队
// @Param   id   path	int true  "id"
// @Success 200 {object} outobjs.OutTeam
// @router /team/:id([0-9]+) [get]
func (c *CommonCPController) GetTeam() {
	id, _ := c.GetInt64(":id")
	if id <= 0 {
		c.Json(libs.NewError("admincp_common_teamget_fail", "GM005_030", "id不能为0", ""))
		return
	}
	bas := &libs.Bas{}
	team := bas.GetTeam(id)
	if team == nil {
		c.Json(libs.NewError("admincp_common_teamget_fail", "GM005_031", "对象不存在", ""))
		return
	}
	c.Json(outobjs.GetOutTeam(team))
}

// @Title 查询战队库
// @Description 查询战队库
// @Param   title   path	string false  "标题"
// @Param   tt   path	int true  "类型"
// @Param   del   path	int false  "是否已被删除(1删除0未删除-1忽略)"
// @Param   page   path	int true  "页"
// @Param   page_size   path	int true  "页大小"
// @Success 200 {object} outobjs.OutTeamPagedList
// @router /team/list [get]
func (c *CommonCPController) GetTeams() {
	title := c.GetString("title")
	tt, _ := c.GetInt("tt")
	del, _ := c.GetInt("del")
	page, _ := c.GetInt("page")
	size, _ := c.GetInt("page_size")
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	bas := &libs.Bas{}
	total, teams := bas.GetTeams(title, libs.TEAM_TYPE(tt), del, page, size)
	outteams := make([]*outobjs.OutTeam, len(teams), len(teams))
	for i, team := range teams {
		outteams[i] = outobjs.GetOutTeam(team)
	}
	outp := &outobjs.OutTeamPagedList{
		CurrentPage: page,
		Totals:      total,
		PageSize:    size,
		Teams:       outteams,
	}
	c.Json(outp)
}

// @Title 新增赛程模式
// @Description 新增赛程模式
// @Param   title   path	string true  "标题"
// @Param   match_id   path	int true  "赛事id"
// @Param   mode_type   path	int true  "类型"
// @Param   displayorder   path	int true  "排序"
// @Param   is_view   path	bool true  "是否显示"
// @Success 200 {object} libs.Error
// @router /matchrace/mode_add [post]
func (c *CommonCPController) MatchModeAdd() {
	title, _ := utils.UrlDecode(c.GetString("title"))
	match_id, _ := c.GetInt64("match_id")
	mode_type, _ := c.GetInt16("mode_type")
	displayorder, _ := c.GetInt("displayorder")
	is_view, _ := c.GetBool("is_view")
	if len(title) == 0 {
		c.Json(libs.NewError("admincp_common_matchmode_add_fail", "GM006_001", "标题不能为空", ""))
		return
	}
	if match_id <= 0 {
		c.Json(libs.NewError("admincp_common_matchmode_add_fail", "GM006_002", "未指定关联赛事", ""))
		return
	}
	if mode_type <= 0 {
		c.Json(libs.NewError("admincp_common_matchmode_add_fail", "GM006_003", "类型错误", ""))
		return
	}
	race := &matchrace.RaceMode{
		MatchId:      match_id,
		ModeType:     matchrace.MODE_TYPE(mode_type),
		DisplayOrder: displayorder,
		IsView:       is_view,
		Title:        title,
	}
	mrs := &matchrace.MatchRaceService{}
	err := mrs.CreateMode(race)
	if err == nil {
		c.Json(libs.NewError("admincp_common_matchmode_add_succ", controllers.RESPONSE_SUCCESS, "创建成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_matchmode_add_fail", "GM006_004", err.Error(), ""))
}

// @Title  更新赛程模式
// @Description 更新赛程模式
// @Param   id   path	int true  "id"
// @Param   title   path	string true  "标题"
// @Param   match_id   path	int true  "赛事id"
// @Param   displayorder   path	int true  "排序"
// @Param   is_view   path	bool true  "是否显示"
// @Success 200 {object} libs.Error
// @router /matchrace/mode_update [post]
func (c *CommonCPController) MatchModeUpdate() {
	id, _ := c.GetInt64("id")
	title, _ := utils.UrlDecode(c.GetString("title"))
	match_id, _ := c.GetInt64("match_id")
	displayorder, _ := c.GetInt("displayorder")
	is_view, _ := c.GetBool("is_view")
	if id <= 0 {
		c.Json(libs.NewError("admincp_common_matchmode_update_fail", "GM006_010", "id错误", ""))
		return
	}
	if len(title) == 0 {
		c.Json(libs.NewError("admincp_common_matchmode_update_fail", "GM006_011", "标题不能为空", ""))
		return
	}
	if match_id <= 0 {
		c.Json(libs.NewError("admincp_common_matchmode_update_fail", "GM006_012", "未指定关联赛事", ""))
		return
	}
	mrs := &matchrace.MatchRaceService{}
	mode := mrs.GetMode(id)
	if mode == nil {
		c.Json(libs.NewError("admincp_common_matchmode_update_fail", "GM006_014", "对象不存在", ""))
		return
	}
	mode.Title = title
	mode.MatchId = match_id
	mode.DisplayOrder = displayorder
	mode.IsView = is_view
	err := mrs.UpdateMode(mode)
	if err == nil {
		c.Json(libs.NewError("admincp_common_matchmode_update_succ", controllers.RESPONSE_SUCCESS, "更新成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_matchmode_update_fail", "GM006_015", err.Error(), ""))
}

// @Title  获取赛程模式
// @Description 获取赛程模式
// @Param   id   path	int true  "id"
// @Success 200 {object} outobjs.OutMatchModeForAdmin
// @router /matchrace/get_mode [get]
func (c *CommonCPController) GetMatchMode() {
	id, _ := c.GetInt64("id")
	mrs := &matchrace.MatchRaceService{}
	mode := mrs.GetMode(id)
	if mode == nil {
		c.Json(libs.NewError("admincp_common_matchmode_get_fail", "GM006_020", "对象不存在", ""))
		return
	}
	out := outobjs.GetOutMatchModeForAdmin(mode)
	c.Json(out)
}

// @Title  获取赛事下的赛程模式
// @Description 获取赛事下的赛程模式
// @Param   match_id   path	int true  "id"
// @Success 200 {object} outobjs.OutMatchMode
// @router /matchrace/get_modes [get]
func (c *CommonCPController) GetMatchModes() {
	match_id, _ := c.GetInt64("match_id")
	if match_id <= 0 {
		c.Json(libs.NewError("admincp_common_matchmode_gets_fail", "GM006_020", "未指定关联赛事", ""))
		return
	}
	mrs := &matchrace.MatchRaceService{}
	modes := mrs.GetModes(match_id)
	outs := []*outobjs.OutMatchMode{}
	for _, mode := range modes {
		outs = append(outs, outobjs.GetOutMatchMode(mode))
	}
	c.Json(outs)
}

// @Title 新增赛程输赢数据
// @Description 新增赛程输赢数据
// @Param   mode_id   path	int true  "模型id"
// @Param   player_id  path	int true  "玩家id"
// @Param   m1   path	int true  "赢输平"
// @Param   m2   path	int true  "赢输平"
// @Param   m3   path	int true  "赢输平"
// @Param   m4   path	int true  "赢输平"
// @Param   m5   path	int true  "赢输平"
// @Param   displayorder   path	int false  "排序"
// @Param   disabled   path	bool false  "关闭"
// @Success 200 {object} libs.Error
// @router /matchrace/recent_add [post]
func (c *CommonCPController) CreateMatchRecent() {
	mode_id, _ := c.GetInt64("mode_id")
	player_id, _ := c.GetInt64("player_id")
	m1, _ := c.GetInt16("m1")
	m2, _ := c.GetInt16("m2")
	m3, _ := c.GetInt16("m3")
	m4, _ := c.GetInt16("m4")
	m5, _ := c.GetInt16("m5")
	displayorder, _ := c.GetInt("displayorder")
	disabled, _ := c.GetBool("disabled")

	if mode_id <= 0 {
		c.Json(libs.NewError("admincp_common_matchrecent_create_fail", "GM006_030", "模型id错误", ""))
		return
	}
	m1t := matchrace.WLP_UNDEFINED
	m2t := matchrace.WLP_UNDEFINED
	m3t := matchrace.WLP_UNDEFINED
	m4t := matchrace.WLP_UNDEFINED
	m5t := matchrace.WLP_UNDEFINED
	if m1 > 0 && m1 < 4 {
		m1t = matchrace.WLP(m1)
	}
	if m2 > 0 && m2 < 4 {
		m2t = matchrace.WLP(m2)
	}
	if m3 > 0 && m3 < 4 {
		m3t = matchrace.WLP(m3)
	}
	if m4 > 0 && m4 < 4 {
		m4t = matchrace.WLP(m4)
	}
	if m5 > 0 && m5 < 4 {
		m5t = matchrace.WLP(m5)
	}
	if player_id <= 0 {
		c.Json(libs.NewError("admincp_common_matchrecent_create_fail", "GM006_031", "玩家id错误", ""))
		return
	}
	recent := &matchrace.MatchRecent{
		ModeId:       mode_id,
		Player:       player_id,
		M1:           m1t,
		M2:           m2t,
		M3:           m3t,
		M4:           m4t,
		M5:           m5t,
		DisplayOrder: displayorder,
		Disabled:     disabled,
	}
	mrs := &matchrace.MatchRaceService{}
	err := mrs.CreateMatchRecent(recent)
	if err == nil {
		c.Json(libs.NewError("admincp_common_matchrecent_create_succ", controllers.RESPONSE_SUCCESS, "创建成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_matchrecent_create_fail", "GM006_032", err.Error(), ""))
}

// @Title 更新赛程输赢数据
// @Description 更新赛程输赢数据
// @Param   id   path	int true  "id"
// @Param   mode_id   path	int true  "模型id"
// @Param   player_id  path	int true  "玩家id"
// @Param   m1   path	int true  "赢输平"
// @Param   m2   path	int true  "赢输平"
// @Param   m3   path	int true  "赢输平"
// @Param   m4   path	int true  "赢输平"
// @Param   m5   path	int true  "赢输平"
// @Param   displayorder   path	int false  "排序"
// @Param   disabled   path	bool false  "关闭"
// @Success 200 {object} libs.Error
// @router /matchrace/recent_update [post]
func (c *CommonCPController) UpdateMatchRecent() {
	id, _ := c.GetInt64("id")
	player_id, _ := c.GetInt64("player_id")
	mode_id, _ := c.GetInt64("mode_id")
	m1, _ := c.GetInt16("m1")
	m2, _ := c.GetInt16("m2")
	m3, _ := c.GetInt16("m3")
	m4, _ := c.GetInt16("m4")
	m5, _ := c.GetInt16("m5")
	displayorder, _ := c.GetInt("displayorder")
	disabled, _ := c.GetBool("disabled")

	m1t := matchrace.WLP_UNDEFINED
	m2t := matchrace.WLP_UNDEFINED
	m3t := matchrace.WLP_UNDEFINED
	m4t := matchrace.WLP_UNDEFINED
	m5t := matchrace.WLP_UNDEFINED
	if m1 > 0 && m1 < 4 {
		m1t = matchrace.WLP(m1)
	}
	if m2 > 0 && m2 < 4 {
		m2t = matchrace.WLP(m2)
	}
	if m3 > 0 && m3 < 4 {
		m3t = matchrace.WLP(m3)
	}
	if m4 > 0 && m4 < 4 {
		m4t = matchrace.WLP(m4)
	}
	if m5 > 0 && m5 < 4 {
		m5t = matchrace.WLP(m5)
	}

	mrs := &matchrace.MatchRaceService{}
	recent := &matchrace.MatchRecent{
		Id:           id,
		ModeId:       mode_id,
		Player:       player_id,
		M1:           m1t,
		M2:           m2t,
		M3:           m3t,
		M4:           m4t,
		M5:           m5t,
		DisplayOrder: displayorder,
		Disabled:     disabled,
	}
	err := mrs.UpdateMatchRecent(recent)
	if err == nil {
		c.Json(libs.NewError("admincp_common_matchrecent_update_succ", controllers.RESPONSE_SUCCESS, "更新成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_matchrecent_update_fail", "GM006_041", err.Error(), ""))
}

// @Title 获取赛程输赢数据
// @Description 获取赛程输赢数据
// @Param   mode_id   path	int true  "模型id"
// @Success 200 {object} outobjs.OutMatchRecentForAdmin
// @router /matchrace/get_recents [get]
func (c *CommonCPController) GetMatchRecents() {
	mode_id, _ := c.GetInt64("mode_id")
	if mode_id <= 0 {
		c.Json(libs.NewError("admincp_common_matchrecent_recents_fail", "GM006_050", "模型id错误", ""))
		return
	}
	mrs := &matchrace.MatchRaceService{}
	recents := mrs.GetMatchRecents(mode_id)
	outs := []*outobjs.OutMatchRecentForAdmin{}
	for _, recent := range recents {
		outs = append(outs, outobjs.GetOutMatchRecentForAdmin(recent))
	}
	c.Json(outs)
}

// @Title 新增赛程小组
// @Description 新增赛程小组
// @Param   title   path	string true  "标题"
// @Param   mode_id   path	int true  "模型id"
// @Param   displayorder   path	int true  "排序"
// @Success 200 {object} libs.Error
// @router /matchrace/group_add [post]
func (c *CommonCPController) CreateMatchGroup() {
	title, _ := utils.UrlDecode(c.GetString("title"))
	mode_id, _ := c.GetInt64("mode_id")
	displayorder, _ := c.GetInt("displayorder")
	if mode_id <= 0 {
		c.Json(libs.NewError("admincp_common_matchgroup_add_fail", "GM006_060", "模型id错误", ""))
		return
	}
	group := &matchrace.MatchGroup{
		ModeId:       mode_id,
		Title:        title,
		DisplayOrder: displayorder,
		PostTime:     time.Now(),
	}
	mrs := &matchrace.MatchRaceService{}
	err := mrs.CreateMatchGroup(group)
	if err == nil {
		c.Json(libs.NewError("admincp_common_matchgroup_add_succ", controllers.RESPONSE_SUCCESS, "创建成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_matchgroup_add_fail", "GM006_061", err.Error(), ""))
}

// @Title 更新赛程小组
// @Description 更新赛程小组
// @Param   id   path	int true  "id"
// @Param   title   path	string true  "标题"
// @Param   displayorder   path	int true  "排序"
// @Success 200 {object} libs.Error
// @router /matchrace/group_update [post]
func (c *CommonCPController) UpdateMatchGroup() {
	title, _ := utils.UrlDecode(c.GetString("title"))
	id, _ := c.GetInt64("id")
	displayorder, _ := c.GetInt("displayorder")
	if id <= 0 {
		c.Json(libs.NewError("admincp_common_matchgroup_update_fail", "GM006_070", "id错误", ""))
		return
	}
	group := &matchrace.MatchGroup{
		Id:           id,
		Title:        title,
		DisplayOrder: displayorder,
		PostTime:     time.Now(),
	}
	mrs := &matchrace.MatchRaceService{}
	err := mrs.UpdateMatchGroup(group)
	if err == nil {
		c.Json(libs.NewError("admincp_common_matchgroup_update_succ", controllers.RESPONSE_SUCCESS, "更新成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_matchgroup_update_fail", "GM006_071", err.Error(), ""))
}

// @Title 获取赛程小组数据
// @Description 获取赛程小组数据
// @Param   mode_id   path	int true  "模型id"
// @Success 200 {object} outobjs.OutMatchGroup
// @router /matchrace/get_groups [get]
func (c *CommonCPController) GetMatchGroups() {
	mode_id, _ := c.GetInt64("mode_id")
	if mode_id <= 0 {
		c.Json(libs.NewError("admincp_common_matchgroup_gets_fail", "GM006_070", "模型id错误", ""))
		return
	}
	mrs := &matchrace.MatchRaceService{}
	mode := mrs.GetMode(mode_id)
	if mode == nil {
		c.Json(libs.NewError("admincp_common_matchgroup_gets_fail", "GM006_071", "模型不存在", ""))
		return
	}
	groups := mrs.GetMatchGroups(mode_id)
	outs := []*outobjs.OutMatchGroup{}
	for _, group := range groups {
		players := mrs.GetMatchGroupPlayers(group.Id)
		outs = append(outs, outobjs.GetOutMatchGroup(mode.MatchId, group, players))
	}
	c.Json(outs)
}

// @Title 新增赛程小组成员
// @Description 新增赛程小组成员
// @Param   group_id   path	int true  "组id"
// @Param   players   path	string true  "成员属性(json)"
// @Success 200 {object} libs.Error
// @router /matchrace/group_players_add [post]
func (c *CommonCPController) CreateMatchGroupPlayers() {
	group_id, _ := c.GetInt64("group_id")
	sPlayer, _ := utils.UrlDecode(c.GetString("players"))

	if group_id <= 0 {
		c.Json(libs.NewError("admincp_common_matchgroupp_create_fail", "GM006_080", "组id错误", ""))
		return
	}
	type PlyerData struct {
		Player       int64 `json:"player"`
		Wins         int16 `json:"wins"`
		Pings        int16 `json:"pings"`
		Loses        int16 `json:"loses`
		Points       int   `json:"points"`
		DisplayOrder int   `json:"displayorder"`
		Outlet       bool  `json:"outlet"`
	}
	var pds []*PlyerData
	err := json.Unmarshal([]byte(sPlayer), pds)
	if err != nil {
		c.Json(libs.NewError("admincp_common_matchgroupp_create_fail", "GM006_081", err.Error(), ""))
		return
	}
	if len(pds) == 0 {
		c.Json(libs.NewError("admincp_common_matchgroup_create_fail", "GM006_082", "添加的成员不能为空", ""))
		return
	}
	mrs := &matchrace.MatchRaceService{}
	for _, p := range pds {
		addobj := &matchrace.MatchGroupPlayer{
			GroupId:      group_id,
			Player:       p.Player,
			Wins:         p.Wins,
			Pings:        p.Pings,
			Loses:        p.Loses,
			Points:       p.Points,
			DisplayOrder: p.DisplayOrder,
			Outlet:       p.Outlet,
		}
		mrs.CreateMatchGroupPlayer(addobj)
	}
	c.Json(libs.NewError("admincp_common_matchgroupp_create_succ", controllers.RESPONSE_SUCCESS, "添加成功", ""))
}

// @Title 更新赛程小组成员
// @Description 更新赛程小组成员
// @Param   id   path	int true  "id"
// @Param   player   path	string true  "成员属性(json)"
// @Success 200 {object} libs.Error
// @router /matchrace/group_players_update [post]
func (c *CommonCPController) UpdateMatchGroupPlayer() {
	id, _ := c.GetInt64("id")
	sPlayer, _ := utils.UrlDecode(c.GetString("players"))

	if id <= 0 {
		c.Json(libs.NewError("admincp_common_matchgroupp_update_fail", "GM006_090", "id错误", ""))
		return
	}
	type PlyerData struct {
		Player       int64 `json:"player"`
		Wins         int16 `json:"wins"`
		Pings        int16 `json:"pings"`
		Loses        int16 `json:"loses`
		Points       int   `json:"points"`
		DisplayOrder int   `json:"displayorder"`
		Outlet       bool  `json:"outlet"`
	}
	var pd *PlyerData
	err := json.Unmarshal([]byte(sPlayer), pd)
	if err != nil || pd == nil {
		c.Json(libs.NewError("admincp_common_matchgroupp_update_fail", "GM006_091", err.Error(), ""))
		return
	}
	mrs := &matchrace.MatchRaceService{}
	p := mrs.GetMatchGroupPlayerForAdmin(id)
	if p == nil {
		c.Json(libs.NewError("admincp_common_matchgroupp_update_fail", "GM006_092", "对象不存在", ""))
		return
	}
	p.Player = pd.Player
	p.Wins = pd.Wins
	p.Pings = pd.Pings
	p.Loses = pd.Loses
	p.Points = pd.Points
	p.DisplayOrder = pd.DisplayOrder
	p.Outlet = pd.Outlet
	err = mrs.UpdateMatchGroupPlayer(p)
	if err == nil {
		c.Json(libs.NewError("admincp_common_matchgroupp_update_succ", controllers.RESPONSE_SUCCESS, "更新成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_matchgroupp_update_fail", "GM006_093", err.Error(), ""))
}

// @Title 新增晋级模型组
// @Description 新增晋级模型组
// @Param   mode_id   path	int true  "模型id"
// @Param   title   path	string true  "标题"
// @Param   icon   path	int true  "图标id"
// @Param   displayorder   path	int true  "排序"
// @Param   t   path	int false  "冠亚军"
// @Success 200 {object} libs.Error
// @router /matchrace/elimin_ms_add [post]
func (c *CommonCPController) CreateMatchEliminMs() {
	mode_id, _ := c.GetInt64("mode_id")
	title, _ := utils.UrlDecode(c.GetString("title"))
	icon, _ := c.GetInt64("icon")
	displayorder, _ := c.GetInt("displayorder")
	t, _ := c.GetInt16("t")

	if mode_id <= 0 {
		c.Json(libs.NewError("admincp_common_matchelimin_add_fail", "GM006_100", "模型id错误", ""))
		return
	}
	if len(title) == 0 {
		c.Json(libs.NewError("admincp_common_matchelimin_add_fail", "GM006_100", "标题不能为空", ""))
		return
	}
	mstype := matchrace.ELIMIN_MSTYPE(t)
	mrs := &matchrace.MatchRaceService{}
	err := mrs.CreateEliminMs(&matchrace.MatchEliminMs{
		ModeId:       mode_id,
		Title:        title,
		PostTime:     time.Now().Unix(),
		Icon:         icon,
		DisplayOrder: displayorder,
		T:            mstype,
	})
	if err == nil {
		c.Json(libs.NewError("admincp_common_matchelimin_add_succ", controllers.RESPONSE_SUCCESS, "添加成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_matchelimin_add_fail", "GM006_101", err.Error(), ""))
}

// @Title 更新晋级模型组
// @Description 更新晋级模型组
// @Param   id   path	int true  "id"
// @Param   title   path	string true  "标题"
// @Param   icon   path	int true  "图标id"
// @Param   displayorder   path	int true  "排序"
// @Param   t   path	int false  "冠亚军"
// @Success 200 {object} libs.Error
// @router /matchrace/elimin_ms_update [post]
func (c *CommonCPController) UpdateEliminMs() {
	id, _ := c.GetInt64("id")
	title, _ := utils.UrlDecode(c.GetString("title"))
	icon, _ := c.GetInt64("icon")
	displayorder, _ := c.GetInt("displayorder")
	t, _ := c.GetInt16("t")

	if id <= 0 {
		c.Json(libs.NewError("admincp_common_matchelimin_update_fail", "GM006_110", "模型id错误", ""))
		return
	}
	if len(title) == 0 {
		c.Json(libs.NewError("admincp_common_matchelimin_update_fail", "GM006_111", "标题不能为空", ""))
		return
	}
	mstype := matchrace.ELIMIN_MSTYPE(t)
	mrs := &matchrace.MatchRaceService{}
	ms := mrs.GetEliminMsForAdmin(id)
	if ms == nil {
		c.Json(libs.NewError("admincp_common_matchelimin_update_fail", "GM006_112", "对象不存在", ""))
		return
	}
	ms.Title = title
	ms.Icon = icon
	ms.DisplayOrder = displayorder
	ms.T = mstype
	err := mrs.UpdateEliminMs(ms)
	if err == nil {
		c.Json(libs.NewError("admincp_common_matchelimin_update_succ", controllers.RESPONSE_SUCCESS, "添加成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_matchelimin_update_fail", "GM006_113", err.Error(), ""))
}

// @Title 删除晋级模型组
// @Description 删除晋级模型组
// @Param   id   path	int true  "id"
// @Success 200 {object} libs.Error
// @router /matchrace/elimin_ms_del [delete]
func (c *CommonCPController) DeleteEliminMs() {
	id, _ := c.GetInt64("id")
	mrs := &matchrace.MatchRaceService{}
	err := mrs.DeleteEliminMs(id)
	if err == nil {
		c.Json(libs.NewError("admincp_common_matchelimin_del_succ", controllers.RESPONSE_SUCCESS, "删除成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_matchelimin_del_fail", "GM006_120", err.Error(), ""))
}

// @Title 获取赛程小组数据
// @Description 获取赛程小组数据
// @Param   mode_id   path	int true  "模型id"
// @Success 200 {object} outobjs.OutMatchEliminMs
// @router /matchrace/get_eliminmss [get]
func (c *CommonCPController) GetEliminMss() {
	mode_id, _ := c.GetInt64("mode_id")
	if mode_id <= 0 {
		c.Json(libs.NewError("admincp_common_matchelimin_gets_fail", "GM006_120", "id错误", ""))
		return
	}
	mrs := &matchrace.MatchRaceService{}
	mss := mrs.GetEliminMss(mode_id)
	out_mss := []*outobjs.OutMatchEliminMs{}
	for _, ms := range mss {
		out_mss = append(out_mss, outobjs.GetOutMatchEliminMs(ms))
	}
	c.Json(out_mss)
}

// @Title 新增晋级模型组对阵
// @Description 新增晋级模型组对阵
// @Param   msid   path	int true  "模型id"
// @Param   vsid   path	int true  "对阵vsid"
// @Param   outletid   path	int true  "出线选手id"
// @Success 200 {object} libs.Error
// @router /matchrace/elimin_vs_add [post]
func (c *CommonCPController) CreateMatchEliminVs() {
	msid, _ := c.GetInt64("msid")
	vsid, _ := c.GetInt64("vsid")
	outletid, _ := c.GetInt64("outletid")
	if msid <= 0 || vsid <= 0 {
		c.Json(libs.NewError("admincp_common_matcheliminvs_add_fail", "GM006_130", "参数错误", ""))
		return
	}
	mrs := &matchrace.MatchRaceService{}
	err := mrs.CreateEliminVs(&matchrace.MatchEliminVs{
		MsId:     msid,
		VsId:     vsid,
		OutletId: outletid,
		PostTime: time.Now().Unix(),
	})
	if err == nil {
		c.Json(libs.NewError("admincp_common_matcheliminvs_add_fail", controllers.RESPONSE_SUCCESS, "添加成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_matcheliminvs_add_fail", "GM006_131", err.Error(), ""))
}

// @Title 更新晋级模型组对阵
// @Description 更新晋级模型组对阵
// @Param   id   path	int true  "id"
// @Param   msid   path	int true  "模型id"
// @Param   vsid   path	int true  "对阵vsid"
// @Param   outletid   path	int true  "出线选手id"
// @Success 200 {object} libs.Error
// @router /matchrace/elimin_vs_update [post]
func (c *CommonCPController) UpdateMatchEliminVs() {
	id, _ := c.GetInt64("id")
	msid, _ := c.GetInt64("msid")
	vsid, _ := c.GetInt64("vsid")
	outletid, _ := c.GetInt64("outletid")
	mrs := &matchrace.MatchRaceService{}
	err := mrs.UpdateEliminVs(&matchrace.MatchEliminVs{
		Id:       id,
		MsId:     msid,
		VsId:     vsid,
		OutletId: outletid,
		PostTime: time.Now().Unix(),
	})
	if err == nil {
		c.Json(libs.NewError("admincp_common_matcheliminvs_update_fail", controllers.RESPONSE_SUCCESS, "更新成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_matcheliminvs_update_fail", "GM006_140", err.Error(), ""))
}

// @Title 删除晋级模型组对阵
// @Description 删除晋级模型组对阵
// @Param   id   path	int true  "id"
// @Success 200 {object} libs.Error
// @router /matchrace/elimin_vs_del [delete]
func (c *CommonCPController) DeleteMatchEliminVs() {
	id, _ := c.GetInt64("id")
	mrs := &matchrace.MatchRaceService{}
	err := mrs.DeleteEliminVs(id)
	if err == nil {
		c.Json(libs.NewError("admincp_common_matcheliminvs_del_succ", controllers.RESPONSE_SUCCESS, "删除成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_matcheliminvs_del_fail", "GM006_150", err.Error(), ""))
}

// @Title 新增对阵数据
// @Description 新增对阵数据
// @Param   a_id   path	int true  "aid"
// @Param   a_img   path	int true  "a图片"
// @Param   a_name   path	string true  "a名称"
// @Param   a_score   path	int false  "a比分"
// @Param   b_id   path	int true  "bid"
// @Param   b_img   path	int true  "b图片"
// @Param   b_name   path	string true  "b名称"
// @Param   b_score   path	int false  "b比分"
// @Param   match_id   path	int true  "赛事id"
// @Param   mode_id   path	int true  "模型id"
// @Param   ref_id   path	int true  "关联id"
// @Success 200 {object} libs.Error
// @router /matchrace/matchvs_add [post]
func (c *CommonCPController) CreateMatchVs() {
	aid, _ := c.GetInt64("a_id")
	aimg, _ := c.GetInt64("a_img")
	aname, _ := utils.UrlDecode(c.GetString("a_name"))
	ascore, _ := c.GetInt16("a_score")
	bid, _ := c.GetInt64("b_id")
	bimg, _ := c.GetInt64("b_img")
	bname, _ := utils.UrlDecode(c.GetString("b_name"))
	bscore, _ := c.GetInt16("b_score")
	matchid, _ := c.GetInt64("match_id")
	modeid, _ := c.GetInt64("mode_id")
	refid, _ := c.GetInt64("ref_id")
	if aimg <= 0 || len(aname) <= 0 || bimg <= 0 || len(bname) <= 0 {
		c.Json(libs.NewError("admincp_common_matchvs_add_fail", "GM006_160", "参数错误", ""))
		return
	}
	if matchid <= 0 || modeid <= 0 {
		c.Json(libs.NewError("admincp_common_matchvs_add_fail", "GM006_161", "参数错误", ""))
		return
	}
	mrs := &matchrace.MatchRaceService{}
	err := mrs.CreateMatchVs(&matchrace.MatchVs{
		A:       aid,
		AName:   aname,
		AImg:    aimg,
		AScore:  ascore,
		B:       bid,
		BName:   bname,
		BImg:    bimg,
		BScore:  bscore,
		MatchId: matchid,
		ModeId:  modeid,
		RefId:   refid,
	})
	if err == nil {
		c.Json(libs.NewError("admincp_common_matchvs_add_succ", controllers.RESPONSE_SUCCESS, "添加成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_matchvs_add_fail", "GM006_162", err.Error(), ""))
}

// @Title 更新对阵数据
// @Description 更新对阵数据
// @Param   id   path	int true  "id"
// @Param   a_id   path	int true  "aid"
// @Param   a_img   path	int true  "a图片"
// @Param   a_name   path	string true  "a名称"
// @Param   a_score   path	int false  "a比分"
// @Param   b_id   path	int true  "bid"
// @Param   b_img   path	int true  "b图片"
// @Param   b_name   path	string true  "b名称"
// @Param   b_score   path	int false  "b比分"
// @Param   match_id   path	int true  "赛事id"
// @Param   mode_id   path	int true  "模型id"
// @Param   ref_id   path	int true  "关联id"
// @Success 200 {object} libs.Error
// @router /matchrace/matchvs_add [post]
func (c *CommonCPController) UpdateMatchVs() {
	id, _ := c.GetInt64("id")
	aid, _ := c.GetInt64("a_id")
	aimg, _ := c.GetInt64("a_img")
	aname, _ := utils.UrlDecode(c.GetString("a_name"))
	ascore, _ := c.GetInt16("a_score")
	bid, _ := c.GetInt64("b_id")
	bimg, _ := c.GetInt64("b_img")
	bname, _ := utils.UrlDecode(c.GetString("b_name"))
	bscore, _ := c.GetInt16("b_score")
	matchid, _ := c.GetInt64("match_id")
	modeid, _ := c.GetInt64("mode_id")
	refid, _ := c.GetInt64("ref_id")
	if aimg <= 0 || len(aname) <= 0 || bimg <= 0 || len(bname) <= 0 {
		c.Json(libs.NewError("admincp_common_matchvs_update_fail", "GM006_170", "参数错误", ""))
		return
	}
	if matchid <= 0 || modeid <= 0 {
		c.Json(libs.NewError("admincp_common_matchvs_update_fail", "GM006_171", "参数错误", ""))
		return
	}
	mrs := &matchrace.MatchRaceService{}
	vs := mrs.GetMatchVs(id)
	if vs == nil {
		c.Json(libs.NewError("admincp_common_matchvs_update_fail", "GM006_172", "参数错误", ""))
		return
	}
	vs.A = aid
	vs.AName = aname
	vs.AImg = aimg
	vs.AScore = ascore
	vs.B = bid
	vs.BName = bname
	vs.BImg = bimg
	vs.BScore = bscore
	vs.MatchId = matchid
	vs.ModeId = modeid
	vs.RefId = refid
	err := mrs.UpdateMatchVs(vs)
	if err == nil {
		c.Json(libs.NewError("admincp_common_matchvs_update_succ", controllers.RESPONSE_SUCCESS, "更新成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_matchvs_update_fail", "GM006_162", err.Error(), ""))
}

// @Title 获取对阵数据
// @Description 获取对阵数据
// @Param   id   path	int true  "id"
// @Success 200 {object} outobjs.OutMatchVs
// @router /matchrace/matchvs [get]
func (c *CommonCPController) GetMatchVs() {
	id, _ := c.GetInt64("id")
	if id <= 0 {
		c.Json(libs.NewError("admincp_common_matchvs_get_fail", "GM006_180", "参数错误", ""))
		return
	}
	mrs := &matchrace.MatchRaceService{}
	vs := mrs.GetMatchVs(id)
	if vs == nil {
		c.Json(libs.NewError("admincp_common_matchvs_get_fail", "GM006_181", "对象不存在", ""))
		return
	}
	c.Json(outobjs.GetOutMatchVs(vs))
}

// @Title 获取对阵数据列表
// @Description 获取对阵数据列表
// @Param   player   path	int true  "选手id"
// @Param   match_id   path	int true  "赛程id"
// @Success 200 {object} outobjs.OutMatchVs
// @router /matchrace/matchvss [get]
func (c *CommonCPController) GetMatchVss() {
	palyerid, _ := c.GetInt64("player")
	matchid, _ := c.GetInt64("match_id")
	if palyerid <= 0 || matchid <= 0 {
		c.Json(libs.NewError("admincp_common_matchvss_get_fail", "GM006_180", "参数错误", ""))
		return
	}
	mrs := &matchrace.MatchRaceService{}
	vss := mrs.GetMatchVss(palyerid, matchid)
	out_vss := []*outobjs.OutMatchVs{}
	for _, vs := range vss {
		out_vss = append(out_vss, outobjs.GetOutMatchVs(vs))
	}

	c.Json(out_vss)
}

// @Title 新增选手资料
// @Description 新增选手资料
// @Param   player_name   path	int true  "选手名称"
// @Param   img_id   path	int true  "图片id"
// @Success 200 {object} libs.Error
// @router /matchrace/player_add [post]
func (c *CommonCPController) CreatePlayer() {
	pname, _ := utils.UrlDecode(c.GetString("player_name"))
	imgid, _ := c.GetInt64("img_id")
	if len(pname) <= 0 || imgid <= 0 {
		c.Json(libs.NewError("admincp_common_player_add_fail", "GM006_190", "参数错误", ""))
		return
	}
	mrs := &matchrace.MatchRaceService{}
	err := mrs.CreatePlayer(&matchrace.MatchPlayer{
		Name:     pname,
		Img:      imgid,
		PostTime: time.Now(),
	})
	if err == nil {
		c.Json(libs.NewError("admincp_common_player_add_succ", controllers.RESPONSE_SUCCESS, "添加成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_player_add_fail", "GM006_191", err.Error(), ""))
}

// @Title 更新选手资料
// @Description 更新选手资料
// @Param   id   path	int true  "选手id"
// @Param   player_name   path	int true  "选手名称"
// @Param   img_id   path	int true  "图片id"
// @Success 200 {object} libs.Error
// @router /matchrace/player_update [post]
func (c *CommonCPController) UpdatePlayer() {
	id, _ := c.GetInt64("id")
	pname, _ := utils.UrlDecode(c.GetString("player_name"))
	imgid, _ := c.GetInt64("img_id")
	if id <= 0 || len(pname) <= 0 || imgid <= 0 {
		c.Json(libs.NewError("admincp_common_player_update_fail", "GM006_200", "参数错误", ""))
		return
	}
	mrs := &matchrace.MatchRaceService{}
	player := mrs.GetPlayer(id)
	if player == nil {
		c.Json(libs.NewError("admincp_common_player_update_fail", "GM006_201", "对象不存在", ""))
		return
	}
	player.Name = pname
	player.Img = imgid
	err := mrs.UpdatePlayer(player)
	if err == nil {
		c.Json(libs.NewError("admincp_common_player_update_succ", controllers.RESPONSE_SUCCESS, "添加成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_common_player_update_fail", "GM006_202", err.Error(), ""))
}

// @Title 获取选手资料
// @Description 获取选手资料
// @Param   id   path	int true  "选手id"
// @Success 200 {object} outobjs.OutMatchPlayer
// @router /matchrace/getplayer [get]
func (c *CommonCPController) GetMatchPlayer() {
	id, _ := c.GetInt64("id")
	mrs := &matchrace.MatchRaceService{}
	player := mrs.GetPlayer(id)
	if player == nil {
		c.Json(libs.NewError("admincp_common_player_get_fail", "GM006_210", "对象不存在", ""))
		return
	}
	c.Json(outobjs.GetOutMatchPlayer(player))
}

// @Title 获取选手们资料
// @Description 获取选手们资料
// @Param   like   path	string false  "关键字"
// @Param   page   path	int true  "页"
// @Param   size   path	int true  "页大小"
// @Success 200 {object} outobjs.OutMatchPlayerPagedList
// @router /matchrace/getplayers [get]
func (c *CommonCPController) GetMatchPlayers() {
	like, _ := utils.UrlDecode(c.GetString("like"))
	page, _ := c.GetInt("page")
	size, _ := c.GetInt("size")
	mrs := &matchrace.MatchRaceService{}
	total, players := mrs.GetPlayers(like, page, size)
	outs := []*outobjs.OutMatchPlayer{}
	for _, player := range players {
		outs = append(outs, outobjs.GetOutMatchPlayer(player))
	}
	outsp := &outobjs.OutMatchPlayerPagedList{
		Total:       total,
		TotalPage:   utils.TotalPages(total, size),
		CurrentPage: page,
		Size:        size,
		Lists:       outs,
	}
	c.Json(outsp)
}
