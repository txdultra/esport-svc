package controllers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"libs"
	"libs/version"
	"os"
	"outobjs"
	"path/filepath"
	"strconv"
	"time"
	"utils"

	beeu "github.com/astaxie/beego/utils"
	"github.com/astaxie/beego/utils/captcha"
)

// 基础信息 API
type CommonController struct {
	BaseController
}

func (c *CommonController) Prepare() {
	c.BaseController.Prepare()
}

func (c *CommonController) URLMapping() {
	c.Mapping("VerifyPic", c.VerifyPic)
	c.Mapping("Games", c.Games)
	c.Mapping("Match", c.Match)
	c.Mapping("Matchs", c.Matchs)
	c.Mapping("Version", c.Version)
	//c.Mapping("ApiHosts", c.ApiHosts)
	c.Mapping("Expressions", c.Expressions)
}

// @Title 随机验证图片
// @Description 随机验证图片(PNG格式),验证码有效期5分钟
// @Param   nums   path  int  true  "字符串数量(最小3，最大6)"
// @Param   h   path  int  true  "图片高度(默认100)"
// @Param   w   path  int  true  "图片宽度(默认100)"
// @Success 200 输出{"code":"abc","img":"Base64图片数据"}
// @router /verify_pic [get]
func (c *CommonController) VerifyPic() {
	challengeNums, _ := c.GetInt("nums")
	height, _ := c.GetInt("h")
	width, _ := c.GetInt("w")

	_nums := int(challengeNums)
	_height := int(height)
	_width := int(width)

	if challengeNums < 3 {
		_nums = 3
	}
	if challengeNums > 6 {
		_nums = 6
	}
	if _height <= 0 {
		_height = 100
	}
	if _width <= 0 {
		_width = 100
	}
	d := beeu.RandomCreateBytes(_nums, []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}...)
	img := captcha.NewImage(d, _width, _height)
	buf := bytes.NewBuffer(nil)
	img.WriteTo(buf)

	digits := ""
	for _, v := range d {
		digits += strconv.Itoa(int(v))
	}

	outMap := make(map[string]string)
	rndStr := utils.RandomStrings(10)
	imgStr := utils.ToBase64(buf.Bytes())
	outMap["code"] = rndStr
	outMap["img"] = imgStr
	//加入缓存
	cache := utils.GetCache()
	cache.Set(rndStr, digits, 5*time.Minute)

	c.Json(outMap)
}

// @Title 获取版本信息
// @Description 获取版本信息,new_version属性为最新版本,无最新版则为空
// @Param   plat   path  string  true  "平台标识(android,ios,wphone)"
// @Param   ver    path  string  true  "客户端当前版本号名称"
// @Param   channel    path  string  false  "渠道"
// @Success 200 {object} outobjs.OutVersion
// @router /version [get]
func (c *CommonController) Version() {
	plat := c.GetString("plat")
	if len(plat) == 0 {
		c.Json(libs.NewError("common_version_parameter_error", "C2001", "plat参数错误", ""))
		return
	}
	ver := c.GetString("ver")
	if len(ver) == 0 {
		c.Json(libs.NewError("common_version_parameter_error", "C2001", "ver参数错误", ""))
		return
	}
	vcs := &version.VCS{}
	err, plt := vcs.ConvertPlatform(plat)
	if err != nil {
		c.Json(libs.NewError("common_version_plat_notexist", "C2002", "平台标识不存在", ""))
		return
	}
	current_ver := vcs.GetClientVersion(plt, ver)
	if current_ver == nil {
		c.Json(libs.NewError("common_version_notexist", "C2003", "版本号不存在", ""))
		return
	}
	lastv := vcs.GetLastClientVersion(plt)

	var last_out *outobjs.OutNewVersion = nil
	if lastv != nil && lastv.Version > current_ver.Version {
		last_out = &outobjs.OutNewVersion{
			Ver:         lastv.Ver,
			VerName:     lastv.VerName,
			Description: lastv.Description,
			Platform:    lastv.Platform,
			IsExpried:   lastv.IsExpried,
			DownloadUrl: lastv.DownloadUrl,
		}
	}
	out := outobjs.OutVersion{
		Ver:         current_ver.Ver,
		VerName:     current_ver.VerName,
		Description: current_ver.Description,
		Platform:    current_ver.Platform,
		IsExpried:   current_ver.IsExpried,
		DownloadUrl: current_ver.DownloadUrl,
		NewVersion:  last_out,
	}
	c.Json(out)
}

// @Title 获取所有游戏
// @Description 游戏列表(返回数组)
// @Success 200 {object} outobjs.OutGame
// @router /games [get]
func (c *CommonController) Games() {
	bas := &libs.Bas{}
	games := bas.Games()
	outs := []*outobjs.OutGame{}
	for _, v := range games {
		if v.Enabled {
			_g := outobjs.GetOutGame(&v)
			outs = append(outs, _g)
		}
	}
	c.Json(outs)
}

// @Title 获取赛事信息
// @Description 赛事信息
// @Param   match_id     path   int  true  "赛事id"
// @Success 200 {object} outobjs.OutMatch
// @router /match [get]
func (c *CommonController) Match() {
	mid, _ := c.GetInt("match_id")
	if mid <= 0 {
		c.Json(libs.NewError("common_parameter", "C1001", "参数错误", ""))
		return
	}
	out_match := outobjs.GetOutMatchById(int(mid))
	if out_match == nil {
		c.Json(libs.NewError("common_match_notexist", "C1002", "赛事不存在", ""))
		return
	}
	c.Json(out_match)
}

// @Title 获取所有赛事列表
// @Description 所有赛事列表
// @Success 200 {object} outobjs.OutMatch
// @router /matchs [get]
func (c *CommonController) Matchs() {
	bas := &libs.Bas{}
	matchs := bas.Matchs()
	outs := []*outobjs.OutMatch{}
	for _, v := range matchs {
		outs = append(outs, outobjs.GetOutMatch(&v))
	}
	c.Json(outs)
}

// @Title 获取表情列表(返回数组)
// @Description 获取表情列表(返回数组)
// @Param   category   path   string  true  "表情类型(默认public)"
// @Success 200 {object} outobjs.OutSmiley
// @router /expressions [get]
func (c *CommonController) Expressions() {
	categroy := c.GetString("category")
	if len(categroy) == 0 {
		categroy = "public"
	}

	//缓存
	cache_key := fmt.Sprintf("front_fast_cache.common_expressions.query:category_%s", categroy)
	c_obj := utils.GetLocalFastExpriesTimePartCache(cache_key)
	if c_obj != nil {
		c.Json(c_obj)
		return
	}

	bas := &libs.Bas{}
	smileies := bas.GetSmileies(categroy)
	out_smileies := []*outobjs.OutSmiley{}
	for _, sm := range smileies {
		if sm.Img > 0 {

		} else if len(sm.ImgPath) > 0 {
			workPath, _ := os.Getwd()
			workPath, _ = filepath.Abs(workPath)
			imgPath := filepath.Join(workPath, sm.ImgPath)
			if utils.FileExists(imgPath) {
				ds, err := ioutil.ReadFile(imgPath)
				if err == nil {
					bins := utils.ToBase64(ds)
					out_smileies = append(out_smileies, &outobjs.OutSmiley{
						Ext:          utils.FileExtName(imgPath),
						Code:         sm.Code,
						Bins:         bins,
						Points:       sm.Points,
						Category:     sm.Category,
						DisplayOrder: sm.DisplayOrder,
					})
				}
			}
		} else {
			continue
		}
	}
	utils.SetLocalFastExpriesTimePartCache(24*time.Hour, cache_key, out_smileies)
	c.Json(out_smileies)
}
