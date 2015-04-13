package message

import (
	"dbs"
	"fmt"
	"time"
	"utils"
	"utils/ssdb"

	"logs"

	"labix.org/v2/mgo/bson"
)

const (
	member_atmsg_box_count = "member_atmsg_box_count:%d"
	member_atmsg_newalert  = "member_atmsg_newalert:%d"
)

func CreateMsg(from_uid int64, to_uid int64, msg_type MSG_TYPE, content string, ref_id string) *MsgData {
	return &MsgData{
		Id:        bson.NewObjectId(),
		FromUid:   from_uid,
		ToUid:     to_uid,
		MsgType:   msg_type,
		Text:      content,
		RefId:     ref_id,
		PostTime:  time.Now(),
		Timestamp: utils.TimeMillisecond(time.Now()),
	}
}

func SendMsgV2(from_uid int64, to_uid int64, msg_type MSG_TYPE, content string, rel_id string, completed_callback func(string)) error {
	msg := CreateMsg(from_uid, to_uid, msg_type, content, rel_id)
	return SendMsg(msg, completed_callback)
}

func SendMsg(msg *MsgData, completed_callback func(string)) error {
	session, col := dbs.MgoC(sns_atmsg_db, sns_atmsg_collection)
	defer session.Close()
	err := col.Insert(msg)
	if err != nil {
		logs.Error(err.Error())
		return fmt.Errorf("发送失败:%s", err)
	}
	atboxcs := fmt.Sprintf(member_atmsg_box_count, msg.ToUid)

	ssdb.New(use_ssdb_message_db).Incr(atboxcs) //@消息计数
	go incrAtMsgAlerts(msg.ToUid)               //用户新消息
	if completed_callback != nil {
		go completed_callback(msg.Id.Hex())
	}
	return nil
}

func incrAtMsgAlerts(uid int64) int {
	new_alert_box := fmt.Sprintf(member_atmsg_newalert, uid)
	c, _ := ssdb.New(use_ssdb_message_db).Incr(new_alert_box)
	return int(c)
}

func NewEventCount(uid int64) int {
	c := 0
	new_alert_box := fmt.Sprintf(member_atmsg_newalert, uid)
	err := ssdb.New(use_ssdb_message_db).Get(new_alert_box, &c)
	if err != nil {
		return 0
	}
	return c
}

func ResetEventCount(uid int64) bool {
	new_alert_box := fmt.Sprintf(member_atmsg_newalert, uid)
	ok, _ := ssdb.New(use_ssdb_message_db).Del(new_alert_box)
	return ok
}

func SendMsgs(msgs []MsgData) {
	session, col := dbs.MgoC(sns_atmsg_db, sns_atmsg_collection)
	defer session.Close()
	for _, msg := range msgs {
		err := col.Insert(msg)
		if err != nil {
			continue
		}
		atboxcs := fmt.Sprintf(member_atmsg_box_count, msg.ToUid)
		ssdb.New(use_ssdb_message_db).Incr(atboxcs) //@消息计数
		go incrAtMsgAlerts(msg.ToUid)               //用户新消息数
	}
}

func EmptyMsgBox(uid int64) error {
	session, col := dbs.MgoC(sns_atmsg_db, sns_atmsg_collection)
	defer session.Close()
	fm := bson.M{}
	fm["to_uid"] = uid
	_, err := col.RemoveAll(fm)
	if err != nil {
		logs.Error(err.Error())
		return fmt.Errorf("清空失败:MSG_EMB")
	}
	atboxcs := fmt.Sprintf(member_atmsg_box_count, uid)
	new_alert_box := fmt.Sprintf(member_atmsg_newalert, uid)
	ssdb.New(use_ssdb_message_db).Del(atboxcs)
	ssdb.New(use_ssdb_message_db).Del(new_alert_box)
	return nil
}

func DelMsg(uid int64, msg_id string) error {
	if !bson.IsObjectIdHex(msg_id) {
		return fmt.Errorf("msg_id格式错误")
	}
	session, col := dbs.MgoC(sns_atmsg_db, sns_atmsg_collection)
	defer session.Close()

	msg := MsgData{}
	err := col.Find(bson.M{"_id": bson.ObjectIdHex(msg_id)}).One(&msg)
	if err != nil {
		return fmt.Errorf("消息不存在")
	}
	atboxcs := fmt.Sprintf(member_atmsg_box_count, msg.ToUid)
	ssdb.New(use_ssdb_message_db).Decr(atboxcs)

	err = col.RemoveId(bson.ObjectIdHex(msg_id))
	if err != nil {
		logs.Error(err.Error())
		return fmt.Errorf("删除失败:MSG_DM")
	}
	return nil
}

func GetMsgs(uid int64, page int, size int, ts time.Time) (int, []*MsgData) {
	msgs := []*MsgData{}
	if uid <= 0 {
		return 0, msgs
	}
	//start := (page - 1) * size
	//end := page * size
	//获取消息总数
	atboxcs := fmt.Sprintf(member_atmsg_box_count, uid)
	msg_totals := 0
	mt_err := ssdb.New(use_ssdb_message_db).Get(atboxcs, &msg_totals)

	session, col := dbs.MgoC(sns_atmsg_db, sns_atmsg_collection)
	defer session.Close()
	fm := bson.M{}
	fm["to_uid"] = uid
	qs := col.Find(fm)
	//ssdb 查询失败
	if mt_err != nil {
		msg_totals, _ = qs.Count()
		ssdb.New(use_ssdb_message_db).Set(atboxcs, msg_totals) //重置总数
	}
	fm["post_time"] = bson.M{"$lt": ts}
	qs = col.Find(fm)
	err := qs.Sort("-post_time").Limit(size).All(&msgs)
	if err != nil {
		return 0, msgs
	}
	return msg_totals, msgs
}
