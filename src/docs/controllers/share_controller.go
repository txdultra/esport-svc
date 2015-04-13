package controllers

import (
	"libs"
	"libs/share"
	"libs/vod"
	"outobjs"
	"strconv"
	"strings"
	"time"
	"utils"
)

// 用户分享 API
type ShareController struct {
	BaseController
}

func (c *ShareController) Prepare() {
	c.BaseController.Prepare()
}

func (c *ShareController) URLMapping() {
	c.Mapping("Publish", c.Publish)
	c.Mapping("Delete", c.Delete)
	c.Mapping("Timeline", c.Timeline)
	c.Mapping("PublicTimeline", c.PublicTimeline)
	c.Mapping("SubscriptionCount", c.SubscriptionCount)
	c.Mapping("MySubscr", c.MySubscr)
	c.Mapping("ShareComment", c.ShareComment)
	c.Mapping("DelShareComment", c.DelShareComment)
}

// @Title 发布分享
// @Description 发布分享 (成功返回error_code:REP000)
// @Param   access_token     path   string  true  "access_token"
// @Param   type     path   int  true  "分享类型(文本=1,视频=2,图片=4)"
// @Param   text     path   string  false  "文本内容"
// @Param   vod_ids     path   string  false  "一个或多个视频id(示例:1001,1002,1023)"
// @Param   pic_ids     path   string  false  "一个或多个图片id(示例:1001,1002,1023),仅支持JPEG、PNG格式"
// @Param   source     path   string  false  "源地址(可忽略)"
// @Success 200 成功返回error_code:REP000
// @router /publish [post]
func (c *ShareController) Publish() {
	_, err := c.ValidateAccessToken()
	if err != nil {
		c.Json(libs.NewError("member_share_premission_denied", UNAUTHORIZED_CODE, "请登陆后重新尝试", ""))
		return
	}
	uid := c.CurrentUid()
	s_type, _ := c.GetInt("type")
	text, _ := utils.UrlDecode(c.GetString("text"))
	vids := c.GetString("vod_ids")
	pics := c.GetString("pic_ids")
	source := c.GetString("source")
	share_type := int(s_type)
	if share_type <= 0 {
		c.Json(libs.NewError("member_share_parameter_fail", "S2002", "type参数错误", ""))
		return
	}
	nts := &share.Shares{}
	//文本形式
	if share_type == int(share.SHARE_KIND_TXT) {
		if len(text) == 0 {
			c.Json(libs.NewError("member_share_parameter_fail", "S2002", "文本内容不能为空", ""))
			return
		}
		_s := share.Share{
			Uid:       uid,
			Source:    source,
			Text:      text,
			ShareType: share_type,
			Type:      share.SHARE_TYPE_ORIGINAL, //默认是原创
			Resources: "",
		}
		err, _id := nts.Create(&_s, true)
		if err != nil {
			c.Json(libs.NewError("member_share_create", "S2003", "发布分享失败", ""))
			return
		}
		c.Json(libs.NewError("member_share_create_succ", RESPONSE_SUCCESS, "分享成功", strconv.FormatInt(_id, 10)))
		return
	}
	if share_type == int(share.SHARE_KIND_PIC) {
		_pic_ids := strings.Split(pics, ",")
		if len(_pic_ids) == 0 {
			c.Json(libs.NewError("member_share_pic_empty", "S2005", "未提交图片或提交格式错误", ""))
			return
		}
		_share_pic_res := []string{}
		for _, _pic := range _pic_ids {
			_picid, err := strconv.ParseInt(_pic, 10, 64)
			if err != nil {
				continue
			}
			if _picid > 0 {
				_share_pic_res = append(_share_pic_res, nts.TranformResource(share.SHARE_KIND_PIC, _pic))
			}
		}
		if len(_share_pic_res) == 0 {
			c.Json(libs.NewError("member_share_pic_empty", "S2005", "未提交图片或提交格式错误", ""))
			return
		}
		_s := share.Share{
			Uid:       uid,
			Source:    source,
			Text:      text,
			ShareType: share_type,
			Type:      share.SHARE_TYPE_ORIGINAL, //默认是原创
			Resources: strings.Join(_share_pic_res, ","),
		}
		err, _id := nts.Create(&_s, true)
		if err != nil {
			c.Json(libs.NewError("member_share_create", "S2006", "发布分享失败", ""))
			return
		}
		c.Json(libs.NewError("member_share_create_succ", RESPONSE_SUCCESS, "分享成功", strconv.FormatInt(_id, 10)))
		return
	}
	//视频形式
	if share_type == int(share.SHARE_KIND_VOD) {
		_vids := strings.Split(vids, ",")
		if len(_vids) == 0 {
			c.Json(libs.NewError("member_share_vod_empty", "S2004", "未提交视频或提交格式错误", ""))
			return
		}
		_share_vod_res := []string{}
		for _, _vid := range _vids {
			_share_vod_res = append(_share_vod_res, nts.TranformResource(share.SHARE_KIND_VOD, _vid))
		}
		_s := share.Share{
			Uid:       uid,
			Source:    source,
			Text:      text,
			ShareType: share_type,
			Type:      share.SHARE_TYPE_ORIGINAL, //默认是原创
			Resources: strings.Join(_share_vod_res, ","),
		}
		err, _id := nts.Create(&_s, true)
		if err != nil {
			c.Json(libs.NewError("member_share_create", "S2003", "发布分享失败", ""))
			return
		}
		c.Json(libs.NewError("member_share_create_succ", RESPONSE_SUCCESS, "分享成功", strconv.FormatInt(_id, 10)))
		return
	}
	c.Json(libs.NewError("member_share_type_notexist", "S2010", "分享类型不存在", ""))
}

// @Title 删除分享
// @Description 删除分享 (成功返回error_code:REP000)
// @Param   access_token     path   string  true  "access_token"
// @Param   share_id     path   string  true  "分享内容id编号"
// @Success 200 成功返回error_code:REP000
// @router /del [post]
func (c *ShareController) Delete() {
	_, err := c.ValidateAccessToken()
	if err != nil {
		c.Json(libs.NewError("member_share_premission_denied", UNAUTHORIZED_CODE, "请登陆后重新尝试", ""))
		return
	}
	share_id := c.GetString("share_id")
	_sid, serr := strconv.ParseInt(share_id, 10, 64)
	if serr != nil {
		c.Json(libs.NewError("member_share_id_fail", "S2101", "share_id格式错误", ""))
		return
	}

	nts := &share.Shares{}
	uid := c.CurrentUid()
	s := nts.Get(_sid)
	if s == nil {
		c.Json(libs.NewError("member_share_notexist", "S2103", "分享记录不存在", ""))
		return
	}
	if s.Uid != uid {
		c.Json(libs.NewError("member_share_permission_denied", "S2104", "您无权删除此分享记录", ""))
		return
	}

	del_err := nts.Delete(_sid)
	if del_err != nil {
		c.Json(libs.NewError("member_share_delete_fail", "S2102", del_err.Error(), ""))
		return
	}
	c.Json(libs.NewError("member_share_delete_succ", RESPONSE_SUCCESS, "成功删除分享记录", ""))
}

// @Title 获取所有的分享消息(朋友圈)
// @Description 获取所有的分享消息(朋友圈)
// @Param   access_token     path   string  true  "access_token"
// @Param   size     path    int  false  "一页行数(默认20)"
// @Param   page     path    int  false  "页数(默认1)"
// @Param   t     path    int  false  "时间戳(每次请求获得的t属性)"
// @Success 200  {object} outobjs.OutSharePageList
// @router /public_timeline [get]
func (c *ShareController) PublicTimeline() {
	current_uid := c.CurrentUid()
	size, _ := c.GetInt("size")
	page, _ := c.GetInt("page")
	timestamp, _ := c.GetInt64("t")

	if current_uid <= 0 {
		c.Json(libs.NewError("member_share_premission_denied", UNAUTHORIZED_CODE, "没有权限查询", ""))
		return
	}

	t := time.Now()
	if timestamp > 0 {
		t = utils.MsToTime(timestamp)
	}
	if size <= 0 {
		size = 20
	}
	if page <= 0 {
		page = 1
	}
	nts := &share.ShareMsgs{}
	total, lst := nts.Gets(current_uid, int(page), int(size), t)
	out_shares := []*outobjs.OutShare{}
	ts := t
	for _, share := range lst {
		out_shares = append(out_shares, c.tranfromOutShare(share, current_uid))
		if share.Ts < utils.TimeMillisecond(ts) {
			ts = utils.MsToTime(share.Ts)
		}
	}
	out := outobjs.OutSharePageList{
		Total:       total,
		TotalPage:   utils.TotalPages(int(total), int(size)),
		CurrentPage: int(page),
		Size:        int(size),
		Lists:       out_shares,
		Time:        utils.TimeMillisecond(ts),
	}
	c.Json(out)
}

// @Title 获取某人发布分享的记录
// @Description 获取某人发布分享的记录
// @Param   access_token     path   string  false  "access_token"
// @Param   uid   	 path    int  true  "查看分享用户的uid"
// @Param   size     path    int  false  "一页行数(默认20)"
// @Param   page     path    int  false  "页数(默认1)"
// @Param   t     path    int  false  "时间戳(每次请求获得的t属性)"
// @Success 200  {object} outobjs.OutSharePageList
// @router /timeline [get]
func (c *ShareController) Timeline() {
	current_uid := c.CurrentUid()
	uid, _ := c.GetInt64("uid")
	size, _ := c.GetInt("size")
	page, _ := c.GetInt("page")
	timestamp, _ := c.GetInt64("t")

	t := time.Now()
	if timestamp > 0 {
		t = utils.MsToTime(timestamp) //t = time.Unix(timestamp, 0)
	}
	if uid <= 0 {
		c.Json(libs.NewError("member_share_show", "S2201", "参数uid不能小于等于0", ""))
		return
	}
	if size <= 0 {
		size = 20
	}
	if page <= 0 {
		page = 1
	}
	nts := &share.Shares{}
	total, lst := nts.Gets(uid, int(page), int(size), t)
	out_shares := []*outobjs.OutShare{}
	ts := t
	for _, share := range lst {
		out_shares = append(out_shares, c.tranfromOutShare(share, current_uid))
		if share.Ts < utils.TimeMillisecond(ts) {
			ts = utils.MsToTime(share.Ts)
		}
	}
	out := outobjs.OutSharePageList{
		Total:       total,
		TotalPage:   utils.TotalPages(int(total), int(size)),
		CurrentPage: int(page),
		Size:        int(size),
		Lists:       out_shares,
		Time:        utils.TimeMillisecond(ts),
	}
	c.Json(out)
}

func (c *ShareController) tranfromOutShare(s *share.Share, suid int64) *outobjs.OutShare {
	out_share := &outobjs.OutShare{}
	out_share.Id = strconv.FormatInt(s.Id, 10)
	out_share.ShareType = s.ShareType
	out_share.Source = s.Source
	out_share.Geo = s.Geo
	out_share.Text = s.Text
	out_share.Ts = s.Ts
	out_share.CreateTime = s.CreateTime
	out_share.RepostCount = s.TransferredCount
	out_share.CommentsCount = s.CommentedCount
	out_share.AttitudesCount = s.AttitudedCount
	out_share.Member = outobjs.GetOutMember(s.Uid, suid)
	//ats := []outobjs.OutScreenNameUid{}
	//for k, v := range s.Ats {
	//	ats = append(ats, outobjs.OutScreenNameUid{
	//		ScreenName: k,
	//		Uid:        v,
	//	})
	//}
	//out_share.Ats = ats
	out_share.Vods = []*outobjs.OutShareVod{}
	out_share.Pics = []*outobjs.OutSharePic{}
	nts := &share.Shares{}
	vs := &vod.Vods{}
	pics := share.NewShareViewPics()
	resources := strings.Split(s.Resources, ",")
	for _, res := range resources {
		t, _id := nts.ResourceTranform(res)
		if t == share.SHARE_KIND_VOD {
			vid, _ := strconv.ParseInt(_id, 10, 64)
			vod := vs.Get(vid, false)
			vodc := vs.GetCount(vid)
			views := 0
			if vodc != nil {
				views = vodc.Views
			}
			if vod != nil {
				out_share.Vods = append(out_share.Vods, &outobjs.OutShareVod{
					Id:           vod.Id,
					Title:        vod.Title,
					ThumbnailPic: file_storage.GetFileUrl(vod.Img),
					BmiddlePic:   file_storage.GetFileUrl(vod.Img),
					OriginalPic:  file_storage.GetFileUrl(vod.Img),
					Views:        views,
					Member:       outobjs.GetOutMember(vod.Uid, suid),
				})
			}
		}
		if t == share.SHARE_KIND_PIC {
			picid, _ := strconv.ParseInt(_id, 10, 64)
			viewPics := pics.Get(picid)
			if len(viewPics) > 0 {
				outpic := &outobjs.OutSharePic{}
				outpic.Id = picid
				for ps, vpic := range viewPics {
					if ps == share.SHARE_PIC_SIZE_ORIGINAL {
						outpic.OriginalPic = file_storage.GetFileUrl(vpic.FileId)
					}
					if ps == share.SHARE_PIC_SIZE_THUMBNAIL {
						outpic.ThumbnailPic = file_storage.GetFileUrl(vpic.FileId)
					}
					if ps == share.SHARE_PIC_SIZE_MIDDLE {
						outpic.BmiddlePic = file_storage.GetFileUrl(vpic.FileId)
					}
				}
				out_share.Pics = append(out_share.Pics, outpic)
			}
		}
	}
	return out_share
}

// @Title 获取订阅收到的记录数
// @Description 获取订阅收到的记录数(成功返回error_code:REP000,error_description:数量[int类型])
// @Param   access_token     path   string  true  "access_token"
// @Success 200  成功返回error_code:REP000,error_description:数量[int类型]
// @router /subscr/count [get]
func (c *ShareController) SubscriptionCount() {
	uid := c.CurrentUid()
	if uid <= 0 {
		c.Json(libs.NewError("member_share_premission_denied", UNAUTHORIZED_CODE, "没有权限查询", ""))
		return
	}
	nts := &share.ShareVodSubcurs{}
	count := nts.NewEventCounts(uid)
	c.Json(libs.NewError("member_share_sub_count", RESPONSE_SUCCESS, strconv.Itoa(count), ""))
}

// @Title 获取订阅收到的记录
// @Description 获取订阅收到的记录(自动清空计数)
// @Param   access_token     path   string  true  "access_token"
// @Param   size     path    int  false  "一页行数(默认20)"
// @Param   page     path    int  false  "页数(默认1)"
// @Param   t     path    int  false  "时间戳(每次请求获得的t属性)"
// @Success 200  {object} outobjs.OutSharePageList
// @router /subscr/my [get]
func (c *ShareController) MySubscr() {
	uid := c.CurrentUid()
	if uid <= 0 {
		c.Json(libs.NewError("member_share_premission_denied", UNAUTHORIZED_CODE, "没有权限查询", ""))
		return
	}
	size, _ := c.GetInt("size")
	page, _ := c.GetInt("page")
	timestamp, _ := c.GetInt64("t")

	t := time.Now()
	if timestamp > 0 {
		t = utils.MsToTime(timestamp) //time.Unix(timestamp, 0)
	}
	if size <= 0 {
		size = 20
	}
	if page <= 0 {
		page = 1
	}
	ts := t
	nts := &share.ShareVodSubcurs{}
	total, lst := nts.Gets(uid, int(page), int(size), t)
	out_shares := []*outobjs.OutShare{}
	for _, share := range lst {
		out_shares = append(out_shares, c.tranfromOutShare(share, uid))
		if share.Ts < utils.TimeMillisecond(ts) {
			ts = utils.MsToTime(share.Ts)
		}
	}
	out := outobjs.OutSharePageList{
		Total:       total,
		TotalPage:   utils.TotalPages(int(total), int(size)),
		CurrentPage: int(page),
		Size:        int(size),
		Lists:       out_shares,
		Time:        utils.TimeMillisecond(ts),
	}
	//自动清空计数器
	nts.ResetNewEventCounts(uid)

	c.Json(out)
}

// @Title 评论分享
// @Description 评论分享
// @Param   access_token     path   string  true  "access_token"
// @Param   sid   	 path    int  true  "分享的id"
// @Param   content     path    int  true  "内容"
// @Param   ruid     path    int  false  "回复某人的uid(留空表示不是回复)"
// @Success 200  成功返回error_code:REP000
// @router /comment [post]
func (c *ShareController) ShareComment() {
	current_uid := c.CurrentUid()
	content, _ := utils.UrlDecode(c.GetString("content"))
	sid, _ := c.GetInt64("sid")
	ruid, _ := c.GetInt64("ruid")

	if sid <= 0 {
		c.Json(libs.NewError("member_share_comment_fail", "S2301", "参数sid不能小于等于0", ""))
		return
	}
	if len(content) == 0 {
		c.Json(libs.NewError("member_share_comment_fail", "S2302", "评论内容不能为空", ""))
		return
	}

	sc := &share.ShareComment{
		Sid:     sid,
		Uid:     current_uid,
		RUid:    ruid,
		Content: content,
	}
	scs := share.NewShareComments()
	id, err := scs.Create(sc)
	if err == nil {
		c.Json(libs.NewError("member_share_comment_succ", RESPONSE_SUCCESS, strconv.FormatInt(id, 10), ""))
		return
	}
	c.Json(libs.NewError("member_share_comment_fail", "S2309", err.Error(), ""))
}

// @Title 删除分享的评论
// @Description 删除分享的评论
// @Param   access_token     path   string  true  "access_token"
// @Param   id   	 path    int  true  "评论id"
// @Success 200  成功返回error_code:REP000
// @router /comment/del [delete]
func (c *ShareController) DelShareComment() {
	current_uid := c.CurrentUid()
	id, _ := c.GetInt64("id")

	if id <= 0 {
		c.Json(libs.NewError("member_share_comment_del_fail", "S2401", "id非法", ""))
		return
	}
	scs := share.NewShareComments()
	err := scs.Delete(id, current_uid)
	if err == nil {
		c.Json(libs.NewError("member_share_comment_del_succ", RESPONSE_SUCCESS, "删除成功", ""))
		return
	}
	c.Json(libs.NewError("member_share_comment_del_fail", "S2302", err.Error(), ""))
}
