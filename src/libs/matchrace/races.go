package matchrace

import (
	"dbs"
	"fmt"
	"time"
	"utils"

	"github.com/astaxie/beego/orm"
)

type MatchRaceService struct{}

func (m *MatchRaceService) getModeCacheKey(id int64) string {
	return fmt.Sprintf("mobile_mr_mode:%d", id)
}

func (m *MatchRaceService) getModesCacheKey(matchId int64) string {
	return fmt.Sprintf("mobile_mr_modes_bymatch:%d", matchId)
}

func (m *MatchRaceService) GetMode(id int64) *RaceMode {
	cache := utils.GetCache()
	race := RaceMode{}
	err := cache.Get(m.getModeCacheKey(id), &race)
	if err == nil {
		return &race
	}
	o := dbs.NewDefaultOrm()
	race.Id = id
	err = o.Read(&race)
	if err != nil {
		return nil
	}
	cache.Set(m.getModeCacheKey(id), race, 72*time.Hour)
	return &race
}

func (m *MatchRaceService) CreateMode(race *RaceMode) error {
	if race.MatchId <= 0 {
		return fmt.Errorf("未设定赛事id")
	}
	if len(race.Title) == 0 {
		return fmt.Errorf("未设置标题")
	}
	o := dbs.NewDefaultOrm()
	id, err := o.Insert(race)
	if err != nil {
		return err
	}
	race.Id = id
	cache := utils.GetCache()
	cache.Set(m.getModeCacheKey(id), *race, 72*time.Hour)
	cache.Delete(m.getModesCacheKey(race.MatchId))
	return nil
}

func (m *MatchRaceService) UpdateMode(race *RaceMode) error {
	if race.MatchId <= 0 {
		return fmt.Errorf("未设定赛事id")
	}
	if len(race.Title) == 0 {
		return fmt.Errorf("未设置标题")
	}
	if race.Id <= 0 {
		return fmt.Errorf("对象id错误")
	}
	o := dbs.NewDefaultOrm()
	_, err := o.Update(race)
	cache := utils.GetCache()
	cache.Delete(m.getModeCacheKey(race.Id))
	return err
}

func (m *MatchRaceService) GetModes(matchId int64) []*RaceMode {
	cache := utils.GetCache()
	res := []int64{}
	err := cache.Get(m.getModesCacheKey(matchId), &res)
	if err != nil {
		o := dbs.NewDefaultOrm()
		var races []*RaceMode
		o.QueryTable(&RaceMode{}).Filter("matchid", matchId).OrderBy("-displayorder").All(&races, "id")
		for _, race := range races {
			res = append(res, race.Id)
		}
		cache.Set(m.getModesCacheKey(matchId), res, 72*time.Hour)
	}
	rms := []*RaceMode{}
	for _, id := range res {
		race := m.GetMode(id)
		if race != nil {
			rms = append(rms, race)
		}
	}
	return rms
}

func (m *MatchRaceService) getPlayerCacheKey(id int64) string {
	return fmt.Sprintf("mobile_mr_player:%d", id)
}

func (m *MatchRaceService) GetPlayer(id int64) *MatchPlayer {
	cache := utils.GetCache()
	mp := MatchPlayer{}
	err := cache.Get(m.getPlayerCacheKey(id), &mp)
	if err == nil {
		return &mp
	}
	o := dbs.NewDefaultOrm()
	mp.Id = id
	err = o.Read(&mp)
	if err != nil {
		return nil
	}
	cache.Set(m.getPlayerCacheKey(id), mp, 72*time.Hour)
	return &mp
}

func (m *MatchRaceService) CreatePlayer(player *MatchPlayer) error {
	if len(player.Name) == 0 {
		return fmt.Errorf("未设置名字")
	}
	o := dbs.NewDefaultOrm()
	id, err := o.Insert(player)
	if err != nil {
		return err
	}
	player.Id = id
	cache := utils.GetCache()
	cache.Set(m.getPlayerCacheKey(id), *player, 72*time.Hour)
	return nil
}

func (m *MatchRaceService) UpdatePlayer(player *MatchPlayer) error {
	if len(player.Name) == 0 {
		return fmt.Errorf("未设置名字")
	}
	if player.Id <= 0 {
		return fmt.Errorf("对象id错误")
	}
	o := dbs.NewDefaultOrm()
	_, err := o.Update(player)
	return err
}

func (m *MatchRaceService) GetPlayers(like string, page int, size int) (int, []*MatchPlayer) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	like = utils.StripSQLInjection(like)
	offset := (page - 1) * size
	o := dbs.NewDefaultOrm()
	objs := []*MatchPlayer{}
	query := o.QueryTable(&MatchPlayer{})
	if len(like) > 0 {
		query = query.Filter("pname__icontains", like)
	}
	total, _ := query.Count()
	query.OrderBy("-id").Offset(offset).Limit(size).All(&objs, "id")
	mps := []*MatchPlayer{}
	for _, obj := range objs {
		mp := m.GetPlayer(obj.Id)
		if mp != nil {
			mps = append(mps, mp)
		}
	}
	return int(total), mps
}

func (m *MatchRaceService) getVssCacheKey(player int64, matchId int64) string {
	return fmt.Sprintf("mobile_mr_vss_p_%d_m_%d", player, matchId)
}

func (m *MatchRaceService) getVsCacheKey(id int64) string {
	return fmt.Sprintf("mobile_mr_vs:%d", id)
}

func (m *MatchRaceService) GetMatchVss(player int64, matchId int64) []*MatchVs {
	cache := utils.GetCache()
	vss := []*MatchVs{}
	err := cache.Get(m.getVssCacheKey(player, matchId), &vss)
	if err == nil {
		return vss
	}
	o := dbs.NewDefaultOrm()
	query := o.QueryTable(&MatchVs{})
	cond := orm.NewCondition()
	cond.And("a", player)
	cond.Or("b", player)
	c2 := cond.AndCond(cond.And("matchid", matchId))
	query.SetCond(c2)
	query.All(&vss)
	cache.Set(m.getVssCacheKey(player, matchId), vss, 24*time.Hour)
	return vss
}

func (m *MatchRaceService) GetMatchVs(id int64) *MatchVs {
	ckey := m.getVsCacheKey(id)
	vs := &MatchVs{}
	cache := utils.GetCache()
	err := cache.Get(ckey, vs)
	if err == nil {
		return vs
	}
	o := dbs.NewDefaultOrm()
	vs.Id = id
	err = o.Read(vs)
	if err != nil {
		return nil
	}
	cache.Set(ckey, *vs, 72*time.Hour)
	return vs
}

func (m *MatchRaceService) CreateMatchVs(vs *MatchVs) error {
	if len(vs.AName) == 0 {
		return fmt.Errorf("A名称不能为空")
	}
	if len(vs.BName) == 0 {
		return fmt.Errorf("B名称不能为空")
	}
	if vs.MatchId <= 0 {
		return fmt.Errorf("关联赛事id未设置")
	}
	if vs.ModeId <= 0 {
		return fmt.Errorf("关联模型id未设置")
	}
	o := dbs.NewDefaultOrm()
	id, err := o.Insert(vs)
	if err != nil {
		return err
	}
	vs.Id = id
	cache := utils.GetCache()
	cache.Delete(m.getVssCacheKey(vs.A, vs.MatchId))
	cache.Delete(m.getVssCacheKey(vs.B, vs.MatchId))
	return nil
}

func (m *MatchRaceService) UpdateMatchVs(vs *MatchVs) error {
	if len(vs.AName) == 0 {
		return fmt.Errorf("A名称不能为空")
	}
	if len(vs.BName) == 0 {
		return fmt.Errorf("B名称不能为空")
	}
	if vs.MatchId <= 0 {
		return fmt.Errorf("关联赛事id未设置")
	}
	if vs.ModeId <= 0 {
		return fmt.Errorf("关联模型id未设置")
	}
	o := dbs.NewDefaultOrm()
	_, err := o.Update(vs)
	cache := utils.GetCache()
	cache.Delete(m.getVssCacheKey(vs.A, vs.MatchId))
	cache.Delete(m.getVssCacheKey(vs.B, vs.MatchId))
	cache.Delete(m.getVsCacheKey(vs.Id))
	return err
}

func (m *MatchRaceService) getRecentsCacheKey(modeId int64) string {
	return fmt.Sprintf("mobile_mr_recents_modeid:%d", modeId)
}

func (m *MatchRaceService) CreateMatchRecent(recent *MatchRecent) error {
	if recent.ModeId <= 0 {
		return fmt.Errorf("关联模型id未设置")
	}
	if recent.Player <= 0 {
		return fmt.Errorf("选手未设置")
	}
	o := dbs.NewDefaultOrm()
	has := o.QueryTable(&MatchRecent{}).Filter("modeid", recent.ModeId).Filter("player", recent.Player).Exist()
	if has {
		return fmt.Errorf("选手已存在模型中")
	}
	id, err := o.Insert(recent)
	if err != nil {
		return err
	}
	recent.Id = id
	cache := utils.GetCache()
	cache.Delete(m.getRecentsCacheKey(recent.ModeId))
	return nil
}

func (m *MatchRaceService) UpdateMatchRecent(recent *MatchRecent) error {
	if recent.ModeId <= 0 {
		return fmt.Errorf("关联模型id未设置")
	}
	if recent.Player <= 0 {
		return fmt.Errorf("选手未设置")
	}
	o := dbs.NewDefaultOrm()
	_, err := o.Update(recent, "player", "m1", "m2", "m3", "m4", "m5", "displayorder", "disabled")
	cache := utils.GetCache()
	cache.Delete(m.getRecentsCacheKey(recent.ModeId))
	return err
}

func (m *MatchRaceService) GetMatchRecents(modeId int64) []*MatchRecent {
	mrs := []*MatchRecent{}
	cache := utils.GetCache()
	err := cache.Get(m.getRecentsCacheKey(modeId), &mrs)
	if err == nil {
		return mrs
	}
	o := dbs.NewDefaultOrm()
	o.QueryTable(&MatchRecent{}).Filter("modeid", modeId).Filter("disabled", false).OrderBy("-displayorder", "-id").All(&mrs)
	cache.Set(m.getRecentsCacheKey(modeId), mrs, 72*time.Hour)
	return mrs
}

func (m *MatchRaceService) getGroupCacheKey(modeId int64) string {
	return fmt.Sprintf("mobile_mr_groups:%d", modeId)
}

func (m *MatchRaceService) CreateMatchGroup(group *MatchGroup) error {
	if group.ModeId <= 0 {
		return fmt.Errorf("关联模型id未设置")
	}
	o := dbs.NewDefaultOrm()
	id, err := o.Insert(group)
	if err != nil {
		return err
	}
	group.Id = id
	cache := utils.GetCache()
	cache.Delete(m.getGroupCacheKey(group.ModeId))
	return nil
}

func (m *MatchRaceService) UpdateMatchGroup(group *MatchGroup) error {
	if group.ModeId <= 0 {
		return fmt.Errorf("关联模型id未设置")
	}
	o := dbs.NewDefaultOrm()
	_, err := o.Update(group, "title", "displayorder", "posttime")
	cache := utils.GetCache()
	cache.Delete(m.getGroupCacheKey(group.ModeId))
	return err
}

func (m *MatchRaceService) GetMatchGroups(modeId int64) []*MatchGroup {
	mgs := []*MatchGroup{}
	cache := utils.GetCache()
	err := cache.Get(m.getGroupCacheKey(modeId), &mgs)
	if err == nil {
		return mgs
	}
	o := dbs.NewDefaultOrm()
	o.QueryTable(&MatchGroup{}).Filter("modeid", modeId).OrderBy("-displayorder", "title", "-id").All(&mgs)
	cache.Set(m.getGroupCacheKey(modeId), mgs, 72*time.Hour)
	return mgs
}

func (m *MatchRaceService) getMGPCacheKey(groupId int64) string {
	return fmt.Sprintf("mobile_mr_gplayers:%d", groupId)
}

func (m *MatchRaceService) CreateMatchGroupPlayer(mgp *MatchGroupPlayer) error {
	if mgp.GroupId <= 0 {
		return fmt.Errorf("组id未设置")
	}
	if mgp.Player <= 0 {
		return fmt.Errorf("选手id未设置")
	}
	o := dbs.NewDefaultOrm()
	id, err := o.Insert(mgp)
	if err != nil {
		return err
	}
	mgp.Id = id
	cache := utils.GetCache()
	cache.Delete(m.getMGPCacheKey(mgp.GroupId))
	return nil
}

func (m *MatchRaceService) UpdateMatchGroupPlayer(mgp *MatchGroupPlayer) error {
	if mgp.GroupId <= 0 {
		return fmt.Errorf("组id未设置")
	}
	if mgp.Player <= 0 {
		return fmt.Errorf("选手id未设置")
	}
	o := dbs.NewDefaultOrm()
	_, err := o.Update(mgp)
	if err != nil {
		return err
	}
	cache := utils.GetCache()
	cache.Delete(m.getMGPCacheKey(mgp.GroupId))
	return nil
}

func (m *MatchRaceService) GetMatchGroupPlayers(groupId int64) []*MatchGroupPlayer {
	mgps := []*MatchGroupPlayer{}
	cache := utils.GetCache()
	err := cache.Get(m.getMGPCacheKey(groupId), &mgps)
	if err == nil {
		return mgps
	}
	o := dbs.NewDefaultOrm()
	o.QueryTable(&MatchGroupPlayer{}).Filter("groupid", groupId).OrderBy("-displayorder", "-points").All(&mgps)
	cache.Set(m.getMGPCacheKey(groupId), mgps, 72*time.Hour)
	return mgps
}

func (m *MatchRaceService) GetMatchGroupPlayerForAdmin(id int64) *MatchGroupPlayer {
	mp := &MatchGroupPlayer{}
	mp.Id = id
	o := dbs.NewDefaultOrm()
	err := o.Read(mp)
	if err == nil {
		return mp
	}
	return nil
}

func (m *MatchRaceService) getEliminMssCacheKey(modeId int64) string {
	return fmt.Sprintf("mobile_mr_eliminmss:%d", modeId)
}

func (m *MatchRaceService) CreateEliminMs(elimin *MatchEliminMs) error {
	if len(elimin.Title) == 0 {
		return fmt.Errorf("标题不能为空")
	}
	if elimin.ModeId <= 0 {
		return fmt.Errorf("模型id错误")
	}
	o := dbs.NewDefaultOrm()
	id, err := o.Insert(elimin)
	if err != nil {
		return err
	}
	elimin.Id = id
	cache := utils.GetCache()
	cache.Delete(m.getEliminMssCacheKey(elimin.ModeId))
	return nil
}

func (m *MatchRaceService) UpdateEliminMs(elimin *MatchEliminMs) error {
	if len(elimin.Title) == 0 {
		return fmt.Errorf("标题不能为空")
	}
	if elimin.ModeId <= 0 {
		return fmt.Errorf("模型id错误")
	}
	if elimin.Id <= 0 {
		return fmt.Errorf("id错误")
	}
	o := dbs.NewDefaultOrm()
	_, err := o.Update(elimin)
	cache := utils.GetCache()
	cache.Delete(m.getEliminMssCacheKey(elimin.ModeId))
	return err
}

func (m *MatchRaceService) DeleteEliminMs(id int64) error {
	o := dbs.NewDefaultOrm()
	o.QueryTable(&MatchEliminMs{}).Filter("id", id).Delete()
	o.QueryTable(&MatchEliminVs{}).Filter("msid", id).Delete()
	cache := utils.GetCache()
	cache.Delete(m.getEliminMssCacheKey(id))
	cache.Delete(m.getEliminVssCacheKey(id))
	return nil
}

func (m *MatchRaceService) GetEliminMss(modeId int64) []*MatchEliminMs {
	ckey := m.getEliminMssCacheKey(modeId)
	cache := utils.GetCache()
	mess := []*MatchEliminMs{}
	err := cache.Get(ckey, &mess)
	if err == nil {
		return mess
	}
	o := dbs.NewDefaultOrm()
	o.QueryTable(&MatchEliminMs{}).Filter("modeid", modeId).OrderBy("-displayorder", "posttime").All(&mess)
	cache.Set(ckey, mess, 72*time.Hour)
	return mess
}

func (m *MatchRaceService) GetEliminMsForAdmin(id int64) *MatchEliminMs {
	ms := &MatchEliminMs{}
	ms.Id = id
	o := dbs.NewDefaultOrm()
	err := o.Read(ms)
	if err != nil {
		return nil
	}
	return ms
}

func (m *MatchRaceService) getEliminVssCacheKey(msid int64) string {
	return fmt.Sprintf("mobile_mr_eliminvss:%d", msid)
}

func (m *MatchRaceService) CreateEliminVs(vs *MatchEliminVs) error {
	if vs.MsId <= 0 {
		return fmt.Errorf("赛程组id错误")
	}
	if vs.VsId <= 0 {
		return fmt.Errorf("对阵id错误")
	}
	vs.PostTime = time.Now().Unix()
	o := dbs.NewDefaultOrm()
	id, err := o.Insert(vs)
	if err != nil {
		return err
	}
	vs.Id = id
	cache := utils.GetCache()
	cache.Delete(m.getEliminVssCacheKey(vs.MsId))
	return nil
}

func (m *MatchRaceService) UpdateEliminVs(vs *MatchEliminVs) error {
	if vs.Id <= 0 {
		return fmt.Errorf("id错误")
	}
	if vs.MsId <= 0 {
		return fmt.Errorf("赛程组id错误")
	}
	if vs.VsId <= 0 {
		return fmt.Errorf("对阵id错误")
	}
	o := dbs.NewDefaultOrm()
	_, err := o.Update(vs)
	if err != nil {
		return err
	}
	cache := utils.GetCache()
	cache.Delete(m.getEliminVssCacheKey(vs.MsId))
	return nil
}

func (m *MatchRaceService) DeleteEliminVs(id int64) error {
	if id <= 0 {
		return fmt.Errorf("id错误")
	}
	o := dbs.NewDefaultOrm()
	vs := &MatchEliminVs{}
	vs.Id = id
	err := o.Read(vs)
	if err != nil {
		fmt.Errorf("对象不存在")
	}
	o.QueryTable(&MatchEliminVs{}).Filter("id", id).Delete()
	cache := utils.GetCache()
	cache.Delete(m.getEliminVssCacheKey(vs.MsId))
	return nil
}

func (m *MatchRaceService) GetEliminVss(msid int64) []*MatchEliminVs {
	ckey := m.getEliminVssCacheKey(msid)
	cache := utils.GetCache()
	mvss := []*MatchEliminVs{}
	err := cache.Get(ckey, &mvss)
	if err == nil {
		return mvss
	}
	o := dbs.NewDefaultOrm()
	o.QueryTable(&MatchEliminVs{}).Filter("msid", msid).All(&mvss)
	cache.Set(ckey, mvss, 72*time.Hour)
	return mvss
}
