package controllers

import (
	"libs"
	"libs/passport"
	"outobjs"
	//"strings"
	//"fmt"
	"time"
	"utils"
)

// 关系模块 API
type FriendShipsController struct {
	AuthorizeController
}

func (c *FriendShipsController) Prepare() {
	c.AuthorizeController.Prepare()
}

func (c *FriendShipsController) URLMapping() {
	c.Mapping("Friends", c.Friends)
	c.Mapping("FriendsP", c.FriendsP)
	c.Mapping("Followers", c.Followers)
	c.Mapping("Show", c.Show)
	c.Mapping("Create", c.Create)
	c.Mapping("Destroy", c.Destroy)
	c.Mapping("Recmds", c.Recmds)
}

// @Title 用户关注列表
// @Description 获取用户的关注列表(拼音首字母列表)
// @Param   access_token   path   string  true  "access_token"
// @Param   uid   path   int  false  "查询的用户uid(填空为查询自己)"
// @Success 200
// @router /friends/all [get]
func (c *FriendShipsController) Friends() {
	uid := c.CurrentUid()
	s_uid, _ := c.GetInt64("uid")
	var f_uid int64 = 0
	if s_uid > 0 && uid != s_uid {
		f_uid = s_uid
	} else {
		f_uid = uid
	}

	fship := &passport.FriendShips{}
	pymaps := fship.FriendIdsByNickNameFirstPY(f_uid)
	type OutPy struct {
		Key     string               `json:"w"`
		Members []*outobjs.OutMember `json:"members"`
	}
	outs := []*OutPy{}
	for _, v := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ*" {
		vs, ok := pymaps[string(v)]
		if ok {
			out_members := []*outobjs.OutMember{}
			for _, _uid := range vs {
				out_member := outobjs.GetOutMember(_uid, uid)
				if out_member != nil {
					out_members = append(out_members, out_member)
				}
			}
			outs = append(outs, &OutPy{
				Key:     string(v),
				Members: out_members,
			})
		}
	}
	c.Json(outs)
}

// @Title 用户关注列表
// @Description 获取用户的关注列表(分页形式)
// @Param   access_token   path   string  true  "access_token"
// @Param   uid   path   int  false  "查看的用户关注表(空为查询自己)"
// @Param   size     path    int  false        "页大小(最大50,默认20)"
// @Param   page     path    int  false        "页码(默认1)"
// @Param   t    path  int  false  "时间戳(每次请求获得的t属性)"
// @Success 200  {object} outobjs.OutFriendList
// @router /friends [get]
func (c *FriendShipsController) FriendsP() {
	current_uid := c.CurrentUid()
	page, _ := c.GetInt("page")
	size, _ := c.GetInt("size")
	uid, _ := c.GetInt64("uid")
	timestamp, _ := c.GetInt64("t")

	t := time.Now()
	if timestamp > 0 {
		t = utils.MsToTime(timestamp) //time.Unix(timestamp, 0)
	}

	show_uid := current_uid
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if uid > 0 {
		show_uid = uid
	}
	fship := &passport.FriendShips{}
	total, uids, last_t := fship.FriendIdsP(show_uid, int(size), int(page), t)
	users := []*outobjs.OutMember{}
	for _, _uid := range uids {
		m := outobjs.GetOutMember(_uid, current_uid)
		users = append(users, m)
	}
	out := outobjs.OutFriendList{
		Total:       total,
		TotalPage:   utils.TotalPages(int(total), int(size)),
		CurrentPage: int(page),
		Size:        int(size),
		Users:       users,
		Time:        last_t,
	}
	c.Json(out)
}

// @Title 用户粉丝列表
// @Description 获取用户粉丝列表(自动清空计数)
// @Param   access_token   path   string  true  "access_token"
// @Param   uid   path   int  false  "查看的用户关注表(空为查询自己)"
// @Param   size     path    int  false        "页大小(最大50,默认20)"
// @Param   page     path    int  false        "页码(默认1)"
// @Param   t    path  int  false  "时间戳(每次请求获得的t属性)"
// @Success 200  {object} outobjs.OutFriendList
// @router /followers [get]
func (c *FriendShipsController) Followers() {
	current_uid := c.CurrentUid()
	page, _ := c.GetInt("page")
	size, _ := c.GetInt("size")
	uid, _ := c.GetInt64("uid")
	timestamp, _ := c.GetInt64("t")

	t := time.Now()
	if timestamp > 0 {
		t = utils.MsToTime(timestamp) //t = time.Unix(timestamp, 0)
	}

	show_uid := current_uid
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if uid > 0 {
		show_uid = uid
	}
	fship := &passport.FriendShips{}
	total, uids, last_t := fship.FollowerIdsP(show_uid, size, page, t)
	users := []*outobjs.OutMember{}
	for _, _uid := range uids {
		m := outobjs.GetOutMember(_uid, current_uid)
		users = append(users, m)
	}
	out := outobjs.OutFriendList{
		Total:       total,
		TotalPage:   utils.TotalPages(int(total), int(size)),
		CurrentPage: int(page),
		Size:        int(size),
		Users:       users,
		Time:        last_t,
	}
	//自动清空计数器
	if uid <= 0 || current_uid == uid {
		fship.ResetNewFollowers(current_uid)
	}

	c.Json(out)
}

// @Title 两个用户关系的详细情况
// @Description 获取两个用户关系的详细情况
// @Param   access_token   path   string  true  "access_token"
// @Param   target_uid   path   int  true  "目标uid"
// @Success 200  {object} outobjs.OutUserSTRelation
// @router /show [get]
func (c *FriendShipsController) Show() {
	uid := c.CurrentUid()
	target_uid, _ := c.GetInt64("target_uid")
	if target_uid <= 0 {
		c.Json(libs.NewError("member_friend_show_fail", "M4001", "uid不能小于0", ""))
		return
	}
	if uid == target_uid {
		c.Json(libs.NewError("member_friend_show_fail", "M4002", "源和目标不能为同一个人", ""))
		return
	}
	ms := passport.NewMemberProvider()
	target := ms.Get(target_uid)
	if target == nil {
		c.Json(libs.NewError("member_friend_show_target_notexist", "M4003", "目标用户不存在", ""))
		return
	}
	source := ms.Get(uid)
	fship := &passport.FriendShips{}
	stmp := fship.Relation(uid, target_uid)
	out := &outobjs.OutUserSTRelation{
		Source: &outobjs.OutUserRelation{
			Uid:        uid,
			NickName:   source.NickName,
			FollowedBy: stmp.Source.FollowedBy,
			Following:  stmp.Source.Following,
		},
		Target: &outobjs.OutUserRelation{
			Uid:        target_uid,
			NickName:   target.NickName,
			FollowedBy: stmp.Target.FollowedBy,
			Following:  stmp.Target.Following,
		},
	}
	c.Json(out)
}

// @Title 关注某用户
// @Description 关注某用户(参数二选一,成功返回error_code:REP000)
// @Param   access_token   path   string  true  "access_token"
// @Param   uid   path   int  true  "需要关注的uid"
// @Param   nick_name   path  string  true  "需要关注的昵称"
// @Success 200  成功返回error_code:REP000
// @router /create [post]
func (c *FriendShipsController) Create() {
	uid := c.CurrentUid()
	target_uid, _ := c.GetInt64("uid")
	target_name, _ := utils.UrlDecode(c.GetString("nick_name"))
	if target_uid > 0 && len(target_name) > 0 {
		c.Json(libs.NewError("member_friend_create_fail", "M4101", "参数只能二选一", ""))
		return
	}
	if len(target_name) > 0 {
		ms := passport.NewMemberProvider()
		_uid := ms.GetUidByNickname(target_name)
		if _uid <= 0 {
			c.Json(libs.NewError("member_friend_create_fail", "M4102", "昵称所对应的用户不存在", ""))
			return
		}
		target_uid = _uid
	}
	fship := &passport.FriendShips{}
	err := fship.FriendTo(uid, target_uid)
	if err == nil {
		c.Json(libs.NewError("member_friend_create_succ", RESPONSE_SUCCESS, "成功关注", ""))
		return
	}
	c.Json(libs.NewError("member_friend_create_fail", "M4103", err.Error(), ""))
}

// @Title 取消关注某用户
// @Description 取消关注某用户 (参数二选一,成功返回error_code:REP000)
// @Param   access_token   path   string  true  "access_token"
// @Param   uid   path   int  true  "取消关注的uid"
// @Param   nick_name   path  string  true  "取消关注的昵称"
// @Success 200  成功返回error_code:REP000
// @router /destroy [post]
func (c *FriendShipsController) Destroy() {
	uid := c.CurrentUid()
	target_uid, _ := c.GetInt64("uid")
	target_name, _ := utils.UrlDecode(c.GetString("nick_name"))
	if target_uid > 0 && len(target_name) > 0 {
		c.Json(libs.NewError("member_friend_create_fail", "M4201", "参数只能二选一", ""))
		return
	}
	if len(target_name) > 0 {
		ms := passport.NewMemberProvider()
		_uid := ms.GetUidByNickname(target_name)
		if _uid <= 0 {
			c.Json(libs.NewError("member_friend_create_fail", "M4202", "昵称所对应的用户不存在", ""))
			return
		}
		target_uid = _uid
	}
	fship := &passport.FriendShips{}
	err := fship.DestroyFriend(uid, target_uid)
	if err == nil {
		c.Json(libs.NewError("member_friend_create_succ", RESPONSE_SUCCESS, "成功取消关注", ""))
		return
	}
	c.Json(libs.NewError("member_friend_create_fail", "M4203", err.Error(), ""))
}

// @Title 推荐关注人员名单
// @Description 推荐关注人员名单(数组)
// @Param   access_token   path   string  true  "access_token"
// @Param   nums   path   int  true  "数量(默认8)"
// @Success 200
// @router /recmds [get]
func (c *FriendShipsController) Recmds() {
	uid := c.CurrentUid()
	nums, _ := c.GetInt("nums")
	if nums <= 0 {
		nums = 8
	}
	fship := &passport.FriendShips{}
	uids := fship.RecmdFriendUids(uid, int(nums))
	outs := []*outobjs.OutMember{}
	for _, _uid := range uids {
		out_m := outobjs.GetOutMember(_uid, uid)
		if out_m != nil {
			outs = append(outs, out_m)
		}
	}
	c.Json(outs)
}
