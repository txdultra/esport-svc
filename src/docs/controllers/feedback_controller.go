package controllers

import (
	"libs"
	"libs/feedback"
	"utils"
)

// 反馈 API
type FeedbackController struct {
	BaseController
}

func (c *FeedbackController) Prepare() {
	c.BaseController.Prepare()
}

func (c *FeedbackController) URLMapping() {
	c.Mapping("Submit", c.Submit)
}

// @Title 提交反馈信息
// @Description 提交反馈信息(error_code:REP000提交成功)
// @Param   access_token   path  string  false  "access_token"
// @Param   category   path  string  false  "分类"
// @Param   title   path  string  false  "反馈标题"
// @Param   content   path  string  true  "反馈内容"
// @Success 200 {object} libs.Error
// @router /submit [post]
func (c *FeedbackController) Submit() {
	uid := c.CurrentUid()
	category, _ := utils.UrlDecode(c.GetString("category"))
	title, _ := utils.UrlDecode(c.GetString("title"))
	content, _ := utils.UrlDecode(c.GetString("content"))
	fb := feedback.Feedback{
		Uid:      uid,
		Category: category,
		Title:    title,
		Content:  content,
	}
	err := feedback.Create(&fb)
	if err == nil {
		c.Json(libs.NewError("feedback_submit_success", RESPONSE_SUCCESS, "成功提交反馈信息", ""))
		return
	}
	c.Json(libs.NewError("feedback_submit_fail", "F1001", "失败:"+err.Error(), ""))
}
