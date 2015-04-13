package routers

import (
	"github.com/astaxie/beego"
)

func init() {
	
	beego.GlobalControllerRouter["controllers:ShareController"] = append(beego.GlobalControllerRouter["controllers:ShareController"],
		beego.ControllerComments{
			"Publish",
			`/publish`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:ShareController"] = append(beego.GlobalControllerRouter["controllers:ShareController"],
		beego.ControllerComments{
			"Delete",
			`/del`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:ShareController"] = append(beego.GlobalControllerRouter["controllers:ShareController"],
		beego.ControllerComments{
			"Timeline",
			`/timeline`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:ShareController"] = append(beego.GlobalControllerRouter["controllers:ShareController"],
		beego.ControllerComments{
			"SubscriptionCount",
			`/subscr/count`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:ShareController"] = append(beego.GlobalControllerRouter["controllers:ShareController"],
		beego.ControllerComments{
			"MySubscr",
			`/subscr/my`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:VodCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:VodCPController"],
		beego.ControllerComments{
			"Add",
			`/add`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:VodCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:VodCPController"],
		beego.ControllerComments{
			"Update",
			`/update`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:VodCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:VodCPController"],
		beego.ControllerComments{
			"Get",
			`/get`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:VodCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:VodCPController"],
		beego.ControllerComments{
			"List",
			`/list`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:VodCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:VodCPController"],
		beego.ControllerComments{
			"UpdateOfGameBatch",
			`/update/ofgame_batch`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:VodCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:VodCPController"],
		beego.ControllerComments{
			"UserVodCenterReptileAdd",
			`/uc/add`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:VodCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:VodCPController"],
		beego.ControllerComments{
			"UserVodCenterReptileDel",
			`/uc/del`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:VodCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:VodCPController"],
		beego.ControllerComments{
			"UserVodCenterReptileScanAll",
			`/uc/scan_all`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:VodCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:VodCPController"],
		beego.ControllerComments{
			"UserVodCenterReptileList",
			`/uc/list`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:LiveController"] = append(beego.GlobalControllerRouter["controllers:LiveController"],
		beego.ControllerComments{
			"PerGet",
			`/personal/:id([0-9]+)`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:LiveController"] = append(beego.GlobalControllerRouter["controllers:LiveController"],
		beego.ControllerComments{
			"PerGets",
			`/personal/list`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:LiveController"] = append(beego.GlobalControllerRouter["controllers:LiveController"],
		beego.ControllerComments{
			"ChannelGet",
			`/channel/:id([0-9]+)`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:LiveController"] = append(beego.GlobalControllerRouter["controllers:LiveController"],
		beego.ControllerComments{
			"ChannelStreams",
			`/channel/streams/:id([0-9]+)`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:LiveController"] = append(beego.GlobalControllerRouter["controllers:LiveController"],
		beego.ControllerComments{
			"SearchProgram",
			`/channel/programs/search`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:LiveController"] = append(beego.GlobalControllerRouter["controllers:LiveController"],
		beego.ControllerComments{
			"Program",
			`/channel/program`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:LiveController"] = append(beego.GlobalControllerRouter["controllers:LiveController"],
		beego.ControllerComments{
			"Programs",
			`/channel/programs`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:LiveController"] = append(beego.GlobalControllerRouter["controllers:LiveController"],
		beego.ControllerComments{
			"SubProgramsByChannel",
			`/channel/subprograms_by_channel`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:LiveController"] = append(beego.GlobalControllerRouter["controllers:LiveController"],
		beego.ControllerComments{
			"RemindSingle",
			`/channel/programs/remind_single`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:LiveController"] = append(beego.GlobalControllerRouter["controllers:LiveController"],
		beego.ControllerComments{
			"RemoveRemindSingle",
			`/channel/programs/remove_remind_single`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["controllers:LiveController"] = append(beego.GlobalControllerRouter["controllers:LiveController"],
		beego.ControllerComments{
			"RemoveReminds",
			`/channel/programs/remove_reminds`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["controllers:LiveController"] = append(beego.GlobalControllerRouter["controllers:LiveController"],
		beego.ControllerComments{
			"SubsReminded",
			`/channel/programs/subs_reminded`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:LiveController"] = append(beego.GlobalControllerRouter["controllers:LiveController"],
		beego.ControllerComments{
			"NewUUID",
			`/online/new_uuid`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:LiveController"] = append(beego.GlobalControllerRouter["controllers:LiveController"],
		beego.ControllerComments{
			"JoinLive",
			`/online/join`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:LiveController"] = append(beego.GlobalControllerRouter["controllers:LiveController"],
		beego.ControllerComments{
			"LeaveLive",
			`/online/leave`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:LiveController"] = append(beego.GlobalControllerRouter["controllers:LiveController"],
		beego.ControllerComments{
			"LiveStreamCallback",
			`/stream/callback`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:UcController"] = append(beego.GlobalControllerRouter["controllers:UcController"],
		beego.ControllerComments{
			"Create",
			`/create`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:LiveCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:LiveCPController"],
		beego.ControllerComments{
			"SupportPlatforms",
			`/live/rep_plats`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:LiveCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:LiveCPController"],
		beego.ControllerComments{
			"LivePersonalAdd",
			`/personal/add`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:LiveCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:LiveCPController"],
		beego.ControllerComments{
			"LivePersonalUpate",
			`/personal/update`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:LiveCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:LiveCPController"],
		beego.ControllerComments{
			"LivePersonalDel",
			`/personal/del`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:LiveCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:LiveCPController"],
		beego.ControllerComments{
			"LivePersonalList",
			`/personal/list`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:LiveCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:LiveCPController"],
		beego.ControllerComments{
			"LivePersonalGet",
			`/personal/:id([0-9]+)`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:LiveCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:LiveCPController"],
		beego.ControllerComments{
			"LiveChannelList",
			`/org/channel/list`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:LiveCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:LiveCPController"],
		beego.ControllerComments{
			"LiveChannelGet",
			`/org/channel/:id([0-9]+)`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:LiveCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:LiveCPController"],
		beego.ControllerComments{
			"LiveChannelAdd",
			`/org/channel/add`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:LiveCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:LiveCPController"],
		beego.ControllerComments{
			"LiveChannelUpdate",
			`/org/channel/update`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:LiveCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:LiveCPController"],
		beego.ControllerComments{
			"LiveChannelStreamList",
			`/org/channel_stream/list`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:LiveCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:LiveCPController"],
		beego.ControllerComments{
			"LiveChannelStreamAdd",
			`/org/channel_stream/add`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:LiveCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:LiveCPController"],
		beego.ControllerComments{
			"LiveChannelStreamUpdate",
			`/org/channel_stream/update`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:LiveCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:LiveCPController"],
		beego.ControllerComments{
			"LiveChannelStreamDel",
			`/org/channel_stream/del`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:LiveCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:LiveCPController"],
		beego.ControllerComments{
			"LiveProgramList",
			`/org/program/list`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:LiveCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:LiveCPController"],
		beego.ControllerComments{
			"LiveProgramGet",
			`/org/program/:id([0-9]+)`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:LiveCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:LiveCPController"],
		beego.ControllerComments{
			"LiveProgramAdd",
			`/org/program/add`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:LiveCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:LiveCPController"],
		beego.ControllerComments{
			"LiveProgramUpdate",
			`/org/program/update`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:LiveCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:LiveCPController"],
		beego.ControllerComments{
			"LiveProgramDelete",
			`/org/program/del`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:LiveCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:LiveCPController"],
		beego.ControllerComments{
			"LiveSubProgramAdd",
			`/org/subprogram/add`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:LiveCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:LiveCPController"],
		beego.ControllerComments{
			"LiveSubProgramUpdate",
			`/org/subprogram/update`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:LiveCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:LiveCPController"],
		beego.ControllerComments{
			"LiveSubProgramDel",
			`/org/subprogram/del`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:LiveCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:LiveCPController"],
		beego.ControllerComments{
			"LiveSubProgramList",
			`/org/subprogram/list`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:LiveCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:LiveCPController"],
		beego.ControllerComments{
			"LiveSubProgramGet",
			`/org/subprogram/:id([0-9]+)`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:CountController"] = append(beego.GlobalControllerRouter["controllers:CountController"],
		beego.ControllerComments{
			"MemberCounts",
			`/all`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:FeedbackController"] = append(beego.GlobalControllerRouter["controllers:FeedbackController"],
		beego.ControllerComments{
			"Submit",
			`/submit`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:ImageController"] = append(beego.GlobalControllerRouter["controllers:ImageController"],
		beego.ControllerComments{
			"Resize",
			`/resize`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:ImageController"] = append(beego.GlobalControllerRouter["controllers:ImageController"],
		beego.ControllerComments{
			"Crop",
			`/crop`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:MessageController"] = append(beego.GlobalControllerRouter["controllers:MessageController"],
		beego.ControllerComments{
			"Mentions",
			`/mentions`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:MessageController"] = append(beego.GlobalControllerRouter["controllers:MessageController"],
		beego.ControllerComments{
			"EmptyMentions",
			`/mentions/empty`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:MessageController"] = append(beego.GlobalControllerRouter["controllers:MessageController"],
		beego.ControllerComments{
			"DelMention",
			`/mentions/del`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:MessageController"] = append(beego.GlobalControllerRouter["controllers:MessageController"],
		beego.ControllerComments{
			"Count",
			`/count`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"Games",
			`/game/list`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"GameAdd",
			`/game/add`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"GameUpdate",
			`/game/update`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"Matchs",
			`/match/list`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"MatchAdd",
			`/match/add`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"MatchUpdate",
			`/match/update`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"RecommendList",
			`/recommend/list`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"RecommendAdd",
			`/recommend/add`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"RecommendUpdate",
			`/recommend/update`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"RecommendDelete",
			`/recommend/del`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"RecommendGet",
			`/recommend/:id([0-9]+)`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:CommentController"] = append(beego.GlobalControllerRouter["controllers:CommentController"],
		beego.ControllerComments{
			"Publish",
			`/publish`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:CommentController"] = append(beego.GlobalControllerRouter["controllers:CommentController"],
		beego.ControllerComments{
			"Gets",
			`/list`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:CommentController"] = append(beego.GlobalControllerRouter["controllers:CommentController"],
		beego.ControllerComments{
			"Get",
			`/get`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:FriendShipsController"] = append(beego.GlobalControllerRouter["controllers:FriendShipsController"],
		beego.ControllerComments{
			"Friends",
			`/friends/all`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:FriendShipsController"] = append(beego.GlobalControllerRouter["controllers:FriendShipsController"],
		beego.ControllerComments{
			"FriendsP",
			`/friends`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:FriendShipsController"] = append(beego.GlobalControllerRouter["controllers:FriendShipsController"],
		beego.ControllerComments{
			"Followers",
			`/followers`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:FriendShipsController"] = append(beego.GlobalControllerRouter["controllers:FriendShipsController"],
		beego.ControllerComments{
			"Show",
			`/show`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:FriendShipsController"] = append(beego.GlobalControllerRouter["controllers:FriendShipsController"],
		beego.ControllerComments{
			"Create",
			`/create`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:FriendShipsController"] = append(beego.GlobalControllerRouter["controllers:FriendShipsController"],
		beego.ControllerComments{
			"Destroy",
			`/destroy`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:FriendShipsController"] = append(beego.GlobalControllerRouter["controllers:FriendShipsController"],
		beego.ControllerComments{
			"Recmds",
			`/recmds`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:FileController"] = append(beego.GlobalControllerRouter["controllers:FileController"],
		beego.ControllerComments{
			"Get",
			`/:id([0-9]+)`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:FileController"] = append(beego.GlobalControllerRouter["controllers:FileController"],
		beego.ControllerComments{
			"Delete",
			`/:id([0-9]+)`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["controllers:FileController"] = append(beego.GlobalControllerRouter["controllers:FileController"],
		beego.ControllerComments{
			"Upload",
			`/upload`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:CollectController"] = append(beego.GlobalControllerRouter["controllers:CollectController"],
		beego.ControllerComments{
			"Add",
			`/add`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:CollectController"] = append(beego.GlobalControllerRouter["controllers:CollectController"],
		beego.ControllerComments{
			"Remove",
			`/remove`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:CollectController"] = append(beego.GlobalControllerRouter["controllers:CollectController"],
		beego.ControllerComments{
			"Removes",
			`/removes`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:CollectController"] = append(beego.GlobalControllerRouter["controllers:CollectController"],
		beego.ControllerComments{
			"Show",
			`/show`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:CollectController"] = append(beego.GlobalControllerRouter["controllers:CollectController"],
		beego.ControllerComments{
			"Gets",
			`/list`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:MemberController"] = append(beego.GlobalControllerRouter["controllers:MemberController"],
		beego.ControllerComments{
			"Register",
			`/register`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:MemberController"] = append(beego.GlobalControllerRouter["controllers:MemberController"],
		beego.ControllerComments{
			"Get",
			`/show`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:MemberController"] = append(beego.GlobalControllerRouter["controllers:MemberController"],
		beego.ControllerComments{
			"Update",
			`/profile/update`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:MemberController"] = append(beego.GlobalControllerRouter["controllers:MemberController"],
		beego.ControllerComments{
			"ShowFoundPwdMobile",
			`/security/found_pwd_mobile`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:MemberController"] = append(beego.GlobalControllerRouter["controllers:MemberController"],
		beego.ControllerComments{
			"SetFoundPwdMobile",
			`/security/found_pwd_mobile`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:MemberController"] = append(beego.GlobalControllerRouter["controllers:MemberController"],
		beego.ControllerComments{
			"Login",
			`/login`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:MemberController"] = append(beego.GlobalControllerRouter["controllers:MemberController"],
		beego.ControllerComments{
			"GetAccessToken",
			`/get_token`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:MemberController"] = append(beego.GlobalControllerRouter["controllers:MemberController"],
		beego.ControllerComments{
			"SetPassword",
			`/set_password`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:MemberController"] = append(beego.GlobalControllerRouter["controllers:MemberController"],
		beego.ControllerComments{
			"SetMail",
			`/set_email`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:MemberController"] = append(beego.GlobalControllerRouter["controllers:MemberController"],
		beego.ControllerComments{
			"MemberGames",
			`/games`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:MemberController"] = append(beego.GlobalControllerRouter["controllers:MemberController"],
		beego.ControllerComments{
			"GetMemberGames",
			`/games`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:MemberController"] = append(beego.GlobalControllerRouter["controllers:MemberController"],
		beego.ControllerComments{
			"Search",
			`/search`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:MemberController"] = append(beego.GlobalControllerRouter["controllers:MemberController"],
		beego.ControllerComments{
			"SetPushId",
			`/set_pushid`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:MemberController"] = append(beego.GlobalControllerRouter["controllers:MemberController"],
		beego.ControllerComments{
			"SetNickname",
			`/set_nickname`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:MemberController"] = append(beego.GlobalControllerRouter["controllers:MemberController"],
		beego.ControllerComments{
			"SetAvatar",
			`/set_avatar`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:MemberController"] = append(beego.GlobalControllerRouter["controllers:MemberController"],
		beego.ControllerComments{
			"FrequentAts",
			`/frequent_ats`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:MemberController"] = append(beego.GlobalControllerRouter["controllers:MemberController"],
		beego.ControllerComments{
			"GetMyConfig",
			`/my_config`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:MemberController"] = append(beego.GlobalControllerRouter["controllers:MemberController"],
		beego.ControllerComments{
			"SetMyConfig",
			`/my_config`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:VideoController"] = append(beego.GlobalControllerRouter["controllers:VideoController"],
		beego.ControllerComments{
			"Modes",
			`/modes`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:VideoController"] = append(beego.GlobalControllerRouter["controllers:VideoController"],
		beego.ControllerComments{
			"ListByGames",
			`/list/by_games`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:VideoController"] = append(beego.GlobalControllerRouter["controllers:VideoController"],
		beego.ControllerComments{
			"List",
			`/list`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:VideoController"] = append(beego.GlobalControllerRouter["controllers:VideoController"],
		beego.ControllerComments{
			"Get",
			`/:id([0-9]+)`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:VideoController"] = append(beego.GlobalControllerRouter["controllers:VideoController"],
		beego.ControllerComments{
			"Flvs",
			`/flvs`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:VideoController"] = append(beego.GlobalControllerRouter["controllers:VideoController"],
		beego.ControllerComments{
			"Recommends",
			`/recommends`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:VideoController"] = append(beego.GlobalControllerRouter["controllers:VideoController"],
		beego.ControllerComments{
			"Download",
			`/download`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:MemberCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:MemberCPController"],
		beego.ControllerComments{
			"VerifyNickName",
			`/verify_nickname`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:MemberCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:MemberCPController"],
		beego.ControllerComments{
			"VerifyUserName",
			`/verify_username`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:MemberCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:MemberCPController"],
		beego.ControllerComments{
			"AddMember",
			`/add`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:MemberCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:MemberCPController"],
		beego.ControllerComments{
			"SetMemberCertifiable",
			`/set_certified`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:MemberCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:MemberCPController"],
		beego.ControllerComments{
			"Search",
			`/search`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:AuthCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:AuthCPController"],
		beego.ControllerComments{
			"GetManagerAccessToken",
			`/get_access_token`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:AuthCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:AuthCPController"],
		beego.ControllerComments{
			"AccessTokenStatus",
			`/access_token/status`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:OpenIDController"] = append(beego.GlobalControllerRouter["controllers:OpenIDController"],
		beego.ControllerComments{
			"Login",
			`/login`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:CommonController"] = append(beego.GlobalControllerRouter["controllers:CommonController"],
		beego.ControllerComments{
			"VerifyPic",
			`/verify_pic`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:CommonController"] = append(beego.GlobalControllerRouter["controllers:CommonController"],
		beego.ControllerComments{
			"Version",
			`/version`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:CommonController"] = append(beego.GlobalControllerRouter["controllers:CommonController"],
		beego.ControllerComments{
			"Games",
			`/games`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:CommonController"] = append(beego.GlobalControllerRouter["controllers:CommonController"],
		beego.ControllerComments{
			"Match",
			`/match`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:CommonController"] = append(beego.GlobalControllerRouter["controllers:CommonController"],
		beego.ControllerComments{
			"Matchs",
			`/matchs`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:CommonController"] = append(beego.GlobalControllerRouter["controllers:CommonController"],
		beego.ControllerComments{
			"Expressions",
			`/expressions`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CacheCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CacheCPController"],
		beego.ControllerComments{
			"CleanCache",
			`/clean`,
			[]string{"get"},
			nil})

}
