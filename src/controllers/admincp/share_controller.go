package admincp

import (
	"controllers"
	"libs"
	"libs/share"
	"strconv"
	"strings"
	"utils"
)

// 视频管理 API
type ShareCPController struct {
	AdminController
}

func (c *ShareCPController) Prepare() {
	c.AdminController.Prepare()
}

// @Title 播主发布分享
// @Description 播主发布分享
// @Param   uids   path	string true  "播主uids,逗号,隔开"
// @Param   vod_ids   path	string true  "分享视频s,逗号,隔开"
// @Param   pic_ids   path	string true  "分享图片s,逗号,隔开"
// @Param   text   path   string  false  "文本内容"
// @Param   type   path   int  true  "分享类型(文本=1,视频=2)"
// @Param   source     path   string  false  "源地址(可忽略)"
// @Success 200 {object} libs.Error
// @router /publishs [post]
func (c *ShareCPController) Publishs() {
	uids_str := c.GetString("uids")
	vodids_str := c.GetString("vod_ids")
	pics := c.GetString("pic_ids")
	text, _ := utils.UrlDecode(c.GetString("text"))
	share_type, _ := c.GetInt("type")
	source := c.GetString("source")

	if share_type <= 0 {
		c.Json(libs.NewError("sys_share_parameter_fail", "GM050_040", "type参数错误", ""))
		return
	}

	uidss := strings.Split(uids_str, ",")
	uids := []int64{}
	for _, uid_str := range uidss {
		_uid, err := strconv.ParseInt(uid_str, 10, 64)
		if err == nil {
			uids = append(uids, _uid)
		}
	}

	nts := &share.Shares{}
	//文本形式
	if share_type == int(share.SHARE_KIND_TXT) {
		if len(text) == 0 {
			c.Json(libs.NewError("sys_share_parameter_fail", "GM050_041", "文本内容不能为空", ""))
			return
		}
		for _, from_uid := range uids {
			_s := share.Share{
				Uid:       from_uid,
				Source:    source,
				Text:      text,
				ShareType: share_type,
				Type:      share.SHARE_TYPE_ORIGINAL, //默认是原创
				Resources: "",
			}
			go nts.Create(&_s, true)
		}
		c.Json(libs.NewError("sys_share_create_succ", controllers.RESPONSE_SUCCESS, "分享成功", ""))
		return
	}
	//视频形式
	if share_type == int(share.SHARE_KIND_VOD) {
		vodidss := strings.Split(vodids_str, ",")
		vod_ids := []int64{}
		for _, vod_str := range vodidss {
			_vid, err := strconv.ParseInt(vod_str, 10, 64)
			if err == nil {
				vod_ids = append(vod_ids, _vid)
			}
		}
		if len(vod_ids) == 0 {
			c.Json(libs.NewError("sys_share_vod_empty", "GM050_042", "未提交视频或提交格式错误", ""))
			return
		}

		for _, from_uid := range uids {
			for _, _vid := range vod_ids {
				_share_vod_res := nts.TranformResource(share.SHARE_KIND_VOD, strconv.FormatInt(_vid, 10))
				if len(_share_vod_res) == 0 {
					continue
				}
				_s := share.Share{
					Uid:       from_uid,
					Source:    source,
					Text:      text,
					ShareType: share_type,
					Type:      share.SHARE_TYPE_ORIGINAL, //默认是原创
					Resources: _share_vod_res,
				}
				go nts.Create(&_s, true)
			}
		}
		c.Json(libs.NewError("sys_member_share_create_succ", controllers.RESPONSE_SUCCESS, "分享成功", ""))
		return
	}
	if share_type == int(share.SHARE_KIND_PIC) {
		_pic_ids := strings.Split(pics, ",")
		if len(_pic_ids) == 0 {
			c.Json(libs.NewError("sys_member_share_pic_empty", "GM050_043", "未提交图片或提交格式错误", ""))
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
			c.Json(libs.NewError("sys_member_share_pic_empty", "GM050_044", "未提交图片或提交格式错误", ""))
			return
		}
		for _, from_uid := range uids {
			_s := share.Share{
				Uid:       from_uid,
				Source:    source,
				Text:      text,
				ShareType: share_type,
				Type:      share.SHARE_TYPE_ORIGINAL, //默认是原创
				Resources: strings.Join(_share_pic_res, ","),
			}
			go nts.Create(&_s, true)
		}
		c.Json(libs.NewError("sys_member_share_create_succ", controllers.RESPONSE_SUCCESS, "分享成功", ""))
		return
	}
	c.Json(libs.NewError("sys_member_share_create_fail", "GM050_049", "没有适配的分享类型", ""))
}
