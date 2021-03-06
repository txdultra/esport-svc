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
	CLIENT_OS_OTHER   CLIENT_OS = "other"
)

type CLIENT_SRC int8

const (
	CLIENT_SRC_APP     CLIENT_SRC = 0
	CLIENT_SRC_WEB     CLIENT_SRC = 1
	CLIENT_SRC_MANAGER CLIENT_SRC = 2
)

type PIC_SIZE int

const (
	PIC_SIZE_ORIGINAL  PIC_SIZE = 1
	PIC_SIZE_THUMBNAIL PIC_SIZE = 2
	PIC_SIZE_MIDDLE    PIC_SIZE = 3
)

type CURRENCY_TYPE int

const (
	CURRENCY_TYPE_CREDIT CURRENCY_TYPE = 1
	CURRENCY_TYPE_RMB    CURRENCY_TYPE = 2
	CURRENCY_TYPE_JING   CURRENCY_TYPE = 4
)

func GetCurrencyName(cur CURRENCY_TYPE) string {
	if cur|CURRENCY_TYPE_CREDIT == CURRENCY_TYPE_CREDIT {
		return "积分"
	}
	if cur|CURRENCY_TYPE_RMB == CURRENCY_TYPE_RMB {
		return "人民币"
	}
	if cur|CURRENCY_TYPE_JING == CURRENCY_TYPE_JING {
		return "竞币"
	}
	return "未知"
}
