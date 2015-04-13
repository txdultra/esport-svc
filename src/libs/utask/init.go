package utask

import (
	"dbs"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

var db_aliasname, db_tbl_pfx string
var use_ssdb_utask_db string
var credit_service_host string

func init() {
	orm.RegisterModel(
		new(TaskGroup),
		new(Task),
		new(TaskEvent),
	)
	register_db()

	use_ssdb_utask_db = beego.AppConfig.String("utask.ssdb.db")
	credit_service_host = beego.AppConfig.String("utask.credit.host")
}

func register_db() {
	db_aliasname = beego.AppConfig.String("utask.db.aliasname")
	if len(db_aliasname) == 0 {
		return
	}
	db_tbl_pfx = beego.AppConfig.String("utask.db.tblpfx")
	db_user := beego.AppConfig.String("utask.db.user")
	db_pwd := beego.AppConfig.String("utask.db.pwd")
	db_host := beego.AppConfig.String("utask.db.host")
	db_port, _ := beego.AppConfig.Int("utask.db.port")
	db_name := beego.AppConfig.String("utask.db.name")
	db_charset := beego.AppConfig.String("utask.db.charset")
	db_protocol := beego.AppConfig.String("utask.db.protocol")
	db_time_local := beego.AppConfig.String("utask.db.time_local")
	db_maxconns, _ := beego.AppConfig.Int("utask.db.maxconns")
	db_maxidels, _ := beego.AppConfig.Int("utask.db.maxidles")
	if db_maxconns <= 0 {
		db_maxconns = 1500
	}
	if db_maxidels <= 0 {
		db_maxidels = 100
	}

	db_addr := db_host + ":" + strconv.Itoa(db_port)
	connection_url := db_user + ":" + db_pwd + "@" + db_protocol + "(" + db_addr + ")/" + db_name + "?charset=" + db_charset + "&loc=" + db_time_local
	dbs.LoadDb(db_aliasname, "mysql", connection_url, db_maxidels, db_maxconns)
}
