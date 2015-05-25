package admincp

import (
	"controllers"
	"fmt"
	"libs"
	"libs/groups"
	"outobjs"
	"strconv"
	"strings"
	"utils"
)

// 组管理 API
type GroupCPController struct {
	AdminController
}

func (c *GroupCPController) Prepare() {
	c.AdminController.Prepare()
}

// @Title 获取小组列表
// @Description 获取小组列表
// @Param   page   path  int  false  "页"
// @Param   size   path  int  false  "页数量"
// @Param   words   path  int  false  "搜索关键字"
// @Param   orderby   path  string  false  "排序规则(recommend默认,hot,fans,official)"
// @Success 200  {object} outobjs.OutGroupPagedList
// @router /group/search [get]
func (c *GroupCPController) GroupSearch() {
	words, _ := utils.UrlDecode(c.GetString("words"))
	page, _ := c.GetInt("page")
	size, _ := c.GetInt("size")
	orderby := c.GetString("orderby")
	var sortBy groups.GP_SEARCH_SORT
	switch orderby {
	case "hot":
		sortBy = groups.GP_SEARCH_SORT_VITALITY
		break
	case "fans":
		sortBy = groups.GP_SEARCH_SORT_USERS
		break
	case "official":
		sortBy = groups.GP_SEARCH_SORT_OFFICIAL
	default:
		sortBy = groups.GP_SEARCH_SORT_DEFAULT
	}

	match_mode := "any"
	if len(words) == 0 {
		match_mode = "all"
	}
	if size <= 0 {
		size = 20
	}
	cfg := groups.GetDefaultCfg()
	gs := groups.NewGroupService(cfg)
	total, groups := gs.Search(words, page, size, match_mode, sortBy, nil, nil)
	out_groups := outobjs.GetOutGroups(groups, 0)
	out_p := &outobjs.OutGroupPagedList{
		CurrentPage: page,
		PageSize:    size,
		TotalPages:  utils.TotalPages(total, size),
		Groups:      out_groups,
		Total:       total,
	}
	c.Json(out_p)
}

// @Title 获取组信息
// @Description 获取组信息
// @Param   group_id  path  int  true  "组id"
// @Success 200 {object} outobjs.OutGroup
// @router /group/get [get]
func (c *GroupCPController) GetGroup() {
	groupid, _ := c.GetInt64("group_id")
	if groupid <= 0 {
		c.Json(libs.NewError("admincp_group_get_fail", "GM040_001", "group_id参数错误", ""))
		return
	}
	gs := groups.NewGroupService(groups.GetDefaultCfg())
	group := gs.Get(groupid)
	if group == nil {
		c.Json(libs.NewError("admincp_group_get_notexist", "GM040_002", "组不存在", ""))
		return
	}
	c.Json(outobjs.GetOutGroup(group, 0))
}

// @Title 建组
// @Description 建组
// @Param   name   path  string  true  "名称"
// @Param   description   path  string  true  "描述"
// @Param   country  path  string  false  "国家"
// @Param   city  path  string  false  "城市"
// @Param   game_ids  path  string  false  "选择游戏(逗号,分隔)"
// @Param   bgimg   path  int  true  "背景图片id"
// @Param   longitude   path  float  false  "经度"
// @Param   latitude   path  float  false  "维度"
// @Param   invite_uids   path  string  false  "用户uids(用,分隔)"
// @Param   belong   path  string  true  "所属类型"
// @Param   uid  path  int  true  "所属uid"
// @Success 200 {object} libs.Error
// @router /group/create [post]
func (c *GroupCPController) CreateGroup() {
	name, _ := utils.UrlDecode(c.GetString("name"))
	desc, _ := utils.UrlDecode(c.GetString("description"))
	country, _ := utils.UrlDecode(c.GetString("country"))
	city, _ := utils.UrlDecode(c.GetString("city"))
	gameids := c.GetString("game_ids")
	bgimg, _ := c.GetInt64("bgimg")
	longitude, _ := c.GetFloat("longitude")
	latitude, _ := c.GetFloat("latitude")
	inviteuids := c.GetString("invite_uids")
	uid, _ := c.GetInt64("uid")
	bl, _ := c.GetInt("belong")

	var belong groups.GROUP_BELONG
	switch bl {
	case int(groups.GROUP_BELONG_MEMBER):
		belong = groups.GROUP_BELONG_MEMBER
		break
	case int(groups.GROUP_BELONG_OFFICIAL):
		belong = groups.GROUP_BELONG_OFFICIAL
		break
	default:
		c.Json(libs.NewError("admincp_group_create_belong_fail", "GM040_010", "所属类型错误", ""))
		return
	}

	inv_uids := []int64{}
	inv_uidss := strings.Split(inviteuids, ",")
	for _, str := range inv_uidss {
		_id, _ := strconv.ParseInt(str, 10, 64)
		if _id > 0 {
			inv_uids = append(inv_uids, _id)
		}
	}

	gs := groups.NewGroupService(groups.GetDefaultCfg())
	group := &groups.Group{
		Name:        name,
		Description: desc,
		Uid:         uid,
		Country:     country,
		City:        city,
		GameIds:     gameids,
		BgImg:       bgimg,
		Belong:      belong,
		Type:        groups.GROUP_TYPE_NORMAL,
		LongiTude:   float32(longitude),
		LatiTude:    float32(latitude),
		InviteUids:  inv_uids,
	}
	err := gs.Create(group)
	if err != nil {
		c.Json(libs.NewError("admincp_group_create_fail", "GM040_011", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("admincp_group_create_success", controllers.RESPONSE_SUCCESS, fmt.Sprintf("%d", group.Id), ""))
}
