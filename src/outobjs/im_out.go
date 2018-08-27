package outobjs

type ChatUserCfg struct {
	NickName  string `json:"nick_name"`
	UserName  string `json:"user_name"`
	Pwd       string `json:"pwd"`
	Domain    string `json:"domain"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	IsVip     bool   `json:"is_vip"`
	IsOC      bool   `json:"is_oc"`
	FontColor int64  `json:"font_color"`
	AllowMsg  bool   `json:"allow_msg"`
	IsAdmin   bool   `json:"is_admin"`
}
