package controllers

import (
	"fmt"
	"libs/lives"
	"libs/reptile"
	"libs/vod"
	"outobjs"
	"reflect"
	"regexp"
	"strings"
	"utils"
)

const (
	MOBILE_AGENT_REGEX = `(android|bb\d+|meego).+mobile|avantgo|bada\/|blackberry|blazer|compal|elaine|fennec|hiptop|iemobile|ip(hone|od)|iris|kindle|lge |maemo|midp|mmp|mobile.+firefox|netfront|opera m(ob|in)i|palm( os)?|phone|p(ixi|re)\/|plucker|pocket|psp|series(4|6)0|symbian|treo|up\.(browser|link)|vodafone|wap|windows ce|xda|xiino`
	MOBILE_DOMAIN      = "m.dianjingquan.cn"
	WWW_DOMAIN         = "www.dianjingquan.cn"
)

type WebController struct {
	BaseController
	IsMobile bool
}

func (c *WebController) URLMapping() {
	c.Mapping("Home", c.Home)
}

func (c *WebController) Prepare() {
	agent := strings.ToLower(c.Ctx.Input.UserAgent())
	mobile, _ := regexp.MatchString(MOBILE_AGENT_REGEX, agent)
	c.IsMobile = mobile
	//domain := strings.ToLower(c.Ctx.Input.Domain())
	//port := c.Ctx.Input.Port()
	//port_str := ""
	//if port != 80 {
	//	port_str = ":" + strconv.Itoa(port)
	//}
	//if mobile && domain != MOBILE_DOMAIN {
	//	redirect_url := c.Ctx.Input.Scheme() + "://" + MOBILE_DOMAIN + port_str + c.Ctx.Input.Url()
	//	c.Ctx.Redirect(302, redirect_url)
	//	c.StopRun()
	//}
	//if !mobile && domain != WWW_DOMAIN {
	//	redirect_url := c.Ctx.Input.Scheme() + "://" + WWW_DOMAIN + port_str + c.Ctx.Input.Url()
	//	c.Ctx.Redirect(302, redirect_url)
	//	c.StopRun()
	//}
}

// @router / [get]
func (c *WebController) Home() {
	if c.IsMobile {
		c.TplNames = "m_home.html"
		return
	}
	c.TplNames = "www_home.html"
}

// @router /privacy_protocol [get]
func (c *WebController) PrivacyProtocol() {
	c.TplNames = "privacy_protocol.html"
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
	}
	if c.IsMobile {
		flvs := vods.GetPlayFlvs(id, true)
		stream_url := ""
		if flvs != nil {
			stream_url = c.m_qx(flvs, reptile.VOD_STREAM_MODE_MSUPER)
		}
		c.Data["Stream"] = stream_url
		c.Data["Poster"] = file_storage.GetFileUrl(vod.Img)
		c.TplNames = "m_vod.html"
		return
	}
	c.TplNames = "www_vod.html"
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
	}
	c.Data["ViewHtml"] = html
	c.Data["Title"] = per.Name
	c.TplNames = "www_plive.html"
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
	html := rep.ViewHtmlOnPc(liveStream.ReptileUrl, 600, 800)
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
	}
	c.Data["ViewHtml"] = html
	c.Data["Title"] = channel.Name
	c.TplNames = "www_plive.html"
}

func (c *WebController) pc_qx(vpf *vod.VideoPlayFlvs, defMode reptile.VOD_STREAM_MODE) *vod.VideoOpt {
	vsts := make(map[reptile.VOD_STREAM_MODE]*vod.VideoOpt)
	for _, opt := range vpf.OptFlvs {
		if opt.Mode == reptile.VOD_STREAM_MODE_STANDARD_SP ||
			opt.Mode == reptile.VOD_STREAM_MODE_HIGH_SP ||
			opt.Mode == reptile.VOD_STREAM_MODE_SUPER_SP ||
			opt.Mode == reptile.VOD_STREAM_MODE_1080P_SP {
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
