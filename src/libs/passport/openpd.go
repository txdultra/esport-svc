package passport

import (
	"dbs"
	"errors"
	"fmt"
	"time"
	"utils"

	"github.com/astaxie/beego/orm"
)

type OpenIDProvider struct{}

func NewOpenIDAuth() *OpenIDProvider {
	return &OpenIDProvider{}
}

func (os *OpenIDProvider) openid_identifier_cachekey(identifier string, openid_mark OPENID_MARK) string {
	return fmt.Sprintf("mobile_OpenID_Identifier:%s_Mark:%d", identifier, byte(openid_mark))
}

func (os *OpenIDProvider) Get(identifier string, openid_mark OPENID_MARK) (*OpenIDOAuth, error) {
	idenstr := utils.StripSQLInjection(identifier)
	cachekey := os.openid_identifier_cachekey(idenstr, openid_mark)
	cache := utils.GetCache()
	var tmp OpenIDOAuth
	err := cache.Get(cachekey, &tmp)
	if err == nil {
		return &tmp, nil
	}
	openid := &OpenIDOAuth{}
	o := dbs.NewDefaultOrm()
	qs := o.QueryTable(openid)
	err = qs.Filter("auth_type", byte(openid_mark)).Filter("auth_identifier", idenstr).One(openid)
	if err == orm.ErrNoRows {
		return nil, errors.New("the identifier not exist")
	}
	if err == nil {
		cache.Set(cachekey, *openid, 24*time.Hour)
		return openid, nil
	}
	return nil, err
}

func (os *OpenIDProvider) Create(auth *OpenIDOAuth) (int, error) {
	cache := utils.GetCache()
	cache_key := os.openid_identifier_cachekey(auth.AuthIdentifier, OPENID_MARK(auth.AuthType))
	var opa OpenIDOAuth
	err := cache.Get(cache_key, &opa)
	if err == nil {
		return 0, errors.New("auth_type and auth_identifier already exist")
	}
	tmp_auth := &OpenIDOAuth{}
	o := dbs.NewDefaultOrm()
	qs := o.QueryTable(tmp_auth)
	err = qs.Filter("auth_type", auth.AuthType).Filter("auth_identifier", auth.AuthIdentifier).One(tmp_auth)
	if err == nil {
		return 0, errors.New("auth_type and auth_identifier already exist")
	}
	id, err := o.Insert(auth)
	if err == nil {
		auth.AuthId = id
		cache.Set(cache_key, *auth, 24*time.Hour)
		return int(id), nil
	}
	return 0, err
}

func (os *OpenIDProvider) Update(auth *OpenIDOAuth) error {
	o := dbs.NewDefaultOrm()
	if num, err := o.Update(auth); err == nil && num == 1 {
		//删除原缓存
		cache := utils.GetCache()
		cache_key := os.openid_identifier_cachekey(auth.AuthIdentifier, OPENID_MARK(auth.AuthType))
		cache.Delete(cache_key)

		cache.Set(cache_key, *auth, 24*time.Hour)
		return nil
	}
	return errors.New("update fail")
}

func (os *OpenIDProvider) Delete(auth *OpenIDOAuth) error {
	o := dbs.NewDefaultOrm()
	_, err := o.Delete(auth)
	cache := utils.GetCache()
	cache_key := os.openid_identifier_cachekey(auth.AuthIdentifier, OPENID_MARK(auth.AuthType))
	cache.Delete(cache_key)
	return err
}
