package groups

import (
	"dbs"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

var db_aliasname string
var use_ssdb_group_db, credit_service_host string

func init() {
	//ssdb tag
	use_ssdb_group_db = beego.AppConfig.String("group.ssdb.db")
	//积分系统地址
	credit_service_host = beego.AppConfig.String("group.credit.host")

	orm.RegisterModel(new(GroupCfg), new(Group), new(Thread), new(Post),
		new(Report),
		new(GroupMemberTable), new(MemberGroupTable), new(PostTable))
	register_db()

	//启动前加载
	beego.AddAPPStartHook(func() error {
		load_tbls() //加载分表数据
		return nil
	})
}

func register_db() {
	db_aliasname = beego.AppConfig.String("group.db.aliasname")
	if len(db_aliasname) == 0 {
		return
	}
	db_user := beego.AppConfig.String("group.db.user")
	db_pwd := beego.AppConfig.String("group.db.pwd")
	db_host := beego.AppConfig.String("group.db.host")
	db_port, _ := beego.AppConfig.Int("group.db.port")
	db_name := beego.AppConfig.String("group.db.name")
	db_charset := beego.AppConfig.String("group.db.charset")
	db_protocol := beego.AppConfig.String("group.db.protocol")
	db_time_local := beego.AppConfig.String("group.db.time_local")
	db_maxconns, _ := beego.AppConfig.Int("group.db.maxconns")
	db_maxidels, _ := beego.AppConfig.Int("group.db.maxidles")
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

func load_tbls() {
	o := dbs.NewOrm(db_aliasname)
	var gmtbls []*GroupMemberTable
	o.QueryTable(&GroupMemberTable{}).All(&gmtbls)
	tbl_mutex.Lock()
	defer tbl_mutex.Unlock()
	for _, tbl := range gmtbls {
		gmTbls[int(tbl.Id)] = tbl.TblName
	}
}