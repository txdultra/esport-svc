package lives

import (
	"dbs"
	"errors"
	"fmt"
	//"github.com/astaxie/beego"
	"libs"
	"libs/search"
	"strconv"
	"time"
	"utils"

	"github.com/astaxie/beego/orm"
)

type LivePrograms struct{}

func (lp *LivePrograms) validate(program *LiveProgram) error {
	//if program.GameId <= 0 {
	//	return errors.New("必须选择所属项目")
	//}
	if len(program.Title) == 0 {
		return errors.New("必须设置标题")
	}
	if program.Date.IsZero() {
		return errors.New("必须设置所属日期")
	}
	if program.StartTime.IsZero() || program.EndTime.IsZero() {
		return errors.New("必须设置开始和结束时间")
	}
	if program.Date.Day() != program.StartTime.Day() {
		return errors.New("日期和开始时间必须在同一天")
	}
	if program.StartTime.After(program.EndTime) {
		return errors.New("结束时间必须大于开始时间")
	}
	if program.MatchId <= 0 {
		return errors.New("必须设置所属赛事")
	}
	if program.DefaultChannelId <= 0 {
		return errors.New("必须设置默认观看频道")
	}
	return nil
}

func (lp *LivePrograms) Create(program *LiveProgram, channelIds []int64) (int64, error) {
	err := lp.validate(program)
	if err != nil {
		return 0, err
	}
	if program.PostTime.IsZero() {
		program.PostTime = time.Now()
	}
	o := dbs.NewDefaultOrm()
	id, err := o.Insert(program)
	if err != nil {
		return 0, err
	}
	program.Id = id

	//同步到搜索
	lp.syncSearchData(program, 0, &o, false)

	cache := utils.GetCache()
	cache.Delete(lp.lpckey(program.Date))

	//添加频道
	lp.CreateProgramChannels(id, channelIds)

	return id, nil
}

func (lp *LivePrograms) lpckey(date time.Time) string {
	return fmt.Sprintf("mobile_programs_day:%d-%d-%d", date.Year(), date.Month(), date.Day())
}

func (lp *LivePrograms) lpkey(id int64) string {
	return fmt.Sprintf("mobile_program_id:%d", id)
}

func (lp *LivePrograms) lppcskey(id int64) string {
	return fmt.Sprintf("mobile_program_%d_channels", id)
}

func (lp *LivePrograms) Update(program *LiveProgram, channelIds []int64) error {
	err := lp.validate(program)
	if err != nil {
		return err
	}
	o := dbs.NewDefaultOrm()
	_, err = o.Update(program)
	if err != nil {
		return err
	}

	//同步到搜索
	lp.syncSearchData(program, 0, &o, false)

	cache := utils.GetCache()
	cache.Delete(lp.lpckey(program.Date))
	cache.Delete(lp.lpkey(program.Id))

	//更新频道
	if len(channelIds) > 0 {
		lp.UpdateProgramChannels(program.Id, channelIds)
	}

	return nil
}

func (lp *LivePrograms) Delete(id int64) error {
	lm := lp.Get(id)
	if lm == nil {
		return errors.New("不存在指定的节目单")
	}
	spms := &LiveSubPrograms{}
	sps := spms.Gets(id)
	if len(sps) > 0 {
		return errors.New("必须删除子节目单")
	}
	o := dbs.NewDefaultOrm()
	_, err := o.Delete(lm)
	if err != nil {
		return err
	}
	lm.Id = id //delete会自动设置主键id为0
	//同步到搜索
	lp.syncSearchData(lm, 0, &o, true)

	cache := utils.GetCache()
	cache.Delete(lp.lpckey(lm.Date))
	cache.Delete(lp.lpkey(lm.Id))
	return nil
}

func (lp *LivePrograms) Get(id int64) *LiveProgram {
	ckey := lp.lpkey(id)
	cache := utils.GetCache()
	lm := LiveProgram{}
	err := cache.Get(ckey, &lm)
	if err == nil {
		return &lm
	}
	o := dbs.NewDefaultOrm()
	lm.Id = id
	err = o.Read(&lm)
	if err == orm.ErrNoRows || err == orm.ErrMissPK || err != nil {
		return nil
	}
	cache.Set(ckey, lm, utils.StrToDuration("12h"))
	return &lm
}

func (lp *LivePrograms) GetsByDate(opt string, date time.Time, channelId int64) []int64 {
	o := dbs.NewDefaultOrm()
	d := fmt.Sprintf("%d-%d-%d", date.Year(), int(date.Month()), date.Day())
	sql := fmt.Sprintf("select distinct(lpc.pid) from live_program_channels lpc join live_program lp on lp.id=lpc.pid where lp.d "+opt+"'%s' ", d)
	if channelId > 0 {
		sql += fmt.Sprintf(" and lpc.cid=%d ", channelId)
	}
	sql += " order by start_time "

	var lists []orm.ParamsList
	lpIds := []int64{}
	_, err := o.Raw(sql).ValuesList(&lists, "pid")
	if err == nil {
		for _, row := range lists {
			_id, err := strconv.ParseInt(row[0].(string), 10, 64)
			//_id := row[0].(int64)
			if err == nil {
				lpIds = append(lpIds, _id)
			}
		}
	}
	return lpIds
}

func (lp *LivePrograms) Gets(date time.Time) []*LiveProgram {
	ckey := lp.lpckey(date)
	cache := utils.GetCache()
	lps := []*LiveProgram{}
	err := cache.Get(ckey, &lps)
	if err == nil {
		return lps
	}
	o := dbs.NewDefaultOrm()
	qs := o.QueryTable(LiveProgram{})
	_, err = qs.Filter("d", date).OrderBy("start_time").All(&lps)
	if err == nil {
		cache.Set(ckey, lps, utils.StrToDuration("6h"))
	}
	return lps
}

//大于等于某日的节目单
func (lp *LivePrograms) getGteDate(gteDate time.Time) []LiveProgram {
	lps := []LiveProgram{}
	o := dbs.NewDefaultOrm()
	qs := o.QueryTable(LiveProgram{})
	qs.Filter("d__gte", gteDate).OrderBy("start_time").All(&lps)
	return lps
}

func (lp *LivePrograms) syncSearchData(program *LiveProgram, cw int, ormer *orm.Ormer, isDel bool) {
	if isDel {
		(*ormer).QueryTable(&LiveProgramSearchData{}).Filter("id", program.Id).Delete()
		return
	}
	var search_data LiveProgramSearchData
	err := (*ormer).QueryTable(&LiveProgramSearchData{}).Filter("id", program.Id).One(&search_data)
	orgs := &LiveOrgs{}
	bas := &libs.Bas{}
	channel := orgs.GetChannel(program.DefaultChannelId)
	game_name := ""
	if program.GameId > 0 {
		game := bas.GetGame(program.GameId)
		game_name = fmt.Sprintf("%s %s", game.Name, game.En)
	}
	if err == orm.ErrNoRows {
		search_data := LiveProgramSearchData{
			ProgramId:    program.Id,
			Title:        program.Title,
			SubTitle:     program.SubTitle,
			ChannelName:  channel.Name,
			GameName:     game_name,
			CustomWeight: cw,
			StartTime:    program.StartTime,
			EndTime:      program.EndTime,
		}
		(*ormer).Insert(&search_data)
	} else {
		search_data.ProgramId = program.Id
		search_data.GameName = game_name
		search_data.CustomWeight = cw
		search_data.StartTime = program.StartTime
		search_data.EndTime = program.EndTime
		(*ormer).Update(&search_data)
	}
}

//搜索
func (lp *LivePrograms) Query(words string, p int, s int, match_mode string, filters []search.SearchFilter, filterRanges []search.FilterRangeInt) (int, []LiveProgram) {
	search_config := &search.SearchOptions{
		Host:    search_program_server,
		Port:    search_program_port,
		Timeout: search_program_timeout,
	}
	offset := (p - 1) * s
	search_config.Offset = offset
	search_config.Limit = s
	search_config.Filters = filters
	search_config.FilterRangeInt = filterRanges
	search_config.MaxMatches = 500
	q_engine := search.NewSearcher(search_config)
	ids, total, err := q_engine.ProgramQuery(words, match_mode)
	programs := []LiveProgram{}
	if err != nil {
		return 0, programs
	}
	for _, id := range ids {
		p := lp.Get(id)
		if p != nil {
			programs = append(programs, *p)
		}
	}
	return total, programs
}

func (lp *LivePrograms) CreateProgramChannels(programId int64, channelIds []int64) error {
	o := dbs.NewDefaultOrm()
	o.Begin()
	for _, cid := range channelIds {
		_, err := o.Raw("insert live_program_channels(pid,cid) values(?,?)", programId, cid).Exec()
		if err != nil {
			o.Rollback()
			return err
		}
	}
	err := o.Commit()
	if err == nil {
		cache := utils.GetCache()
		cache.Delete(lp.lppcskey(programId))
	}
	return err
}

func (lp *LivePrograms) UpdateProgramChannels(programId int64, channelIds []int64) error {
	o := dbs.NewDefaultOrm()
	o.Begin()
	o.Raw(fmt.Sprintf("delete from live_program_channels where pid=%d", programId)).Exec()
	for _, cid := range channelIds {
		_, err := o.Raw("insert live_program_channels(pid,cid) values(?,?)", programId, cid).Exec()
		if err != nil {
			o.Rollback()
			return err
		}
	}
	err := o.Commit()
	if err == nil {
		cache := utils.GetCache()
		cache.Delete(lp.lppcskey(programId))
	}
	return err
}

func (lp *LivePrograms) GetChannelIds(programId int64) []int64 {
	ckey := lp.lppcskey(programId)
	cache := utils.GetCache()
	cids := []int64{}
	err := cache.Get(ckey, &cids)
	if err == nil {
		return cids
	}
	sql := fmt.Sprintf("select cid from live_program_channels where pid=%d", programId)
	o := dbs.NewDefaultOrm()
	var maps []orm.Params
	num, err := o.Raw(sql).Values(&maps)
	if err == nil && num > 0 {
		for _, row := range maps {
			str_id := row["cid"].(string)
			_id, err := strconv.ParseInt(str_id, 10, 64)
			if err == nil {
				cids = append(cids, _id)
			}
		}
	}
	if err == nil {
		cache.Set(ckey, cids, utils.StrToDuration("6h"))
	}
	return cids
}

//二级节目单

type LiveSubPrograms struct{}

func (lsp *LiveSubPrograms) validate(sp *LiveSubProgram) error {
	if sp.ViewType == 0 {
		return errors.New("未选择类型")
	}
	if sp.ProgramId <= 0 {
		return errors.New("未设定对应节目单")
	}
	lp := &LivePrograms{}
	lm := lp.Get(int64(sp.ProgramId))
	if lm == nil {
		return errors.New("不存在对应的节目单")
	}
	if sp.StartTime.Before(lm.Date) {
		return errors.New("开始时间错误")
	}
	if sp.ViewType == LIVE_SUBPROGRAM_VIEW_VS {
		if len(sp.Vs1Name) == 0 || len(sp.Vs2Name) == 0 {
			return errors.New("未设置对阵对象名称")
		}
		if sp.GameId <= 0 {
			return errors.New("未设定关联项目")
		}
	}
	if sp.ViewType == LIVE_SUBPROGRAM_VIEW_SINGLE {
		if len(sp.Title) == 0 {
			return errors.New("未设置标题")
		}
	}
	return nil
}

func (lsp *LiveSubPrograms) lpskey(programId int64) string {
	return fmt.Sprintf("mobile_program_vss_pid:%d", programId)
}

func (lsp *LiveSubPrograms) lpkey(id int64) string {
	return fmt.Sprintf("mobile_subprogram_id:%d", id)
}

func (lsp *LiveSubPrograms) Create(sp *LiveSubProgram) (int64, error) {
	err := lsp.validate(sp)
	if err != nil {
		return 0, err
	}
	if sp.PostTime.IsZero() {
		sp.PostTime = time.Now()
	}
	o := dbs.NewDefaultOrm()
	id, err := o.Insert(sp)
	if err != nil {
		return 0, err
	}
	sp.Id = id
	cache := utils.GetCache()
	cache.Delete(lsp.lpskey(sp.ProgramId))
	//同步到搜索
	//lsp.syncSearchData(sp, 0, &o, false)
	lsp.ResetProgramTimeEdges(sp.ProgramId)
	//开启提醒线程
	notice_service := NewProgramNoticeService()
	notice_service.StartNoticeTimer(id)
	return id, nil
}

//重置主节目单开始和结束时间
func (lsp *LiveSubPrograms) ResetProgramTimeEdges(programId int64) {
	allsubs := lsp.Gets(programId)
	lps := &LivePrograms{}
	program := lps.Get(programId)
	if program == nil {
		return
	}
	if len(allsubs) == 0 {
		program.EndTime = program.StartTime
		lps.Update(program, nil)
		return
	}
	startTime := allsubs[0].StartTime
	endTime := allsubs[0].EndTime
	for _, sub := range allsubs {
		if startTime.After(sub.StartTime) {
			startTime = sub.StartTime
		}
		if endTime.Before(sub.EndTime) {
			endTime = sub.EndTime
		}
	}
	program.StartTime = startTime
	program.EndTime = endTime
	lps.Update(program, nil)
}

func (lsp *LiveSubPrograms) IsLocked(id int64) bool {
	_dbsp := lsp.Get(id)
	if _dbsp == nil {
		return true
	}
	if _dbsp.StartTime.Before(time.Now().Add(program_lock_dur_str)) {
		return true
	}
	return false
}

func (lsp *LiveSubPrograms) IsLiving(id int64) bool {
	_dbsp := lsp.Get(id)
	if _dbsp == nil {
		return false
	}
	if _dbsp.StartTime.Before(time.Now()) && _dbsp.EndTime.After(time.Now()) {
		return true
	}
	return false
}

func (lsp *LiveSubPrograms) Update(sp LiveSubProgram) error {
	err := lsp.validate(&sp)
	if err != nil {
		return err
	}
	if lsp.IsLocked(sp.Id) {
		return errors.New("节目单已被锁定,不能修改")
	}

	o := dbs.NewDefaultOrm()
	_, err = o.Update(&sp)
	if err != nil {
		return err
	}
	cache := utils.GetCache()
	cache.Delete(lsp.lpkey(sp.Id))
	cache.Delete(lsp.lpskey(sp.ProgramId))

	//同步到搜索
	lsp.ResetProgramTimeEdges(sp.ProgramId)
	//重置提醒时间
	notice_service := NewProgramNoticeService()
	notice_service.ResetNoticeTimer(sp.Id, sp.StartTime)
	return nil
}

//更改二级节目单(用于锁定中的二级节目单)
func (lsp *LiveSubPrograms) UpdateLockingSubProgram(sp LiveSubProgram) error {
	err := lsp.validate(&sp)
	if err != nil {
		return err
	}
	o := dbs.NewDefaultOrm()
	_, err = o.Update(&sp, "gid", "vs1", "vs1_uid", "vs1_img", "vs2", "vs2_uid", "vs2_img", "title", "img", "end_time")
	if err != nil {
		return err
	}
	cache := utils.GetCache()
	cache.Delete(lsp.lpkey(sp.Id))
	cache.Delete(lsp.lpskey(sp.ProgramId))

	//同步到搜索
	lsp.ResetProgramTimeEdges(sp.ProgramId)
	return nil
}

func (lsp *LiveSubPrograms) Delete(id int64) error {
	lm := lsp.Get(id)
	if lm == nil {
		return errors.New("不存在指定的预告")
	}
	if lsp.IsLocked(id) {
		return errors.New("节目单已被锁定,不能删除")
	}
	o := dbs.NewDefaultOrm()
	_, err := o.Delete(lm)
	if err != nil {
		return err
	}
	lm.Id = id
	cache := utils.GetCache()
	cache.Delete(lsp.lpskey(lm.ProgramId))
	cache.Delete(lsp.lpkey(lm.Id))

	//同步到搜索
	//lsp.syncSearchData(lm, 0, &o, true)
	lsp.ResetProgramTimeEdges(lm.ProgramId)
	//删除提醒线程
	notice_service := NewProgramNoticeService()
	//notice_service.StopNoticeTimer(lm.Id)
	notice_service.RemoveAllSubscribeEvent(lm.Id)
	return nil
}

func (lsp *LiveSubPrograms) Get(id int64) *LiveSubProgram {
	cache := utils.GetCache()
	lp := LiveSubProgram{}
	err := cache.Get(lsp.lpkey(id), &lp)
	if err == nil {
		return &lp
	}
	o := dbs.NewDefaultOrm()
	lp.Id = id
	err = o.Read(&lp)
	if err == orm.ErrNoRows || err == orm.ErrMissPK || err != nil {
		return nil
	}
	cache.Set(lsp.lpkey(id), lp, utils.StrToDuration("6h"))
	return &lp
}

func (lsp *LiveSubPrograms) Gets(programId int64) []*LiveSubProgram {
	cache := utils.GetCache()
	lpids := []int64{}
	err := cache.Get(lsp.lpskey(programId), &lpids)
	lps := []*LiveSubProgram{}
	if err == nil {
		for _, lpid := range lpids {
			lp := lsp.Get(lpid)
			if lp != nil {
				lps = append(lps, lp)
			}
		}
		return lps
	}
	o := dbs.NewDefaultOrm()
	_, err = o.QueryTable(LiveSubProgram{}).Filter("pid", programId).OrderBy("start_time", "id").All(&lps)
	if err != nil {
		return lps
	}
	for _, lp := range lps {
		lpids = append(lpids, lp.Id)
		cache.Set(lsp.lpkey(lp.Id), *lp, utils.StrToDuration("6h"))
	}
	cache.Set(lsp.lpskey(programId), lpids, utils.StrToDuration("6h"))
	return lps
}
