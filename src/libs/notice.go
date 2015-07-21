package libs

import (
	"dbs"
	"errors"
	"logs"
	//"github.com/astaxie/beego"
	"labix.org/v2/mgo/bson"
	//"strconv"
	"sync"
	"time"
	//"utils"
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////
//v2.0
////////////////////////////////////////////////////////////////////////////////
type NoticeEvent struct {
	Id        bson.ObjectId `bson:"_id"`
	EventId   int           `bson:"event_id"`
	FromId    int64         `bson:"from_id"`
	EventTime time.Time     `bson:"event_time"`
	AddTime   time.Time     `bson:"add_time"`
}

func NewNoticeStorager(db string, collection string) *NoticeStorager {
	storager := &NoticeStorager{}
	storager.storager_db = db
	storager.storager_collection = collection
	return storager
}

type NoticeStorager struct {
	storager_db         string
	storager_collection string
}

func (n *NoticeStorager) Subscribe(e *NoticeEvent) error {
	if e.EventId == 0 {
		return errors.New("属性SubId不能为0")
	}
	if e.FromId == 0 {
		return errors.New("属性FromId不能为0")
	}
	if len(e.Id) > 0 {
		return errors.New("Id属性不能设置值")
	}
	_id := bson.NewObjectId()
	session, c := dbs.MgoC(n.storager_db, n.storager_collection)
	defer session.Close()
	e.Id = _id
	e.AddTime = time.Now()
	err := c.Insert(&e)
	if err != nil {
		return errors.New("插入错误:" + err.Error())
	}
	return nil
}

func (n *NoticeStorager) UpdateObserver(e *NoticeEvent) error {
	if e.EventId == 0 {
		return errors.New("属性SubId不能为0")
	}
	if e.FromId == 0 {
		return errors.New("属性FromId不能为0")
	}
	if len(e.Id) == 0 {
		return errors.New("Id属性必须设置值")
	}
	session, c := dbs.MgoC(n.storager_db, n.storager_collection)
	defer session.Close()
	err := c.Update(bson.M{"_id": e.Id},
		bson.M{"$set": bson.M{
			"event_id":   e.EventId,
			"from_id":    e.FromId,
			"event_time": e.EventTime,
			"add_time":   e.AddTime,
		}})
	if err != nil {
		return errors.New("更新失败:" + err.Error())
	}
	return nil
}

func (n *NoticeStorager) RemoveObserver(fromId int64, eventId int) error {
	if eventId <= 0 {
		return errors.New("eventId 不能小于等于0")
	}
	session, c := dbs.MgoC(n.storager_db, n.storager_collection)
	defer session.Close()
	_, err := c.RemoveAll(bson.M{"from_id": fromId, "event_id": eventId})
	if err != nil {
		return errors.New("更新失败:" + err.Error())
	}
	return nil
}

func (n *NoticeStorager) RemoveEventAllObservers(eventId int) error {
	if eventId <= 0 {
		return errors.New("eventId 不能小于等于0")
	}
	session, c := dbs.MgoC(n.storager_db, n.storager_collection)
	defer session.Close()
	_, err := c.RemoveAll(bson.M{"event_id": eventId})
	if err != nil {
		return errors.New("更新失败:" + err.Error())
	}
	return nil
}

func (n *NoticeStorager) Observer(id string) *NoticeEvent {
	session, c := dbs.MgoC(n.storager_db, n.storager_collection)
	defer session.Close()
	result := NoticeEvent{}
	err := c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
	if err != nil {
		return nil
	}
	return &result
}

func (n *NoticeStorager) GetObserver(eventId int, fromId int64) *NoticeEvent {
	session, c := dbs.MgoC(n.storager_db, n.storager_collection)
	defer session.Close()
	result := NoticeEvent{}
	err := c.Find(bson.M{"event_id": eventId, "from_id": fromId}).One(&result)
	if err != nil {
		return nil
	}
	return &result
}

func (n *NoticeStorager) GetObservers(p int, s int, condition map[string]interface{}, sort string) (int, []NoticeEvent) {
	session, c := dbs.MgoC(n.storager_db, n.storager_collection)
	defer session.Close()
	fm := bson.M(condition)
	offset := (p - 1) * s
	var ns []NoticeEvent
	qs := c.Find(fm)
	counts, _ := qs.Count()
	if len(sort) > 0 {
		qs = qs.Sort(sort)
	}
	err := qs.Skip(offset).Limit(s).All(&ns)
	if err != nil {
		return 0, []NoticeEvent{}
	}
	return counts, ns
}

//通知处理类
type Noticer struct{}

var noticeTasks map[string]*time.Timer = make(map[string]*time.Timer)
var locker sync.RWMutex

func NewNoticer() *Noticer {
	return &Noticer{}
}

func (n *Noticer) SubscribeNotices(uid int64, orginalEventIds []int, eventIds []int, eventTime time.Time, storager *NoticeStorager) error {
	for _, eventId := range orginalEventIds {
		storager.RemoveObserver(uid, eventId)
	}
	for _, eventId := range eventIds {
		sb := storager.GetObserver(eventId, uid)
		if sb == nil {
			_sb := &NoticeEvent{
				EventId:   eventId,
				EventTime: eventTime,
				FromId:    uid,
			}
			storager.Subscribe(_sb)
		}
	}
	return nil
}

func (n *Noticer) eventIdentifier(eventId int, eventKey string) string {
	return fmt.Sprintf("%s_%d", eventKey, eventId)
}

func (n *Noticer) StartNoticeTimer(eventId int, eventKey string, afterDuration time.Duration, storager *NoticeStorager,
	msqCfg *MsqQueueConfig, buildMsg func(*NoticeEvent) *MsqMessage) error {

	event_identifier := n.eventIdentifier(eventId, eventKey)
	_, ok := noticeTasks[event_identifier]
	if ok {
		return errors.New("已经存在等待中的任务")
	}
	pusher := NewPushManager()
	pstate := pusher.State(eventId, eventKey)
	if pstate == PUSH_STATE_NOSET {
		pusher.SetState(eventId, eventKey, PUSH_STATE_READY)
	} else if pstate != PUSH_STATE_READY {
		return errors.New("通知线程已在:" + string(pstate) + "状态")
	}
	locker.Lock()
	defer locker.Unlock()
	if _, ok := noticeTasks[event_identifier]; !ok {
		timer := time.AfterFunc(afterDuration, func() {
			pn := storager
			msq_config := msqCfg
			msq := NewMsq()
			msvc, err := msq.CreateMsqService(msq_config)
			if err != nil {
				logs.Errorf("初始化队列服务错误:", err.Error())
				return
			}
			pusher.SetState(eventId, eventKey, PUSH_STATE_SENDING)
			i := 1
			err = msvc.BeginMulti()
			if err != nil {
				logs.Errorf("开启队列批量作业错误:", err.Error())
				return
			}
			defer msvc.EndMulti()
			for {
				_, events := pn.GetObservers(i, 200, map[string]interface{}{"event_id": eventId}, "")
				if len(events) == 0 {
					break
				}
				msgs := make([]*MsqMessage, len(events), len(events))
				for j, event := range events {
					msg := buildMsg(&event)
					if msg == nil {
						continue
					}
					msgs[j] = msg
				}
				err = msvc.MultiSend(msgs)
				if err != nil {
					logs.Errorf("队列消息批量发送失败:", err.Error())
				}
				i++
			}
			pusher.SetState(eventId, eventKey, PUSH_STATE_SENDED)
			delete(noticeTasks, event_identifier) //删除任务引用
		})
		noticeTasks[event_identifier] = timer
	} else {
		return errors.New("已经存在等待中的任务")
	}
	return nil
}

func (n *Noticer) ResetNoticeTimer(eventId int, eventKey string, afterDuration time.Duration) error {
	locker.Lock()
	defer locker.Unlock()
	event_identifier := n.eventIdentifier(eventId, eventKey)
	timer, ok := noticeTasks[event_identifier]
	if !ok {
		return errors.New("计划任务集中的计划不存在")
	}
	if timer.Reset(afterDuration) {
		return nil
	}
	return errors.New("重设计划时间错误")
}

func (n *Noticer) StopNoticeTimer(eventId int, eventKey string) error {
	locker.Lock()
	defer locker.Unlock()
	event_identifier := n.eventIdentifier(eventId, eventKey)
	timer, ok := noticeTasks[event_identifier]
	if !ok {
		return errors.New("计划任务集中的计划不存在")
	}
	ok = timer.Stop()
	if ok {
		return nil
	}
	return errors.New("无法停止")
}
