package dbs

import (
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var mongodb_addrs string
var mongodb_session_consistency string
var mongodb_session_refresh bool

func init() {
	// 设置为 UTC 时间
	orm.DefaultTimeLoc = time.Local
	orm.Debug = true
	orm.DefaultRowsLimit = 10000
	// set default database
	// read app.conf parameters

	default_db_aliasname := beego.AppConfig.String("db.default.aliasname")
	if len(default_db_aliasname) == 0 {
		return
	}
	default_db_user := beego.AppConfig.String("db.default.user")
	default_db_pwd := beego.AppConfig.String("db.default.pwd")
	default_db_host := beego.AppConfig.String("db.default.host")
	default_db_port, _ := beego.AppConfig.Int("db.default.port")
	default_db_name := beego.AppConfig.String("db.default.name")
	default_db_charset := beego.AppConfig.String("db.default.charset")
	default_db_protocol := beego.AppConfig.String("db.default.protocol")
	default_db_time_local := beego.AppConfig.String("db.default.time_local")
	default_db_maxconns, _ := beego.AppConfig.Int("db.default.maxconns")
	default_db_maxidels, _ := beego.AppConfig.Int("db.default.maxidles")
	if default_db_maxconns <= 0 {
		default_db_maxconns = 1500
	}
	if default_db_maxidels <= 0 {
		default_db_maxidels = 100
	}
	//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	default_db_addr := default_db_host + ":" + strconv.Itoa(default_db_port)

	default_connection_url := default_db_user + ":" + default_db_pwd + "@" + default_db_protocol + "(" + default_db_addr + ")/" + default_db_name + "?charset=" + default_db_charset + "&loc=" + default_db_time_local
	orm.RegisterDataBase(default_db_aliasname, "mysql", default_connection_url, default_db_maxidels, default_db_maxconns)

	//monogodb
	mongodb_addrs = beego.AppConfig.String("mongodb.addrs")
	mongodb_session_consistency = beego.AppConfig.String("mongodb.session.consistency")
	mongodb_session_refresh, _ = beego.AppConfig.Bool("mongodb.session.refresh")

	//uc
	if ok, _ := beego.AppConfig.Bool("db.uc.open"); ok {
		uc_db_user := beego.AppConfig.String("db.uc.user")
		uc_db_pwd := beego.AppConfig.String("db.uc.pwd")
		uc_db_host := beego.AppConfig.String("db.uc.host")
		uc_db_port, _ := beego.AppConfig.Int("db.uc.port")
		uc_db_name := beego.AppConfig.String("db.uc.name")
		uc_db_charset := beego.AppConfig.String("db.uc.charset")
		uc_db_protocol := beego.AppConfig.String("db.uc.protocol")
		uc_db_addr := uc_db_host + ":" + strconv.Itoa(uc_db_port)
		uc_connection_url := uc_db_user + ":" + uc_db_pwd + "@" + uc_db_protocol + "(" + uc_db_addr + ")/" + uc_db_name + "?charset=" + uc_db_charset
		orm.RegisterDataBase("uc", "mysql", uc_connection_url, default_db_maxidels)
	}
}

func NewDefaultOrm() orm.Ormer {
	return orm.NewOrm()
}

func NewOrm(dbalias string) orm.Ormer {
	o := orm.NewOrm()
	o.Using(dbalias)
	return o
}

func NewUcOrm() orm.Ormer {
	o := orm.NewOrm()
	o.Using("uc")
	return o
}

func NewAuthOrm() orm.Ormer {
	o := orm.NewOrm()
	o.Using("auth")
	return o
}

func LoadDb(aliasName, providerName, connectionUrl string, maxIdles int, maxConns int) {
	orm.RegisterDataBase(aliasName, providerName, connectionUrl, maxIdles, maxConns)
}
