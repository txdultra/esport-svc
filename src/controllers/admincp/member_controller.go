package admincp

import (
	"controllers"
	"fmt"
	"libs"
	"libs/passport"
	"libs/vars"
	"outobjs"
	"strconv"
	"strings"
	"time"
	"utils"
)

// 用户管理 API
type MemberCPController struct {
	AdminController
}

func (c *MemberCPController) Prepare() {
	c.AdminController.Prepare()
}

////////////////////////////////////////////////////////////////////////////////
//添加后台管理人员
////////////////////////////////////////////////////////////////////////////////

// @Title 获取所有角色
// @Description 获取所有角色
// @Success 200  {object} outobjs.OutRole
// @router /get_roles [get]
func (c *MemberCPController) GetRoles() {
	service := passport.NewRoleService()
	roles := service.Roles()
	out_roles := make([]*outobjs.OutRole, len(roles), len(roles))
	for i := range roles {
		out_roles[i] = &outobjs.OutRole{
			Id:       roles[i].Id,
			RoleName: roles[i].RoleName,
			Icon:     roles[i].Icon,
			Enabled:  roles[i].Enabled,
		}
	}
	c.Json(out_roles)
}

// @Title 角色人员
// @Description 角色人员
// @Param   role_id   path	int true  "角色id"
// @Success 200  {object} outobjs.OutMemberRole
// @router /get_rolemembers [get]
func (c *MemberCPController) GetRoleMembers() {
	roleid, _ := c.GetInt("role_id")
	if roleid <= 0 {
		c.Json(libs.Error{"member_get_rolemembers_fail", "GM020_161", "角色id非法", ""})
		return
	}
	service := passport.NewRoleService()
	mrs := service.RoleMembers(roleid)
	out_mrs := make([]*outobjs.OutMemberRole, len(mrs), len(mrs))
	for i := range mrs {
		role, err := service.Role(mrs[i].RoleId)
		if err != nil {
			continue
		}
		out_role := &outobjs.OutRole{
			Id:       role.Id,
			RoleName: role.RoleName,
			Icon:     role.Icon,
			Enabled:  role.Enabled,
		}
		out_mrs[i] = &outobjs.OutMemberRole{
			Id:       mrs[i].Id,
			Uid:      mrs[i].Uid,
			Member:   outobjs.GetOutSimpleMember(mrs[i].Uid),
			RoleId:   mrs[i].RoleId,
			Role:     out_role,
			PostTime: mrs[i].PostTime,
			Expries:  mrs[i].Expries,
		}
	}
	c.Json(out_mrs)
}

// @Title 添加编辑
// @Description 添加编辑
// @Param   uid   path	int true  "要添加的用户uid"
// @Param   roles   path	string true  "roles(1|10000,2|20000)"
// @Param   plat   path	string false  "neotv"
// @Param   plat_uid   path	int false  "平台uid"
// @Success 200  {object} libs.Error
// @router /set_roles [post]
func (c *MemberCPController) SetRoles() {
	uid, _ := c.GetInt64("uid")
	roles := strings.Split(c.GetString("roles"), ",")
	plat := c.GetString("plat")
	plat_uid, _ := c.GetInt64("plat_uid")
	mp := passport.NewMemberProvider()
	member := mp.Get(uid)
	if member == nil {
		c.Json(libs.Error{"member_set_roles_fail", "GM020_101", "用户不存在", ""})
		return
	}
	pms := &passport.PlatManagers{}
	if len(plat) > 0 && plat_uid > 0 {
		exist := pms.ExistBind(uid)
		if exist {
			c.Json(libs.Error{"member_set_roles_fail", "GM020_103", "uid已被绑定", ""})
			return
		}
	}

	service := passport.NewRoleService()
	r_ex := make(map[int]int)
	for _, r := range roles {
		rp := strings.Split(r, "|")
		if len(rp) != 2 {
			continue
		}
		_role_id, err := strconv.Atoi(rp[0])
		_exprise, err := strconv.Atoi(rp[1])
		if err == nil {
			r_ex[_role_id] = _exprise
		}
	}
	succ, rlt_err := service.SetMemberRoles(uid, r_ex)
	if succ {

		if len(plat) > 0 && plat_uid > 0 {
			pms.AddPlatManager(uid, plat_uid, plat)
		}
		c.Json(libs.Error{"member_set_roles_succ", controllers.RESPONSE_SUCCESS, "授权成功", ""})
		return
	}
	c.Json(libs.Error{"member_set_roles_fail", "GM020_102", "授权失败:" + rlt_err.Error(), ""})
}

////////////////////////////////////////////////////////////////////////////////
//用户
////////////////////////////////////////////////////////////////////////////////

// @Title 校验播主昵称是否有效
// @Description 校验播主昵称是否有效
// @Param   nick_name   path	string true  "播主昵称"
// @Success 200  {object} libs.Error
// @router /verify_nickname [get]
func (c *MemberCPController) VerifyNickName() {
	nick_name, _ := utils.UrlDecode(c.GetString("nick_name"))
	if len(nick_name) == 0 {
		c.Json(libs.Error{"nick_name_verify_fail", "GM020_011", "参数不能为空", ""})
		return
	}
	mp := passport.NewMemberProvider()
	err := mp.VerifyNickname(0, nick_name)
	if err != nil {
		c.Json(err)
		return
	}
	c.Json(libs.Error{"nick_name_verify_succ", controllers.RESPONSE_SUCCESS, "可以使用", ""})
}

// @Title 校验播主用户名是否有效
// @Description 校验播主用户名是否有效
// @Param   user_name   path	string true  "播主昵称"
// @Success 200  {object} libs.Error
// @router /verify_username [get]
func (c *MemberCPController) VerifyUserName() {
	user_name, _ := utils.UrlDecode(c.GetString("user_name"))
	mp := passport.NewMemberProvider()
	_, err := mp.VerifyUserName(user_name)
	if err != nil {
		c.Json(err)
		return
	}
	c.Json(libs.Error{"user_name_verify_succ", controllers.RESPONSE_SUCCESS, "可以使用", ""})
}

// @Title 添加播主
// @Description 添加播主
// @Param   user_name   path	string true  "播主用户名"
// @Param   nick_name   path	string true  "播主昵称"
// @Param   avatar  path	int true  "播主头像"
// @Param   certified  path	bool false  "是否认证"
// @Param   certified_reson  path	string false  "是否认证"
// @Param   official_certified  path  bool false  "是否官方认证"
// @Param   game_ids   path	string true  "喜好游戏"
// @Success 200  {object} libs.Error
// @router /add [post]
func (c *MemberCPController) AddMember() {
	user_name, _ := utils.UrlDecode(c.GetString("user_name"))
	nick_name, _ := utils.UrlDecode(c.GetString("nick_name"))
	avatar, _ := c.GetInt64("avatar")
	certified, _ := c.GetBool("certified")
	certified_reson, _ := utils.UrlDecode(c.GetString("certified_reson"))
	official_certified, _ := c.GetBool("official_certified")
	ip := c.Ctx.Input.IP()
	game_ids := c.GetString("game_ids")

	mp := passport.NewMemberProvider()
	err := mp.VerifyNickname(0, nick_name)
	if err != nil {
		c.Json(libs.NewError("member_register_fail", "GM020_020", err.Error(), ""))
		return
	}

	_, terr := mp.VerifyUserName(user_name)
	if terr != nil {
		c.Json(terr)
		return
	}

	member := passport.Member{}
	member.UserName = user_name
	member.Password = "def_pwd123"
	member.CreateIP = utils.IpToInt(ip)
	//member.Avatar = avatar
	member.Certified = certified
	member.CertifiedReason = certified_reson
	member.OfficialCertified = official_certified
	member.Src = vars.CLIENT_SRC_MANAGER

	uid, lerr := mp.Create(member, 0, 0)
	if lerr != nil {
		c.Json(lerr)
		return
	}
	if uid <= 0 {
		c.Json(libs.NewError("member_register_fail", "GM020_021", "创建播主失败", ""))
		return
	}
	mp.SetMemberAvatar(uid, avatar)

	err = mp.SetNickname(uid, nick_name)
	if err != nil {
		c.Json(libs.NewError("member_set_nickname_fail", "GM020_022", err.Error(), ""))
		return
	}

	game_ids_arr := strings.Split(game_ids, ",")
	gameIds := []int{}
	for _, gid := range game_ids_arr {
		_id, err := strconv.Atoi(gid)
		if err == nil {
			gameIds = append(gameIds, _id)
		}
	}
	mp.UpdateMemberGames(uid, gameIds)

	c.Json(libs.NewError("member_cpcreate_succ", controllers.RESPONSE_SUCCESS, strconv.FormatInt(uid, 10), ""))
}

// @Title 设置播主认证
// @Description 设置播主认证
// @Param   uid   path	int true  "被认证uid"
// @Param   certified  path	bool true  "是否VIP"
// @Param   certified_reson  path	string false  "是否VIP"
// @Param   official_certified  path  bool false  "是否官方认证"
// @Success 200  {object} libs.Error
// @router /set_certified [post]
func (c *MemberCPController) SetMemberCertifiable() {
	uid, _ := c.GetInt64("uid")
	certified, _ := c.GetBool("certified")
	certified_reson, _ := utils.UrlDecode(c.GetString("certified_reson"))
	officialCertified, _ := c.GetBool("official_certified")
	if uid <= 0 {
		c.Json(libs.NewError("member_set_certified_fail", "GM020_031", "uid参数错误", ""))
		return
	}
	mp := passport.NewMemberProvider()
	err := mp.SetMemberCertified(uid, certified, certified_reson, officialCertified)
	if err != nil {
		c.Json(libs.NewError("member_set_certified_fail", "GM020_032", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("member_set_certified_succ", controllers.RESPONSE_SUCCESS, "成功设置", ""))
}

// @Title 校验播主用户名是否有效
// @Description 校验播主用户名是否有效
// @Param   uid   path	int true  "uid"
// @Param   avatar   path	int true  "播主头像"
// @Param   game_ids   path	string true  "喜好游戏"
// @Success 200  {object} libs.Error
// @router /update [post]
func (c *MemberCPController) UpdateMember() {
	uid, _ := c.GetInt64("uid")
	avatar, _ := c.GetInt64("avatar")
	game_ids := c.GetString("game_ids")
	if uid <= 0 {
		c.Json(libs.NewError("member_update_fail", "GM020_041", "uid参数错误", ""))
		return
	}
	if avatar < 0 {
		avatar = 0
	}
	mp := passport.NewMemberProvider()
	member := mp.Get(uid)
	if member == nil {
		c.Json(libs.NewError("member_update_fail", "GM020_042", "用户不存在", ""))
		return
	}

	if avatar > 0 {
		if member.Avatar != avatar {
			_, err := mp.SetMemberAvatar(uid, avatar)
			if err != nil {
				c.Json(libs.NewError("member_update_fail", "GM020_043", err.Error(), ""))
				return
			}
		}
	} else {
		member.Avatar = 0
		err := mp.Update(*member)
		if err != nil {
			c.Json(libs.NewError("member_update_fail", "GM020_043", err.Error(), ""))
			return
		}
	}

	game_ids_arr := strings.Split(game_ids, ",")
	gameIds := []int{}
	for _, gid := range game_ids_arr {
		_id, err := strconv.Atoi(gid)
		if err == nil {
			gameIds = append(gameIds, _id)
		}
	}
	err := mp.UpdateMemberGames(uid, gameIds)
	if err != nil {
		c.Json(libs.NewError("member_update_fail", "GM020_043", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("member_update_succ", controllers.RESPONSE_SUCCESS, "操作成功", ""))
}

// @Title 重设昵称
// @Description 重设昵称
// @Param   uid   path	int true  "uid"
// @Param   nick_name   path	string true  "昵称名称"
// @Success 200  {object} libs.Error
// @router /reset_nickname [post]
func (c *MemberCPController) ResetNickName() {
	uid, _ := c.GetInt64("uid")
	nick_name, _ := utils.UrlDecode(c.GetString("nick_name"))
	mp := passport.NewMemberProvider()
	err := mp.SetNickname(uid, nick_name)
	if err != nil {
		c.Json(libs.NewError("member_reset_nickname_fail", "GM020_051", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("member_reset_nickname_succ", controllers.RESPONSE_SUCCESS, "操作成功", ""))
}

// @Title 获取用户游戏喜好
// @Description 获取用户游戏喜好(返回数组)
// @Param   uid   path   int  true  "用户id"
// @Success 200 {object} outobjs.OutMemberGame
// @router /games [get]
func (c *MemberCPController) GetMemberGames() {
	uid, _ := c.GetInt64("uid")
	out_mgs := []*outobjs.OutMemberGame{}
	if uid <= 0 {
		c.Json(out_mgs)
		return
	}
	provider := passport.NewMemberProvider()
	mgs := provider.MemberGames(uid)
	for _, mg := range mgs {
		out_game := outobjs.GetOutGameById(mg.GameId)
		if out_game == nil {
			continue
		}
		out_mgs = append(out_mgs, &outobjs.OutMemberGame{
			AddTime: mg.PostTime,
			Game:    out_game,
		})
	}
	c.Json(out_mgs)
}

// @Title 查询播主
// @Description 查询播主
// @Param   query   path	string true  "关键字"
// @Param   certified  path	bool false  "vip用户"
// @Param   official_certified  path	bool false  " 官方认证"
// @Param   page  path	int true  "page"
// @Param   size  path	int false  "页数量(默认20)"
// @Success 200  {object} outobjs.OutMemberPageList
// @router /search [get]
func (c *MemberCPController) Search() {
	query, _ := utils.UrlDecode(c.GetString("query"))
	page, _ := c.GetInt("page")
	size, _ := c.GetInt("size")
	certified, _ := c.GetBool("certified")
	official_certified, _ := c.GetBool("official_certified")
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	//match_mode := c.GetString("match_mode")
	//timestamp, _ := c.GetInt("t")
	//if len(match_mode) == 0 {
	//	match_mode = "any"
	//}
	//if len(query) == 0 {
	//	match_mode = "all"
	//}
	//t := time.Now()
	//if timestamp > 0 {
	//	t = time.Unix(timestamp, 0)
	//}

	mp := passport.NewMemberProvider()
	total, uids := mp.QueryForAdmin(query, certified, official_certified, int(page), int(size))
	out_members := []*outobjs.OutMember{}
	for _, uid := range uids {
		m := outobjs.GetOutMember(uid, 0)
		if m != nil {
			out_members = append(out_members, m)
		}
	}

	out := outobjs.OutMemberPageList{
		Total:       total,
		TotalPage:   utils.TotalPages(int(total), int(size)),
		CurrentPage: int(page),
		Size:        int(size),
		Lists:       out_members,
		Time:        time.Now().Unix(),
	}
	c.Json(out)
}

// @Title 操作积分
// @Description 操作积分
// @Param   uid path int true  "uid"
// @Param   type path string true  "赠送类型(p=积分,j=竞币)"
// @Param   nums  path	int true  "数量"
// @Param   desc  path	string false  "描述"
// @Param   action_pwd  path	string true  "操作密码"
// @Success 200  {object} libs.Error
// @router /action_credit [post]
func (c *MemberCPController) ActionCredit() {
	current_uid := c.CurrentUid()
	uid, _ := c.GetInt64("uid")
	t := c.GetString("type")
	nums, _ := c.GetInt64("nums")
	desc, _ := utils.UrlDecode(c.GetString("desc"))
	pwd := c.GetString("action_pwd")
	if pwd != fmt.Sprintf("djq_%d", current_uid) {
		c.Json(libs.NewError("member_action_credit_fail", "GM020_074", "操作密码错误", ""))
		return
	}
	if uid <= 0 {
		c.Json(libs.NewError("member_action_credit_fail", "GM020_070", "参数错误", ""))
		return
	}
	if len(t) == 0 || (t != "p" && t != "j") {
		c.Json(libs.NewError("member_action_credit_fail", "GM020_071", "参数错误", ""))
		return
	}
	if nums == 0 {
		c.Json(libs.NewError("member_action_credit_fail", "GM020_072", "参数错误", ""))
		return
	}
	ptype := vars.CURRENCY_TYPE_CREDIT
	if t == "j" {
		ptype = vars.CURRENCY_TYPE_JING
	}
	mp := passport.NewMemberProvider()
	no, err := mp.ActionCredit(current_uid, uid, ptype, nums, desc)
	if err != nil {
		c.Json(libs.NewError("member_action_credit_fail", "GM020_075", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("member_action_credit_success", controllers.RESPONSE_SUCCESS, "操作成功:"+no, ""))
}
