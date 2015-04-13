package sequoiadb

/*
#cgo CFLAGS:  -I /opt/sequoiadb/include
#cgo LDFLAGS: -L /opt/sequoiadb/lib -lsdbc

#include "client.h"

int getStringArrayLen(char** list){
	return strlen(*list)-1;
}

char** newStringArrayPtr(){
	return malloc(sizeof(char*));
}

void freeStringArrayPtr(char** list){
	free(list);
}

char* newStringPtr(int n){
	return malloc(sizeof(char)*n);
}

void freeStringPtr(char* p){
	return free(p);
}

*/
import "C"

import (
	"reflect"
	"unsafe"
)

type bson C.bson
type bson_iterator C.bson_iterator
type bson_oid_t C.bson_oid_t

func bson_create() *bson {
	obj := C.bson_create()
	return (*bson)(unsafe.Pointer(obj))
}

func bson_dispose(b *bson) {
	C.bson_dispose((*C.bson)(unsafe.Pointer(b)))
}

func bson_size(b *bson) int {
	rc := C.bson_size((*C.bson)(unsafe.Pointer(b)))
	return int(rc)
}

func bson_buffer_size(b *bson) int {
	rc := C.bson_buffer_size((*C.bson)(unsafe.Pointer(b)))
	return int(rc)
}

func bson_print(b *bson) {
	C.bson_print((*C.bson)(unsafe.Pointer(b)))
}

func bson_sprint_length(b *bson) int {
	rc := C.bson_sprint_length((*C.bson)(unsafe.Pointer(b)))
	return int(rc)
}

func bson_sprint(size int, b *bson) (string, int) {
	var s C.int = (C.int)(size)
	var buf *C.char = C.newStringPtr(s)
	rc := C.bson_sprint(buf, s, (*C.bson)(unsafe.Pointer(b)))
	defer C.freeStringPtr(buf)
	return C.GoString(buf), int(rc)
}

func bson_sprint_iterator(i *bson_iterator, delChar byte) ([]string, int, int) {
	var buf **C.char = C.newStringArrayPtr()
	var num C.int
	rc := C.bson_sprint_iterator(buf, &num, (*C.bson_iterator)(unsafe.Pointer(i)), (C.char)(delChar))
	if int(rc) == 0 {
		var pbuf []*C.char
		header := (*reflect.SliceHeader)(unsafe.Pointer(&pbuf))
		header.Data = uintptr(unsafe.Pointer(buf))
		header.Cap = int(C.getStringArrayLen(buf))
		header.Len = int(C.getStringArrayLen(buf))
		buffer := []string{}
		for _, i := range pbuf {
			buffer = append(buffer, C.GoString(i))
		}
		return buffer, int(num), 0
	}
	return nil, 0, int(rc)
}

func bson_sprint_length_iterator(i *bson_iterator) int {
	rc := C.bson_sprint_length_iterator((*C.bson_iterator)(unsafe.Pointer(i)))
	return int(rc)
}

func bson_init(b *bson) {
	C.bson_init((*C.bson)(unsafe.Pointer(b)))
}

func bson_init_data(b *bson, data *string) int {
	rc := C.bson_init_data((*C.bson)(unsafe.Pointer(b)), (*C.char)(unsafe.Pointer(data)))
	return int(rc)
}

func bson_init_finished_data(b *bson, data string) int {
	rc := C.bson_init_finished_data((*C.bson)(unsafe.Pointer(b)), (*C.char)(C.CString(data)))
	return int(rc)
}

func bson_init_size(b *bson, size int) {
	C.bson_init_size((*C.bson)(unsafe.Pointer(b)), (C.int)(size))
}

func bson_append_start_object(b *bson, name string) int {
	rc := C.bson_append_start_object((*C.bson)(unsafe.Pointer(b)), (*C.char)(C.CString(name)))
	return int(rc)
}

func bson_append_finish_object(b *bson) int {
	rc := C.bson_append_finish_object((*C.bson)(unsafe.Pointer(b)))
	return int(rc)
}

func bson_append_start_array(b *bson, name string) int {
	rc := C.bson_append_start_array((*C.bson)(unsafe.Pointer(b)), (*C.char)(C.CString(name)))
	return int(rc)
}

func bson_append_finish_array(b *bson) int {
	rc := C.bson_append_finish_array((*C.bson)(unsafe.Pointer(b)))
	return int(rc)
}

func bson_append_string(b *bson, key string, value string) int {
	rc := C.bson_append_string((*C.bson)(unsafe.Pointer(b)), (*C.char)(C.CString(key)), (*C.char)(C.CString(value)))
	return int(rc)
}

func bson_append_int(b *bson, key string, value int) int {
	rc := C.bson_append_int((*C.bson)(unsafe.Pointer(b)), (*C.char)(C.CString(key)), (C.int)(value))
	return int(rc)
}

func bson_append_double(b *bson, key string, value float64) int {
	rc := C.bson_append_double((*C.bson)(unsafe.Pointer(b)), (*C.char)(C.CString(key)), (C.double)(value))
	return int(rc)
}

func bson_finish(b *bson) int {
	rc := C.bson_finish((*C.bson)(unsafe.Pointer(b)))
	return int(rc)
}

func bson_empty(b *bson) *bson {
	rc := C.bson_empty((*C.bson)(unsafe.Pointer(b)))
	return (*bson)(rc)
}

func bson_destroy(b *bson) {
	C.bson_destroy((*C.bson)(unsafe.Pointer(b)))
}

func bson_copy(out *bson, in *bson) int {
	rc := C.bson_copy((*C.bson)(unsafe.Pointer(out)), (*C.bson)(unsafe.Pointer(in)))
	return int(rc)
}

func bson_iterator_create() *bson_iterator {
	rc := C.bson_iterator_create()
	return (*bson_iterator)(unsafe.Pointer(rc))
}

func bson_iterator_init(i *bson_iterator, b *bson) {
	C.bson_iterator_init((*C.bson_iterator)(unsafe.Pointer(i)), (*C.bson)(unsafe.Pointer(b)))
}

/*
bson_iterator_type, bson_find, bson_iterator_next 函数的返回值如下:
-1	Min key
0	End of object
1	Floatint point
2	UTF-8 string
3	Embedded document
4	Array
5	Binary data
6	Undefined-Deprecated
7	Objectid
8	BOOlean
9	UTC datetime
10	Null value
11	Regular expression
12	DBPointer-Deprecated
13	JavaScript code
14	Symbol-Deprecated
15	JavaScript code w/scope
16	32-bit integer
17	Timestamp
18	64-bit integer
127	Max key
*/
func bson_find(it *bson_iterator, obj *bson, name string) int {
	rc := C.bson_find((*C.bson_iterator)(unsafe.Pointer(it)), (*C.bson)(unsafe.Pointer(obj)), (*C.char)(C.CString(name)))
	return int(rc)
}

func bson_iterator_next(i *bson_iterator) int {
	rc := C.bson_iterator_next((*C.bson_iterator)(unsafe.Pointer(i)))
	return int(rc)
}

func bson_iterator_type(i *bson_iterator) int {
	rc := C.bson_iterator_type((*C.bson_iterator)(unsafe.Pointer(i)))
	return int(rc)
}

func bson_iterator_key(i *bson_iterator) string {
	rc := C.bson_iterator_key((*C.bson_iterator)(unsafe.Pointer(i)))
	return C.GoString(rc)
}

func bson_iterator_value(i *bson_iterator) string {
	rc := C.bson_iterator_value((*C.bson_iterator)(unsafe.Pointer(i)))
	return C.GoString(rc)
}

func bson_iterator_double(i *bson_iterator) float64 {
	rc := C.bson_iterator_double((*C.bson_iterator)(unsafe.Pointer(i)))
	return float64(rc)
}

func bson_iterator_oid(i *bson_iterator) *bson_oid_t {
	rc := C.bson_iterator_oid((*C.bson_iterator)(unsafe.Pointer(i)))
	return (*bson_oid_t)(unsafe.Pointer(rc))
}

func bson_iterator_int(i *bson_iterator) int {
	rc := C.bson_iterator_int((*C.bson_iterator)(unsafe.Pointer(i)))
	return int(rc)
}

func bson_iterator_long(i *bson_iterator) uint64 {
	rc := C.bson_iterator_long((*C.bson_iterator)(unsafe.Pointer(i)))
	return uint64(rc)
}

func bson_iterator_bool(i *bson_iterator) int {
	rc := C.bson_iterator_bool((*C.bson_iterator)(unsafe.Pointer(i)))
	return int(rc)
}

func bson_iterator_string(i *bson_iterator) string {
	rc := C.bson_iterator_string((*C.bson_iterator)(unsafe.Pointer(i)))
	return C.GoString(rc)
}

func bson_iterator_string_len(i *bson_iterator) int {
	rc := C.bson_iterator_string_len((*C.bson_iterator)(unsafe.Pointer(i)))
	return int(rc)
}

func bson_iterator_code(i *bson_iterator) string {
	rc := C.bson_iterator_code((*C.bson_iterator)(unsafe.Pointer(i)))
	return C.GoString(rc)
}

func bson_iterator_code_scope(i *bson_iterator, scope *bson) {
	C.bson_iterator_code_scope((*C.bson_iterator)(unsafe.Pointer(i)), (*C.bson)(unsafe.Pointer(scope)))
}

func bson_iterator_date(i *bson_iterator) uint64 {
	rc := C.bson_iterator_date((*C.bson_iterator)(unsafe.Pointer(i)))
	return uint64(rc)
}

func bson_iterator_subiterator(i *bson_iterator, sub *bson_iterator) {
	C.bson_iterator_subiterator((*C.bson_iterator)(unsafe.Pointer(i)), (*C.bson_iterator)(unsafe.Pointer(sub)))
}

func bson_oid_from_string(str string) *bson_oid_t {
	var oid *C.bson_oid_t
	C.bson_oid_gen(oid)
	C.bson_oid_from_string(oid, (*C.char)(C.CString(str)))
	return (*bson_oid_t)(unsafe.Pointer(oid))
}

func bson_oid_to_string(oid *bson_oid_t) string {
	var buf *C.char = C.newStringPtr(30)
	C.bson_oid_to_string((*C.bson_oid_t)(unsafe.Pointer(oid)), buf)
	return C.GoString(buf)
}

func bson_oid_gen(oid *bson_oid_t) {
	C.bson_oid_gen((*C.bson_oid_t)(unsafe.Pointer(oid)))
}

func toBson(obj *bson, data interface{}) {
	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Map {
		keys := value.MapKeys()
		for i := 0; i < len(keys); i++ {
			switchMapType(obj, keys[i], value.MapIndex(keys[i]))
		}
	} else if value.Kind() == reflect.Slice {
		for i := 0; i < value.Len(); i++ {
			switchMapType(obj, reflect.ValueOf(i), value.Index(i))
		}
	}
}

func switchMapType(obj *bson, key reflect.Value, value reflect.Value) {
	n, _ := key.Interface().(string)
	switch reflect.ValueOf(value.Interface()).Kind() {
	case reflect.Int:
		m, _ := value.Interface().(int)
		bson_append_int(obj, n, m)
	case reflect.Float64:
		m, _ := value.Interface().(float64)
		bson_append_double(obj, n, m)
	case reflect.String:
		m, _ := value.Interface().(string)
		bson_append_string(obj, n, m)
	case reflect.Slice:
		bson_append_start_array(obj, n)
		toBson(obj, value.Interface())
		bson_append_finish_array(obj)
	case reflect.Map:
		bson_append_start_object(obj, n)
		toBson(obj, value.Interface())
		bson_append_finish_object(obj)
	}
}

func toMap(data interface{}, i *bson_iterator) {
	value := reflect.Indirect(reflect.ValueOf(data))
	if value.Kind() == reflect.Map {
		for {
			if bson_iterator_next(i) == 0 {
				break
			}
			switchBsonType(data, i, false)
		}
	} else if value.Kind() == reflect.Slice {
		for {
			if bson_iterator_next(i) == 0 {
				break
			}
			switchBsonType(data, i, true)
		}
	}
}

func switchBsonType(data interface{}, i *bson_iterator, state bool) {
	s := bson_iterator_create()
	r := reflect.Indirect(reflect.ValueOf(data))
	k := reflect.ValueOf(bson_iterator_key(i))
	var v reflect.Value
	switch bson_iterator_type(i) {
	case 1:
		v = reflect.ValueOf(bson_iterator_double(i))
	case 2:
		v = reflect.ValueOf(bson_iterator_string(i))
	case 3:
		bson_iterator_subiterator(i, s)
		tmp1 := map[string]interface{}{}
		toMap(&tmp1, s)
		v = reflect.ValueOf(tmp1)
	case 4:
		bson_iterator_subiterator(i, s)
		tmp2 := []interface{}{}
		toMap(&tmp2, s)
		v = reflect.ValueOf(tmp2)
	case 7:
		oid := bson_iterator_oid(i)
		v = reflect.ValueOf(bson_oid_to_string(oid))
	}
	if state == true {
		r.Set(reflect.Append(r, v))
	} else {
		r.SetMapIndex(k, v)
	}
}
