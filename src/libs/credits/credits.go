package credits

import (
	"dbs"
	"fmt"
	"libs/dlock"
	"logs"
	"time"
	"utils"
	"utils/ssdb"

	"github.com/astaxie/beego/orm"
)

type CreditService struct{}

func NewCreditService() *CreditService {
	return &CreditService{}
}

func (c *CreditService) creditKey(uid int64) string {
	return fmt.Sprintf("credit_%d", uid)
}

func (c *CreditService) GetCredit(uid int64) int64 {
	var i int64 = 0
	err := ssdb.New(use_ssdb_credit_db).Get(c.creditKey(uid), &i)
	if err != nil {
		logs.Errorf("credit service get credit fail:%s", err.Error())
		return 0
	}
	return i
}

func (c *CreditService) GetCreditRecord(no string) *CreditRecord {
	record_tbl, _, err := credit_record_tbl_byno(no)
	if err != nil {
		return nil
	}
	var cr CreditRecord
	sql := fmt.Sprintf("select no,uid,changes,points,pre_no,pre_points,des,act,ref,ref_id,ts,state from %s where no=?", record_tbl)
	o := dbs.NewOrm(credit_db)
	err = o.Raw(sql, no).QueryRow(&cr)
	if err != nil {
		return nil
	}
	return &cr
}

func (c *CreditService) OperCredit(op *OperationCreditParameter) *Result {
	switch op.Operation {
	case CREDIT_OPERATION_INCR:
		return c.incrCredit(op)
	case CREDIT_OPERATION_DECR:
		return c.decrCredit(op)
	case CREDIT_OPERATION_LOCKINCR, CREDIT_OPERATION_LOCKDECR:
		return c.lockCredit(op)
	case CREDIT_OPERATION_ROLLBACKLOCK:
		return c.rblockCredit(op)
	case CREDIT_OPERATION_UNLOCK:
		return c.unlockCredit(op)
	default:
		return &Result{
			RemainingPoints: 0,
			State:           OPERATION_CREDIT_STATE_UNDEFINED,
			Error:           fmt.Errorf("未实现操作方法"),
		}
	}
}

func (c *CreditService) lockKey(uid int64) string {
	return fmt.Sprintf("u_%d", uid)
}

func (c *CreditService) incrCredit(op *OperationCreditParameter) *Result {
	if op.Points == 0 || op.Uid <= 0 {
		return &Result{
			RemainingPoints: 0,
			State:           OPERATION_CREDIT_STATE_PARAMFAIL,
			Error:           fmt.Errorf("参数错误"),
		}
	}

	//分布式锁
	lock := dlock.NewDistributedLock()
	locker, err := lock.Lock(c.lockKey(op.Uid))
	if err != nil {
		logs.Errorf("credit service increment get lock fail:%s", err.Error())
		return &Result{
			RemainingPoints: 0,
			State:           OPERATION_CREDIT_STATE_SYSBUSY,
			Error:           fmt.Errorf("系统繁忙,请稍后再试"),
		}
	}
	defer locker.Unlock()

	//处理逻辑
	record_tbl := hash_credit_records_tbl(op.Uid)
	var cr CreditRecord
	o := dbs.NewOrm(credit_db)
	sql := fmt.Sprintf("select no,uid,changes,points,pre_no,pre_points,des,act,ref,ref_id,ts from %s where uid=%d order by ts desc limit 1", record_tbl, op.Uid)
	err = o.Raw(sql).QueryRow(&cr)

	var pre_no string = ""
	var pre_points int64 = 0
	var orginal_points int64 = 0
	if err == nil {
		pre_no = cr.No
		pre_points = cr.Points
		orginal_points = cr.Points
	}
	//开启事务
	o.Begin()
	sql = fmt.Sprintf("insert %s(no,uid,changes,points,pre_no,pre_points,des,act,ref,ref_id,ts) values(?,?,?,?,?,?,?,?,?,?,?)", record_tbl)
	ts := utils.TimeMillisecond(time.Now())
	no := credit_record_no(op.Uid)
	now_points := orginal_points + int64(op.Points)
	_, err = o.Raw(sql,
		no,
		op.Uid,
		op.Points,
		now_points,
		pre_no,
		pre_points,
		op.Desc,
		string(CREDIT_OPERATION_INCR),
		op.Ref,
		op.RefId,
		ts,
	).Exec()
	if err != nil {
		o.Rollback()
		logs.Errorf("credit service increment insert %s table fail:%s", record_tbl, err.Error())
		return &Result{
			RemainingPoints: orginal_points,
			State:           OPERATION_CREDIT_STATE_ERROR,
			Error:           err,
		}
	}
	var maps []orm.Params
	num, err := o.Raw("select uid from user_credits where uid=?", op.Uid).Values(&maps)
	if err == nil && num > 0 {
		sql = fmt.Sprintf("update user_credits set total_points=total_points+%d,rem_points=%d,last_no=?,last_ts=%d where uid=?", op.Points, now_points, ts)
		o.Raw(sql, no, op.Uid).Exec()
	} else {
		sql = "insert user_credits(uid,total_points,rem_points,lockincr_points,lockdecr_points,last_no,last_ts) values(?,?,?,?,?,?,?)"
		o.Raw(sql, op.Uid, op.Points, op.Points, 0, 0, no, ts).Exec()
	}
	err = o.Commit()
	if err != nil {
		o.Rollback()
		logs.Errorf("credit service increment commit transaction fail:%s", err.Error())
		return &Result{
			RemainingPoints: orginal_points,
			State:           OPERATION_CREDIT_STATE_ERROR,
			Error:           err,
		}
	}
	//更新高速缓存
	c.updateCredit(op.Uid, now_points)
	return &Result{
		No:              no,
		RemainingPoints: now_points,
		State:           OPERATION_CREDIT_STATE_SUCCESS,
		Error:           nil,
	}
}

func (c *CreditService) decrCredit(op *OperationCreditParameter) *Result {
	if op.Points == 0 || op.Uid <= 0 {
		return &Result{
			RemainingPoints: 0,
			State:           OPERATION_CREDIT_STATE_PARAMFAIL,
			Error:           fmt.Errorf("参数错误"),
		}
	}

	//分布式锁
	lock := dlock.NewDistributedLock()
	locker, err := lock.Lock(c.lockKey(op.Uid))
	if err != nil {
		logs.Errorf("credit service increment get lock fail:%s", err.Error())
		return &Result{
			RemainingPoints: 0,
			State:           OPERATION_CREDIT_STATE_SYSBUSY,
			Error:           fmt.Errorf("系统繁忙,请稍后再试"),
		}
	}
	defer locker.Unlock()

	//处理逻辑
	record_tbl := hash_credit_records_tbl(op.Uid)
	var cr CreditRecord
	o := dbs.NewOrm(credit_db)
	sql := fmt.Sprintf("select no,uid,changes,points,pre_no,pre_points,des,act,ref,ref_id,ts from %s where uid=%d order by ts desc limit 1", record_tbl, op.Uid)
	err = o.Raw(sql).QueryRow(&cr)

	var pre_no string = ""
	var pre_points int64 = 0
	var orginal_points int64 = 0
	if err == nil {
		pre_no = cr.No
		pre_points = cr.Points
		orginal_points = cr.Points
	}
	if pre_points < int64(op.Points) {
		return &Result{
			RemainingPoints: pre_points,
			State:           OPERATION_CREDIT_STATE_UNDER,
			Error:           fmt.Errorf("积分不足"),
		}
	}
	//开启事务
	o.Begin()
	sql = fmt.Sprintf("insert %s(no,uid,changes,points,pre_no,pre_points,des,act,ref,ref_id,ts) values(?,?,?,?,?,?,?,?,?,?,?)", record_tbl)
	ts := utils.TimeMillisecond(time.Now())
	no := credit_record_no(op.Uid)
	now_points := orginal_points - int64(op.Points)
	_, err = o.Raw(sql,
		no,
		op.Uid,
		op.Points,
		now_points,
		pre_no,
		pre_points,
		op.Desc,
		string(CREDIT_OPERATION_DECR),
		op.Ref,
		op.RefId,
		ts,
	).Exec()
	if err != nil {
		o.Rollback()
		logs.Errorf("credit service decrement insert %s table fail:%s", record_tbl, err.Error())
		return &Result{
			RemainingPoints: orginal_points,
			State:           OPERATION_CREDIT_STATE_ERROR,
			Error:           err,
		}
	}
	var maps []orm.Params
	num, err := o.Raw("select uid from user_credits where uid=?", op.Uid).Values(&maps)
	if err == nil && num > 0 {
		sql = fmt.Sprintf("update user_credits set rem_points=%d,last_no=?,last_ts=%d where uid=?", now_points, ts)
		o.Raw(sql, no, op.Uid).Exec()
	} else {
		sql = "insert user_credits(uid,total_points,rem_points,lockincr_points,lockdecr_points,last_no,last_ts) values(?,?,?,?,?,?,?)"
		o.Raw(sql, op.Uid, op.Points, op.Points, 0, 0, no, ts).Exec()
	}
	err = o.Commit()
	if err != nil {
		o.Rollback()
		logs.Errorf("credit service decrement commit transaction fail:%s", err.Error())
		return &Result{
			RemainingPoints: orginal_points,
			State:           OPERATION_CREDIT_STATE_ERROR,
			Error:           err,
		}
	}
	//更新高速缓存
	c.updateCredit(op.Uid, now_points)
	return &Result{
		No:              no,
		RemainingPoints: now_points,
		State:           OPERATION_CREDIT_STATE_SUCCESS,
		Error:           nil,
	}
}

func (c *CreditService) lockCredit(op *OperationCreditParameter) *Result {
	if op.Points == 0 || op.Uid <= 0 {
		return &Result{
			RemainingPoints: 0,
			State:           OPERATION_CREDIT_STATE_PARAMFAIL,
			Error:           fmt.Errorf("参数错误"),
		}
	}

	//分布式锁
	lock := dlock.NewDistributedLock()
	locker, err := lock.Lock(c.lockKey(op.Uid))
	if err != nil {
		logs.Errorf("credit service increment get lock fail:%s", err.Error())
		return &Result{
			RemainingPoints: 0,
			State:           OPERATION_CREDIT_STATE_SYSBUSY,
			Error:           fmt.Errorf("系统繁忙,请稍后再试"),
		}
	}
	defer locker.Unlock()

	//处理逻辑
	record_tbl := hash_credit_records_tbl(op.Uid)
	var cr CreditRecord
	o := dbs.NewOrm(credit_db)
	sql := fmt.Sprintf("select no,uid,changes,points,pre_no,pre_points,des,act,ref,ref_id,ts from %s where uid=%d order by ts desc limit 1", record_tbl, op.Uid)
	err = o.Raw(sql).QueryRow(&cr)

	var pre_no string = ""
	var pre_points int64 = 0
	var orginal_points int64 = 0
	if err == nil {
		pre_no = cr.No
		pre_points = cr.Points
		orginal_points = cr.Points
	}
	if op.Operation == CREDIT_OPERATION_LOCKDECR {
		if pre_points < int64(op.Points) {
			return &Result{
				RemainingPoints: pre_points,
				State:           OPERATION_CREDIT_STATE_UNDER,
				Error:           fmt.Errorf("积分不足"),
			}
		}
	}
	//开启事务
	o.Begin()
	sql = fmt.Sprintf("insert %s(no,uid,changes,points,pre_no,pre_points,des,act,ref,ref_id,ts,state) values(?,?,?,?,?,?,?,?,?,?,?,?)", record_tbl)
	ts := utils.TimeMillisecond(time.Now())
	no := credit_record_no(op.Uid)
	now_points := orginal_points
	if op.Operation == CREDIT_OPERATION_LOCKDECR {
		now_points = orginal_points - int64(op.Points)
	}
	_, err = o.Raw(sql,
		no,
		op.Uid,
		op.Points,
		now_points,
		pre_no,
		pre_points,
		op.Desc,
		op.Operation,
		op.Ref,
		op.RefId,
		ts,
		CREDIT_RECORD_STATE_LOCKED,
	).Exec()
	if err != nil {
		o.Rollback()
		logs.Errorf("credit service lockp insert %s table fail:%s", record_tbl, err.Error())
		return &Result{
			RemainingPoints: orginal_points,
			State:           OPERATION_CREDIT_STATE_ERROR,
			Error:           err,
		}
	}
	var maps []orm.Params
	num, err := o.Raw("select uid from user_credits where uid=?", op.Uid).Values(&maps)
	if err == nil && num > 0 {
		if op.Operation == CREDIT_OPERATION_LOCKINCR {
			sql = fmt.Sprintf("update user_credits set rem_points=%d,lockincr_points=lockincr_points+%d,last_no=?,last_ts=%d where uid=?", now_points, op.Points, ts)
		} else if op.Operation == CREDIT_OPERATION_LOCKDECR {
			sql = fmt.Sprintf("update user_credits set rem_points=%d,lockdecr_points=lockdecr_points+%d,last_no=?,last_ts=%d where uid=?", now_points, op.Points, ts)
		}
		o.Raw(sql, no, op.Uid).Exec()
	} else {
		if op.Operation == CREDIT_OPERATION_LOCKINCR {
			sql = "insert user_credits(uid,total_points,rem_points,lockincr_points,lockdecr_points,last_no,last_ts) values(?,?,?,?,?,?,?)"
			o.Raw(sql, op.Uid, 0, 0, op.Points, 0, no, ts).Exec()
		} else if op.Operation == CREDIT_OPERATION_LOCKDECR {
			sql = "insert user_credits(uid,total_points,rem_points,lockincr_points,lockdecr_points,last_no,last_ts) values(?,?,?,?,?,?,?)"
			o.Raw(sql, op.Uid, 0, 0, 0, op.Points, no, ts).Exec()
		}
	}
	err = o.Commit()
	if err != nil {
		o.Rollback()
		logs.Errorf("credit service lockp commit transaction fail:%s", err.Error())
		return &Result{
			RemainingPoints: orginal_points,
			State:           OPERATION_CREDIT_STATE_ERROR,
			Error:           err,
		}
	}
	//更新高速缓存
	c.updateCredit(op.Uid, now_points)
	return &Result{
		No:              no,
		RemainingPoints: now_points,
		State:           OPERATION_CREDIT_STATE_SUCCESS,
		Error:           nil,
	}
	return nil
}

func (c *CreditService) checkUnlockCreditState(r *CreditRecord) *Result {
	if r == nil {
		return &Result{
			RemainingPoints: 0,
			State:           OPERATION_CREDIT_STATE_OPNOEXIST,
			Error:           fmt.Errorf("操作单号不存在"),
		}
	}
	if r.State != CREDIT_RECORD_STATE_LOCKED { //订单已完成
		return &Result{
			RemainingPoints: 0,
			State:           OPERATION_CREDIT_STATE_COMPLETED,
			Error:           fmt.Errorf("操作单号已完成"),
		}
	}
	return nil
}

func (c *CreditService) rblockCredit(op *OperationCreditParameter) *Result {
	if len(op.No) == 0 {
		return &Result{
			RemainingPoints: 0,
			State:           OPERATION_CREDIT_STATE_PARAMFAIL,
			Error:           fmt.Errorf("参数错误"),
		}
	}

	//订单号校验
	record_tbl, uid, err := credit_record_tbl_byno(op.No)
	if err != nil {
		return &Result{
			RemainingPoints: 0,
			State:           OPERATION_CREDIT_STATE_NOFAIL,
			Error:           fmt.Errorf("非法单号"),
		}
	}

	//分布式锁
	lock := dlock.NewDistributedLock()
	locker, err := lock.Lock(c.lockKey(uid))
	if err != nil {
		logs.Errorf("credit service rollback unlock get lock fail:%s", err.Error())
		return &Result{
			RemainingPoints: 0,
			State:           OPERATION_CREDIT_STATE_SYSBUSY,
			Error:           fmt.Errorf("系统繁忙,请稍后再试"),
		}
	}
	defer locker.Unlock()

	//校验订单状态
	lock_cr := c.GetCreditRecord(op.No)
	crlt := c.checkUnlockCreditState(lock_cr)
	if crlt != nil {
		return crlt
	}

	var last_cr CreditRecord
	sql := fmt.Sprintf("select no,uid,changes,points,pre_no,pre_points,des,act,ref,ref_id,ts from %s where uid=%d order by ts desc limit 1", record_tbl, lock_cr.Uid)
	o := dbs.NewOrm(credit_db)
	err = o.Raw(sql).QueryRow(&last_cr)

	var pre_no string = ""
	var pre_points int64 = 0
	var orginal_points int64 = 0
	if err == nil {
		pre_no = last_cr.No
		pre_points = last_cr.Points
		orginal_points = last_cr.Points
	}

	//开启事务
	o.Begin()
	sql = fmt.Sprintf("insert %s(no,uid,changes,points,pre_no,pre_points,des,act,ref,ref_id,ts,state) values(?,?,?,?,?,?,?,?,?,?,?,?)", record_tbl)
	ts := utils.TimeMillisecond(time.Now())
	no := credit_record_no(lock_cr.Uid)
	now_points := orginal_points
	if lock_cr.Act == CREDIT_OPERATION_LOCKDECR {
		now_points = orginal_points + int64(lock_cr.Changes)
	} else if lock_cr.Act == CREDIT_OPERATION_LOCKINCR {
		now_points = orginal_points
	}
	_, err = o.Raw(sql,
		no,
		lock_cr.Uid,
		lock_cr.Changes,
		now_points,
		pre_no,
		pre_points,
		fmt.Sprintf("回滚锁单:%s", op.No),
		op.Operation,
		"rblock",
		op.No,
		ts,
		CREDIT_RECORD_STATE_COMPLETED,
	).Exec()
	if err != nil {
		o.Rollback()
		logs.Errorf("credit service rblock insert %s table fail:%s", record_tbl, err.Error())
		return &Result{
			RemainingPoints: orginal_points,
			State:           OPERATION_CREDIT_STATE_ERROR,
			Error:           err,
		}
	}
	if lock_cr.Act == CREDIT_OPERATION_LOCKINCR {
		sql = fmt.Sprintf("update user_credits set rem_points=%d,lockincr_points=lockincr_points-%d,last_no=?,last_ts=%d where uid=?", now_points, lock_cr.Changes, ts)
	} else if lock_cr.Act == CREDIT_OPERATION_LOCKDECR {
		sql = fmt.Sprintf("update user_credits set rem_points=%d,lockdecr_points=lockdecr_points-%d,last_no=?,last_ts=%d where uid=?", now_points, lock_cr.Changes, ts)
	}
	o.Raw(sql, no, lock_cr.Uid).Exec()
	_, err = o.Raw(fmt.Sprintf("update %s set state=? where no=?", record_tbl), CREDIT_RECORD_STATE_RBLOCKED, op.No).Exec()
	if err != nil {
		o.Rollback()
		logs.Errorf("credit service rblock update %s table fail:%s", record_tbl, err.Error())
		return &Result{
			RemainingPoints: orginal_points,
			State:           OPERATION_CREDIT_STATE_ERROR,
			Error:           err,
		}
	}
	err = o.Commit()
	if err != nil {
		o.Rollback()
		logs.Errorf("credit service rblock commit transaction fail:%s", err.Error())
		return &Result{
			RemainingPoints: orginal_points,
			State:           OPERATION_CREDIT_STATE_ERROR,
			Error:           err,
		}
	}
	//更新高速缓存
	c.updateCredit(lock_cr.Uid, now_points)
	return &Result{
		No:              no,
		RemainingPoints: now_points,
		State:           OPERATION_CREDIT_STATE_SUCCESS,
		Error:           nil,
	}
	return nil
}

func (c *CreditService) unlockCredit(op *OperationCreditParameter) *Result {
	if len(op.No) == 0 {
		return &Result{
			RemainingPoints: 0,
			State:           OPERATION_CREDIT_STATE_PARAMFAIL,
			Error:           fmt.Errorf("参数错误"),
		}
	}

	//订单号校验
	record_tbl, uid, err := credit_record_tbl_byno(op.No)
	if err != nil {
		return &Result{
			RemainingPoints: 0,
			State:           OPERATION_CREDIT_STATE_NOFAIL,
			Error:           fmt.Errorf("非法单号"),
		}
	}

	//分布式锁
	lock := dlock.NewDistributedLock()
	locker, err := lock.Lock(c.lockKey(uid))
	if err != nil {
		logs.Errorf("credit service unlock get lock fail:%s", err.Error())
		return &Result{
			RemainingPoints: 0,
			State:           OPERATION_CREDIT_STATE_SYSBUSY,
			Error:           fmt.Errorf("系统繁忙,请稍后再试"),
		}
	}
	defer locker.Unlock()

	//校验订单状态
	lock_cr := c.GetCreditRecord(op.No)
	crlt := c.checkUnlockCreditState(lock_cr)
	if crlt != nil {
		return crlt
	}

	var last_cr CreditRecord
	sql := fmt.Sprintf("select no,uid,changes,points,pre_no,pre_points,des,act,ref,ref_id,ts from %s where uid=%d order by ts desc limit 1", record_tbl, lock_cr.Uid)
	o := dbs.NewOrm(credit_db)
	err = o.Raw(sql).QueryRow(&last_cr)

	var pre_no string = ""
	var pre_points int64 = 0
	var orginal_points int64 = 0
	if err == nil {
		pre_no = last_cr.No
		pre_points = last_cr.Points
		orginal_points = last_cr.Points
	}

	//开启事务
	o.Begin()
	sql = fmt.Sprintf("insert %s(no,uid,changes,points,pre_no,pre_points,des,act,ref,ref_id,ts,state) values(?,?,?,?,?,?,?,?,?,?,?,?)", record_tbl)
	ts := utils.TimeMillisecond(time.Now())
	no := credit_record_no(lock_cr.Uid)
	now_points := orginal_points
	if lock_cr.Act == CREDIT_OPERATION_LOCKDECR {
		now_points = orginal_points //锁单时已减去
	} else if lock_cr.Act == CREDIT_OPERATION_LOCKINCR {
		now_points = orginal_points + lock_cr.Changes
	}
	_, err = o.Raw(sql,
		no,
		lock_cr.Uid,
		lock_cr.Changes,
		now_points,
		pre_no,
		pre_points,
		fmt.Sprintf("完成锁单:%s", op.No),
		op.Operation,
		"unlock",
		op.No,
		ts,
		CREDIT_RECORD_STATE_COMPLETED,
	).Exec()
	if err != nil {
		o.Rollback()
		logs.Errorf("credit service unlock insert %s table fail:%s", record_tbl, err.Error())
		return &Result{
			RemainingPoints: orginal_points,
			State:           OPERATION_CREDIT_STATE_ERROR,
			Error:           err,
		}
	}
	if lock_cr.Act == CREDIT_OPERATION_LOCKINCR {
		sql = fmt.Sprintf("update user_credits set total_points=total_points+%d,rem_points=%d,lockincr_points=lockincr_points-%d,last_no=?,last_ts=%d where uid=?", lock_cr.Changes, now_points, lock_cr.Changes, ts)
	} else if lock_cr.Act == CREDIT_OPERATION_LOCKDECR {
		sql = fmt.Sprintf("update user_credits set rem_points=%d,lockdecr_points=lockdecr_points-%d,last_no=?,last_ts=%d where uid=?", now_points, lock_cr.Changes, ts)
	}
	o.Raw(sql, no, lock_cr.Uid).Exec()
	_, err = o.Raw(fmt.Sprintf("update %s set state=? where no=?", record_tbl), CREDIT_RECORD_STATE_COMPLETED, op.No).Exec()
	if err != nil {
		o.Rollback()
		logs.Errorf("credit service unlock update %s table fail:%s", record_tbl, err.Error())
		return &Result{
			RemainingPoints: orginal_points,
			State:           OPERATION_CREDIT_STATE_ERROR,
			Error:           err,
		}
	}
	err = o.Commit()
	if err != nil {
		o.Rollback()
		logs.Errorf("credit service unlock commit transaction fail:%s", err.Error())
		return &Result{
			RemainingPoints: orginal_points,
			State:           OPERATION_CREDIT_STATE_ERROR,
			Error:           err,
		}
	}
	//更新高速缓存
	c.updateCredit(lock_cr.Uid, now_points)
	return &Result{
		No:              no,
		RemainingPoints: now_points,
		State:           OPERATION_CREDIT_STATE_SUCCESS,
		Error:           nil,
	}
	return nil
}

func (c *CreditService) updateCredit(uid int64, points int64) {
	err := ssdb.New(use_ssdb_credit_db).Set(c.creditKey(uid), points)
	if err != nil {
		logs.Errorf("credit service update credit fail:%s", err.Error())
	}
}
