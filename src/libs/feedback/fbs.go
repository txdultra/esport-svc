package feedback

import (
	"dbs"
	"fmt"
	"time"
	"utils"
)

func Create(fb *Feedback) error {
	_tmp := utils.StripSQLInjection(fb.Content)
	if _tmp != fb.Content {
		return fmt.Errorf("内容包含非法字符")
	}
	fb.PostTime = time.Now()
	o := dbs.NewDefaultOrm()
	_, err := o.Insert(fb)
	return err
}

func Gets(page int, size int) (int, []*Feedback) {
	o := dbs.NewDefaultOrm()
	offset := (page - 1) * size
	count, _ := o.QueryTable(Feedback{}).Count()
	var fds []*Feedback
	o.QueryTable(Feedback{}).OrderBy("-post_time").Limit(size, offset).All(&fds)
	return int(count), fds
}
