package credits

import (
	"dbs"
	"fmt"
	"regexp"
	"strconv"
	"sync"
	"time"
	"utils"
)

var creditRecordTbls = make(map[string]bool)
var n_locker *sync.Mutex = new(sync.Mutex)

type Credit struct {
	Uid              int64  `orm:"column(uid)"`
	TotalPoints      int64  `orm:"column(total_points)"`
	RemainingPoints  int64  `orm:"column(rem_points)"`
	LockedIncrPoints int64  `orm:"column(lockincr_points)"`
	LockedDecrPoints int64  `orm:"column(lockdecr_points)"`
	LastNo           string `orm:"column(last_no)"`
	LastTs           int64  `orm:"column(last_ts)"`
}

type CreditRecord struct {
	No          string                  `orm:"column(no)"`
	Uid         int64                   `orm:"column(uid)"`
	Changes     int64                   `orm:"column(changes)"`
	Points      int64                   `orm:"column(points)"`
	PreNo       string                  `orm:"column(pre_no)"`
	PrePoints   int64                   `orm:"column(pre_points)"`
	Description string                  `orm:"column(des)"`
	Act         CREDIT_OPERATION_ACTION `orm:"column(act)"`
	Ref         string                  `orm:"column(ref)"`
	RefId       string                  `orm:"column(ref_id)"`
	Ts          int64                   `orm:"column(ts)"`
	State       CREDIT_RECORD_STATE     `orm:"column(state)"`
}

type CREDIT_RECORD_STATE string

const (
	CREDIT_RECORD_STATE_COMPLETED CREDIT_RECORD_STATE = "completed"
	CREDIT_RECORD_STATE_RBLOCKED  CREDIT_RECORD_STATE = "rblocked"
	CREDIT_RECORD_STATE_LOCKED    CREDIT_RECORD_STATE = "locked"
)

type CREDIT_OPERATION_ACTION string

const (
	CREDIT_OPERATION_LOCKINCR     CREDIT_OPERATION_ACTION = "lockincr"
	CREDIT_OPERATION_LOCKDECR     CREDIT_OPERATION_ACTION = "lockdecr"
	CREDIT_OPERATION_ROLLBACKLOCK CREDIT_OPERATION_ACTION = "rblock"
	CREDIT_OPERATION_UNLOCK       CREDIT_OPERATION_ACTION = "unlock"
	CREDIT_OPERATION_INCR         CREDIT_OPERATION_ACTION = "incr"
	CREDIT_OPERATION_DECR         CREDIT_OPERATION_ACTION = "decr"
)

type OperationCreditParameter struct {
	No        string
	Uid       int64
	Points    uint64
	Desc      string
	Operation CREDIT_OPERATION_ACTION
	Ref       string
	RefId     string
}

type OPERATION_CREDIT_STATE string

const (
	OPERATION_CREDIT_STATE_UNDEFINED OPERATION_CREDIT_STATE = "undefined"
	OPERATION_CREDIT_STATE_UNDER     OPERATION_CREDIT_STATE = "under"
	OPERATION_CREDIT_STATE_PARAMFAIL OPERATION_CREDIT_STATE = "paramfail"
	OPERATION_CREDIT_STATE_NOFAIL    OPERATION_CREDIT_STATE = "nofail"
	OPERATION_CREDIT_STATE_ERROR     OPERATION_CREDIT_STATE = "error"
	OPERATION_CREDIT_STATE_COMPLETED OPERATION_CREDIT_STATE = "completed"
	OPERATION_CREDIT_STATE_SYSBUSY   OPERATION_CREDIT_STATE = "sysbusy"
	OPERATION_CREDIT_STATE_OPNOEXIST OPERATION_CREDIT_STATE = "opnoexist"
	OPERATION_CREDIT_STATE_SUCCESS   OPERATION_CREDIT_STATE = "success"
)

type Result struct {
	No              string                 //订单号
	RemainingPoints int64                  //剩余积分
	State           OPERATION_CREDIT_STATE //操作状态
	Error           error                  //错误
}

//积分日志记录分表
func credit_records_tbl_tag(uid int64) string {
	mod := uid % 999
	mtag := ""
	if mod > 99 {
		mtag = fmt.Sprintf("%d", mod)
	} else if mod < 10 {
		mtag = fmt.Sprintf("00%d", mod)
	} else {
		mtag = fmt.Sprintf("0%d", mod)
	}
	return mtag
}

func hash_credit_records_tbl(uid int64) string {
	mtag := credit_records_tbl_tag(uid)
	tbl := credit_records_tbl_bytag(mtag)
	if _, ok := creditRecordTbls[tbl]; ok {
		return tbl
	}
	return credit_create_tbl(tbl)
}

func credit_records_tbl_bytag(tag string) string {
	return fmt.Sprintf("%s_%s", credit_records_tbl_pfx, tag)
}

func credit_create_tbl(tbl string) string {
	n_locker.Lock()
	defer n_locker.Unlock()
	if _, ok := creditRecordTbls[tbl]; ok {
		return tbl
	}
	o := dbs.NewOrm(credit_db)
	create_tbl_sql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(no char(32) NOT NULL,
	  uid int(11) NOT NULL,
	  changes int(11) NOT NULL,
	  points int(11) NOT NULL,
	  pre_no char(32) NOT NULL,
	  pre_points int(11) NOT NULL,
	  des varchar(200) NOT NULL,
	  act char(10) NOT NULL,
	  ref char(10) NOT NULL,
	  ref_id varchar(30) NOT NULL,
	  ts bigint(15) NOT NULL,
	  state char(10) NOT NULL DEFAULT 'completed',
	  PRIMARY KEY (no),
	  KEY idx_uid_ts (uid,ts) USING BTREE) ENGINE=InnoDB DEFAULT CHARSET=utf8`, tbl)
	_, err := o.Raw(create_tbl_sql).Exec()
	if err == nil {
		creditRecordTbls[tbl] = true
	}
	return tbl
}

func credit_record_tbl_byno(no string) (tbl string, uid int64, err error) {
	re := regexp.MustCompile(`[0-9]+`)
	us := re.FindAllString(no, -1)
	if len(us) != 3 {
		return "", 0, fmt.Errorf("非法编号")
	}
	_uid, _err := strconv.ParseInt(us[1], 10, 64)
	if _err != nil {
		return "", 0, fmt.Errorf("非法编号")
	}
	return credit_records_tbl_bytag(us[0]), _uid, nil
}

func credit_record_no(uid int64) string {
	tag := credit_records_tbl_tag(uid)
	utag := fmt.Sprintf("%010d", uid)
	return fmt.Sprintf("%s%s%sN%d", tag, NO_pfx, utag, utils.TimeMillisecond(time.Now()))
}
