package comment

import (
	"dbs"
	"fmt"
	"libs/hook"
	"libs/stat"
	"libs/vars"
	"logs"

	"labix.org/v2/mgo/bson"
	//"regexp"
	//"strconv"
	//"strings"
	"libs/message"
	"libs/passport"
	"time"
	"utils"
	"utils/ssdb"
)

const (
	comment_at_screen_fmt    = "@%s"
	comment_attext_fmt       = "@%s:"
	comment_attext_relay_fmt = " //%s"

	//comment_inner_at_fmt   = "[@:%d]"
	//comment_inner_at_regex = `\[@:(\d+)\]`
)

type Commentor struct {
	mod string
}

func NewCommentor(mod string) *Commentor {
	c := &Commentor{}
	c.mod = mod
	return c
}

func (c *Commentor) Create(data map[string]interface{}, atId func(string) int64, atName func(int64) string, msgNotice bool, msgType vars.MSG_TYPE) (error, string) {
	comment := c.dataToComment(data)
	if comment.RefId <= 0 {
		return fmt.Errorf("RefId 属性不能为0"), ""
	}
	if comment.FromId <= 0 {
		return fmt.Errorf("FromId 属性不能为0"), ""
	}
	if len(comment.Text) == 0 {
		return fmt.Errorf("Text 属性不能为空"), ""
	}
	fromName := atName(comment.FromId)
	if len(fromName) == 0 {
		return fmt.Errorf("FromId 必须有名称"), ""
	}

	var replyComment *Comment
	if len(comment.Pid) > 0 && bson.IsObjectIdHex(comment.Pid) {
		replyComment = c.Get(comment.Pid)
	}
	comment.WbText = fmt.Sprintf(comment_attext_fmt, fromName) + comment.Text
	atIds := make(map[string]int64)
	atNames := utils.ExtractAts(comment.Text)
	for _, _at := range atNames {
		_id := atId(_at)
		if _id > 0 {
			_atname := fmt.Sprintf(comment_at_screen_fmt, _at)
			if _, ok := atIds[_atname]; ok {
				continue
			}
			atIds[_atname] = _id
		}
	}
	//需要发送消息的uids
	msgAtUids := []int64{}
	msg_existed := false
	for _, _uid := range atIds {
		if replyComment != nil && replyComment.FromId == _uid {
			msg_existed = true
		}
		msgAtUids = append(msgAtUids, _uid)
	}
	//回复的对象也发送一条消息
	if replyComment != nil && !msg_existed {
		msgAtUids = append(msgAtUids, replyComment.FromId)
	}

	//@计数器
	go func() {
		for _, atuid := range atIds {
			passport.IncrAtCount(comment.FromId, atuid)
		}
	}()
	fromAtName := fmt.Sprintf(comment_at_screen_fmt, fromName)
	comment.Ats = atIds
	comment.Ats[fromAtName] = comment.FromId //加入发送者
	comment.ID = bson.NewObjectId()
	comment.PostTime = time.Now()
	comment.Nano = time.Now().UnixNano()
	if replyComment != nil {
		comment.Reply = &ReplyComment{
			ID:        replyComment.ID.Hex(),
			FromId:    replyComment.FromId,
			Position:  replyComment.Position,
			Text:      replyComment.Text,
			PostTime:  replyComment.PostTime,
			IP:        replyComment.IP,
			Longitude: replyComment.Longitude,
			Latitude:  replyComment.Latitude,
			Ats:       replyComment.Ats,
			WbText:    replyComment.WbText,
		}
		comment.WbText += fmt.Sprintf(comment_attext_relay_fmt, replyComment.WbText)
		for k, v := range comment.Reply.Ats {
			if _, ok := comment.Ats[k]; ok {
				continue
			}
			comment.Ats[k] = v
		}
	}

	collection_name, _ := commonts_collections[c.mod]
	session, col := dbs.MgoC(comments_db_name, collection_name)
	defer session.Close()
	query := col.Find(bson.M{"ref_id": comment.RefId})
	query.Sort("-position")
	iter := query.Iter()
	r := Comment{}
	ok := iter.Next(&r)
	if ok {
		comment.Position = r.Position + 1
	} else {
		comment.Position = 1
	}

	//过滤敏感字
	comment.Text = utils.CensorWords(comment.Text)

	err := col.Insert(comment)
	if err != nil {
		logs.Error(err.Error())
		return fmt.Errorf("插入错误:%s", err), ""
	}
	if msgNotice {
		//处理@通知
		go func() {
			for _, atuid := range msgAtUids {
				if atuid != comment.FromId {
					message.SendMsgV2(comment.FromId, atuid, msgType, comment.Text, comment.ID.Hex(), nil)
				}
			}
		}()
	}
	stat.GetCounter(c.mod).DoC(comment.RefId, 1, "comments")
	go func() {
		com_bs, err := bson.Marshal(*comment)
		if err == nil {
			ssdb.New(use_ssdb_comment_db).Set(c.ckey(comment.ID.Hex()), com_bs)
		}
	}()

	//钩子事件
	go func() {
		if c.mod == "vod" {
			hook.Do("vod_comment", comment.FromId, 1)
		}
	}()
	return nil, comment.ID.Hex()
}

func (c *Commentor) dataToComment(data map[string]interface{}) *Comment {
	comment := new(Comment)
	if val, ok := data["reply_id"]; ok {
		comment.Pid = val.(string)
	}
	if val, ok := data["ref_id"]; ok {
		comment.RefId = val.(int64)
	}
	if val, ok := data["from_id"]; ok {
		comment.FromId = val.(int64)
	}
	if val, ok := data["title"]; ok {
		comment.Title = val.(string)
	}
	if val, ok := data["text"]; ok {
		comment.Text = val.(string)
	}
	if val, ok := data["ip"]; ok {
		comment.IP = val.(string)
	}
	if val, ok := data["longitude"]; ok {
		comment.Longitude = val.(float64)
	}
	if val, ok := data["latitude"]; ok {
		comment.Latitude = val.(float64)
	}
	if val, ok := data["allow_reply"]; ok {
		comment.AllowReply = val.(bool)
	}
	return comment
}

func (c *Commentor) ckey(id string) string {
	return fmt.Sprintf("mobile_comment:%s", id)
}

func (c *Commentor) Get(id string) *Comment {
	if !bson.IsObjectIdHex(id) {
		return nil
	}
	var cmt *Comment
	var bs []byte
	err := ssdb.New(use_ssdb_comment_db).Get(c.ckey(id), &bs)
	if err == nil {
		err = bson.Unmarshal(bs, &cmt)
		if err == nil {
			return cmt
		}
	}

	collection_name, _ := commonts_collections[c.mod]
	session, col := dbs.MgoC(comments_db_name, collection_name)
	defer session.Close()
	err = col.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&cmt)
	if err == nil {
		com_bs, err := bson.Marshal(*cmt)
		if err == nil {
			ssdb.New(use_ssdb_comment_db).Set(c.ckey(id), com_bs)
		}
		return cmt
	}
	return nil
}

func (c *Commentor) Gets(mod string, refId int64, p int, s int, t time.Time, nano int64) (int, []*Comment, int64) {
	collection_name, _ := commonts_collections[c.mod]
	session, col := dbs.MgoC(comments_db_name, collection_name)
	defer session.Close()
	offset := (p - 1) * s
	var cms []*Comment
	fm := bson.M{}
	if !t.IsZero() {
		fm["post_time"] = bson.M{"$lt": t}
	}
	fm["ref_id"] = refId
	qs := col.Find(fm)
	counts, _ := qs.Count()
	err := qs.Sort("-post_time").Skip(offset).Limit(s).All(&cms)
	if err != nil {
		logs.Error(err.Error())
	}
	_last_nano := int64(0)
	if len(cms) > 0 {
		_last_idx := len(cms) - 1
		_last_nano = cms[_last_idx].Nano
	}
	return counts, cms, _last_nano
}
