package outobjs

type OutMemberNewCount struct {
	NewFollowers     int `json:"followers"`
	NewSubscrs       int `json:"subscrs"`
	NewMsg           int `json:"msgs"`
	NewLiveSubscrs   int `json:"program_subscrs"`
	NewShareNotices  int `json:"share_notices"`
	LastNewShareMsgs int `json:"share_new_msgs"`
}
