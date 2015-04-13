package passport

type ThridUserInfo struct {
	UserName  string
	NickName  string
	Uid       string
	FigureUrl string
}

type IThirdAuth interface {
	GetUserInfo(token, openid string) (*ThridUserInfo, error)
}
