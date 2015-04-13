package lives

import (
	"dbs"
	"errors"
	"fmt"
	"libs"
	"libs/collect"
	"libs/passport"
	"libs/reptile"
	"libs/search"
	"logs"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"utils"
	"utils/redis"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
)

//客户端参与的抓取方法
func ClientProxyReptile(rep_support reptile.REP_SUPPORT, parameter string, state string) (clientReqUrl string, nextState string, err error) {
	t, ok := reptile.LIVE_REPTILE_MODULES[string(rep_support)]
	if ok {
		_obj := reflect.New(t).Interface()
		if service, ok := _obj.(reptile.IReptileLiveClientProxy); ok {
			return service.ProxyReptile(parameter, state)
		}
	}
	return "", reptile.REP_CALLBACK_COMMAND_EXIT, errors.New("没有支持的方法")
}

var per_lives map[int64]*LivePerson = make(map[int64]*LivePerson)
var per_lives_deleted map[int64]bool = make(map[int64]bool)
var mutex *sync.RWMutex = new(sync.RWMutex)

const (
	living_personal_list = "mobile_living_personal_list"
)

type LivePers struct{}

func (lv *LivePers) init() {
	o := dbs.NewDefaultOrm()
	var plives []*LivePerson
	_, err := o.QueryTable(&LivePerson{}).All(&plives)
	if err == nil {
		mutex.Lock()
		for _, p := range plives {
			per_lives[p.Id] = p
		}
		mutex.Unlock()
	}
	lv.loadTask() //启动直播状态抓取任务
}

func (lv *LivePers) loadTask() {
	spec := "0 */10 * * * *"
	_spec := beego.AppConfig.String("reptile.status.spec")
	if len(_spec) > 0 {
		spec = _spec
	}
	tskName := "personal_live_status_task"
	toolbox.AddTask(tskName, toolbox.NewTask(tskName, spec, func() error {
		fmt.Println("personal live status task running...")

		//更新在线人数
		go lv.upateSearchDataOnlines()

		for id, _ := range per_lives {
			v := lv.Get(id)

			sup := v.Rep
			t, ok := reptile.LIVE_REPTILE_MODULES[string(sup)]
			if !ok {
				continue
			}
			service := reflect.New(t).Interface().(reptile.ILiveStatus)
			status, err := service.GetStatus(v.ReptileUrl)
			if err != nil {
				v.ReptileDes += fmt.Sprintf("%s:%s\n\r", time.Now().String(), err.Error())
			} else {
				v.ReptileDes = ""
			}
			if status != v.LiveStatus { //和原状态相同则不更新
				v.LiveStatus = status
				v.LastTime = time.Now()
				lv.Update(*v)

				//由直播变非直播状态,重置在线人数;非直播变直播状态，添加伪在线
				oc := &OnlineCounter{}
				if status == reptile.LIVE_STATUS_NOTHING {
					go oc.Reset(LIVE_TYPE_PERSONAL, int(v.Id))
				} else {
					go oc.AddRandomViewers(LIVE_TYPE_PERSONAL, int(v.Id), v.ShowOnlineMin, v.ShowOnlineMax)
				}
			}
			switch status {
			case reptile.LIVE_STATUS_LIVING:
				lv.AddLivingChannelToList(v.Id)
			case reptile.LIVE_STATUS_NOTHING:
				lv.RemoveLivingChannelInList(v.Id)
			}

			//只在直播状态时更新直播流
			if status == reptile.LIVE_STATUS_LIVING && reptile.LiveRepMethod(v.Rep) == reptile.REP_METHOD_DIRECT { //&& v.RepMethod == reptile.REP_METHOD_DIRECT {
				lv.reptiling(*v)
			}
		}
		return nil
	}))
}

func (lv *LivePers) AddLivingChannelToList(id int64) {
	redis.HSet(nil, living_personal_list, strconv.FormatInt(id, 10), 1)
}

func (lv *LivePers) RemoveLivingChannelInList(id int64) {
	redis.HDel(nil, living_personal_list, strconv.FormatInt(id, 10))
}

func (lv *LivePers) perCacheKey(id int64) string {
	return fmt.Sprintf("mobile_live_personal_id:%d", id)
}

func (lv *LivePers) Get(id int64) *LivePerson {
	ckey := lv.perCacheKey(id)
	cache := utils.GetCache()
	lp := LivePerson{}
	err := cache.Get(ckey, &lp)
	if err == nil {
		return &lp
	}
	//减少db读取
	mutex.RLock()
	if _, ok := per_lives_deleted[id]; ok {
		mutex.RUnlock()
		return nil
	}
	mutex.RUnlock()

	o := dbs.NewDefaultOrm()
	lp.Id = id
	err = o.Read(&lp)
	if err == orm.ErrNoRows || err == orm.ErrMissPK || err != nil {
		mutex.Lock()
		defer mutex.Unlock()
		per_lives_deleted[id] = true
		return nil
	}
	cache.Set(ckey, lp, utils.StrToDuration("12h"))
	return &lp
}

func (lv *LivePers) perGamesCacheKey(id int64) string {
	return fmt.Sprintf("mobile_live_personal_gameids_id:%d", id)
}

func (lv *LivePers) GetOfGames(id int64) []int {
	ckey := lv.perGamesCacheKey(id)
	cache := utils.GetCache()
	lp := []int{}
	err := cache.Get(ckey, &lp)
	if err == nil {
		return lp
	}
	o := dbs.NewDefaultOrm()
	type M struct {
		Gid int
	}
	var row_ids []M
	var ids []int
	_, err = o.Raw("SELECT gid FROM live_personal_games WHERE per_id = ?", id).QueryRows(&row_ids)
	if err == nil {
		for _, m := range row_ids {
			ids = append(ids, m.Gid)
		}
	}
	cache.Set(ckey, ids, utils.StrToDuration("12h"))
	return ids
}

func (lv *LivePers) Create(per LivePerson, ofGames []int, doRep bool) (int64, error) {
	if len(per.Name) == 0 {
		return 0, errors.New("名称不能为空")
	}
	if len(per.ReptileUrl) == 0 {
		return 0, errors.New("抓取地址不能为空")
	}
	if per.Uid <= 0 {
		return 0, errors.New("必须有所属主播用户")
	}
	if per.Succ {
		return 0, errors.New("未抓取前不能设为成功")
	}
	if per.PostTime.IsZero() {
		per.PostTime = time.Now()
	}
	if per.LastTime.IsZero() {
		per.LastTime = time.Now()
	}
	repSupport, err := reptile.Get_REP_SUPPORT(per.ReptileUrl)
	if err != nil {
		return 0, errors.New("不支持此地址抓取")
	}
	per.Rep = repSupport
	per.LiveStatus = reptile.LIVE_STATUS_NOTHING
	//per.RepMethod = reptile.LiveRepMethod(repSupport)

	bas := &libs.Bas{}
	for _, ofgid := range ofGames {
		g := bas.GetGame(ofgid)
		if g == nil {
			return 0, errors.New("某个所属游戏不存在")
		}
	}

	o := dbs.NewDefaultOrm()

	count, _ := o.QueryTable(LivePerson{}).Filter("reptile_url", per.ReptileUrl).Count()
	if count > 0 {
		return 0, errors.New("已存在相同的个人直播")
	}

	o.Begin()
	id, err := o.Insert(&per)
	if err != nil {
		logs.Error(err.Error())
		o.Rollback()
		return 0, err
	}

	for _, gid := range ofGames {
		_, err := o.Raw("insert live_personal_games(per_id,gid) values(?,?)", id, gid).Exec()
		if err != nil {
			logs.Error(err.Error())
			o.Rollback()
			return 0, err
		}
	}
	//保存进搜索库中
	lv.syncSearchData(&per, ofGames, 0, &o)

	err = o.Commit()
	if err != nil {
		logs.Error(err.Error())
		o.Rollback()
		return 0, errors.New("提交事务时错误:" + err.Error())
	}

	per.Id = id
	cache := utils.GetCache()
	ckey := lv.perCacheKey(per.Id)
	mutex.Lock()
	per_lives[per.Id] = &per
	delete(per_lives_deleted, id)
	mutex.Unlock()
	cache.Set(ckey, per, utils.StrToDuration("12h"))
	if doRep {
		lv.reptiling(per) //创建成功后直接抓取
	}
	return id, nil
}

func (lv *LivePers) upateSearchDataOnlines() {
	o := dbs.NewDefaultOrm()
	counter := &OnlineCounter{}
	lsd := LiveSearchData{}
	tbl := lsd.TableName()
	p, _ := o.Raw("UPDATE " + tbl + " SET onlines = ? WHERE id = ?").Prepare()
	defer p.Close()
	for id, _ := range per_lives {
		onlines := counter.GetChannelCounts(LIVE_TYPE_PERSONAL, int(id))
		p.Exec(onlines, id)
	}
}

func (lv *LivePers) syncSearchData(live *LivePerson, gameIds []int, cw int, ormer *orm.Ormer) {
	ofSearchGameNames := []string{}
	bases := libs.Bas{}
	for _, gameId := range gameIds {
		game := bases.GetGame(gameId)
		if game != nil {
			ofSearchGameNames = append(ofSearchGameNames, fmt.Sprintf("%s_%s", game.Name, game.En))
		}
	}
	gids := []string{}
	for _, gid := range gameIds {
		gids = append(gids, strconv.Itoa(gid))
	}
	ppt := passport.NewMemberProvider()
	ofAnchor := ppt.GetNicknameByUid(live.Uid)
	var search_data LiveSearchData
	err := (*ormer).QueryTable(&LiveSearchData{}).Filter("id", live.Id).One(&search_data)
	if err == orm.ErrNoRows {
		search_data := LiveSearchData{
			LiveId:              live.Id,
			PersonalChannelName: live.Name,
			AnchorName:          ofAnchor,
			Status:              reptile.LIVE_STATUS_NOTHING,
			GameNames:           strings.Join(ofSearchGameNames, ","),
			GameIds:             strings.Join(gids, ","),
			CustomWeight:        cw,
			Enabled:             live.Enabled,
		}
		(*ormer).Insert(&search_data)
	} else {
		search_data.PersonalChannelName = live.Name
		search_data.AnchorName = ofAnchor
		search_data.Status = live.LiveStatus
		search_data.GameNames = strings.Join(ofSearchGameNames, ",")
		search_data.GameIds = strings.Join(gids, ",")
		search_data.CustomWeight = cw
		search_data.Enabled = live.Enabled
		(*ormer).Update(&search_data)
	}
}

func (lv *LivePers) reptiling(per LivePerson) error {
	t, ok := reptile.LIVE_REPTILE_MODULES[string(per.Rep)]
	if !ok {
		return errors.New("不存在对应的抓取模块")
	}
	rep := reflect.New(t).Interface().(reptile.IReptileLive)
	stream, err := rep.Reptile(per.ReptileUrl)
	if err == nil {
		per.StreamUrl = stream
		per.Succ = true
		per.LastTime = time.Now()
	} else {
		per.Succ = false
		per.LastTime = time.Now()
		per.ReptileDes = err.Error()
	}
	o := dbs.NewDefaultOrm()
	o.Update(&per)
	//更新缓存
	cache := utils.GetCache()
	ckey := lv.perCacheKey(per.Id)
	mutex.Lock()
	per_lives[per.Id] = &per
	mutex.Unlock()

	cache.Replace(ckey, per, utils.StrToDuration("12h"))
	return err
}

func (lv *LivePers) Reptiling(liveId int64) error {
	live := lv.Get(liveId)
	if live == nil {
		return errors.New("个人直播频道不存在")
	}
	return lv.reptiling(*live)
}

func (lv *LivePers) UpdateGames(liveId int64, ofGames []int) error {
	per := lv.Get(liveId)
	if per == nil {
		return errors.New("个人直播频道不存在")
	}
	o := dbs.NewDefaultOrm()
	o.Begin()
	_, err := o.Raw("delete from live_personal_games where per_id=?", per.Id).Exec()
	if err != nil {
		logs.Error(err.Error())
		o.Rollback()
		return err
	}
	for _, gid := range ofGames {
		_, err := o.Raw("insert live_personal_games(per_id,gid) values(?,?)", per.Id, gid).Exec()
		if err != nil {
			logs.Error(err.Error())
			o.Rollback()
			return err
		}
	}
	//同步到搜索库中
	lv.syncSearchData(per, ofGames, 0, &o)

	err = o.Commit()
	if err != nil {
		logs.Error(err.Error())
		o.Rollback()
		return errors.New("提交事务时错误:" + err.Error())
	}
	//缓存更新
	cache := utils.GetCache()
	cache.Delete(lv.perGamesCacheKey(per.Id))
	return nil
}

func (lv *LivePers) Update(per LivePerson) error {
	if len(per.Name) == 0 {
		return errors.New("名称不能为空")
	}
	if len(per.ReptileUrl) == 0 {
		return errors.New("抓取地址不能为空")
	}
	if per.Uid <= 0 {
		return errors.New("必须有所属主播用户")
	}
	if per.PostTime.IsZero() {
		per.PostTime = time.Now()
	}
	if per.LastTime.IsZero() {
		per.LastTime = time.Now()
	}
	repSupport, err := reptile.Get_REP_SUPPORT(per.ReptileUrl)
	if err != nil {
		return errors.New("不支持此地址抓取")
	}
	per.Rep = repSupport
	//per.RepMethod = reptile.LiveRepMethod(repSupport)

	o := dbs.NewDefaultOrm()
	o.Begin()
	_, err = o.Update(&per)
	if err != nil {
		logs.Error(err.Error())
		o.Rollback()
		return err
	}
	//同步到搜索库中
	ofGames := lv.GetOfGames(per.Id)
	lv.syncSearchData(&per, ofGames, 0, &o)

	err = o.Commit()
	if err != nil {
		o.Rollback()
		return errors.New("提交事务时错误:" + err.Error())
	}
	//缓存更新
	cache := utils.GetCache()
	cache.Replace(lv.perCacheKey(per.Id), per, utils.StrToDuration("12h"))
	mutex.Lock()
	per_lives[per.Id] = &per
	mutex.Unlock()

	//删除正在直播列表
	if !per.Enabled {
		lv.RemoveLivingChannelInList(per.Id)
	}

	return nil
}

func (lv *LivePers) Delete(id int64) error {
	o := dbs.NewDefaultOrm()
	o.Begin()

	o.QueryTable(&LiveSearchData{}).Filter("id", id).Delete()          //删除搜索
	o.Raw("delete from live_personal_games where per_id=?", id).Exec() //删除所属游戏
	o.QueryTable(&LivePerson{}).Filter("id", id).Delete()

	err := o.Commit()
	if err != nil {
		logs.Error(err.Error())
		o.Rollback()
		return errors.New("提交事务时错误:" + err.Error())
	}
	collect.DeleteAll(strconv.Itoa(int(id)), string(collect.COLLECT_RELTYPE_PLIVE))
	cache := utils.GetCache()
	cache.Delete(lv.perCacheKey(id))
	cache.Delete(lv.perGamesCacheKey(id))
	mutex.Lock()
	delete(per_lives, id)
	mutex.Unlock()

	//删除正在直播列表
	lv.RemoveLivingChannelInList(id)

	return nil
}

func (lv *LivePers) persCacheKey(gid int, status reptile.LIVE_STATUS) string {
	return fmt.Sprintf("mobile_live_personal_ids:gid:%d_status:%d", gid, status)
}

func (lv *LivePers) ListForAdmin(query string, gameId int, p int, s int) (int, []*LivePerson) {
	sql := "select {0} from live_personal lp join live_personal_games lpg on lpg.per_id = lp.id where 1=1 "
	if gameId > 0 {
		sql += " and lpg.gid = " + strconv.Itoa(gameId) + " "
	} else {
		sql = "select {0} from live_personal lp where 1=1" //有未编组游戏归类的个直
	}
	if len(query) > 0 {
		sql += " and lp.title like '%" + query + "%' "
	}
	if p <= 0 {
		p = 1
	}
	total := 0
	o := dbs.NewDefaultOrm()
	var lists []orm.ParamsList
	total_sql := strings.Replace(sql, "{0}", "count(lp.id)", -1)
	num, err := o.Raw(total_sql).ValuesList(&lists)
	if err == nil && num > 0 {
		total, _ = strconv.Atoi(lists[0][0].(string))
	}

	sql += " order by post_time desc limit " + strconv.Itoa((p-1)*s) + "," + strconv.Itoa(s)
	var maps []orm.Params
	lpIds := []int64{}
	select_sql := strings.Replace(sql, "{0}", "lp.id", -1)
	num, err = o.Raw(select_sql).Values(&maps)
	if err == nil && num > 0 {
		for _, row := range maps {
			str_id := row["id"].(string)
			_id, err := strconv.ParseInt(str_id, 10, 64)
			if err == nil {
				lpIds = append(lpIds, _id)
			}
		}
	}
	lps := []*LivePerson{}
	for _, id := range lpIds {
		lp := lv.Get(id)
		if lp != nil {
			lps = append(lps, lp)
		}
	}
	return total, lps
}

//gid = 0 all lives, status = nil all status
//func (lv *LivePers) Gets(p int, s int, filterFunc func(*LivePerson) bool) (*libs.PL, error) {
//	mutex.RLock()
//	pids := []int64{}
//	utils.Try(func() {
//		if filterFunc != nil {
//			for k, v := range per_lives {
//				ok := filterFunc(v)
//				if !ok {
//					break
//				} else {
//					pids = append(pids, k)
//				}
//			}
//		} else {
//			for k, _ := range per_lives {
//				pids = append(pids, k)
//			}
//		}
//	}, func(_ interface{}) {}, func() {
//		mutex.RUnlock()
//	})

//	//分页
//	offset := (p - 1) * s
//	end := p * s
//	var _pids []int64
//	if len(pids) > offset && len(pids) < end {
//		_pids = pids[offset:]
//	}
//	if len(pids) > offset && len(pids) >= end {
//		_pids = pids[offset:end]
//	}
//	lives := []LivePerson{}
//	for _, _id := range _pids {
//		lp := lv.Get(_id)
//		if lp != nil {
//			lives = append(lives, *lp)
//		}
//	}
//	pl := &libs.PL{
//		P:     p,
//		S:     s,
//		Total: len(pids),
//		Type:  reflect.TypeOf(lives),
//		List:  lives,
//	}
//	return pl, nil
//}

func (lv *LivePers) Livings() []*LivePerson {
	ids, _ := redis.HKeys(nil, living_personal_list)
	lps := []*LivePerson{}
	for _, idstr := range ids {
		_id, err := strconv.ParseInt(idstr, 10, 64)
		if err == nil {
			lp := lv.Get(_id)
			if lp != nil {
				lps = append(lps, lp)
			}
		}
	}
	return lps
}

//搜索
func (lv *LivePers) Query(words string, p int, s int, match_mode string, filters []search.SearchFilter, filterRanges []search.FilterRangeInt) (int, []*LivePerson) {
	search_config := &search.SearchOptions{
		Host:    search_live_server,
		Port:    search_live_port,
		Timeout: search_live_timeout,
	}
	offset := (p - 1) * s
	search_config.Offset = offset
	search_config.Limit = s
	search_config.Filters = filters
	search_config.FilterRangeInt = filterRanges
	search_config.MaxMatches = 500
	q_engine := search.NewSearcher(search_config)
	ids, total, err := q_engine.LiveQuery(words, match_mode)
	lives := []*LivePerson{}
	if err != nil {
		return 0, lives
	}
	for _, id := range ids {
		live := lv.Get(id)
		if live != nil {
			lives = append(lives, live)
		}
	}
	return total, lives
}

//收藏接口
func (lv LivePers) CompleCollectible(c *collect.Collectible) error {
	_id, _ := strconv.ParseInt(c.RelId, 10, 64)
	if _id <= 0 {
		return fmt.Errorf("参数转换错误")
	}
	channel := lv.Get(_id)
	if channel == nil {
		return fmt.Errorf("不存在指定的个人直播频道")
	}
	c.PreviewContent = channel.Name
	c.PreviewImg = channel.Img
	return nil
}

////////////////////////////////////////////////////////////////////////////////
type LiveOrgs struct{}

func (org *LiveOrgs) ck() string {
	return "mobile_live_channels_all"
}

func (org *LiveOrgs) CreateChannel(channel *LiveChannel) (int64, error) {
	if len(channel.Name) == 0 {
		return 0, errors.New("名称不能为空")
	}
	o := dbs.NewDefaultOrm()
	id, err := o.Insert(channel)
	if err != nil {
		logs.Error(err.Error(), channel)
		return 0, err
	}
	channel.Id = id
	cache := utils.GetCache()
	cache.Delete(org.ck())
	return id, nil
}

func (org *LiveOrgs) UpdateChannel(channel *LiveChannel) error {
	o := dbs.NewDefaultOrm()
	num, err := o.Update(channel)
	if err != nil {
		logs.Error(err.Error(), channel)
		return err
	}
	if num == 0 {
		return errors.New("不存在对应的频道")
	}
	cache := utils.GetCache()
	cache.Delete(org.ck())
	return nil
}

func (org *LiveOrgs) GetChannel(id int64) *LiveChannel {
	lcs := org.GetChannels()
	for _, lc := range lcs {
		if lc.Id == id {
			return lc
		}
	}
	return nil
}

func (org *LiveOrgs) GetChannels() []*LiveChannel {
	cache := utils.GetCache()
	_chs := []*LiveChannel{}
	err := cache.Get(org.ck(), &_chs)
	if err == nil {
		return _chs
	}
	o := dbs.NewDefaultOrm()
	_, err = o.QueryTable(LiveChannel{}).All(&_chs)
	if err != nil {
		return _chs
	}
	cache.Set(org.ck(), _chs, utils.StrToDuration("12h"))
	return _chs
}

func (org *LiveOrgs) GetChannelsByUid(uid int64) []*LiveChannel {
	lcs := org.GetChannels()
	result := []*LiveChannel{}
	for _, lc := range lcs {
		if lc.Uid == uid {
			result = append(result, lc)
		}
	}
	return result
}

//收藏接口
func (org LiveOrgs) CompleCollectible(c *collect.Collectible) error {
	_id, _ := strconv.ParseInt(c.RelId, 10, 64)
	if _id <= 0 {
		return fmt.Errorf("参数转换错误")
	}
	channel := org.GetChannel(_id)
	if channel == nil {
		return fmt.Errorf("不存在指定的机构频道")
	}
	c.PreviewContent = channel.Name
	c.PreviewImg = channel.Img
	return nil
}

///////////////////
//stream
///////////////////
func (org *LiveOrgs) sk(channelId int64) string {
	return fmt.Sprintf("mobile_channel_streams:%d", channelId)
}

func (org *LiveOrgs) sk_bysid(streamId int64) string {
	return fmt.Sprintf("mobile_channel_stream_id:%d", streamId)
}

func (org *LiveOrgs) CreateStream(stream *LiveStream) (int64, error) {
	if stream.ChannelId <= 0 {
		return 0, errors.New("未指定对应的频道")
	}
	channel := org.GetChannel(stream.ChannelId)
	if channel == nil {
		return 0, errors.New("指定的频道不存在")
	}
	if len(stream.StreamUrl) == 0 && len(stream.ReptileUrl) == 0 {
		return 0, errors.New("抓取地址不能为空")
	}
	repSupport, err := reptile.Get_REP_SUPPORT(stream.ReptileUrl)
	if err != nil {
		return 0, errors.New("不支持此地址抓取")
	}
	if stream.LastTime.IsZero() {
		stream.LastTime = time.Now()
	}
	stream.Rep = repSupport

	o := dbs.NewDefaultOrm()
	id, err := o.Insert(stream)
	if err != nil {
		logs.Error(err.Error(), stream)
		return 0, err
	}
	stream.Id = id
	if stream.Default {
		org.SetDefaultStream(stream)
	}
	channel.Childs += 1
	org.UpdateChannel(channel)
	cache := utils.GetCache()
	cache.Delete(org.sk(stream.ChannelId))
	return id, nil
}

func (org *LiveOrgs) SetDefaultStream(defStream *LiveStream) {
	streams := org.GetStreams(defStream.ChannelId)
	o := dbs.NewDefaultOrm()
	cache := utils.GetCache()
	for _, strm := range streams {
		if strm.Id != defStream.Id && strm.Default {
			strm.Default = false
			o.Update(&strm)
			cache.Delete(org.sk(strm.ChannelId))
			cache.Delete(org.sk_bysid(strm.Id))
		}
	}
}

func (org *LiveOrgs) UpdateStream(stream *LiveStream) error {
	repSupport, err := reptile.Get_REP_SUPPORT(stream.ReptileUrl)
	if err != nil {
		return errors.New("不支持此地址抓取")
	}
	stream.Rep = repSupport

	o := dbs.NewDefaultOrm()
	_, err = o.Update(stream)
	if err != nil {
		return err
	}
	cache := utils.GetCache()
	cache.Delete(org.sk(stream.ChannelId))
	cache.Delete(org.sk_bysid(stream.Id))

	if stream.Default {
		org.SetDefaultStream(stream)
	}
	return nil
}

func (org *LiveOrgs) DeleteStream(streamId int64) error {
	o := dbs.NewDefaultOrm()
	stream := &LiveStream{}
	stream.Id = streamId
	err := o.Read(stream)
	if err == orm.ErrNoRows || err == orm.ErrMissPK || err != nil {
		return errors.New("指定删除的流渠道不存在")
	}
	channel := org.GetChannel(stream.ChannelId)
	channel.Childs -= 1
	org.UpdateChannel(channel)
	o.Delete(stream)
	cache := utils.GetCache()
	cache.Delete(org.sk(stream.ChannelId))
	cache.Delete(org.sk_bysid(streamId))
	return nil
}

func (org *LiveOrgs) GetStreams(channelId int64) []*LiveStream {
	cache := utils.GetCache()
	streams := []*LiveStream{}
	err := cache.Get(org.sk(channelId), &streams)
	if err == nil {
		return streams
	}
	o := dbs.NewDefaultOrm()
	_, err = o.QueryTable(LiveStream{}).Filter("pid", channelId).OrderBy("-def").All(&streams)
	if err != nil {
		return streams
	}
	cache.Set(org.sk(channelId), streams, utils.StrToDuration("12h"))
	return streams
}

func (org *LiveOrgs) GetStream(streamId int64) *LiveStream {
	cache := utils.GetCache()
	var stream LiveStream
	err := cache.Get(org.sk_bysid(streamId), &stream)
	if err == nil {
		return &stream
	}
	o := dbs.NewDefaultOrm()
	err = o.QueryTable(LiveStream{}).Filter("id", streamId).One(&stream)
	if err != nil {
		return nil
	}
	cache.Set(org.sk_bysid(streamId), stream, utils.StrToDuration("12h"))
	return &stream
}
