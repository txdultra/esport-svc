package controllers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"libs"
	//"libs/matchrace"
	"libs/matchrace"
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
	c.Mapping("HomeAd", c.HomeAd)
	c.Mapping("MatchRaces", c.MatchRaces)
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
			Ver:          lastv.Ver,
			VerName:      lastv.VerName,
			Description:  lastv.Description,
			Platform:     lastv.Platform,
			IsExpried:    lastv.IsExpried,
			DownloadUrl:  lastv.DownloadUrl,
			AllowVodDown: lastv.AllowVodDown,
		}
	}
	out := outobjs.OutVersion{
		Ver:          current_ver.Ver,
		VerName:      current_ver.VerName,
		Description:  current_ver.Description,
		Platform:     current_ver.Platform,
		IsExpried:    current_ver.IsExpried,
		DownloadUrl:  current_ver.DownloadUrl,
		AllowVodDown: current_ver.AllowVodDown,
		NewVersion:   last_out,
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

// @Title 获取最新首页广告
// @Description 获取最新首页广告
// @Success 200 {object} outobjs.OutHomeAd
// @router /homead [get]
func (c *CommonController) HomeAd() {
	//缓存
	cache_key := "front_fast_cache.common_homead.last"
	c_obj := utils.GetLocalFastExpriesTimePartCache(cache_key)
	if c_obj != nil {
		c.Json(c_obj)
		return
	}
	bas := &libs.Bas{}
	ad := bas.LastNewHomeAd()
	if ad == nil {
		c.Json(libs.NewError("common_homead_notnew", "C1010", "没有最新的广告", ""))
		return
	}
	if ad.EndTime.Before(time.Now()) {
		c.Json(libs.NewError("common_homead_notnew", "C1010", "没有最新的广告", ""))
		return
	}
	out_c := &outobjs.OutHomeAd{
		Id:      ad.Id,
		Title:   ad.Title,
		Img:     ad.Img,
		ImgUrl:  file_storage.GetFileUrl(ad.Img),
		Action:  ad.Action,
		Args:    ad.Args,
		Waits:   ad.Waits,
		EndTime: ad.EndTime,
	}
	utils.SetLocalFastExpriesTimePartCache(1*time.Minute, cache_key, out_c)
	c.Json(out_c)
}

// @Title 赛程
// @Description 赛程(返回数组)
// @Param   match_id    path  int  true  "赛事id"
// @Success 200 {object} outobjs.OutMatchMode
// @router /match_races [get]
func (c *CommonController) MatchRaces() {
	matchId, _ := c.GetInt64("match_id")
	if matchId <= 0 {
		c.Json(libs.NewError("common_match_race", "C1020", "赛事id错误", ""))
		return
	}
	//races := &matchrace.MatchRaceService{}
	//rs := races.GetModes(matchId)
	outs := []*outobjs.OutMatchMode{
		&outobjs.OutMatchMode{
			Id:       1,
			MatchId:  1,
			ModeType: matchrace.MODE_TYPE_GROUP,
			Title:    "组模型测试",
			Groups: []*outobjs.OutMatchGroup{
				&outobjs.OutMatchGroup{
					Id:    1,
					Title: "小组A",
					Players: []*outobjs.OutMatchGroupPlayer{
						&outobjs.OutMatchGroupPlayer{
							Id:       1,
							GroupId:  1,
							PlayerId: 1,
							Player: &outobjs.OutMatchPlayer{
								Id:     1,
								Name:   "选手A",
								Img:    0,
								ImgUrl: "",
							},
							Wins:   3,
							Pings:  0,
							Loses:  0,
							Points: 9,
							Outlet: true,
							Vss: []*outobjs.OutMatchVs{
								&outobjs.OutMatchVs{
									Id:      1,
									A:       1,
									AName:   "选手A",
									AImg:    0,
									AImgUrl: "",
									AScore:  1,
									AOutlet: false,
									B:       2,
									BName:   "选手B",
									BImg:    0,
									BImgUrl: "",
									BScore:  0,
									BOutlet: false,
								},
								&outobjs.OutMatchVs{
									Id:      2,
									A:       1,
									AName:   "选手A",
									AImg:    0,
									AImgUrl: "",
									AScore:  2,
									AOutlet: false,
									B:       3,
									BName:   "选手C",
									BImg:    0,
									BImgUrl: "",
									BScore:  0,
									BOutlet: false,
								},
								&outobjs.OutMatchVs{
									Id:      3,
									A:       1,
									AName:   "选手A",
									AImg:    0,
									AImgUrl: "",
									AScore:  2,
									AOutlet: false,
									B:       4,
									BName:   "选手D",
									BImg:    0,
									BImgUrl: "",
									BScore:  0,
									BOutlet: false,
								},
							},
						},
						&outobjs.OutMatchGroupPlayer{
							Id:       2,
							GroupId:  1,
							PlayerId: 2,
							Player: &outobjs.OutMatchPlayer{
								Id:     2,
								Name:   "选手B",
								Img:    0,
								ImgUrl: "",
							},
							Wins:   2,
							Pings:  0,
							Loses:  1,
							Points: 6,
							Outlet: true,
							Vss: []*outobjs.OutMatchVs{
								&outobjs.OutMatchVs{
									Id:      4,
									A:       1,
									AName:   "选手A",
									AImg:    0,
									AImgUrl: "",
									AScore:  1,
									AOutlet: false,
									B:       2,
									BName:   "选手B",
									BImg:    0,
									BImgUrl: "",
									BScore:  0,
									BOutlet: false,
								},
								&outobjs.OutMatchVs{
									Id:      5,
									A:       2,
									AName:   "选手B",
									AImg:    0,
									AImgUrl: "",
									AScore:  2,
									AOutlet: false,
									B:       3,
									BName:   "选手C",
									BImg:    0,
									BImgUrl: "",
									BScore:  0,
									BOutlet: false,
								},
								&outobjs.OutMatchVs{
									Id:      6,
									A:       2,
									AName:   "选手B",
									AImg:    0,
									AImgUrl: "",
									AScore:  2,
									AOutlet: false,
									B:       4,
									BName:   "选手D",
									BImg:    0,
									BImgUrl: "",
									BScore:  0,
									BOutlet: false,
								},
							},
						},
						&outobjs.OutMatchGroupPlayer{
							Id:       3,
							GroupId:  1,
							PlayerId: 3,
							Player: &outobjs.OutMatchPlayer{
								Id:     3,
								Name:   "选手C",
								Img:    0,
								ImgUrl: "",
							},
							Wins:   0,
							Pings:  2,
							Loses:  1,
							Points: 2,
							Outlet: false,
							Vss: []*outobjs.OutMatchVs{
								&outobjs.OutMatchVs{
									Id:      7,
									A:       1,
									AName:   "选手C",
									AImg:    0,
									AImgUrl: "",
									AScore:  1,
									AOutlet: false,
									B:       2,
									BName:   "选手B",
									BImg:    0,
									BImgUrl: "",
									BScore:  1,
									BOutlet: false,
								},
								&outobjs.OutMatchVs{
									Id:      8,
									A:       2,
									AName:   "选手B",
									AImg:    0,
									AImgUrl: "",
									AScore:  0,
									AOutlet: false,
									B:       3,
									BName:   "选手C",
									BImg:    0,
									BImgUrl: "",
									BScore:  0,
									BOutlet: false,
								},
								&outobjs.OutMatchVs{
									Id:      9,
									A:       2,
									AName:   "选手C",
									AImg:    0,
									AImgUrl: "",
									AScore:  1,
									AOutlet: false,
									B:       4,
									BName:   "选手D",
									BImg:    0,
									BImgUrl: "",
									BScore:  1,
									BOutlet: false,
								},
							},
						},
					},
				},
				&outobjs.OutMatchGroup{
					Id:    2,
					Title: "小组B",
					Players: []*outobjs.OutMatchGroupPlayer{
						&outobjs.OutMatchGroupPlayer{
							Id:       4,
							GroupId:  2,
							PlayerId: 4,
							Player: &outobjs.OutMatchPlayer{
								Id:     4,
								Name:   "选手D",
								Img:    0,
								ImgUrl: "",
							},
							Wins:   0,
							Pings:  0,
							Loses:  0,
							Points: 0,
							Outlet: false,
							Vss:    []*outobjs.OutMatchVs{},
						},
						&outobjs.OutMatchGroupPlayer{
							Id:       5,
							GroupId:  2,
							PlayerId: 5,
							Player: &outobjs.OutMatchPlayer{
								Id:     5,
								Name:   "选手E",
								Img:    0,
								ImgUrl: "",
							},
							Wins:   0,
							Pings:  0,
							Loses:  0,
							Points: 0,
							Outlet: false,
							Vss:    []*outobjs.OutMatchVs{},
						},
						&outobjs.OutMatchGroupPlayer{
							Id:       6,
							GroupId:  2,
							PlayerId: 6,
							Player: &outobjs.OutMatchPlayer{
								Id:     4,
								Name:   "选手F",
								Img:    0,
								ImgUrl: "",
							},
							Wins:   0,
							Pings:  0,
							Loses:  0,
							Points: 0,
							Outlet: false,
							Vss:    []*outobjs.OutMatchVs{},
						},
						&outobjs.OutMatchGroupPlayer{
							Id:       7,
							GroupId:  2,
							PlayerId: 7,
							Player: &outobjs.OutMatchPlayer{
								Id:     7,
								Name:   "选手G",
								Img:    0,
								ImgUrl: "",
							},
							Wins:   0,
							Pings:  0,
							Loses:  0,
							Points: 0,
							Outlet: false,
							Vss:    []*outobjs.OutMatchVs{},
						},
					},
				},
				&outobjs.OutMatchGroup{
					Id:      3,
					Title:   "小组C",
					Players: []*outobjs.OutMatchGroupPlayer{},
				},
				&outobjs.OutMatchGroup{
					Id:      4,
					Title:   "小组D",
					Players: []*outobjs.OutMatchGroupPlayer{},
				},
				&outobjs.OutMatchGroup{
					Id:      5,
					Title:   "小组E",
					Players: []*outobjs.OutMatchGroupPlayer{},
				},
				&outobjs.OutMatchGroup{
					Id:      6,
					Title:   "小组F",
					Players: []*outobjs.OutMatchGroupPlayer{},
				},
				&outobjs.OutMatchGroup{
					Id:      7,
					Title:   "小组G",
					Players: []*outobjs.OutMatchGroupPlayer{},
				},
				&outobjs.OutMatchGroup{
					Id:      8,
					Title:   "小组H",
					Players: []*outobjs.OutMatchGroupPlayer{},
				},
			},
		},
		&outobjs.OutMatchMode{
			Id:       2,
			MatchId:  1,
			ModeType: matchrace.MODE_TYPE_RECENT,
			Title:    "胜负平测试",
			Recents: []*outobjs.OutMatchRecent{
				&outobjs.OutMatchRecent{
					Id:       1,
					PlayerId: 1,
					Player: &outobjs.OutMatchPlayer{
						Id:   1,
						Name: "选手名称A",
					},
					M1: matchrace.WLP_L,
					M2: matchrace.WLP_L,
					M3: matchrace.WLP_W,
					M4: matchrace.WLP_UNDEFINED,
					M5: matchrace.WLP_W,
				},
				&outobjs.OutMatchRecent{
					Id:       2,
					PlayerId: 2,
					Player: &outobjs.OutMatchPlayer{
						Id:   2,
						Name: "选手名称B",
					},
					M1: matchrace.WLP_UNDEFINED,
					M2: matchrace.WLP_W,
					M3: matchrace.WLP_W,
					M4: matchrace.WLP_P,
					M5: matchrace.WLP_W,
				},
				&outobjs.OutMatchRecent{
					Id:       3,
					PlayerId: 3,
					Player: &outobjs.OutMatchPlayer{
						Id:   3,
						Name: "选手名称C",
					},
					M1: matchrace.WLP_L,
					M2: matchrace.WLP_L,
					M3: matchrace.WLP_UNDEFINED,
					M4: matchrace.WLP_W,
					M5: matchrace.WLP_W,
				},
				&outobjs.OutMatchRecent{
					Id:       4,
					PlayerId: 4,
					Player: &outobjs.OutMatchPlayer{
						Id:   4,
						Name: "选手名称D",
					},
					M1: matchrace.WLP_W,
					M2: matchrace.WLP_L,
					M3: matchrace.WLP_P,
					M4: matchrace.WLP_W,
					M5: matchrace.WLP_L,
				},
				&outobjs.OutMatchRecent{
					Id:       5,
					PlayerId: 5,
					Player: &outobjs.OutMatchPlayer{
						Id:   5,
						Name: "选手名称E",
					},
					M1: matchrace.WLP_W,
					M2: matchrace.WLP_L,
					M3: matchrace.WLP_L,
					M4: matchrace.WLP_L,
					M5: matchrace.WLP_W,
				},
				&outobjs.OutMatchRecent{
					Id:       6,
					PlayerId: 6,
					Player: &outobjs.OutMatchPlayer{
						Id:   6,
						Name: "选手名称F",
					},
					M1: matchrace.WLP_L,
					M2: matchrace.WLP_UNDEFINED,
					M3: matchrace.WLP_W,
					M4: matchrace.WLP_L,
					M5: matchrace.WLP_W,
				},
				&outobjs.OutMatchRecent{
					Id:       7,
					PlayerId: 7,
					Player: &outobjs.OutMatchPlayer{
						Id:   7,
						Name: "选手名称G",
					},
					M1: matchrace.WLP_P,
					M2: matchrace.WLP_W,
					M3: matchrace.WLP_W,
					M4: matchrace.WLP_P,
					M5: matchrace.WLP_P,
				},
				&outobjs.OutMatchRecent{
					Id:       8,
					PlayerId: 8,
					Player: &outobjs.OutMatchPlayer{
						Id:   8,
						Name: "选手名称H",
					},
					M1: matchrace.WLP_L,
					M2: matchrace.WLP_W,
					M3: matchrace.WLP_W,
					M4: matchrace.WLP_W,
					M5: matchrace.WLP_W,
				},
				&outobjs.OutMatchRecent{
					Id:       9,
					PlayerId: 9,
					Player: &outobjs.OutMatchPlayer{
						Id:   9,
						Name: "选手名称I",
					},
					M1: matchrace.WLP_P,
					M2: matchrace.WLP_L,
					M3: matchrace.WLP_W,
					M4: matchrace.WLP_P,
					M5: matchrace.WLP_W,
				},
				&outobjs.OutMatchRecent{
					Id:       10,
					PlayerId: 10,
					Player: &outobjs.OutMatchPlayer{
						Id:   10,
						Name: "选手名称J",
					},
					M1: matchrace.WLP_L,
					M2: matchrace.WLP_L,
					M3: matchrace.WLP_L,
					M4: matchrace.WLP_L,
					M5: matchrace.WLP_L,
				},
			},
		},
		&outobjs.OutMatchMode{
			Id:       3,
			MatchId:  1,
			ModeType: matchrace.MODE_TYPE_ELIMIN,
			Title:    "晋级模型测试",
			Elimins: []*outobjs.OutMatchEliminMs{
				&outobjs.OutMatchEliminMs{
					Id:    1,
					Title: "4强淘汰赛",
					T:     matchrace.ELIMIN_MSTYPE_STANDARD,
					Evs: []*outobjs.OutMatchEliminVs{
						&outobjs.OutMatchEliminVs{
							Id:   1,
							VsId: 1,
							MsId: 1,
							Vs: &outobjs.OutMatchVs{
								Id:      10,
								A:       1,
								AName:   "选手A",
								AImg:    0,
								AImgUrl: "",
								AScore:  1,
								AOutlet: true,
								B:       2,
								BName:   "选手B",
								BImg:    0,
								BImgUrl: "",
								BScore:  0,
								BOutlet: false,
							},
						},
						&outobjs.OutMatchEliminVs{
							Id:   1,
							VsId: 1,
							MsId: 1,
							Vs: &outobjs.OutMatchVs{
								Id:      11,
								A:       3,
								AName:   "选手C",
								AImg:    0,
								AImgUrl: "",
								AScore:  0,
								AOutlet: false,
								B:       4,
								BName:   "选手D",
								BImg:    0,
								BImgUrl: "",
								BScore:  2,
								BOutlet: true,
							},
						},
						&outobjs.OutMatchEliminVs{
							Id:   1,
							VsId: 1,
							MsId: 1,
							Vs: &outobjs.OutMatchVs{
								Id:      1,
								A:       5,
								AName:   "选手E",
								AImg:    0,
								AImgUrl: "",
								AScore:  0,
								AOutlet: false,
								B:       6,
								BName:   "选手F",
								BImg:    0,
								BImgUrl: "",
								BScore:  2,
								BOutlet: true,
							},
						},
						&outobjs.OutMatchEliminVs{
							Id:   1,
							VsId: 1,
							MsId: 1,
							Vs: &outobjs.OutMatchVs{
								Id:      1,
								A:       7,
								AName:   "选手G",
								AImg:    0,
								AImgUrl: "",
								AScore:  1,
								AOutlet: true,
								B:       8,
								BName:   "选手H",
								BImg:    0,
								BImgUrl: "",
								BScore:  0,
								BOutlet: false,
							},
						},
					},
				},
				&outobjs.OutMatchEliminMs{
					Id:    2,
					Title: "半决赛",
					T:     matchrace.ELIMIN_MSTYPE_THIRD,
					Evs: []*outobjs.OutMatchEliminVs{
						&outobjs.OutMatchEliminVs{
							Id:   1,
							VsId: 1,
							MsId: 1,
							Vs: &outobjs.OutMatchVs{
								Id:      10,
								A:       1,
								AName:   "选手B",
								AImg:    0,
								AImgUrl: "",
								AScore:  1,
								AOutlet: true,
								B:       2,
								BName:   "选手C",
								BImg:    0,
								BImgUrl: "",
								BScore:  0,
								BOutlet: false,
							},
						},
						&outobjs.OutMatchEliminVs{
							Id:   1,
							VsId: 1,
							MsId: 1,
							Vs: &outobjs.OutMatchVs{
								Id:      11,
								A:       3,
								AName:   "选手D",
								AImg:    0,
								AImgUrl: "",
								AScore:  0,
								AOutlet: false,
								B:       4,
								BName:   "选手E",
								BImg:    0,
								BImgUrl: "",
								BScore:  2,
								BOutlet: true,
							},
						},
					},
				},
				&outobjs.OutMatchEliminMs{
					Id:    3,
					Title: "决赛",
					T:     matchrace.ELIMIN_MSTYPE_CHAMPION,
					Evs: []*outobjs.OutMatchEliminVs{
						&outobjs.OutMatchEliminVs{
							Id:   1,
							VsId: 1,
							MsId: 1,
							Vs: &outobjs.OutMatchVs{
								Id:      12,
								A:       1,
								AName:   "选手B",
								AImg:    0,
								AImgUrl: "",
								AScore:  1,
								AOutlet: true,
								B:       2,
								BName:   "选手E",
								BImg:    0,
								BImgUrl: "",
								BScore:  0,
								BOutlet: false,
							},
						},
					},
				},
			},
		},
	}
	//	for _, r := range rs {
	//		_out := outobjs.GetOutMatchMode(r)
	//		if _out != nil {
	//			outs = append(outs, _out)
	//		}
	//	}
	c.Json(outs)
}
