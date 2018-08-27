package pushapi

import (
	"encoding/json"
	//"fmt"
	"logs"
	"sort"
	"strconv"
	"time"
	"utils"

	"github.com/astaxie/beego/httplib"
)

type BAIDU_PUSH_TYPE int

const (
	BAIDU_PUSH_TYPE_SINGLE BAIDU_PUSH_TYPE = 1
	BAIDU_PUSH_TYPE_GROUP  BAIDU_PUSH_TYPE = 2
	BAIDU_PUSH_TYPE_ALL    BAIDU_PUSH_TYPE = 3
)

type BAIDU_DEVICE_TYPE int

const (
	BAIDU_DEVICE_TYPE_BROWER  BAIDU_DEVICE_TYPE = 1
	BAIDU_DEVICE_TYPE_PC      BAIDU_DEVICE_TYPE = 2
	BAIDU_DEVICE_TYPE_ANDROID BAIDU_DEVICE_TYPE = 3
	BAIDU_DEVICE_TYPE_IOS     BAIDU_DEVICE_TYPE = 4
	BAIDU_DEVICE_TYPE_WP      BAIDU_DEVICE_TYPE = 5
)

type BAIDU_MSG_TYPE int

const (
	BAIDU_MSG_TYPE_MSG    BAIDU_MSG_TYPE = 0
	BAIDU_MSG_TYPE_NOTICE BAIDU_MSG_TYPE = 1
)

type BAIDU_DEPLOY_STATUS int

const (
	BAIDU_DEPLOY_STATUS_DEV BAIDU_DEPLOY_STATUS = 1
	BAIDU_DEPLOY_STATUS_REL BAIDU_DEPLOY_STATUS = 2
)

//type BaiduPushMsgRequestParameter struct {
//	UserId         string
//	PushType       BAIDU_PUSH_TYPE
//	ChannelId      uint
//	Tag            string
//	DeviceType     BAIDU_DEVICE_TYPE
//	MessageType    BAIDU_MSG_TYPE
//	Messages       string
//	MsgKeys        string
//	MessageExpires uint
//	DeployStatus   BAIDU_DEPLOY_STATUS
//	Timestamp      uint
//}

//返回json对象
type BaiduResponseParams struct {
	SuccessAmount int      `json:"success_amount"`
	MsgIds        []string `json:"msgids"`
}
type BaiduResult struct {
	RequestId int64               `json:"request_id"`
	RepParams BaiduResponseParams `json:"response_params"`
}

const (
	BAIDU_PUSHMSG_URL = "http://channel.api.duapp.com/rest/2.0/channel/channel"
)

//type BaiduMessageBody struct {
//	Title                string            `json:"title"`
//	Description          string            `json:"description"`
//	AndroidOpenType      string            `json:"open_type"`
//	AndroidUrl           string            `json:"url"`
//	AndroidPKGContent    string            `json:"pkg_content"`
//	AndroidCustomContent map[string]string `json:"custom_content"`
//	IOSAps               map[string]string `json:"aps"`
//	IOSKey1              string            `json:"key1"`
//	IOSKey2              string            `json:"key2"`
//	IOSKey3              string            `json:"key3"`
//	IOSKey4              string            `json:"key4"`
//	IOSKey5              string            `json:"key5"`
//}

func NewBaiduPusher(apiKey string, secretKey string) *BaiduPusher {
	return &BaiduPusher{
		apiKey,
		secretKey,
	}
}

type BaiduPusher struct {
	ApiKey    string
	SecretKey string
}

func (bp *BaiduPusher) buildSign(httpMethod string, url string, params map[string]string) string {
	ms := NewMapSorter(params)
	sort.Sort(ms)
	str := httpMethod + url
	for _, v := range ms {
		str += v.Key + "=" + v.Val
	}
	str += bp.SecretKey
	sign := utils.Md5(utils.UrlEncode(str))
	return sign
}

func (bp *BaiduPusher) PushMsg(userId string, pushType BAIDU_PUSH_TYPE, channelId string, tag string, deviceType BAIDU_DEVICE_TYPE,
	messageType BAIDU_MSG_TYPE, messages map[string]interface{}, expries time.Time, deploy_status BAIDU_DEPLOY_STATUS, args ...string) (int64, error) {

	params := make(map[string]string)
	params["method"] = "push_msg"
	params["apikey"] = bp.ApiKey
	params["push_type"] = strconv.Itoa(int(pushType))
	if len(userId) > 0 {
		params["user_id"] = userId
	}
	if len(channelId) > 0 {
		params["channel_id"] = channelId
	}
	if len(tag) > 0 {
		params["tag"] = tag
	}
	params["device_type"] = strconv.Itoa(int(deviceType))
	params["message_type"] = strconv.Itoa(int(messageType))
	msg_body, _ := json.Marshal(messages)
	params["messages"] = string(msg_body)
	params["msg_keys"] = utils.RandomStrings(32)
	params["deploy_status"] = strconv.Itoa(int(deploy_status))
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	params["timestamp"] = timestamp
	sign := bp.buildSign("POST", BAIDU_PUSHMSG_URL, params)
	req := httplib.Post(BAIDU_PUSHMSG_URL)
	for k, v := range params {
		req.Param(k, v)
	}
	req.Param("sign", sign)

	var result BaiduResult
	err := req.ToJSON(&result)
	if err != nil {
		logs.Errorf("baidu push msg_err:%s", err.Error())
		return 0, err
	}
	return result.RequestId, nil
}

////////////////////////////////////////////////////////////////////////////////
//排序用接口
type MapSorter []Item

type Item struct {
	Key string
	Val string
}

func NewMapSorter(m map[string]string) MapSorter {
	ms := make(MapSorter, 0, len(m))
	for k, v := range m {
		ms = append(ms, Item{k, v})
	}
	return ms
}

func (ms MapSorter) Len() int {
	return len(ms)
}

func (ms MapSorter) Less(i, j int) bool {
	return ms[i].Key < ms[j].Key // 按键排序
}

func (ms MapSorter) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}
