package routers

import (
	"github.com/astaxie/beego"
)

func init() {
	
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

	beego.GlobalControllerRouter["controllers:UserTaskController"] = append(beego.GlobalControllerRouter["controllers:UserTaskController"],
		beego.ControllerComments{
			"All",
			`/all`,
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
			"GetUrl",
			`/url/:id([0-9]+)`,
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
			"FlvsCallback",
			`/flvs_callback`,
			[]string{"post"},
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

	beego.GlobalControllerRouter["controllers:VideoController"] = append(beego.GlobalControllerRouter["controllers:VideoController"],
		beego.ControllerComments{
			"DownloadClarities",
			`/download/clarities`,
			[]string{"get"},
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

	beego.GlobalControllerRouter["controllers:ImageController"] = append(beego.GlobalControllerRouter["controllers:ImageController"],
		beego.ControllerComments{
			"Blur",
			`/blur`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:ShareController"] = append(beego.GlobalControllerRouter["controllers:ShareController"],
		beego.ControllerComments{
			"Publish",
			`/publish`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:ShareController"] = append(beego.GlobalControllerRouter["controllers:ShareController"],
		beego.ControllerComments{
			"Get",
			`/:id([0-9]+)`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:ShareController"] = append(beego.GlobalControllerRouter["controllers:ShareController"],
		beego.ControllerComments{
			"Delete",
			`/del`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:ShareController"] = append(beego.GlobalControllerRouter["controllers:ShareController"],
		beego.ControllerComments{
			"PublicTimeline",
			`/public_timeline`,
			[]string{"get"},
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

	beego.GlobalControllerRouter["controllers:ShareController"] = append(beego.GlobalControllerRouter["controllers:ShareController"],
		beego.ControllerComments{
			"ShareComments",
			`/comments`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:ShareController"] = append(beego.GlobalControllerRouter["controllers:ShareController"],
		beego.ControllerComments{
			"ShareComment",
			`/comment`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:ShareController"] = append(beego.GlobalControllerRouter["controllers:ShareController"],
		beego.ControllerComments{
			"DelShareComment",
			`/comment/del`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["controllers:ShareController"] = append(beego.GlobalControllerRouter["controllers:ShareController"],
		beego.ControllerComments{
			"LastNewMsg",
			`/last_new`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:ShareController"] = append(beego.GlobalControllerRouter["controllers:ShareController"],
		beego.ControllerComments{
			"DelNotices",
			`/notice_del`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["controllers:ShareController"] = append(beego.GlobalControllerRouter["controllers:ShareController"],
		beego.ControllerComments{
			"EmptyNotices",
			`/notice_empty`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["controllers:ShareController"] = append(beego.GlobalControllerRouter["controllers:ShareController"],
		beego.ControllerComments{
			"NoticeCount",
			`/notice_count`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:ShareController"] = append(beego.GlobalControllerRouter["controllers:ShareController"],
		beego.ControllerComments{
			"NewNotices",
			`/new_notices`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:ShareController"] = append(beego.GlobalControllerRouter["controllers:ShareController"],
		beego.ControllerComments{
			"Notices",
			`/notices`,
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

	beego.GlobalControllerRouter["controllers:CountController"] = append(beego.GlobalControllerRouter["controllers:CountController"],
		beego.ControllerComments{
			"MemberCounts",
			`/all`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:UcController"] = append(beego.GlobalControllerRouter["controllers:UcController"],
		beego.ControllerComments{
			"Create",
			`/create`,
			[]string{"post"},
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
			"SetMemberBackgroundImg",
			`/profile/set_bgimg`,
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
			"Logout",
			`/logout`,
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
			"MemberGameSingle",
			`/games_single`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:MemberController"] = append(beego.GlobalControllerRouter["controllers:MemberController"],
		beego.ControllerComments{
			"RemoveMemberGameSingle",
			`/games_single`,
			[]string{"delete"},
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
			"GetOriginalAvatar",
			`/get_original_avatar`,
			[]string{"get"},
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

	beego.GlobalControllerRouter["controllers:ShopController"] = append(beego.GlobalControllerRouter["controllers:ShopController"],
		beego.ControllerComments{
			"GetProvinces",
			`/provinces`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:ShopController"] = append(beego.GlobalControllerRouter["controllers:ShopController"],
		beego.ControllerComments{
			"GetAreas",
			`/areas`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:ShopController"] = append(beego.GlobalControllerRouter["controllers:ShopController"],
		beego.ControllerComments{
			"ShowItems",
			`/items_show`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:ShopController"] = append(beego.GlobalControllerRouter["controllers:ShopController"],
		beego.ControllerComments{
			"GetItem",
			`/item/:id([0-9]+)`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:ShopController"] = append(beego.GlobalControllerRouter["controllers:ShopController"],
		beego.ControllerComments{
			"GetOrders",
			`/orders`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:ShopController"] = append(beego.GlobalControllerRouter["controllers:ShopController"],
		beego.ControllerComments{
			"GetOrder",
			`/order`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:ShopController"] = append(beego.GlobalControllerRouter["controllers:ShopController"],
		beego.ControllerComments{
			"Stocks",
			`/stocks/:id([0-9]+)`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:ShopController"] = append(beego.GlobalControllerRouter["controllers:ShopController"],
		beego.ControllerComments{
			"Buy",
			`/buy`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:ShopController"] = append(beego.GlobalControllerRouter["controllers:ShopController"],
		beego.ControllerComments{
			"OrderCancel",
			`/order/cancel`,
			[]string{"post"},
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
			"BothFriends",
			`/both_friends`,
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

	beego.GlobalControllerRouter["controllers:OpenIDController"] = append(beego.GlobalControllerRouter["controllers:OpenIDController"],
		beego.ControllerComments{
			"Login",
			`/login`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:WebController"] = append(beego.GlobalControllerRouter["controllers:WebController"],
		beego.ControllerComments{
			"Home",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:WebController"] = append(beego.GlobalControllerRouter["controllers:WebController"],
		beego.ControllerComments{
			"DownApp",
			`/download`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:WebController"] = append(beego.GlobalControllerRouter["controllers:WebController"],
		beego.ControllerComments{
			"Down",
			`/down`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:WebController"] = append(beego.GlobalControllerRouter["controllers:WebController"],
		beego.ControllerComments{
			"Feedback",
			`/feedback`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers:WebController"] = append(beego.GlobalControllerRouter["controllers:WebController"],
		beego.ControllerComments{
			"PrivacyProtocol",
			`/privacy_protocol`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:WebController"] = append(beego.GlobalControllerRouter["controllers:WebController"],
		beego.ControllerComments{
			"VodPlay",
			`/vod/:id([0-9]+)`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:WebController"] = append(beego.GlobalControllerRouter["controllers:WebController"],
		beego.ControllerComments{
			"VodStream",
			`/vod/:id([0-9]+)/stream`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:WebController"] = append(beego.GlobalControllerRouter["controllers:WebController"],
		beego.ControllerComments{
			"PeronalLive",
			`/plive/:id([0-9]+)`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers:WebController"] = append(beego.GlobalControllerRouter["controllers:WebController"],
		beego.ControllerComments{
			"JigouLive",
			`/jlive/:id([0-9]+)`,
			[]string{"get"},
			nil})

}
