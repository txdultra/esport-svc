package passport

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type WeixinThridOpenID struct{}

const (
	WEIXIN_API_URL_GET_USER_INFO = "https://api.weixin.qq.com/sns/userinfo"
)

func (qt *WeixinThridOpenID) GetUserInfo(token, openid string) (*ThridUserInfo, error) {
	url := WEIXIN_API_URL_GET_USER_INFO + "?access_token=" + token + "&openid=" + openid
	resp, err := http.Get(url)
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	defer func() {
		if resp != nil && !resp.Close {
			resp.Body.Close()
		}
	}()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var f interface{}
		err = json.Unmarshal(body, &f)
		fmt.Println(string(body)) //debug
		if err != nil {
			return nil, err
		}
		datas := f.(map[string]interface{})
		if rlt, ok := datas["errcode"].(string); ok {
			return nil, errors.New("查询用户信息错误，返回数据:" + rlt)
		}
		userInfo := &ThridUserInfo{"", datas["nickname"].(string), "", datas["headimgurl"].(string)}
		return userInfo, nil
	}
	return nil, errors.New("Weixin服务器返回状态StatusCode非200")
}
