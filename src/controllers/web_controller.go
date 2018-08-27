package controllers

import (
	"bytes"
	"fmt"
	"libs/feedback"
	"libs/lives"
	"libs/reptile"
	"libs/version"
	"libs/vod"
	"outobjs"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"utils"
)

const (
	MOBILE_AGENT_REGEX = `(android|bb\d+|meego).+mobile|avantgo|bada\/|blackberry|blazer|compal|elaine|fennec|hiptop|iemobile|ip(hone|od)|iris|kindle|lge |maemo|midp|mmp|mobile.+firefox|netfront|opera m(ob|in)i|palm( os)?|phone|p(ixi|re)\/|plucker|pocket|psp|series(4|6)0|symbian|treo|up\.(browser|link)|vodafone|wap|windows ce|xda|xiino|mobile`
	MOBILE_DOMAIN      = "m.dianjingquan.cn"
	WWW_DOMAIN         = "www.dianjingquan.cn"
)

type WebController struct {
	BaseController
	IsMobile bool
	IosDown  string
	AndDown  string
}

func (c *WebController) URLMapping() {
	c.Mapping("Home", c.Home)
}

func (c *WebController) Prepare() {
	agent := strings.ToLower(c.Ctx.Input.UserAgent())
	mobile, _ := regexp.MatchString(MOBILE_AGENT_REGEX, agent)
	c.IsMobile = mobile

	domain := strings.ToLower(c.Ctx.Input.Domain())
	port := c.Ctx.Input.Port()
	port_str := ""
	if port != 80 {
		port_str = ":" + strconv.Itoa(port)
	}
	if mobile && domain != MOBILE_DOMAIN {
		redirect_url := c.Ctx.Input.Scheme() + "://" + MOBILE_DOMAIN + port_str + c.Ctx.Input.URL()
		c.Ctx.Redirect(302, redirect_url)
		c.StopRun()
	}
	if !mobile && domain != WWW_DOMAIN {
		redirect_url := c.Ctx.Input.Scheme() + "://" + WWW_DOMAIN + port_str + c.Ctx.Input.URL()
		c.Ctx.Redirect(302, redirect_url)
		c.StopRun()
	}

	vcs := &version.VCS{}
	iosv := vcs.GetLastClientVersion(version.MOBILE_PLATFORM_APPLE)
	andv := vcs.GetLastClientVersion(version.MOBILE_PLATFORM_ANDROID)
	if iosv != nil {
		c.Data["IOS_DOWN"] = iosv.DownloadUrl
	}
	if andv != nil {
		c.Data["AND_DOWN"] = andv.DownloadUrl
	}
	if strings.Index(agent, "micromessenger") >= 0 {
		c.Data["IN_WEIXIN"] = true
	}
	plat := c.getPlat(agent)
	if plat == version.MOBILE_PLATFORM_ANDROID {
		c.Data["IS_ANDROID"] = true
	} else if plat == version.MOBILE_PLATFORM_APPLE {
		c.Data["IS_IOS"] = true
	}
}

// @router / [get]
func (c *WebController) Home() {
	if c.IsMobile {
		c.TplName = "m_home.html"
		return
	}
	c.TplName = "www_home.html"
}

func (c *WebController) getPlat(agent string) version.MOBILE_PLATFORM {
	var plat version.MOBILE_PLATFORM
	if strings.Index(agent, "iphone") >= 0 || strings.Index(agent, "mac") >= 0 || strings.Index(agent, "ipad") >= 0 || strings.Index(agent, "ipad") >= 0 {
		plat = version.MOBILE_PLATFORM_APPLE
	}
	if strings.Index(agent, "android") >= 0 || strings.Index(agent, "linux") >= 0 {
		plat = version.MOBILE_PLATFORM_ANDROID
	}
	return plat
}

// @router /download [get]
func (c *WebController) DownApp() {
	agent := strings.ToLower(c.Ctx.Input.UserAgent())
	plat := c.getPlat(agent)
	//	var plat version.MOBILE_PLATFORM
	//	if strings.Index(agent, "iphone") >= 0 || strings.Index(agent, "mac") >= 0 || strings.Index(agent, "ipad") >= 0 || strings.Index(agent, "ipad") >= 0 {
	//		plat = version.MOBILE_PLATFORM_APPLE
	//	}
	//	if strings.Index(agent, "android") >= 0 || strings.Index(agent, "linux") >= 0 {
	//		plat = version.MOBILE_PLATFORM_ANDROID
	//	}
	if plat == "" {
		c.Redirect("/", 302)
		c.StopRun()
	}
	vcs := &version.VCS{}
	lastv := vcs.GetLastClientVersion(plat)
	if lastv == nil {
		c.Redirect("/", 302)
		c.StopRun()
	}
	//if c.IsMobile {
	c.Redirect(lastv.DownloadUrl, 302)
	//	c.StopRun()
	//}
	//c.Redirect("/", 302)
}

// @router /down [get]
func (c *WebController) Down() {
	c.Redirect("/", 302)
}

// @router /feedback [post]
func (c *WebController) Feedback() {
	category, _ := utils.UrlDecode(c.GetString("category"))
	title, _ := utils.UrlDecode(c.GetString("title"))
	content, _ := utils.UrlDecode(c.GetString("content"))
	contact, _ := utils.UrlDecode(c.GetString("contact"))
	source, _ := utils.UrlDecode(c.GetString("source"))
	if len(source) == 0 {
		source = "web"
	}
	type OutResult struct {
		Code    int    `json:"code"`
		Content string `json:"content"`
	}
	if len(content) == 0 {
		c.Ctx.WriteString(`<script type="text/javascript">parent.callback("内容不能为空!")</script>`)
		return
	}
	f, h, _ := c.GetFile("file")
	file_name := h.Filename
	buf := &bytes.Buffer{}
	_, err := buf.ReadFrom(f)
	if err != nil {
		c.Ctx.WriteString(`<script type="text/javascript">parent.callback("上传的图片格式错误!")</script>`)
		return
	}
	node, err := file_storage.SaveFile(buf.Bytes(), file_name, 0)
	if node.ExtName != "jpg" && node.ExtName != "jpeg" && node.ExtName != "png" {
		c.Ctx.WriteString(`<script type="text/javascript">parent.callback("上传的图片格式错误!")</script>`)
		return
	}
	if err != nil {
		c.Ctx.WriteString(`<script type="text/javascript">parent.callback("上传的图片格式错误!")</script>`)
		return
	}
	fb := feedback.Feedback{
		Uid:      0,
		Category: category,
		Title:    title,
		Content:  content,
		Img:      node.FileId,
		Contact:  contact,
		Source:   source,
	}
	err = feedback.Create(&fb)
	if err == nil {
		c.Ctx.WriteString(`<script type="text/javascript">parent.callback("提交成功")</script>`)
		return
	}
	c.Ctx.WriteString(`<script type="text/javascript">parent.callback("提交失败!")</script>`)
}

// @router /privacy_protocol [get]
func (c *WebController) PrivacyProtocol() {
	c.TplName = "privacy_protocol.html"
}

// @router /vod/:id([0-9]+) [get]
func (c *WebController) VodPlay() {
	id, _ := c.GetInt64(":id")
	vods := &vod.Vods{}
	vod := vods.Get(id, true)
	if vod == nil {
		c.Redirect("/", 302)
		c.StopRun()
	}
	c.Data["Vod"] = vod
	out_m := outobjs.GetOutMember(vod.Uid, 0)
	if out_m != nil {
		c.Data["MAvatorUrl"] = out_m.AvatarUrl
		c.Data["MName"] = out_m.NickName
		c.Data["MVods"] = out_m.Vods
		c.Data["MFans"] = out_m.Fans
		c.Data["MGil"] = out_m.Credits
		c.Data["MVip"] = out_m.IsCertified
	}
	if c.IsMobile {
		flvs := vods.GetPlayFlvs(id, true)
		stream_url := ""
		if flvs != nil {
			stream_url = c.m_qx(flvs, reptile.VOD_STREAM_MODE_MSUPER)
		}
		c.Data["Stream"] = stream_url
		c.Data["Poster"] = file_storage.GetFileUrl(vod.Img)
		c.TplName = "m_vod.html"
		return
	}
	c.TplName = "www_vod.html"
}

// @router /vod/:id([0-9]+)/stream [get]
func (c *WebController) VodStream() {
	id, _ := c.GetInt64(":id")
	vods := &vod.Vods{}
	flvs := vods.GetPlayFlvs(id, true)
	if flvs == nil {
		return
	}
	xml := `<?xml version="1.0" encoding="utf-8"?>`
	xml += `<ckplayer><flashvars>{h->4}</flashvars>`
	opt := c.pc_qx(flvs, reptile.VOD_STREAM_MODE_SUPER_SP)
	if opt != nil {
		for _, flv := range opt.Flvs {
			xml += fmt.Sprintf(`<video><file>%s</file><size>%d</size><seconds>%d</seconds></video>`, utils.XmlEscape(flv.Url), flv.Size, flv.Seconds)
		}
	}
	xml += `</ckplayer>`
	c.Ctx.WriteString(xml)
}

// @router /plive/:id([0-9]+) [get]
func (c *WebController) PeronalLive() {
	if c.IsMobile {
		c.Redirect("/", 302)
		c.StopRun()
	}

	id, _ := c.GetInt64(":id")
	lps := &lives.LivePers{}
	per := lps.Get(id)
	if per == nil {
		c.Redirect("/", 302)
		c.StopRun()
	}
	t, ok := reptile.LIVE_REPTILE_MODULES[string(per.Rep)]
	if !ok {
		c.Redirect("/", 302)
		c.StopRun()
	}
	rep := reflect.New(t).Interface().(reptile.ILiveViewOnPc)
	html := rep.ViewHtmlOnPc(per.ReptileUrl, 800, 500)
	if len(html) == 0 {
		c.Redirect("/", 302)
		c.StopRun()
	}
	out_m := outobjs.GetOutMember(per.Uid, 0)
	if out_m != nil {
		c.Data["MAvatorUrl"] = out_m.AvatarUrl
		c.Data["MName"] = out_m.NickName
		c.Data["MVods"] = out_m.Vods
		c.Data["MFans"] = out_m.Fans
		c.Data["MGil"] = out_m.Credits
		c.Data["MVip"] = out_m.IsCertified
	}
	c.Data["ViewHtml"] = html
	c.Data["Title"] = per.Name
	c.TplName = "www_plive.html"
}

// @router /jlive/:id([0-9]+) [get]
func (c *WebController) JigouLive() {
	if c.IsMobile {
		c.Redirect("/", 302)
		c.StopRun()
	}

	id, _ := c.GetInt64(":id")
	orgs := &lives.LiveOrgs{}
	channel := orgs.GetChannel(id)
	if channel == nil {
		c.Redirect("/", 302)
		c.StopRun()
	}
	streams := orgs.GetStreams(id)
	var liveStream *lives.LiveStream
	for _, st := range streams {
		if st.Default {
			liveStream = st
		}
	}
	if liveStream == nil {
		c.Redirect("/", 302)
		c.StopRun()
	}
	t, ok := reptile.LIVE_REPTILE_MODULES[string(liveStream.Rep)]
	if !ok {
		c.Redirect("/", 302)
		c.StopRun()
	}
	rep := reflect.New(t).Interface().(reptile.ILiveViewOnPc)
	html := rep.ViewHtmlOnPc(liveStream.ReptileUrl, 800, 600)
	if len(html) == 0 {
		c.Redirect("/", 302)
		c.StopRun()
	}
	out_m := outobjs.GetOutMember(channel.Uid, 0)
	if out_m != nil {
		c.Data["MAvatorUrl"] = out_m.AvatarUrl
		c.Data["MName"] = out_m.NickName
		c.Data["MVods"] = out_m.Vods
		c.Data["MFans"] = out_m.Fans
		c.Data["MGil"] = out_m.Credits
		c.Data["MVip"] = out_m.IsCertified
	}
	c.Data["ViewHtml"] = html
	c.Data["Title"] = channel.Name
	c.TplName = "www_plive.html"
}

func (c *WebController) pc_qx2(opts []vod.VideoOpt, defMode reptile.VOD_STREAM_MODE) *vod.VideoOpt {
	vsts := make(map[reptile.VOD_STREAM_MODE]*vod.VideoOpt)
	for _, opt := range opts {
		if opt.Mode == reptile.VOD_STREAM_MODE_STANDARD_SP ||
			opt.Mode == reptile.VOD_STREAM_MODE_HIGH_SP ||
			opt.Mode == reptile.VOD_STREAM_MODE_SUPER_SP ||
			opt.Mode == reptile.VOD_STREAM_MODE_1080P_SP ||
			opt.Mode == reptile.VOD_STREAM_MODE_M1080P ||
			opt.Mode == reptile.VOD_STREAM_MODE_MSUPER ||
			opt.Mode == reptile.VOD_STREAM_MODE_MHIGH ||
			opt.Mode == reptile.VOD_STREAM_MODE_MSTD {
			if len(opt.Flvs) > 0 {
				vsts[opt.Mode] = &opt
			}
		}
	}
	if opt, ok := vsts[defMode]; ok {
		return opt
	}
	if opt, ok := vsts[reptile.VOD_STREAM_MODE_1080P_SP]; ok {
		return opt
	}
	if opt, ok := vsts[reptile.VOD_STREAM_MODE_SUPER_SP]; ok {
		return opt
	}
	if opt, ok := vsts[reptile.VOD_STREAM_MODE_HIGH_SP]; ok {
		return opt
	}
	if opt, ok := vsts[reptile.VOD_STREAM_MODE_STANDARD_SP]; ok {
		return opt
	}
	return nil
}

func (c *WebController) pc_qx(vpf *vod.VideoPlayFlvs, defMode reptile.VOD_STREAM_MODE) *vod.VideoOpt {
	vsts := make(map[reptile.VOD_STREAM_MODE]*vod.VideoOpt)
	for _, opt := range vpf.OptFlvs {
		if opt.Mode == reptile.VOD_STREAM_MODE_STANDARD_SP ||
			opt.Mode == reptile.VOD_STREAM_MODE_HIGH_SP ||
			opt.Mode == reptile.VOD_STREAM_MODE_SUPER_SP ||
			opt.Mode == reptile.VOD_STREAM_MODE_1080P_SP ||
			opt.Mode == reptile.VOD_STREAM_MODE_M1080P ||
			opt.Mode == reptile.VOD_STREAM_MODE_MSUPER ||
			opt.Mode == reptile.VOD_STREAM_MODE_MHIGH ||
			opt.Mode == reptile.VOD_STREAM_MODE_MSTD {
			if len(opt.Flvs) > 0 {
				vsts[opt.Mode] = &opt
			}
		}
	}
	if opt, ok := vsts[defMode]; ok {
		return opt
	}
	if opt, ok := vsts[reptile.VOD_STREAM_MODE_1080P_SP]; ok {
		return opt
	}
	if opt, ok := vsts[reptile.VOD_STREAM_MODE_M1080P]; ok {
		return opt
	}
	if opt, ok := vsts[reptile.VOD_STREAM_MODE_SUPER_SP]; ok {
		return opt
	}
	if opt, ok := vsts[reptile.VOD_STREAM_MODE_MSUPER]; ok {
		return opt
	}
	if opt, ok := vsts[reptile.VOD_STREAM_MODE_HIGH_SP]; ok {
		return opt
	}
	if opt, ok := vsts[reptile.VOD_STREAM_MODE_MHIGH]; ok {
		return opt
	}
	if opt, ok := vsts[reptile.VOD_STREAM_MODE_STANDARD_SP]; ok {
		return opt
	}
	if opt, ok := vsts[reptile.VOD_STREAM_MODE_MSTD]; ok {
		return opt
	}
	return nil
}

func (c *WebController) m_qx(vpf *vod.VideoPlayFlvs, defMode reptile.VOD_STREAM_MODE) string {
	vsts := make(map[reptile.VOD_STREAM_MODE]string)
	for _, opt := range vpf.OptFlvs {
		if opt.Mode == reptile.VOD_STREAM_MODE_MSTD ||
			opt.Mode == reptile.VOD_STREAM_MODE_MSTD_MP4 ||
			opt.Mode == reptile.VOD_STREAM_MODE_MHIGH ||
			opt.Mode == reptile.VOD_STREAM_MODE_MSUPER {
			if len(opt.Flvs) > 0 {
				vsts[opt.Mode] = opt.Flvs[0].Url
			}
		}
	}
	if url, ok := vsts[defMode]; ok {
		return url
	}
	if url, ok := vsts[reptile.VOD_STREAM_MODE_MSUPER]; ok {
		return url
	}
	if url, ok := vsts[reptile.VOD_STREAM_MODE_MHIGH]; ok {
		return url
	}
	if url, ok := vsts[reptile.VOD_STREAM_MODE_MSTD_MP4]; ok {
		return url
	}
	if url, ok := vsts[reptile.VOD_STREAM_MODE_MSTD]; ok {
		return url
	}
	return ""
}
