package groups

import (
	"dbs"
	"fmt"
	"libs"
	"strconv"
	"time"
)

const (
	tjKind = "group_status"
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
			tj.SetJob(tjName(group.Id), tjKind, strconv.FormatInt(group.Id, 10), 1000, time.Unix(group.EndTime, 0))
		}
	}
	go tjWorker()
}

func tjWorker() {
	tj := libs.NewTimerJobSys(group_timerjob_host)
	tj.RunGrabJob(tjKind, tjHandler)
}

func tjHandler(jobName, args string, otherArgs ...interface{}) (schedLater int, err error) {
	return 0, nil
}
