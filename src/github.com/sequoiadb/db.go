package sequoiadb

import (
	"encoding/json"
	"log"
	"strings"
)

type Bson struct {
	Data *bson
}

type Session struct {
	handle sdbConnectionHandle
}

type Cursor struct {
	cursor sdbCursorHandle
}

func jsonToMap(str string) interface{} {
	if strings.HasPrefix(str, "[") == true {
		data := []interface{}{}
		err := json.Unmarshal([]byte(str), &data)
		if err != nil {
			log.Println("error:", err)
		}
		return data
	}
	data := map[string]interface{}{}
	err := json.Unmarshal([]byte(str), &data)
	if err != nil {
		log.Println("error:", err)
	}
	return data
}

func Connect(Host string, Name string, Usr string, Passwd string) (*Session, int) {
	hd, rc := sdbConnect(Host, Name, Usr, Passwd)
	if rc != 0 {
		return nil, rc
	}
	return &Session{handle: hd}, rc
}

func NewBson() *Bson {
	return &Bson{Data: bson_create()}
}

func (b *Bson) Destroy() {
	bson_destroy(b.Data)
}

func (b *Bson) Print() {
	bson_print(b.Data)
}

func (b *Bson) Sprint() string {
	str, _ := bson_sprint(bson_sprint_length(b.Data), b.Data)
	return str
}

func (b *Bson) Marshal() interface{} {
	str, _ := bson_sprint(bson_sprint_length(b.Data), b.Data)
	return jsonToMap(str)
}

func (c *Cursor) Next() (interface{}, int) {
	b := NewBson()
	defer b.Destroy()
	rc := sdbNext(c.cursor, b.Data)
	if rc != 0 {
		return nil, rc
	}
	return b.Marshal(), 0
}

func (c *Cursor) Current() (interface{}, int) {
	b := NewBson()
	defer b.Destroy()
	rc := sdbCurrent(c.cursor, b.Data)
	if rc != 0 {
		return nil, rc
	}
	return b.Marshal(), 0
}

func (c *Cursor) CloseCursor() int {
	rc := sdbCloseCursor(c.cursor)
	return rc
}

func (c *Cursor) ReleaseCursor() {
	sdbReleaseCursor(c.cursor)
}

func (c *Cursor) CloseAllCursors(s *Session) int {
	rc := sdbCloseAllCursors(s.handle)
	return rc
}

func (s *Session) Query(sql string) *Cursor {
	hd, _ := sdbExec(s.handle, sql)
	return &Cursor{cursor: hd}
}

func (s *Session) Exec(sql string) int {
	return sdbExecUpdate(s.handle, sql)
}

func (s *Session) TransactionBegin() int {
	return sdbTransactionBegin(s.handle)
}

func (s *Session) TransactionCommit() int {
	return sdbTransactionCommit(s.handle)
}

func (s *Session) TransactionRollback() int {
	return sdbTransactionRollback(s.handle)
}

func (s *Session) Disconnect() {
	sdbDisconnect(s.handle)
}

func (s *Session) ReleaseConnection() {
	sdbReleaseConnection(s.handle)
}
