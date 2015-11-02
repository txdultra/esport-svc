package vod

import (
	"dbs"
	"errors"
	"fmt"
	"reflect"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/orm"
	//"github.com/yunge/sphinx"
	"libs"
	"libs/collect"
	"libs/passport"
	"libs/reptile"
	"libs/search"
	"libs/stat"
	"math/rand"
	"strconv"
	"time"
	"utils"
	"utils/cache"
	"utils/ssdb"

	"labix.org/v2/mgo/bson"
)

var vod_view_counter_fmt string = "vod_view_%d"
var vod_download_counter_fmt string = "vod_down_%d"
var vod_view_counter_seed int
var vod_comment_counts_fmt string = "vod_comment_counts:%d"

type Vods struct{}

//重置新消息列表 IEventCounter interface
//func (v Vods) ResetEventCount(uid int64) bool {
//	return message.ResetEventCount(uid, MSG_TYPE_COMMENT)
//}

////IEventCounter interface
//func (v Vods) NewEventCount(uid int64) int {
//	return message.NewEventCount(uid, MSG_TYPE_COMMENT)
//}

func (v *Vods) Create(vod *Video, doRep bool) (int64, error) {
	if len(vod.Url) == 0 {
		return 0, errors.New("Url属性不能为空")
	}
	if vod.Uid <= 0 {
		return 0, errors.New("必须设置视频归属主播")
	}
	rep := reptile.ReptileService(vod.Source)
	if rep == nil {
		return 0, errors.New("不存在可抓取的源")
	}
	if v.GetByKey(vod.Key()) != nil {
		return 0, errors.New("已存在相同的链接")
	}
	o := dbs.NewDefaultOrm()
	o.Begin()
	defer func() {
		if err := recover(); err != nil {
			o.Rollback()
		}
	}()
	if vod.LastUpdateTime.IsZero() {
		vod.LastUpdateTime = time.Now()
	}
	//vod.Seconds = 0 //float32(vs.TotalSeconds)
	vod.Dkey = vod.Key()
	vod.AddTime = time.Now()
	vid, err := o.Insert(vod)
	if err != nil {
		o.Rollback()
		return 0, errors.New("保存数据失败:" + err.Error())
	}
	vcount := VideoCount{
		VideoId: vid,
	}
	_, err = o.Insert(&vcount)
	if err != nil {
		o.Rollback()
		return 0, errors.New("保存视频统计数据失败:" + err.Error())
	}
	//提交事务
	err = o.Commit()
	if err != nil {
		o.Rollback()
		return 0, errors.New("提交事务时发生错误:" + err.Error())
	}
	if doRep {
		go v.Reptile(vid)
	}
	vod.Id = vid
	return vid, nil
}

func (v *Vods) videoCacheKey(id int64) string {
	return fmt.Sprintf("mobile_video_id:%d", id)
}

func (v *Vods) Get(id int64, checkExpried bool) *Video {
	if id <= 0 {
		return nil
	}
	cache := utils.GetCache()
	vod := &Video{}
	err := cache.Get(v.videoCacheKey(id), vod)
	if err != nil {
		o := dbs.NewDefaultOrm()
		vod.Id = id
		err = o.Read(vod)
		if err == orm.ErrNoRows || err == orm.ErrMissPK || err != nil {
			return nil
		}
		v.SetVodToCache(vod)
	}
	//过期检查
	if checkExpried {
		vod.ExpriedAndReptile()
	}
	return vod
}

func (v *Vods) GetByKey(key string) *Video {
	cache := utils.GetCache()
	vod := &Video{}
	err := cache.Get(key, vod)
	if err == nil {
		return vod
	}
	o := dbs.NewDefaultOrm()
	qs := o.QueryTable(vod)
	err = qs.Filter("dkey", key).One(vod)
	if err == orm.ErrNoRows || err != nil {
		return nil
	}
	cache.Set(key, *vod, 96*time.Hour)
	return vod
}

func (v *Vods) Update(vod *Video) error {
	_vod := v.Get(vod.Id, false)
	if _vod == nil {
		return errors.New("原视频不存在")
	}
	if _vod.GameId == 0 && vod.GameId > 0 { //(作假)原视频未归类并且更新视频已归类
		vc := v.GetCount(vod.Id)
		if vc != nil {
			rint := v.rangeViews(vod.GameId)
			vc.Views += rint
			v.UpdateCount(vc)
		}
	}

	rep := reptile.ReptileService(vod.Source)
	if rep == nil {
		return errors.New("不存在可抓取的源")
	}
	if _vod.Dkey != vod.Key() {
		_tmp := v.GetByKey(vod.Key())
		if _tmp != nil {
			return errors.New("已存在相同的视频原地址")
		}
		vod.LastUpdateTime = time.Now().Add(-12 * time.Hour) //设置了不同的url需要重新抓取
	}
	o := dbs.NewDefaultOrm()
	num, err := o.Update(vod)
	if err == nil {
		if num > 0 {
			v.SetVodToCache(vod)
			return nil
		} else {
			return errors.New("未进行更新改动")
		}
	}
	return err
}

func (v *Vods) SetVodToCache(vod *Video) {
	cache := utils.GetCache()
	cache.Set(v.videoCacheKey(vod.Id), *vod, 96*time.Hour)
	cache.Set(vod.Dkey, *vod, 96*time.Hour)
}

func (v *Vods) countCacheKey(videoId int64) string {
	return fmt.Sprintf("mobile_video_counts_videoid:%d", videoId)
}

//统计
func (v *Vods) GetCount(videoId int64) *VideoCount {
	cache := utils.GetCache()
	count := &VideoCount{}
	err := cache.Get(v.countCacheKey(videoId), count)
	if err != nil {
		o := dbs.NewDefaultOrm()
		err = o.QueryTable(count).Filter("vid", videoId).One(count)
		if err == orm.ErrNoRows || err == orm.ErrMultiRows || err != nil {
			return nil
		}
		cache.Set(v.countCacheKey(videoId), *count, 96*time.Hour)
	}
	return count
}

func (v *Vods) rangeViews(gameId int) int {
	rint := 0
	if gameId <= 0 {
		return 0
	}
	switch gameId {
	case 7: //lol
		rint, _ = utils.IntRange(500, 5000)
		break
	case 4: //sc2
		rint, _ = utils.IntRange(300, 4000)
		break
	case 5: //dota2
		rint, _ = utils.IntRange(200, 2000)
		break
	case 3: //lushi
		rint, _ = utils.IntRange(100, 1000)
		break
	case 6: //fengbao
		rint, _ = utils.IntRange(50, 500)
		break
	default: //默认
		rint, _ = utils.IntRange(10, 200)
		break
	}
	return rint
}

func (v *Vods) UpdateCount(count *VideoCount) {
	o := dbs.NewDefaultOrm()
	o.Update(count)
	cache := utils.GetCache()
	cache.Delete(v.countCacheKey(count.VideoId))
}

//views 计数器
func (v Vods) DoC(id int64, n int, term string) {
	switch term {
	case "views":
		stat.NewCountHplper().IncreaseCount(fmt.Sprintf(vod_view_counter_fmt, id), n,
			func(i int) bool {
				j := rand.Intn(vod_view_counter_seed)
				//fmt.Println(i, "====", j)
				if i >= j {
					return true
				}
				return false
			},
			func(i int) bool {
				v := &Vods{}
				count := v.GetCount(id)
				if count != nil {
					count.Views += i
					v.UpdateCount(count)
					return true
				}
				return false
			})
		break
	case "comments":
		key := fmt.Sprintf(vod_comment_counts_fmt, id)
		_c, _ := ssdb.New(use_ssdb_vod_db).Incr(key)
		j := rand.Int63n(10)
		if j%3 == 0 {
			count := v.GetCount(id)
			if count != nil {
				count.Comments = int(_c)
				v.UpdateCount(count)
			}
		}
		break
	case "download":
		stat.NewCountHplper().IncreaseCount(fmt.Sprintf(vod_download_counter_fmt, id), n,
			func(i int) bool {
				j := rand.Intn(30)
				if i >= j {
					return true
				}
				return false
			},
			func(i int) bool {
				v := &Vods{}
				count := v.GetCount(id)
				if count != nil {
					count.Downloads += i
					v.UpdateCount(count)
					return true
				}
				return false
			})
	default:
		break
	}
}

func (v Vods) GetC(id int64, term string) int {
	count := v.GetCount(id)
	if count == nil {
		return 0
	}
	switch term {
	case "views":
		return count.Views
	case "comments":
		return count.Comments
	case "download":
		return count.Downloads
	default:
		return 0
	}
}

func (v *Vods) flvsCacheKey(videoId int64) string {
	return fmt.Sprintf("mobile_video_flvs_videoid:%d", videoId)
}

func (v *Vods) GetPlayFlvs(videoId int64, wait bool) *VideoPlayFlvs {
	if videoId <= 0 {
		return nil
	}
	if wait {
		libs.ChanSyncWait(vod_chan_sync_key(videoId))
	}
	go stat.GetCounter(MOD_NAME).DoC(videoId, 1, "views")
	cache := utils.GetCache()

	var flvs *VideoPlayFlvs
	var bs []byte
	err := cache.Get(v.flvsCacheKey(videoId), &bs)
	if err == nil {
		err = bson.Unmarshal(bs, &flvs)
		if err == nil {
			return flvs
		}
	}

	session, col := dbs.MgoC(MOD_NAME, VOD_FLVS_COLLECTION_NAME)
	defer session.Close()
	err = col.Find(bson.M{"video_id": videoId}).One(&flvs)
	if err == nil {
		com_bs, err := bson.Marshal(*flvs)
		if err == nil {
			cache.Set(v.flvsCacheKey(videoId), com_bs, 1*time.Hour)
		}
		return flvs
	}
	return nil
}

func (v *Vods) GetM3u8Flvs(videoId int64, wait bool) []VideoOpt {
	vod := v.Get(videoId, false)
	if vod == nil {
		return nil
	}
	flvs := v.GetPlayFlvs(videoId, wait)
	if flvs == nil {
		return nil
	}

	convert_modes := make(map[reptile.VOD_STREAM_MODE]reptile.VOD_STREAM_MODE)
	convert_modes[reptile.VOD_STREAM_MODE_MSTD] = reptile.VOD_STREAM_MODE_STANDARD_SP
	convert_modes[reptile.VOD_STREAM_MODE_MHIGH] = reptile.VOD_STREAM_MODE_HIGH_SP
	convert_modes[reptile.VOD_STREAM_MODE_MSUPER] = reptile.VOD_STREAM_MODE_SUPER_SP

	rep := reptile.ReptileService(vod.Source)
	if rep == nil {
		return nil
	}
	vos := []VideoOpt{}
	for mode, toMode := range convert_modes {
		for _, opt := range flvs.OptFlvs {
			if opt.Mode == mode {
				if len(opt.Flvs) > 0 {
					vseg, err := rep.M3u8ToSegs(opt.Flvs[0].Url)
					if err != nil {
						continue
					}
					vo := VideoOpt{
						Size:    opt.Size,
						Mode:    toMode,
						Seconds: opt.Seconds,
						Flvs:    []VideoFlv{},
					}
					for _, seg := range vseg {
						vo.Flvs = append(vo.Flvs, VideoFlv{
							Url:     seg.Url,
							No:      int16(seg.No),
							Size:    seg.Size,
							Seconds: seg.Seconds,
						})
					}
					vo.N = len(vo.Flvs)
					vos = append(vos, vo)
				}
			}
		}
	}
	if len(vos) > 0 {
		for _, vo := range vos {
			existed := false
			for _, opt := range flvs.OptFlvs {
				if opt.Mode == vo.Mode {
					existed = true
				}
			}
			if !existed {
				flvs.OptFlvs = append(flvs.OptFlvs, vo)
			}
		}
		session, col := dbs.MgoC(MOD_NAME, VOD_FLVS_COLLECTION_NAME)
		defer session.Close()
		col.RemoveId(flvs.ID)
		err := col.Insert(*flvs)
		if err == nil {
			com_bs, err := bson.Marshal(*flvs)
			if err == nil {
				cache.Set(v.flvsCacheKey(videoId), com_bs, 1*time.Hour)
			}
		}
	}
	return vos
}

func (v *Vods) VodStreamToVideoPlayFlvs(videoId int64, vs *reptile.VodStreams) *VideoPlayFlvs {
	if vs == nil {
		return nil
	}
	optFlvs := []VideoOpt{}
	if len(vs.Segs) == 0 {
		return nil
	}
	for k, fvs := range vs.Segs {
		size, _ := vs.StreamSizes[k]
		opt := VideoOpt{
			N:       len(fvs),
			Size:    size,
			Mode:    k,
			Seconds: float32(vs.TotalSeconds),
		}
		//快速插入
		flvs := []VideoFlv{}
		for _, flv := range fvs {
			vf := VideoFlv{
				Url:     flv.Url,
				No:      int16(flv.No),
				Size:    flv.Size,
				Seconds: flv.Seconds,
			}
			flvs = append(flvs, vf)
		}
		opt.Flvs = flvs
		optFlvs = append(optFlvs, opt)
	}
	vpf := &VideoPlayFlvs{}
	vpf.ID = bson.NewObjectId()
	vpf.VideoId = videoId
	vpf.OptFlvs = optFlvs
	return vpf
}

func (v *Vods) Reptile(videoId int64) error {
	vod := v.Get(videoId, false)
	if vod == nil {
		return errors.New("不存在此视频")
	}
	if time.Now().Sub(vod.LastUpdateTime).Seconds() < 15 { //防止密集抓取
		return errors.New("15秒内不能同时抓取...")
	}
	rep := reptile.ReptileService(vod.Source)
	vs, err := rep.Reptile(vod.Url)
	if err != nil {
		return errors.New("抓取失败:" + err.Error())
	}
	session, col := dbs.MgoC(MOD_NAME, VOD_FLVS_COLLECTION_NAME)
	defer session.Close()
	col.Remove(bson.M{"video_id": videoId})
	optFlvs := []VideoOpt{}
	if len(vs.Segs) == 0 {
		return errors.New("未抓取到任何数据")
	}
	for k, fvs := range vs.Segs {
		size, _ := vs.StreamSizes[k]
		opt := VideoOpt{
			N:       len(fvs),
			Size:    size,
			Mode:    k,
			Seconds: float32(vs.TotalSeconds),
		}
		//快速插入
		flvs := []VideoFlv{}
		for _, flv := range fvs {
			vf := VideoFlv{
				Url:     flv.Url,
				No:      int16(flv.No),
				Size:    flv.Size,
				Seconds: flv.Seconds,
			}
			flvs = append(flvs, vf)
		}
		opt.Flvs = flvs
		optFlvs = append(optFlvs, opt)
	}
	vpf := VideoPlayFlvs{}
	vpf.ID = bson.NewObjectId()
	vpf.VideoId = videoId
	vpf.OptFlvs = optFlvs
	err = col.Insert(vpf)
	if err != nil {
		return fmt.Errorf("插入错误:%s", err)
	}
	o := dbs.NewDefaultOrm()
	vod.LastUpdateTime = time.Now()
	vod.Seconds = float32(vs.TotalSeconds)
	_, err = o.Update(vod)
	//提交事务
	//err = o.Commit()
	if err != nil {
		return errors.New("提交更新时发生错误:" + err.Error())
	}
	//清空缓存
	cache := utils.GetCache()
	cache.Delete(v.videoCacheKey(vod.Id))
	cache.Delete(vod.Dkey)
	cache.Delete(v.flvsCacheKey(videoId))
	return nil
}

func (v *Vods) Query(words string, p int, s int, match_mode string, filters []search.SearchFilter, filterRanges []search.FilterRangeInt) (int, []*Video) {
	search_config := &search.SearchOptions{
		Host:    search_vod_server,
		Port:    search_vod_port,
		Timeout: search_vod_timeout,
	}

	offset := (p - 1) * s
	search_config.Offset = offset
	search_config.Limit = s
	search_config.Filters = filters
	search_config.FilterRangeInt = filterRanges
	search_config.MaxMatches = 500
	q_engine := search.NewSearcher(search_config)
	ids, total, err := q_engine.VideoQuery(words, match_mode)
	vods := []*Video{}
	if err != nil {
		return 0, vods
	}
	for _, id := range ids {
		vod := v.Get(id, false)
		if vod != nil {
			vods = append(vods, vod)
		}
	}
	return total, vods
}

func (v *Vods) DbQuery(parameters map[string]string, p int, s int) (int, []*Video) {
	o := dbs.NewDefaultOrm()
	query := o.QueryTable(&Video{})
	for k, v := range parameters {
		query = query.Filter(k, v)
	}
	total, err := query.Count()
	if err != nil {
		return 0, []*Video{}
	}
	videos := []*Video{}
	var lists []orm.ParamsList
	offset := (p - 1) * s
	query.OrderBy("-id").Limit(s, offset).ValuesList(&lists, "id")
	for _, m := range lists {
		//id, err := strconv.ParseInt(m[0].(int64), 10, 64)
		id := m[0].(int64)
		if err == nil && id > 0 {
			vod := v.Get(id, false)
			if vod != nil {
				videos = append(videos, vod)
			}
		}
	}
	return int(total), videos
}

//收藏接口
func (v Vods) CompleCollectible(c *collect.Collectible) error {
	_id, _ := strconv.ParseInt(c.RelId, 10, 64)
	if _id <= 0 {
		return fmt.Errorf("参数转换错误")
	}
	vod := v.Get(_id, false)
	if vod == nil {
		return fmt.Errorf("不存在指定的视频")
	}
	c.PreviewContent = vod.Title
	c.PreviewImg = vod.Img
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//专辑 vodIds->id,no

func (v *Vods) CreatePlaylist(p *VideoPlaylist, vodIds map[int64]int) (int64, error) {
	if p.Uid <= 0 {
		return 0, errors.New("专辑必须设置归属主播")
	}
	if len(p.Title) == 0 {
		return 0, errors.New("专辑必须有标题")
	}
	//	if p.Img == 0 {
	//		return 0, errors.New("必须有一张标题图片")
	//	}
	for vid, _ := range vodIds {
		if vid <= 0 {
			return 0, errors.New("视频编号不能小于0")
		}
	}
	p.Vods = len(vodIds)
	if p.PostTime.IsZero() {
		p.PostTime = time.Now()
	}
	o := dbs.NewDefaultOrm()
	id, err := o.Insert(p)
	if err != nil {
		return 0, errors.New("创建视频专辑失败，请联系管理员")
	}
	p.Id = id
	v.UpdatePlaylistVods(p.Id, vodIds)
	return id, nil
}

func (v *Vods) vplCacheKey(plid int64) string {
	return fmt.Sprintf("mobile_video_playlist_id:%d", plid)
}

func (v *Vods) GetPlaylist(id int64) *VideoPlaylist {
	vkey := v.vplCacheKey(id)
	tmp := &VideoPlaylist{}
	cache := utils.GetCache()
	err := cache.Get(vkey, tmp)
	if err == nil {
		return tmp
	}
	vpl := &VideoPlaylist{}
	vpl.Id = id
	o := dbs.NewDefaultOrm()
	err = o.Read(vpl)
	if err == orm.ErrNoRows || err == orm.ErrMissPK || err != nil {
		return nil
	}
	cache.Set(vkey, vpl, 5*time.Hour)
	return vpl
}

func (v *Vods) GetPlaylistsForAdmin(p int, s int) (int, []*VideoPlaylist) {
	if p <= 0 {
		p = 1
	}
	if s <= 0 {
		s = 20
	}
	offset := (p - 1) * s
	var lst []*VideoPlaylist
	o := dbs.NewDefaultOrm()
	query := o.QueryTable(&VideoPlaylist{})
	total, _ := query.Count()
	query.OrderBy("-post_time").Limit(s, offset).All(&lst)
	return int(total), lst
}

func (v *Vods) vplvodsCacheKey(plsId int64) string {
	return fmt.Sprintf("mobile_video_playlist_vods_vid:%d", plsId)
}

func (v *Vods) GetPlaylistVods(plsId int64, p int, s int) (int, []*Video) {
	ckey := v.vplvodsCacheKey(plsId)
	pl := v.GetPlaylist(plsId)
	if pl == nil {
		return 0, []*Video{}
	}
	if p <= 0 {
		p = 1
	}
	if s <= 0 {
		s = 20
	}
	offset := (p - 1) * s
	end := p * s
	cclient := ssdb.New(use_ssdb_vod_db)
	objs, _ := cclient.Zrange(ckey, offset, end, reflect.TypeOf(int64(0)))
	videos := []*Video{}
	for _, obj := range objs {
		_id := *(obj.(*int64))
		vod := v.Get(_id, false)
		if vod != nil {
			videos = append(videos, vod)
		}
	}
	return pl.Vods, videos
}

func (v *Vods) GetPlaylistPlvsForAdmin(plsId int64, p int, s int) (int, []*VideoPlaylistVod) {
	offset := (p - 1) * s
	o := dbs.NewDefaultOrm()
	query := o.QueryTable(&VideoPlaylistVod{})
	total, _ := query.Count()
	var lst []*VideoPlaylistVod
	query.Filter("pid", plsId).OrderBy("no").Limit(s, offset).All(&lst)
	return int(total), lst
}

func (v *Vods) UpdatePlaylistVods(plsId int64, vodIds map[int64]int) error {
	pl := v.GetPlaylist(plsId)
	if pl == nil {
		return fmt.Errorf("专辑不存在")
	}
	o := dbs.NewDefaultOrm()
	o.QueryTable(&VideoPlaylistVod{}).Filter("pid", plsId).Delete()

	cclient := ssdb.New(use_ssdb_vod_db)
	cclient.Zclear(v.vplvodsCacheKey(plsId)) //删除原列表

	inserts := []VideoPlaylistVod{}
	for vid, no := range vodIds {
		inserts = append(inserts, VideoPlaylistVod{
			PlaylistId: plsId,
			VideoId:    vid,
			No:         no,
		})
	}
	rows, err := o.InsertMulti(50, inserts)

	pl.Vods = int(rows)
	o.Update(pl)

	ks := []interface{}{}
	vs := []int64{}
	for k, v := range vodIds {
		ks = append(ks, k)
		vs = append(vs, int64(v))
	}
	cclient.MultiZadd(v.vplvodsCacheKey(plsId), ks, vs)

	cache := utils.GetCache()
	cache.Delete(v.vplCacheKey(plsId))
	return err
}

//func (v *Vods) UpdatePLsVodNos(plsId int64, vodIds map[int64]int) error {
//	o := dbs.NewDefaultOrm()
//	vpvs := []*VideoPlaylistVod{}
//	qs := o.QueryTable(&VideoPlaylistVod{})
//	qs.Filter("pid", plsId).All(&vpvs)
//	get := func(vid int64) (bool, *VideoPlaylistVod) {
//		for _, vp := range vpvs {
//			if vp.VideoId == vid {
//				return true, vp
//			}
//		}
//		return false, nil
//	}
//	o.Begin()
//	for vid, no := range vodIds {
//		ok, p := get(vid)
//		if ok {
//			if p.No != no {
//				p.No = no
//				o.Update(p)
//			}
//		}
//	}
//	err := o.Commit()
//	if err != nil {
//		o.Rollback()
//		return errors.New("提交事务时发生错误:" + err.Error())
//	}
//	cache := utils.GetCache()
//	cache.Delete(v.vplvodsCacheKey(plsId))
//	return nil
//}

//func (v *Vods) RemovePlsVod(plsId int64, videoId int64) error {
//	o := dbs.NewDefaultOrm()
//	pls := &VideoPlaylist{}
//	pls.Id = plsId
//	err := o.Read(pls)
//	if err == orm.ErrNoRows || err == orm.ErrMissPK || err != nil {
//		return errors.New("指定的专辑不存在")
//	}

//	qs := o.QueryTable(&VideoPlaylistVod{})
//	num, err := qs.Filter("pid", plsId).Filter("vid", videoId).Delete()
//	if err != nil {
//		return err
//	}
//	if num == 1 {
//		pls.Vods -= 1
//		o.Update(pls)
//		cache := utils.GetCache()
//		cache.Delete(v.vplvodsCacheKey(plsId))
//		cache.Delete(v.vplCacheKey(plsId))
//	}
//	return nil
//}

//func (v *Vods) AppenedPlsVods(plsId int64, vodIds map[int64]int) error {
//	if len(vodIds) == 0 {
//		return nil
//	}
//	o := dbs.NewDefaultOrm()
//	pls := &VideoPlaylist{}
//	pls.Id = plsId
//	err := o.Read(pls)
//	if err == orm.ErrNoRows || err == orm.ErrMissPK || err != nil {
//		return errors.New("指定的专辑不存在")
//	}

//	vpvs := []*VideoPlaylistVod{}
//	qs := o.QueryTable(&VideoPlaylistVod{})
//	qs.Filter("pid", plsId).OrderBy("-no").All(&vpvs)
//	o.Begin()
//	has := func(vid int64) bool {
//		for _, vp := range vpvs {
//			if vp.VideoId == vid {
//				return true
//			}
//		}
//		return false
//	}
//	apds := 0
//	for vid, no := range vodIds {
//		ok := has(vid)
//		if !ok {
//			vp := &VideoPlaylistVod{
//				PlaylistId: plsId,
//				VideoId:    vid,
//				No:         no,
//			}
//			o.Insert(vp)
//			apds++
//		}
//	}
//	pls.Vods += apds
//	o.Update(pls)
//	err = o.Commit()
//	if err != nil {
//		o.Rollback()
//		return errors.New("提交事务时发生错误:" + err.Error())
//	}
//	cache := utils.GetCache()
//	cache.Delete(v.vplvodsCacheKey(plsId))
//	cache.Delete(v.vplCacheKey(plsId))
//	return nil
//}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//个人空间抓取
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
var reptile_uc_task_lock bool

type VodUcenterReptile struct{}

func (v *VodUcenterReptile) Create(uc VodUcenter) (int64, error) {
	if len(uc.SiteUrl) == 0 {
		return 0, errors.New("SiteUrl不能为空")
	}
	if uc.Uid <= 0 {
		return 0, errors.New("必须设置视频归属主播")
	}
	o := dbs.NewDefaultOrm()
	qs := o.QueryTable(&VodUcenter{})
	cs, err := qs.Filter("uid", uc.Uid).Filter("source", uc.Source).Count()
	if err != nil || cs > 0 {
		return 0, errors.New("已存在相同的空间地址")
	}
	uc.LastTime = time.Now().Add(-12 * time.Hour)
	id, err := o.Insert(&uc)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (v *VodUcenterReptile) Get(id int64) *VodUcenter {
	if id <= 0 {
		return nil
	}
	var uc VodUcenter
	o := dbs.NewDefaultOrm()
	qs := o.QueryTable(&VodUcenter{})
	err := qs.Filter("id", uc.Uid).One(&uc)
	if err != nil {
		return nil
	}
	return &uc
}

func (v *VodUcenterReptile) Gets(page int, size int) (int, []VodUcenter) {
	o := dbs.NewDefaultOrm()
	offset := (page - 1) * size
	var list []VodUcenter
	total, _ := o.QueryTable(&VodUcenter{}).OrderBy("-create_time").Count()
	o.QueryTable(&VodUcenter{}).OrderBy("-create_time").Limit(size, offset).All(&list)
	return int(total), list
}

func (v *VodUcenterReptile) Delete(id int64) error {
	o := dbs.NewDefaultOrm()
	qs := o.QueryTable(&VodUcenter{})
	_, err := qs.Filter("id", id).Delete()
	return err
}

func (v *VodUcenterReptile) Update(uc VodUcenter) error {
	o := dbs.NewDefaultOrm()
	ups, err := o.Update(&uc)
	if err != nil {
		return err
	}
	if ups <= 0 {
		return errors.New("不存在需要更新的对象")
	}
	return nil
}

func (v *VodUcenterReptile) ScanAllReptile(uid int64) error {
	o := dbs.NewDefaultOrm()
	uc := VodUcenter{}
	err := o.QueryTable(&VodUcenter{}).Filter("uid", uid).One(&uc)
	if err == orm.ErrMultiRows {
		return fmt.Errorf("存在多条记录")
	}
	if err == orm.ErrNoRows {
		return fmt.Errorf("没有此uid记录")
	}
	uc.ScanAll = true
	ups, err := o.Update(&uc)
	if err != nil {
		return err
	}
	if ups <= 0 {
		return errors.New("无法设置更新对象")
	}
	return nil
}

//切换到新用户
func (v *VodUcenterReptile) ChangeToUser(ucid int64, newUid int64, syncOldVods bool) error {
	o := dbs.NewDefaultOrm()
	var uc VodUcenter
	err := o.QueryTable(&VodUcenter{}).Filter("id", ucid).One(&uc)
	if err == orm.ErrNoRows {
		return fmt.Errorf("没有此记录")
	}

	mps := passport.NewMemberProvider()
	member := mps.Get(newUid)
	if member == nil {
		return fmt.Errorf("新的用户不存在")
	}
	if reptile_uc_task_lock {
		return fmt.Errorf("抓取正在执行，请稍后再试")
	}
	reptile_uc_task_lock = true
	defer func() {
		reptile_uc_task_lock = false
	}()
	if uc.Uid == newUid {
		return fmt.Errorf("抓取空间的用户不能和新用户是同一个")
	}
	existed := o.QueryTable(&VodUcenter{}).Filter("uid", newUid).Exist()
	if existed {
		return fmt.Errorf("新用户已存在抓取列表中")
	}
	orginalUid := uc.Uid
	uc.Uid = newUid
	_, err = o.Update(&uc)
	if err != nil {
		return err
	}

	if syncOldVods {
		size := 300
		for {
			var vods []*Video
			_v := &Video{}
			o.QueryTable(_v).Filter("uid", orginalUid).OrderBy("add_time").Limit(size).All(&vods)
			if len(vods) == 0 {
				break
			}
			_vods := &Vods{}
			p, _ := o.Raw("UPDATE " + _v.TableName() + " SET uid = ? WHERE id = ?").Prepare()
			for _, v := range vods {
				_, err = p.Exec(newUid, v.Id)
				if err == nil {
					v.Uid = newUid
					_vods.SetVodToCache(v)
				}
			}
			p.Close() // 别忘记关闭 statement
		}
	}

	return nil
}

func (v *VodUcenterReptile) ReptileTask() {
	if reptile_uc_task_lock {
		return
	}
	reptile_uc_task_lock = true
	defer func() {
		reptile_uc_task_lock = false
	}()
	reptile_ucenter_intervals, _ := beego.AppConfig.Int("reptile.ucenter.intervals")
	if reptile_ucenter_intervals <= 0 {
		reptile_ucenter_intervals = 5
	}
	o := dbs.NewDefaultOrm()
	qs := o.QueryTable(&VodUcenter{})
	vcs := []VodUcenter{}
	qs.OrderBy("last_time").Limit(reptile_ucenter_intervals).All(&vcs)
	for _, vc := range vcs {
		task := func() {
			save := func(rvds []*reptile.RVodData) bool {
				if len(rvds) == 0 {
					return false
				}
				for _, vd := range rvds {
					vod := &Video{
						Title:          vd.Title,
						Url:            vd.PlayUrl,
						LastUpdateTime: vd.PostTime,
						PostTime:       vd.PostTime,
						Mt:             false,
						Source:         vc.Source,
						Uid:            vc.Uid,
					}
					vss := &Vods{}
					_v := vss.GetByKey(vod.Key())
					if _v != nil { //再下面的视频已被抓取过，无需继续循环,除非设定全局扫描
						if vc.ScanAll || _v.Mt { //已有手动添加的略过
							continue
						} else {
							return false
						}
					}
					if len(vd.Image) > 0 {
						req := httplib.Get(vd.Image)
						req.SetTimeout(3*time.Minute, 3*time.Minute)
						img_data, err := req.Bytes()
						if err == nil {
							file_name := fmt.Sprintf("%s.%s", vod.Key(), vd.ImgExt)
							file_storage := libs.NewFileStorage()
							node, err := file_storage.SaveFile(img_data, file_name, 0)
							if err == nil {
								vod.Img = node.FileId
							}
						}
					}
					vss.Create(vod, false)
				}
				return true
			}
			rep_service := reptile.ReptileUcService(vc.Source)
			if rep_service != nil {
				rep_service.Reptile(vc.SiteUrl, save)
			}
		}
		finally := func() {
			vc.LastTime = time.Now()
			vc.ScanAll = false
			o.Update(&vc)
		}
		utils.Try(task, func(interface{}) {}, finally)
	}
}
