package passport

import (
	"time"
)

type MemberRole struct {
	Id       int64
	Uid      int64
	RoleId   int64
	PostTime time.Time
	Expries  int //0为不过期
}

func (self *MemberRole) TableName() string {
	return "common_member_roles"
}

func (self *MemberRole) TableEngine() string {
	return "INNODB"
}

func (self *MemberRole) TableUnique() [][]string {
	return [][]string{
		[]string{"Uid", "RoleId"},
	}
}

type Role struct {
	Id       int64
	RoleName string
	Icon     string
	Enabled  bool
}

func (self *Role) TableName() string {
	return "common_roles"
}

func (self *Role) TableEngine() string {
	return "INNODB"
}

type Authority struct {
	Id       int64
	AuthName string
	AuthKey  string
}

func (self *Authority) TableName() string {
	return "common_authorities"
}

func (self *Authority) TableEngine() string {
	return "INNODB"
}
