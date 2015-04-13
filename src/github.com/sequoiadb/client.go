package sequoiadb

/*
#cgo CFLAGS:  -I /opt/sequoiadb/include
#cgo LDFLAGS: -L /opt/sequoiadb/lib -lsdbc

#include "client.h"
*/
import "C"

import (
	"unsafe"
)

type sdbConnectionHandle C.sdbConnectionHandle
type sdbCSHandle C.sdbCSHandle
type sdbCollectionHandle C.sdbCollectionHandle
type sdbCursorHandle C.sdbCursorHandle

//=================================================

func sdbConnect(Host string, Name string, Usr string, Passwd string) (sdbConnectionHandle, int) {
	var conn C.sdbConnectionHandle = 0
	rc := C.sdbConnect((*C.CHAR)(C.CString(Host)), (*C.CHAR)(C.CString(Name)), (*C.CHAR)(C.CString(Usr)), (*C.CHAR)(C.CString(Passwd)), &conn)
	return sdbConnectionHandle(conn), int(rc)
}

func sdbDisconnect(conn sdbConnectionHandle) {
	C.sdbDisconnect((C.sdbConnectionHandle)(conn))
}

func sdbReleaseConnection(conn sdbConnectionHandle) {
	C.sdbReleaseConnection((C.sdbConnectionHandle)(conn))
}

//=================================================

func sdbGetCollectionSpace(conn sdbConnectionHandle, name string) (sdbCSHandle, int) {
	var hd C.sdbCSHandle = 0
	rc := C.sdbGetCollectionSpace((C.sdbConnectionHandle)(conn), (*C.CHAR)(C.CString(name)), &hd)
	return sdbCSHandle(hd), int(rc)
}

func sdbReleaseCS(cHandle sdbCSHandle) {
	C.sdbReleaseCS((C.sdbCSHandle(cHandle)))
}

//=================================================

func sdbGetCollection(cHandle sdbConnectionHandle, CollectionFullName string) (sdbCollectionHandle, int) {
	var hd C.sdbCollectionHandle = 0
	rc := C.sdbGetCollection((C.sdbConnectionHandle)(cHandle), (*C.CHAR)(C.CString(CollectionFullName)), &hd)
	return sdbCollectionHandle(hd), int(rc)
}

func sdbGetCollection1(cHandle sdbCSHandle, CollectionName string) (sdbCollectionHandle, int) {
	var handle C.sdbCollectionHandle = 0
	rc := C.sdbGetCollection1((C.sdbCSHandle(cHandle)), (*C.CHAR)(C.CString(CollectionName)), &handle)
	return sdbCollectionHandle(handle), int(rc)
}

func sdbReleaseCollection(cHandle sdbCollectionHandle) {
	C.sdbReleaseCollection((C.sdbCollectionHandle)(cHandle))
}

func sdbGetIndexes(cHandle sdbCollectionHandle, pIndexName string) (sdbCursorHandle, int) {
	var handld C.sdbCursorHandle = 0
	rc := C.sdbGetIndexes((C.sdbCollectionHandle)(cHandle), (*C.CHAR)(C.CString(pIndexName)), &handld)
	return sdbCursorHandle(handld), int(rc)
}

func sdbGetCount(cHandle sdbCollectionHandle, condition *bson, hint *bson) (int64, int) {
	var count C.SINT64 = 0
	rc := C.sdbGetCount1((C.sdbCollectionHandle)(cHandle), (*C.bson)(unsafe.Pointer(condition)), (*C.bson)(unsafe.Pointer(hint)), &count)
	return int64(count), int(rc)
}

func sdbInsert(cHandle sdbCollectionHandle, obj *bson) int {
	rc := C.sdbInsert((C.sdbCollectionHandle)(cHandle), (*C.bson)(unsafe.Pointer(obj)))
	return int(rc)
}

func sdbUpdate(cHandle sdbCollectionHandle, rule *bson, condition *bson, hint *bson) int {
	rc := C.sdbUpdate((C.sdbCollectionHandle)(cHandle), (*C.bson)(unsafe.Pointer(rule)), (*C.bson)(unsafe.Pointer(condition)), (*C.bson)(unsafe.Pointer(hint)))
	return int(rc)
}

func sdbUpsert(cHandle sdbCollectionHandle, rule *bson, condition *bson, hint *bson) int {
	rc := C.sdbUpsert((C.sdbCollectionHandle)(cHandle), (*C.bson)(unsafe.Pointer(rule)), (*C.bson)(unsafe.Pointer(condition)), (*C.bson)(unsafe.Pointer(hint)))
	return int(rc)
}

func sdbDelete(cHandle sdbCollectionHandle, condition *bson, hint *bson) int {
	rc := C.sdbDelete((C.sdbCollectionHandle)(cHandle), (*C.bson)(unsafe.Pointer(condition)), (*C.bson)(unsafe.Pointer(hint)))
	return int(rc)
}

func sdbQuery(cHandle sdbCollectionHandle, condition *bson, fields *bson, orderBy *bson, hint *bson, skip int32, limit int32, flag int32) (sdbCursorHandle, int) {
	var handle C.sdbCursorHandle = 0
	cHandle2 := (C.sdbCollectionHandle)(cHandle)
	condition2 := (*C.bson)(unsafe.Pointer(condition))
	fields2 := (*C.bson)(unsafe.Pointer(fields))
	orderBy2 := (*C.bson)(unsafe.Pointer(orderBy))
	hint2 := (*C.bson)(unsafe.Pointer(hint))
	rc := C.sdbQuery1(cHandle2, condition2, fields2, orderBy2, hint2, (C.INT64)(skip), (C.INT64)(limit), (C.INT32)(flag), &handle)
	return sdbCursorHandle(handle), int(rc)
}

func sdbExplain(cHandle sdbCollectionHandle, condition *bson, fields *bson, orderBy *bson, hint *bson, skip int32, limit int32, flag int32, options *bson) (sdbCursorHandle, int) {
	var handle C.sdbCursorHandle = 0
	cHandle2 := (C.sdbCollectionHandle)(cHandle)
	condition2 := (*C.bson)(unsafe.Pointer(condition))
	fields2 := (*C.bson)(unsafe.Pointer(fields))
	orderBy2 := (*C.bson)(unsafe.Pointer(orderBy))
	hint2 := (*C.bson)(unsafe.Pointer(hint))
	options2 := (*C.bson)(unsafe.Pointer(options))
	rc := C.sdbExplain(cHandle2, condition2, fields2, orderBy2, hint2, (C.INT32)(flag), (C.INT64)(skip), (C.INT64)(limit), options2, &handle)
	return sdbCursorHandle(handle), int(rc)
}

//=================================================

func sdbNext(cHandle sdbCursorHandle, obj *bson) int {
	rc := C.sdbNext((C.sdbCursorHandle)(cHandle), (*C.bson)(unsafe.Pointer(obj)))
	return int(rc)
}

func sdbCurrent(cHandle sdbCursorHandle, obj *bson) int {
	rc := C.sdbCurrent((C.sdbCursorHandle)(cHandle), (*C.bson)(unsafe.Pointer(obj)))
	return int(rc)
}

func sdbCloseCursor(cHandle sdbCursorHandle) int {
	rc := C.sdbCloseCursor((C.sdbCursorHandle)(cHandle))
	return int(rc)
}

func sdbReleaseCursor(cHandle sdbCursorHandle) {
	C.sdbReleaseCursor((C.sdbCursorHandle)(cHandle))
}

func sdbCloseAllCursors(cHandle sdbConnectionHandle) int {
	rc := C.sdbCloseAllCursors((C.sdbConnectionHandle)(cHandle))
	return int(rc)
}

func sdbExec(cHandle sdbConnectionHandle, sql string) (sdbCursorHandle, int) {
	var result C.sdbCursorHandle = 0
	rc := C.sdbExec((C.sdbConnectionHandle)(cHandle), (*C.CHAR)(C.CString(sql)), &result)
	return sdbCursorHandle(result), int(rc)
}

func sdbExecUpdate(cHandle sdbConnectionHandle, sql string) int {
	rc := C.sdbExecUpdate((C.sdbConnectionHandle)(cHandle), (*C.CHAR)(C.CString(sql)))
	return int(rc)
}

//=================================================

func sdbTransactionBegin(cHandle sdbConnectionHandle) int {
	rc := C.sdbTransactionBegin((C.sdbConnectionHandle)(cHandle))
	return int(rc)
}

func sdbTransactionCommit(cHandle sdbConnectionHandle) int {
	rc := C.sdbTransactionCommit((C.sdbConnectionHandle)(cHandle))
	return int(rc)
}

func sdbTransactionRollback(cHandle sdbConnectionHandle) int {
	rc := C.sdbTransactionRollback((C.sdbConnectionHandle)(cHandle))
	return int(rc)
}

//=================================================
