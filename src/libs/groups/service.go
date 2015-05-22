package groups

import (
	"bytes"
	"dbs"
	"encoding/json"
	"fmt"
	"image"
	"libs"
	credit_client "libs/credits/client"
	credit_proxy "libs/credits/proxy"
	//"libs/dlock"
	"libs/message"
	"libs/passport"
	"libs/search"
	"logs"
	"reflect"
	"strconv"
	"sync"
	"time"
	"utils"
	"utils/cache"
	"utils/ssdb"

	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/orm"
	"github.com/disintegration/imaging"
)

type BaseService struct {
	cfg *GroupCfg
}

////////////////////////////////////////////////////////////////////////////////
//组设定配置
////////////////////////////////////////////////////////////////////////////////
var cfgLock *sync.RWMutex = new(sync.RWMutex)

var default_cfg *GroupCfg

func getCfgCacheKey(id int64) string {
	return fmt.Sprintf("group.config_%d", id)
}

func GetGroupCfg(id int64) *GroupCfg {
	ckey := getCfgCacheKey(id)
	cfg := GroupCfg{}
	cache := utils.GetCache()
	err := cache.Get(ckey, &cfg)
	if err == nil {
		return &cfg
	}
	cfg.Id = id
	o := dbs.NewOrm(db_aliasname)
	err = o.Read(&cfg)
	if err != nil {
		return nil
	}
	cache.Add(ckey, cfg, 24*time.Hour)
	return &cfg
}

func GetDefaultCfg() *GroupCfg {
	cfgLock.RLock()
	defer cfgLock.RUnlock()
	if default_cfg == nil {
		default_cfg = GetGroupCfg(int64(group_setting_id))
	}
	return default_cfg
}

func ResetDefaultCfg() {
	cfgLock.Lock()
	defer cfgLock.Unlock()
	default_cfg = nil
}

func UpdateGroupCfg(cfg *GroupCfg) error {
	o := dbs.NewOrm(db_aliasname)
	_, err := o.Update(cfg)
	if err != nil {
		return err
	}
	data, err := json.Marshal(cfg)
	if err == nil {
		watcher.Write(watcher_path, data)
	}
	cache.Replace(getCfgCacheKey(cfg.Id), *cfg, 24*time.Hour)
	return nil
}

////////////////////////////////////////////////////////////////////////////////
//公用方法
////////////////////////////////////////////////////////////////////////////////
func savePostTableToDb(addTblId int, addTblName string, ptrTable interface{}) {
	o := dbs.NewOrm(db_aliasname)
	create_tbl_sql := fmt.Sprintf(`select TABLE_NAME from INFORMATION_SCHEMA.TABLES where TABLE_SCHEMA='%s' and TABLE_NAME='%s'`, db_name, addTblName)
	var maps []orm.Params
	o.Raw(create_tbl_sql).Values(&maps)
	//已存在表
	o.ReadOrCreate(ptrTable, "Id")
}

var file_storage = libs.NewWeedFsFileStorage()

func fileData(fid int64) ([]byte, error) {
	fileUrl := file_storage.GetPrivateFileUrl(fid)
	request := httplib.Get(fileUrl)
	request.SetTimeout(1*time.Minute, 1*time.Minute)
	data, err := request.Bytes()
	if err != nil {
		logs.Errorf("group pic get file data fail:%s", err.Error())
		return nil, err
	}
	return data, nil
}

func imgResize(data []byte, file *libs.File, size int) (int64, error) {
	if size <= 0 {
		return file.Id, nil
	}
	if file.Width <= size && file.Height <= size {
		return file.Id, nil
	}
	if file.Width >= file.Height && file.Width > size {
		buf := bytes.NewBuffer(data)
		srcImg, _, err := image.Decode(buf)
		if err != nil {
			logs.Errorf("group pic thumbnail by width's ratio data convert to image fail:%s", err.Error())
			return 0, err
		}
		var dstImage image.Image
		dstImage = imaging.Resize(srcImg, size, 0, imaging.Lanczos)
		fileData, err := utils.ImageToBytes(dstImage, file.OriginalName)
		if err != nil {
			logs.Errorf("group pic thumbnail by width's ratio image to bytes fail:%s", err.Error())
			return 0, err
		}
		ratio := float32(file.Width) / float32(size)
		resizeName := fmt.Sprintf("%s-r%dx%d.%s", utils.FileName(file.OriginalName), size, int(float32(file.Height)*ratio), file.ExtName)
		node, err := file_storage.SaveFile(fileData, resizeName, file.Id)
		if err != nil {
			logs.Errorf("group pic thumbnail by width's ratio save image fail:%s", err.Error())
			return 0, err
		}
		return node.FileId, err
	}
	if file.Width < file.Height && file.Height > size {
		buf := bytes.NewBuffer(data)
		srcImg, _, err := image.Decode(buf)
		if err != nil {
			logs.Errorf("group pic thumbnail by height's ratio data convert to image fail:%s", err.Error())
			return 0, err
		}
		var dstImage image.Image
		dstImage = imaging.Resize(srcImg, 0, size, imaging.Lanczos)
		fileData, err := utils.ImageToBytes(dstImage, file.OriginalName)
		if err != nil {
			logs.Errorf("group pic thumbnail by height's ratio image to bytes fail:%s", err.Error())
			return 0, err
		}
		ratio := float32(file.Height) / float32(size)
		resizeName := fmt.Sprintf("%s-r%dx%d.%s", utils.FileName(file.OriginalName), int(float32(file.Width)*ratio), size, file.ExtName)
		node, err := file_storage.SaveFile(fileData, resizeName, file.Id)
		if err != nil {
			logs.Errorf("group pic thumbnail by height's ratio save image fail:%s", err.Error())
			return 0, err
		}
		return node.FileId, err
	}
	return file.Id, nil
}

////////////////////////////////////////////////////////////////////////////////
//组服务
////////////////////////////////////////////////////////////////////////////////
const (
	group_name_sets             = "group.name_sets"
	group_member_joined_sortset = "group.member_%d_joined"
	group_joined_member_sortset = "group.joined_%d_uids"
	group_invited_set           = "group.invited_%d_%d_uids"
	group_index_name            = "group_idx"
	count_action_set            = "group_count_incr_%s"
)

type GP_PROPERTY string

const (
	GP_PROPERTY_MEMBERS GP_PROPERTY = "members"
	GP_PROPERTY_THREADS GP_PROPERTY = "threads"
)

var GP_PROPERTY_ALL = []GP_PROPERTY{
	GP_PROPERTY_MEMBERS,
	GP_PROPERTY_THREADS,
}

type GP_SEARCH_SORT int

const (
	GP_SEARCH_SORT_DEFAULT  GP_SEARCH_SORT = 1
	GP_SEARCH_SORT_USERS    GP_SEARCH_SORT = 2
	GP_SEARCH_SORT_VITALITY GP_SEARCH_SORT = 3
	GP_SEARCH_SORT_OFFICIAL GP_SEARCH_SORT = 4
	GP_SEARCH_SORT_ENDTIME  GP_SEARCH_SORT = 5
)

var tbl_mutex *sync.Mutex = new(sync.Mutex)
var mgTbls map[int]string = make(map[int]string)
var gmTbls map[int]string = make(map[int]string)

func NewGroupService(cfg *GroupCfg) *GroupService {
	service := &GroupService{}
	service.cfg = cfg
	return service
}

type GroupService struct {
	BaseService
}

//重置新消息列表 IEventCounter interface
func (s GroupService) ResetEventCount(uid int64) bool {
	return message.ResetEventCount(uid, MSG_TYPE_MESSAGE)
}

//IEventCounter interface
func (s GroupService) NewEventCount(uid int64) int {
	return message.NewEventCount(uid, MSG_TYPE_MESSAGE)
}

func (s *GroupService) GetCacheKey(id int64) string {
	return fmt.Sprintf("group_model_%d", id)
}

func (s *GroupService) AllowMaxCreateLimits(uid int64) int {
	ps := passport.NewMemberProvider()
	member := ps.Get(uid)
	if member == nil {
		return 0
	}
	max_limit := s.cfg.CreateGroupMaxCount
	if member.Certified {
		max_limit = s.cfg.CreateGroupCertifiedMaxCount
	}
	return max_limit
}

func (s *GroupService) CheckMemberNewGroupPass(uid int64, belong GROUP_BELONG) error {
	max_limit := s.AllowMaxCreateLimits(uid)
	//验证创建个数限制
	o := dbs.NewOrm(db_aliasname)
	nums, _ := o.QueryTable(&Group{}).Filter("uid", uid).Exclude("status", GROUP_STATUS_CLOSED).Count()
	if int(nums) >= max_limit {
		return fmt.Errorf("已超过可以创建小组的数量限制")
	}

	//验证积分
	if belong == GROUP_BELONG_MEMBER {
		var needPoint int64 = 0 //需要积分
		needPoint = s.cfg.CreateGroupBasePoint
		client, transport, err := credit_client.NewClient(credit_service_host)
		if err != nil {
			return fmt.Errorf("查询积分失败:001")
		}
		defer func() {
			if transport != nil {
				transport.Close()
			}
		}()
		points, err := client.GetCredit(uid)
		if err != nil {
			return fmt.Errorf("查询积分失败:002")
		}
		if needPoint > points {
			return fmt.Errorf("积分不足")
		}
	}
	return nil
}

func (s *GroupService) VerifyNewGroup(group *Group) error {
	if s.cfg == nil {
		return fmt.Errorf("未设置配置对象")
	}
	nameRunes := []rune(group.Name)
	if len(nameRunes) > s.cfg.GroupNameLen {
		return fmt.Errorf(fmt.Sprintf("组名称不能大于%d个字符", s.cfg.GroupNameLen))
	}
	_name := utils.CensorWords(group.Name)
	if _name != group.Name {
		return fmt.Errorf("小组名称包含屏蔽字")
	}
	descRunes := []rune(group.Description)
	if len(descRunes) > s.cfg.GroupDescMaxLen {
		return fmt.Errorf(fmt.Sprintf("组描述内容不能大于%d个字符", s.cfg.GroupDescMaxLen))
	}
	if len(descRunes) <= s.cfg.GroupDescMinLen {
		return fmt.Errorf(fmt.Sprintf("组描述内容不能小于%d个字符", s.cfg.GroupDescMinLen))
	}
	if group.Uid <= 0 {
		return fmt.Errorf("未设置创建人")
	}
	if len(group.GameIds) == 0 {
		return fmt.Errorf("未选择游戏分类")
	}
	//验证是否同名
	has, _ := ssdb.New(use_ssdb_group_db).Hexists(group_name_sets, group.Name)
	if has {
		return fmt.Errorf("组名称已存在")
	}
	o := dbs.NewOrm(db_aliasname)
	nums, _ := o.QueryTable(&Group{}).Filter("groupname", group.Name).Count()
	if int(nums) > 0 {
		return fmt.Errorf("组名称已存在")
	}

	//认证用户创建个数不同
	max_limit := s.AllowMaxCreateLimits(group.Uid)
	//验证创建个数限制
	nums, _ = o.QueryTable(&Group{}).Filter("uid", group.Uid).Exclude("status", GROUP_STATUS_CLOSED).Count()
	if int(nums) >= max_limit {
		return fmt.Errorf("已超过可以创建小组的数量限制")
	}

	//验证积分
	if group.Belong == GROUP_BELONG_MEMBER {
		var needPoint int64 = 0 //需要积分
		needPoint = s.cfg.CreateGroupBasePoint
		client, transport, err := credit_client.NewClient(credit_service_host)
		if err != nil {
			return fmt.Errorf("查询积分失败:001")
		}
		defer func() {
			if transport != nil {
				transport.Close()
			}
		}()
		points, err := client.GetCredit(group.Uid)
		if err != nil {
			return fmt.Errorf("查询积分失败:002")
		}
		if needPoint > points {
			return fmt.Errorf("积分不足")
		}
	}
	return nil
}

//新建加入组分表
func (s *GroupService) GetGMTable(uid int64) (int, string, error) {
	str_uid := strconv.FormatInt(uid, 10)
	str_pfx := str_uid[:1] //1位
	i_tag, err := strconv.Atoi(str_pfx)
	if err != nil {
		return 0, "", err
	}
	if name, ok := gmTbls[i_tag]; ok {
		return i_tag, name, nil
	}
	tbl_mutex.Lock()
	defer tbl_mutex.Unlock()
	if name, ok := gmTbls[i_tag]; ok {
		return i_tag, name, nil
	}
	tbl := fmt.Sprintf("group_members_%d", i_tag)

	o := dbs.NewOrm(db_aliasname)
	create_tbl_sql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(groupid int(11) NOT NULL,
	  uid int(11) NOT NULL,
	  ts int(11) NOT NULL,
	  PRIMARY KEY (groupid,uid)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`, tbl)
	o.Raw(create_tbl_sql).Exec()
	gmTbls[i_tag] = tbl

	savePostTableToDb(i_tag, tbl, &GroupMemberTable{
		Id:      int64(i_tag),
		TblName: tbl,
		Ts:      time.Now().Unix(),
	})
	return i_tag, tbl, nil
}

//新建用户加入的组分表
func (s *GroupService) GetMGTable(uid int64) (int, string, error) {
	i_tag := 0
	var err error = nil
	if uid < 10 {
		i_tag = int(uid)
	} else {
		str_uid := strconv.FormatInt(uid, 10)
		str_pfx := str_uid[:2] //2位
		i_tag, err = strconv.Atoi(str_pfx)
	}
	if err != nil {
		return 0, "", err
	}
	if name, ok := mgTbls[i_tag]; ok {
		return i_tag, name, nil
	}
	tbl_mutex.Lock()
	defer tbl_mutex.Unlock()
	if name, ok := mgTbls[i_tag]; ok {
		return i_tag, name, nil
	}
	tbl := fmt.Sprintf("member_groups_%d", i_tag)

	o := dbs.NewOrm(db_aliasname)
	create_tbl_sql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(uid int(11) NOT NULL,
	  groupid int(11) NOT NULL,
	  ts int(11) NOT NULL,
	  PRIMARY KEY (uid,groupid)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`, tbl)
	o.Raw(create_tbl_sql).Exec()
	mgTbls[i_tag] = tbl

	savePostTableToDb(i_tag, tbl, &MemberGroupTable{
		Id:      int64(i_tag),
		TblName: tbl,
		Ts:      time.Now().Unix(),
	})
	return i_tag, tbl, nil
}

//设置新建组分表属性
func (s *GroupService) setGroupTableId(group *Group) error {
	group.ThreadTableId = 1 //默认,数据量大后分表
	tblId, _, err := s.GetGMTable(group.Uid)
	if tblId > 0 {
		group.MembersTableId = tblId
	}
	return err
}

func (s *GroupService) thumbnailImgResize(group *Group) error {
	//图片压缩
	file := file_storage.GetFile(group.BgImg)
	if file == nil {
		return fmt.Errorf("背景图片不存在")
	}
	data, err := fileData(group.BgImg)
	if err != nil {
		fmt.Println("---------------------------------", err)
		logs.Errorf("group view pic err:%s", err.Error())
		return fmt.Errorf("获取背景图片失败")
	}
	fid, err := imgResize(data, file, group_pic_thumbnail_w)
	if err != nil {
		fmt.Println("---------------------------------", err)
		return fmt.Errorf("压缩图片失败")
	}
	group.Img = fid //设置小图
	return nil
}

func (s *GroupService) setSearchKeyword(group *Group) error {
	gids := group.GameIDs()
	if len(gids) == 0 {
		return nil
	}
	str := ""
	bas := &libs.Bas{}
	for _, gid := range gids {
		game := bas.GetGame(gid)
		if game != nil {
			str += fmt.Sprintf("%s,%s,", game.Name, game.En)
		}
	}
	group.SearchKeyword = str
	return nil
}

func (s *GroupService) Create(group *Group) error {
	err := s.VerifyNewGroup(group)
	if err != nil {
		return err
	}
	group.CreateTime = time.Now().Unix()
	group.Country = ""
	group.City = ""
	if s.cfg == nil {
		return fmt.Errorf("未设置配置对象")
	}
	switch group.Belong {
	case GROUP_BELONG_MEMBER:
		group.Status = GROUP_STATUS_RECRUITING
		group.StartTime = time.Now().Unix()
		addHours := time.Duration(s.cfg.CreateGroupRecruitDay * 24)
		group.EndTime = time.Now().Add(addHours * time.Hour).Unix()
		group.MinUsers = s.cfg.CreateGroupMinUsers
	case GROUP_BELONG_OFFICIAL:
		group.Status = GROUP_STATUS_OPENING
		group.MinUsers = 0
	}

	//设置分表
	err = s.setGroupTableId(group)
	if err != nil {
		return fmt.Errorf("建组失败:001")
	}

	//图片压缩
	if group.BgImg > 0 {
		err = s.thumbnailImgResize(group)
		if err != nil {
			return err
		}
	}

	//设置搜索关键字
	s.setSearchKeyword(group)

	//扣除积分
	client, transport, err := credit_client.NewClient(credit_service_host)
	if err != nil {
		return fmt.Errorf("扣除积分失败:001")
	}
	defer func() {
		if transport != nil {
			transport.Close()
		}
	}()
	cr, err := client.Do(&credit_proxy.OperationCreditParameter{
		Uid:    group.Uid,
		Points: s.cfg.CreateGroupBasePoint,
		Action: credit_proxy.OPERATION_ACTOIN_LOCKDECR,
	})
	if err != nil {
		return fmt.Errorf("扣除积分失败:002")
	}
	//记录积分订单号
	group.OrderNo = cr.No

	o := dbs.NewOrm(db_aliasname)
	id, err := o.Insert(group)
	if err != nil {
		//退回积分
		client.Do(&credit_proxy.OperationCreditParameter{
			No:     cr.No,
			Action: credit_proxy.OPERATION_ACTOIN_UNLOCK,
		})
		logs.Errorf("建组失败:%+v", err)
		return fmt.Errorf("建组失败:002")
	}
	group.Id = id

	//加入定时任务
	if group.Belong == GROUP_BELONG_MEMBER {
		go tjSetJob(group)
	}

	//加入组名hash集合
	ssdb.New(use_ssdb_group_db).Hset(group_name_sets, group.Name, group.Id)

	cache := utils.GetCache()
	cache.Add(s.GetCacheKey(id), *group, 1*time.Hour)

	//更新我的缓存
	cache.Delete(s.getMyGroupIdsCacheKey(group.Uid))
	//增加计数
	NewMemberCountService().ActionCountProperty([]MC_PROPERTY{MC_PROPERTY_GROUPS}, []COUNT_ACTION{COUNT_ACTION_INCR}, group.Uid)

	//邀请好友
	go s.Invite(group.Uid, group.Id, group.InviteUids)

	return nil
}

func (s *GroupService) Delete(groupId int64) error {
	group := s.Get(groupId)
	if group == nil {
		return fmt.Errorf("组已被删除")
	}

	gjmkey := fmt.Sprintf(group_joined_member_sortset, groupId)
	//清理用户加入组的记录
	vals, _ := ssdb.New(use_ssdb_group_db).Zrange(gjmkey, 0, -1, reflect.TypeOf(int64(0)))
	for _, val := range vals {
		uid := *(val.(*int64))
		if uid <= 0 {
			continue
		}
		s.Exit(uid, groupId) //强制用户离开
	}

	o := dbs.NewOrm(db_aliasname)
	_, err := o.Delete(group)
	if err != nil {
		return fmt.Errorf("删除失败:001")
	}

	cache := utils.GetCache()
	group.IsDeleted = true
	cache.Set(s.GetCacheKey(groupId), *group, 24*time.Hour)
	cache.Delete(s.getMyGroupIdsCacheKey(group.Uid))

	//删除计数
	NewMemberCountService().ActionCountProperty([]MC_PROPERTY{MC_PROPERTY_GROUPS}, []COUNT_ACTION{COUNT_ACTION_DECR}, group.Uid)

	group.Status = GROUP_STATUS_CLOSED
	go s.UpdateSearchEngineAttrs([]*Group{group}, []string{"status"})

	return nil
}

func (s *GroupService) Close(groupId int64) error {
	group := s.Get(groupId)
	if group == nil {
		return fmt.Errorf("组不存在")
	}
	if group.Status == GROUP_STATUS_CLOSED {
		return fmt.Errorf("组已被关闭")
	}
	o := dbs.NewOrm(db_aliasname)
	group.Status = GROUP_STATUS_CLOSED
	group.StartTime = 0
	group.EndTime = 0
	o.Update(group)
	cache := utils.GetCache()
	cache.Replace(s.GetCacheKey(groupId), *group, 1*time.Hour)
	cache.Delete(s.getMyGroupIdsCacheKey(group.Uid))
	//减少计数
	NewMemberCountService().ActionCountProperty([]MC_PROPERTY{MC_PROPERTY_GROUPS}, []COUNT_ACTION{COUNT_ACTION_INCR}, group.Uid)
	go s.UpdateBaseSearchEngineAttrs([]*Group{group})
	return nil
}

func (s *GroupService) Update(group *Group) error {
	if s.cfg == nil {
		return fmt.Errorf("未设置配置对象")
	}
	if group.Status == GROUP_STATUS_CLOSED {
		return fmt.Errorf("小组已被关闭")
	}
	if len(group.Name) > s.cfg.GroupNameLen {
		return fmt.Errorf(fmt.Sprintf("小组名称不能大于%d个字符", s.cfg.GroupNameLen))
	}
	_name := utils.CensorWords(group.Name)
	if _name != group.Name {
		return fmt.Errorf("小组名称包含屏蔽字")
	}
	if len(group.Description) > s.cfg.GroupDescMaxLen {
		return fmt.Errorf(fmt.Sprintf("小组描述内容不能大于%d个字符", s.cfg.GroupDescMaxLen))
	}
	if len(group.Description) < s.cfg.GroupDescMinLen {
		return fmt.Errorf(fmt.Sprintf("小组描述内容不能小于%d个字符", s.cfg.GroupDescMinLen))
	}
	if group.Uid <= 0 {
		return fmt.Errorf("未设置创建人")
	}
	if len(group.GameIds) == 0 {
		return fmt.Errorf("未选择游戏分类")
	}
	if group.BgImg == 0 {
		return fmt.Errorf("未设置背景图")
	}
	//新建后更新间隔
	if group.CreateTime+update_group_limit_seconds > time.Now().Unix() {
		return fmt.Errorf(fmt.Sprintf("新建小组%d秒内不允许修改", update_group_limit_seconds))
	}

	//图片压缩
	src_group := s.Get(group.Id)
	if src_group.BgImg != group.BgImg {
		err := s.thumbnailImgResize(group)
		if err != nil {
			return err
		}
	}

	o := dbs.NewOrm(db_aliasname)
	_, err := o.Update(group, "groupname", "description", "uid", "country", "city", "gameids", "displayorder", "img", "bgimg",
		"belong", "type", "searchkeyword")
	if err != nil {
		return err
	}
	cache := utils.GetCache()
	cache.Replace(s.GetCacheKey(group.Id), *group, 1*time.Hour)
	return nil
}

func (s *GroupService) Get(groupId int64) *Group {
	cache := utils.GetCache()
	group := Group{}
	err := cache.Get(s.GetCacheKey(groupId), &group)
	if err == nil {
		if group.IsDeleted {
			return nil
		}
		return &group
	}
	group.Id = groupId
	o := dbs.NewOrm(db_aliasname)
	err = o.Read(&group)
	if err != nil {
		group.IsDeleted = true
		cache.Add(s.GetCacheKey(groupId), group, 24*time.Hour)
		return nil
	}
	cache.Add(s.GetCacheKey(groupId), group, 1*time.Hour)
	return &group
}

func (s *GroupService) AsyncActionCount(groupId int64, gps []GP_PROPERTY, incrs []int) {
	group := s.Get(groupId)
	if group == nil {
		return
	}
	if len(gps) != len(incrs) {
		panic("fields和incrs的数量必须一致")
	}

	cclient := ssdb.New(use_ssdb_group_db)
	for i, field := range gps {
		_key := fmt.Sprintf(count_action_set, string(field))
		cclient.Zincrby(_key, groupId, int64(incrs[i]))
	}
}

func (s *GroupService) ActionCount(groupId int64, gps []GP_PROPERTY, incrs []int) {
	group := s.Get(groupId)
	if group == nil {
		return
	}
	if len(gps) != len(incrs) {
		panic("fields和incrs的数量必须一致")
	}

	params := make(orm.Params)
	for i, field := range gps {
		switch field {
		case GP_PROPERTY_MEMBERS:
			group.Members += incrs[i]
			params["members"] = orm.ColValue(orm.Col_Add, incrs[i])
		case GP_PROPERTY_THREADS:
			group.Threads += incrs[i]
			params["threads"] = orm.ColValue(orm.Col_Add, incrs[i])
		default:
		}
	}

	o := dbs.NewOrm(db_aliasname)
	_, err := o.QueryTable(&Group{}).Filter("id", groupId).Update(params)
	if err != nil {
		logs.Errorf("小组计数数据库更新失败:%s", err.Error())
		return
	}
	cache := utils.GetCache()
	cache.Replace(s.GetCacheKey(group.Id), *group, 1*time.Hour)

	//更新搜索引擎文档属性
	s.UpdateBaseSearchEngineAttrs([]*Group{group})
}

func (s *GroupService) Join(uid int64, groupId int64) error {
	group := s.Get(groupId)
	if group == nil {
		return fmt.Errorf("小组不存在")
	}
	if group.Status == GROUP_STATUS_CLOSED {
		return fmt.Errorf("小组已关闭,不允许加入")
	}
	if group.Uid == uid {
		return fmt.Errorf("无法加入自己创建的组")
	}
	mjoined_key := fmt.Sprintf(group_member_joined_sortset, uid)     //用户加入小组的集合key
	gjoined_key := fmt.Sprintf(group_joined_member_sortset, groupId) //小组中加入用户的集合key

	score, _ := ssdb.New(use_ssdb_group_db).Zscore(mjoined_key, groupId)
	if score > 0 {
		return fmt.Errorf("已加入该小组,不能重复加入")
	}

	_, gmTbl, _ := s.GetGMTable(group.Uid)
	_, mgTbl, _ := s.GetMGTable(uid)
	ts := time.Now().Unix()
	o := dbs.NewOrm(db_aliasname)
	o.Begin()
	sql := fmt.Sprintf("insert %s(groupid,uid,ts) values(?,?,?)", gmTbl)
	_, err := o.Raw(sql, groupId, uid, ts).Exec()
	sql = fmt.Sprintf("insert %s(uid,groupid,ts) values(?,?,?)", mgTbl)
	_, err = o.Raw(sql, uid, groupId, ts).Exec()
	if err != nil {
		o.Rollback()
		return fmt.Errorf("加入小组失败:001")
	}
	ssdb_client := ssdb.New(use_ssdb_group_db)
	//用户加入的小组
	_, err = ssdb_client.Zadd(mjoined_key, groupId, ts)
	//小组中加入的用户
	_, err = ssdb_client.Zadd(gjoined_key, uid, ts)
	if err != nil {
		o.Rollback()
		return fmt.Errorf("加入小组失败:002")
	}
	err = o.Commit()
	if err != nil {
		ssdb_client.Zrem(mjoined_key, groupId)
		ssdb_client.Zrem(gjoined_key, uid)
	}

	s.AsyncActionCount(groupId, []GP_PROPERTY{GP_PROPERTY_MEMBERS}, []int{1})

	//计数
	go NewMemberCountService().ActionCountProperty([]MC_PROPERTY{MC_PROPERTY_JOINS}, []COUNT_ACTION{COUNT_ACTION_INCR}, group.Uid)
	//检查是否完成招募
	go s.CompletedInvited(group, true)

	return nil
}

func (s *GroupService) Exit(uid int64, groupId int64) error {
	group := s.Get(groupId)
	if group == nil {
		return fmt.Errorf("小组不存在")
	}
	if group.Uid == uid {
		return fmt.Errorf("无法离开自己创建的组")
	}
	if group.Status == GROUP_STATUS_CLOSED {
		return fmt.Errorf("小组已关闭,不允许退出")
	}
	mjoined_key := fmt.Sprintf(group_member_joined_sortset, uid)     //用户加入小组的集合key
	gjoined_key := fmt.Sprintf(group_joined_member_sortset, groupId) //小组中加入用户的集合key

	_, gmTbl, _ := s.GetGMTable(group.Uid)
	_, mgTbl, _ := s.GetMGTable(uid)
	o := dbs.NewOrm(db_aliasname)
	sql := fmt.Sprintf("delete from %s where groupid=? and uid=?", gmTbl)
	o.Raw(sql, groupId, uid).Exec()
	sql = fmt.Sprintf("delete from %s where uid=? and groupid=?", mgTbl)
	o.Raw(sql, uid, groupId).Exec()

	ssdb_client := ssdb.New(use_ssdb_group_db)
	//用户离开的小组
	ssdb_client.Zrem(mjoined_key, groupId)
	//小组中去除用户
	ssdb_client.Zrem(gjoined_key, uid)

	s.AsyncActionCount(groupId, []GP_PROPERTY{GP_PROPERTY_MEMBERS}, []int{-1})

	//计数
	go NewMemberCountService().ActionCountProperty([]MC_PROPERTY{MC_PROPERTY_JOINS}, []COUNT_ACTION{COUNT_ACTION_DECR}, group.Uid)

	//检查是否到达人数预警线
	go s.LessThanLimitUsers(group)

	return nil
}

//func (s *GroupService) UpdateMemberLastAction(groupId int64, uid int64, t time.Time) {
//	mjoined_key := fmt.Sprintf(group_member_joined_sortset, uid) //用户加入小组的集合key
//	ts, err := ssdb.New(use_ssdb_group_db).Zscore(mjoined_key, groupId)
//	if ts == 0 || err != nil {
//		return
//	}
//	if ts > t.Unix() {
//		return
//	}
//	gap := t.Unix() - ts
//	ssdb.New(use_ssdb_group_db).Zincrby(mjoined_key, groupId, gap)
//}

func (s *GroupService) MyJoins(uid int64, page int, size int) (int, []*Group) {
	mjoined_key := fmt.Sprintf(group_member_joined_sortset, uid) //用户加入小组的集合key
	client := ssdb.New(use_ssdb_group_db)
	count, _ := client.Zcard(mjoined_key)
	vals, _ := client.Zrevrange(mjoined_key, (page-1)*size, page*size, reflect.TypeOf(int64(0)))
	ids := make([]int64, len(vals), len(vals))
	for i, val := range vals {
		ids[i] = *(val.(*int64))
	}

	groups := s.multiGets(ids)

	return count, groups
}

func (s *GroupService) MyAllJoinGroupIds(uid int64) map[int64]int64 {
	mjoined_key := fmt.Sprintf(group_member_joined_sortset, uid) //用户加入小组的集合key
	client := ssdb.New(use_ssdb_group_db)
	vals, _ := client.Zscan(mjoined_key, 0, time.Now().Unix(), 1<<10, reflect.TypeOf(int64(0)))
	ids := make(map[int64]int64)
	for _, val := range vals {
		id := *(val.(*int64))
		ids[id] = 0
	}
	return ids
}

func (s *GroupService) IsJoined(uid int64, groupId int64) bool {
	if uid <= 0 {
		return false
	}
	group := s.Get(groupId)
	if group == nil {
		return false
	}
	if group.Uid == uid {
		return true
	}
	mjoined_key := fmt.Sprintf(group_member_joined_sortset, uid) //用户加入小组的集合key
	client := ssdb.New(use_ssdb_group_db)
	joined, _ := client.Zexists(mjoined_key, groupId)
	return joined
}

func (s *GroupService) multiGets(groupIds []int64) []*Group {
	if len(groupIds) == 0 {
		return []*Group{}
	}
	gkeys := make([]string, len(groupIds), len(groupIds))
	for i, id := range groupIds {
		gkeys[i] = s.GetCacheKey(id)
	}
	cache := utils.GetCache()
	getter, _ := cache.GetMulti(gkeys...)

	groups := []*Group{}
	for _, id := range groupIds {
		var _group Group
		err := getter.Get(s.GetCacheKey(id), &_group)
		if err == nil && !_group.IsDeleted {
			groups = append(groups, &_group)
			continue
		}
		_g := s.Get(id)
		if _g != nil {
			groups = append(groups, _g)
		}
	}
	return groups
}

func (s *GroupService) getMyGroupIdsCacheKey(uid int64) string {
	return fmt.Sprintf("group_my_%d", uid)
}

func (s *GroupService) MyGroups(uid int64) []*Group {
	ckey := s.getMyGroupIdsCacheKey(uid)
	cache := utils.GetCache()
	var ids []int64
	err := cache.Get(ckey, &ids)
	if err == nil {
		return s.multiGets(ids)
	}
	_g := &Group{}
	o := dbs.NewOrm(db_aliasname)
	var maps []orm.Params
	num, err := o.Raw(fmt.Sprintf("SELECT id FROM %s WHERE status <> ? and uid=?", _g.TableName()), GROUP_STATUS_CLOSED, uid).Values(&maps)
	ids = []int64{}
	if err == nil && num > 0 {
		for i := range maps {
			_strid := maps[i]["id"].(string)
			_id, _ := strconv.ParseInt(_strid, 10, 64)
			ids = append(ids, _id)
		}
	}
	cache.Add(ckey, ids, 6*time.Hour)
	groups := s.multiGets(ids)
	return groups
}

func (s *GroupService) Search(words string, page int, size int, match_mode string, sort GP_SEARCH_SORT,
	filters []search.SearchFilter, filterRanges []search.FilterRangeInt) (int, []*Group) {
	search_config := &search.SearchOptions{
		Host:    search_server,
		Port:    search_port,
		Timeout: search_timeout,
	}
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size
	search_config.Offset = offset
	search_config.Limit = size
	search_config.Filters = filters
	search_config.FilterRangeInt = filterRanges
	search_config.MaxMatches = 500

	sorts := []string{"@weight DESC"}
	switch sort {
	case GP_SEARCH_SORT_USERS:
		sorts = append(sorts, "members DESC")
		break
	case GP_SEARCH_SORT_OFFICIAL:
		sorts = append(sorts, "belong DESC")
		sorts = append(sorts, "vitality DESC")
		break
	case GP_SEARCH_SORT_VITALITY:
		sorts = append(sorts, "vitality DESC")
		break
	case GP_SEARCH_SORT_ENDTIME:
		sorts = append(sorts, "endtime ASC")
		break
	default:
		sorts = append(sorts, "displayorder DESC")
		sorts = append(sorts, "vitality DESC")
	}

	sph := search.NewSearcher(search_config)
	ids, total, err := sph.Query(words, sorts, group_index_name, match_mode)
	if err != nil {
		return 0, []*Group{}
	}
	groups := s.multiGets(ids)
	return total, groups
}

func (s *GroupService) GroupUserCount(groupId int64) int {
	gjoined_key := fmt.Sprintf(group_joined_member_sortset, groupId)
	count, _ := ssdb.New(use_ssdb_group_db).Zcard(gjoined_key)
	return int(count)
}

func (s *GroupService) ChangeStatus(group *Group) error {
	o := dbs.NewOrm(db_aliasname)
	_, err := o.Update(group, "status")
	if err != nil {
		return fmt.Errorf("更新状态失败")
	}
	cache := utils.GetCache()
	cache.Delete(s.GetCacheKey(group.Id))
	return nil
}

func (s *GroupService) LessThanLimitUsers(group *Group) error {
	//警示和关闭中不检查，官方账号不检查
	if group == nil ||
		group.Status == GROUP_STATUS_RECRUITING ||
		group.Status == GROUP_STATUS_LOWMEMBER ||
		group.Status == GROUP_STATUS_CLOSED ||
		group.Belong == GROUP_BELONG_OFFICIAL {
		return fmt.Errorf("未在需要检查的状态中")
	}
	count := s.GroupUserCount(group.Id)
	if count < group.MinUsers {

		group.Status = GROUP_STATUS_LOWMEMBER
		group.StartTime = time.Now().Unix()
		addHours := time.Duration(s.cfg.CreateGroupRecruitDay * 24)
		group.EndTime = time.Now().Add(addHours * time.Hour).Unix()
		err := s.ChangeStatus(group)
		if err != nil {
			return err
		}

		//设置定时任务
		go tjSetJob(group)
	}
	return nil
}

func (s *GroupService) CompletedInvited(group *Group, doTj bool) (bool, error) {
	//开放和关闭中不检查，官方账号不检查
	if group == nil ||
		group.Status == GROUP_STATUS_OPENING ||
		group.Status == GROUP_STATUS_CLOSED ||
		group.Belong == GROUP_BELONG_OFFICIAL {
		return false, fmt.Errorf("未在需要检查的状态中")
	}
	count := s.GroupUserCount(group.Id)
	if count >= group.MinUsers {

		group.Status = GROUP_STATUS_OPENING
		group.StartTime = 0
		group.EndTime = 0
		err := s.ChangeStatus(group)
		if err != nil {
			return false, err
		}
		//正式扣除积分
		client, transport, err := credit_client.NewClient(credit_service_host)
		if err != nil {
			go tjRemoveJob(group)
			return false, fmt.Errorf("扣除正式积分失败")
		}
		defer func() {
			if transport != nil {
				transport.Close()
			}
		}()
		client.Do(&credit_proxy.OperationCreditParameter{
			No:     group.OrderNo,
			Action: credit_proxy.OPERATION_ACTOIN_UNLOCK,
		})

		//删除定时任务
		if doTj {
			go tjRemoveJob(group)
		}
	}
	return true, nil
}

func (s *GroupService) InviteFriends(uid int64, groupId int64) map[string][]*InviteMember {
	fs := &passport.FriendShips{}
	friendsPy := fs.BothFriendUidsPy(uid)
	ckey := fmt.Sprintf(group_invited_set, groupId, uid)
	cclient := ssdb.New(use_ssdb_group_db)
	objs, _ := cclient.Zscan(ckey, 0, time.Now().UnixNano()/1000, 1<<10, reflect.TypeOf(int64(0)))
	invmaps := make(map[int64]bool)
	joined_uids := []interface{}{}
	for _, obj := range objs {
		_uid := *(obj.(*int64))
		invmaps[_uid] = true
	}
	for _, friends := range friendsPy {
		for _, uid := range friends {
			joined_uids = append(joined_uids, uid)
		}
	}
	//已加入组的
	gjoined_key := fmt.Sprintf(group_joined_member_sortset, groupId) //小组中加入用户的集合key
	joinedObjs, _ := cclient.MultiZget(gjoined_key, joined_uids, reflect.TypeOf(int64(0)))
	joineds := make(map[int64]int64)
	for obj, score := range joinedObjs {
		_uid := *(obj.(*int64))
		joineds[_uid] = score
	}

	maps := make(map[string][]*InviteMember)
	for c, uids := range friendsPy {
		ims := []*InviteMember{}
		for _, uid := range uids {
			ed := false
			jd := false
			if _, ok := invmaps[uid]; ok {
				ed = true
			}
			if _, ok := joineds[uid]; ok {
				jd = true
			}
			ims = append(ims, &InviteMember{
				Uid:     uid,
				Invited: ed,
				Joined:  jd,
			})
		}
		maps[c] = ims
	}
	return maps
}

func (s *GroupService) Invite(uid int64, groupId int64, inviteUids []int64) {
	if len(inviteUids) == 0 {
		return
	}
	ckey := fmt.Sprintf(group_invited_set, groupId, uid)
	vals := make([]interface{}, len(inviteUids), len(inviteUids))
	scores := make([]int64, len(inviteUids), len(inviteUids))
	for i := range inviteUids {
		vals[i] = inviteUids[i]
		scores[i] = time.Now().UnixNano() / 1000
	}
	ssdb.New(use_ssdb_group_db).MultiZadd(ckey, vals, scores)

	go func() {
		ps := passport.NewMemberProvider()
		author := ps.Get(uid)
		if author == nil {
			return
		}
		group := s.Get(groupId)
		if group == nil {
			return
		}
		for _, iUid := range inviteUids {
			invMsg := fmt.Sprintf("%s邀请你加入%s小组", author.NickName, group.Name)
			message.SendMsgV2(uid, iUid, MSG_TYPE_INVITED, invMsg, strconv.FormatInt(groupId, 10), nil)
		}
	}()
}

func (s *GroupService) UpdateSearchEngineAttrs(groups []*Group, attrs []string) error {
	if len(groups) == 0 || len(attrs) == 0 {
		return nil
	}

	valuesA := [][]interface{}{}
	for _, group := range groups {
		vals := []interface{}{}
		vals = append(vals, uint64(group.Id))
		for _, attr := range attrs {
			switch attr {
			case "members":
				vals = append(vals, group.Members)
				break
			case "threads":
				vals = append(vals, group.Threads)
				break
			case "displayorder":
				vals = append(vals, group.DisplarOrder)
				break
			case "status":
				vals = append(vals, int(group.Status))
				break
			case "belong":
				vals = append(vals, int(group.Belong))
				break
			case "type":
				vals = append(vals, int(group.Type))
				break
			case "vitality":
				vals = append(vals, group.Vitality)
				break
			case "recommend":
				recommend := 0
				if group.Recommend {
					recommend = 1
				}
				vals = append(vals, recommend)
				break
			case "starttime":
				vals = append(vals, int(group.StartTime))
				break
			case "endtime":
				vals = append(vals, int(group.EndTime))
				break
			}
		}
		valuesA = append(valuesA, vals)
	}
	search_config := &search.SearchOptions{
		Host:    search_server,
		Port:    search_port,
		Timeout: search_timeout,
	}
	sph := search.NewSearcher(search_config)
	_, err := sph.UpdateAttributes(group_index_name, attrs, valuesA)
	if err != nil {
		logs.Errorf("group update search's attribute normal fail:%v", err)
		return err
	}
	_, err = sph.FlushAttributes()
	if err != nil {
		logs.Errorf("group flush search's attribute fail:%v", err)
		return err
	}

	//	//有bug
	//	if mva {
	//		attrsB := []string{"gameids"}
	//		valuesB := [][]interface{}{}
	//		for _, group := range groups {
	//			valuesB = append(valuesB, []interface{}{
	//				uint64(group.Id),
	//				group.GameIDs(),
	//			})
	//		}
	//		_, err = sph.UpdateAttributes(group_index_name, attrsB, valuesB)
	//		if err != nil {
	//			logs.Errorf("group update search's attribute mva fail:%v", err)
	//			return err
	//		}
	//	}

	//	_, err = sph.FlushAttributes()
	//	if err != nil {
	//		logs.Errorf("group flush search's attribute fail:%v", err)
	//	}
	return nil
}

func (s *GroupService) UpdateBaseSearchEngineAttrs(groups []*Group) error {
	attrs := []string{"members", "threads", "displayorder", "status", "belong",
		"type", "recommend", "starttime", "endtime"}
	return s.UpdateSearchEngineAttrs(groups, attrs)
}

////////////////////////////////////////////////////////////////////////////////
// 计数
////////////////////////////////////////////////////////////////////////////////
type COUNT_ACTION int

const (
	COUNT_ACTION_INCR COUNT_ACTION = 1
	COUNT_ACTION_DECR COUNT_ACTION = 2
)

type MC_PROPERTY int

const (
	MC_PROPERTY_GROUPS    MC_PROPERTY = 1
	MC_PROPERTY_TODINGS   MC_PROPERTY = 2
	MC_PROPERTY_TOCAIS    MC_PROPERTY = 4
	MC_PROPERTY_FROMDINGS MC_PROPERTY = 8
	MC_PROPERTY_FROMCAIS  MC_PROPERTY = 16
	MC_PROPERTY_POSTS     MC_PROPERTY = 32
	MC_PROPERTY_THREADS   MC_PROPERTY = 64
	MC_PROPERTY_JOINS     MC_PROPERTY = 128
	MC_PROPERTY_SHARES    MC_PROPERTY = 256
	MC_PROPERTY_REPORTS   MC_PROPERTY = 512
)

func NewMemberCountService() *MemberCountService {
	return &MemberCountService{}
}

type MemberCountService struct{}

func (c *MemberCountService) getCacheKey(uid int64) string {
	return fmt.Sprintf("group_member_count:%d", uid)
}

func (c *MemberCountService) GetCount(uid int64) *MemberCount {
	ckey := c.getCacheKey(uid)
	cache := utils.GetCache()
	mc := &MemberCount{}
	err := cache.Get(ckey, mc)
	if err == nil {
		return mc
	}
	mc.Uid = uid
	o := dbs.NewOrm(db_aliasname)
	err = o.Read(mc)
	if err == nil {
		cache.Add(ckey, *mc, 24*time.Hour)
	} else {
		mc.LastTs = time.Now().Unix()
		o.Insert(mc)
		cache.Add(ckey, *mc, 24*time.Hour)
	}
	return mc
}

func (c *MemberCountService) ActionCountProperty(mcs []MC_PROPERTY, actions []COUNT_ACTION, uid int64) error {
	if len(mcs) != len(actions) {
		panic("mcs和actions的数量必须一致")
	}
	mc := c.GetCount(uid)

	params := make(orm.Params)
	for i, property := range mcs {
		switch property {
		case MC_PROPERTY_GROUPS:
			if actions[i] == COUNT_ACTION_INCR {
				mc.Groups++
				params["groups"] = orm.ColValue(orm.Col_Add, 1)
			} else {
				mc.Groups--
				params["groups"] = orm.ColValue(orm.Col_Except, 1)
			}
		case MC_PROPERTY_TODINGS:
			if actions[i] == COUNT_ACTION_INCR {
				mc.ToDings++
				params["todings"] = orm.ColValue(orm.Col_Add, 1)
			} else {
				mc.ToDings--
				params["todings"] = orm.ColValue(orm.Col_Except, 1)
			}
		case MC_PROPERTY_TOCAIS:
			if actions[i] == COUNT_ACTION_INCR {
				mc.ToCais++
				params["tocais"] = orm.ColValue(orm.Col_Add, 1)
			} else {
				mc.ToCais--
				params["tocais"] = orm.ColValue(orm.Col_Except, 1)
			}
		case MC_PROPERTY_FROMDINGS:
			if actions[i] == COUNT_ACTION_INCR {
				mc.FromDings++
				params["fromdings"] = orm.ColValue(orm.Col_Add, 1)
			} else {
				mc.FromDings--
				params["fromdings"] = orm.ColValue(orm.Col_Except, 1)
			}
		case MC_PROPERTY_FROMCAIS:
			if actions[i] == COUNT_ACTION_INCR {
				mc.FromCais++
				params["fromcais"] = orm.ColValue(orm.Col_Add, 1)
			} else {
				mc.FromCais--
				params["fromcais"] = orm.ColValue(orm.Col_Except, 1)
			}
		case MC_PROPERTY_POSTS:
			if actions[i] == COUNT_ACTION_INCR {
				mc.Posts++
				params["posts"] = orm.ColValue(orm.Col_Add, 1)
			} else {
				mc.Posts--
				params["posts"] = orm.ColValue(orm.Col_Except, 1)
			}
		case MC_PROPERTY_THREADS:
			if actions[i] == COUNT_ACTION_INCR {
				mc.Threads++
				params["threads"] = orm.ColValue(orm.Col_Add, 1)
			} else {
				mc.Threads--
				params["threads"] = orm.ColValue(orm.Col_Except, 1)
			}
		case MC_PROPERTY_JOINS:
			if actions[i] == COUNT_ACTION_INCR {
				mc.Joins++
				params["joins"] = orm.ColValue(orm.Col_Add, 1)
			} else {
				mc.Joins--
				params["joins"] = orm.ColValue(orm.Col_Except, 1)
			}
		case MC_PROPERTY_SHARES:
			if actions[i] == COUNT_ACTION_INCR {
				mc.Shares++
				params["shares"] = orm.ColValue(orm.Col_Add, 1)
			} else {
				mc.Shares--
				params["shares"] = orm.ColValue(orm.Col_Except, 1)
			}
		case MC_PROPERTY_REPORTS:
			if actions[i] == COUNT_ACTION_INCR {
				mc.Reports++
				params["reports"] = orm.ColValue(orm.Col_Add, 1)
			} else {
				mc.Reports--
				params["reports"] = orm.ColValue(orm.Col_Except, 1)
			}
		default:
		}
	}

	mc.LastTs = time.Now().Unix()

	o := dbs.NewOrm(db_aliasname)
	_, err := o.QueryTable(&MemberCount{}).Filter("uid", uid).Update(params)
	if err != nil {
		logs.Errorf("用户计数数据库更新失败:%s", err.Error())
		return err
	}
	cache := utils.GetCache()
	cache.Replace(c.getCacheKey(uid), *mc, 12*time.Hour)
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// 帖子
////////////////////////////////////////////////////////////////////////////////

type TH_PROPERTY string

const (
	TH_PROPERTY_FAVORITES TH_PROPERTY = "favorites"
	TH_PROPERTY_REPLIES   TH_PROPERTY = "replies"
	TH_PROPERTY_SHARES    TH_PROPERTY = "shares"
	TH_PROPERTY_VIEWS     TH_PROPERTY = "views"
)

var postTbls map[int]string = make(map[int]string)
var thcQs map[int64]*time.Timer = make(map[int64]*time.Timer)
var thcQmutex *sync.Mutex = new(sync.Mutex)
var thcLastQs map[int64]*time.Timer = make(map[int64]*time.Timer)
var thcLastQmutex *sync.Mutex = new(sync.Mutex)

const (
	group_post_table_fmt = "post_%d"
)

func NewPostId(tblId int) string {
	ms := time.Now().UnixNano() / 1000
	return fmt.Sprintf("%03dA%d", tblId, ms)
}

func PostTableId(postId string) int {
	if len(postId) != 20 {
		return 0
	}
	_strid := postId[:3]
	_id, _ := strconv.Atoi(_strid)
	return _id
}

func getPostTableId(thread *Thread) (int, string, error) {
	str_groupid := strconv.FormatInt(thread.GroupId, 10)
	str_pfx := str_groupid[:1] //1位
	i_tag, err := strconv.Atoi(str_pfx)
	if err != nil {
		return 0, "", fmt.Errorf("参数错误001")
	}
	if name, ok := postTbls[i_tag]; ok {
		return i_tag, name, nil
	}
	tbl_mutex.Lock()
	defer tbl_mutex.Unlock()
	if name, ok := postTbls[i_tag]; ok {
		return i_tag, name, nil
	}
	tbl := fmt.Sprintf(group_post_table_fmt, i_tag)

	o := dbs.NewOrm(db_aliasname)
	create_tbl_sql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(id char(20) NOT NULL,
	  tid int(11) NOT NULL,
	  authorid int(11) NOT NULL,
	  subject char(30) NOT NULL,
	  dateline int(11) NOT NULL,
	  message mediumtext NOT NULL,
	  ip char(15) NOT NULL,
	  invisible tinyint(1) NOT NULL,
	  ding int(11) NOT NULL,
	  cai int(11) NOT NULL,
	  position int(11) NOT NULL,
	  replyid char(20) NOT NULL,
	  replyuid int(11) NOT NULL,
	  replyposition int(11) NOT NULL,
	  img int(11) NOT NULL,
	  resources varchar(500) NOT NULL,
	  longitude double(7,5) NOT NULL,
	  latitude double(7,5) NOT NULL,
	  fromdev tinyint(2) NOT NULL,
	  PRIMARY KEY (id),
	  UNIQUE KEY idx_thread_position (tid,position)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`, tbl)
	o.Raw(create_tbl_sql).Exec()
	postTbls[i_tag] = tbl

	savePostTableToDb(i_tag, tbl, &PostTable{
		Id:      int64(i_tag),
		TblName: tbl,
		Ts:      time.Now().Unix(),
	})
	return i_tag, tbl, nil
}

func NewThreadService(cfg *GroupCfg) *ThreadService {
	ts := &ThreadService{}
	ts.cfg = cfg
	return ts
}

type ThreadService struct {
	BaseService
}

func (c *ThreadService) getCacheKey(threadId int64) string {
	return fmt.Sprintf("group_thread_%d", threadId)
}

func (c *ThreadService) verifyThread(thread *Thread) error {
	if c.cfg == nil {
		return fmt.Errorf("未设定配置")
	}
	if len(thread.Subject) == 0 {
		return fmt.Errorf("标题不能为空")
	}
	if thread.GroupId <= 0 {
		return fmt.Errorf("未指定所属组")
	}
	if thread.AuthorId <= 0 {
		return fmt.Errorf("未指定发表人")
	}
	gs := NewGroupService(c.cfg)
	group := gs.Get(thread.GroupId)
	if group == nil || group.Status == GROUP_STATUS_CLOSED {
		return fmt.Errorf("组不存在")
	}
	if group.Status == GROUP_STATUS_RECRUITING {
		return fmt.Errorf("组正在招募中,不能发帖")
	}
	return nil
}

func (c *ThreadService) Create(thread *Thread, post *Post) error {
	err := c.verifyThread(thread)
	if err != nil {
		return err
	}
	ps := passport.NewMemberProvider()
	author := ps.Get(thread.AuthorId)
	if author == nil {
		return fmt.Errorf("创建人不存在")
	}

	thread.Author = author.NickName
	thread.DateLine = time.Now().Unix()
	thread.LastPost = time.Now().Unix()
	thread.Status = c.cfg.NewThreadDefaultStatus
	thread.Subject = utils.CensorWords(thread.Subject) //过滤关键字

	//设置post表id
	tblId, _, err := getPostTableId(thread)
	if err != nil {
		return fmt.Errorf("创建帖子失败:001")
	}
	thread.PostTableId = tblId

	o := dbs.NewOrm(db_aliasname)
	id, err := o.Insert(thread)
	if err != nil {
		logs.Errorf("新建帖子失败:%s", err.Error())
		return fmt.Errorf("创建帖子失败:002")
	}
	thread.Id = id
	//插入楼主评论
	post.ThreadId = id
	post.Id = NewPostId(int(id))
	err = NewPostService(c.cfg).Create(post)
	if err != nil {
		o.Delete(thread)
		logs.Errorf("新建帖子失败:%s", err.Error())
		return fmt.Errorf("创建帖子失败:003")
	}

	thread.Lordpid = post.Id
	thread.Img = post.Img
	o.Update(thread)

	ssdb.New(use_ssdb_group_db).Set(c.getCacheKey(id), thread)

	//计数
	go NewMemberCountService().ActionCountProperty([]MC_PROPERTY{MC_PROPERTY_THREADS}, []COUNT_ACTION{COUNT_ACTION_INCR}, thread.AuthorId)
	go func() {
		gs := NewGroupService(c.cfg)
		gs.AsyncActionCount(thread.GroupId, []GP_PROPERTY{GP_PROPERTY_THREADS}, []int{1})
	}()
	return err
}

func (c *ThreadService) Get(threadId int64) *Thread {
	var thread Thread
	cclient := ssdb.New(use_ssdb_group_db)
	ckey := c.getCacheKey(threadId)
	err := cclient.Get(ckey, &thread)
	if err == nil {
		return &thread
	}
	thread = Thread{
		Id: threadId,
	}
	o := dbs.NewOrm(db_aliasname)
	err = o.Read(&thread)
	if err == nil {
		cclient.Set(ckey, &thread)
		return &thread
	}
	return nil
}

func (c *ThreadService) ActionProperty(ths []TH_PROPERTY, vals []int64, threadId int64) {
	if len(ths) != len(vals) {
		return
	}
	for i, th := range ths {
		ckey := fmt.Sprintf("group_thread_%d_%s_count", threadId, string(th))
		ssdb.New(use_ssdb_group_db).Incrby(ckey, vals[i])
	}

	if _, ok := thcQs[threadId]; ok {
		return
	}

	thcQmutex.Lock()
	defer thcQmutex.Unlock()
	if _, ok := thcQs[threadId]; ok {
		return
	}
	//延迟更新数据库计数
	refId := threadId
	thcQs[threadId] = time.AfterFunc(1*time.Minute, func() {
		thread := c.Get(refId)
		if thread == nil {
			return
		}
		cclient := ssdb.New(use_ssdb_group_db)
		fmtstr := "group_thread_%d_%s_count"
		var i int64 = 0
		cclient.Get(fmt.Sprintf(fmtstr, refId, string(TH_PROPERTY_FAVORITES)), &i)
		thread.Favorites = int(i)
		i = 0
		cclient.Get(fmt.Sprintf(fmtstr, refId, string(TH_PROPERTY_REPLIES)), &i)
		thread.Replies = int(i)
		i = 0
		cclient.Get(fmt.Sprintf(fmtstr, refId, string(TH_PROPERTY_SHARES)), &i)
		thread.Shares = int(i)
		i = 0
		cclient.Get(fmt.Sprintf(fmtstr, refId, string(TH_PROPERTY_VIEWS)), &i)
		thread.Views = int(i)

		o := dbs.NewOrm(db_aliasname)
		_, err := o.Update(thread)
		if err == nil {
			cclient.Set(c.getCacheKey(refId), *thread)
		}
		thcQmutex.Lock()
		defer thcQmutex.Unlock()
		delete(thcQs, refId)
	})
}

func (c *ThreadService) Gets(groupId int64, page int, size int, uid int64) (int, []*Thread) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	threads := []*Thread{}
	group := NewGroupService(c.cfg).Get(groupId)
	if group == nil {
		return 0, threads
	}
	o := dbs.NewOrm(db_aliasname)
	var maps []orm.Params
	_, err := o.QueryTable(&Thread{}).Filter("groupid", groupId).OrderBy("-lastpost").Limit(size).Offset((page-1)*size).Values(&maps, "id")
	if err == nil {
		strids := make([]string, len(maps), len(maps))
		for i, m := range maps {
			_id := m["Id"].(int64)
			strids[i] = c.getCacheKey(_id)
		}
		objs := ssdb.New(use_ssdb_group_db).MultiGet(strids, reflect.TypeOf(Thread{}))
		for _, t := range objs {
			threads = append(threads, t.(*Thread))
		}
		return group.Threads, threads
	}
	return 0, threads
}

func (c *ThreadService) SetLastPost(post *Post) {

	type LastInfo struct {
		LastId      string
		LastPost    int64
		LastPostUid int64
	}

	ckey := fmt.Sprintf("group_thread_%d_last", post.ThreadId)
	ssdb.New(use_ssdb_group_db).Set(ckey, LastInfo{
		LastId:      post.Id,
		LastPost:    post.DateLine,
		LastPostUid: post.AuthorId,
	})

	if _, ok := thcLastQs[post.ThreadId]; ok {
		return
	}

	thcLastQmutex.Lock()
	defer thcLastQmutex.Unlock()
	if _, ok := thcLastQs[post.ThreadId]; ok {
		return
	}

	refId := post.ThreadId
	thcLastQs[refId] = time.AfterFunc(10*time.Second, func() {
		ckey := fmt.Sprintf("group_thread_%d_last", refId)
		var lastInfo LastInfo
		cclient := ssdb.New(use_ssdb_group_db)
		err := cclient.Get(ckey, &lastInfo)
		if err != nil {
			return
		}

		mp := passport.NewMemberProvider()
		member := mp.Get(lastInfo.LastPostUid)
		if member == nil {
			return
		}
		_t := c.Get(refId)
		if _t == nil {
			return
		}

		_t.LastId = lastInfo.LastId
		_t.LastPost = lastInfo.LastPost
		_t.LastPoster = member.NickName
		_t.LastPostUid = lastInfo.LastPostUid

		thread_ckey := c.getCacheKey(post.ThreadId)
		cclient.Set(thread_ckey, _t)

		o := dbs.NewOrm(db_aliasname)
		o.Update(_t, "lastpost", "lastposter", "lastpostuid", "lastid")

		thcLastQmutex.Lock()
		defer thcLastQmutex.Unlock()
		delete(thcLastQs, refId)
		cclient.Del(ckey)
	})
}

////////////////////////////////////////////////////////////////////////////////
// 评论
////////////////////////////////////////////////////////////////////////////////
type PostRes struct {
	ImgResource []*ImgRes `json:"imgs"`
}

type ImgRes struct {
	ThumbnailImgId int64 `json:"t_id"`
	BmiddleImgId   int64 `json:"b_id"`
	OriginalImgId  int64 `json:"o_id"`
}

type POST_ACTION int

const (
	POST_ACTION_DING        POST_ACTION = 1
	POST_ACTION_CAI         POST_ACTION = 2
	POST_ACTION_CANCEL_DING POST_ACTION = 3
	POST_ACTION_CANCEL_CAI  POST_ACTION = 4
)

type POST_ACTIONTAG int

const (
	POST_ACTIONTAG_DING POST_ACTIONTAG = 1
	POST_ACTIONTAG_CAI  POST_ACTIONTAG = -1
)

type POST_ORDERBY string

const (
	POST_ORDERBY_ASC  POST_ORDERBY = "asc"
	POST_ORDERBY_DESC POST_ORDERBY = "desc"
)

const (
	thread_post_sets              = "group.thread_%d_post_set"
	thread_post_ding_sets         = "group.thread_%d_post_ding_set"
	thread_post_user_ding_hashtbl = "group.thread_%d_post_ding_uid_%d_hashtbl" //postid <=> 1,0 (1顶)
)

func NewPostService(cfg *GroupCfg) *PostService {
	ps := &PostService{}
	ps.cfg = cfg
	return ps
}

type PostService struct {
	BaseService
}

func (s *PostService) getCacheKey(id string) string {
	return fmt.Sprintf("group_post_%s", id)
}

func (s *PostService) getImgResCacheKey(imgId int64) string {
	return fmt.Sprintf("group_post_img_%d", imgId)
}

func (s *PostService) Get(id string) *Post {
	ckey := s.getCacheKey(id)
	var p Post
	cclient := ssdb.New(use_ssdb_group_db)
	err := cclient.Get(ckey, &p)
	if err == nil {
		return &p
	}
	o := dbs.NewOrm(db_aliasname)
	tblId := PostTableId(id)
	tbl := fmt.Sprintf(group_post_table_fmt, tblId)
	var post *Post
	sql := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", tbl)
	err = o.Raw(sql, id).QueryRow(post)
	if err == nil {
		cclient.Set(ckey, post)
		return post
	}
	return nil
}

func (s *PostService) Action(postId string, uid int64, action POST_ACTION) error {
	post := s.Get(postId)
	if post == nil {
		return fmt.Errorf("评论不存在")
	}

	if s.Actioned(postId, post.ThreadId, uid, action) && action == POST_ACTION_DING {
		return fmt.Errorf("已顶过此评论")
	}

	ckey := fmt.Sprintf(thread_post_ding_sets, post.ThreadId)
	var i int64 = 0
	switch action {
	case POST_ACTION_DING:
		i = 1
		break
	case POST_ACTION_CANCEL_DING:
		i = -1
		break
	default:
	}
	cclient := ssdb.New(use_ssdb_group_db)
	_c, err := cclient.Zincrby(ckey, post.Id, i)
	if err == nil && _c < 0 {
		cclient.Zincrby(ckey, post.Id, -_c)
	}

	//用户顶记录
	user_action_key := fmt.Sprintf(thread_post_user_ding_hashtbl, post.ThreadId, uid)
	switch action {
	case POST_ACTION_DING:
		cclient.Hset(user_action_key, post.Id, int64(POST_ACTIONTAG_DING))
		break
	case POST_ACTION_CANCEL_DING:
		cclient.Hdel(user_action_key, post.Id)
		break
	default:
	}
	return nil
}

func (s *PostService) Actioned(postId string, threadId int64, uid int64, action POST_ACTION) bool {
	user_action_key := fmt.Sprintf(thread_post_user_ding_hashtbl, threadId, uid)
	cclient := ssdb.New(use_ssdb_group_db)
	has, _ := cclient.Hexists(user_action_key, postId)
	return has
}

func (s *PostService) GetPostActionCounts(threadId int64, postIds []string, action POST_ACTIONTAG) map[string]int {
	ckey := fmt.Sprintf(thread_post_ding_sets, threadId)
	cclient := ssdb.New(use_ssdb_group_db)
	vals := make([]interface{}, len(postIds), len(postIds))
	for i, id := range postIds {
		vals[i] = id
	}
	result := make(map[string]int)
	maps, err := cclient.MultiZget(ckey, vals, reflect.TypeOf(""))
	if err != nil {
		return result
	}
	for k, v := range maps {
		id := *(k.(*string))
		result[id] = int(v)
	}
	return result
}

func (s *PostService) GetTops(threadId int64, tops int, pt POST_ACTIONTAG) []*Post {
	ckey := fmt.Sprintf(thread_post_ding_sets, threadId)
	var objs []interface{}
	cclient := ssdb.New(use_ssdb_group_db)
	switch pt {
	case POST_ACTIONTAG_DING:
		objs, _ = cclient.Zrscan(ckey, 1<<32, -1<<32, tops, reflect.TypeOf(""))
	}
	posts := []*Post{}
	if len(objs) == 0 {
		return posts
	}
	post_keys := []string{}
	for _, obj := range objs {
		id, ok := obj.(*string)
		if ok {
			post_keys = append(post_keys, s.getCacheKey(*id))
		}
	}
	post_objs := cclient.MultiGet(post_keys, reflect.TypeOf(Post{}))
	for _, pobj := range post_objs {
		posts = append(posts, pobj.(*Post))
	}
	return posts
}

func (s *PostService) MemberThreadPostActionTags(threadId int64, uid int64) map[string]POST_ACTIONTAG {
	maps := make(map[string]POST_ACTIONTAG)
	if uid <= 0 {
		return maps
	}
	user_action_key := fmt.Sprintf(thread_post_user_ding_hashtbl, threadId, uid)
	kvs, err := ssdb.New(use_ssdb_group_db).Hgetall(user_action_key)
	if err != nil {
		return maps
	}
	for i := 0; i < len(kvs); i += 2 {
		k := kvs[i]
		v, err := strconv.Atoi(kvs[i+1])
		if err != nil {
			continue
		}
		maps[k] = POST_ACTIONTAG(v)
	}
	return maps
}

func (s *PostService) Invisible(postId string, do bool) error {
	post := s.Get(postId)
	if post == nil {
		return fmt.Errorf("评论不存在")
	}
	o := dbs.NewOrm(db_aliasname)
	tblid := PostTableId(postId)
	tbl := fmt.Sprintf(group_post_table_fmt, tblid)
	sql := fmt.Sprintf(`update %s set invisible=? where id=?`, tbl)
	_, err := o.Raw(sql, true, postId).Exec()
	if err != nil {
		return err
	}
	post.Invisible = true
	ssdb.New(use_ssdb_group_db).Set(s.getCacheKey(postId), *post)
	return nil
}

func (s *PostService) imgToRes(imgId int64, picSizes []libs.PIC_SIZE) *ImgRes {
	if imgId <= 0 {
		return nil
	}
	file := file_storage.GetFile(imgId)
	if file == nil {
		return nil
	}
	data, err := fileData(imgId)
	if err != nil {
		fmt.Println("---------------------------------", err)
		logs.Errorf("post view pic err:%s", err.Error())
		return nil
	}
	ir := &ImgRes{}
	for _, spsize := range picSizes {
		size := 0
		switch spsize {
		case libs.PIC_SIZE_MIDDLE:
			size = group_pic_middle_w
			break
		case libs.PIC_SIZE_THUMBNAIL:
			size = group_pic_thumbnail_w
			break
		default:
			size = 0
			break
		}
		if spsize != libs.PIC_SIZE_ORIGINAL {
			fid, err := imgResize(data, file, size)
			if err != nil {
				fmt.Println("---------------------------------", err)
				continue
			}
			if spsize == libs.PIC_SIZE_MIDDLE {
				ir.BmiddleImgId = fid
			}
			if spsize == libs.PIC_SIZE_THUMBNAIL {
				ir.ThumbnailImgId = fid
			}
		} else {
			ir.OriginalImgId = imgId
		}
	}
	cache := utils.GetCache()
	cache.Set(s.getImgResCacheKey(imgId), *ir, 92*time.Hour)
	return ir
}

func (s *PostService) GetSrcToRes(resource string) *PostRes {
	var pr *PostRes
	err := json.Unmarshal([]byte(resource), &pr)
	if err == nil {
		return pr
	}
	return nil
}

func (s *PostService) MaxPosition(threadId int64) int {
	i, _ := ssdb.New(use_ssdb_group_db).Zcard(fmt.Sprintf(thread_post_sets, threadId))
	return i
}

func (s *PostService) Create(post *Post) error {
	if post.ThreadId <= 0 {
		return fmt.Errorf("帖子不存在")
	}
	if post.AuthorId <= 0 {
		return fmt.Errorf("未指定评论人")
	}
	tservice := NewThreadService(s.cfg)
	thread := tservice.Get(post.ThreadId)
	if thread == nil {
		return fmt.Errorf("帖子不存在")
	}

	var reply *Post = nil
	if len(post.ReplyId) > 0 {
		reply = s.Get(post.ReplyId)
		if reply == nil {
			return fmt.Errorf("回复的评论不存在")
		}
		post.ReplyPosition = reply.Position
		post.ReplyUid = reply.AuthorId
	}

	pr := &PostRes{}
	//处理图片数据
	if len(post.ImgIds) > 0 {
		irs := []*ImgRes{}
		var wait sync.WaitGroup
		for _, imgid := range post.ImgIds {
			_id := imgid
			wait.Add(1)
			go func() {
				defer wait.Done()
				file := file_storage.GetFile(_id)
				if file == nil {
					return
				}
				if file.ExtName == "gif" { //gif不支持
					return
				}
				ir := s.imgToRes(_id, []libs.PIC_SIZE{
					libs.PIC_SIZE_ORIGINAL,
					libs.PIC_SIZE_THUMBNAIL,
					libs.PIC_SIZE_MIDDLE,
				})
				if ir != nil {
					irs = append(irs, ir)
				}
			}()
		}
		wait.Wait()
		//原始顺序排列
		imgrs := []*ImgRes{}
		for _, imgid := range post.ImgIds {
			for _, ir := range irs {
				if ir.OriginalImgId == imgid {
					imgrs = append(imgrs, ir)
					break
				}
			}
		}
		pr.ImgResource = imgrs
	}
	//处理视频数据

	_jd, err := json.Marshal(pr)
	if err == nil {
		post.Resources = string(_jd)
	}
	post.Subject = utils.CensorWords(post.Subject) //过滤关键字
	post.Message = utils.CensorWords(post.Message) //过滤关键字
	post.DateLine = time.Now().Unix()
	post.Position = s.MaxPosition(post.ThreadId) + 1 //设置楼层
	post.Id = NewPostId(thread.PostTableId)          //设置id
	if len(pr.ImgResource) > 0 {                     //头图片,默认第一张
		post.Img = pr.ImgResource[0].ThumbnailImgId
	}

	tbl := fmt.Sprintf(group_post_table_fmt, thread.PostTableId)
	o := dbs.NewOrm(db_aliasname)
	sql := fmt.Sprintf(`insert %s(id,tid,authorid,subject,dateline,message,ip,invisible,ding,cai,position,
		replyid,replyuid,replyposition,img,resources,longitude,latitude,fromdev)
		values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, tbl)
	_, err = o.Raw(sql, post.Id,
		post.ThreadId,
		post.AuthorId,
		post.Subject,
		post.DateLine,
		post.Message,
		post.Ip,
		post.Invisible,
		post.Ding,
		post.Cai,
		post.Position,
		post.ReplyId,
		post.ReplyUid,
		post.ReplyPosition,
		post.Img,
		post.Resources,
		post.LongiTude,
		post.LatiTude,
		post.FromDevice).Exec()
	if err != nil {
		return fmt.Errorf("发表评论失败:001")
	}
	cclient := ssdb.New(use_ssdb_group_db)
	cclient.Zadd(fmt.Sprintf(thread_post_sets, post.ThreadId), post.Id, int64(post.Position))
	cclient.Set(s.getCacheKey(post.Id), post)

	//回复
	if reply != nil && reply.AuthorId != post.AuthorId {
		go func() {
			message.SendMsgV2(post.AuthorId, reply.AuthorId, MSG_TYPE_REPLY, post.Message, post.Id, nil)
		}()
	}

	//@消息
	atNames := utils.ExtractAts(post.Message)
	if len(atNames) > 0 {
		go func() {
			ps := passport.NewMemberProvider()
			for _, atName := range atNames {
				atuid := ps.GetUidByNickname(atName)
				if atuid > 0 {
					message.SendMsgV2(post.AuthorId, atuid, MSG_TYPE_MESSAGE, thread.Subject, post.Id, nil)
				}
			}
		}()
	}

	if post.Position > 1 {
		go tservice.ActionProperty([]TH_PROPERTY{TH_PROPERTY_REPLIES}, []int64{1}, post.ThreadId)
		go tservice.SetLastPost(post)
	}

	return nil
}

func (s *PostService) Gets(threadId int64, page int, size int, orderby POST_ORDERBY, onlylz bool) (int, []*Post) {
	tservice := NewThreadService(s.cfg)
	thread := tservice.Get(threadId)
	posts := []*Post{}
	if thread == nil {
		return 0, posts
	}
	type Result struct {
		Id string
	}
	offset := (page - 1) * size

	tbl := fmt.Sprintf(group_post_table_fmt, thread.PostTableId)
	o := dbs.NewOrm(db_aliasname)
	total := 0
	if onlylz {
		var rs []Result
		sql := fmt.Sprintf("SELECT id FROM %s where tid=? and authorid=%d and position > 1 order by position %s limit %d,%d", tbl, thread.AuthorId, orderby, offset, size)
		_, err := o.Raw(sql, threadId).QueryRows(&rs)
		if err != nil {
			return 0, posts
		}

		//总数
		sql = fmt.Sprintf("SELECT count(id) as counts FROM %s where tid=? and authorid=%d and position > 1", tbl, thread.AuthorId)
		type TotalResult struct {
			Counts int
		}
		var tr TotalResult
		o.Raw(sql, threadId).QueryRow(&tr)
		total = tr.Counts

		keys := make([]string, len(rs), len(rs))
		for i, r := range rs {
			keys[i] = s.getCacheKey(r.Id)
		}
		objs := ssdb.New(use_ssdb_group_db).MultiGet(keys, reflect.TypeOf(Post{}))
		for _, t := range objs {
			posts = append(posts, t.(*Post))
		}
	} else {
		var rs []Result
		sql := fmt.Sprintf("SELECT id FROM %s where tid=? and position > 1 order by position %s limit %d,%d", tbl, orderby, offset, size)
		_, err := o.Raw(sql, threadId).QueryRows(&rs)
		if err != nil {
			return 0, posts
		}

		total, _ = ssdb.New(use_ssdb_group_db).Zcard(fmt.Sprintf(thread_post_sets, threadId))

		keys := make([]string, len(rs), len(rs))
		for i, r := range rs {
			keys[i] = s.getCacheKey(r.Id)
		}
		objs := ssdb.New(use_ssdb_group_db).MultiGet(keys, reflect.TypeOf(Post{}))
		for _, t := range objs {
			posts = append(posts, t.(*Post))
		}
	}
	return total, posts
}

////////////////////////////////////////////////////////////////////////////////
// 举报
////////////////////////////////////////////////////////////////////////////////
type ReportService struct{}

func NewReportService() *ReportService {
	return &ReportService{}
}

func (r *ReportService) Create(report *Report) error {
	o := dbs.NewOrm(db_aliasname)
	_, err := o.Insert(report)
	return err
}

func (r *ReportService) Gets(page int, size int, c REPORT_CATEGORY) ([]*Report, int) {
	o := dbs.NewOrm(db_aliasname)
	query := o.QueryTable(&Report{})
	if int(c) > 0 {
		query.Filter("c", int(c))
	}
	total, _ := query.Count()
	offset := (page - 1) * size
	var lst []*Report
	query.OrderBy("-ts").Limit(size, offset).All(&lst)
	return lst, int(total)
}
