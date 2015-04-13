package passport

import (
	"libs"
)

type IUserProvider interface {
	Create(member Member, longiTude, latiTude float32) (int64, *libs.Error)
	ResetPassword(uid int64, orginalPwd, newPwd string) (int, *libs.Error)
	LoginByName(userName, password string) (int, *libs.Error)
	LoginByEmail(email, password string) (int, *libs.Error)
}
