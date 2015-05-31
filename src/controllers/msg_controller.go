package controllers

import (
	//"fmt"
	"libs"
	"libs/comment"
	"libs/groups"
	"libs/message"
	"libs/share"
	"libs/vars"
	"libs/vod"
	"outobjs"
	"strconv"
	"time"
	"utils"
)

// 消息模块 API
type MessageController struct {
	AuthorizeController
}

func (c *MessageController) Prepare() {
	c.AuthorizeController.Prepare()
}

func (c *MessageController) URLMapping() {
	c.Mapping("Mentions", c.Mentions)
	c.Mapping("EmptyMentions", c.EmptyMentions)
	c.Mapping("DelMention", c.DelMention)
	c.Mapping("Count", c.Count)
}

// @Title 获取@当前用户的最新消息
// @Description 获取@当前用户的最新消息(自动清空计数;消息类型: comment,vod,text)
// @Param   access_token   path   string  true  "access_token"
// @Param   page   path   int  false  "页(默认1)"
// @Param   size   path   int  false  "页数量(默认20)"
// @Param   msg_type   path   string  true  "消息类型(vod,share,group,sys)"
// @Param   t    path  int  false  "时间戳(每次请求获得的t属性)"
// @Success 200 {object} outobjs.OutAtMsgPageList
// @router /mentions [get]
func (c *MessageController) Mentions() {
	uid := c.CurrentUid()
	size, _ := c.GetInt("size")
	page, _ := c.GetInt("page")
	timestamp, _ := c.GetInt64("t")
	msg_type := c.GetString("msg_type")

	if len(msg_type) == 0 {
		c.Json(libs.NewError("member_msgbox_gets_fail", "MSG1011", "未指定需要获取的消息类型", ""))
		return
	}
	mt := c.transformMsgType(msg_type)
	if mt == "" {
		c.Json(libs.NewError("member_msgbox_gets_fail", "MSG1012", "获取消息的类型不存在", ""))
		return
	}
	t := time.Now()
	if timestamp > 0 {
		t = time.Unix(timestamp, 0)
	}
	if size <= 0 {
		size = 20
	}
	if page <= 0 {
		page = 1
	}
	total, msgs := message.GetMsgs(uid, int(page), int(size), t, mt)
	out_msgs := []*outobjs.OutAtMsg{}
	ts := t
	for _, m := range msgs {
		out_m := &outobjs.OutAtMsg{
			FromUid:    m.FromUid,
			ToUid:      m.ToUid,
			FromMember: outobjs.GetOutSimpleMember(m.FromUid),
			MsgType:    m.MsgType,
			Text:       m.Text,
			RefId:      m.RefId,
			PostTime:   m.PostTime,
			FriendTime: utils.FriendTime(m.PostTime),
		}
		c.transformObj(m.RefId, m.MsgType, out_m)
		out_msgs = append(out_msgs, out_m)
		if m.PostTime.Before(ts) {
			ts = m.PostTime
		}
	}
	out_list := outobjs.OutAtMsgPageList{
		Total:       total,
		TotalPage:   utils.TotalPages(int(total), int(size)),
		CurrentPage: int(page),
		Size:        int(size),
		Lists:       out_msgs,
		Time:        ts.Unix(),
	}
	//自动清空计数
	message.ResetEventCount(uid, mt)
	c.Json(out_list)
}

func (c *MessageController) transformMsgType(submitName string) vars.MSG_TYPE {
	switch submitName {
	case "vod":
		return vod.MSG_TYPE_COMMENT
	case "share":
		return share.MSG_TYPE_TEXT
	case "group":
		return groups.MSG_TYPE_MESSAGE
	default:
		return vars.MSG_TYPE_SYS
	}
}

func (c *MessageController) transformObj(refId string, msgType vars.MSG_TYPE, out *outobjs.OutAtMsg) {
	switch msgType {
	case vod.MSG_TYPE_COMMENT:
		cmtt := comment.NewCommentor("vod")
		cmt := cmtt.Get(refId)
		if cmt != nil {
			out.Comment = outobjs.GetOutComment(cmt)
		}
		break
	}
}

// @Title 清空@消息
// @Description 清空某人的@消息(成功返回error_code:REP000)
// @Param   access_token   path   string  true  "access_token"
// @Param   msg_type   path   string  true  "消息类型(vod,share,group)"
// @Success 200 成功返回error_code:REP000
// @router /mentions/empty [post]
func (c *MessageController) EmptyMentions() {
	uid := c.CurrentUid()
	msg_type := c.GetString("msg_type")
	if len(msg_type) == 0 {
		c.Json(libs.NewError("member_msgbox_empty_fail", "MSG1002", "未指定需要清空的消息类型", ""))
		return
	}
	mt := c.transformMsgType(msg_type)
	if mt == "" {
		c.Json(libs.NewError("member_msgbox_empty_fail", "MSG1003", "清空消息的类型不存在", ""))
		return
	}
	err := message.EmptyMsgBox(uid, mt)
	if err == nil {
		c.Json(libs.NewError("member_msgbox_empty_succ", RESPONSE_SUCCESS, "成功清空", ""))
		return
	}
	c.Json(libs.NewError("member_msgbox_empty_fail", "MSG1001", err.Error(), ""))
}

// @Title 清空@消息
// @Description 清空某人的@消息(成功返回error_code:REP000)
// @Param   access_token   path   string  true  "access_token"
// @Param   msg_id   path   string  true  "消息id"
// @Param   msg_type   path   string  true  "消息类型(vod,share,group)"
// @Success 200 成功返回error_code:REP000
// @router /mentions/del [post]
func (c *MessageController) DelMention() {
	msg_id := c.GetString("msg_id")
	if len(msg_id) == 0 {
		c.Json(libs.NewError("member_msgbox_del_fail", "MSG1102", "必须提供msg_id参数", ""))
		return
	}
	msg_type := c.GetString("msg_type")
	if len(msg_type) == 0 {
		c.Json(libs.NewError("member_msgbox_del_fail", "MSG1103", "未指定需要删除的消息类型", ""))
		return
	}
	mt := c.transformMsgType(msg_type)
	if mt == "" {
		c.Json(libs.NewError("member_msgbox_del_fail", "MSG1104", "删除消息的类型不存在", ""))
		return
	}

	uid := c.CurrentUid()
	err := message.DelMsg(uid, msg_id, mt)
	if err == nil {
		c.Json(libs.NewError("member_msgbox_del_succ", RESPONSE_SUCCESS, "成功清空", ""))
		return
	}
	c.Json(libs.NewError("member_msgbox_del_fail", "MSG1101", err.Error(), ""))
}

// @Title 获取@当前用户的新消息数量
// @Description 获取@当前用户的新消息数量(成功返回error_code:REP000,error_description:数量[int类型])
// @Param   access_token   path   string  true  "access_token"
// @Param   msg_type   path   string  true  "消息类型(vod,share,group)"
// @Success 200 成功返回error_code:REP000,error_description:数量[int类型]
// @router /count [get]
func (c *MessageController) Count() {
	uid := c.CurrentUid()
	msg_type := c.GetString("msg_type")
	if len(msg_type) == 0 {
		c.Json(libs.NewError("member_msgbox_count_fail", "MSG1110", "未指定需要清空的消息类型", ""))
		return
	}
	mt := c.transformMsgType(msg_type)
	if mt == "" {
		c.Json(libs.NewError("member_msgbox_count_fail", "MSG1111", "清空消息的类型不存在", ""))
		return
	}
	count := message.NewEventCount(uid, mt)
	c.Json(libs.NewError("member_msgbox_count", RESPONSE_SUCCESS, strconv.Itoa(count), ""))
}
