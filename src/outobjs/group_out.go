package outobjs

import (
	"libs/groups"
	"time"
	"utils"
)

type OutGroupSetting struct {
	GroupNameLen    int    `json:"groupname_len"`
	GroupDescMaxLen int    `json:"groupdesc_maxlen"`
	GroupDescMinLen int    `json:"groupdesc_minlen"`
	DeductPoint     int64  `json:"deduct_point"`
	MinUsers        int    `json:"min_users"`
	LimitDay        int    `json:"limit_day"`
	GroupClause     string `json:"group_clause"`
}

type OutGroup struct {
	Id               int64               `json:"id"`
	Name             string              `json:"name"`
	Description      string              `json:"description"`
	OfUid            int64               `json:"ofuid"`
	OfMember         *OutSimpleMember    `json:"ofmember"`
	CreateTime       time.Time           `json:"create_time"`
	CreateFriendTime string              `json:"create_friendtime"`
	MemberCount      int                 `json:"member_count"`
	ThreadCount      int                 `json:"thread_count"`
	Country          string              `json:"country"`
	City             string              `json:"city"`
	Games            []*OutGame          `json:"games"`
	Status           groups.GROUP_STATUS `json:"status"`
	Type             groups.GROUP_TYPE   `json:"type"`
	Belong           groups.GROUP_BELONG `json:"belong"`
	ImgId            int64               `json:"img_id"`
	ImgUrl           string              `json:"img_url"`
	BgImgId          int64               `json:"bgimg_id"`
	BgImgUrl         string              `json:"bgimg_url"`
	Vitality         int                 `json:"vitality"`
	LongiTude        float32             `json:"longi_tude"`
	LatiTude         float32             `json:"lati_tude"`
	Recommend        bool                `json:"recommend"`
	StartTime        time.Time           `json:"start_time"`
	EndTime          time.Time           `json:"end_time"`
	RemainSeconds    int64               `json:"remain_seconds"`
	MinUsers         int                 `json:"min_users"`
	IsJoined         bool                `json:"is_joined"`
	DisplayOrder     int                 `json:"displayorder"`
	CreditNo         string              `json:"credit_no"`
}

type OutInviteMember struct {
	Uid     int64            `json:"uid"`
	Member  *OutSimpleMember `json:"member"`
	Invited bool             `json:"invited"`
	Joined  bool             `json:"joined"`
}

type OutMyGroups struct {
	MaxAllowGroupCount int         `json:"maxallow_group_count"`
	Groups             []*OutGroup `json:"groups"`
}

type OutGroupPagedList struct {
	CurrentPage int         `json:"current_page"`
	PageSize    int         `json:"page_size"`
	TotalPages  int         `json:"total_pages"`
	Total       int         `json:"total"`
	Groups      []*OutGroup `json:"groups"`
}

type OutGroupMember struct {
	Uid        int64            `json:"uid"`
	Member     *OutSimpleMember `json:"member"`
	JoinedTime time.Time        `json:"joined_time"`
}

type OutThreadPagedList struct {
	CurrentPage int          `json:"current_page"`
	PageSize    int          `json:"page_size"`
	Threads     []*OutThread `json:"threads"`
}

type OutThread struct {
	Id                 int64                `json:"id"`
	GroupId            int64                `json:"group_id"`
	Subject            string               `json:"subject"`
	AuthorId           int64                `json:"author_id"`
	Author             string               `json:"author"`
	AuthorMember       *OutSimpleMember     `json:"author_member"`
	CreateTime         time.Time            `json:"create_time"`
	CreateFriendTime   string               `json:"create_friendtime"`
	LastPostTime       time.Time            `json:"last_posttime"`
	LastPostFriendTime string               `json:"last_postfriendtime"`
	LastPostId         string               `json:"last_postid"`
	LastPostMember     *OutSimpleMember     `json:"last_postmember"`
	Views              int                  `json:"views"`
	Replies            int                  `json:"replies"`
	Shares             int                  `json:"shares"`
	Favs               int                  `json:"favorites"`
	Status             groups.THREAD_STATUS `json:"status"`
	ImgId              int64                `json:"img_id"`
	ImgUrl             string               `json:"img_url"`
	Closed             bool                 `json:"closed"`
	Highlight          bool                 `json:"highlight"`
	Heats              int                  `json:"heats"`
	DisplayOrder       int                  `json:"displayorder"`
}

type OutThreadPagedListForAdmin struct {
	CurrentPage int          `json:"current_page"`
	PageSize    int          `json:"page_size"`
	Total       int          `json:"total"`
	Pages       int          `json:"pages"`
	Threads     []*OutThread `json:"threads"`
}

type OutPostPagedList struct {
	CurrentPage int        `json:"current_page"`
	TotalPages  int        `json:"total_pages"`
	Total       int        `json:"total"`
	PageSize    int        `json:"page_size"`
	Posts       []*OutPost `json:"posts"`
	MaxDingPost *OutPost   `json:"max_ding_post"`
	MaxCaiPost  *OutPost   `json:"max_cai_post"`
	Thread      *OutThread `json:"thread"`
	JoinedGroup bool       `json:"joined_group"`
	LordPost    *OutPost   `json:"lord_post"`
}

type OutSinglePost struct {
	Thread      *OutThread `json:"thread"`
	JoinedGroup bool       `json:"joined_group"`
	Post        *OutPost   `json:"post"`
}

type OutPost struct {
	Id               string           `json:"id"`
	ThreadId         int64            `json:"thread_id"`
	AuthorId         int64            `json:"author_id"`
	AuthorMember     *OutSimpleMember `json:"author_member"`
	Subject          string           `json:"subject"`
	CreateTime       time.Time        `json:"create_time"`
	CreateFriendTime string           `json:"create_friendtime"`
	Message          string           `json:"msg"`
	Ip               string           `json:"ip"`
	Invisible        bool             `json:"invisible"`
	Ding             int              `json:"ding"`
	Cai              int              `json:"cai"`
	Position         int              `json:"position"`
	ReplyId          string           `json:"reply_id"`
	ReplyUid         int64            `json:"reply_uid"`
	ReplyMember      *OutSimpleMember `json:"reply_member"`
	ReplyPosition    int              `json:"reply_position"`
	Pics             []*OutPicture    `json:"imgs"`
	LongiTude        float32          `json:"longi_tude"`
	LatiTude         float32          `json:"lati_tude"`
	Dinged           bool             `json:"dinged"`
	Caied            bool             `json:"caied"`
}

type OutReport struct {
	Id      int64                  `json:"id"`
	RefId   string                 `json:"refid"`
	C       groups.REPORT_CATEGORY `json:"category"`
	Ts      time.Time              `json:"post_time"`
	RefTxt  string                 `json:"reftxt"`
	PostUid int64                  `json:"post_uid"`
	Member  *OutSimpleMember       `json:"post_member"`
	Msg     string                 `json:"msg"`
}

type OutReportPagedList struct {
	CurrentPage int          `json:"current_page"`
	TotalPages  int          `json:"total_pages"`
	PageSize    int          `json:"page_size"`
	Reports     []*OutReport `json:"reports"`
}

//type OutUserPostNote struct {
//	Id       int64                 `json:"id"`
//	Uid      int64                 `json:"uid"`
//	ThreadId int64                 `json:"thread_id"`
//	PostId   string                `json:"post_id"`
//	Subject  string                `json:"subject"`
//	DateLine time.Time             `json:"dateline"`
//	T        groups.USER_POST_TYPE `json:"t"`
//}

type OutUserPostNotePagedList struct {
	CurrentPage int          `json:"current_page"`
	Ts          int64        `json:"ts"`
	Threads     []*OutThread `json:"threads"`
}

type OutUserPostNoteCount struct {
	Publishs int       `json:"publishs"`
	Replys   int       `json:"replys"`
	LastTime time.Time `json:"lasttime"`
}

type OutGroupConfig struct {
	Id                           int64                `json:"id"`
	GroupNameLen                 int                  `json:"groupname_len"`
	GroupDescMaxLen              int                  `json:"groupdesc_maxlen"`
	GroupDescMinLen              int                  `json:"groupdesc_minlen"`
	CreateGroupBasePoint         int64                `json:"creategroup_basepoint"`
	CreateGroupRate              float32              `json:"creategroup_rate"`
	CreateGroupMinUsers          int                  `json:"creategroup_minusers"`
	CreateGroupRecruitDay        int                  `json:"creategroup_recruitday"`
	CreateGroupMaxCount          int                  `json:"creategroup_maxcount"`
	CreateGroupCertifiedMaxCount int                  `json:"creategroup_certifiedmaxcount"`
	CreateGroupClause            string               `json:"creategroup_clause"`
	ReportOptions                string               `json:"report_options"`
	NewThreadDefaultStatus       groups.THREAD_STATUS `json:"newthread_defstatus"`
}

func GetOutGroup(group *groups.Group, concernUid int64) *OutGroup {
	outgames := []*OutGame{}
	for _, gid := range group.GameIDs() {
		outgames = append(outgames, GetOutGameById(gid))
	}
	var remain_seconds int64 = 0
	if time.Now().Unix() < group.EndTime {
		remain_seconds = group.EndTime - time.Now().Unix()
	}
	isJoined := false
	if concernUid > 0 {
		gs := groups.NewGroupService(groups.GetDefaultCfg())
		isJoined = gs.IsJoined(concernUid, group.Id)
	}

	return &OutGroup{
		Id:               group.Id,
		Name:             group.Name,
		Description:      group.Description,
		OfUid:            group.Uid,
		OfMember:         GetOutSimpleMember(group.Uid),
		CreateTime:       time.Unix(group.CreateTime, 0),
		CreateFriendTime: utils.FriendTime(time.Unix(group.CreateTime, 0)),
		MemberCount:      group.Members,
		ThreadCount:      group.Threads,
		Country:          group.Country,
		City:             group.City,
		Games:            outgames,
		Status:           group.Status,
		Type:             group.Type,
		Belong:           group.Belong,
		ImgId:            group.Img,
		ImgUrl:           file.GetFileUrl(group.Img),
		BgImgId:          group.BgImg,
		BgImgUrl:         file.GetFileUrl(group.BgImg),
		Vitality:         group.Vitality,
		LongiTude:        group.LongiTude,
		LatiTude:         group.LatiTude,
		Recommend:        group.Recommend,
		StartTime:        time.Unix(group.StartTime, 0),
		EndTime:          time.Unix(group.EndTime, 0),
		RemainSeconds:    remain_seconds,
		MinUsers:         group.MinUsers,
		IsJoined:         isJoined,
		DisplayOrder:     group.DisplarOrder,
		CreditNo:         group.OrderNo,
	}
}

func GetOutGroups(gps []*groups.Group, concernUid int64) []*OutGroup {
	outgroups := make([]*OutGroup, len(gps), len(gps))
	for i, _g := range gps {
		outgroups[i] = GetOutGroup(_g, 0)
	}
	if concernUid > 0 {
		gs := groups.NewGroupService(groups.GetDefaultCfg())
		joineds := gs.MyAllJoinGroupIds(concernUid)
		for _, outg := range outgroups {
			if _, ok := joineds[outg.Id]; ok {
				outg.IsJoined = true
			}
		}
	}
	return outgroups
}

func GetOutThread(thread *groups.Thread) *OutThread {
	return &OutThread{
		Id:                 thread.Id,
		GroupId:            thread.GroupId,
		Subject:            thread.Subject,
		AuthorId:           thread.AuthorId,
		Author:             thread.Author,
		AuthorMember:       GetOutSimpleMember(thread.AuthorId),
		CreateTime:         time.Unix(thread.DateLine, 0),
		CreateFriendTime:   utils.FriendTime(time.Unix(thread.DateLine, 0)),
		LastPostTime:       time.Unix(thread.LastPost, 0),
		LastPostFriendTime: utils.FriendTime(time.Unix(thread.LastPost, 0)),
		LastPostId:         thread.LastId,
		LastPostMember:     GetOutSimpleMember(thread.LastPostUid),
		Views:              thread.Views,
		Replies:            thread.Replies,
		Shares:             thread.Shares,
		Favs:               thread.Favorites,
		Status:             thread.Status,
		ImgId:              thread.Img,
		ImgUrl:             file.GetFileUrl(thread.Img),
		Closed:             thread.Closed,
		Highlight:          thread.Highlight,
		Heats:              thread.Heats,
		DisplayOrder:       thread.DisplayOrder,
	}
}

func GetOutPost(post *groups.Post, res *groups.PostRes) *OutPost {
	var pics []*OutPicture
	if res != nil && len(res.ImgResource) > 0 {
		pics = []*OutPicture{}
		for _, img := range res.ImgResource {
			pics = append(pics, &OutPicture{
				Id:           img.OriginalImgId,
				ThumbnailPic: file.GetFileUrl(img.ThumbnailImgId),
				BmiddlePic:   file.GetFileUrl(img.BmiddleImgId),
				OriginalPic:  file.GetFileUrl(img.OriginalImgId),
			})
		}
	}
	message := post.Message
	if post.Invisible {
		message = ""
		pics = []*OutPicture{}
	}
	return &OutPost{
		Id:               post.Id,
		ThreadId:         post.ThreadId,
		AuthorId:         post.AuthorId,
		AuthorMember:     GetOutSimpleMember(post.AuthorId),
		Subject:          post.Subject,
		CreateTime:       time.Unix(post.DateLine, 0),
		CreateFriendTime: utils.FriendTime(time.Unix(post.DateLine, 0)),
		Message:          message,
		Ip:               post.Ip,
		Invisible:        post.Invisible,
		Ding:             post.Ding,
		Cai:              post.Cai,
		Position:         post.Position,
		ReplyId:          post.ReplyId,
		ReplyUid:         post.ReplyUid,
		ReplyMember:      GetOutSimpleMember(post.ReplyUid),
		ReplyPosition:    post.ReplyPosition,
		Pics:             pics,
		LongiTude:        post.LongiTude,
		LatiTude:         post.LatiTude,
	}
}
