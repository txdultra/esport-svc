package passport

import (
	//"utils"
	//"debug"
	"encoding/json"
	"errors"
	//"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type QQThridOpenID struct{}

const (
	//QQ_API_URL_ME            = "https://graph.z.qq.com/moc2/me" //wap
	QQ_API_URL_GET_USER_INFO = "https://openmobile.qq.com/user/get_simple_userinfo"
	QQ_API_KEY               = "ef5hmQ1urH549X5I"
	QQ_API_CONSUMER_KEY      = "1101974652"
)

func (qt *QQThridOpenID) GetUserInfo(token, openid string) (*ThridUserInfo, error) {
	url := QQ_API_URL_GET_USER_INFO + "?access_token=" + token + "&oauth_consumer_key=" + QQ_API_CONSUMER_KEY + "&openid=" + openid
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
		if err != nil {
			return nil, err
		}
		datas := f.(map[string]interface{})
		ret := datas["ret"].(float64)
		if int(ret) != 0 {
			return nil, errors.New("查询用户信息错误，返回ret" + strconv.Itoa(int(ret)))
		}
		userInfo := &ThridUserInfo{"", datas["nickname"].(string), "", datas["figureurl_qq_2"].(string)}
		return userInfo, nil
	}
	return nil, errors.New("QQ服务器返回状态StatusCode非200")
}
