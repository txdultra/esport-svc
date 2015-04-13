package outobjs

type OutFriendList struct {
	Total       int          `json:"total"`
	TotalPage   int          `json:"pages"`
	CurrentPage int          `json:"current_page"`
	Size        int          `json:"size"`
	Users       []*OutMember `json:"users"`
	Time        int64        `json:"t"`
}

type OutUserSTRelation struct {
	Source *OutUserRelation `json:"source"`
	Target *OutUserRelation `json:"target"`
}

type OutUserRelation struct {
	Uid        int64  `json:"uid"`
	NickName   string `json:"nick_name"`
	FollowedBy bool   `json:"followed_by"`
	Following  bool   `json:"following"`
}
