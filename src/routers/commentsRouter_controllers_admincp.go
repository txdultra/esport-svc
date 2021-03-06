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

	beego.GlobalControllerRouter["controllers/admincp:CacheCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CacheCPController"],
		beego.ControllerComments{
			"CleanCache",
			`/clean`,
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

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"HomeAdAdd",
			`/homead/add`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"GetLastHomeAd",
			`/homead/get`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"DeleteHomeAd",
			`/homead/del`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"TeamAdd",
			`/team/add`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"TeamUpdate",
			`/team/update`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"DelTeam",
			`/team/remove`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"GetTeam",
			`/team/:id([0-9]+)`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"GetTeams",
			`/team/list`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"MatchModeAdd",
			`/matchrace/mode_add`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"MatchModeUpdate",
			`/matchrace/mode_update`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"GetMatchMode",
			`/matchrace/get_mode`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"GetMatchModes",
			`/matchrace/get_modes`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"CreateMatchRecent",
			`/matchrace/recent_add`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"UpdateMatchRecent",
			`/matchrace/recent_update`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"GetMatchRecents",
			`/matchrace/get_recents`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"CreateMatchGroup",
			`/matchrace/group_add`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"UpdateMatchGroup",
			`/matchrace/group_update`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"GetMatchGroups",
			`/matchrace/get_groups`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"CreateMatchGroupPlayers",
			`/matchrace/group_players_add`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"UpdateMatchGroupPlayer",
			`/matchrace/group_players_update`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"CreateMatchEliminMs",
			`/matchrace/elimin_ms_add`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"UpdateEliminMs",
			`/matchrace/elimin_ms_update`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"DeleteEliminMs",
			`/matchrace/elimin_ms_del`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"GetEliminMss",
			`/matchrace/get_eliminmss`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"CreateMatchEliminVs",
			`/matchrace/elimin_vs_add`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"UpdateMatchEliminVs",
			`/matchrace/elimin_vs_update`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"DeleteMatchEliminVs",
			`/matchrace/elimin_vs_del`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"CreateMatchVs",
			`/matchrace/matchvs_add`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"UpdateMatchVs",
			`/matchrace/matchvs_add`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"GetMatchVs",
			`/matchrace/matchvs`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"GetMatchVss",
			`/matchrace/matchvss`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"CreatePlayer",
			`/matchrace/player_add`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"UpdatePlayer",
			`/matchrace/player_update`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"GetMatchPlayer",
			`/matchrace/getplayer`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:CommonCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:CommonCPController"],
		beego.ControllerComments{
			"GetMatchPlayers",
			`/matchrace/getplayers`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:FeedbackCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:FeedbackCPController"],
		beego.ControllerComments{
			"List",
			`/list`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:GroupCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:GroupCPController"],
		beego.ControllerComments{
			"GetConfig",
			`/config/get`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:GroupCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:GroupCPController"],
		beego.ControllerComments{
			"UpdateConfig",
			`/config/update`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:GroupCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:GroupCPController"],
		beego.ControllerComments{
			"GroupSearch",
			`/group/search`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:GroupCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:GroupCPController"],
		beego.ControllerComments{
			"GetGroup",
			`/group/get`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:GroupCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:GroupCPController"],
		beego.ControllerComments{
			"CreateGroup",
			`/group/create`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:GroupCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:GroupCPController"],
		beego.ControllerComments{
			"CloseGroup",
			`/group/close`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:GroupCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:GroupCPController"],
		beego.ControllerComments{
			"UpdateGroup",
			`/group/update`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:GroupCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:GroupCPController"],
		beego.ControllerComments{
			"CloseThread",
			`/thread/close`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:GroupCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:GroupCPController"],
		beego.ControllerComments{
			"SetThreadOrder",
			`/thread/set_displayorder`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:GroupCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:GroupCPController"],
		beego.ControllerComments{
			"GetThread",
			`/thread/get`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:GroupCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:GroupCPController"],
		beego.ControllerComments{
			"GetThreads",
			`/thread/list`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:GroupCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:GroupCPController"],
		beego.ControllerComments{
			"CreateThread",
			`/thread/submit`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:GroupCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:GroupCPController"],
		beego.ControllerComments{
			"CreatePost",
			`/post/submit`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:GroupCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:GroupCPController"],
		beego.ControllerComments{
			"ClosePost",
			`/post/invisible`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:GroupCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:GroupCPController"],
		beego.ControllerComments{
			"GetPosts",
			`/post/list`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:GroupCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:GroupCPController"],
		beego.ControllerComments{
			"Reports",
			`/report/list`,
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
			"LiveProgramLast",
			`/org/program/last`,
			[]string{"get"},
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

	beego.GlobalControllerRouter["controllers/admincp:LiveCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:LiveCPController"],
		beego.ControllerComments{
			"CloseLockingSubProgram",
			`/org/locking_subprogram/close`,
			[]string{"post"},
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

	beego.GlobalControllerRouter["controllers/admincp:MemberCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:MemberCPController"],
		beego.ControllerComments{
			"ActionCredit",
			`/action_credit`,
			[]string{"post"},
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

	beego.GlobalControllerRouter["controllers/admincp:ShopCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:ShopCPController"],
		beego.ControllerComments{
			"OrderTransport",
			`/order_transport`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:ShopCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:ShopCPController"],
		beego.ControllerComments{
			"UpdateTransportNo",
			`/update_transport`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:ShopCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:ShopCPController"],
		beego.ControllerComments{
			"BuildTickets",
			`/item/build_tickets`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:ShopCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:ShopCPController"],
		beego.ControllerComments{
			"AddTag",
			`/tag/add`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:ShopCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:ShopCPController"],
		beego.ControllerComments{
			"UpdateTag",
			`/tag/update`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:ShopCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:ShopCPController"],
		beego.ControllerComments{
			"GetTag",
			`/tag/:id([0-9]+)`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:ShopCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:ShopCPController"],
		beego.ControllerComments{
			"GetTags",
			`/tag/all`,
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

	beego.GlobalControllerRouter["controllers/admincp:VodCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:VodCPController"],
		beego.ControllerComments{
			"VodPlayLists",
			`/playlists`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:VodCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:VodCPController"],
		beego.ControllerComments{
			"VodPlayListGet",
			`/playlist`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:VodCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:VodCPController"],
		beego.ControllerComments{
			"CreatePlaylist",
			`/playlist/create`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:VodCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:VodCPController"],
		beego.ControllerComments{
			"UpdatePlaylistVods",
			`/playlist/update_vods`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["controllers/admincp:VodCPController"] = append(beego.GlobalControllerRouter["controllers/admincp:VodCPController"],
		beego.ControllerComments{
			"PlaylistVods",
			`/playlist/vods`,
			[]string{"get"},
			nil})

}
