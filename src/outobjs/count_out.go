package outobjs

type OutMemberNewCount struct {
	NewFollowers     int `json:"followers"`
	NewSubscrs       int `json:"subscrs"`
	NewMsgs          int `json:"msgs"`
	ShareMsgs        int `json:"share_msgs"`
	VodMsgs          int `json:"vod_msgs"`
	NewLiveSubscrs   int `json:"program_subscrs"`
	NewShareNotices  int `json:"share_notices"`
	LastNewShareMsgs int `json:"share_new_msgs"`
	NewGroupMsgs     int `json:"group_msgs"`
}
