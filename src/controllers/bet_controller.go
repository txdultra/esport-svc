package controllers

// 菠菜 API
type BetController struct {
	BaseController
}

func (c *BetController) Prepare() {
	c.BaseController.Prepare()
}

func (c *BetController) URLMapping() {
}

// @Title 菠菜系统设定
// @Description 菠菜系统设定
// @Success 200 {object} outobjs.OutBetSetting
// @router /setting [get]
func (c *BetController) Setting() {

}

// @Title 菠菜游戏
// @Description 菠菜游戏(返回数组)
// @Success 200 {object} outobjs.OutBetGame
// @router /games [get]
func (c *BetController) Games() {

}

// @Title 菠菜赛事列表
// @Description 菠菜赛事列表
// @Param   access_token  path  string  false  "access_token"
// @Param   game_id   path  int  true  "游戏id"
// @Success 200 {object} outobjs.OutBetMatchList
// @router /match_list [get]
func (c *BetController) MatchList() {

}

// @Title 菠菜赛事对阵列表
// @Description 菠菜赛事对阵列表(返回数组)
// @Param   access_token  path  string  false  "access_token"
// @Param   match_id   path  int  true  "赛事id"
// @Success 200 {object} outobjs.OutBetCompletetion
// @router /match_completetions [get]
func (c *BetController) MatchCompletetions() {

}

// @Title 菠菜对阵
// @Description 菠菜对阵
// @Param   access_token  path  string  false  "access_token"
// @Success 200 {object} outobjs.OutBetCompletetion
// @router /completetion/:id([0-9]+) [get]
func (c *BetController) Getompletetions() {

}

// @Title 对阵下的菠菜模型列表
// @Description 对阵下的菠菜模型列表(返回数组)
// @Param   access_token  path  string  false  "access_token"
// @Param   id   path  int  true  "对阵id"
// @Success 200 {object} outobjs.OutBetModel
// @router /bocai_models [get]
func (c *BetController) BocaiModels() {

}

// @Title 我参加的菠菜对阵列表
// @Description 我参加的菠菜对阵列表(返回数组)
// @Param   access_token  path  string  true  "access_token"
// @Param   t   path  int  false  "类型(0=结束前,1=结束后)"
// @Param   page   path  int  false  "页数"
// @Success 200 {object} outobjs.OutBetCompletetionPagedList
// @router /my_bocai [get]
func (c *BetController) MyBocais() {

}

// @Title 我的菠菜统计
// @Description 我的菠菜统计
// @Param   access_token  path  string  true  "access_token"
// @Success 200 {object} outobjs.OutBetUserStats
// @router /my_stats [get]
func (c *BetController) MyStats() {

}

// @Title 菠菜说明
// @Description 菠菜说明
// @Success 200 {object} outobjs.OutBetExplain
// @router /explain [get]
func (c *BetController) Explain() {

}

// @Title 下注
// @Description 下注
// @Param   access_token  path  string  true  "access_token"
// @Param   obj_id  path  int  true  "下注的对象id"
// @Param   stakes  path  int  true  "下注金额(非倍数)"
// @Success 200 {object} libs.Error
// @router /bet [post]
func (c *BetController) Bet() {

}
