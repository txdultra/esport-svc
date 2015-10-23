//error prefix "M4"
package passport

import (
	"bytes"
	"dbs"
	"fmt"
	"image"
	"libs"
	"libs/credits/proxy"
	"libs/search"
	"libs/vars"
	"log"
	"logs"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"utils"
	"utils/ssdb"

	"github.com/astaxie/beego/orm"
	"github.com/thrift"

	"github.com/astaxie/beego/httplib"
	"github.com/disintegration/imaging"

	"libs/hook"
)

// 中国大陆手机号码正则匹配, 不是那么太精细
// 只要是 13,14,15,18 开头的 11 位数字就认为是中国手机号
//chinaMobilePattern = `^1[3458][0-9]{9}$`
// 用户昵称的正则匹配, 合法的字符有 0-9, A-Z, a-z, _, 汉字
// 字符 '_' 只能出现在中间且不能重复, 如 "__"
//nicknamePattern = `^[a-z0-9A-Z\p{Han}]+(_[a-z0-9A-Z\p{Han}]+)*$`
// 用户名的正则匹配, 合法的字符有 0-9, A-Z, a-z, _
// 第一个字母不能为 _, 0-9
// 最后一个字母不能为 _, 且 _ 不能连续
//namePattern = `^[a-zA-Z][a-z0-9A-Z]*(_[a-z0-9A-Z]+)*$`
// 电子邮箱的正则匹配, 考虑到各个网站的 mail 要求不一样, 这里匹配比较宽松
// 邮箱用户名可以包含 0-9, A-Z, a-z, -, _, .
// 开头字母不能是 -, _, .
// 结尾字母不能是 -, _, .
// -, _, . 这三个连接字母任意两个不能连续, 如不能出现 --, __, .., -_, -., _.
// 邮箱的域名可以包含 0-9, A-Z, a-z, -
// 连接字符 - 只能出现在中间, 不能连续, 如不能 --
// 支持多级域名, x@y.z, x@y.z.w, x@x.y.z.w.e
//mailPattern = `"^[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?$"`
//密码 ^(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z]{6,12}$

type LOGIN_ACTION_STATUS string

const (
	LOGIN_ACTION_STATUS_USERNAME_INVALID LOGIN_ACTION_STATUS = "用户名无效"
	LOGIN_ACTION_STATUS_PASSWORD_INVALID LOGIN_ACTION_STATUS = "密码无效"
	LOGIN_ACTION_STATUS_PASSWORD_ERROR   LOGIN_ACTION_STATUS = "用户名或密码错误"
	LOGIN_ACTION_STATUS_USER_NOTEXIST    LOGIN_ACTION_STATUS = "用户不存在"
	LOGIN_ACTION_STATUS_NOTSET_PWD       LOGIN_ACTION_STATUS = "未设置密码"
	LOGIN_ACTION_STATUS_CHECKPWD_SUCC    LOGIN_ACTION_STATUS = "验证密码成功"

	cache_nickname_uid_hashmap = "mobile_member_nickname|uid_hashmap:%d"
	cache_uid_nickname_hashmap = "mobile_member_uid|nickname_hashmap:%d"
)

//用户名验证规则
type REGISTER_USERNAME_RULE int

const (
	REGISTER_USERNAME_RULE_DENY   REGISTER_USERNAME_RULE = 9999
	REGISTER_USERNAME_RULE_MAIL   REGISTER_USERNAME_RULE = 1
	REGISTER_USERNAME_RULE_MOBILE REGISTER_USERNAME_RULE = 2
	REGISTER_USERNAME_RULE_NORMAL REGISTER_USERNAME_RULE = 4
)

//允许注册的用户规则
var AllowUserNameRules []REGISTER_USERNAME_RULE = []REGISTER_USERNAME_RULE{
	REGISTER_USERNAME_RULE_MOBILE,
	REGISTER_USERNAME_RULE_MAIL,
	REGISTER_USERNAME_RULE_NORMAL,
}

var (
	_nick_uid_map       map[string]int64 = make(map[string]int64)
	_uid_nick_map       map[int64]string = make(map[int64]string)
	_noset_nick_uid_map map[string]int64 = make(map[string]int64)
	_noset_uid_nick_map map[int64]string = make(map[int64]string)
	rwmutex             *sync.RWMutex    = new(sync.RWMutex)
	g_locker            *sync.Mutex      = new(sync.Mutex)
	ugameTbls           map[string]bool  = make(map[string]bool)
)

type MemberProvider struct{}

func NewMemberProvider() *MemberProvider {
	return &MemberProvider{}
}

func (m *MemberProvider) buildPassword(source string, salt string) string {
	return utils.Md5(utils.Md5(source) + salt)
}

func (m *MemberProvider) GetUserNameRule(userName string) REGISTER_USERNAME_RULE {
	injName := utils.StripSQLInjection(userName)
	if len(injName) == 0 {
		return REGISTER_USERNAME_RULE_DENY
	}
	if len(injName) != len(userName) {
		return REGISTER_USERNAME_RULE_DENY
	}
	if matched, _ := regexp.MatchString(libs.EMAIL_REGEX, injName); matched {
		return REGISTER_USERNAME_RULE_MAIL
	}
	if matched, _ := regexp.MatchString(libs.MOBILE_PHONE, injName); matched {
		return REGISTER_USERNAME_RULE_MOBILE
	}
	if matched, _ := regexp.MatchString(libs.USR_NAME_REGEX, injName); matched {
		if len(injName) >= default_username_minlen && len(injName) < default_username_maxlen {
			return REGISTER_USERNAME_RULE_NORMAL
		}
	}
	return REGISTER_USERNAME_RULE_DENY
}

func (m *MemberProvider) VerifyUserName(userName string) (REGISTER_USERNAME_RULE, *libs.Error) {
	//校验
	name_rule := m.GetUserNameRule(userName)
	if name_rule == REGISTER_USERNAME_RULE_DENY {
		return name_rule, &libs.Error{"user_name_verify_fail", "M4001", "用户名使用了不被允许的格式", ""}
	}
	rule_allowed := false
	for _, allow_rule := range AllowUserNameRules {
		if allow_rule == name_rule {
			rule_allowed = true
			break
		}
	}
	if !rule_allowed {
		return name_rule, &libs.Error{"user_name_verify_fail", "M4001", "用户名使用了不被允许的格式", ""}
	}
	exist := m.ExistUserName(userName)
	if exist {
		return name_rule, &libs.Error{"user_name_verify_fail", "M4011", "用户名已被占用", ""}
	}
	return name_rule, nil
}

func (m *MemberProvider) SyncVerifyMemberState(uid int64) {
	state := m.GetState(uid)
	if state == nil {
		return
	}
	fs := &FriendShips{}
	fans := fs.FollowerCounts(uid)
	friends := fs.FriendCounts(uid)
	if fans != state.Fans || friends != state.Friends {
		//有可能ssdb返回失败
		if fans > 0 || friends > 0 {
			state.Fans = fans
			state.Friends = friends
			m.UpdateState(state)
		}
	}
}

func (m *MemberProvider) Create(member Member, longiTude, latiTude float32) (int64, *libs.Error) {
	//校验
	name_rule, terr := m.VerifyUserName(member.UserName)
	if terr != nil {
		return 0, terr
	}

	////email格式
	//if name_rule == REGISTER_USERNAME_RULE_MAIL {
	//	var maps []orm.Params
	//	num, _ = o.Raw("select uid from common_member_byemail where email = ?", member.UserName).Values(&maps)
	//	if num > 0 {
	//		return 0, &libs.Error{"user_name_verify_fail", "M4012", "邮箱已被占用", ""}
	//	}
	//}
	////手机号格式
	//if name_rule == REGISTER_USERNAME_RULE_MOBILE {
	//	var maps []orm.Params
	//	num, _ = o.Raw("select uid from common_member_bymobile where mobile = ?", member.UserName).Values(&maps)
	//	if num > 0 {
	//		return 0, &libs.Error{"user_name_verify_fail", "M4013", "手机号已被占用", ""}
	//	}
	//}
	////普通格式
	//if name_rule == REGISTER_USERNAME_RULE_MAIL {
	//	exist := o.QueryTable(&Member{}).Filter("user_name", member.UserName).Exist()
	//	if exist {
	//		return 0, &libs.Error{"user_name_verify_fail", "M4011", "用户名已被占用", ""}
	//	}
	//}

	if len(member.Password) > 0 {
		salt := utils.RandomStrings(5)
		member.Salt = salt
		member.Password = m.buildPassword(member.Password, salt)
	}
	//分配唯一序列号
	member.MemberIdentifier = utils.RandomStrings(32)
	member.CreateTime = time.Now().Unix()
	o := dbs.NewDefaultOrm()
	//特殊处理
	if member.RegMode != MEMBER_REGISTER_OPENID {
		if name_rule == REGISTER_USERNAME_RULE_MAIL {
			if len(member.Email) == 0 {
				member.Email = member.UserName
			}
			member.RegMode = MEMBER_REGISTER_EMAIL
		} else if name_rule == REGISTER_USERNAME_RULE_MOBILE {
			member.RegMode = MEMBER_REGISTER_MOBILE
		} else {
			member.RegMode = MEMBER_REGISTER_NORMAL
		}
	}

	err := o.Begin()
	if err != nil {
		return 0, &libs.Error{"user_create_trans_fail", "M4004", "启动事务处理错误", ""}
	}
	defer func() {
		if e := recover(); e != nil {
			o.Rollback()
			log.Println(e)
		}
	}()
	uid, err := o.Insert(&member)
	if err == nil && uid > 0 {
		mstate := new(MemberState)
		mstate.Uid = uid
		mstate.LastLoginIP = member.CreateIP
		mstate.LastLoginTime = time.Now().Unix()
		mstate.LongiTude = longiTude
		mstate.LatiTude = latiTude
		o.Insert(mstate)
		//member profile
		profile := new(MemberProfile)
		profile.Uid = uid
		if name_rule == REGISTER_USERNAME_RULE_MOBILE {
			profile.FoundPwdMobile = member.UserName
		}
		o.Insert(profile)
	} else {
		o.Rollback()
		return 0, &libs.Error{"user_create_fail", "M4005", "创建用户错误", ""}
	}
	err = o.Commit()
	if err == nil {
		hook.Do("new_register_user", uid, 1)
		return uid, nil
	}
	return 0, &libs.Error{"user_create_commit_fail", "M4006", "创建用户错误", ""}
}

func (m *MemberProvider) ExistUserName(userName string) bool {
	o := dbs.NewDefaultOrm()
	exist := o.QueryTable(&Member{}).Filter("user_name", userName).Exist()
	if exist {
		return true
	}
	return false
}

func (m *MemberProvider) VerifyUcUserName(userName string) bool {
	injName := utils.StripSQLInjection(userName)
	if len(injName) == 0 {
		return false
	}
	if len(injName) != len(userName) {
		return false
	}
	if len(injName) <= default_username_minlen || len(injName) > default_username_maxlen {
		return false
	}
	if matched, err := regexp.MatchString(libs.USR_NAME_REGEX, injName); err != nil || !matched {
		return false
	}

	exist := m.ExistUserName(injName)
	if exist {
		return false
	}
	return true
}

func (m *MemberProvider) VerifyUserNameOfEmail(userName string) bool {
	injName := utils.StripSQLInjection(userName)
	if len(injName) == 0 {
		return false
	}
	if len(injName) != len(userName) {
		return false
	}
	if matched, err := regexp.MatchString(libs.EMAIL_REGEX, injName); err != nil || !matched {
		return false
	}

	exist := m.ExistUserName(injName)
	if exist {
		return false
	}
	return true
}

func (m *MemberProvider) VerifyUserNameOfMobile(userName string) bool {
	injName := utils.StripSQLInjection(userName)
	if len(injName) == 0 {
		return false
	}
	if len(injName) != 11 {
		return false
	}
	if matched, err := regexp.MatchString(libs.MOBILE_PHONE, injName); err != nil || !matched {
		return false
	}

	exist := m.ExistUserName(injName)
	if exist {
		return false
	}
	return true
}

func (m *MemberProvider) CheckLoginPassword(name string, pwd string) (LOGIN_ACTION_STATUS, *Member) {
	_name := utils.StripSQLInjection(name)
	if len(_name) == 0 {
		return LOGIN_ACTION_STATUS_USERNAME_INVALID, nil
	}
	if len(_name) != len(name) {
		return LOGIN_ACTION_STATUS_USERNAME_INVALID, nil
	}
	if len(pwd) == 0 {
		return LOGIN_ACTION_STATUS_PASSWORD_INVALID, nil
	}
	usr := m.GetByUserName(_name)
	if usr == nil {
		return LOGIN_ACTION_STATUS_USER_NOTEXIST, nil
	}
	if len(usr.Password) == 0 {
		return LOGIN_ACTION_STATUS_NOTSET_PWD, nil
	}
	_encode_pwd := m.buildPassword(pwd, usr.Salt)
	if _encode_pwd == usr.Password {
		return LOGIN_ACTION_STATUS_CHECKPWD_SUCC, usr
	}
	return LOGIN_ACTION_STATUS_PASSWORD_ERROR, nil
}

func (m *MemberProvider) ResetPassword(uid int64, pwd string) error {
	if len(pwd) < 6 {
		return fmt.Errorf("密码必须大于等于6位")
	}
	if len(pwd) > 16 {
		return fmt.Errorf("密码不能大于等于16位")
	}
	//matched, _ := regexp.MatchString(`^.*(\d[a-zA-Z]|[a-zA-Z]\d).*$`, pwd) //未实现
	//if !matched {
	//	return fmt.Errorf("密码必须包含字母数字")
	//}
	member := m.Get(uid)
	if member == nil {
		return fmt.Errorf("用户不存在")
	}
	salt := utils.RandomStrings(5)
	member.Salt = salt
	member.Password = m.buildPassword(pwd, salt)
	return m.Update(*member)
}

func (m *MemberProvider) ResetMail(uid int64, mail string) error {
	matched, err := regexp.MatchString(libs.EMAIL_REGEX, mail)
	if !matched || err != nil {
		return fmt.Errorf("邮箱格式错误")
	}
	member := m.Get(uid)
	if member == nil {
		return fmt.Errorf("用户不存在")
	}
	if strings.ToLower(member.Email) == strings.ToLower(mail) {
		return fmt.Errorf("新邮箱和原邮箱相同,无需重新设置")
	}
	o := dbs.NewDefaultOrm()
	exist := o.QueryTable(&Member{}).Filter("email", mail).Exist()
	if exist {
		return fmt.Errorf("已存在相同的邮箱")
	}
	member.Email = mail
	return m.Update(*member)
}

func (m *MemberProvider) CheckMemberUnCompletedProcs(uid int64, version string) []string {
	ck := fmt.Sprintf("member_procs_completed_%s_%d", version, uid)

	ok, _ := ssdb.New(use_ssdb_passport_db).Exists(ck)
	if ok {
		return nil
	}

	member := m.Get(uid)
	if member == nil {
		return nil
	}
	unprocs := []string{}
	if len(member.NickName) == 0 {
		unprocs = append(unprocs, "unset_nickname")
	}
	select_gameids := m.MemberGames(uid)
	if len(select_gameids) == 0 {
		unprocs = append(unprocs, "unset_games")
	}

	fs := &FriendShips{}
	friend_counts := fs.FriendCounts(uid)
	if friend_counts == 0 {
		unprocs = append(unprocs, "unset_friend")
	}

	//缓存
	if len(unprocs) == 0 {
		ssdb.New(use_ssdb_passport_db).Set(ck, 1)
	}
	return unprocs
}

func (m *MemberProvider) SetMemberCertified(uid int64, certified bool, certified_reson string) error {
	u := m.Get(uid)
	if u == nil {
		return fmt.Errorf("用户不存在")
	}
	u.Certified = certified
	u.CertifiedReason = certified_reson
	return m.Update(*u)
}

func (m *MemberProvider) Get(uid int64) *Member {
	if uid <= 0 {
		return nil
	}
	ckey := m.cacheKeyByUid(uid)
	cache := utils.GetCache()
	member := Member{}
	err := cache.Get(ckey, &member)
	if err == nil {
		return &member
	}
	o := dbs.NewDefaultOrm()
	member.Uid = uid
	err = o.Read(&member)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		return nil
	}
	cache.Set(ckey, member, 120*time.Hour)
	return &member
}

func (m *MemberProvider) GetByUserName(name string) *Member {
	_name := utils.StripSQLInjection(name)
	if len(_name) == 0 {
		return nil
	}
	if len(_name) != len(name) {
		return nil
	}
	ckey := m.cacheKeyByName(_name)
	cache := utils.GetCache()
	member := Member{}
	err := cache.Get(ckey, &member)
	if err == nil {
		return &member
	}
	o := dbs.NewDefaultOrm()
	err = o.QueryTable(&member).Filter("user_name", _name).One(&member)
	if err == orm.ErrMultiRows {
		log.Println("Returned Multi Rows Not One")
		return nil
	}
	if err == orm.ErrNoRows {
		return nil
	}
	cache.Set(ckey, member, 120*time.Hour)
	return &member
}

func (m *MemberProvider) Update(member Member) error {
	_oldm := m.Get(member.Uid)
	if _oldm == nil {
		return fmt.Errorf("原用户不存在")
	}
	if _oldm.NickName != member.NickName {
		return fmt.Errorf("不能通过此方法修改昵称")
	}
	o := dbs.NewDefaultOrm()
	_, err := o.Update(&member)
	if err != nil {
		return fmt.Errorf("更新失败:" + err.Error())
	}
	cache := utils.GetCache()
	cache.Replace(m.cacheKeyByName(member.UserName), member, 120*time.Hour)
	cache.Replace(m.cacheKeyByUid(member.Uid), member, 120*time.Hour)
	return nil
}

//更新用户喜好游戏(搜索用)
func (m *MemberProvider) UpdateMemberGids(uid int64, gameIds []int) {
	u := m.Get(uid)
	if u == nil {
		return
	}
	u.Gids = u.ConvertToGids(gameIds)
	o := dbs.NewDefaultOrm()
	o.Update(u, "gids")
	cache := utils.GetCache()
	cache.Replace(m.cacheKeyByName(u.UserName), *u, 120*time.Hour)
	cache.Replace(m.cacheKeyByUid(u.Uid), *u, 120*time.Hour)
}

func (m *MemberProvider) UpdatePushConfig(uid int64, pushProxy int, channelId string, pushId string, deviceType vars.CLIENT_OS) error {
	u := m.Get(uid)
	if u == nil {
		return fmt.Errorf("用户不存在")
	}
	u.PushId = pushId
	u.PushChannelId = channelId
	u.PushProxy = pushProxy
	u.DeviceType = deviceType
	return m.Update(*u)
}

func (m *MemberProvider) EmptyPushConfig(uid int64) error {
	u := m.Get(uid)
	if u == nil {
		return fmt.Errorf("用户不存在")
	}
	u.PushId = ""
	u.PushChannelId = ""
	u.PushProxy = 0
	return m.Update(*u)
}

func (m *MemberProvider) VerifyNickname(uid int64, nickName string) error {
	if len(nickName) == 0 {
		return fmt.Errorf("昵称不能为空")
	}
	matched, err := regexp.MatchString(`^[a-z0-9A-Z\p{Han}]+(_[a-z0-9A-Z\p{Han}]+)*$`, nickName)
	if !matched || err != nil {
		return fmt.Errorf("只能包含数字,字母,中文,_")
	}
	runes := []rune(nickName)
	lens := 0
	for _, r := range runes {
		if utils.IsChineseChar(r) {
			lens += 2
		} else {
			lens += 1
		}
	}
	if lens > 12 || lens < 2 {
		return fmt.Errorf("字符数不能大于12或小于2")
	}

	_str := utils.CensorWords(nickName)
	if _str != nickName {
		return fmt.Errorf("昵称含有敏感字符")
	}

	_uid := m.GetUidByNickname(nickName)
	if uid > 0 {
		if _uid > 0 {
			if _uid != uid {
				return fmt.Errorf("新昵称已被占用")
			}
		}
	} else if _uid > 0 {
		return fmt.Errorf("新昵称已被占用")
	}
	return nil
}

func (m *MemberProvider) SetNickname(uid int64, nickName string) error {
	err := m.VerifyNickname(uid, nickName)
	if err != nil {
		return err
	}

	o := dbs.NewDefaultOrm()
	nm := MemberNickName{}
	nm.Uid = uid
	err = o.Read(&nm)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		_n := MemberNickName{
			Uid:      uid,
			NickName: nickName,
		}
		_, err = o.Insert(&_n)
		if err != nil {
			return err
		}
		err = m.setCacheNameUidMap(nickName, uid)
		if err != nil {
			o.Delete(&_n)
			return fmt.Errorf("更新失败:cache1")
		}
	} else {
		_old_nick := nm.NickName
		nm.NickName = nickName
		_, err := o.Update(&nm, "nick_name")
		if err != nil {
			return fmt.Errorf("更新失败:old_nick")
		}
		err = m.setCacheNameUidMap(nickName, uid)
		if err != nil {
			nm.NickName = _old_nick
			o.Update(&nm, "nick_name")
			return fmt.Errorf("更新失败:cache2")
		}
	}
	member := m.Get(uid)
	if member != nil {
		member.NickName = nickName
		_, err := o.Update(member, "nick_name")
		if err != nil {
			return fmt.Errorf("更新失败:nt_nick")
		}
		cache := utils.GetCache()
		cache.Replace(m.cacheKeyByName(member.UserName), member, 120*time.Hour)
		cache.Replace(m.cacheKeyByUid(member.Uid), member, 120*time.Hour)
	}
	return nil
}

//预览图
func (m *MemberProvider) avatarThumbnail(fileId int64, maxPx int) (*libs.FileNode, error) {
	file_storage := libs.NewFileStorage()
	file := file_storage.GetFileNode(fileId)
	if file == nil {
		return nil, fmt.Errorf("原图片不存在")
	}
	request := httplib.Get(file_storage.GetPrivateFileUrl(fileId))
	request.SetTimeout(1*time.Minute, 1*time.Minute)
	data, err := request.Bytes()
	if err != nil {
		logs.Errorf("member avatar pic get file data fail:%s", err.Error())
		return nil, err
	}
	buf := bytes.NewBuffer(data)
	srcImg, _, err := image.Decode(buf)
	if err != nil {
		logs.Errorf("member avatar pic thumbnail by width's ratio data convert to image fail:%s", err.Error())
		return nil, err
	}
	if file.Width >= file.Height && file.Width > maxPx {
		var dstImage image.Image
		dstImage = imaging.Thumbnail(srcImg, maxPx, maxPx, imaging.Lanczos)
		fileData, err := utils.ImageToBytes(dstImage, file.OriginalName, file.ExtName)
		if err != nil {
			logs.Errorf("member avatar pic thumbnail by width's ratio image to bytes fail:%s", err.Error())
			return nil, err
		}
		ratio := float32(file.Width) / float32(maxPx)
		resizeName := fmt.Sprintf("%s-r%dx%d.%s", utils.FileName(file.OriginalName), maxPx, int(float32(file.Height)*ratio), file.ExtName)
		node, err := file_storage.SaveFile(fileData, resizeName, file.FileId)
		if err != nil {
			logs.Errorf("member avatar pic thumbnail by width's ratio save image fail:%s", err.Error())
			return nil, err
		}
		return node, nil
	}
	if file.Width < file.Height && file.Height > maxPx {
		var dstImage image.Image
		dstImage = imaging.Thumbnail(srcImg, maxPx, maxPx, imaging.Lanczos)
		fileData, err := utils.ImageToBytes(dstImage, file.OriginalName, file.ExtName)
		if err != nil {
			logs.Errorf("member avatar  pic thumbnail by height's ratio image to bytes fail:%s", err.Error())
			return nil, err
		}
		ratio := float32(file.Height) / float32(maxPx)
		resizeName := fmt.Sprintf("%s-r%dx%d.%s", utils.FileName(file.OriginalName), int(float32(file.Width)*ratio), maxPx, file.ExtName)
		node, err := file_storage.SaveFile(fileData, resizeName, file.FileId)
		if err != nil {
			logs.Errorf("member avatar  pic thumbnail by height's ratio save image fail:%s", err.Error())
			return nil, err
		}
		return node, nil
	}
	return file, nil
}

func (m *MemberProvider) cacheKeyMemberAvatar(uid int64, picSize vars.PIC_SIZE) string {
	return fmt.Sprintf("member_avatar_thumbnail_pic:%d_ps:%d", uid, picSize)
}

func (m *MemberProvider) GetMemberAvatar(uid int64, picSize vars.PIC_SIZE) *MemberAvatar {
	member := m.Get(uid)
	if member == nil {
		return nil
	}
	if member.Avatar <= 0 {
		return nil
	}
	cache := utils.GetCache()
	ma := &MemberAvatar{}
	key := m.cacheKeyMemberAvatar(uid, picSize)
	err := cache.Get(key, ma)
	if err == nil {
		return ma
	}
	o := dbs.NewDefaultOrm()
	err = o.Raw("select uid,fid,size,w,h,ts from common_member_avatars where uid=? and size=?", uid, picSize).QueryRow(ma)
	if err == nil {
		cache.Set(key, *ma, 96*time.Hour)
		return ma
	} else {
		ma, err := m.SetMemberAvatar(uid, member.Avatar)
		if err == nil {
			cache.Set(key, *ma, 96*time.Hour)
		}
		return ma
	}
	return nil
}

//设置头像，自动截取预览图
func (m *MemberProvider) SetMemberAvatar(uid int64, srcFileId int64) (*MemberAvatar, error) {
	member := m.Get(uid)
	if member == nil {
		return nil, fmt.Errorf("用户不存在")
	}
	file_storage := libs.NewFileStorage()
	file := file_storage.GetFile(srcFileId)
	if file == nil {
		return nil, fmt.Errorf("原图片不存在")
	}
	newFileNode, err := m.avatarThumbnail(srcFileId, 200)
	if err != nil {
		return nil, err
	}
	ma := &MemberAvatar{
		Uid:    uid,
		FileId: srcFileId,
		Size:   vars.PIC_SIZE_ORIGINAL,
		Width:  file.Width,
		Height: file.Height,
		Ts:     utils.TimeMillisecond(time.Now()),
	}
	o := dbs.NewDefaultOrm()
	o.Raw("delete from common_member_avatars where uid=? and size=?", uid, vars.PIC_SIZE_ORIGINAL).Exec()
	o.Raw("insert into common_member_avatars(uid,fid,size,w,h,ts) values(?,?,?,?,?,?)",
		ma.Uid,
		ma.FileId,
		ma.Size,
		ma.Width,
		ma.Height,
		ma.Ts).Exec()

	member.Avatar = newFileNode.FileId
	err = m.Update(*member)
	if err != nil {
		logs.Errorf("member set avatar update obj fail:%s", err.Error())
		return nil, err
	}
	cache := utils.GetCache()
	cache.Delete(m.cacheKeyMemberAvatar(uid, vars.PIC_SIZE_ORIGINAL))

	//hook
	hook.Do("upload_avatar", uid, 1)

	return ma, nil
}

func (m *MemberProvider) cacheKeyByName(name string) string {
	return fmt.Sprintf("mobile_member_get_byname:%s", name)
}

func (m *MemberProvider) cacheKeyByUid(uid int64) string {
	return fmt.Sprintf("mobile_member_get_byuid:%d", uid)
}

func (m *MemberProvider) cacheKeyState(uid int64) string {
	return fmt.Sprintf("mobile_member_get_state_uid:%d", uid)
}

func (m *MemberProvider) cacheKeyProfile(uid int64) string {
	return fmt.Sprintf("mobile_member_profile_uid:%d", uid)
}

//设置昵称到缓存
func (m *MemberProvider) setCacheNameUidMap(nickName string, uid int64) error {
	name_key := nickName
	uid_key := strconv.FormatInt(uid, 10)
	ckey := m.nickname_uid_hashmap_cachekey(nickName)
	_, err := ssdb.New(use_ssdb_passport_db).Hset(ckey, name_key, uid)
	if err != nil {
		return err
	}
	ckey = m.uid_nickname_hashmap_cachekey(uid)
	_, err = ssdb.New(use_ssdb_passport_db).Hset(ckey, uid_key, nickName)
	if err != nil {
		return err
	}
	delete(_noset_nick_uid_map, nickName)
	delete(_noset_uid_nick_map, uid)
	m.setLocalCacheNameUidMap(nickName, uid)
	return nil
}

//取模分区key
func (m *MemberProvider) nickname_uid_hashmap_cachekey(name string) string {
	runes := []rune(name)
	md := runes[0] % 99
	return fmt.Sprintf(cache_nickname_uid_hashmap, md)
}

//取模分区key
func (m *MemberProvider) uid_nickname_hashmap_cachekey(uid int64) string {
	md := uid % 99
	return fmt.Sprintf(cache_uid_nickname_hashmap, md)
}

func (m *MemberProvider) getInMapOfUid(nickName string) int64 {
	var uid int64 = 0
	ckey := m.nickname_uid_hashmap_cachekey(nickName)
	err := ssdb.New(use_ssdb_passport_db).Hget(ckey, nickName, &uid)
	if err == nil {
		return uid
	}
	return 0
}

func (m *MemberProvider) getInMapOfName(uid int64) string {
	name := ""
	uid_key := strconv.FormatInt(uid, 10)
	ckey := m.uid_nickname_hashmap_cachekey(uid)
	err := ssdb.New(use_ssdb_passport_db).Hget(ckey, uid_key, &name)
	if err == nil {
		return name
	}
	return ""
}

func (m *MemberProvider) setLocalCacheNameUidMap(nickName string, uid int64) {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	_nick_uid_map[nickName] = uid
	_uid_nick_map[uid] = nickName
}

func (m *MemberProvider) getLocalInMapOfUid(nickName string) int64 {
	rwmutex.RLock()
	defer rwmutex.RUnlock()
	if i, ok := _nick_uid_map[nickName]; ok {
		return i
	}
	return 0
}

func (m *MemberProvider) getLocalInMapOfName(uid int64) string {
	rwmutex.RLock()
	defer rwmutex.RUnlock()
	if n, ok := _uid_nick_map[uid]; ok {
		return n
	}
	return ""
}

func (m *MemberProvider) GetUidByNickname(nickName string) int64 {
	_name := utils.StripSQLInjection(nickName)
	if len(_name) == 0 {
		return 0
	}
	if len(_name) != len(nickName) {
		return 0
	}
	if !libs.OpenDistributed() { //分布式下不适用
		if _, ok := _noset_nick_uid_map[_name]; ok { //未设置昵称的
			return 0
		}
	}
	uid := m.getLocalInMapOfUid(_name)
	if uid > 0 {
		return uid
	}
	uid = m.getInMapOfUid(_name)
	if uid > 0 {
		m.setLocalCacheNameUidMap(_name, uid) //重设本地map
		return uid
	}
	o := dbs.NewDefaultOrm()
	nmn := MemberNickName{}
	err := o.QueryTable(&nmn).Filter("nick_name", _name).One(&nmn)
	if err == orm.ErrNoRows {
		_noset_nick_uid_map[_name] = -1 //防止反复查询数据库
		return 0
	}
	m.setCacheNameUidMap(_name, nmn.Uid)
	return nmn.Uid
}

func (m *MemberProvider) GetNicknameByUid(uid int64) string {
	if uid <= 0 {
		return ""
	}
	if !libs.OpenDistributed() { //分布式下不适用
		if _, ok := _noset_uid_nick_map[uid]; ok { //未设置昵称的
			return ""
		}
	}
	nickName := m.getLocalInMapOfName(uid)
	if len(nickName) > 0 {
		return nickName
	}
	nickName = m.getInMapOfName(uid)
	if len(nickName) > 0 {
		m.setLocalCacheNameUidMap(nickName, uid) //重设本地map
		return nickName
	}
	nmn := MemberNickName{}
	o := dbs.NewDefaultOrm()
	err := o.QueryTable(&nmn).Filter("uid", uid).One(&nmn)
	if err == orm.ErrNoRows {
		_noset_uid_nick_map[uid] = "noset" //防止反复查询数据库
		return ""
	}
	m.setCacheNameUidMap(nmn.NickName, uid)
	return nmn.NickName
}

//views 计数器
func (m MemberProvider) DoC(id int64, n int, term string) {
	state := m.GetState(id)
	if state == nil {
		return
	}
	switch term {
	case "vods":
		if state.Vods+n >= 0 {
			m.UpdateStateOnTerm(id, n, "vods")
		}
		return
	case "fans":
		if state.Fans+n >= 0 {
			m.UpdateStateOnTerm(id, n, "fans")
		}
		return
	case "friends":
		if state.Friends+n >= 0 {
			m.UpdateStateOnTerm(id, n, "friends")
		}
	case "logins":
		m.UpdateStateOnTerm(id, n, "logins")
	case "notes":
		if state.Notes+n >= 0 {
			m.UpdateStateOnTerm(id, n, "notes")
		}
	default:
		return
	}
}

func (m MemberProvider) GetC(id int64, term string) int {
	state := m.GetState(id)
	if state == nil {
		return 0
	}
	switch term {
	case "vods":
		return state.Vods
	case "fans":
		return state.Fans
	case "friends":
		return state.Friends
	case "logins":
		return state.Logins
	case "notes":
		return state.Notes
	default:
		return 0
	}
}

func (m *MemberProvider) GetState(uid int64) *MemberState {
	cache := utils.GetCache()
	ckey := m.cacheKeyState(uid)
	state := MemberState{}
	err := cache.Get(ckey, &state)
	if err == nil {
		return &state
	}
	o := dbs.NewDefaultOrm()
	state.Uid = uid
	err = o.Read(&state)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		return nil
	}
	cache.Set(ckey, state, 120*time.Hour)
	return &state
}

func (m *MemberProvider) UpdateState(state *MemberState) error {
	o := dbs.NewDefaultOrm()
	_, err := o.Update(state)
	ckey := m.cacheKeyState(state.Uid)
	utils.GetCache().Delete(ckey)
	return err
}

func (m *MemberProvider) UpdateStateOnTerm(uid int64, n int, term string) error {
	o := dbs.NewDefaultOrm()
	ms := &MemberState{}
	sql := fmt.Sprintf("update %s set %s=%s+%d where uid=%d", ms.TableName(), term, term, n, uid)
	_, err := o.Raw(sql).Exec()
	ckey := m.cacheKeyState(uid)
	utils.GetCache().Delete(ckey)
	return err
}

func (m *MemberProvider) GetProfile(uid int64) *MemberProfile {
	cache := utils.GetCache()
	ckey := m.cacheKeyProfile(uid)
	profile := MemberProfile{}
	err := cache.Get(ckey, &profile)
	if err == nil {
		return &profile
	}
	o := dbs.NewDefaultOrm()
	profile.Uid = uid
	err = o.Read(&profile)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		return nil
	}
	cache.Set(ckey, profile, 120*time.Hour)
	return &profile
}

func (m *MemberProvider) UpdateProfile(profile MemberProfile) error {
	o := dbs.NewDefaultOrm()
	_, err := o.Update(&profile)
	ckey := m.cacheKeyProfile(profile.Uid)
	utils.GetCache().Replace(ckey, profile, 120*time.Hour)
	return err
}

func (m *MemberProvider) SetMemberBackgroundImg(uid int64, bgId int64) (fn *libs.FileNode, err error) {
	file_storage := libs.NewFileStorage()
	file := file_storage.GetFileNode(bgId)
	if file == nil {
		return nil, fmt.Errorf("原图片不存在")
	}
	request := httplib.Get(file_storage.GetPrivateFileUrl(bgId))
	request.SetTimeout(1*time.Minute, 1*time.Minute)
	data, err := request.Bytes()
	if err != nil {
		logs.Errorf("member background pic get file data fail:%s", err.Error())
		return nil, err
	}
	buf := bytes.NewBuffer(data)
	srcImg, _, err := image.Decode(buf)
	if err != nil {
		logs.Errorf("member background pic convert to image fail:%s", err.Error())
		return nil, err
	}
	sigma := 3.0
	dstImage := imaging.Blur(srcImg, sigma)
	fileData, err := utils.ImageToBytes(dstImage, file.OriginalName, file.ExtName)
	if err != nil {
		logs.Errorf("member background pic convert fail:%s", err.Error())
		return nil, err
	}
	blurName := fmt.Sprintf("%s-blur-%f.%s", utils.FileName(file.OriginalName), sigma, file.ExtName)
	node, err := file_storage.SaveFile(fileData, blurName, file.FileId)
	if err == nil {
		profile := m.GetProfile(uid)
		if profile != nil {
			profile.BackgroundImg = node.FileId
			o := dbs.NewDefaultOrm()
			o.Update(profile, "bg")
			ckey := m.cacheKeyProfile(profile.Uid)
			utils.GetCache().Replace(ckey, profile, 120*time.Hour)
		}
		return node, nil
	}
	return nil, err
}

//搜索用户
func (m *MemberProvider) Query(words string, p int, s int, match_mode string, sortBy string, filters []search.SearchFilter, filterRanges []search.FilterRangeInt) (int, []int64) {
	search_config := &search.SearchOptions{
		Host:    search_member_server,
		Port:    search_member_port,
		Timeout: search_member_timeout,
	}
	offset := (p - 1) * s
	search_config.Offset = offset
	search_config.Limit = s
	search_config.Filters = filters
	search_config.FilterRangeInt = filterRanges
	search_config.MaxMatches = 500
	q_engine := search.NewSearcher(search_config)
	ids, total, err := q_engine.MemberQuery(words, match_mode, sortBy)
	if err != nil {
		return 0, []int64{}
	}
	return total, ids
}

func (m *MemberProvider) QueryForAdmin(words string, certified bool, p int, s int) (int, []int64) {
	offset := (p - 1) * s
	o := dbs.NewDefaultOrm()
	query := o.QueryTable(Member{}).Filter("nick_name__icontains", words)
	if certified {
		query = query.Filter("certified", true)
	}
	total, _ := query.Count()
	uids := []int64{}
	var maps []orm.Params
	num, err := query.OrderBy("-uid").Limit(s, offset).Values(&maps, "uid")
	if err == nil && num > 0 {
		for _, row := range maps {
			_id := row["Uid"].(int64)
			//_uid, err := strconv.ParseInt(_id, 10, 64)
			//if err == nil {
			uids = append(uids, _id)
			//}
		}
	}
	return int(total), uids
}

//用户订阅游戏分表
func (m *MemberProvider) hash_ugs_tbl(uid int64, pfx_table string) string {
	pfx := strconv.Itoa(int(uid % 9))
	tbl := fmt.Sprintf("%s_%s", pfx_table, pfx)
	if _, ok := ugameTbls[tbl]; ok {
		return tbl
	}
	g_locker.Lock()
	defer g_locker.Unlock()
	if _, ok := ugameTbls[tbl]; ok {
		return tbl
	}
	orm := dbs.NewDefaultOrm()
	create_tbl_sql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
	  uid int(11) NOT NULL,
      gid int(11) NOT NULL,
 	  post_time datetime NOT NULL,
	  PRIMARY KEY (uid,gid),
	  KEY idx_uid (uid)) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`, tbl)
	_, err := orm.Raw(create_tbl_sql).Exec()
	if err == nil {
		ugameTbls[tbl] = true
	}
	return tbl
}

func (m *MemberProvider) cacheKeyGames(uid int64) string {
	return fmt.Sprintf("mobile_member_games_uid:%d", uid)
}

func (m *MemberProvider) MemberGames(uid int64) []*MemberGame {
	ckey := m.cacheKeyGames(uid)
	mgs := []*MemberGame{}
	cache := utils.GetCache()
	err := cache.Get(ckey, &mgs)
	if err == nil {
		return mgs
	}
	o := dbs.NewDefaultOrm()
	tbl := m.hash_ugs_tbl(uid, "common_member_games")
	num, err := o.Raw("select cmg.* FROM "+tbl+" cmg join common_games cg on cmg.gid=cg.id WHERE cmg.uid = ? order by cg.display_order desc", uid).QueryRows(&mgs)
	if err == nil && num > 0 {
		cache.Add(m.cacheKeyGames(uid), mgs, 60*time.Hour)
		return mgs
	}
	return mgs
}

func (m *MemberProvider) UpdateMemberGameSingle(uid int64, gameId int) error {
	if uid <= 0 || gameId <= 0 {
		return fmt.Errorf("参数错误")
	}
	bas := &libs.Bas{}
	games := bas.Games()
	hasGid := false
	for _, game := range games {
		if game.Id == gameId {
			hasGid = true
		}
	}
	if !hasGid {
		return fmt.Errorf("不存在对应的游戏")
	}
	subrGames := m.MemberGames(uid)
	for _, g := range subrGames {
		if g.GameId == gameId { //已订阅
			return nil
		}
	}
	t := time.Now()
	tbl := m.hash_ugs_tbl(uid, "common_member_games")
	o := dbs.NewDefaultOrm()
	_, err := o.Raw("insert into "+tbl+"(uid,gid,post_time) values(?,?,?)", uid, gameId, t).Exec()
	if err == nil {
		subrGames = append(subrGames, &MemberGame{
			Uid:      uid,
			GameId:   gameId,
			PostTime: t,
		})
		//更新用户库内容
		gids := []int{}
		for _, g := range subrGames {
			gids = append(gids, g.GameId)
		}
		m.UpdateMemberGids(uid, gids)

		cache := utils.GetCache()
		cache.Replace(m.cacheKeyGames(uid), subrGames, 60*time.Hour)
	}
	return nil
}

func (m *MemberProvider) RemoveMemberGameSingle(uid int64, gameId int) error {
	if uid <= 0 || gameId <= 0 {
		return fmt.Errorf("参数错误")
	}
	bas := &libs.Bas{}
	games := bas.Games()
	hasGameId := false
	for _, game := range games {
		if game.Id == gameId {
			hasGameId = true
		}
	}
	if !hasGameId {
		return fmt.Errorf("不存在对应的游戏")
	}
	tbl := m.hash_ugs_tbl(uid, "common_member_games")
	o := dbs.NewDefaultOrm()
	_, err := o.Raw("delete from "+tbl+" where uid=? and gid=?", uid, gameId).Exec()
	if err == nil {
		subrGames := m.MemberGames(uid)
		//更新用户库内容
		newSubGames := []*MemberGame{}
		gids := []int{}
		for _, g := range subrGames {
			if g.GameId != gameId {
				newSubGames = append(newSubGames, g)
				gids = append(gids, g.GameId)
			}
		}
		m.UpdateMemberGids(uid, gids)

		cache := utils.GetCache()
		cache.Replace(m.cacheKeyGames(uid), newSubGames, 60*time.Hour)
	}
	return nil
}

//更新用户喜好游戏
func (m *MemberProvider) UpdateMemberGames(uid int64, gameIds []int) error {
	if uid <= 0 || len(gameIds) == 0 {
		return fmt.Errorf("参数错误")
	}
	//检查是否和原来一致,减少数据库更新操作
	subrGames := m.MemberGames(uid)
	if len(subrGames) == len(gameIds) {
		needups := []int{}
		for _, sg := range subrGames {
			hasup := false
			for _, _gid := range gameIds {
				if sg.GameId == _gid {
					hasup = true
					break
				}
			}
			if !hasup {
				needups = append(needups, sg.GameId)
			}
		}
		if len(needups) == 0 {
			return nil
		}
	}

	bas := &libs.Bas{}
	for _, gid := range gameIds {
		game := bas.GetGame(gid)
		if game == nil {
			return fmt.Errorf("某个游戏编号不存在")
		}
	}
	tbl := m.hash_ugs_tbl(uid, "common_member_games")
	o := dbs.NewDefaultOrm()
	o.Begin()
	defer func() {
		if e := recover(); e != nil {
			o.Rollback()
			fmt.Println(e)
		}
	}()
	_, err := o.Raw("delete from "+tbl+" where uid = ?", uid).Exec()
	if err != nil {
		o.Rollback()
		return fmt.Errorf("内部错误:UPMG100")
	}
	mgs := []*MemberGame{}
	p, err := o.Raw("insert into " + tbl + "(uid,gid,post_time) values(?,?,?)").Prepare()
	if err != nil {
		o.Rollback()
		return fmt.Errorf("内部错误:UPMG101")
	}
	for _, gid := range gameIds {
		t := time.Now()
		_, err = p.Exec(uid, gid, t)
		if err != nil {
			o.Rollback()
			return fmt.Errorf("内部错误:UPMG102")
		}
		mgs = append(mgs, &MemberGame{
			Uid:      uid,
			GameId:   gid,
			PostTime: t,
		})
	}
	p.Close()
	err = o.Commit()
	if err != nil {
		o.Rollback()
		return fmt.Errorf("内部错误:UPMG103")
	}

	//更新用户库内容
	m.UpdateMemberGids(uid, gameIds)

	//主库
	//go func() {
	//	o := dbs.NewDefaultOrm()
	//	o.Raw("delete from common_member_games where uid=?", uid).Exec()
	//	p, _ := o.Raw("insert into common_member_games(uid,gid,post_time) values(?,?,?)").Prepare()
	//	for _, gid := range gameIds {
	//		t := time.Now()
	//		p.Exec(uid, gid, t)
	//	}
	//	p.Close()
	//}()

	cache := utils.GetCache()
	cache.Replace(m.cacheKeyGames(uid), mgs, 60*time.Hour)
	return nil
}

func (m *MemberProvider) getCreditHost(priceType vars.CURRENCY_TYPE) string {
	switch priceType {
	case vars.CURRENCY_TYPE_CREDIT:
		return credit_service_host
	case vars.CURRENCY_TYPE_JING:
		return jing_service_host
	default:
		return ""
	}
}

//操作积分
func (m *MemberProvider) ActionCredit(fuid int64, uid int64, ptype vars.CURRENCY_TYPE, credits int64, desc string) (string, error) {
	if credits == 0 {
		return "", fmt.Errorf("积分不能等于0")
	}
	if fuid <= 0 {
		return "", fmt.Errorf("操作人不存在")
	}
	action := proxy.OPERATION_ACTOIN_INCR
	if credits < 0 {
		action = proxy.OPERATION_ACTOIN_DECR
		credits = -credits
	}
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	host := m.getCreditHost(ptype)
	transport, err := thrift.NewTSocket(host)
	if err != nil {
		return "", fmt.Errorf("error resolving address:%s", err.Error())
	}

	useTransport := transportFactory.GetTransport(transport)
	client := proxy.NewCreditServiceClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		return "", fmt.Errorf("Error opening socket to %s;error:%s", host, err.Error())
	}
	defer transport.Close()

	param := &proxy.OperationCreditParameter{
		Uid:    uid,
		Points: credits,
		Desc:   desc,
		Action: action,
		Ref:    "passport",
		RefId:  strconv.FormatInt(fuid, 10),
	}

	result, err := client.Do(param)
	if result.State == proxy.OPERATION_STATE_SUCCESS {
		return result.No, nil
	}
	return "", err
}
