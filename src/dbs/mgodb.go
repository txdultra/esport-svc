package dbs

import (
	"labix.org/v2/mgo"
)

func MgoSession() *mgo.Session {
	// 连接数据库,使用后关闭Session
	session, err := mgo.Dial(mongodb_addrs)
	if err != nil {
		panic(err)
	}
	if mongodb_session_consistency == "Eventual" {
		session.SetMode(mgo.Eventual, mongodb_session_refresh)
	} else if mongodb_session_consistency == "Strong" {
		session.SetMode(mgo.Strong, mongodb_session_refresh)
	} else {
		session.SetMode(mgo.Monotonic, mongodb_session_refresh)
	}
	return session
}

func MgoC(db string, c string) (*mgo.Session, *mgo.Collection) {
	session := MgoSession()
	return session, session.DB(db).C(c)
}
