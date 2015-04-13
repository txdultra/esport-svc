package libs

import (
	"dbs"
	"errors"
	"github.com/astaxie/beego/orm"
)

//推送数据对象
type PushMsgData struct {
	ToUid      int64
	ToUserName string
	ToNickName string
	PushId     string
	ChannelId  string
	DeviceType CLIENT_OS
	Title      string
	Content    string
	Category   string
	PushType   MSG_PUSH_TYPE
}

//推送的数据类型
type MSG_PUSH_TYPE string

const (
	MSG_PUSH_TYPE_NOTICE    MSG_PUSH_TYPE = "notice"
	MSG_PUSH_TYPE_MSG       MSG_PUSH_TYPE = "msg"
	MSG_PUSH_TYPE_RICHMEDIA MSG_PUSH_TYPE = "richmedia"
)

type PushManager struct{}

func NewPushManager() *PushManager {
	return &PushManager{}
}

func (p *PushManager) State(eventId int, eventClass string) PUSH_STATE {
	o := dbs.NewDefaultOrm()
	var _state PushState
	err := o.QueryTable(&PushState{}).Filter("event_id", eventId).Filter("event_ct", eventClass).One(&_state)
	if err == orm.ErrNoRows {
		return PUSH_STATE_NOSET
	}
	return _state.State
}

func (p *PushManager) SetState(eventId int, eventClass string, state PUSH_STATE) error {
	if eventId <= 0 || len(eventClass) == 0 {
		return errors.New("参数错误")
	}
	o := dbs.NewDefaultOrm()
	var _state PushState
	err := o.QueryTable(&PushState{}).Filter("event_id", eventId).Filter("event_ct", eventClass).One(&_state)
	if err == orm.ErrNoRows {
		s := PushState{
			EventId:    eventId,
			EventClass: eventClass,
			State:      state,
		}
		_, err := o.Insert(&s)
		if err != nil {
			return err
		}
	} else {
		_state.EventId = eventId
		_state.EventClass = eventClass
		_state.State = state
		_, err := o.Update(&_state)
		if err != nil {
			return err
		}
	}
	return nil
}
