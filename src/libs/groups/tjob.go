package groups

import (
	"dbs"
	"fmt"
	"libs"
	"strconv"
	"time"
)

const (
	tjKind = "group_status_job"
)

func tjName(groupId int64) string {
	return fmt.Sprintf("group_tj_%d", groupId)
}

func tjInit() {
	o := dbs.NewOrm(db_aliasname)
	var groups []*Group
	o.QueryTable(&Group{}).Filter("status__in", int(GROUP_STATUS_RECRUITING), int(GROUP_STATUS_LOWMEMBER)).All(&groups)
	if len(groups) > 0 {
		tj := libs.NewTimerJobSys(group_timerjob_host)
		for _, group := range groups {
			schdAt := group.EndTime
			if group.EndTime < time.Now().Unix() {
				schdAt = time.Now().Unix() + 60
			}
			tj.SetJob(tjName(group.Id), tjKind, strconv.FormatInt(group.Id, 10), 1000, time.Unix(schdAt, 0))
		}
	}
	go tjWorker()
}

func tjWorker() {
	tj := libs.NewTimerJobSys(group_timerjob_host)
	tj.RunGrabJob(tjKind, tjHandler)
}

func tjHandler(jobName, args string, otherArgs ...interface{}) (schedLater int, err error) {
	groupId, err := strconv.ParseInt(args, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("参数错误")
	}
	gs := NewGroupService(GetDefaultCfg())
	group := gs.Get(groupId)
	if group == nil {
		return 0, nil
	}
	if group.Belong == GROUP_BELONG_OFFICIAL {
		return 0, nil
	}
	switch group.Status {
	case GROUP_STATUS_RECRUITING:
		gs.Delete(group.Id)
		break
	case GROUP_STATUS_LOWMEMBER:
		gs.Close(group.Id)
		break
	default:
		break
	}
	return 0, nil
}

func tjSetJob(group *Group) error {
	name := tjName(group.Id)
	tj := libs.NewTimerJobSys(group_timerjob_host)
	return tj.SetJob(name, tjKind, strconv.FormatInt(group.Id, 10), 1000, time.Unix(group.EndTime, 0))
}

func tjRemoveJob(group *Group) error {
	name := tjName(group.Id)
	tj := libs.NewTimerJobSys(group_timerjob_host)
	return tj.RemoveJob(name, tjKind)
}
