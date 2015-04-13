//error prefix "P1"
package passport

import (
	"dbs"
	"utils"
	//"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"libs"
	"regexp"
)

func NewUcMemberProvider() *UcMemberProvider {
	return &UcMemberProvider{}
}

type UcMemberProvider struct{}

func (this *UcMemberProvider) Create(member *Member, longiTude, latiTude float32) (int64, *libs.Error) {
	//
	if !this.verifyUcUserName(member.UserName) {
		return -1, &libs.Error{"user_name_verify_fail", "P1101", "用户名非法", ""}
	}
	if !this.verifyUcPassword(member.Password) {
		return -2, &libs.Error{"user_password_verify_fail", "P1102", "密码非法", ""}
	}
	if !this.verifyUcEmail(member.Email) {
		return -3, &libs.Error{"user_email_verify_fail", "P1103", "邮箱非法", ""}
	}

	uname := member.UserName
	mail := member.Email
	regdate := member.CreateTime
	salt := utils.RandomStrings(6)
	pwd := this.makeUcPwd(member.Password, salt)
	regip := "hidden"
	if member.CreateIP > 0 {
		regip = utils.IntToIp(member.CreateIP)
	}

	o := dbs.NewUcOrm()
	res, err := o.Raw("insert uc_members(username,password,email,regip,regdate,salt) values(?,?,?,?,?,?)", uname, pwd, mail, regip, regdate, salt).Exec()
	if err == nil {
		id, err := res.LastInsertId()
		if err == nil {
			return id, nil
		}
		return id, &libs.Error{"user_create_db_fail", "P1199", err.Error(), ""}
	}
	return -99, &libs.Error{"user_create_fail", "P1199", "未知错误", ""}
}

func (this *UcMemberProvider) ResetPassword(uid int64, orginalPwd, newPwd string) (int, *libs.Error) {
	if !this.verifyUcPassword(newPwd) {
		return -1, &libs.Error{"user_password_verify_fail", "P1102", "密码非法", ""}
	}

	o := dbs.NewUcOrm()
	var maps []orm.Params
	_, err := o.Raw("select password,salt from uc_members where uid=?", uid).Values(&maps)
	if err == orm.ErrNoRows {
		return -2, &libs.Error{"user_not_exist", "P1104", "用户不存在", ""}
	}
	orpwd := this.makeUcPwd(orginalPwd, maps[0]["salt"].(string))
	if orpwd != maps[0]["password"].(string) {
		return -3, &libs.Error{"user_orginal_password_error", "P1105", "原密码错误", ""}
	}

	salt := utils.RandomStrings(6)
	pwd := this.makeUcPwd(newPwd, salt)
	res, err := o.Raw("update uc_members set password=?,salt=? where uid=?", pwd, salt, uid).Exec()
	if err == nil {
		if ns, _ := res.RowsAffected(); ns > 0 {
			return 0, nil
		}
		return -99, &libs.Error{"user_sys_error", "P1199", "未知错误", ""}
	}
	return -5, &libs.Error{"user_db_error", "P1198", err.Error(), ""}
}

func (this *UcMemberProvider) LoginByName(userName, password string) (int, *libs.Error) {
	injName := utils.StripSQLInjection(userName)
	o := dbs.NewUcOrm()
	var maps []orm.Params
	_, err := o.Raw("select password,salt from uc_members where username=?", injName).Values(&maps)
	if err == orm.ErrNoRows {
		return -1, &libs.Error{"user_not_exist", "P1104", "用户不存在", ""}
	}
	orpwd := this.makeUcPwd(password, maps[0]["salt"].(string))
	if orpwd == maps[0]["password"].(string) {
		return 0, nil
	}
	return -2, &libs.Error{"user_password_error", "P1106", "密码错误", ""}
}

func (this *UcMemberProvider) LoginByEmail(email, password string) (int, *libs.Error) {
	injMail := utils.StripSQLInjection(email)
	o := dbs.NewUcOrm()
	var maps []orm.Params
	_, err := o.Raw("select password,salt from uc_members where email=?", injMail).Values(&maps)
	if err == orm.ErrNoRows {
		return -1, &libs.Error{"user_not_exist", "P1107", "邮箱不存在", ""}
	}
	orpwd := this.makeUcPwd(password, maps[0]["salt"].(string))
	if orpwd == maps[0]["password"].(string) {
		return 0, nil
	}
	return -2, &libs.Error{"user_password_error", "P1106", "密码错误", ""}
}

func (this *UcMemberProvider) verifyUcUserName(userName string) bool {
	_ = beego.AppConfig.String("uc_username_mask_url")
	unmin, _ := beego.AppConfig.Int("uc_username_minlen")
	unmax, _ := beego.AppConfig.Int("uc_username_maxlen")
	injName := utils.StripSQLInjection(userName)
	if len(injName) == 0 {
		return false
	}
	if len(injName) != len(userName) {
		return false
	}
	if len(injName) <= unmin || len(injName) > unmax {
		return false
	}
	if matched, err := regexp.MatchString(libs.USR_NAME_REGEX, injName); err != nil || !matched {
		return false
	}

	o := dbs.NewUcOrm()
	var uid int64 = 0
	err := o.Raw("select uid from uc_members where username=?", injName).QueryRow(&uid)
	if err != orm.ErrNoRows {
		return false
	}
	return true
}

func (this *UcMemberProvider) verifyUcPassword(password string) bool {
	if len(password) == 0 {
		return false
	}
	return true
}

func (this *UcMemberProvider) verifyUcEmail(email string) bool {
	if matched, err := regexp.MatchString(libs.EMAIL_REGEX, email); err != nil || !matched {
		return false
	}

	o := dbs.NewUcOrm()
	var uid int64 = 0
	err := o.Raw("select uid from uc_members where email=?", email).QueryRow(&uid)
	if err != orm.ErrNoRows {
		return false
	}
	return true
}

func (this *UcMemberProvider) makeUcPwd(password, salt string) string {
	return utils.Md5(utils.Md5(password + salt))
}
