package pushapi

const (
	GETUI_API_URL = "http://sdk.open.api.igexin.com/apiex.htm"
)

type GetuiMessage struct {
	IsOffline         bool
	OfflineExpireTime int64
	PushNetWorkType   int
	Data              string
}

type GetuiTarget struct {
	AppId    string
	ClientId string
	Alias    string
}

type GetuiPusher struct {
	ApiKey       string
	SecretKey    string
	MasterSecret string
}

func NewGetuiPusher(apiKey string, secretKey string, masterSecret string) *GetuiPusher {
	return &GetuiPusher{
		apiKey,
		secretKey,
		masterSecret,
	}
}

func (gt *GetuiPusher) PushMsg(userId string, msg *GetuiMessage, target *GetuiTarget, requestId string) (string, error) {

	return "", nil
}
