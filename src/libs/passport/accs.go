package passport

//error prefix "A1"
import (
	"dbs"
	//"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"libs"
	"time"
	"utils"
)

type AuthorizationCode struct {
	Code       string
	Expires    int
	CreateTime time.Time
}

func NewAccessTokenService() *AccessTokenService {
	return &AccessTokenService{}
}

type AccessTokenService struct{}

func (ats *AccessTokenService) codeCacheKey(code string) string {
	return "mobile_authcode_" + code
}

func (ats *AccessTokenService) accessTokenCacheKey(access_token string) string {
	return "mobile_access_token:" + access_token
}

func (ats *AccessTokenService) NewCode() string {
	authorization_code_expries_seconds, err := beego.AppConfig.Int("authorization_code_expries_seconds")
	if err != nil {
		authorization_code_expries_seconds = 300
	}
	code := utils.RandomStrings(8)
	authcode := AuthorizationCode{code, authorization_code_expries_seconds, time.Now()}
	cache_key := ats.codeCacheKey(code)
	cache := utils.GetCache()
	dur := utils.IntToDuration(authorization_code_expries_seconds, "s")
	cache.Set(cache_key, authcode, dur)
	return code
}

func (ats *AccessTokenService) GetCode(code string) *AuthorizationCode {
	cache_key := ats.codeCacheKey(code)
	cache := utils.GetCache()
	authcode := AuthorizationCode{}
	err := cache.Get(cache_key, &authcode)
	if err == nil {
		return &authcode
	}
	return nil
}

func (ats *AccessTokenService) GetAccessTokenNoCode(client_ip, client_id, client_secret string, uid int64) (*AccessToken, *libs.Error) {
	if uid <= 0 {
		return nil, libs.NewError("authorize_uid_error", "A1098", "uid 非法", "")
	}
	cache := utils.GetCache()
	token := AccessToken{}
	o := dbs.NewDefaultOrm()
	qs := o.QueryTable(&token)
	err := qs.Filter("uid", uid).One(&token)
	if err == orm.ErrNoRows {
		token.Uid = uid
		token.AccessToken = utils.RandomStrings(32)
		token.ExpiresIn = authorization_access_token_expries_seconds
		token.LastTime = time.Now()
		token.LoginIp = utils.IpToInt(client_ip)
		token.App = 1 //默认
		id, i_err := o.Insert(&token)
		if i_err == nil {
			token.Id = id
			//缓存access token
			cache_key := ats.accessTokenCacheKey(token.AccessToken)
			dur := utils.IntToDuration(authorization_access_token_expries_seconds, "s")
			cache.Set(cache_key, token, dur)
			return &token, nil
		}
		return nil, libs.NewError("authorize_sys_error", "A1099", "system error", "")
	}
	if err == nil {
		//删除原缓存
		cache_key := ats.accessTokenCacheKey(token.AccessToken)
		cache.Delete(cache_key)

		token.AccessToken = utils.RandomStrings(32)
		token.LastTime = time.Now()
		token.LoginIp = utils.IpToInt(client_ip)
		cache_key = ats.accessTokenCacheKey(token.AccessToken)
		if num, u_err := o.Update(&token); u_err == nil && num == 1 {
			cache.Set(cache_key, token, utils.IntToDuration(token.ExpiresIn, "s"))
			return &token, nil
		}
	}
	return nil, libs.NewError("authorize_sys_error", "A1099", "system error", "")
}

func (ats *AccessTokenService) RefreshAccessTokenExpries(token *AccessToken) {
	orignalTime := token.LastTime
	token.LastTime = time.Now()
	cache_key := ats.accessTokenCacheKey(token.AccessToken)
	o := dbs.NewDefaultOrm()
	if num, u_err := o.Update(token); u_err == nil && num == 1 {
		cache := utils.GetCache()
		cache.Set(cache_key, *token, utils.IntToDuration(token.ExpiresIn, "s"))
	} else {
		token.LastTime = orignalTime
	}
}

func (ats *AccessTokenService) RevokeAccessToken(token *AccessToken) {
	cache_key := ats.accessTokenCacheKey(token.AccessToken)
	o := dbs.NewDefaultOrm()
	o.Delete(token)
	cache := utils.GetCache()
	cache.Delete(cache_key)
}

func (ats *AccessTokenService) GetAccessToken(code, client_ip, client_id, client_secret string, uid int64) (*AccessToken, *libs.Error) {
	acode := ats.GetCode(code)
	if acode == nil {
		return nil, libs.NewError("authorize_code_not_exist", "A1001", "authorize code not exist", "")
	}
	duration := utils.IntToDuration(acode.Expires, "s")
	exprieTime := acode.CreateTime.Add(duration)
	if exprieTime.Before(time.Now()) {
		return nil, libs.NewError("authorize_code_timeout", "A1002", "authorize code timeout", "")
	}

	return ats.GetAccessTokenNoCode(client_ip, client_id, client_secret, uid)
}

func (ats *AccessTokenService) GetTokenObj(access_token string) (*AccessToken, *libs.Error) {
	_access_token := utils.StripSQLInjection(access_token)
	cache := utils.GetCache()

	_temp_expired := ""
	cache.Get(_access_token, &_temp_expired)
	if _temp_expired == "A1004" {
		return nil, libs.NewError("authorize_access_token_not_exist", "A1004", "access_token not exist", "")
	}

	cache_key := ats.accessTokenCacheKey(_access_token)
	accessToken := AccessToken{}
	err := cache.Get(cache_key, &accessToken)
	if err == nil {
		if ats.accessTokenExpired(&accessToken) {
			return nil, libs.NewError("authorize_access_token_timeout", "A1003", "timeout", "")
		}
		return &accessToken, nil
	}

	o := dbs.NewDefaultOrm()
	qs := o.QueryTable(&accessToken)
	err = qs.Filter("access_token", _access_token).One(&accessToken)
	if err == nil {
		cache.Set(cache_key, accessToken, utils.IntToDuration(accessToken.ExpiresIn, "s"))
		if ats.accessTokenExpired(&accessToken) {
			return nil, libs.NewError("authorize_access_token_timeout", "A1003", "timeout", "")
		}
		return &accessToken, nil
	}
	if err == orm.ErrNoRows {
		cache.Set(_access_token, "A1004", 12*time.Hour)
		return nil, libs.NewError("authorize_access_token_not_exist", "A1004", "access_token not exist", "")
	}
	return nil, libs.NewError("authorize_sys_error", "A1099", "system error", "")
}

//是否过期
func (ats *AccessTokenService) accessTokenExpired(accessToken *AccessToken) bool {
	duration := time.Second * time.Duration(accessToken.ExpiresIn)
	exprieTime := accessToken.LastTime.Add(duration)
	//fmt.Println(exprieTime, "->", time.Now())
	if exprieTime.Before(time.Now()) {
		return true
	}
	return false
}
