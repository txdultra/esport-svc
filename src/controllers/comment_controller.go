package controllers

import (
	"fmt"
	"libs"
	"libs/comment"
	"libs/passport"
	"libs/vod"
	"outobjs"
	"time"
	"utils"
)

// 评论模块 API
type CommentController struct {
	BaseController
}

func (c *CommentController) Prepare() {
	c.BaseController.Prepare()
}

func (c *CommentController) URLMapping() {
	c.Mapping("Publish", c.Publish)
	c.Mapping("Gets", c.Gets)
}

// @Title 发表评论
// @Description 发表评论(成功返回error_code:REP000,error_description:新评论编号)
// @Param   access_token  path  string  true  "access_token"
// @Param   mod   path  string  true  "所属模块(vod)"
// @Param   ref_id    path  int  true  "关联id(比如视频id号等)"
// @Param   reply_id    path  string  false  "回复评论(非回复则为空)"
// @Param   title    path  string  false  "标题(无则为空)"
// @Param   content    path  string  true  "回复内容,@系统自动处理"
// @Param   longitude    path  float  false  "经度"
// @Param   latitude    path  float  false  "纬度"
// @Param   is_msg    path  bool  false  "是否给@的用户发送消息(默认为true)"
// @Success 200 {object} libs.Error
// @router /publish [post]
func (c *CommentController) Publish() {
	uid := c.CurrentUid()
	if uid <= 0 {
		c.Json(libs.NewError("comment_publish_premission_denied", UNAUTHORIZED_CODE, "必须登录后才能评论", ""))
		return
	}
	mod := c.GetString("mod")
	ref_id, _ := c.GetInt64("ref_id")
	reply_id := c.GetString("reply_id")
	title, _ := utils.UrlDecode(c.GetString("title"))
	content, _ := utils.UrlDecode(c.GetString("content"))
	longitude, _ := c.GetFloat("longitude")
	latitude, _ := c.GetFloat("latitude")
	is_msg, err := c.GetBool("is_msg")
	if err != nil {
		is_msg = true
	}
	ip := c.Ctx.Input.IP()

	intervalKey := fmt.Sprintf("mobile_comment_user_interval:%d", uid)
	cache := utils.GetLocalCache()
	_tmp := 0
	err = cache.Get(intervalKey, &_tmp)
	if err == nil {
		c.Json(libs.NewError("comment_mod_empty", "C5006", "不能在过断的时间内连续评论", ""))
		return
	}

	if len(mod) == 0 {
		c.Json(libs.NewError("comment_mod_empty", "C5001", "mod参数不能为空", ""))
		return
	}
	if ref_id <= 0 {
		c.Json(libs.NewError("comment_mod_refid", "C5002", "ref_id参数必须大于0", ""))
		return
	}
	if len(content) == 0 {
		c.Json(libs.NewError("comment_mod_content", "C5004", "content参数不能为空", ""))
		return
	}
	if len(comment.GetCollectionName(mod)) == 0 {
		c.Json(libs.NewError("comment_mod_notexist", "C5005", "mod评论模块不存在", ""))
		return
	}
	commentor := comment.NewCommentor(mod)
	data := map[string]interface{}{
		"ref_id":      ref_id,
		"from_id":     int64(uid),
		"reply_id":    reply_id,
		"text":        content,
		"ip":          ip,
		"longitude":   longitude,
		"latitude":    latitude,
		"title":       title,
		"allow_reply": true,
	}
	var msgType = vod.MSG_TYPE_COMMENT
	ms := passport.NewMemberProvider()
	err, _id := commentor.Create(data,
		ms.GetUidByNickname,
		ms.GetNicknameByUid,
		is_msg,
		msgType)
	if err == nil {
		cache.Set(intervalKey, 1, comment.CommentPulishInterval())
		c.Json(libs.NewError("comment_mod_publish_succ", RESPONSE_SUCCESS, _id, ""))
		return
	}
	c.Json(libs.NewError("comment_mod_publish_fail", "C5010", err.Error(), ""))
}

// @Title 评论列表
// @Description 获取评论列表
// @Param   mod   path  string  true  "所属模块(vod)"
// @Param   ref_id    path  int  true  "关联id(比如视频id号等)"
// @Param   page    path  int  false  "页(默认1)"
// @Param   size    path  int  false  "页行数(默认20)"
// @Param   t 		path  int  false  "最后更改时间戳,由服务器生成,首次调用会回传t参数时间戳,客户端每次提交时带上,第一页时留空"
// @Param   nano    path  int  false "服务器端大数据分页使用,最后次获取数据时的last_nano值,如果第一页传空(此项暂时无视)"
// @Success 200 {object} outobjs.OutCommentPageList
// @router /list [get]
func (c *CommentController) Gets() {
	mod := c.GetString("mod")
	ref_id, _ := c.GetInt64("ref_id")
	page, _ := c.GetInt("page")
	size, _ := c.GetInt("size")
	timestamp, _ := c.GetInt64("t")
	nano, _ := c.GetInt64("nano")

	t := time.Now()
	if timestamp > 0 {
		t = time.Unix(timestamp, 0)
	}
	if len(mod) == 0 {
		c.Json(libs.NewError("comment_mod_empty", "C5001", "mod参数不能为空", ""))
		return
	}
	if ref_id <= 0 {
		c.Json(libs.NewError("comment_mod_refid", "C5002", "ref_id参数必须大于0", ""))
		return
	}
	if len(comment.GetCollectionName(mod)) == 0 {
		c.Json(libs.NewError("comment_mod_notexist", "C5005", "mod评论模块不存在", ""))
		return
	}
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	commentor := comment.NewCommentor(mod)
	total, comments, last_nano := commentor.Gets(mod, ref_id, page, size, t, nano)
	list := []*outobjs.OutComment{}
	for _, cmt := range comments {
		out_c := outobjs.GetOutComment(cmt)
		list = append(list, out_c)
	}
	out_list := outobjs.OutCommentPageList{
		Total:       total,
		TotalPage:   utils.TotalPages(int(total), int(size)),
		CurrentPage: int(page),
		Size:        int(size),
		Time:        t.Unix(),
		Lists:       list,
		LastNano:    last_nano,
	}
	c.Json(out_list)
}

// @Title 单条评论
// @Description 获取单条评论
// @Param   mod   path  string  true  "所属模块(vod)"
// @Param   id    path  string  true  "评论id"
// @Success 200 {object} outobjs.OutComment
// @router /get [get]
func (c *CommentController) Get() {
	mod := c.GetString("mod")
	id := c.GetString("id")
	if len(mod) == 0 {
		c.Json(libs.NewError("comment_mod_empty", "C5011", "mod参数不能为空", ""))
		return
	}
	if len(id) <= 0 {
		c.Json(libs.NewError("comment_mod_id", "C5012", "id参数错误", ""))
		return
	}
	if len(comment.GetCollectionName(mod)) == 0 {
		c.Json(libs.NewError("comment_mod_notexist", "C5013", "mod评论模块不存在", ""))
		return
	}
	commentor := comment.NewCommentor(mod)
	comment := commentor.Get(id)
	if comment == nil {
		c.Json(libs.NewError("comment_mod_refid", "C5014", "评论不存在", ""))
		return
	}
	c.Json(outobjs.GetOutComment(comment))
}
