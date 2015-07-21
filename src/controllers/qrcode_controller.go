package controllers

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

}
