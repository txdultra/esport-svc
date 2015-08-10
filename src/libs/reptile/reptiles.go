package reptile

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"libs/reptile/douyuapi"
	"logs"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"
	"utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/m3u8"
	"github.com/thrift"
)

//抓取支持项
type REP_SUPPORT string

const (
	REP_SUPPORT_17173  REP_SUPPORT = "17173"
	REP_SUPPORT_PPLIVE REP_SUPPORT = "pplive"
	REP_SUPPORT_TWITCH REP_SUPPORT = "twitch"
	REP_SUPPORT_FYZB   REP_SUPPORT = "fyzb"
	REP_SUPPORT_DOUYU  REP_SUPPORT = "douyu"
	REP_SUPPORT_QQ     REP_SUPPORT = "qq"
	REP_SUPPORT_ZHANQI REP_SUPPORT = "zhanqi"
	REP_SUPPORT_HUOMAO REP_SUPPORT = "huomao"
	REP_SUPPORT_MLGTV  REP_SUPPORT = "mlgtv" //majorleaguegaming
	REP_SUPPORT_HITBOX REP_SUPPORT = "hitbox"
	REP_SUPPORT_QQOPEN REP_SUPPORT = "qqopen"
)

type REP_METHOD string

const (
	REP_METHOD_DIRECT    REP_METHOD = "direct"
	REP_METHOD_SERANDCLT REP_METHOD = "server_client"
)

const (
	//REP_CALLBACK_COMMAND_REPLY     = "reply"
	REP_CALLBACK_COMMAND_EXIT      = "exit"
	REP_CALLBACK_COMMAND_302       = "302"
	REP_CALLBACK_COMMAND_COMPLETED = "completed"
	REP_CALLBACK_COMMAND_POST      = "post"
)

type LIVE_STATUS int

const (
	LIVE_STATUS_LIVING  LIVE_STATUS = 1
	LIVE_STATUS_NOTHING LIVE_STATUS = 2
)

var LIVE_REPTILE_MODULES = map[string]reflect.Type{
	string(REP_SUPPORT_17173):  reflect.TypeOf(R17173Live{}),
	string(REP_SUPPORT_DOUYU):  reflect.TypeOf(DouyuLive{}),
	string(REP_SUPPORT_QQ):     reflect.TypeOf(QQLive{}),
	string(REP_SUPPORT_ZHANQI): reflect.TypeOf(ZhanqiLive{}),
	string(REP_SUPPORT_HUOMAO): reflect.TypeOf(HuomaoLive{}),
	//string(REP_SUPPORT_PPLIVE): reflect.TypeOf(PPTVLive{}),
	string(REP_SUPPORT_TWITCH): reflect.TypeOf(TwitchTVLive{}),
	//string(REP_SUPPORT_FYZB):   reflect.TypeOf(FyzbLive{}),
	string(REP_SUPPORT_MLGTV):  reflect.TypeOf(MLGTVLive{}),
	string(REP_SUPPORT_HITBOX): reflect.TypeOf(HitboxLive{}),
	string(REP_SUPPORT_QQOPEN): reflect.TypeOf(QQOpenLive{}),
}

func SupportReptilePlatforms() []string {
	plats := []string{}
	for k, _ := range LIVE_REPTILE_MODULES {
		plats = append(plats, k)
	}
	return plats
}

func RegisterReptileModule(name string, t reflect.Type) error {
	if _, ok := LIVE_REPTILE_MODULES[name]; ok {
		return errors.New("已存在相同名称的组件")
	}
	LIVE_REPTILE_MODULES[name] = t
	return nil
}

func LiveRepMethod(support REP_SUPPORT) REP_METHOD {
	switch support {
	case REP_SUPPORT_17173, REP_SUPPORT_DOUYU, REP_SUPPORT_QQ, REP_SUPPORT_ZHANQI,
		REP_SUPPORT_HUOMAO, REP_SUPPORT_TWITCH, REP_SUPPORT_MLGTV, REP_SUPPORT_HITBOX, REP_SUPPORT_QQOPEN:
		return REP_METHOD_SERANDCLT
	default:
		return REP_METHOD_DIRECT
	}
}

func Get_REP_SUPPORT(url string) (REP_SUPPORT, error) {
	var lowerUrl = strings.ToLower(url)
	matched, _ := regexp.MatchString("http://v.17173.com/live/(\\w+)/(\\w+)", lowerUrl)
	if matched {
		return REP_SUPPORT_17173, nil
	}
	matched, _ = regexp.MatchString("http://www.douyutv.com/(\\w+)", lowerUrl)
	if matched {
		return REP_SUPPORT_DOUYU, nil
	}
	matched, _ = regexp.MatchString("http://v.qq.com/live/game/(\\d+).html", lowerUrl)
	if matched {
		return REP_SUPPORT_QQ, nil
	}
	matched, _ = regexp.MatchString("http://www.zhanqi.tv/(\\w+)", lowerUrl)
	if matched {
		return REP_SUPPORT_ZHANQI, nil
	}
	matched, _ = regexp.MatchString("http://www.huomaotv.com/live/(\\w+)", lowerUrl)
	if matched {
		return REP_SUPPORT_HUOMAO, nil
	}
	matched, _ = regexp.MatchString("http://www.twitch.tv/(\\w+)", lowerUrl)
	if matched {
		return REP_SUPPORT_TWITCH, nil
	}
	matched, _ = regexp.MatchString("http://tv.majorleaguegaming.com/channel/(\\w+)", lowerUrl)
	if matched {
		return REP_SUPPORT_MLGTV, nil
	}
	matched, _ = regexp.MatchString("http://www.hitbox.tv/(\\w+)", lowerUrl)
	if matched {
		return REP_SUPPORT_HITBOX, nil
	}
	matched, _ = regexp.MatchString("qq_id=(\\d+)", lowerUrl)
	if matched {
		return REP_SUPPORT_QQOPEN, nil
	}
	return "", errors.New("not exist support rep")
}

type ILiveViewOnPc interface {
	ViewHtmlOnPc(url string, width int, height int) string
}

//点播抓取接口
type IReptile interface {
	Reptile(url string) (*VodStreams, error)
	M3u8ToSegs(m3u8url string) ([]VodSeg, error)
}

//点播回调抓取模式接口
type ICallbackReptile interface {
	CallbackReptile(parameter string, url string, cmd string) (vs *VodStreams, clientReqUrl string, nextCmd string, err error)
}

//直播抓取
type IReptileLive interface {
	Reptile(parameter string) (string, error)
}

//通过客户端直播抓取
type IReptileLiveClientProxy interface {
	ProxyReptile(parameter string, cmd string) (clientReqUrl string, nextCmd string, err error)
}

type ILiveStatus interface {
	GetStatus(parameter string) (LIVE_STATUS, error)
}

//空间抓取
type IUserReptile interface {
	Reptile(url string, saveVodFunc func([]*RVodData) bool) error
	ValidateUrl(url string) error
}

var reptile_status_douyu_url, reptile_status_17173_url, reptile_status_fengyun_url, reptile_status_twitchtv_urls, reptile_status_zhanqi_url, reptile_status_huomao_url string

func init() {
	reptile_status_douyu_url = beego.AppConfig.String("reptile.status.douyu.url")
	reptile_status_17173_url = beego.AppConfig.String("reptile.status.17173.url")
	reptile_status_fengyun_url = beego.AppConfig.String("reptile.status.fengyun.url")
	reptile_status_twitchtv_urls = beego.AppConfig.String("reptile.status.twitchtv.urls")
	reptile_status_zhanqi_url = beego.AppConfig.String("reptile.status.zhanqi.url")
	reptile_status_huomao_url = beego.AppConfig.String("reptile.status.huomao.url")
}

////////////////////////////////////////////////////////////////////////////////
//douyu
////////////////////////////////////////////////////////////////////////////////
type DouyuLive struct{}

func (d DouyuLive) ViewHtmlOnPc(url string, width int, height int) string {
	var roomId string
	cache := utils.GetCache()
	err := cache.Get("live_douyu_url_onpc:"+url, &roomId)
	if err != nil {
		str, _err := httplib.Get(url).String()
		if _err != nil {
			return ""
		}
		re := regexp.MustCompile(`"room_id":(\d+)`)
		us := re.FindAllString(str, -1)
		if len(us) == 0 {
			return ""
		}
		mstr := us[0]
		__idx := strings.LastIndex(mstr, ":")
		roomId = mstr[__idx+1:]
		cache.Set("live_douyu_url_onpc:"+url, roomId, 10*time.Hour)
	}
	return fmt.Sprintf(`<embed width="%d" height="%d" type="application/x-shockwave-flash" `+
		`allowfullscreeninteractive="true" allowfullscreen="true" wmode="transparent" bgcolor="javascript:;ffffff"`+
		` quality="high" src="http://staticlive.douyu.tv/common/share/play.swf?room_id=%s" allowscriptaccess="always" `+
		`allownetworking="all">`, width, height, roomId)
}

func (d DouyuLive) GetStatus(parameter string) (LIVE_STATUS, error) {
	__index := strings.LastIndex(parameter, "/")
	if __index <= 0 {
		return LIVE_STATUS_NOTHING, errors.New("抓取地址格式错误")
	}
	live_url := parameter[__index+1:]
	//本地缓存,防止同一时间多次获取数据
	cache := utils.GetLocalCache()
	cacheKey := "reptile_status_douyu_content"
	var content string
	cache.Get(cacheKey, &content)
	if len(content) == 0 {
		for i := 1; i <= 3; i++ {
			req_url := fmt.Sprintf(reptile_status_douyu_url, (i-1)*100, 100)
			req := httplib.Get(req_url)
			req.SetTimeout(3*time.Minute, 3*time.Minute)
			_cnt, err := req.String()
			if err != nil {
				return LIVE_STATUS_NOTHING, err
			}
			content += strings.ToLower(_cnt)
		}
		cache.Set(cacheKey, content, 5*time.Minute)
	}
	if strings.Contains(content, fmt.Sprintf(`href="/%s"`, strings.ToLower(live_url))) {
		return LIVE_STATUS_LIVING, nil
	}
	return LIVE_STATUS_NOTHING, nil
}

//客户端参与的特别处理模式
func (d DouyuLive) ProxyReptile(parameter string, cmd string) (clientReqUrl string, nextCmd string, err error) {
	lcmd := strings.ToLower(cmd)
	//timestamp := time.Now().Unix() * 1000

	if lcmd == "get_stream_params" {
		json_data, err := utils.NewJson([]byte(parameter))
		if err != nil {
			return "", REP_CALLBACK_COMMAND_EXIT, err
		}
		live_stream_url := ""
		live_stream_path := ""
		if data, ok := json_data.CheckGet("data"); ok {
			if _d, _ok := data.CheckGet("rtmp_multi_bitrate"); _ok {
				if __d, __ok := _d.CheckGet("middle"); __ok {
					live_stream_path, _ = __d.String()
				}
			}
			if len(live_stream_path) == 0 {
				if _d, _ok := data.CheckGet("rtmp_live"); _ok {
					live_stream_path, _ = _d.String()
				}
			}
			if _d, _ok := data.CheckGet("rtmp_url"); _ok {
				live_stream_url, _ = _d.String()
			}
		}
		if len(live_stream_url) == 0 || len(live_stream_path) == 0 {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("回传数据错误:005")
		}
		client_url := live_stream_url + "/" + live_stream_path
		return client_url, REP_CALLBACK_COMMAND_302, nil
	}
	if lcmd == REP_CALLBACK_COMMAND_302 {
		return parameter, REP_CALLBACK_COMMAND_COMPLETED, nil
	}

	//第一步交互(默认)
	support, _ := Get_REP_SUPPORT(parameter)
	if support != REP_SUPPORT_DOUYU {
		return "", REP_CALLBACK_COMMAND_EXIT, errors.New("抓取地址格式错误:001")
	}

	var roomId string
	cache := utils.GetCache()
	crr := cache.Get("live_douyu_url:"+parameter, &roomId)
	if crr != nil {
		str, _err := httplib.Get(parameter).String()
		if _err != nil {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("抓取地址格式错误:002")
		}
		re := regexp.MustCompile(`"room_id":(\d+)`)
		us := re.FindAllString(str, -1)
		if len(us) == 0 {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("抓取地址格式错误:003")
		}
		mstr := us[0]
		__idx := strings.LastIndex(mstr, ":")
		roomId = mstr[__idx+1:]
		cache.Set("live_douyu_url:"+parameter, roomId, 10*time.Hour)
	}
	if len(roomId) == 0 {
		return "", REP_CALLBACK_COMMAND_EXIT, errors.New("抓取地址格式错误:004")
	}

	//远程服务调用
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transport, _ := thrift.NewTSocketTimeout(douyu_api_host, 1*time.Minute)
	client := douyuapi.NewDouyuApiServiceClientFactory(transport, protocolFactory)
	if err := transport.Open(); err != nil {
		return "", REP_CALLBACK_COMMAND_EXIT, errors.New("远程服务器错误:005")
	}
	defer transport.Close()
	url, err := client.GetData(roomId, parameter)
	if len(url) == 0 || err != nil {
		return "", REP_CALLBACK_COMMAND_EXIT, errors.New("远程服务器错误:006")
	}
	return url, "get_stream_params", nil
}

func (d DouyuLive) Reptile(parameter string) (string, error) {
	return "", errors.New("此方法无效")
}

////////////////////////////////////////////////////////////////////////////////
//17173
////////////////////////////////////////////////////////////////////////////////
type R17173Live struct{}

func (r R17173Live) ViewHtmlOnPc(url string, width int, height int) string {
	sps := strings.Split(url, "/")
	if len(sps) <= 2 {
		return ""
	}
	return fmt.Sprintf(`<embed width="%d" height="%d" type="application/x-shockwave-flash" `+
		`allowfullscreeninteractive="true" allowfullscreen="always" wmode="transparent" bgcolor="javascript:;ffffff"`+
		` quality="high" src="http://v.17173.com/live/player/Player_stream_customOut.swf?url=http://v.17173.com/live/%s/%s" allowscriptaccess="always" `+
		`allownetworking="all">`, width, height, sps[len(sps)-2], sps[len(sps)-1])
}

func (r R17173Live) GetStatus(parameter string) (LIVE_STATUS, error) {
	__index := strings.LastIndex(parameter, "/")
	if __index <= 0 {
		return LIVE_STATUS_NOTHING, errors.New("抓取地址格式错误")
	}
	live_room_id := parameter[__index+1:]
	//本地缓存,防止同一时间多次获取数据
	cache := utils.GetLocalCache()
	cacheKey := "reptile_status_17173_content"
	var content string
	cache.Get(cacheKey, &content)
	if len(content) == 0 {
		req := httplib.Get(reptile_status_17173_url)
		req.SetTimeout(3*time.Minute, 3*time.Minute)
		content, err := req.String()
		if err != nil {
			return LIVE_STATUS_NOTHING, err
		}
		cache.Set(cacheKey, content, 3*time.Minute)
	}
	if strings.Contains(content, live_room_id) {
		return LIVE_STATUS_LIVING, nil
	}
	return LIVE_STATUS_NOTHING, nil
}

//客户端参与的特别处理模式
func (r R17173Live) ProxyReptile(parameter string, cmd string) (clientReqUrl string, nextCmd string, err error) {
	lcmd := strings.ToLower(cmd)
	timestamp := time.Now().Unix() * 1000
	if lcmd == "get_stream_params" {
		rp_url := "http://gslb.v.17173.com/live?ckey=%s&cdntype=1&optimal=1&t=%d&name=%s&sip=&prot=3&cid=%d&ver=2.0"
		json_data, err := utils.NewJson([]byte(parameter))
		if err != nil {
			return "", REP_CALLBACK_COMMAND_EXIT, err
		}
		if _, ok := json_data.CheckGet("obj"); !ok {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("回传数据错误:005")
		}
		j_data := json_data.Get("obj").Get("liveInfo").Get("live")
		ckey, _ := j_data.Get("ckey").String()
		cId, _ := j_data.Get("cId").Int64()
		streamName, _ := j_data.Get("streamName").String()

		if len(ckey) == 0 { //未抓取到直接退出
			return "", REP_CALLBACK_COMMAND_EXIT, nil
		}
		return fmt.Sprintf(rp_url, ckey, timestamp, streamName, cId), "get_play_stream_url", nil
	}
	if lcmd == "get_play_stream_url" {
		live_json, err := utils.NewJson([]byte(parameter))
		if err != nil {
			return "", REP_CALLBACK_COMMAND_EXIT, err
		}
		live_url, _ := live_json.Get("url").String()
		return live_url, REP_CALLBACK_COMMAND_COMPLETED, nil
	}
	//第一步交互(默认)
	support, _ := Get_REP_SUPPORT(parameter)
	if support != REP_SUPPORT_17173 {
		return "", REP_CALLBACK_COMMAND_EXIT, errors.New("抓取地址格式错误:001")
	}
	live_info_url_fmt := "http://v.17173.com/live/l_jsonData.action?liveRoomId=%s&_=%d"
	__index := strings.LastIndex(parameter, "/")
	if __index <= 0 {
		return "", REP_CALLBACK_COMMAND_EXIT, errors.New("抓取地址格式错误:002")
	}
	live_room_id := parameter[__index+1:]
	return fmt.Sprintf(live_info_url_fmt, live_room_id, timestamp), "get_stream_params", nil
}

func (r R17173Live) Reptile(parameter string) (string, error) {
	rp_url := "http://gslb.v.17173.com/live?ckey=%s&cdntype=1&optimal=1&t=%d&name=%s&sip=&prot=3&cid=%d&ver=2.0"
	live_info_url_fmt := "http://v.17173.com/live/l_jsonData.action?liveRoomId=%s&_=%d"
	__index := strings.LastIndex(parameter, "/")
	if __index <= 0 {
		return "", errors.New("抓取地址格式错误")
	}
	live_room_id := parameter[__index+1:]
	timestamp := time.Now().Unix() * 1000
	req := httplib.Get(fmt.Sprintf(live_info_url_fmt, live_room_id, timestamp))

	live_info, err := req.Bytes()
	if err == nil {
		json_data, err := utils.NewJson(live_info)
		if err != nil {
			return "live/l_jsonData.action地址返回错误信息", err
		}
		j_data := json_data.Get("obj").Get("liveInfo").Get("live")
		ckey, _ := j_data.Get("ckey").String()
		cId, _ := j_data.Get("cId").Int64()
		streamName, _ := j_data.Get("streamName").String()
		gs_live_url := fmt.Sprintf(rp_url, ckey, timestamp, streamName, cId)
		live_info, err := httplib.Get(gs_live_url).Bytes()
		if err == nil {
			live_json, _ := utils.NewJson(live_info)
			live_url, _ := live_json.Get("url").String()
			return live_url, nil
		}
		return "", err
	}
	return "", err
}

////////////////////////////////////////////////////////////////////////////////
//QQ
////////////////////////////////////////////////////////////////////////////////
type QQLive struct{}

func (q QQLive) ViewHtmlOnPc(url string, width int, height int) string {
	return ""
}

func (q QQLive) GetStatus(parameter string) (LIVE_STATUS, error) {
	return LIVE_STATUS_NOTHING, nil
}

func (q QQLive) Reptile(parameter string) (string, error) {
	return "", errors.New("此方法无效")
}

//客户端参与的特别处理模式
func (r QQLive) ProxyReptile(parameter string, cmd string) (clientReqUrl string, nextCmd string, err error) {
	lcmd := strings.ToLower(cmd)
	if lcmd == "get_stream_params" {
		json_data, err := utils.NewJson([]byte(parameter))
		if err != nil {
			return "", REP_CALLBACK_COMMAND_EXIT, err
		}
		if _, ok := json_data.CheckGet("playurl"); !ok {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("回传数据错误:005")
		}
		stream_url, _ := json_data.Get("playurl").String()
		if len(stream_url) == 0 { //未抓取到直接退出
			return "", REP_CALLBACK_COMMAND_EXIT, nil
		}
		return stream_url, REP_CALLBACK_COMMAND_COMPLETED, nil
	}

	//第一步交互(默认)
	support, _ := Get_REP_SUPPORT(parameter)
	if support != REP_SUPPORT_QQ {
		return "", REP_CALLBACK_COMMAND_EXIT, errors.New("抓取地址格式错误:001")
	}
	cache := utils.GetLocalCache()
	var content string
	err = cache.Get(fmt.Sprintf("reptile_qqlive_page:%s", parameter), &content)
	if err != nil {
		req := httplib.Get(parameter)
		req.SetTimeout(3*time.Minute, 3*time.Minute)
		content, err = req.String()
		if err != nil || len(content) == 0 {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("抓取地址格式错误:002")
		}
		cache.Set(fmt.Sprintf("reptile_qqlive_page:%s", parameter), content, 12*time.Hour)
	}
	//playid:'100601300'
	re := regexp.MustCompile(`playid:'(\d+)'`)
	us := re.FindAllString(content, -1)
	if len(us) == 0 {
		return "", REP_CALLBACK_COMMAND_EXIT, errors.New("抓取地址格式错误:003")
	}
	mstr := us[0]
	mstr = strings.Replace(mstr, "'", "", -1)
	__idx := strings.LastIndex(mstr, ":")
	cnlid := mstr[__idx+1:]
	client_url := "http://info.zb.qq.com/?cnlid=" + cnlid + "&host=qq.com&cmd=2&qq=0&guid=" + utils.RandomStrings(32) +
		"&txvjsv=2.0&stream=1&debug=&ip=&system=1&sdtfrom=113"
	return client_url, "get_stream_params", nil
}

////////////////////////////////////////////////////////////////////////////////
//ZHANQI
////////////////////////////////////////////////////////////////////////////////
type ZhanqiLive struct{}

func (q ZhanqiLive) ViewHtmlOnPc(url string, width int, height int) string {
	cache := utils.GetCache()
	var roomId string
	err := cache.Get(fmt.Sprintf("live_zhanqi_url_onpc:%s", url), &roomId)
	if err != nil {
		req := httplib.Get(url)
		req.SetTimeout(3*time.Minute, 3*time.Minute)
		content, err := req.String()
		if err != nil || len(content) == 0 {
			return ""
		}
		re := regexp.MustCompile(`oRoom = {\"id\":\"(\d+)\"`)
		us := re.FindAllString(content, -1)
		if len(us) == 0 {
			return ""
		}
		mstr := us[0]
		mstr = strings.Replace(mstr, "\"", "", -1)
		__idx := strings.LastIndex(mstr, ":")
		roomId = mstr[__idx+1:]
		cache.Set(fmt.Sprintf("live_zhanqi_url_onpc:%s", url), roomId, 12*time.Minute)
	}
	return fmt.Sprintf(`<iframe width="%d" height="%d" frameborder="0" scrolling="no" src="http://www.zhanqi.tv/live/embed?roomId=%s&fhost=other"></iframe>`, width, height, roomId)
}

func (q ZhanqiLive) GetStatus(parameter string) (LIVE_STATUS, error) {
	__index := strings.LastIndex(parameter, "/")
	if __index <= 0 {
		return LIVE_STATUS_NOTHING, errors.New("抓取地址格式错误")
	}
	live_room_id := parameter[__index+1:]
	//本地缓存,防止同一时间多次获取数据
	cache := utils.GetLocalCache()
	cacheKey := "reptile_status_zhanqi_content"
	var content string
	cache.Get(cacheKey, &content)
	if len(content) == 0 {
		req := httplib.Get(reptile_status_zhanqi_url)
		req.SetTimeout(3*time.Minute, 3*time.Minute)
		content, err := req.String()
		if err != nil {
			return LIVE_STATUS_NOTHING, err
		}
		cache.Set(cacheKey, content, 3*time.Minute)
	}
	//"code":"14723","domain":"jiezou"
	type RoomObj struct {
		Code   string `json:"code"`
		Domain string `json:"domain"`
		Status string `json:"status"`
	}
	type DataObj struct {
		Cnt   int        `json:"cnt"`
		Rooms []*RoomObj `json:"rooms"`
	}
	type LiveObj struct {
		Code    int      `json:"code"`
		Data    *DataObj `json:"data"`
		Message string   `json:"message"`
	}
	lobj := &LiveObj{}
	err := json.Unmarshal([]byte(content), &lobj)
	if err == nil && lobj.Data != nil && lobj.Data.Rooms != nil {
		for _, room := range lobj.Data.Rooms {
			if room.Code == live_room_id || room.Domain == live_room_id {
				if room.Status == "4" {
					return LIVE_STATUS_LIVING, nil
				} else {
					return LIVE_STATUS_NOTHING, nil
				}
			}
		}
	}

	//code_pttn := fmt.Sprintf(`"code":"%s"`, live_room_id)
	//domain_pttn := fmt.Sprintf(`"domain":"%s"`, live_room_id)
	//if strings.Contains(content, code_pttn) {
	//	return LIVE_STATUS_LIVING, nil
	//}
	//if strings.Contains(content, domain_pttn) {
	//	return LIVE_STATUS_LIVING, nil
	//}
	return LIVE_STATUS_NOTHING, nil
}

func (q ZhanqiLive) Reptile(parameter string) (string, error) {
	return "", errors.New("此方法无效")
}

//客户端参与的特别处理模式
func (r ZhanqiLive) ProxyReptile(parameter string, cmd string) (clientReqUrl string, nextCmd string, err error) {
	//第一步交互(默认)
	support, _ := Get_REP_SUPPORT(parameter)
	if support != REP_SUPPORT_ZHANQI {
		return "", REP_CALLBACK_COMMAND_EXIT, errors.New("抓取地址格式错误:001")
	}

	cache := utils.GetLocalCache()
	var client_url string
	err = cache.Get(fmt.Sprintf("reptile_zhanqi_stream700_url:%s", parameter), &client_url)
	if err != nil {
		req := httplib.Get(parameter)
		req.SetTimeout(3*time.Minute, 3*time.Minute)
		content, err := req.String()
		if err != nil || len(content) == 0 {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("抓取地址格式错误:003")
		}
		re := regexp.MustCompile(`"videoIdKey":"([a-z0-9A-Z_/=\\\-]+)"`)
		us := re.FindAllString(content, -1)
		if len(us) == 0 {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("抓取地址格式错误:004")
		}
		mstr := us[0]
		mstr = strings.Replace(mstr, "\"", "", -1)
		__idx := strings.LastIndex(mstr, ":")
		cid := mstr[__idx+1:]

		client_url = fmt.Sprintf("http://dlhls.cdn.zhanqi.tv/zqlive/%s_700/index.m3u8", cid)
		req1 := httplib.Get(client_url)
		rep1, _ := req1.Response()
		defer rep1.Body.Close()
		if rep1.StatusCode == 200 {
			cache.Set(fmt.Sprintf("reptile_zhanqi_stream_url:%s", parameter), client_url, 2*time.Hour)
			return client_url, REP_CALLBACK_COMMAND_COMPLETED, nil
		}
		client_url = fmt.Sprintf("http://dlhls.cdn.zhanqi.tv/zqlive/%s.m3u8", cid)
		req2 := httplib.Get(client_url)
		rep2, _ := req2.Response()
		defer rep2.Body.Close()
		if rep2.StatusCode == 200 {
			cache.Set(fmt.Sprintf("reptile_zhanqi_stream_url:%s", parameter), client_url, 2*time.Hour)
			return client_url, REP_CALLBACK_COMMAND_COMPLETED, nil
		}
		client_url = fmt.Sprintf("rtmp://dlrtmp.cdn.zhanqi.tv/zqlive&VideoID=%s", cid)
		cache.Set(fmt.Sprintf("reptile_zhanqi_stream700_url:%s", parameter), client_url, 2*time.Hour)
	}
	return client_url, REP_CALLBACK_COMMAND_COMPLETED, nil
}

////////////////////////////////////////////////////////////////////////////////
//HUOMAO
////////////////////////////////////////////////////////////////////////////////
type HuomaoLive struct{}

func (q HuomaoLive) ViewHtmlOnPc(url string, width int, height int) string {
	__index := strings.LastIndex(url, "/")
	if __index <= 0 {
		return ""
	}
	return fmt.Sprintf(`<iframe width="%d" height="%d" frameborder="0" scrolling="no" src="http://www.huomaotv.com/index.php?c=outplayer&live_id=%d"></iframe>`, width, height, url[__index+1])
}

func (q HuomaoLive) GetStatus(parameter string) (LIVE_STATUS, error) {
	__index := strings.LastIndex(parameter, "/")
	if __index <= 0 {
		return LIVE_STATUS_NOTHING, errors.New("抓取地址格式错误")
	}
	live_id := parameter[__index+1:]
	//本地缓存,防止同一时间多次获取数据
	cache := utils.GetLocalCache()
	cacheKey := "reptile_status_houmao_content"

	htmls := []string{}
	cache.Get(cacheKey, &htmls)
	if len(htmls) == 0 {
		for i := 1; i <= 3; i++ {
			url := fmt.Sprintf(reptile_status_huomao_url, i)
			req := httplib.Get(url)
			req.SetTimeout(3*time.Minute, 3*time.Minute)
			content, err := req.String()
			if err != nil || len(content) == 0 {
				continue
			}
			htmls = append(htmls, content)
		}
		cache.Set(cacheKey, htmls, 5*time.Minute)
	}

	status := LIVE_STATUS_NOTHING
	for _, html := range htmls {
		buf := bytes.NewBuffer(nil)
		buf.WriteString(html)
		doc, err := goquery.NewDocumentFromReader(buf)
		if err == nil {
			doc.Find(".VOD_pic").Each(func(i int, s *goquery.Selection) {
				href, _ := s.Find(".play_btn").Attr("href")
				if href == fmt.Sprintf(`/live/%s`, live_id) {
					living := s.Find(".up_offline").Length()
					if living == 0 {
						status = LIVE_STATUS_LIVING
						return
					}
				}
			})
		}
	}
	return status, nil
}

func (q HuomaoLive) Reptile(parameter string) (string, error) {
	return "", errors.New("此方法无效")
}

//客户端参与的特别处理模式
func (r HuomaoLive) ProxyReptile(parameter string, cmd string) (clientReqUrl string, nextCmd string, err error) {
	lcmd := strings.ToLower(cmd)
	if lcmd == REP_CALLBACK_COMMAND_302 {
		return parameter, REP_CALLBACK_COMMAND_COMPLETED, nil
	}
	//第一步交互(默认)
	support, _ := Get_REP_SUPPORT(parameter)
	if support != REP_SUPPORT_HUOMAO {
		return "", REP_CALLBACK_COMMAND_EXIT, errors.New("抓取地址格式错误:001")
	}

	cache := utils.GetLocalCache()
	var client_url string
	err = cache.Get(fmt.Sprintf("reptile_huomao_url:%s", parameter), &client_url)
	if err != nil {
		req := httplib.Get(parameter)
		req.SetTimeout(3*time.Minute, 3*time.Minute)
		content, err := req.String()
		if err != nil || len(content) == 0 {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("抓取地址格式错误:003")
		}
		re := regexp.MustCompile(`video_name = '([a-z0-9A-Z_/=\\\-]+)'`)
		us := re.FindAllString(content, -1)
		if len(us) == 0 {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("抓取地址格式错误:004")
		}
		mstr := us[0]
		mstr = strings.Replace(mstr, "'", "", -1)
		mstr = strings.Replace(mstr, " ", "", -1)
		__idx := strings.LastIndex(mstr, "=")
		vids := mstr[__idx+1:]

		req = httplib.Post("http://www.huomaotv.com/swf/live_data")
		req.SetTimeout(3*time.Minute, 3*time.Minute)
		req.Param("streamtype", "live")
		req.Param("VideoIDS", vids)
		bytes, err := req.Bytes()
		if err != nil || len(bytes) == 0 {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("抓取地址格式错误:005")
		}
		//分析数据
		json_data, err := utils.NewJson(bytes)
		if err != nil {
			return "", REP_CALLBACK_COMMAND_EXIT, err
		}
		if _, ok := json_data.CheckGet("streamList"); !ok {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("回传数据错误:005")
		}
		arr, _ := json_data.Get("streamList").Array()
		if len(arr) == 0 { //未抓取到直接退出
			return "", REP_CALLBACK_COMMAND_EXIT, nil
		}
		stream_url := ""
		searchStream := func(data []interface{}, t string) (bool, string) {
			for _, ar := range data {
				stream := ar.(map[string]interface{})
				t, ok := stream["type"]
				if ok {
					if t.(string) == t {
						return true, stream["url"].(string)
					}
				}
			}
			return false, ""
		}

		ok, stream_url := searchStream(arr, "HD")
		if ok {
			cache.Set(fmt.Sprintf("reptile_huomao_url:%s", parameter), stream_url, 2*time.Hour)
			return stream_url, REP_CALLBACK_COMMAND_COMPLETED, nil
		}
		ok, stream_url = searchStream(arr, "SD")
		if ok {
			cache.Set(fmt.Sprintf("reptile_huomao_url:%s", parameter), stream_url, 2*time.Hour)
			return stream_url, REP_CALLBACK_COMMAND_COMPLETED, nil
		}
		return "", REP_CALLBACK_COMMAND_EXIT, nil
	}
	return client_url, REP_CALLBACK_COMMAND_COMPLETED, nil
}

////////////////////////////////////////////////////////////////////////////////
//PPTV
////////////////////////////////////////////////////////////////////////////////
type PPTVLive struct{}

func (r PPTVLive) GetStatus(parameter string) (LIVE_STATUS, error) {
	return LIVE_STATUS_NOTHING, errors.New("未实现此功能")
}

func (r PPTVLive) Reptile(parameter string) (string, error) {

	return "http://web-play.pptv.com/web-m3u8-" + parameter + ".m3u8?type=m3u8.web.phone&playback=0&kk=", nil

	//ifurl := "http://web-play.pptv.com/webplay3-0-" + url + ".xml?version=4&type=live&kk=" + utils.MakeRndStrs(10)
	//req := httplib.Get(ifurl)
	//req.Header("Referer", "http://pub.pptv.com/player/iframe/index.html?#w=800&h=622&id="+url)
	//req.Header("User-Agent", "Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_2 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Mobile/8C134")
	//rspstr, _ := req.String()
	//xml, _ := utils.LoadByXml(rspstr)
	////fmt.Println(err)
	////fmt.Println(xml.Node("channel").Node("stream").Node("item"))
	//rid, ok := xml.Node("channel").Node("stream").Node("item").AttrValue("rid")
	//ip := xml.Node("dt").Node("sh").Value
	//if ok && len(ip) > 0 {
	//	return fmt.Sprintf("http://%s/%s.m3u8?type=m3u8.web.phone", ip, rid), nil
	//}
	//return "", errors.New("解析失败")
}

////////////////////////////////////////////////////////////////////////////////
//Twitch
////////////////////////////////////////////////////////////////////////////////
type TwitchTVLive struct{}

func (r TwitchTVLive) ViewHtmlOnPc(url string, width int, height int) string {
	__index := strings.LastIndex(url, "/")
	if __index <= 0 {
		return ""
	}
	return fmt.Sprintf(`<iframe src="http://www.twitch.tv/%s/embed" frameborder="0" scrolling="no" height="%d" width="%d"></iframe>`,
		url[__index+1:], height, width)
}

func (r TwitchTVLive) ProxyReptile(parameter string, cmd string) (clientReqUrl string, nextCmd string, err error) {
	lcmd := strings.ToLower(cmd)
	if lcmd == "get_stream_params" {
		jdata, err := utils.NewJson([]byte(parameter))
		if err != nil {
			return "", REP_CALLBACK_COMMAND_EXIT, err
		}
		if _, ok := jdata.CheckGet("token"); !ok {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("数据错误:001")
		}
		token, _ := jdata.Get("token").String()

		tjson, err := utils.NewJson([]byte(token))
		if err != nil {
			return "", REP_CALLBACK_COMMAND_EXIT, err
		}
		room_id, _ := tjson.Get("channel").String()

		sig, _ := jdata.Get("sig").String()
		stream_url := fmt.Sprintf("http://usher.justin.tv/api/channel/hls/%s.m3u8?sig=%s&allow_source=true&player=twitchweb&token=%s", room_id, sig, utils.UrlEncode(token))
		return stream_url, "download_file_content", nil
	}
	if lcmd == "download_file_content" {
		buf := bytes.NewBufferString(parameter)
		p, listType, err := m3u8.Decode(*buf, true)
		if err != nil {
			logs.Errorf("twitch m3u8 decode fail:%s", err.Error())
			fmt.Println(err.Error())
			return "", REP_CALLBACK_COMMAND_EXIT, nil
		}
		vpuris := make(map[string]string)
		switch listType {
		case m3u8.MASTER:
			masterpl := p.(*m3u8.MasterPlaylist)
			for _, v := range masterpl.Variants {
				vpuris[v.VariantParams.Video] = v.URI
			}
			if uri, ok := vpuris["medium"]; ok {
				return uri, REP_CALLBACK_COMMAND_COMPLETED, nil
			}
			if uri, ok := vpuris["low"]; ok {
				return uri, REP_CALLBACK_COMMAND_COMPLETED, nil
			}
			if uri, ok := vpuris["high"]; ok {
				return uri, REP_CALLBACK_COMMAND_COMPLETED, nil
			}
			if uri, ok := vpuris["chunked"]; ok {
				return uri, REP_CALLBACK_COMMAND_COMPLETED, nil
			}
			return "", REP_CALLBACK_COMMAND_EXIT, nil
		default:
			return "", REP_CALLBACK_COMMAND_EXIT, nil
		}
	}
	__index := strings.LastIndex(parameter, "/")
	if __index <= 0 {
		return "", REP_CALLBACK_COMMAND_EXIT, errors.New("抓取地址格式错误")
	}
	room_id := parameter[__index+1:]
	url := fmt.Sprintf("http://api.twitch.tv/api/channels/%s/access_token?on_site=1", room_id)
	return url, "get_stream_params", nil
}

func (r TwitchTVLive) GetStatus(parameter string) (LIVE_STATUS, error) {
	_idx := strings.LastIndex(parameter, "/")
	if _idx <= 0 {
		return LIVE_STATUS_NOTHING, errors.New("抓取地址格式错误")
	}
	who := parameter[_idx+1:]
	//本地缓存,防止同一时间多次获取数据
	cache := utils.GetLocalCache()
	cacheKey := "reptile_status_twtichtv_content"
	var content string
	cache.Get(cacheKey, &content)
	if len(content) == 0 {
		urls := strings.Split(reptile_status_twitchtv_urls, ";")
		for _, url := range urls {
			if len(strings.TrimSpace(url)) == 0 {
				continue
			}
			_c, _ := httplib.Get(url).String()
			content += _c
		}
		cache.Set(cacheKey, content, 1*time.Minute)
	}
	if strings.Contains(content, who) {
		return LIVE_STATUS_LIVING, nil
	}
	return LIVE_STATUS_NOTHING, nil
}

func (r TwitchTVLive) Reptile(parameter string) (string, error) {
	return "", errors.New("此方法无效")
}

////////////////////////////////////////////////////////////////////////////////
//FYZB
////////////////////////////////////////////////////////////////////////////////
type FyzbLive struct{}

func (fy FyzbLive) GetStatus(parameter string) (LIVE_STATUS, error) {
	_idx := strings.LastIndex(parameter, "/")
	who := parameter[_idx+1:]
	//本地缓存,防止同一时间多次获取数据
	cache := utils.GetLocalCache()
	cacheKey := "reptile_status_fengyun_content"
	var content string
	cache.Get(cacheKey, &content)
	if len(content) == 0 {
		content, err := httplib.Get(reptile_status_fengyun_url).String()
		if err != nil {
			return LIVE_STATUS_NOTHING, err
		}
		cache.Set(cacheKey, content, 1*time.Minute)
	}

	buf := bytes.NewBuffer(nil)
	buf.WriteString(content)
	doc, err := goquery.NewDocumentFromReader(buf)
	if err != nil {
		return LIVE_STATUS_NOTHING, err
	}
	status := LIVE_STATUS_NOTHING
	doc.Find(".channel-link").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		clss, _ := s.Attr("class")
		if strings.Contains(href, who) && !strings.Contains(clss, "bad") {
			status = LIVE_STATUS_LIVING
		}
	})
	return status, nil
}

func (fy FyzbLive) Reptile(parameter string) (string, error) {
	m_url := strings.Replace(parameter, "www", "m", 1)
	req := httplib.Get(m_url)
	req.Header("User-Agent", "Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_2 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Mobile/8C134")
	rep, _ := req.Response()
	doc, err := goquery.NewDocumentFromResponse(rep)
	if err == nil {
		if val, ok := doc.Find("video").Attr("src"); ok {
			return val, nil
		}
	}
	return "", errors.New("解析失败")
}

////////////////////////////////////////////////////////////////////////////////
//majorleaguegaming
////////////////////////////////////////////////////////////////////////////////
type MLGTVLive struct{}

func (r MLGTVLive) ViewHtmlOnPc(url string, width int, height int) string {
	return ""
}

func (r MLGTVLive) ProxyReptile(parameter string, cmd string) (clientReqUrl string, nextCmd string, err error) {
	lcmd := strings.ToLower(cmd)
	if lcmd == "download_file_content" {
		buf := bytes.NewBufferString(parameter)
		p, listType, err := m3u8.Decode(*buf, true)
		if err != nil {
			logs.Errorf("majorleaguegaming m3u8 decode fail:%s", err.Error())
			fmt.Println(err.Error())
			return "", REP_CALLBACK_COMMAND_EXIT, nil
		}
		vpuris := make(map[string]string)
		switch listType {
		case m3u8.MASTER:
			masterpl := p.(*m3u8.MasterPlaylist)
			for _, v := range masterpl.Variants {
				vpuris[v.VariantParams.Resolution] = v.URI
			}
			fmt.Println(vpuris)
			if uri, ok := vpuris["854x480"]; ok {
				return uri, REP_CALLBACK_COMMAND_COMPLETED, nil
			}
			if uri, ok := vpuris["640x360"]; ok {
				return uri, REP_CALLBACK_COMMAND_COMPLETED, nil
			}
			if uri, ok := vpuris["1280x720"]; ok {
				return uri, REP_CALLBACK_COMMAND_COMPLETED, nil
			}
			return "", REP_CALLBACK_COMMAND_EXIT, nil
		default:
			return "", REP_CALLBACK_COMMAND_EXIT, nil
		}
	}

	support, _ := Get_REP_SUPPORT(parameter)
	if support != REP_SUPPORT_MLGTV {
		return "", REP_CALLBACK_COMMAND_EXIT, errors.New("抓取地址格式错误:00")
	}

	cache := utils.GetLocalCache()
	var m3u8_url string
	ckey := "live_mlgtv_m3u8url:" + parameter
	crr := cache.Get(ckey, &m3u8_url)
	if crr != nil {
		//第一步
		doc, _ := goquery.NewDocument(parameter)
		sel := doc.Find("meta").FilterFunction(func(i int, s *goquery.Selection) bool {
			val, _ := s.Attr("name")
			if val == "twitter:image" {
				return true
			}
			return false
		})
		if sel == nil {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("抓取地址格式错误01")
		}
		metaUrl, existed := sel.Attr("content")
		if !existed {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("抓取地址格式错误02")
		}
		_ss := strings.Split(metaUrl, "/")
		if len(_ss) < 7 {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("抓取地址格式错误03")
		}
		channle_id := _ss[6]
		stream_all_data, err := httplib.Get("http://streamapi.majorleaguegaming.com/service/streams/all").String()
		if err != nil {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("连接失败01")
		}
		_stream_json, err := utils.NewJson([]byte(stream_all_data))
		if err != nil {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("解析json失败00")
		}
		_j1, _ok := _stream_json.CheckGet("data")
		if !_ok {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("解析json失败01")
		}
		_j2, _ok := _j1.CheckGet("items")
		if !_ok {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("解析json失败02")
		}
		channel_streams, err := _j2.Array()
		if err != nil {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("解析json失败03")
		}
		stream_name := ""
		for _, cs := range channel_streams {
			csm := cs.(map[string]interface{})
			if csm == nil {
				return "", REP_CALLBACK_COMMAND_EXIT, errors.New("解析json失败04")
			}
			cid := csm["channel_id"].(json.Number).String()
			if cid == channle_id {
				stream_name = csm["stream_name"].(string)
				break
			}
		}
		if len(stream_name) == 0 {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("解析json失败05")
		}
		//m3u8 url
		url := fmt.Sprintf("http://streamapi.majorleaguegaming.com/service/streams/playback/%s?format=all", stream_name)
		m3u8_data, err := httplib.Get(url).String()
		if err != nil {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("抓取地址格式错误04")
		}
		_m3u8_json, err := utils.NewJson([]byte(m3u8_data))
		if err != nil {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("解析json失败06")
		}
		_m3u8_j1, _ok := _m3u8_json.CheckGet("data")
		if !_ok {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("解析json失败07")
		}
		_m3u8_j2, _ok := _m3u8_j1.CheckGet("items")
		if !_ok {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("解析json失败08")
		}
		stream_urls, err := _m3u8_j2.Array()
		if err != nil {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("解析json失败09")
		}
		for _, su := range stream_urls {
			csm := su.(map[string]interface{})
			if csm == nil {
				return "", REP_CALLBACK_COMMAND_EXIT, errors.New("解析json失败10")
			}
			cid := csm["channel_id"].(json.Number).String()
			fm := csm["format"].(string)
			if cid == channle_id && fm == "hls" {
				m3u8_url = csm["url"].(string)
				break
			}
		}
		cache.Set(ckey, m3u8_url, 1*time.Hour)
	}
	if len(m3u8_url) == 0 {
		return "", REP_CALLBACK_COMMAND_EXIT, errors.New("解析json失败11")
	}
	return m3u8_url, "download_file_content", nil
}

func (r MLGTVLive) GetStatus(parameter string) (LIVE_STATUS, error) {
	return LIVE_STATUS_NOTHING, nil
}

func (r MLGTVLive) Reptile(parameter string) (string, error) {
	return "", errors.New("此方法无效")
}

////////////////////////////////////////////////////////////////////////////////
//hitbox
////////////////////////////////////////////////////////////////////////////////
type HitboxLive struct{}

func (r HitboxLive) ViewHtmlOnPc(url string, width int, height int) string {
	return ""
}

func (r HitboxLive) ProxyReptile(parameter string, cmd string) (clientReqUrl string, nextCmd string, err error) {
	//	http://edge.hls.dt.hitbox.tv/hls/ectv_360p/index.m3u8
	//	http://edge.hls.dt.hitbox.tv/hls/ectv_480p/index.m3u8
	//	http://edge.hls.dt.hitbox.tv/hls/ectv/index.m3u8
	support, _ := Get_REP_SUPPORT(parameter)
	if support != REP_SUPPORT_HITBOX {
		return "", REP_CALLBACK_COMMAND_EXIT, errors.New("抓取地址格式错误:00")
	}

	cache := utils.GetLocalCache()
	ckey := fmt.Sprintf("reptile_hitbox_stream_url:%s", parameter)
	var m3u8url string
	err = cache.Get(ckey, &m3u8url)
	if err != nil {
		__index := strings.LastIndex(parameter, "/")
		if __index <= 0 {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("抓取地址格式错误")
		}
		room_name := parameter[__index+1:]
		room_name = strings.Trim(room_name, " ")

		qxs := []string{"_480p", "", "_360p"}
		for _, qx := range qxs {
			m3u8url = fmt.Sprintf("http://edge.hls.dt.hitbox.tv/hls/%s%s/index.m3u8", room_name, qx)
			req := httplib.Get(m3u8url)
			rep, _ := req.Response()
			if rep.StatusCode == 200 {
				cache.Set(ckey, m3u8url, 1*time.Hour)
				rep.Body.Close()
				return m3u8url, REP_CALLBACK_COMMAND_COMPLETED, nil
			}
			rep.Body.Close()
		}
	}
	return m3u8url, REP_CALLBACK_COMMAND_COMPLETED, nil
}

func (r HitboxLive) GetStatus(parameter string) (LIVE_STATUS, error) {
	return LIVE_STATUS_NOTHING, nil
}

func (r HitboxLive) Reptile(parameter string) (string, error) {
	return "", errors.New("此方法无效")
}

////////////////////////////////////////////////////////////////////////////////
//qqopen
////////////////////////////////////////////////////////////////////////////////
type QQOpenLive struct{}

func (r QQOpenLive) ViewHtmlOnPc(url string, width int, height int) string {
	return ""
}

func (r QQOpenLive) ProxyReptile(parameter string, cmd string) (clientReqUrl string, nextCmd string, err error) {
	support, _ := Get_REP_SUPPORT(parameter)
	if support != REP_SUPPORT_QQOPEN {
		return "", REP_CALLBACK_COMMAND_EXIT, errors.New("抓取地址格式错误:00")
	}

	v, err := url.ParseRequestURI(parameter)
	if err != nil {
		return "", REP_CALLBACK_COMMAND_EXIT, errors.New("抓取地址格式错误:01")
	}
	u := v.RawQuery
	vals, _ := url.ParseQuery(u)
	roomId := vals.Get("qq_id")
	if len(roomId) == 0 {
		return "", REP_CALLBACK_COMMAND_EXIT, errors.New("抓取地址格式错误:02")
	}
	roomId = fmt.Sprintf("qqopen_%s", roomId)

	cache := utils.GetLocalCache()
	ckey := fmt.Sprintf("reptile_qqopen_stream_url:%s", parameter)
	var playurl string
	err = cache.Get(ckey, &playurl)
	if err != nil {
		//远程服务调用
		protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
		transport, _ := thrift.NewTSocketTimeout(douyu_api_host, 1*time.Minute)
		client := douyuapi.NewDouyuApiServiceClientFactory(transport, protocolFactory)
		if err := transport.Open(); err != nil {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("远程服务器错误:005")
		}
		defer transport.Close()
		url, err := client.GetData(roomId, parameter)
		if len(url) == 0 || err != nil {
			return "", REP_CALLBACK_COMMAND_EXIT, errors.New("远程服务器错误:006")
		}
		return url, REP_CALLBACK_COMMAND_COMPLETED, nil
	}
	return playurl, REP_CALLBACK_COMMAND_COMPLETED, nil
}

func (r QQOpenLive) GetStatus(parameter string) (LIVE_STATUS, error) {
	return LIVE_STATUS_NOTHING, nil
}

func (r QQOpenLive) Reptile(parameter string) (string, error) {
	return "", errors.New("此方法无效")
}
