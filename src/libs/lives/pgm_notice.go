package lives

import (
	"dbs"
	"errors"
	"fmt"
	//"github.com/astaxie/beego"
	"libs"
	"strconv"

	"labix.org/v2/mgo/bson"
	//"sync"
	"encoding/json"
	"libs/hook"
	"libs/passport"
	"time"
	"utils"
	"utils/ssdb"
)

//节目单通知
var program_notice_db, program_notice_collection, program_notice_orginal_collection string

//var noticeTasks map[int64]*time.Timer = make(map[int64]*time.Timer)
var min_duration float64 = 10.0 //最小开启提醒间隔,5
//var locker sync.RWMutex

//var once sync.Once

const (
	member_program_new_notice_counts = "member_program_new_notice_counts:%d"
)

//提醒策略
type PROGRAM_NOTICE_STRATEGY int

const (
	PROGRAM_NOTICE_STRATEGY_SEQUENCE PROGRAM_NOTICE_STRATEGY = 1
)

type ProgramNoticeStrategy interface {
	CalculateEffectiveSubscriptionIds(pid int64, scribeIds []int64) []int64
}

type SequenceProgramNoticeStrategy struct{}

func FilterProgramNoticeStrategy(strategy PROGRAM_NOTICE_STRATEGY) ProgramNoticeStrategy {
	return &SequenceProgramNoticeStrategy{}
}

type sequenceProgramNoticeData struct {
	Pm      *LiveSubProgram
	Seleted bool
}

func (s *SequenceProgramNoticeStrategy) CalculateEffectiveSubscriptionIds(pid int64, scribeIds []int64) []int64 {
	subpm := &LiveSubPrograms{}
	subps := subpm.Gets(pid)
	spds := []sequenceProgramNoticeData{}
	//已排序分类的所有
	for _, pm := range subps {
		selected := false
		for _, id := range scribeIds {
			if id == pm.Id {
				selected = true
			}
		}
		spds = append(spds, sequenceProgramNoticeData{
			Pm:      pm,
			Seleted: selected,
		})
	}
	//根据进行需要提醒项目筛选
	result := []int64{}
	prev := false
	for _, sp := range spds {
		if sp.Seleted {
			if !prev {
				result = append(result, sp.Pm.Id)
				prev = true
			}
		} else {
			prev = false
		}
	}
	return result
}

//用户订阅提醒容器名
func member_subscribe_program_ids_set_name(uid int64) string {
	return fmt.Sprintf("member_subscribe_programs_table:%d", uid)
}

//节目单提醒
type ProgramNoticer struct{}

func NewProgramNoticeService() *ProgramNoticer {
	return &ProgramNoticer{}
}

func (n *ProgramNoticer) addToCache(uid int64, ids []int64) {
	addids := []interface{}{}
	addtms := []int64{}
	t := time.Now()
	for _, id := range ids {
		addids = append(addids, id)
		addtms = append(addtms, t.Unix())
	}
	ssdb.New(use_ssdb_live_db).MultiZadd(member_subscribe_program_ids_set_name(uid), addids, addtms)
	end := t.Add(-30 * 24 * time.Hour) //一月前的数据清除
	ssdb.New(use_ssdb_live_db).Zremrangebyscore(member_subscribe_program_ids_set_name(uid), int64(0), int64(end.Unix()))
}

func (n *ProgramNoticer) remOnCache(uid int64, ids []int64) {
	args := []interface{}{}
	for _, id := range ids {
		args = append(args, id)
	}
	ssdb.New(use_ssdb_live_db).MultiZdel(member_subscribe_program_ids_set_name(uid), args)
}

//保存用户原始提醒列表
func (n *ProgramNoticer) SaveOrginalNotices(uid int64, programId int64, subscribeIds []int64) error {
	session, c := dbs.MgoC(program_notice_db, program_notice_orginal_collection)
	defer session.Close()
	ogn := CommitOriginalNotices{}
	err := c.Find(bson.M{"event_id": programId, "from_id": uid}).One(&ogn)
	if err != nil {
		ogn.EventId = programId
		ogn.FromId = uid
		ogn.RefIds = subscribeIds
		ogn.LastTime = time.Now()
		ogn.ID = bson.NewObjectId()
		err = c.Insert(&ogn)
		if err != nil {
			return errors.New("插入错误:" + err.Error())
		}
		//添加新提醒单到缓冲库
		n.addToCache(uid, subscribeIds)
	} else {
		//删除缓冲库中的原提醒单
		n.remOnCache(uid, ogn.RefIds)

		err = c.Update(bson.M{"_id": ogn.ID},
			bson.M{"$set": bson.M{
				"event_id":  ogn.EventId,
				"from_id":   ogn.FromId,
				"refids":    subscribeIds,
				"last_time": time.Now(),
			}})
		if err != nil {
			return errors.New("更新失败:" + err.Error())
		}
		//添加新提醒单到缓冲库
		n.addToCache(uid, subscribeIds)
	}
	return nil
}

func (n *ProgramNoticer) SaveOrginalNoticeSingle(uid int64, programId int64, subscribeId int64) error {
	originalNotices := n.GetOrginalNotices(uid, programId)
	if originalNotices == nil {
		return n.SaveOrginalNotices(uid, programId, []int64{subscribeId})
	} else {
		for _, refId := range originalNotices.RefIds {
			if refId == subscribeId {
				return nil //已订阅
			}
		}
		originalNotices.RefIds = append(originalNotices.RefIds, subscribeId)
		return n.SaveOrginalNotices(uid, programId, originalNotices.RefIds)
	}
}

func (n *ProgramNoticer) RemoveOrginalNoticeSingle(uid int64, programId int64, subscribeId int64) error {
	originalNotices := n.GetOrginalNotices(uid, programId)
	if originalNotices == nil {
		return nil
	} else {
		refIds := []int64{}
		for _, refId := range originalNotices.RefIds {
			if refId != subscribeId {
				refIds = append(refIds, refId)
			}
		}
		return n.SaveOrginalNotices(uid, programId, refIds)
	}
}

func (n *ProgramNoticer) GetOrginalNotices(uid int64, programId int64) *CommitOriginalNotices {
	session, c := dbs.MgoC(program_notice_db, program_notice_orginal_collection)
	defer session.Close()
	ogn := &CommitOriginalNotices{}
	err := c.Find(bson.M{"event_id": programId, "from_id": uid}).One(ogn)
	if err == nil {
		return ogn
	}
	return nil
}

func (n *ProgramNoticer) validateLock(subProgramId int64) bool {
	lsp := &LiveSubPrograms{}
	return lsp.IsLocked(subProgramId)
}

////////////////////////////////////////////////////////////////////////////////
//v2.0
////////////////////////////////////////////////////////////////////////////////
//subscribeEventIds 只传需要提醒的id(二级节目单id),必须同属于一个programId下
func (n *ProgramNoticer) SubscribeNotice(uid int64, programId int64, subscribeEventIds []int64) error {
	lp := &LivePrograms{}
	program := lp.Get(programId)
	if program == nil {
		return errors.New("不存在节目单")
	}
	lsp := &LiveSubPrograms{}
	for _, sid := range subscribeEventIds {
		_sp := lsp.Get(sid)
		if _sp == nil || _sp.ProgramId != programId {
			return errors.New("id:" + strconv.Itoa(int(sid)) + "二级节目单不存在或不属于同一上级节目单")
		}
	}
	for _, sid := range subscribeEventIds {
		if n.validateLock(sid) {
			return errors.New("某些节目单已超过可设置时间,不能设置提醒")
		}
	}
	sn := libs.NewNoticeStorager(program_notice_db, program_notice_collection)
	originalNotices := n.GetOrginalNotices(uid, programId)
	if originalNotices != nil {
		for _, eventId := range originalNotices.RefIds {
			sn.RemoveObserver(uid, int(eventId))
		}
		//删除缓冲库中原id
		n.remOnCache(uid, originalNotices.RefIds)
	}
	//添加新的提醒
	for _, eid := range subscribeEventIds {
		_sp := lsp.Get(eid)
		sn.Subscribe(&libs.NoticeEvent{
			EventId:   int(eid),
			EventTime: _sp.StartTime,
			FromId:    uid,
			AddTime:   time.Now(),
		})
	}
	//添加新的入缓冲库
	n.addToCache(uid, subscribeEventIds)
	return nil
}

//是否已订阅
func (n *ProgramNoticer) IsSubsuribed(uid int64, subscribeEventId int64) bool {
	ok, _ := ssdb.New(use_ssdb_live_db).Zexists(member_subscribe_program_ids_set_name(uid), subscribeEventId)
	return ok
}

func (n *ProgramNoticer) SubscribeNoticeSingle(uid int64, programId int64, subscribeEventId int64) error {
	lp := &LivePrograms{}
	program := lp.Get(programId)
	if program == nil {
		return errors.New("不存在节目单")
	}
	lsp := &LiveSubPrograms{}
	_sp := lsp.Get(subscribeEventId)
	if _sp == nil || _sp.ProgramId != programId {
		return errors.New("提醒对象不存在或不属于同一上级节目单")
	}
	if n.validateLock(subscribeEventId) {
		return errors.New("节目单已超过可设置时间,不能设置提醒")
	}
	sn := libs.NewNoticeStorager(program_notice_db, program_notice_collection)
	noticeEvent := sn.GetObserver(int(subscribeEventId), uid)
	if noticeEvent != nil {
		return nil //已在提醒列表
	}

	//添加新的提醒
	sn.Subscribe(&libs.NoticeEvent{
		EventId:   int(subscribeEventId),
		FromId:    uid,
		EventTime: _sp.StartTime, //事件触发时间点
		AddTime:   time.Now(),
	})
	//添加新的入缓冲库
	n.addToCache(uid, []int64{subscribeEventId})
	//新提醒计数器
	n.incrMemberNotice(uid, 1)
	//hook
	go hook.Do("live_sub_program", uid, 1)
	return nil
}

func (n *ProgramNoticer) GetSubscribeNotices(uid int64, afterTime time.Time, page int, size int) (int, []*LiveSubProgram) {
	if uid <= 0 {
		return 0, []*LiveSubProgram{}
	}
	sn := libs.NewNoticeStorager(program_notice_db, program_notice_collection)
	condition := make(map[string]interface{})
	condition["from_id"] = uid
	condition["event_time"] = bson.M{"$gte": afterTime}
	count, list := sn.GetObservers(page, size, condition, "event_time")
	lsp := &LiveSubPrograms{}
	sps := []*LiveSubProgram{}
	for _, e := range list {
		subp := lsp.Get(int64(e.EventId))
		if subp == nil {
			count--
			continue
		}
		sps = append(sps, subp)
	}
	return count, sps
}

func (n *ProgramNoticer) RemoveSubscribeNoticeSingle(uid int64, programId int64, subscribeEventId int64) error {
	lp := &LivePrograms{}
	program := lp.Get(programId)
	if program == nil {
		return errors.New("不存在节目单")
	}
	lsp := &LiveSubPrograms{}
	_sp := lsp.Get(subscribeEventId)
	if _sp == nil || _sp.ProgramId != programId {
		return errors.New("提醒对象不存在或不属于同一上级节目单")
	}
	if n.validateLock(subscribeEventId) {
		return errors.New("节目单已超过可设置时间,不能设置提醒")
	}
	sn := libs.NewNoticeStorager(program_notice_db, program_notice_collection)
	sn.RemoveObserver(uid, int(subscribeEventId))
	//删除缓冲库中原id
	n.remOnCache(uid, []int64{subscribeEventId})
	//减少提醒计数器
	n.decrMemberNotice(uid, 1)
	return nil
}

func (n *ProgramNoticer) RemoveAllSubscribeEvent(eventId int64) error {
	sn := libs.NewNoticeStorager(program_notice_db, program_notice_collection)
	sn.RemoveEventAllObservers(int(eventId))
	err := n.StopNoticeTimer(eventId)
	return err
}

func (n *ProgramNoticer) incrMemberNotice(uid int64, cs int64) int {
	key := fmt.Sprintf(member_program_new_notice_counts, uid)
	c, _ := ssdb.New(use_ssdb_live_db).Incrby(key, cs)
	return int(c)
}

func (n *ProgramNoticer) decrMemberNotice(uid int64, cs int64) int {
	key := fmt.Sprintf(member_program_new_notice_counts, uid)
	c, _ := ssdb.New(use_ssdb_live_db).Decrby(key, cs)
	return int(c)
}

func (n *ProgramNoticer) NewEventCount(uid int64) int {
	c := 0
	key := fmt.Sprintf(member_program_new_notice_counts, uid)
	err := ssdb.New(use_ssdb_live_db).Get(key, &c)
	if err != nil {
		return 0
	}
	return c
}

func (n *ProgramNoticer) ResetEventCount(uid int64) bool {
	key := fmt.Sprintf(member_program_new_notice_counts, uid)
	ok, _ := ssdb.New(use_ssdb_live_db).Del(key)
	return ok
}

func (n *ProgramNoticer) StartNoticeTimer(subProgramId int64) error {
	if subProgramId <= 0 {
		return errors.New("subProgramId参数不能小于0")
	}
	lsp := &LiveSubPrograms{}
	sprm := lsp.Get(subProgramId)
	if sprm == nil {
		return errors.New("参数提供的节目预告不存在")
	}
	if sprm.StartTime.Before(time.Now()) {
		return errors.New("提醒时间已过,添加后也无效果")
	}
	dur := sprm.StartTime.Sub(time.Now())
	if dur.Minutes() < min_duration && dur.Minutes() > 0.0 {
		return errors.New("过于接近开始时间")
	}
	_sub_t := utils.IntToDuration(int(min_duration), "m")
	dur = sprm.StartTime.Sub(time.Now().Add(_sub_t)) //提前duration分钟
	//if subProgramId == 20 {                          //debug
	//	dur = 2 * time.Minute
	//}

	msq_config := &libs.MsqQueueConfig{
		MsqConnUrl: libs.MsqConnectionUrl(), //beego.AppConfig.String("msq.rabbitmq.addr"),
		QueueName:  libs.PUSH_MSQ_QUEUE_NAME,
		Durable:    true,
		QueueMode:  libs.MsqWorkQueueMode,
		AutoAck:    true,
		MsqId:      utils.RandomStrings(25),
	}

	pms := &LivePrograms{}
	psms := &LiveSubPrograms{}
	mbs := passport.NewMemberProvider()
	buildMsgFunc := func(en *libs.NoticeEvent) *libs.MsqMessage {
		msg := &libs.MsqMessage{}
		msg.Ts = utils.TimeMillisecond(time.Now())
		msg_txt := ""
		psm := psms.Get(int64(en.EventId))
		if psm == nil {
			return nil
		}
		pm := pms.Get(psm.ProgramId)
		if pm == nil {
			return nil
		}
		member := mbs.Get(en.FromId)
		if member == nil {
			return nil
		}
		if psm.ViewType == LIVE_SUBPROGRAM_VIEW_VS {
			msg_txt = fmt.Sprintf("%s %s VS %s 即将开始", pm.Title, psm.Vs1Name, psm.Vs2Name)
		}
		if psm.ViewType == LIVE_SUBPROGRAM_VIEW_SINGLE {
			msg_txt = fmt.Sprintf("%s %s 即将开始", pm.Title, psm.Title)
		}
		msg.DataType = "json"
		msg_data := libs.PushMsgData{
			ToUid:      en.FromId,
			ToUserName: member.UserName,
			ToNickName: member.NickName,
			DeviceType: member.DeviceType,
			PushId:     member.PushId,
			ChannelId:  member.PushChannelId,
			Title:      "节目单提醒",
			Content:    msg_txt,
			Category:   "live_program_notice",
			PushType:   libs.MSG_PUSH_TYPE_NOTICE,
		}

		//fmt.Println("--------------------------------------pgm_notice:", msg_data)

		body, _ := json.Marshal(msg_data)
		msg.Data = body
		msg.Ts = utils.TimeMillisecond(psm.StartTime)
		msg.MsgType = libs.MSG_TYPE_MPUSH
		return msg
	}
	noticer := libs.NewNoticer()
	err := noticer.StartNoticeTimer(int(subProgramId), event_program_key, dur,
		libs.NewNoticeStorager(program_notice_db, program_notice_collection), msq_config, buildMsgFunc)
	return err
}

func (n *ProgramNoticer) ResetNoticeTimer(subProgramId int64, t time.Time) error {
	lsp := &LiveSubPrograms{}
	sprm := lsp.Get(subProgramId)
	if sprm == nil {
		return errors.New("参数提供的节目预告不存在")
	}
	if t.Before(time.Now()) {
		return errors.New("提醒时间已过,添加后也无效果")
	}
	dur := t.Sub(time.Now())
	if dur.Minutes() < min_duration && dur.Minutes() > 0.0 {
		return errors.New("过于接近开始时间")
	}
	noticer := libs.NewNoticer()
	err := noticer.ResetNoticeTimer(int(subProgramId), event_program_key, dur)
	return err
}

func (n *ProgramNoticer) StopNoticeTimer(subProgramId int64) error {
	noticer := libs.NewNoticer()
	err := noticer.StopNoticeTimer(int(subProgramId), event_program_key)
	return err
}

func (n *ProgramNoticer) loadNoticeTimers() error {
	lps := &LivePrograms{}
	sps := &LiveSubPrograms{}
	programs := lps.getGteDate(time.Now())
	for _, program := range programs {
		sprograms := sps.Gets(program.Id)
		for _, sp := range sprograms {
			n.StartNoticeTimer(sp.Id)
		}
	}
	return nil
}
