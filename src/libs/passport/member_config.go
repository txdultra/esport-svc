package passport

import (
	"dbs"
	"encoding/json"
	"fmt"
	"utils/ssdb"
)

type MemberConfigs struct{}

func (mc *MemberConfigs) ckey(uid int64) string {
	return fmt.Sprintf("mobile_member_config:%d", uid)
}

func (mc *MemberConfigs) GetConfig(uid int64) *MemberConfigAttrs {
	var ucfg MemberConfigAttrs
	err := ssdb.New(use_ssdb_passport_db).Get(mc.ckey(uid), &ucfg)
	if err != nil {
		return NewMemberConfigAttrs()
	}
	return &ucfg
}

func (mc *MemberConfigs) SetConfig(uid int64, attrs *MemberConfigAttrs) error {
	btys, err := json.Marshal(*attrs)
	if err != nil {
		return err
	}
	mcfg := MemberConfig{
		Uid:     uid,
		Setting: string(btys),
	}
	o := dbs.NewDefaultOrm()
	count, _ := o.QueryTable(mcfg).Filter("uid", uid).Count()
	if count > 0 {
		o.Update(&mcfg)
	} else {
		o.Insert(&mcfg)
	}
	ssdb.New(use_ssdb_passport_db).Set(mc.ckey(uid), *attrs)
	return nil
}
