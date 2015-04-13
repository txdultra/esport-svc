package shop

import (
	"dbs"
	"strconv"
	"sync"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

var db_aliasname string
var provinces map[string]*Province = make(map[string]*Province)
var citys map[string][]*City = make(map[string][]*City)
var areas map[string][]*Area = make(map[string][]*Area)
var once *sync.Once = new(sync.Once)

func init() {
	orm.RegisterModel(
		new(Item),
		new(Order),
		new(OrderIncerment),
		new(OrderTransport),
		new(ItemCode),
		new(OrderItemSnap),
		new(Province),
		new(City),
		new(Area),
	)
	register_db()
	//初始化地区
	init_areas()
}

func init_areas() {
	//启动前加载
	beego.AddAPPStartHook(func() error {
		once.Do(func() {
			o := dbs.NewOrm(db_aliasname)
			var _ps []*Province
			var _cs []*City
			var _as []*Area
			o.QueryTable(&Province{}).All(&_ps)
			o.QueryTable(&City{}).All(&_cs)
			o.QueryTable(&Area{}).All(&_as)

			for _, p := range _ps {
				provinces[p.ProvinceID] = p
			}
			for _, c := range _cs {
				if ls, ok := citys[c.Father]; ok {
					ls = append(ls, c)
					citys[c.Father] = ls
				} else {
					ls = []*City{}
					ls = append(ls, c)
					citys[c.Father] = ls
				}
			}
			for _, a := range _as {
				if as, ok := areas[a.Father]; ok {
					as = append(as, a)
					areas[a.Father] = as
				} else {
					as = []*Area{}
					as = append(as, a)
					areas[a.Father] = as
				}
			}
		})
		return nil
	})

}

func register_db() {
	db_aliasname = beego.AppConfig.String("shop.db.aliasname")
	if len(db_aliasname) == 0 {
		return
	}
	db_user := beego.AppConfig.String("shop.db.user")
	db_pwd := beego.AppConfig.String("shop.db.pwd")
	db_host := beego.AppConfig.String("shop.db.host")
	db_port, _ := beego.AppConfig.Int("shop.db.port")
	db_name := beego.AppConfig.String("shop.db.name")
	db_charset := beego.AppConfig.String("shop.db.charset")
	db_protocol := beego.AppConfig.String("shop.db.protocol")
	db_time_local := beego.AppConfig.String("shop.db.time_local")
	db_maxconns, _ := beego.AppConfig.Int("shop.db.maxconns")
	db_maxidels, _ := beego.AppConfig.Int("shop.db.maxidles")
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
