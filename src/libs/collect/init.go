package collect

import (
	"github.com/astaxie/beego"
)

var collect_db, collect_collection, use_ssdb_collect_db string

func init() {
	collect_db = beego.AppConfig.String("collect.db")
	collect_collection = beego.AppConfig.String("collect.collection")
	use_ssdb_collect_db = beego.AppConfig.String("ssdb.collect.db")
	if len(collect_db) == 0 {
		panic("未配置collect.db参数")
	}
	if len(collect_collection) == 0 {
		panic("未配置collect.collection参数")
	}
	if len(use_ssdb_collect_db) == 0 {
		panic("未配置ssdb.collect.db参数")
	}
}

type COLLECT_RELTYPE string

const (
	COLLECT_RELTYPE_VOD   COLLECT_RELTYPE = "vod"
	COLLECT_RELTYPE_PLIVE COLLECT_RELTYPE = "p_live"
	COLLECT_RELTYPE_JLIVE COLLECT_RELTYPE = "j_live"
)
