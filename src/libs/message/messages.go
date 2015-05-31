package message

import (
	"dbs"
	"fmt"
	"libs/vars"
	"time"
	"utils"
	"utils/ssdb"

	"logs"

	"labix.org/v2/mgo/bson"
)

func CreateMsg(fromUid int64, toUid int64, msgType vars.MSG_TYPE, content string, refId string) *MsgData {
	return &MsgData{
		Id:        bson.NewObjectId(),
		FromUid:   fromUid,
		ToUid:     toUid,
		MsgType:   msgType,
		Text:      content,
		RefId:     refId,
		PostTime:  time.Now(),
		Timestamp: utils.TimeMillisecond(time.Now()),
	}
}

func SendMsg(msg *MsgData, completedCallback func(string)) error {
	config := getMsgStorageConfig(msg.MsgType)
	session, col := dbs.MgoC(config.DbName, config.TableName)
	defer session.Close()
	err := col.Insert(msg)
	if err != nil {
		logs.Error(err.Error())
		return fmt.Errorf("发送失败:%s", err)
	}
	atboxcs := fmt.Sprintf(config.MailboxCountCacheName, msg.ToUid)

	ssdb.New(config.CacheDb).Incr(atboxcs) //@消息计数
	go incrAtMsgAlerts(msg.ToUid, config)  //用户新消息
	if completedCallback != nil {
		go completedCallback(msg.Id.Hex())
	}
	return nil
}

func SendMsgV2(fromUid int64, toUid int64, msgType vars.MSG_TYPE, content string, relId string, completedCallback func(string)) error {
	msg := CreateMsg(fromUid, toUid, msgType, content, relId)
	return SendMsg(msg, completedCallback)
}

func incrAtMsgAlerts(uid int64, msc *MsgStorageConfig) int {
	new_alert_box := fmt.Sprintf(msc.NewMsgCountCacheName, uid)
	c, _ := ssdb.New(msc.CacheDb).Incr(new_alert_box)
	return int(c)
}

func NewEventCount(uid int64, msgType vars.MSG_TYPE) int {
	config := getMsgStorageConfig(msgType)
	c := 0
	new_alert_box := fmt.Sprintf(config.NewMsgCountCacheName, uid)
	err := ssdb.New(config.CacheDb).Get(new_alert_box, &c)
	if err != nil {
		return 0
	}
	return c
}

func ResetEventCount(uid int64, msgType vars.MSG_TYPE) bool {
	config := getMsgStorageConfig(msgType)
	new_alert_box := fmt.Sprintf(config.NewMsgCountCacheName, uid)
	ok, _ := ssdb.New(config.CacheDb).Del(new_alert_box)
	return ok
}

func SendMsgs2(msgs []MsgData) error {
	if len(msgs) == 0 {
		return nil
	}
	mt := msgs[0].MsgType
	config := getMsgStorageConfig(mt)
	tblName := config.TableName
	for _, msg := range msgs {
		_mt := getMsgStorageConfig(msg.MsgType)
		if tblName != _mt.TableName {
			return fmt.Errorf("批量发送的消息类型对应的表必须一致")
		}
	}
	session, col := dbs.MgoC(config.DbName, config.TableName)
	defer session.Close()
	for _, msg := range msgs {
		err := col.Insert(msg)
		if err != nil {
			continue
		}
		atboxcs := fmt.Sprintf(config.MailboxCountCacheName, msg.ToUid)
		ssdb.New(config.CacheDb).Incr(atboxcs) //@消息计数
		go incrAtMsgAlerts(msg.ToUid, config)  //用户新消息数
	}
	return nil
}

func EmptyMsgBox(uid int64, msgType vars.MSG_TYPE) error {
	config := getMsgStorageConfig(msgType)
	session, col := dbs.MgoC(config.DbName, config.TableName)
	defer session.Close()
	fm := bson.M{}
	fm["to_uid"] = uid
	_, err := col.RemoveAll(fm)
	if err != nil {
		logs.Error(err.Error())
		return fmt.Errorf("清空失败:MSG_EMB")
	}
	atboxcs := fmt.Sprintf(config.MailboxCountCacheName, uid)
	new_alert_box := fmt.Sprintf(config.NewMsgCountCacheName, uid)
	ssdb.New(config.CacheDb).Del(atboxcs)
	ssdb.New(config.CacheDb).Del(new_alert_box)
	return nil
}

func DelMsg(uid int64, msg_id string, msgType vars.MSG_TYPE) error {
	config := getMsgStorageConfig(msgType)
	if !bson.IsObjectIdHex(msg_id) {
		return fmt.Errorf("msg_id格式错误")
	}
	session, col := dbs.MgoC(config.DbName, config.TableName)
	defer session.Close()

	msg := MsgData{}
	err := col.Find(bson.M{"_id": bson.ObjectIdHex(msg_id)}).One(&msg)
	if err != nil {
		return fmt.Errorf("消息不存在")
	}
	atboxcs := fmt.Sprintf(config.MailboxCountCacheName, msg.ToUid)
	ssdb.New(config.CacheDb).Decr(atboxcs)

	err = col.RemoveId(bson.ObjectIdHex(msg_id))
	if err != nil {
		logs.Error(err.Error())
		return fmt.Errorf("删除失败:MSG_DM")
	}
	return nil
}

func GetMsgs(uid int64, page int, size int, ts time.Time, msgType vars.MSG_TYPE) (int, []*MsgData) {
	msgs := []*MsgData{}
	if uid <= 0 {
		return 0, msgs
	}
	config := getMsgStorageConfig(msgType)
	//start := (page - 1) * size
	//end := page * size
	//获取消息总数
	atboxcs := fmt.Sprintf(config.MailboxCountCacheName, uid)
	msg_totals := 0
	mt_err := ssdb.New(config.CacheDb).Get(atboxcs, &msg_totals)

	session, col := dbs.MgoC(config.DbName, config.TableName)
	defer session.Close()
	fm := bson.M{}
	fm["to_uid"] = uid
	qs := col.Find(fm)
	//ssdb 查询失败
	if mt_err != nil {
		msg_totals, _ = qs.Count()
		ssdb.New(config.CacheDb).Set(atboxcs, msg_totals) //重置总数
	}
	fm["post_time"] = bson.M{"$lt": ts}
	qs = col.Find(fm)
	err := qs.Sort("-post_time").Limit(size).All(&msgs)
	if err != nil {
		return 0, msgs
	}
	return msg_totals, msgs
}

//thrift impl
type MessageServiceImpl struct{}

func (m *MessageServiceImpl) Send(fromUid int64, toUid int64, msgType string, content string, refId string) (err error) {
	_, ok := msgtype_storage_maps[vars.MSG_TYPE(msgType)]
	if !ok {
		return fmt.Errorf("消息类型不存在")
	}
	return SendMsgV2(fromUid, toUid, vars.MSG_TYPE(msgType), content, refId, nil)
}
