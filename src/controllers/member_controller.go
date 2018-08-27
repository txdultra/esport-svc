package controllers

import (
	"fmt"
	//"fmt"
	"encoding/json"
	"libs"
	"libs/passport"
	"libs/stat"
	"libs/vars"
	"outobjs"
	"strconv"
	"strings"
	"time"
	"utils"
)

// 用户 API
type MemberController struct {
	BaseController
	m_provider *passport.MemberProvider
}

func (c *MemberController) Prepare() {
	c.BaseController.Prepare()
	c.m_provider = passport.NewMemberProvider()
}

func (c *MemberController) URLMapping() {
	c.Mapping("Register", c.Register)
	c.Mapping("Get", c.Get)
	c.Mapping("Update", c.Update)
	c.Mapping("ShowFoundPwdMobile", c.ShowFoundPwdMobile)
	c.Mapping("SetFoundPwdMobile", c.SetFoundPwdMobile)
	c.Mapping("Login", c.Login)
	c.Mapping("Logout", c.Logout)
	c.Mapping("Search", c.Search)
	c.Mapping("MemberGames", c.MemberGames)
	c.Mapping("MemberGameSingle", c.MemberGameSingle)
	c.Mapping("RemoveMemberGameSingle", c.RemoveMemberGameSingle)
	c.Mapping("GetMemberGames", c.GetMemberGames)
	c.Mapping("SetPushId", c.SetPushId)
	c.Mapping("IsSetPushId", c.IsSetPushId)
	c.Mapping("SetNickname", c.SetNickname)
	c.Mapping("SetAvatar", c.SetAvatar)
	c.Mapping("FrequentAts", c.FrequentAts)
	c.Mapping("GetMyConfig", c.GetMyConfig)
	c.Mapping("SetMyConfig", c.SetMyConfig)
	c.Mapping("SetMemberBackgroundImg", c.SetMemberBackgroundImg)
	c.Mapping("ResetPwd", c.ResetPwd)
	c.Mapping("GetByUserName", c.GetByUserName)
}

// @Title 注册用户
// @Description 注册用户(成功注册后请调用set_pushid接口绑定第三方推送账号)
// @Param   user_name   path  string  true  "email或手机号形式"
// @Param   pwd    path  string  true  "密码"
// @Param   mobile_iden   path  string  false  "用户手机唯一序列号"
// @Param   longi   path  float  false  "经度"
// @Param   lati   path  float  false  "纬度"
// @Param   vfy_code  path  string  false  "验证码随机字符串"
// @Param   vfy_str   path  string  false "验证码"
// @Param   src   path  int  false "来源(0=app,1=web 默认0)"
// @Success 200 {object} outobjs.OutAccessToken
// @router /register [post]
func (c *MemberController) Register() {
	user_name := c.GetString("user_name")
	pwd := c.GetString("pwd")
	mobile_iden := c.GetString("mobile_iden")
	longiTude, _ := c.GetFloat("longi")
	latiTude, _ := c.GetFloat("lati")
	vfy_code := c.GetString("vfy_code")
	vfy_str := c.GetString("vfy_str")
	ip := c.Ctx.Input.IP()
	src, _ := c.GetInt8("src")
	if len(vfy_code) > 0 {
		lerr := c.checkVerifyCode(vfy_code, vfy_str)
		if lerr != nil {
			c.Json(lerr)
			return
		}
	}
	longi := float32(longiTude)
	lati := float32(latiTude)

	if len(pwd) < passport.MemberPasswordMinLen {
		c.Json(libs.NewError("member_register_fail", "M3902", fmt.Sprintf("密码必须大于等于%d位", passport.MemberPasswordMinLen), ""))
		return
	}
	if len(pwd) > passport.MemberPasswordMaxLen {
		c.Json(libs.NewError("member_register_fail", "M3902", fmt.Sprintf("密码必须小于等于%d位", passport.MemberPasswordMaxLen), ""))
		return
	}
	if src > 1 || src < 0 {
		c.Json(libs.NewError("member_register_fail", "M3903", "来源编号错误", ""))
		return
	}

	member := passport.Member{}
	member.UserName = user_name
	member.Password = pwd
	member.MobileIdentifier = mobile_iden
	member.CreateIP = utils.IpToInt(ip)
	member.Src = vars.CLIENT_SRC(src)
	member_provider := passport.NewMemberProvider()
	uid, err := member_provider.Create(member, longi, lati)
	if err != nil {
		c.Json(err)
		return
	}
	if uid <= 0 {
		c.Json(libs.NewError("member_register_fail", "M3901", "创建用户失败", ""))
		return
	}
	access_service := passport.NewAccessTokenService()
	acc_token, terr := access_service.GetAccessTokenNoCode(ip, "", "", uid)
	if terr == nil {
		c.Json(outobjs.GetAccessToken(acc_token))
		return
	}
	c.Json(terr)
}

// @Title 获取用户信息
// @Description 获取用户信息 (uid,昵称二选一)
// @Param   access_token   path   string  true  "access_token"
// @Param   uid   path  int  false  "用户id"
// @Param   nick_name   path  string  false  "用户昵称"
// @Success 200
// @router /show [get]
func (c *MemberController) Get() {
	currentUid := c.CurrentUid()
	if currentUid <= 0 {
		c.Json(libs.NewError("member_show_premission_denied", UNAUTHORIZED_CODE, "你没有操作的权限", ""))
		return
	}
	uid, _ := c.GetInt64("uid")
	nick_name, _ := utils.UrlDecode(c.GetString("nick_name"))
	if uid <= 0 && len(nick_name) == 0 {
		c.Json(libs.NewError("member_show_parameters", "M3001", "uid或nick_name参数必须一个有值", ""))
		return
	}
	var member *passport.Member
	if uid > 0 {
		member = c.m_provider.Get(uid)
	} else {
		_uid := c.m_provider.GetUidByNickname(nick_name)
		if _uid > 0 {
			uid = _uid
			member = c.m_provider.Get(_uid)
		}
	}
	if member == nil {
		c.Json(libs.NewError("member_show_parameters", "M3002", "指定的uid或nick_name用户不存在", ""))
		return
	}
	ms := c.m_provider.GetState(uid)
	mp := c.m_provider.GetProfile(uid)
	if ms == nil || mp == nil {
		c.Json(libs.NewError("member_show_fail", "M3099", "系统数据发生错误,请联系管理员", ""))
		return
	}
	out := &outobjs.OutMemberInfo{
		OutMember:        outobjs.GetOutMember(uid, currentUid),
		Logins:           ms.Logins,
		LastLoginIP:      utils.IntToIp(ms.LastLoginIP),
		LastLoginTime:    time.Unix(ms.LastLoginTime, 0),
		LongiTude:        ms.LongiTude,
		LatiTude:         ms.LatiTude,
		Vods:             ms.Vods,
		Fans:             ms.Fans,
		Notes:            ms.Notes,
		Count1:           ms.Count1,
		Count2:           ms.Count2,
		Count3:           ms.Count3,
		Count4:           ms.Count4,
		Count5:           ms.Count5,
		BackgroundImg:    mp.BackgroundImg,
		BackgroundImgUrl: file_storage.GetFileUrl(mp.BackgroundImg),
		Description:      mp.Description,
		Gender:           mp.Gender,
		RealName:         mp.RealName,
		BirthYear:        mp.BirthYear,
		BirthMonth:       mp.BirthMonth,
		BirthDay:         mp.BirthDay,
		Mobile:           utils.StringReplace(mp.Mobile, 5, 9, "*"),
		IDCard:           mp.IDCard,
		QQ:               mp.QQ,
		Alipay:           mp.Alipay,
		Youku:            mp.Youku,
		Bio:              mp.Bio,
		Field1:           mp.Field1,
		Field2:           mp.Field2,
		Field3:           mp.Field3,
		Field4:           mp.Field4,
		Field5:           mp.Field5,
	}
	c.Json(out)
}

// @Title 更新用户信息
// @Description 更新用户信息(error_code=REP000表示更新成功)
// @Param   access_token   path   string  true  "access_token"
// @Param   description   path  string  false  "用户描述"
// @Param   gender   path  string  false  "用户性别(m为男,n女 默认m)"
// @Param   real_name   path  string  false  "用户姓名"
// @Param   birth_year   path  int  false  "出生年"
// @Param   birth_month   path  int  false  "出生月"
// @Param   birth_day   path  int  false  "出生天"
// @Param   mobile   path  string  false  "手机号"
// @Param   id_card   path  string  false  "身份证号"
// @Param   qq   path  string  false  "qq"
// @Param   alipay   path  string  false  "alipy"
// @Param   youku   path  string  false  "youku"
// @Param   bio   path  string  false  "个人简历"
// @Param   field1   path  string  false  "其他1"
// @Param   field2   path  string  false  "其他2"
// @Param   field3   path  string  false  "其他3"
// @Param   field4   path  string  false  "其他4"
// @Param   field5   path  string  false  "其他5"
// @Success 200
// @router /profile/update [post]
func (c *MemberController) Update() {
	current_uid := c.CurrentUid()

	if current_uid <= 0 {
		c.Json(libs.NewError("member_update_profile_premission_denied", UNAUTHORIZED_CODE, "你没有操作的权限", ""))
		return
	}
	desc, _ := utils.UrlDecode(c.GetString("description"))
	gender := c.GetString("gender")
	real_name, _ := utils.UrlDecode(c.GetString("real_name"))
	birth_year, _ := c.GetInt("birth_year")
	birth_month, _ := c.GetInt("birth_month")
	birth_day, _ := c.GetInt("birth_day")
	mobile := c.GetString("mobile")
	id_card := c.GetString("id_card")
	qq := c.GetString("qq")
	alipay := c.GetString("alipay")
	youku := c.GetString("youku")
	bio := c.GetString("bio")
	field1, _ := utils.UrlDecode(c.GetString("field1"))
	field2, _ := utils.UrlDecode(c.GetString("field2"))
	field3, _ := utils.UrlDecode(c.GetString("field3"))
	field4, _ := utils.UrlDecode(c.GetString("field4"))
	field5, _ := utils.UrlDecode(c.GetString("field5"))
	_birth_year := int(birth_year)
	_birth_month := int(birth_month)
	_birth_day := int(birth_day)

	if gender != "m" && gender != "n" {
		c.Json(libs.NewError("member_update_profile_gender", "M3101", "性别格式错误", ""))
		return
	}
	if _birth_year < 0 || _birth_year > time.Now().Year() {
		c.Json(libs.NewError("member_update_profile_birthyear", "M3102", "年龄年格式错误", ""))
		return
	}
	if _birth_month < 0 || _birth_month > 12 {
		c.Json(libs.NewError("member_update_profile_birthmonth", "M3103", "年龄月格式错误", ""))
		return
	}
	if _birth_day < 0 || _birth_day > 31 {
		c.Json(libs.NewError("member_update_profile_birthday", "M3104", "年龄日格式错误", ""))
		return
	}

	mp := c.m_provider.GetProfile(current_uid)
	if mp == nil {
		c.Json(libs.NewError("member_update_profile_member_notexist", "M3105", "用户不存在", ""))
		return
	}
	mp.Description = desc
	mp.Gender = gender
	mp.RealName = real_name
	mp.BirthYear = _birth_year
	mp.BirthMonth = _birth_month
	mp.BirthDay = _birth_day
	mp.Mobile = mobile
	mp.IDCard = id_card
	mp.QQ = qq
	mp.Alipay = alipay
	mp.Youku = youku
	mp.Bio = bio
	mp.Field1 = field1
	mp.Field2 = field2
	mp.Field3 = field3
	mp.Field4 = field4
	mp.Field5 = field5
	err := c.m_provider.UpdateProfile(*mp)
	if err == nil {
		c.Json(libs.NewError("member_update_profile_success", RESPONSE_SUCCESS, "更新成功", ""))
		return
	}
	c.Json(libs.NewError("member_update_profile_fail", "M3108", "更新错误", ""))
}

// @Title 更新用户空间背景图
// @Description 更新用户空间背景图
// @Param   access_token   path   string  true  "access_token"
// @Param   file_id   path  int  true  "上传文件后的fileid"
// @Success 200 {object} outobjs.OutFileUrl
// @router /profile/set_bgimg [post]
func (c *MemberController) SetMemberBackgroundImg() {
	file_id, _ := c.GetInt64("file_id")
	if file_id <= 0 {
		c.Json(libs.NewError("member_set_backgroundimg_fail", "M3230", "文件id错误", ""))
		return
	}
	current_uid := c.CurrentUid()
	if current_uid <= 0 {
		c.Json(libs.NewError("member_set_backgroundimg_premission_denied", UNAUTHORIZED_CODE, "你没有操作的权限", ""))
		return
	}
	mp := passport.NewMemberProvider()
	fn, err := mp.SetMemberBackgroundImg(current_uid, file_id)
	if err != nil {
		c.Json(libs.NewError("member_set_backgroundimg_fail", "M3231", err.Error(), ""))
		return
	}
	ofu := &outobjs.OutFileUrl{
		FileId:  fn.FileId,
		FileUrl: file_storage.GetFileUrl(fn.FileId),
	}
	c.Json(ofu)
}

// @Title 找回密码设定的手机号
// @Description 更新用户信息(error_code=REP000表示获取成功,error_content中带有手机号,可能为空)
// @Param   access_token   path   string  true  "access_token"
// @Success 200
// @router /security/found_pwd_mobile [get]
func (c *MemberController) ShowFoundPwdMobile() {
	current_uid := c.CurrentUid()

	if current_uid <= 0 {
		c.Json(libs.NewError("member_update_foundpwdmobile_premission_denied", UNAUTHORIZED_CODE, "你没有操作的权限", ""))
		return
	}
	mp := c.m_provider.GetProfile(current_uid)
	if mp == nil {
		c.Json(libs.NewError("member_update_foundpwdmobile_member_notexist", "M3203", "用户不存在", ""))
		return
	}
	if len(mp.FoundPwdMobile) == 0 {
		c.Json(libs.NewError("member_update_foundpwdmobile", RESPONSE_SUCCESS, "", ""))
		return
	}
	c.Json(libs.NewError("member_update_foundpwdmobile", RESPONSE_SUCCESS, utils.StringReplace(mp.FoundPwdMobile, 5, 9, "*"), ""))
}

// @Title 更新找回密码使用的手机号
// @Description 更新用户信息(error_code=REP000表示更新成功)
// @Param   access_token   path   string  true  "access_token"
// @Param   mobile   path  string  false  "找回密码所使用的手机号"
// @Success 200 返回对象的error_code=REP000表示更新成功
// @router /security/found_pwd_mobile [post]
func (c *MemberController) SetFoundPwdMobile() {
	current_uid := c.CurrentUid()
	if current_uid <= 0 {
		c.Json(libs.NewError("member_update_foundpwdmobile_premission_denied", UNAUTHORIZED_CODE, "你没有操作的权限", ""))
		return
	}
	found_mobile := c.GetString("mobile")
	if !utils.IsMobile(found_mobile) {
		c.Json(libs.NewError("member_update_foundpwdmobile_format_error", "M3302", "手机号格式错误", ""))
		return
	}
	mp := c.m_provider.GetProfile(current_uid)
	if mp == nil {
		c.Json(libs.NewError("member_update_foundpwdmobile_member_notexist", "M3303", "用户不存在", ""))
		return
	}
	mp.FoundPwdMobile = found_mobile
	err := c.m_provider.UpdateProfile(*mp)
	if err == nil {
		c.Json(libs.NewError("member_update_foundpwdmobile_success", RESPONSE_SUCCESS, "更新成功", ""))
		return
	}
	c.Json(libs.NewError("member_update_foundpwdmobile_fail", "M3308", "更新错误", ""))
}

// @Title 用户名密码登录
// @Description 用户名密码登录(成功登录后请调用set_pushid接口绑定第三方推送账号)
// @Param   user_name   path   string  true  "用户名"
// @Param   pwd   path  string  true  "密码"
// @Param   vfy_code  path  string  false  "验证码随机字符串"
// @Param   vfy_str   path  string  false "验证码"
// @Success 200 {object} outobjs.OutAccessToken
// @router /login [post]
func (c *MemberController) Login() {
	user_name := c.GetString("user_name")
	pwd := c.GetString("pwd")
	vfy_code := c.GetString("vfy_code")
	vfy_str := c.GetString("vfy_str")
	ip := c.Ctx.Input.IP()
	if len(vfy_code) > 0 {
		lerr := c.checkVerifyCode(vfy_code, vfy_str)
		if lerr != nil {
			c.Json(lerr)
			return
		}
	}
	login_status, member := c.m_provider.CheckLoginPassword(user_name, strings.Trim(pwd, " "))
	if login_status != passport.LOGIN_ACTION_STATUS_CHECKPWD_SUCC {
		c.Json(libs.NewError("member_login_error", "M3404", string(login_status), ""))
		return
	}
	uid := member.Uid
	access_service := passport.NewAccessTokenService()
	acc_token, terr := access_service.GetAccessTokenNoCode(ip, "", "", uid)
	if terr == nil {
		go stat.GetCounter(passport.MOD_NAME).DoC(uid, 1, "logins")
		data := outobjs.GetAccessToken(acc_token)
		c.Json(data)
		return
	}
	c.Json(terr)
}

func (c *MemberController) checkVerifyCode(code string, source string) *libs.Error {
	cache := utils.GetCache()
	str := ""
	err := cache.Get(code, &str)
	if err != nil {
		return libs.NewError("member_login_verifycode_invalid", "M3402", "验证码随机字符串无效", "")
	}
	cache.Delete(code)
	if str != source {
		return libs.NewError("member_login_verifycode_fail", "M3403", "验证码错误", "")
	}
	return nil
}

// @Title 用户注销登录
// @Description 用户注销登录
// @Param   access_token   path   string  true  "access_token"
// @Success 200 返回对象的error_code=REP000表示更新成功
// @router /logout [post]
func (c *MemberController) Logout() {
	uid := c.CurrentUid()
	if uid <= 0 {
		c.Json(libs.NewError("member_logout_premission_denied", UNAUTHORIZED_CODE, "access_token过期", ""))
		return
	}
	//清空绑定的推送id
	//member_provider := passport.NewMemberProvider()
	//member_provider.EmptyPushConfig(uid)
	access_service := passport.NewAccessTokenService()
	access_token, _ := access_service.GetTokenObj(c.GetString("access_token"))
	if access_token != nil {
		access_service.RevokeAccessToken(access_token)
	}
	c.Json(libs.NewError("member_logout_success", RESPONSE_SUCCESS, "成功注销", ""))
}

// @Title 查询用户access_token的授权相关信息
// @Description 查询用户access_token的授权相关信息
// @Param   access_token   path   string  true  "access_token"
// @Success 200 {object} outobjs.OutAccessToken
// @router /get_token [post]
func (c *MemberController) GetAccessToken() {
	uid := c.CurrentUid()
	if uid <= 0 {
		c.Json(libs.NewError("member_gettoken_premission_denied", UNAUTHORIZED_CODE, "access_token过期", ""))
		return
	}
	token := c.GetString("access_token")
	access_service := passport.NewAccessTokenService()
	accessToken, err := access_service.GetTokenObj(token)
	if err == nil {
		c.Json(outobjs.GetAccessToken(accessToken))
		return
	}
	c.Json(err)
}

// @Title 设置用户登录密码
// @Description 设置用户登录密码
// @Param   access_token   path   string  true  "access_token"
// @Param   pwd   path  string  true  "密码(必须大于等于6位的字母和数字组成)"
// @Param   vfy_code  path  string  false  "验证码随机字符串"
// @Param   vfy_str   path  string  false "验证码"
// @Success 200 返回对象的error_code=REP000表示更新成功
// @router /set_password [post]
func (c *MemberController) SetPassword() {
	uid := c.CurrentUid()
	if uid <= 0 {
		c.Json(libs.NewError("member_setpwd_premission_denied", UNAUTHORIZED_CODE, "未登录状态不能重设密码", ""))
		return
	}
	pwd := c.GetString("pwd")
	vfy_code := c.GetString("vfy_code")
	vfy_str := c.GetString("vfy_str")
	//ip := c.Ctx.Input.IP()
	if len(vfy_code) > 0 {
		lerr := c.checkVerifyCode(vfy_code, vfy_str)
		if lerr != nil {
			c.Json(lerr)
			return
		}
	}
	err := c.m_provider.ResetPassword(uid, pwd)
	if err != nil {
		c.Json(libs.NewError("member_setpwd_fail", "M3502", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("member_setpwd_succ", RESPONSE_SUCCESS, "重设密码成功", ""))
}

// @Title 设置用户邮箱
// @Description 设置用户邮箱
// @Param   access_token   path   string  true  "access_token"
// @Param   email   path  string  true  "邮箱"
// @Param   vfy_code  path  string  false  "验证码随机字符串"
// @Param   vfy_str   path  string  false "验证码"
// @Success 200 返回对象的error_code=REP000表示更新成功
// @router /set_email [post]
func (c *MemberController) SetMail() {
	uid := c.CurrentUid()
	if uid <= 0 {
		c.Json(libs.NewError("member_setmail_premission_denied", UNAUTHORIZED_CODE, "未登录状态不能重设邮箱", ""))
		return
	}
	mail := c.GetString("email")
	vfy_code := c.GetString("vfy_code")
	vfy_str := c.GetString("vfy_str")
	//ip := c.Ctx.Input.IP()
	if len(vfy_code) > 0 {
		lerr := c.checkVerifyCode(vfy_code, vfy_str)
		if lerr != nil {
			c.Json(lerr)
			return
		}
	}
	err := c.m_provider.ResetMail(uid, mail)
	if err != nil {
		c.Json(libs.NewError("member_setmail_fail", "M3602", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("member_setmail_succ", RESPONSE_SUCCESS, "重设邮箱成功", ""))
}

// @Title 设置用户游戏喜好(单个提交)
// @Description 设置用户游戏喜好(单个提交)
// @Param   access_token   path   string  true  "access_token"
// @Param   game_id  path  string  true  "游戏编号"
// @Success 200 返回对象的error_code=REP000表示设置成功
// @router /games_single [post]
func (c *MemberController) MemberGameSingle() {
	current_uid := c.CurrentUid()
	game_id, _ := c.GetInt("game_id")
	if current_uid <= 0 {
		c.Json(libs.NewError("member_update_membergames_premission_denied", UNAUTHORIZED_CODE, "你没有操作的权限", ""))
		return
	}
	if game_id <= 0 {
		c.Json(libs.NewError("member_update_membergames_fail", "M3705", "参数错误", ""))
		return
	}
	ms := passport.NewMemberProvider()
	err := ms.UpdateMemberGameSingle(current_uid, game_id)
	if err != nil {
		c.Json(libs.NewError("member_update_membergames_fail", "M3706", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("member_update_membergames_succ", RESPONSE_SUCCESS, "更新成功", ""))
}

// @Title 删除用户游戏喜好(单个提交)
// @Description 删除用户游戏喜好(单个提交)
// @Param   access_token   path   string  true  "access_token"
// @Param   game_id  path  string  true  "游戏编号"
// @Success 200 返回对象的error_code=REP000表示设置成功
// @router /games_single [delete]
func (c *MemberController) RemoveMemberGameSingle() {
	current_uid := c.CurrentUid()
	game_id, _ := c.GetInt("game_id")
	if current_uid <= 0 {
		c.Json(libs.NewError("member_remove_membergames_premission_denied", UNAUTHORIZED_CODE, "你没有操作的权限", ""))
		return
	}
	if game_id <= 0 {
		c.Json(libs.NewError("member_remove_membergames_fail", "M3707", "参数错误", ""))
		return
	}
	ms := passport.NewMemberProvider()
	err := ms.RemoveMemberGameSingle(current_uid, game_id)
	if err != nil {
		c.Json(libs.NewError("member_remove_membergames_fail", "M3708", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("member_remove_membergames_succ", RESPONSE_SUCCESS, "更新成功", ""))
}

// @Title 设置用户游戏喜好
// @Description 设置用户游戏喜好
// @Param   access_token   path   string  true  "access_token"
// @Param   game_ids   path  string  true  "游戏编号以,逗号分隔"
// @Success 200 返回对象的error_code=REP000表示设置成功
// @router /games [post]
func (c *MemberController) MemberGames() {
	current_uid := c.CurrentUid()
	if current_uid <= 0 {
		c.Json(libs.NewError("member_update_membergames_premission_denied", UNAUTHORIZED_CODE, "你没有操作的权限", ""))
		return
	}
	game_ids := c.GetString("game_ids")
	game_ids_arr := strings.Split(game_ids, ",")
	gameIds := []int{}
	for _, gid := range game_ids_arr {
		_id, err := strconv.Atoi(gid)
		if err == nil {
			gameIds = append(gameIds, _id)
		}
	}
	ms := passport.NewMemberProvider()
	err := ms.UpdateMemberGames(current_uid, gameIds)
	if err != nil {
		c.Json(libs.NewError("member_update_membergames_fail", "M3702", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("member_update_membergames_succ", RESPONSE_SUCCESS, "更新成功", ""))
}

// @Title 获取用户游戏喜好
// @Description 获取用户游戏喜好(返回数组)
// @Param   uid   path   int  true  "用户id"
// @Success 200 {object} outobjs.OutMemberGame
// @router /games [get]
func (c *MemberController) GetMemberGames() {
	uid, _ := c.GetInt64("uid")
	out_mgs := []*outobjs.OutMemberGame{}
	if uid <= 0 {
		c.Json(out_mgs)
		return
	}
	bas := &libs.Bas{}
	games := bas.Games()

	provider := passport.NewMemberProvider()
	mgs := provider.MemberGames(uid)

	for _, game := range games {
		var memberGame *passport.MemberGame = nil
		for _, mg := range mgs {
			if mg.GameId == game.Id {
				memberGame = mg
			}
		}
		if game.ForcedzSel {
			pTime := game.PostTime
			if memberGame != nil {
				pTime = memberGame.PostTime
			}
			out_game := outobjs.GetOutGame(&game)
			out_mgs = append(out_mgs, &outobjs.OutMemberGame{
				AddTime: pTime,
				Game:    out_game,
			})
			continue
		}
		if memberGame == nil {
			continue
		}
		out_game := outobjs.GetOutGame(&game)
		if out_game == nil || !out_game.Enabled {
			continue
		}
		out_mgs = append(out_mgs, &outobjs.OutMemberGame{
			AddTime: memberGame.PostTime,
			Game:    out_game,
		})
	}

	//	for _, mg := range mgs {
	//		out_game := outobjs.GetOutGameById(mg.GameId)
	//		if out_game == nil || !out_game.Enabled {
	//			continue
	//		}
	//		out_mgs = append(out_mgs, &outobjs.OutMemberGame{
	//			AddTime: mg.PostTime,
	//			Game:    out_game,
	//		})
	//	}
	c.Json(out_mgs)
}

// @Title 获取用户列表
// @Description 用户列表(返回数组)
// @Param   access_token   path   string  false  "access_token"
// @Param   query    path    string  false  "查询字符串"
// @Param   size     path    int  false  "一页行数(默认20)"
// @Param   page     path    int  false  "页数(默认1)"
// @Param   match_mode path  string   false  "搜索模式(all,any,phrase,boolean,extended,fullscan,extended2),默认为any"
// @Success 200 {object} outobjs.OutMemberPageList
// @router /search [get]
func (c *MemberController) Search() {
	current_uid := c.CurrentUid()
	querystr, _ := utils.UrlDecode(c.GetString("query"))
	size, _ := c.GetInt("size")
	page, _ := c.GetInt("page")
	match_mode := c.GetString("match_mode")
	timestamp, _ := c.GetInt64("t")
	if len(match_mode) == 0 {
		match_mode = "any"
	}
	t := time.Now()
	if timestamp > 0 {
		t = time.Unix(timestamp, 0)
	}

	//1分钟间隔缓存,有用户关注状态(无法缓存)
	//query_cache_key := fmt.Sprintf("front_fast_cache.members.query:words_%s_p_%d_s_%d_mode_%s_t_%s",
	//	querystr, page, size, match_mode)
	//c_obj := utils.GetLocalFastTimePartCache(t, query_cache_key, utils.CACHE_INTERVAL_TIME_TYPE_MINUTE)
	//if c_obj != nil {
	//	c.Json(c_obj)
	//	return
	//}

	if size <= 0 {
		size = 20
	}
	if page <= 0 {
		page = 1
	}
	mp := passport.NewMemberProvider()
	total, uids := mp.Query(querystr, int(page), int(size), match_mode, "", nil, nil)
	out_members := []*outobjs.OutMember{}
	for _, uid := range uids {
		m := outobjs.GetOutMember(uid, current_uid)
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
		Time:        t.Unix(),
	}
	//utils.SetLocalFastTimePartCache(t, query_cache_key, utils.CACHE_INTERVAL_TIME_TYPE_MINUTE, out)
	c.Json(out)
}

// @Title 提交在第三方推送绑定时的用户id
// @Description 提交在第三方推送绑定时的用户id
// @Param   access_token   path   string  true  "access_token"
// @Param   proxy_id    path    int  true  "第三方标识(1表示百度云推)"
// @Param   channel_id    path   string  false  " 推送通道ID"
// @Param   user_id    path    string  true  "第三方用户id"
// @Param   device_type  path   string  true  "客户端系统标识(andriod,ios,wp)"
// @Success 200
// @router /set_pushid [post]
func (c *MemberController) SetPushId() {
	current_uid := c.CurrentUid()
	if current_uid <= 0 {
		c.Json(libs.NewError("member_set_pushid_premission_denied", UNAUTHORIZED_CODE, "你没有操作的权限", ""))
		return
	}
	proxy_id, _ := c.GetInt("proxy_id")
	channel_id := c.GetString("channel_id")
	user_id := c.GetString("user_id")
	device_type := c.GetString("device_type")

	if proxy_id != 1 {
		c.Json(libs.NewError("member_set_pushid_parameter_error", "M3801", "proxy_id必须设置为1", ""))
		return
	}
	if len(user_id) == 0 {
		c.Json(libs.NewError("member_set_pushid_parameter_error", "M3802", "user_id不能为空", ""))
		return
	}
	device_type_lower := strings.ToLower(device_type)
	if device_type_lower != string(vars.CLIENT_OS_ANDROID) && device_type_lower != string(vars.CLIENT_OS_IOS) && device_type_lower != string(vars.CLIENT_OS_WP) {
		c.Json(libs.NewError("member_set_pushid_parameter_error", "M3803", "客户端系统标识不被支持", ""))
		return
	}
	mp := passport.NewMemberProvider()
	err := mp.UpdatePushConfig(current_uid, int(proxy_id), channel_id, user_id, vars.CLIENT_OS(device_type_lower))
	if err != nil {
		c.Json(libs.NewError("member_set_pushid_parameter_fail", "M3805", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("member_set_pushid_succ", RESPONSE_SUCCESS, "更新成功", ""))
}

// @Title 是否已绑定第三方推送
// @Description 是否已绑定第三方推送
// @Param   access_token   path   string  true  "access_token"
// @Success 200 已绑定返回error_code:REP000
// @router /isset_pushid [get]
func (c *MemberController) IsSetPushId() {
	member := c.CurrentUser()
	if member == nil {
		c.Json(libs.NewError("member_isset_pushid_premission_denied", UNAUTHORIZED_CODE, "你没有操作的权限", ""))
		return
	}
	if member.PushProxy > 0 {
		c.Json(libs.NewError("member_isset_pushid_success", RESPONSE_SUCCESS, "已绑定推送服务", ""))
		return
	}
	c.Json(libs.NewError("member_isset_pushid_fail", REPSONSE_FAIL, "未绑定推送服务", ""))
}

// @Title 设置用户昵称
// @Description 设置用户昵称(已设置的不允许修改)
// @Param   access_token   path   string  true  "access_token"
// @Param   nick_name    path    string  true  "昵称(正则^[0-9a-zA-Z\u4e00-\u9fa5_]*$)"
// @Success 200
// @router /set_nickname [post]
func (c *MemberController) SetNickname() {
	member := c.CurrentUser()
	if member == nil {
		c.Json(libs.NewError("member_set_nickname_premission_denied", UNAUTHORIZED_CODE, "你没有操作的权限", ""))
		return
	}
	nickname, _ := utils.UrlDecode(c.GetString("nick_name"))
	if len(member.NickName) > 0 {
		c.Json(libs.NewError("member_set_nickname_notallow", "M3820", "你已设置昵称,不能重复修改", ""))
		return
	}
	mp := passport.NewMemberProvider()
	err := mp.SetNickname(member.Uid, nickname)
	if err != nil {
		c.Json(libs.NewError("member_set_nickname_fail", "M3821", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("member_set_nickname_succ", RESPONSE_SUCCESS, "成功设置昵称", ""))
}

// @Title 设置用户头像
// @Description 设置用户头像
// @Param   access_token   path   string  true  "access_token"
// @Param   file_id    path    string  true  "file_id"
// @Success 200  {object} libs.Error
// @router /set_avatar [post]
func (c *MemberController) SetAvatar() {
	file_id, _ := c.GetInt64("file_id")
	if file_id <= 0 {
		c.Json(libs.NewError("member_set_avatar_fail", "M3830", "文件id错误", ""))
		return
	}
	member := c.CurrentUser()
	if member == nil {
		c.Json(libs.NewError("member_set_avatar_premission_denied", UNAUTHORIZED_CODE, "你没有操作的权限", ""))
		return
	}
	mp := passport.NewMemberProvider()
	_, err := mp.SetMemberAvatar(member.Uid, file_id)
	if err != nil {
		c.Json(libs.NewError("member_set_avatar_fail", "M3831", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("member_set_avatar_succ", RESPONSE_SUCCESS, "成功设置头像", ""))
}

// @Title 获取用户头像原图
// @Description 获取用户头像原图
// @Param   uid    path    string  true  "uid"
// @Success 200  {object} outobjs.OutFileUrl
// @router /get_original_avatar [get]
func (c *MemberController) GetOriginalAvatar() {
	uid, _ := c.GetInt64("uid")
	mp := passport.NewMemberProvider()
	avatar := mp.GetMemberAvatar(uid, vars.PIC_SIZE_ORIGINAL)
	if avatar == nil {
		c.Json(libs.NewError("member_get_avatar_fail", "M3835", "用户未设置头像", ""))
		return
	}
	ofu := &outobjs.OutFileUrl{
		FileId:  avatar.FileId,
		FileUrl: file_storage.GetFileUrl(avatar.FileId),
	}
	c.Json(ofu)
}

// @Title 用户常用@对象列表
// @Description 用户常用@对象列表
// @Param   access_token   path   string  true  "access_token"
// @Param   size    path   int  false  "默认10"
// @Success 200
// @router /frequent_ats [get]
func (c *MemberController) FrequentAts() {
	current_uid := c.CurrentUid()
	if current_uid <= 0 {
		c.Json(libs.NewError("member_set_pushid_premission_denied", UNAUTHORIZED_CODE, "你没有操作的权限", ""))
		return
	}
	size, _ := c.GetInt("size")
	if size <= 0 {
		size = 10
	}

	freqs := passport.GetFrequentAts(current_uid, int(size))
	type OutFreqAt struct {
		NickName string `json:"nick_name"`
		Uid      int64  `json:"uid"`
	}
	mp := passport.NewMemberProvider()
	outs := []*OutFreqAt{}
	for _, toUid := range freqs {
		nick_name := mp.GetNicknameByUid(toUid)
		if len(nick_name) == 0 {
			continue
		}
		outs = append(outs, &OutFreqAt{
			NickName: nick_name,
			Uid:      toUid,
		})
	}
	c.Json(outs)
}

// @Title 用户配置信息
// @Description 用户配置信息
// @Param   access_token   path   string  true  "access_token"
// @Success 200
// @router /my_config [get]
func (c *MemberController) GetMyConfig() {
	current_uid := c.CurrentUid()
	if current_uid <= 0 {
		c.Json(libs.NewError("member_set_pushid_premission_denied", UNAUTHORIZED_CODE, "你没有操作的权限", ""))
		return
	}
	cfg := &passport.MemberConfigs{}
	mycfg := cfg.GetConfig(current_uid)
	c.Json(mycfg)
}

// @Title 设置用户配置信息
// @Description 设置用户配置信息
// @Param   access_token   path   string  true  "access_token"
// @Param   config   path   string  true  "json格式字符串"
// @Success 200
// @router /my_config [post]
func (c *MemberController) SetMyConfig() {
	current_uid := c.CurrentUid()
	if current_uid <= 0 {
		c.Json(libs.NewError("member_set_pushid_premission_denied", UNAUTHORIZED_CODE, "你没有操作的权限", ""))
		return
	}
	str := c.GetString("config")
	if len(str) == 0 {
		c.Json(libs.NewError("member_set_config_fail", "M3901", "config不能为空", ""))
		return
	}
	var mycfg passport.MemberConfigAttrs
	err := json.Unmarshal([]byte(str), &mycfg)
	if err != nil {
		c.Json(libs.NewError("member_set_config_fail", "M3902", "json格式错误:"+err.Error(), ""))
		return
	}
	cfg := &passport.MemberConfigs{}
	err = cfg.SetConfig(current_uid, &mycfg)
	if err != nil {
		c.Json(libs.NewError("member_set_config_fail", "M3903", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("member_set_config_succ", RESPONSE_SUCCESS, "成功设置配置", ""))
}

// @Title 设置用户配置信息
// @Description 设置用户配置信息
// @Param   ts	   path  int  true "时间戳"
// @Param   sign   path   string  true  "签名(5分钟有效)"
// @Param   uid   path   int  true  "用户id"
// @Param   pwd   path   string true "新密码"
// @Success 200
// @router /reset_pwd [post]
func (c *MemberController) ResetPwd() {
	ts, _ := c.GetInt64("ts")
	sign := c.GetString("sign")
	uid, _ := c.GetInt64("uid")
	pwd := c.GetString("pwd")
	if ts <= time.Now().Add(-5*time.Minute).Unix() {
		c.Json(libs.NewError("member_reset_pwd_fail", "M3921", "时间戳已经过期", ""))
		return
	}
	if uid <= 0 {
		c.Json(libs.NewError("member_reset_pwd_fail", "M3922", "uid不能小于等于0", ""))
		return
	}
	if len(pwd) == 0 || len(sign) == 0 {
		c.Json(libs.NewError("member_reset_pwd_fail", "M3923", "参数错误", ""))
		return
	}
	verify_sign := fmt.Sprintf("%d_%d_%s_7297e4cc1568906c07", ts, uid, pwd)
	if verify_sign != sign {
		c.Json(libs.NewError("member_reset_pwd_fail", "M3924", "签名错误", ""))
		return
	}
	mp := passport.NewMemberProvider()
	err := mp.ResetPassword(uid, pwd)
	if err != nil {
		c.Json(libs.NewError("member_reset_pwd_fail", "M3925", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("member_reset_pwd_fail", RESPONSE_SUCCESS, "成功更改密码", ""))
}

// @Title 根据用户名获取用户信息
// @Description 根据用户名获取用户信息
// @Param   user_name	   path  string  true "用户名"
// @Success 200 {object} outobjs.OutSimpleMember
// @router /get_by_username [get]
func (c *MemberController) GetByUserName() {
	user_name, _ := utils.UrlDecode(c.GetString("user_name"))
	user_name = utils.StripSQLInjection(user_name)
	if len(user_name) == 0 {
		c.Json(libs.NewError("member_get_byusername_fail", "M3931", "用户名不能为空", ""))
		return
	}
	mp := passport.NewMemberProvider()
	usr := mp.GetByUserName(user_name)
	if usr == nil {
		c.Json(libs.NewError("member_get_byusername_fail", "M3932", "用户不存在", ""))
		return
	}
	c.Json(outobjs.GetOutSimpleMember(usr.Uid))
}
