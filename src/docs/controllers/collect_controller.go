package controllers

import (
	"libs"
	"libs/collect"
	"libs/lives"
	"libs/vod"
	"outobjs"
	"strconv"
	"strings"
	"time"
	"utils"
)

// 收藏模块 API
type CollectController struct {
	AuthorizeController
}

func (c *CollectController) Prepare() {
	c.AuthorizeController.Prepare()
}

func (c *CollectController) URLMapping() {
	c.Mapping("Add", c.Add)
	c.Mapping("Remove", c.Remove)
	c.Mapping("Show", c.Show)
	c.Mapping("Gets", c.Gets)
}

// @Title 添加收藏
// @Description 添加收藏(成功返回error_code:REP000,error_description:收藏顺序号)
// @Param   access_token  path  string  true  "access_token"
// @Param   mod   path  string  true  "所属模块(vod:视频,j_live:机构直播,p_live:个人直播)"
// @Param   ref_id    path  int  true  "关联id(比如视频id号等)"
// @Success 200 {object} libs.Error
// @router /add [post]
func (c *CollectController) Add() {
	uid := c.CurrentUid()
	if uid <= 0 {
		c.Json(libs.NewError("member_collect_premission_denied", UNAUTHORIZED_CODE, "必须登录后才能收藏", ""))
		return
	}
	mod := c.GetString("mod")
	refid := c.GetString("ref_id")
	if len(mod) == 0 {
		c.Json(libs.NewError("member_collect_add_fail", "C7001", "mod参数不能为空", ""))
		return
	}
	if len(refid) == 0 {
		c.Json(libs.NewError("member_collect_add_fail", "C7002", "ref_id参数不能为空", ""))
		return
	}
	err, id := collect.Create(uid, refid, mod)
	if err != nil {
		c.Json(libs.NewError("member_collect_add_fail", "C7005", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("member_collect_add_succ", RESPONSE_SUCCESS, id, ""))
}

// @Title 删除收藏
// @Description 删除收藏(成功返回error_code:REP000)
// @Param   access_token  path  string  true  "access_token"
// @Param   id   path  string  true  "需要删除的收藏id"
// @Success 200 {object} libs.Error
// @router /remove [post]
func (c *CollectController) Remove() {
	id := c.GetString("id")
	uid := c.CurrentUid()
	if uid <= 0 {
		c.Json(libs.NewError("member_collect_premission_denied", UNAUTHORIZED_CODE, "必须登录后才能移除收藏", ""))
		return
	}
	if len(id) == 0 {
		c.Json(libs.NewError("member_collect_remove_fail", "C7101", "id参数不能为空", ""))
		return
	}
	err := collect.Delete(uid, id)
	if err != nil {
		c.Json(libs.NewError("member_collect_remove_fail", "C7102", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("member_collect_remove_succ", RESPONSE_SUCCESS, "成功移除收藏", ""))
}

// @Title 批量删除收藏
// @Description 删除收藏(成功返回error_code:REP000)
// @Param   access_token  path  string  true  "access_token"
// @Param   ids   path  string  true  "需要删除的收藏ids,用,隔开"
// @Success 200 {object} libs.Error
// @router /removes [post]
func (c *CollectController) Removes() {
	idstr := c.GetString("ids")
	uid := c.CurrentUid()
	if uid <= 0 {
		c.Json(libs.NewError("member_collect_premission_denied", UNAUTHORIZED_CODE, "必须登录后才能移除收藏", ""))
		return
	}
	if len(idstr) == 0 {
		c.Json(libs.NewError("member_collect_removes_fail", "C7201", "ids参数不能为空", ""))
		return
	}
	ids := strings.Split(idstr, ",")
	collect.Deletes(uid, ids)
	c.Json(libs.NewError("member_collect_removes_succ", RESPONSE_SUCCESS, "成功移除收藏", ""))
}

// @Title 查询是否收藏
// @Description 查询是否收藏(已收藏返回error_code:REP000,未收藏REP001,description:收藏号)
// @Param   access_token  path  string  true  "access_token"
// @Param   mod   path  string  true  "所属模块(vod:视频,j_live:机构直播,p_live:个人直播)"
// @Param   ref_id    path  int  true  "关联id(比如视频id号等)"
// @Success 200 {object} libs.Error
// @router /show [get]
func (c *CollectController) Show() {
	uid := c.CurrentUid()
	if uid <= 0 {
		c.Json(libs.NewError("member_collect_premission_denied", UNAUTHORIZED_CODE, "必须登录后才能查询收藏状态", ""))
		return
	}
	mod := c.GetString("mod")
	refid := c.GetString("ref_id")
	if len(mod) == 0 {
		c.Json(libs.NewError("member_collect_show_fail", "C7301", "mod参数不能为空", ""))
		return
	}
	if len(refid) == 0 {
		c.Json(libs.NewError("member_collect_show_fail", "C7302", "ref_id参数不能为空", ""))
		return
	}
	ed := collect.IsCollcetd(uid, refid, mod)
	if !ed {
		c.Json(libs.NewError("member_collect_show_unfav", REPSONSE_FAIL, "未收藏", ""))
		return
	}
	fav := collect.Get(uid, refid, mod)
	if fav == nil {
		c.Json(libs.NewError("member_collect_show_unfav", REPSONSE_FAIL, "未收藏", ""))
		return
	}
	c.Json(libs.NewError("member_collect_show_faved", RESPONSE_SUCCESS, fav.Id.Hex(), ""))
}

// @Title 获取收藏列表
// @Description 获取收藏列表
// @Param   access_token  path  string  true  "access_token"
// @Param   size    path  int  false  "页行数(默认20)"
// @Param   ts    path  int  false  "时间戳(每次请求获得的t属性)"
// @Success 200 {object} outobjs.OutCollectiblePageList
// @router /list [get]
func (c *CollectController) Gets() {
	uid := c.CurrentUid()
	page, _ := c.GetInt("page")
	size, _ := c.GetInt("size")
	timestamp, _ := c.GetInt64("ts")

	t := time.Now()
	if timestamp > 0 {
		t = time.Unix(timestamp, 0)
	}
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	total, collects := collect.Gets(uid, int(page), int(size), t)
	list := []*outobjs.OutCollectible{}
	for _, cc := range collects {
		out_c := &outobjs.OutCollectible{
			Id:             cc.Id.Hex(),
			RelId:          cc.RelId,
			RelType:        cc.RelType,
			PreviewContent: cc.PreviewContent,
			PreviewImgId:   cc.PreviewImg,
			PreviewImgUrl:  file_storage.GetFileUrl(cc.PreviewImg),
			CreateTime:     cc.CreateTime,
		}
		c.transformObj(cc.RelId, cc.RelType, out_c)
		list = append(list, out_c)
		if cc.CreateTime.Before(t) {
			t = cc.CreateTime
		}
	}
	out_list := outobjs.OutCollectiblePageList{
		Total:       total,
		TotalPage:   utils.TotalPages(int(total), int(size)),
		CurrentPage: int(page),
		Size:        int(size),
		Time:        t.Unix(),
		Lists:       list,
	}
	c.Json(out_list)
}

var vst *vod.Vods = &vod.Vods{}
var pvt *lives.LivePers = &lives.LivePers{}
var jvt *lives.LiveOrgs = &lives.LiveOrgs{}

//vod:视频,j_live:机构直播,p_live:个人直播
func (c *CollectController) transformObj(refId string, refType string, out *outobjs.OutCollectible) {
	switch refType {
	case "vod":
		_id, err := strconv.ParseInt(refId, 10, 64)
		if err == nil {
			v := vst.Get(_id, false)
			if v != nil {
				out.Vod = *outobjs.GetOutVideoInfo(v)
			}
		}
		break
	case "p_live":
		_id, err := strconv.ParseInt(refId, 10, 64)
		if err == nil {
			p := pvt.Get(_id)
			if p != nil {
				out.Per = *outobjs.GetOutPersonalLive(p)
			}
		}
		break
	case "j_live":
		_id, err := strconv.ParseInt(refId, 10, 64)
		if err == nil {
			p := jvt.GetChannel(_id)
			if p != nil {
				out.Org = *outobjs.GetOutLiveChannel(p, 0)
			}
		}
		break
	}
}
