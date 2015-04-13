package feedback

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(
		//直播
		new(Feedback),
	)
}
