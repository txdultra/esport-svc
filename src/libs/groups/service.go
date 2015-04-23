package groups

import (
	"dbs"
	"fmt"
	credit_client "libs/credits/client"
	"logs"
	"strconv"
	"sync"
	"time"
	"utils"
	"utils/cache"
	"utils/ssdb"

	"github.com/astaxie/beego/orm"
)

////////////////////////////////////////////////////////////////////////////////
//组设定配置
////////////////////////////////////////////////////////////////////////////////
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

func UpdateGroupCfg(cfg *GroupCfg) error {
	o := dbs.NewOrm(db_aliasname)
	_, err := o.Update(cfg)
	if err != nil {
		return err
	}
	cache.Replace(getCfgCacheKey(cfg.Id), *cfg, 24*time.Hour)
	return nil
}

////////////////////////////////////////////////////////////////////////////////
//组服务
////////////////////////////////////////////////////////////////////////////////
const (
	group_members_count_field   = "group_members_count"
	group_name_sets             = "group.name_sets"
	group_member_joined_sortset = "group.member_%d_joined"
	group_joined_member_sortset = "group.joined_%d_uids"
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
	cfg *GroupCfg
}

func (s *GroupService) getGroupCacheKey(id int64) string {
	return fmt.Sprintf("group_model_%d", id)
}

func (s *GroupService) VerifyNewGroup(group *Group) error {
	if s.cfg == nil {
		return fmt.Errorf("未设置配置对象")
	}
	if len(group.Name) > s.cfg.GroupNameLen {
		return fmt.Errorf(fmt.Sprintf("组名称不能大于%d个字符", s.cfg.GroupNameLen))
	}
	if len(group.Description) > s.cfg.GroupDescLen {
		return fmt.Errorf(fmt.Sprintf("组描述内容不能大于%d个字符", s.cfg.GroupDescLen))
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
	//验证是否同名
	has, err := ssdb.New(use_ssdb_group_db).Hexists(group_name_sets, group.Name)
	if !has || err != nil {
		return fmt.Errorf("组名称已存在")
	}

	//验证积分
	if group.Belong == GROUP_BELONG_MEMBER {
		var needPoint int64 = 0 //需要积分
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
	  PRIMARY KEY (groupid,uid)) ENGINE=InnoDB DEFAULT CHARSET=utf8`, tbl)
	o.Raw(create_tbl_sql).Exec()
	gmTbls[i_tag] = tbl
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
	  PRIMARY KEY (uid,groupid)) ENGINE=InnoDB DEFAULT CHARSET=utf8`, tbl)
	o.Raw(create_tbl_sql).Exec()
	mgTbls[i_tag] = tbl
	return i_tag, tbl, nil
}

func (s *GroupService) setGroupTableId(group *Group) error {
	group.ThreadTableId = 1 //默认,数据量大后分表
	tblId, _, err := s.GetGMTable(group.Uid)
	if tblId > 0 {
		group.MembersTableId = tblId
	}
	return err
}

func (s *GroupService) CreateGroup(group *Group) error {
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

	o := dbs.NewOrm(db_aliasname)
	id, err := o.Insert(group)
	if err != nil {
		return fmt.Errorf("建组失败:002")
	}
	group.Id = id
	cache := utils.GetCache()
	cache.Add(s.getGroupCacheKey(id), *group, 1*time.Hour)
	return nil
}

func (s *GroupService) GetGroup(groupId int64) *Group {
	cache := utils.GetCache()
	group := Group{}
	err := cache.Get(s.getGroupCacheKey(groupId), &group)
	if err == nil {
		return &group
	}
	group.Id = groupId
	o := dbs.NewOrm(db_aliasname)
	err = o.Read(&group)
	if err != nil {
		return nil
	}
	cache.Add(s.getGroupCacheKey(groupId), group, 1*time.Hour)
	return &group
}

func (s *GroupService) UpdateGroup(group *Group) error {
	o := dbs.NewOrm(db_aliasname)
	_, err := o.Update(group, "groupname", "description", "uid", "country", "city", "gameids", "displayorder", "img", "bgimg",
		"belong", "type", "searchkeyword")
	if err != nil {
		return err
	}
	cache := utils.GetCache()
	cache.Replace(s.getGroupCacheKey(group.Id), *group, 1*time.Hour)
	return nil
}

func (s *GroupService) doGroupCounts(groupId int64, fields []string, incrs []int) {
	group := s.GetGroup(groupId)
	if group == nil {
		return
	}
	if len(fields) != len(incrs) {
		panic("fields和incrs的数量必须一致")
	}
	params := make(orm.Params)
	for i, field := range fields {
		switch field {
		case group_members_count_field:
			group.Members += incrs[i]
			params["members"] = orm.ColValue(orm.Col_Add, incrs[i])
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
	cache.Replace(s.getGroupCacheKey(group.Id), *group, 1*time.Hour)
}

func (s *GroupService) MemberJoinGroup(uid int64, groupId int64) error {
	group := s.GetGroup(groupId)
	if group == nil {
		return fmt.Errorf("小组不存在")
	}
	if group.Status == GROUP_STATUS_CLOSED {
		return fmt.Errorf("小组已关闭,不允许加入")
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
	sql = fmt.Sprint("insert %s(uid,groupid,ts) values(?,?,?)", mgTbl)
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

	s.doGroupCounts(groupId, []string{group_members_count_field}, []int{1})

	return nil
}

func (s *GroupService) MemberExitGroup(uid int64, groupId int64) error {
	group := s.GetGroup(groupId)
	if group == nil {
		return fmt.Errorf("小组不存在")
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
	sql = fmt.Sprint("delete from %s where uid=? and groupid=?", mgTbl)
	o.Raw(sql, uid, groupId).Exec()

	ssdb_client := ssdb.New(use_ssdb_group_db)
	//用户离开的小组
	ssdb_client.Zrem(mjoined_key, groupId)
	//小组中去除用户
	ssdb_client.Zrem(gjoined_key, uid)

	s.doGroupCounts(groupId, []string{group_members_count_field}, []int{-1})

	return nil
}

func (s *GroupService) UpdateMemberLastActionGroup(groupId int64, uid int64, t time.Time) {
	mjoined_key := fmt.Sprintf(group_member_joined_sortset, uid) //用户加入小组的集合key
	ts, err := ssdb.New(use_ssdb_group_db).Zscore(mjoined_key, groupId)
	if ts == 0 || err != nil {
		return
	}
	if ts > t.Unix() {
		return
	}
	gap := t.Unix() - ts
	ssdb.New(use_ssdb_group_db).Zincrby(mjoined_key, groupId, gap)
}

//func (s *GroupService) updateGroupProperty(groupId int64, column string, val interface{}) error {
//	o := dbs.NewOrm(db_aliasname)
//	_, err := o.QueryTable(&Group{}).Filter("id", groupId).Update(orm.Params{
//		column: val,
//	})
//	return err
//}

//func (s *GroupService) UpdateGroupMemberCount(groupId int64, nums int) error {
//	err := s.updateGroupProperty(groupId, "members", orm.ColValue(orm.Col_Add, nums))
//	if err != nil {
//		return err
//	}
//	gp := s.GetGroup(groupId)
//	if gp != nil {
//		gp.Members += nums
//		cache := utils.GetCache()
//		cache.Replace(s.getGroupCacheKey(groupId), *gp, 12*time.Hour)
//	}
//	return nil
//}

//func (s *GroupService) UpdateGroupStatus(groupId int64, status GROUP_STATUS) error {
//	err := s.updateGroupProperty(groupId, "status", status)
//	if err != nil {
//		return err
//	}
//	gp := s.GetGroup(groupId)
//	if gp != nil {
//		gp.Status = status
//		cache := utils.GetCache()
//		cache.Replace(s.getGroupCacheKey(groupId), *gp, 12*time.Hour)
//	}
//	return nil
//}

//func (s *GroupService) UpdateGroupRecommend(groupId int64, recommend bool) error {
//	err := s.updateGroupProperty(groupId, "recommend", recommend)
//	if err != nil {
//		return err
//	}
//	gp := s.GetGroup(groupId)
//	if gp != nil {
//		gp.Recommend = recommend
//		cache := utils.GetCache()
//		cache.Replace(s.getGroupCacheKey(groupId), *gp, 12*time.Hour)
//	}
//	return nil
//}
