package libs

import (
	"fmt"
	//"bytes"
	"dbs"
	"encoding/json"
	"errors"
	//"fmt"
	"logs"
	"utils"
	//"github.com/astaxie/beego"
	"reflect"

	"github.com/astaxie/beego/orm"
	"github.com/streadway/amqp"
	//"strings"
	"sync"
	"time"
)

//注册队列处理handler
var msq_process_handlers = map[string]reflect.Type{}

func RegisterMsqProcessHandler(name string, handlerType reflect.Type) {
	msq_process_handlers[name] = handlerType
}

const (
	MSQ_RECEIVE_CLOSED    = "closed"
	MSQ_RECEIVE_SUCC      = "ok"
	MSQ_RECEIVE_COMPLETED = "completed"
	MSQ_RECEIVE_FAIL      = "fail"
)

type MsqMode string

const (
	MsqWorkQueueMode  MsqMode = "WORKQUEUE"
	MsqExcDirectMode  MsqMode = "direct"
	MsqExcTopicMode   MsqMode = "topic"
	MsqExcHeadersMode MsqMode = "headers"
	MsqExcFanoutMode  MsqMode = "fanout"
)

type MsqQueueConfig struct {
	MsqId      string //唯一
	MsqConnUrl string
	QueueName  string
	Durable    bool
	//Exclusive  bool
	AutoAck    bool
	Exchange   string
	RoutingKey string
	QueueMode  MsqMode
}

//消息类型
type MSG_TYPE string

const (
	MSG_TYPE_MPUSH   MSG_TYPE = "mobile_push"
	MSG_TYPE_DBBATCH MSG_TYPE = "db_batch"
)

type MsqMessage struct {
	DataType string   `json:"data_type"`
	Ts       int64    `json:"ts"`
	MsgType  MSG_TYPE `json:"msg_type"`
	Data     []byte   `json:"data"`
}

type IMsqService interface {
	MsqId() string
	Config() *MsqQueueConfig
	Init(config *MsqQueueConfig) error
	Send(msg *MsqMessage) error
	BeginMulti() error
	EndMulti() error
	MultiSend(msgs []*MsqMessage) error
	Receive(handler IMsqMsgProcesser) (<-chan string, error) //完成时发送ok
	Receiving(handler IMsqMsgProcesser) error                //接收进行时
	DelQueue() error
}

type IMsqMsgProcesser interface {
	Do(msg *MsqMessage) error
}

//rabbitMQ service
type AmpqMsqService struct {
	_config *MsqQueueConfig
	conn    *amqp.Connection
	channel *amqp.Channel
}

func (s *AmpqMsqService) MsqId() string {
	return s._config.MsqId
}

func (s *AmpqMsqService) Config() *MsqQueueConfig {
	return s._config
}

func (s *AmpqMsqService) Init(config *MsqQueueConfig) error {
	s._config = config
	_, _, err := s.createConn(config, nil)
	return err
}

func (s *AmpqMsqService) DelQueue() error {
	conn, c, err := s.createConn(s._config, nil)
	if err != nil {
		return errors.New("connection.open: " + err.Error())
	}
	defer conn.Close()
	defer c.Close()
	queue := s._config.QueueName
	if len(s._config.Exchange) > 0 {
		queue = s._config.Exchange
	}
	if _, err := c.QueueDelete(queue, false, true, false); err != nil {
		return errors.New("err purge:" + err.Error())
	}
	return nil
}

func (s *AmpqMsqService) createConn(config *MsqQueueConfig, args amqp.Table) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(config.MsqConnUrl)
	if err != nil {
		return nil, nil, errors.New("connection.open: " + err.Error())
	}
	c, err := conn.Channel()
	if err != nil {
		defer conn.Close()
		return nil, nil, errors.New("channel.open: " + err.Error())
	}
	switch config.QueueMode {
	case MsqExcDirectMode:
		if err := c.ExchangeDeclare(config.Exchange, string(config.QueueMode), config.Durable, false, false, false, args); err != nil {
			defer conn.Close()
			defer c.Close()
			return nil, nil, errors.New("exchange.declare destination: " + err.Error())
		}
	case MsqWorkQueueMode:
		if _, err := c.QueueDeclare(config.QueueName, config.Durable, false, false, false, args); err != nil {
			defer conn.Close()
			defer c.Close()
			return nil, nil, errors.New("queue.declare source: " + err.Error())
		}
	default:
		defer conn.Close()
		defer c.Close()
		return nil, nil, errors.New("not exist queue's mode")
	}
	return conn, c, nil
}

func (s *AmpqMsqService) EndMulti() error {
	if s.channel != nil {
		s.channel.Close()
	}
	if s.conn != nil {
		s.conn.Close()
	}
	return nil
}

func (s *AmpqMsqService) BeginMulti() error {
	_conn, _c, err := s.createConn(s._config, nil)
	if err != nil {
		return errors.New(err.Error())
	}
	s.conn = _conn
	s.channel = _c
	return nil
}

func (s *AmpqMsqService) MultiSend(msgs []*MsqMessage) error {
	if s.conn == nil || s.channel == nil {
		return fmt.Errorf("批量发送必须先调用BeginMulti方法")
	}
	deliveryMode := amqp.Persistent
	if !s._config.Durable {
		deliveryMode = amqp.Transient
	}
	for _, msg := range msgs {
		body, err := json.Marshal(msg)
		if err != nil {
			continue
		}
		_msg := amqp.Publishing{
			DeliveryMode: deliveryMode,
			Timestamp:    time.Now(),
			Body:         body,
		}
		rkey := s._config.QueueName
		if len(s._config.Exchange) > 0 {
			rkey = s._config.RoutingKey
		}
		s.channel.Publish(s._config.Exchange, rkey, false, false, _msg)
	}
	return nil
}

func (s *AmpqMsqService) Send(msg *MsqMessage) error {
	conn, c, err := s.createConn(s._config, nil)
	if err != nil {
		return errors.New(err.Error())
	}
	defer conn.Close()
	defer c.Close()
	deliveryMode := amqp.Persistent
	if !s._config.Durable {
		deliveryMode = amqp.Transient
	}
	body, err := json.Marshal(msg)
	if err != nil {
		return errors.New("转换消息数据错误")
	}
	_msg := amqp.Publishing{
		DeliveryMode: deliveryMode,
		Timestamp:    time.Now(),
		Body:         body,
	}
	rkey := s._config.QueueName
	if len(s._config.Exchange) > 0 {
		rkey = s._config.RoutingKey
	}
	err = c.Publish(s._config.Exchange, rkey, false, false, _msg)
	return err
}

func (s *AmpqMsqService) Receive(handler IMsqMsgProcesser) (<-chan string, error) {
	conn, c, err := s.createConn(s._config, nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	cc := make(chan string)
	queueName := s._config.QueueName
	rkey := ""
	if len(s._config.Exchange) > 0 {
		queueName = s._config.Exchange
		rkey = s._config.RoutingKey
	}
	deliver, err := c.Consume(queueName, rkey, s._config.AutoAck, false, false, false, nil)
	if err != nil {
		defer conn.Close()
		defer c.Close()
		return nil, errors.New("接收队列消息错误:" + err.Error())
	}
	go func(channel chan string, mqconn *amqp.Connection, mqc *amqp.Channel) {
		defer mqconn.Close()
		defer mqc.Close()
		for {
			msg, ok := <-deliver
			if !ok {
				channel <- MSQ_RECEIVE_CLOSED //已无消息或通道关闭
				return
			} else {
				channel <- MSQ_RECEIVE_SUCC
			}
			utils.Try(func() {
				var m MsqMessage
				err = json.Unmarshal(msg.Body, &m)
				if err == nil {
					err = handler.Do(&m)
					if err != nil { //打印错误,写入log
						logs.Errorf("msq fail:%s", err.Error())
					}
				}
				if s._config.AutoAck {
					msg.Ack(false)
				}
			}, func(e interface{}) {
				//记录错误,保证继续执行
			}, func() {
			})
		}
	}(cc, conn, c)
	return (<-chan string)(cc), nil
}

func (s *AmpqMsqService) Receiving(handler IMsqMsgProcesser) error {
	return errors.New("未实现此方法")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//Msq任务执行调度集
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var MSQ_USE_DRIVER string
var MSQ_DRIVERS = map[string]reflect.Type{
	"amqp": reflect.TypeOf(AmpqMsqService{}),
}
var MSQ_SERVICE = map[string]IMsqService{}
var MSQ_TASK = map[string]*time.Timer{}
var MSQ_MUTEX = new(sync.Mutex)
var MSQ_TIMER_RUNNING = map[string]bool{}
var MSQ_TIMER_MUTEX = new(sync.Mutex)
var MSQ_RELOAD_ONCE sync.Once

type Msq struct {
	MsqAdaptor string
}

func NewMsq() *Msq {
	msq := &Msq{}
	return msq
}

func InitMsqTask() {
	MSQ_RELOAD_ONCE.Do(func() {
		msq := &Msq{}
		msq.reloadTask()
	})
}

//开启进程或重置后使用
func (t *Msq) reloadTask() {
	msqdb := &MsqDb{}
	statuss := []string{string(MSQ_STATE_READY), string(MSQ_STATE_RUNNING)}
	tasks := msqdb.Gets(statuss)
	if tasks == nil {
		return
	}
	for _, tak := range tasks {
		tt := t.GetTask(tak.MsqId)
		if tt == nil { //处理已销毁的任务
			if tak.ScheduleTime.Before(time.Now().Add(-5 * time.Minute)) {
				tak.Status = MSQ_STATE_DISCARD
				tak.Des += time.Now().String() + ":重置任务时,任务指定时间已过;"
				msqdb.Update(tak)
				continue
			}
			data := []byte(tak.CfgJson)
			cfg := &MsqQueueConfig{}
			err := json.Unmarshal(data, cfg)
			if err != nil {
				tak.Des += time.Now().String() + ":cfg字符串转换失败;"
				msqdb.Update(tak)
				continue
			}
			_, err = t.CreateMsqService(cfg)
			if err != nil {
				tak.Des += time.Now().String() + ":启动服务失败;"
				msqdb.Update(tak)
				continue
			}
			_, ok := msq_process_handlers[tak.ConsumerType]
			if !ok {
				tak.Des += time.Now().String() + ":不存在指定的队列处理对象;"
				msqdb.Update(tak)
			}
			scheduleTime := tak.ScheduleTime
			if tak.ScheduleTime.Before(time.Now()) {
				scheduleTime = tak.ScheduleTime.Add(10 * time.Second)
			}
			t.AddTask(tak.MsqId, scheduleTime, uint8(tak.Consumers), tak.ConsumerType, tak.CompletedDel)
		}
	}
}

func (t *Msq) CreateMsqService(config *MsqQueueConfig) (IMsqService, error) {
	if len(config.MsqId) < 10 {
		return nil, errors.New("IdentityId不能少于10个字符,已保证唯一性")
	}
	msqId := config.MsqId
	MSQ_MUTEX.Lock()
	defer MSQ_MUTEX.Unlock()
	if s, ok := MSQ_SERVICE[msqId]; ok {
		return s, nil
	}
	_t, has := MSQ_DRIVERS[MSQ_USE_DRIVER]
	if !has {
		return nil, errors.New("msq adaptor not exist")
	}
	var service IMsqService = reflect.New(_t).Interface().(IMsqService)
	service.Init(config)
	MSQ_SERVICE[msqId] = service
	return service, nil
}

func (t *Msq) GetMsqService(msqId string) (IMsqService, error) {
	MSQ_MUTEX.Lock()
	defer MSQ_MUTEX.Unlock()
	if s, ok := MSQ_SERVICE[msqId]; ok {
		return s, nil
	}
	return nil, errors.New("msq not exist")
}

func (t *Msq) AddTask(msqId string, scheduleTime time.Time, consumers uint8, consumerHandlerType string, completedDelQueue bool) error {
	MSQ_MUTEX.Lock()
	defer MSQ_MUTEX.Unlock()
	if consumers == 0 {
		return errors.New("consumers can't equal zero")
	}
	if scheduleTime.Before(time.Now()) {
		return errors.New("schedule time can't less now")
	}
	if _, ok := MSQ_TASK[msqId]; ok {
		return errors.New("already exist msq task")
	}
	service, ok := MSQ_SERVICE[msqId]
	if !ok {
		return errors.New("first create msq service")
	}
	handler_type, ok := msq_process_handlers[consumerHandlerType]
	if !ok {
		return errors.New("consumerHandlerType not exist")
	}
	fun := func() {
		var wg sync.WaitGroup
		//记录任务运行时
		msqdb := &MsqDb{}
		m := msqdb.GetByMsqId(msqId)
		if m != nil {
			m.Status = MSQ_STATE_RUNNING
			msqdb.Update(m)
		}
		//正在执行
		MSQ_TIMER_MUTEX.Lock()
		MSQ_TIMER_RUNNING[msqId] = true
		MSQ_TIMER_MUTEX.Unlock()

		for i := 0; i < int(consumers); i++ {
			var handler IMsqMsgProcesser = reflect.New(handler_type).Interface().(IMsqMsgProcesser)
			wg.Add(1)
			go func() {
				c, err := service.Receive(handler)
				if err != nil { //处理错误，记录log
					logs.Errorf("msq add task fail:%s", err.Error())
					wg.Done()
				} else {
					select {
					case <-c:
						wg.Done()
					}
				}
			}()
		}
		wg.Wait()
		MSQ_TIMER_MUTEX.Lock()
		MSQ_TIMER_RUNNING[msqId] = false
		MSQ_TIMER_MUTEX.Unlock()
		//任务完成,删除队列
		if m != nil {
			m.Status = MSQ_STATE_COMPLETED
			msqdb.Update(m)
		}
		if completedDelQueue {
			service.DelQueue()
		}
	}
	dur := scheduleTime.Sub(time.Now())
	timer := time.AfterFunc(dur, fun)
	MSQ_TASK[msqId] = timer

	//记录任务
	msqdb := &MsqDb{}
	m := msqdb.GetByMsqId(msqId)
	if m == nil {
		bys, _ := json.Marshal(service.Config())
		m = &Msqtor{
			MsqId:        msqId,
			ScheduleTime: scheduleTime,
			CfgJson:      string(bys),
			ConsumerType: consumerHandlerType,
			Consumers:    int16(consumers),
			Status:       MSQ_STATE_READY,
			CreateTime:   time.Now(),
			CompletedDel: completedDelQueue,
		}
		msqdb.Create(m)
	}
	return nil
}

func (t *Msq) ResetTask(msqId string, scheduleTime time.Time, handler IMsqMsgProcesser) error {
	MSQ_MUTEX.Lock()
	defer MSQ_MUTEX.Unlock()
	if scheduleTime.Before(time.Now()) {
		return errors.New("schedule time can't less now")
	}
	task, ok := MSQ_TASK[msqId]
	if ok {
		return errors.New("already exist msq task")
	}
	dur := scheduleTime.Sub(time.Now())
	ok = task.Reset(dur)
	if !ok {
		return errors.New("reset fail")
	}
	//记录任务
	msqdb := &MsqDb{}
	m := msqdb.GetByMsqId(msqId)
	if m != nil {
		m.ScheduleTime = scheduleTime
		msqdb.Update(m)
	}
	return nil
}

func (t *Msq) DelTask(msqId string) error {
	MSQ_MUTEX.Lock()
	defer MSQ_MUTEX.Unlock()
	timer, ok := MSQ_TASK[msqId]
	if !ok {
		return errors.New("不存在指定任务")
	}

	MSQ_TIMER_MUTEX.Lock()
	defer MSQ_TIMER_MUTEX.Unlock()
	run, ok := MSQ_TIMER_RUNNING[msqId]
	if ok && run {
		return errors.New("任务正在执行")
	}

	timer.Stop()
	delete(MSQ_TASK, msqId)
	delete(MSQ_TIMER_RUNNING, msqId)

	//记录任务
	msqdb := &MsqDb{}
	msqdb.Delete(msqId)

	return nil
}

func (t *Msq) GetTask(msqId string) *time.Timer {
	timer, ok := MSQ_TASK[msqId]
	if !ok {
		return nil
	}
	return timer
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//Msq持久化
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type MsqDb struct{}

func (d *MsqDb) Create(msqtor *Msqtor) (int64, error) {
	_m := d.GetByMsqId(msqtor.MsqId)
	if _m != nil {
		return 0, errors.New("已存在MsqId")
	}
	o := dbs.NewDefaultOrm()
	id, err := o.Insert(msqtor)
	if err == nil {
		msqtor.Id = id
		return id, nil
	}
	return 0, err
}

func (d *MsqDb) Get(id int64) *Msqtor {
	msqtor := &Msqtor{}
	msqtor.Id = id
	o := dbs.NewDefaultOrm()
	err := o.Read(msqtor)
	if err == orm.ErrNoRows {
		return nil
	} else if err == orm.ErrMissPK {
		return nil
	}
	return msqtor
}

func (d *MsqDb) GetByMsqId(msqId string) *Msqtor {
	msqtor := &Msqtor{}
	o := dbs.NewDefaultOrm()
	qs := o.QueryTable(msqtor)
	err := qs.Filter("msqid", msqId).One(msqtor)
	if err == orm.ErrMultiRows {
		return nil
	}
	if err == orm.ErrNoRows {
		return nil
	}
	return msqtor
}

func (d *MsqDb) Update(msqtor *Msqtor) error {
	o := dbs.NewDefaultOrm()
	if num, err := o.Update(msqtor); err == nil {
		if num == 1 {
			return nil
		}
		return errors.New("MsqDb:" + msqtor.MsqId + "更新失败")
	}
	return errors.New("MsqDb:" + msqtor.MsqId + "更新失败")
}

func (d *MsqDb) Delete(msqId string) error {
	msqtor := &Msqtor{}
	o := dbs.NewDefaultOrm()
	qs := o.QueryTable(msqtor)
	_, err := qs.Filter("msqid", msqId).Delete()
	return err
}

func (d *MsqDb) Gets(statuss []string) []*Msqtor {
	msqtor := &Msqtor{}
	var msqtors []*Msqtor
	o := dbs.NewDefaultOrm()
	qs := o.QueryTable(msqtor)
	num, err := qs.Filter("status__in", statuss).All(&msqtors)
	if err != nil || num == 0 {
		return nil
	}
	return msqtors
}
