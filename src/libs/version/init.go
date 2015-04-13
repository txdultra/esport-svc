package version

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(
		new(ClientVersion),
	)
}
