package controllers

import (
	"libs/qrcode"
	"outobjs"
	"utils"
)

// 扫码 API
type QRCodeController struct {
	BaseController
}

func (c *QRCodeController) Prepare() {
	c.BaseController.Prepare()
}

func (c *QRCodeController) URLMapping() {
	c.Mapping("Scan", c.Scan)
}

// @Title 扫码
// @Description 扫码
// @Param   access_token  path  string  false  "access_token"
// @Param   code   path  string  true  "码"
// @Success 200 {object} outobjs.OutQRCodeResult
// @router /scan [post]
func (c *QRCodeController) Scan() {
	fromUid := c.CurrentUid()
	code, _ := utils.UrlDecode(c.GetString("code"))
	if len(code) == 0 {
		c.Json(&outobjs.OutQRCodeResult{
			Result: "fail",
			Msg:    "二维码错误",
		})
		return
	}
	result, err := qrcode.DecodeCode(fromUid, code)
	if err != nil {
		c.Json(&outobjs.OutQRCodeResult{
			Result: "fail",
			Msg:    err.Error(),
		})
		return
	}
	c.Json(&outobjs.OutQRCodeResult{
		Mod:    result.Mod,
		Action: result.Action,
		Result: result.Result,
		Msg:    result.Msg,
		Args:   result.Args,
	})
}
