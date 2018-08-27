package matchrace

import "github.com/astaxie/beego/orm"

func init() {
	orm.RegisterModel(
		new(RaceMode), new(MatchPlayer), new(MatchRecent),
		new(MatchGroup), new(MatchGroupPlayer), new(MatchVs),
		new(MatchEliminMs), new(MatchEliminVs),
	)
}
