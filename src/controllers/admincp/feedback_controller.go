package admincp

import (
	"libs"
	"libs/feedback"
	"outobjs"
	"utils"
)

// 回馈管理模块 API
type FeedbackCPController struct {
	AdminController
	storage libs.IFileStorage
}

func (c *FeedbackCPController) Prepare() {
	c.AdminController.Prepare()
	c.storage = libs.NewFileStorage()
}

// @Title 用户回馈信息
// @Description 用户回馈信息
// @Param   page   path	int true  "页"
// @Param   size   path	int true  "页大小"
// @Success 200 {object} outobjs.OutFeedbackPagedList
// @router /list [get]
func (c *FeedbackCPController) List() {
	page, _ := c.GetInt("page")
	size, _ := c.GetInt("size")
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	total, fds := feedback.Gets(page, size)
	outs := []*outobjs.OutFeedback{}
	for _, fd := range fds {
		outs = append(outs, &outobjs.OutFeedback{
			Id:       fd.Id,
			Uid:      fd.Uid,
			Member:   outobjs.GetOutMember(fd.Uid, 0),
			Category: fd.Category,
			PostTime: fd.PostTime,
			Title:    fd.Title,
			Content:  fd.Content,
			Img:      fd.Img,
			ImgUrl:   c.storage.GetFileUrl(fd.Img),
			Contact:  fd.Contact,
			Source:   fd.Source,
		})
	}
	pds := &outobjs.OutFeedbackPagedList{
		Total:       total,
		TotalPage:   utils.TotalPages(total, size),
		CurrentPage: page,
		Size:        size,
		Lists:       outs,
	}
	c.Json(pds)
}
