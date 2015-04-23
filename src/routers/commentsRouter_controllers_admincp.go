package routers

import (
	"github.com/astaxie/beego"
)

func init() {
	
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
			"LiveLockingSubProgramUpdate",
			`/org/locking_subprogram/update`,
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

	beego.GlobalControllerRouter["controllers/admincp:CacheCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CacheCPController"],
		beego.ControllerComments{
			"CleanCache",
			`/clean`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:FeedbackCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:FeedbackCPController"],
		beego.ControllerComments{
			"List",
			`/list`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:ShareCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:ShareCPController"],
		beego.ControllerComments{
			"Publishs",
			`/publishs`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:ShopCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:ShopCPController"],
		beego.ControllerComments{
			"GetItem",
			`/item`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:ShopCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:ShopCPController"],
		beego.ControllerComments{
			"GetItems",
			`/items`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:ShopCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:ShopCPController"],
		beego.ControllerComments{
			"AddItem",
			`/item/add`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:ShopCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:ShopCPController"],
		beego.ControllerComments{
			"UpdateItem",
			`/item/update`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:ShopCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:ShopCPController"],
		beego.ControllerComments{
			"DeleteItem",
			`/item/del`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:ShopCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:ShopCPController"],
		beego.ControllerComments{
			"AddItemStock",
			`/item/add_stock`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:ShopCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:ShopCPController"],
		beego.ControllerComments{
			"AddItemCodes",
			`/item/add_codes`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:ShopCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:ShopCPController"],
		beego.ControllerComments{
			"GetOrders",
			`/orders`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:ShopCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:ShopCPController"],
		beego.ControllerComments{
			"GetOrder",
			`/order`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:ShopCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:ShopCPController"],
		beego.ControllerComments{
			"SetOrderStatus",
			`/order_status`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:ShopCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:ShopCPController"],
		beego.ControllerComments{
			"Snap",
			`/snap`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:UTaskCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:UTaskCPController"],
		beego.ControllerComments{
			"AddGroup",
			`/add_group`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:UTaskCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:UTaskCPController"],
		beego.ControllerComments{
			"UpdateGroup",
			`/update_group`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:UTaskCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:UTaskCPController"],
		beego.ControllerComments{
			"DelGroup",
			`/del_group`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:UTaskCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:UTaskCPController"],
		beego.ControllerComments{
			"GetAllGroups",
			`/get_group_all`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:UTaskCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:UTaskCPController"],
		beego.ControllerComments{
			"GetGroup",
			`/get_group`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:UTaskCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:UTaskCPController"],
		beego.ControllerComments{
			"AddTask",
			`/add_task`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:UTaskCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:UTaskCPController"],
		beego.ControllerComments{
			"UpdateTask",
			`/update_task`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:UTaskCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:UTaskCPController"],
		beego.ControllerComments{
			"DelTask",
			`/del_task`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:UTaskCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:UTaskCPController"],
		beego.ControllerComments{
			"GetTasks",
			`/get_tasks`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:UTaskCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:UTaskCPController"],
		beego.ControllerComments{
			"GetTask",
			`/get_task`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:UTaskCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:UTaskCPController"],
		beego.ControllerComments{
			"GetEventNames",
			`/get_eventnames`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:UTaskCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:UTaskCPController"],
		beego.ControllerComments{
			"GetTaskTimers",
			`/get_tasktimers`,
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

	beego.GlobalControllerRouter["controllers/admincp:VodCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:VodCPController"],
		beego.ControllerComments{
			"UserVodCenterChangeUser",
			`/uc/change_user`,
			[]string{"post"},
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

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"Versions",
			`/versions`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"Version",
			`/version`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"VersionAdd",
			`/version/add`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"VersionUpdate",
			`/version/update`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"VersionDel",
			`/version/del`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:MemberCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:MemberCPController"],
		beego.ControllerComments{
			"GetRoles",
			`/get_roles`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:MemberCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:MemberCPController"],
		beego.ControllerComments{
			"GetRoleMembers",
			`/get_rolemembers`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:MemberCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:MemberCPController"],
		beego.ControllerComments{
			"SetRoles",
			`/set_roles`,
			[]string{"post"},
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
			"UpdateMember",
			`/update`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:MemberCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:MemberCPController"],
		beego.ControllerComments{
			"ResetNickName",
			`/reset_nickname`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:MemberCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:MemberCPController"],
		beego.ControllerComments{
			"GetMemberGames",
			`/games`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:MemberCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:MemberCPController"],
		beego.ControllerComments{
			"Search",
			`/search`,
			[]string{"get"},
			nil})

}
