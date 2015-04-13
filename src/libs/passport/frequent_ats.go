package passport

import (
	"fmt"
	"reflect"
	"utils/ssdb"
)

func at_count_cachekey(fromUid int64) string {
	return fmt.Sprintf("member_frequent_at_list_%d", fromUid)
}

//累加@次数
func IncrAtCount(fromUid int64, toUid int64) {
	if fromUid <= 0 || toUid <= 0 {
		return
	}
	ck := at_count_cachekey(fromUid)
	ssdb.New(use_ssdb_passport_db).Zincrby(ck, toUid, 1)
}

//返回最多@的人员map[uid]昵称,从大到小
func GetFrequentAts(fromUid int64, size int) []int64 {
	ck := at_count_cachekey(fromUid)
	result, _ := ssdb.New(use_ssdb_passport_db).Zrevrange(ck, 0, size, reflect.TypeOf(int64(0)))
	list := []int64{}
	for _, obj := range result {
		toUid := *(obj.(*int64))
		list = append(list, toUid)
	}
	return list
}
