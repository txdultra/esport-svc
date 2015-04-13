package comment

import (
	"strings"
	"time"
	"utils"

	"github.com/astaxie/beego"
)

var comments_pulish_interval time.Duration
var comments_db_name, use_ssdb_comment_db string
var commonts_collections map[string]string = make(map[string]string)

func init() {
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//comment 配置文件参数初始化
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////
	comments_pulish_interval_str := beego.AppConfig.String("comments.pulish.interval")
	if len(comments_pulish_interval_str) > 0 {
		comments_pulish_interval = utils.StrToDuration(comments_pulish_interval_str)
	} else {
		comments_pulish_interval = 3 * time.Second
	}
	comments_db_name = beego.AppConfig.String("comments.db.name")
	if len(comments_db_name) == 0 {
		panic("配置文件:comments.db.name 必须设置值")
	}
	commonts_collections_str := beego.AppConfig.String("commonts.collections")
	commonts_cols := strings.Split(commonts_collections_str, ";")
	for _, str := range commonts_cols {
		_nv := strings.Split(str, ":")
		if len(_nv) == 2 {
			_name := _nv[0]
			_val := _nv[1]
			commonts_collections[_name] = _val
		}
	}
	use_ssdb_comment_db = beego.AppConfig.String("ssdb.comment.db")
}

func CommentPulishInterval() time.Duration {
	return comments_pulish_interval
}

func GetCollectionName(name string) string {
	if val, ok := commonts_collections[name]; ok {
		return val
	}
	return ""
}
