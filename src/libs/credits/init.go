package credits

import (
	"dbs"
	"strconv"

	"github.com/astaxie/beego"
)

var use_ssdb_credit_db, credit_db, credit_records_tbl_pfx string

func init() {
	//ssdb tag
	use_ssdb_credit_db = beego.AppConfig.String("credit.ssdb.db")
	//db
	register_credit_db()
}

func register_credit_db() {
	credit_db = beego.AppConfig.String("credit.db.aliasname")
	if len(credit_db) == 0 {
		return
	}
	credit_records_tbl_pfx = beego.AppConfig.String("credit.db.tblpfx")
	db_user := beego.AppConfig.String("credit.db.user")
	db_pwd := beego.AppConfig.String("credit.db.pwd")
	db_host := beego.AppConfig.String("credit.db.host")
	db_port, _ := beego.AppConfig.Int("credit.db.port")
	db_name := beego.AppConfig.String("credit.db.name")
	db_charset := beego.AppConfig.String("credit.db.charset")
	db_protocol := beego.AppConfig.String("credit.db.protocol")
	db_time_local := beego.AppConfig.String("credit.db.time_local")
	db_maxconns, _ := beego.AppConfig.Int("credit.db.maxconns")
	db_maxidels, _ := beego.AppConfig.Int("credit.db.maxidles")
	if db_maxconns <= 0 {
		db_maxconns = 1500
	}
	if db_maxidels <= 0 {
		db_maxidels = 100
	}

	db_addr := db_host + ":" + strconv.Itoa(db_port)
	connection_url := db_user + ":" + db_pwd + "@" + db_protocol + "(" + db_addr + ")/" + db_name + "?charset=" + db_charset + "&loc=" + db_time_local
	dbs.LoadDb(credit_db, "mysql", connection_url, db_maxidels, db_maxconns)
}
