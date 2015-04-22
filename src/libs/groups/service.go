package groups

import (
	"dbs"
	"fmt"
	credit_client "libs/credits/client"
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
	return fmt.Sprintf("group_config_%d", id)
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
	group_name_sets = "group_name_sets"
)

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

func (s *GroupService) setGroupTableId(group *Group) error {
	group.ThreadTableId = 1 //默认,数据量大后分表
	group.
	return nil
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
	s.setGroupTableId(group)

	o := dbs.NewOrm(db_aliasname)
	id, err := o.Insert(group)
	if err != nil {
		return fmt.Errorf("建组失败:001")
	}
	group.Id = id
	cache := utils.GetCache()
	cache.Add(s.getGroupCacheKey(id), *group, 12*time.Hour)
	return nil
}

func (s *GroupService) GetGroup(groupId int64) *Group {
	cache := utils.GetCache()
	group := Group{}
	err := cache.Get(s.getGroupCacheKey(groupId), &group)
	if err != nil {
		return nil
	}
	group.Id = groupId
	o := dbs.NewOrm(db_aliasname)
	err = o.Read(&group)
	if err != nil {
		return nil
	}
	cache.Add(s.getGroupCacheKey(groupId), group, 12*time.Hour)
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
	cache.Replace(s.getGroupCacheKey(group.Id), *group, 12*time.Hour)
	return nil
}

func (s *GroupService) updateGroupProperty(groupId int64, column string, val interface{}) error {
	o := dbs.NewOrm(db_aliasname)
	_, err := o.QueryTable(&Group{}).Filter("id", groupId).Update(orm.Params{
		column: val,
	})
	return err
}

func (s *GroupService) UpdateGroupMemberCount(groupId int64, nums int) error {
	err := s.updateGroupProperty(groupId, "members", orm.ColValue(orm.Col_Add, nums))
	if err != nil {
		return err
	}
	gp := s.GetGroup(groupId)
	if gp != nil {
		gp.Members += nums
		cache := utils.GetCache()
		cache.Replace(s.getGroupCacheKey(groupId), *gp, 12*time.Hour)
	}
	return nil
}

func (s *GroupService) UpdateGroupStatus(groupId int64, status GROUP_STATUS) error {
	err := s.updateGroupProperty(groupId, "status", status)
	if err != nil {
		return err
	}
	gp := s.GetGroup(groupId)
	if gp != nil {
		gp.Status = status
		cache := utils.GetCache()
		cache.Replace(s.getGroupCacheKey(groupId), *gp, 12*time.Hour)
	}
	return nil
}

func (s *GroupService) UpdateGroupRecommend(groupId int64, recommend bool) error {
	err := s.updateGroupProperty(groupId, "recommend", recommend)
	if err != nil {
		return err
	}
	gp := s.GetGroup(groupId)
	if gp != nil {
		gp.Recommend = recommend
		cache := utils.GetCache()
		cache.Replace(s.getGroupCacheKey(groupId), *gp, 12*time.Hour)
	}
	return nil
}
