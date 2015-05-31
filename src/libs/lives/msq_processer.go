package lives

import (
	"encoding/json"
	"fmt"
	"libs"
	"libs/passport"
	"libs/pushapi"
	"libs/vars"
	"logs"
	"time"
	//"strconv"
	"utils"
)

//用户配置信息
var mcfg *passport.MemberConfigs = &passport.MemberConfigs{}
var baidu_pusher *pushapi.BaiduPusher

type ProgramNoticeMsqProcesser struct{}

func (ProgramNoticeMsqProcesser) Run() (<-chan string, error) {
	msq_config := &libs.MsqQueueConfig{
		MsqConnUrl: libs.MsqConnectionUrl(),
		QueueName:  libs.PUSH_MSQ_QUEUE_NAME,
		Durable:    true,
		QueueMode:  libs.MsqWorkQueueMode,
		AutoAck:    true,
		MsqId:      utils.RandomStrings(25),
	}
	msq := libs.NewMsq()
	msvc, err := msq.CreateMsqService(msq_config)
	if err != nil {
		logs.Error(err.Error())
		return nil, fmt.Errorf("初始化队列服务错误:" + err.Error())
	}
	c, err := msvc.Receive(ProgramNoticeMsqProcesser{})
	return c, err
}

//IMsqMsgProcesser 接口
func (p ProgramNoticeMsqProcesser) Do(msg *libs.MsqMessage) error {
	if msg.DataType == "json" {
		var pmd libs.PushMsgData
		err := json.Unmarshal(msg.Data, &pmd)
		if err == nil {
			//用户关闭推送
			if !(mcfg.GetConfig(pmd.ToUid).AllowPush) {
				return nil
			}
			if pmd.Category == "live_program_notice" {
				go p.PushMsgToBaiduApi(&pmd)
			}
			return nil
		}
		return err
	}
	return fmt.Errorf("不能处理" + msg.DataType + "格式的队列消息")
}

func (ProgramNoticeMsqProcesser) PushMsgToBaiduApi(msg *libs.PushMsgData) (int64, error) {
	userId := msg.PushId
	pushType := pushapi.BAIDU_PUSH_TYPE_SINGLE
	channelId := msg.ChannelId

	if len(userId) == 0 {
		return 0, fmt.Errorf("未指定推送编号user_id")
	}

	if baidu_pusher == nil {
		baidu_pusher = pushapi.NewBaiduPusher(push_baidu_apikey, push_baidu_secret)
	}

	tag := ""
	messages := make(map[string]interface{})
	var deviceType pushapi.BAIDU_DEVICE_TYPE
	var messageType pushapi.BAIDU_MSG_TYPE

	switch msg.DeviceType {
	case vars.CLIENT_OS_ANDROID:
		deviceType = pushapi.BAIDU_DEVICE_TYPE_ANDROID
		messageType = pushapi.BAIDU_MSG_TYPE_NOTICE
		messages["title"] = msg.Title
		messages["description"] = msg.Content
	case vars.CLIENT_OS_IOS:
		deviceType = pushapi.BAIDU_DEVICE_TYPE_IOS
		messageType = pushapi.BAIDU_MSG_TYPE_NOTICE
		messages["description"] = msg.Content
	case vars.CLIENT_OS_WP:
		deviceType = pushapi.BAIDU_DEVICE_TYPE_WP
		messageType = pushapi.BAIDU_MSG_TYPE_NOTICE
	}
	expries_time := time.Now().Add(10 * time.Minute)
	reqid, err := baidu_pusher.PushMsg(userId, pushType, channelId, tag, deviceType, messageType, messages, expries_time, pushapi.BAIDU_DEPLOY_STATUS_REL)
	return reqid, err
}
