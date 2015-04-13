package libs

import "reflect"

const (
	APP_SERVER_VER          = "1.001"
	USR_NAME_REGEX          = `^[\w|u4e00-u9fa5]+$`
	EMAIL_REGEX             = "^[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?$"
	MOBILE_PHONE            = `^1[3458][0-9]{9}$`
	MEMBER_DEFAULT_PASSWORD = "ntv_Mobile_APP_V7!@#"
	PUSH_MSQ_QUEUE_NAME     = "push_msg_queue"
)

type CLIENT_OS string

const (
	CLIENT_OS_ANDROID CLIENT_OS = "android"
	CLIENT_OS_IOS     CLIENT_OS = "ios"
	CLIENT_OS_WP      CLIENT_OS = "wp"
)

type PIC_SIZE int

const (
	PIC_SIZE_ORIGINAL  PIC_SIZE = 1
	PIC_SIZE_THUMBNAIL PIC_SIZE = 2
	PIC_SIZE_MIDDLE    PIC_SIZE = 3
)

type PL struct {
	P     int
	S     int
	Total int
	Type  reflect.Type
	List  interface{}
}

//Msg后台处理进程接口
type IMsqMsgProcessTasker interface {
	Run() (<-chan string, error)
}

type IEventCounter interface {
	ResetEventCount(uid int64) bool
	NewEventCount(uid int64) int
}

type IEventInDecrCounter interface {
	IncrEventCount(uid int64) int
	DecrEventCount(uid int64) int
}

//type TestHandlerMsg struct{}

//func (m *TestHandlerMsg) Do(msg *MsqMessage) error {
//	defer func() {
//		if e := recover(); e != nil {
//			//log.Println("func error")
//			//return fmt.Errorf("func error")
//		}
//	}()
//	fmt.Println(string(msg.MsgType))
//	return nil
//}
