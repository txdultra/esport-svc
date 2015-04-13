package passport

import (
	"dbs"
	"errors"
	"fmt"
	"strings"
	"time"
	"utils"

	"github.com/astaxie/beego/orm"
)

var RoleCacheKeys = map[string]string{
	"role_cachekey":            "mobile_member_roles",
	"role_authorities":         "mobile_member_authorities",
	"role_member_member_roles": "mobile_member_roles_uid",
	"role_member_role_members": "mobile_member_role_members_roleid",
	"role_role_authorities":    "mobile_member_role_authorities",
}

type RoleService struct{}

func NewRoleService() *RoleService {
	return &RoleService{}
}

func (r *RoleService) Roles() []*Role {
	cache := utils.GetCache()
	cache_key, _ := RoleCacheKeys["role_cachekey"]
	tmp := &[]*Role{}
	err := cache.Get(cache_key, tmp)
	if err == nil {
		return *tmp
	}

	var roles []*Role
	o := dbs.NewDefaultOrm()
	_, err = o.QueryTable(&Role{}).All(&roles)
	if err == nil {
		cache.Set(cache_key, roles, utils.StrToDuration("2h"))
	}
	return roles
}

func (r *RoleService) Role(id int64) (*Role, error) {
	roles := r.Roles()
	for _, role := range roles {
		if role.Id == id {
			return role, nil
		}
	}
	return nil, errors.New("不存在对应的角色")
}

func (r *RoleService) RoleByName(name string) (*Role, error) {
	name_lower := strings.ToLower(name)
	roles := r.Roles()
	for _, role := range roles {
		if strings.ToLower(role.RoleName) == name_lower {
			return role, nil
		}
	}
	return nil, errors.New("不存在对应的角色")
}

func (r *RoleService) Authorities() []*Authority {
	cache := utils.GetCache()
	cache_key, _ := RoleCacheKeys["role_authorities"]
	tmp := &[]*Authority{}
	err := cache.Get(cache_key, tmp)
	if err == nil {
		return *tmp
	}

	var auths []*Authority
	o := dbs.NewDefaultOrm()
	_, err = o.QueryTable(&Authority{}).All(&auths)
	if err == nil {
		cache.Set(cache_key, auths, utils.StrToDuration("2h"))
	}
	return auths
}

func (r *RoleService) AuthorityByKey(key string) *Authority {
	auths := r.Authorities()
	for _, a := range auths {
		if strings.ToLower(a.AuthKey) == strings.ToLower(key) {
			return a
		}
	}
	return nil
}

func (r *RoleService) AuthorityById(id int) *Authority {
	auths := r.Authorities()
	for _, a := range auths {
		if a.Id == int64(id) {
			return a
		}
	}
	return nil
}

func (r *RoleService) MemberRoles(uid int64) []*MemberRole {
	cache := utils.GetCache()
	cache_key, _ := RoleCacheKeys["role_member_member_roles"]
	cache_key = fmt.Sprintf("%s:%d", cache_key, uid)
	tmp := &[]*MemberRole{}
	err := cache.Get(cache_key, tmp)
	if err == nil {
		return *tmp
	}

	var mrs []*MemberRole
	o := dbs.NewDefaultOrm()
	_, err = o.QueryTable(&MemberRole{}).Filter("uid", uid).All(&mrs)
	if err == nil {
		cache.Set(cache_key, mrs, utils.StrToDuration("5h"))
	}
	return mrs
}

func (r *RoleService) VerifyMemberInRoles(uid int64, roleNames []string) (bool, error) {
	for _, rname := range roleNames {
		role, _ := r.RoleByName(rname)
		if role == nil {
			continue
		}
		member_roles := r.MemberRoles(uid)
		for _, mr := range member_roles {
			if mr.RoleId == role.Id {
				if !role.Enabled {
					return false, fmt.Errorf("此角色已关闭")
				}
				if mr.Expries == 0 { //永不过期
					return true, nil
				}
				duration := time.Second * time.Duration(mr.Expries)
				exprieTime := mr.PostTime.Add(duration)
				if exprieTime.After(time.Now()) {
					return true, nil
				}
				return false, fmt.Errorf("角色授权已过期")
			}
		}
	}
	return false, fmt.Errorf("用户不在对应的角色列表中")
}

func (r *RoleService) RoleMembers(role_id int) []*MemberRole {
	cache := utils.GetCache()
	cache_key, _ := RoleCacheKeys["role_member_role_members"]
	cache_key = fmt.Sprintf("%s:%d", cache_key, role_id)
	tmp := &[]*MemberRole{}
	lst := cache.Get(cache_key, tmp)
	if lst == nil {
		return *tmp
	}

	var mrs []*MemberRole
	o := dbs.NewDefaultOrm()
	_, err := o.QueryTable(&MemberRole{}).Filter("role_id", role_id).All(&mrs)
	if err == nil {
		cache.Set(cache_key, mrs, utils.StrToDuration("5h"))
	}
	return mrs
}

func (r *RoleService) clearMemberRoleCache(uid int64) {
	cache := utils.GetCache()
	if uid > 0 {
		cache_key, _ := RoleCacheKeys["role_member_member_roles"]
		cache_key = fmt.Sprintf("%s:%d", cache_key, uid)
		cache.Delete(cache_key)
	}
}

func (r *RoleService) clearRoleMemberCache(role_id int) {
	cache := utils.GetCache()
	if role_id > 0 {
		cache_key, _ := RoleCacheKeys["role_member_role_members"]
		cache_key = fmt.Sprintf("%s:%d", cache_key, role_id)
		cache.Delete(cache_key)
	}
}

func (r *RoleService) clearRoleAuthorities(role_id int) {
	cache := utils.GetCache()
	if role_id > 0 {
		cache_key, _ := RoleCacheKeys["role_role_authorities"]
		cache_key = fmt.Sprintf("%s:%d", cache_key, role_id)
		cache.Delete(cache_key)
	}
}

func (r *RoleService) SetMemberRoles(uid int64, role_expries map[int]int) (bool, error) {
	roles := r.Roles()
	for k, _ := range role_expries {
		has := false
		for _, r := range roles {
			if r.Id == int64(k) {
				has = true
				break
			}
		}
		if !has {
			return false, errors.New(fmt.Sprintf("权限编号%d不存在", k))
		}
	}
	o := dbs.NewDefaultOrm()
	err := o.Begin()
	if err != nil {
		return false, errors.New("启动事务处理失败")
	}
	_, err = o.Raw("delete from common_member_roles where uid = ?", uid).Exec()
	if err != nil {
		return false, err
	}
	for k, v := range role_expries {
		mr := &MemberRole{
			Uid:      uid,
			RoleId:   int64(k),
			PostTime: time.Now(),
			Expries:  v,
		}
		_, err := o.Insert(mr)
		if err != nil {
			o.Rollback()
			return false, err
		}
	}
	err = o.Commit()
	if err != nil {
		return false, err
	}
	//clear cache
	r.clearMemberRoleCache(uid)
	for k, _ := range role_expries {
		r.clearRoleMemberCache(k)
	}
	return true, nil
}

func (r *RoleService) RoleAuthorities(role_id int) []*Authority {
	cache := utils.GetCache()
	cache_key, _ := RoleCacheKeys["role_role_authorities"]
	cache_key = fmt.Sprintf("%s:%d", cache_key, role_id)
	tmp := &[]*Authority{}
	err := cache.Get(cache_key, tmp)
	if err == nil {
		return *tmp
	}

	var list orm.ParamsList
	o := dbs.NewDefaultOrm()
	mrs := []*Authority{}
	_, err = o.Raw("SELECT auth_id FROM common_role_authority WHERE role_id = ?", role_id).ValuesFlat(&list)
	if err == nil {
		for _, v := range list {
			if vi, ok := v.(int); ok {
				if m := r.AuthorityById(vi); m != nil {
					mrs = append(mrs, m)
				}
			}
		}
		cache.Set(cache_key, mrs, utils.StrToDuration("5h"))
	}
	return mrs
}

////////////////////////////////////////////////////////////////////////////////
//多平台管理者映射
////////////////////////////////////////////////////////////////////////////////
type PlatManagers struct{}

func (m *PlatManagers) ck(from_uid int64, plat string) string {
	return fmt.Sprintf("mobile_manager_member_map:uid:%d:plat:%s", from_uid, plat)
}

func (m *PlatManagers) GetManagerMap(from_uid int64, plat string) *ManageMemberMap {
	cache := utils.GetCache()
	cache_key := m.ck(from_uid, plat)
	var tmp ManageMemberMap
	err := cache.Get(cache_key, &tmp)
	if err == nil {
		return &tmp
	}
	o := dbs.NewDefaultOrm()
	err = o.QueryTable(&ManageMemberMap{}).Filter("from_uid", from_uid).Filter("plat", plat).One(&tmp)
	if err != nil {
		return nil
	}
	cache.Set(cache_key, tmp, 2*time.Hour)
	return &tmp
}

func (m *PlatManagers) ExistBind(uid int64) bool {
	o := dbs.NewDefaultOrm()
	has := o.QueryTable(&ManageMemberMap{}).Filter("to_uid", uid).Exist()
	return has
}

func (m *PlatManagers) AddPlatManager(uid int64, from_uid int64, plat string) error {
	o := dbs.NewDefaultOrm()
	has := o.QueryTable(&ManageMemberMap{}).Filter("from_uid", from_uid).Filter("plat", plat).Exist()
	if has {
		o.QueryTable(&ManageMemberMap{}).Filter("from_uid", from_uid).Filter("plat", plat).Delete()
	}
	o.Raw("insert into common_manage_members(from_uid,to_uid,plat) value(?,?,?)", from_uid, uid, plat).Exec()
	cache := utils.GetCache()
	cache_key := m.ck(from_uid, plat)
	cache.Delete(cache_key)
	return nil
}
