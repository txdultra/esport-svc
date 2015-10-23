package admincp

import (
	"controllers"
	"libs"
	//"libs/passport"
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
// @Success 200 {object} libs.Error
// @router /game/add [post]
func (c *CommonCPController) GameAdd() {
	name, _ := utils.UrlDecode(c.GetString("name"))
	en, _ := utils.UrlDecode(c.GetString("en"))
	img_id, _ := c.GetInt64("img_id")
	enabled, _ := c.GetBool("enabled")
	display_order, _ := c.GetInt64("display_order")
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
// @Param   display_order   path	int false  "排序"
// @Success 200 {object} libs.Error
// @router /game/update [post]
func (c *CommonCPController) GameUpdate() {
	id, _ := c.GetInt64("id")
	name, _ := utils.UrlDecode(c.GetString("name"))
	en, _ := utils.UrlDecode(c.GetString("en"))
	img_id, _ := c.GetInt64("img_id")
	enabled, _ := c.GetBool("enabled")
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
		c.Json(libs.NewError("admincp_common_homead_dle_fail", "GM004_070", "id非法", ""))
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
