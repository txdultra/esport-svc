package vars

const (
	MESSAGE_SYS_ID = 0
)

type MSG_TYPE string

const (
	MSG_TYPE_SYS MSG_TYPE = "system:msg"
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
