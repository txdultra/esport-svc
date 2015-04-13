package admincp

import (
	"controllers"
	//"fmt"
	"libs"
	"libs/passport"
)

type AdminController struct {
	controllers.AuthorizeController
}

func (c *AdminController) Prepare() {
	c.AuthorizeController.Prepare()

	current_uid := c.Access_Token.Uid
	role_service := passport.NewRoleService()
	isTrue, err := role_service.VerifyMemberInRoles(current_uid, []string{"administrators", "editor"}) //暂统一后台权限
	if !isTrue {
		out_err := libs.NewError("unauthorized", controllers.UNAUTHORIZED_CODE, "您未进行后台操作授权:"+err.Error(), "")
		c.Json(out_err)
		c.StopRun()
	}
}

func (c *AdminController) CurrentMemberRoles() []*passport.Role {
	current_uid := c.Access_Token.Uid
	role_service := passport.NewRoleService()
	mrs := role_service.MemberRoles(current_uid)
	roles := []*passport.Role{}
	for _, mr := range mrs {
		role, _ := role_service.Role(int64(mr.RoleId))
		if role != nil {
			roles = append(roles, role)
		}
	}
	return roles
}

func (c *AdminController) IsCurrentMemberInRole(roleNames []string) bool {
	current_uid := c.Access_Token.Uid
	role_service := passport.NewRoleService()
	ok, _ := role_service.VerifyMemberInRoles(current_uid, roleNames)
	return ok
}
