package feedback

import (
	"time"
)

type Feedback struct {
	Id       int
	Uid      int64
	Category string
	PostTime time.Time
	Title    string
	Content  string
	Img      int64
	Contact  string
	Source   string
}

func (self *Feedback) TableName() string {
	return "common_feedback"
}

func (self *Feedback) TableEngine() string {
	return "INNODB"
}
