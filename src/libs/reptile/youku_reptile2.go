package reptile

import (
	"bytes"
	enjson "encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"utils"

	"github.com/m3u8"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
)

const (
	YOUKU_M3U8_FMT   = "http://pl.youku.com/playlist/m3u8?ctype=12&ep=%s&ev=1&keyframe=1&oip=%d&sid=%s&token=%s&type=%s&vid=%s"
	VIDEOIDS_URL_FMT = "http://v.youku.com/player/getPlayList/VideoIDS/%s/Pf/4/ctype/12/ev/1"
)

var last_reptile_time_v2 time.Time = time.Now()
var locker_v2 *sync.Mutex = new(sync.Mutex)

type YoukuReptileV2 struct {
	//charCode map[string]string
}

func NewYoukuReptileV2() *YoukuReptileV2 {
	yr := &YoukuReptileV2{}
	//yr.charCode = make(map[string]string)
	return yr
}

func (y *YoukuReptileV2) getVid(url string) (string, error) {
	re := regexp.MustCompile("http://v.youku.com/v_show/id_([0-9A-Za-z_=]+)\\.html")
	matchs := re.FindSubmatch([]byte(url))
	if len(matchs) != 2 {
		return "", errors.New("youku url format fail")
	}
	yk_id := string(matchs[1])
	__idx := strings.Index(yk_id, "_")
	if __idx > 0 {
		yk_id = yk_id[0:__idx]
	}
	return yk_id, nil
}

func (y *YoukuReptileV2) encoder(a string, b []byte, base64 bool) string {
	result := ""
	byteR := []byte{}
	f := 0
	bs := make([]int, 256, 256)
	for i := 0; i < 256; i++ {
		bs[i] = i
	}
	for i := 0; i < 256; i++ {
		f = (f + bs[i] + int(a[i%len(a)])) % 256
		temp := bs[i]
		bs[i] = bs[f]
		bs[f] = temp
	}
	f, h := 0, 0
	for i := 0; i < len(b); i++ {
		h = (h + 1) % 256
		f = (f + bs[h]) % 256
		temp := bs[h]
		bs[h] = bs[f]
		bs[f] = temp
		_b := int(b[i]) ^ bs[(bs[h]+bs[f])%256]
		byteR = append(byteR, byte(_b))
		result += string(_b)
	}
	if base64 {
		result = utils.ToBase64(byteR)
	}
	return result
}

func (y *YoukuReptileV2) getEp(vid string, ep string) (newEp string, token string, sid string) {
	template1 := "becaf9be"
	template2 := "bf7e5f01"
	bytes, _ := utils.FromBase64(ep)
	ep = string(bytes)
	_temp := y.encoder(template1, bytes, false)
	parts := strings.Split(_temp, "_")
	sid = parts[0]
	token = parts[1]
	whole := fmt.Sprintf("%s_%s_%s", sid, vid, token)
	newbytes := []byte(whole)
	newEp = y.encoder(template2, newbytes, true)
	return
}

func (y *YoukuReptileV2) Reptile(url string) (*VodStreams, error) {
	vid, err := y.getVid(url)
	if err != nil {
		return nil, err
	}
	videoids_url := fmt.Sprintf(VIDEOIDS_URL_FMT, vid)
	data, err := httplib.Get(videoids_url).String()
	if err != nil {
		return nil, err
	}
	return y.reptile(data, vid)
}

func (y *YoukuReptileV2) CallbackReptile(parameter string, url string, cmd string) (vs *VodStreams, clientReqUrl string, nextCmd string, err error) {
	lcmd := strings.ToLower(cmd)
	if lcmd == "client_get" {
		_vid, _err := y.getVid(url)
		if _err != nil {
			return nil, "", REP_CALLBACK_COMMAND_EXIT, err
		}
		_vs, _err := y.reptile(parameter, _vid)
		return _vs, "", REP_CALLBACK_COMMAND_COMPLETED, _err
	}
	vid, err := y.getVid(url)
	if err != nil {
		return nil, "", "", err
	}
	videoids_url := fmt.Sprintf(VIDEOIDS_URL_FMT, vid)
	return nil, videoids_url, "client_get", nil
}

func (y *YoukuReptileV2) reptile(data string, vid string) (*VodStreams, error) {
	locker_v2.Lock()
	if time.Now().Sub(last_reptile_time_v2) < 1*time.Second { //1秒钟的间隔
		time.Sleep(1 * time.Second)
		last_reptile_time_v2 = time.Now()
	}
	locker_v2.Unlock()

	json, err := utils.NewJson([]byte(data))
	if err != nil {
		return nil, err
	}

	//start 校验数据
	if _, ok := json.CheckGet("data"); !ok {
		return nil, fmt.Errorf("解析错误1")
	}
	_das, err := json.Get("data").Array()
	if err != nil {
		return nil, err
	}
	if len(_das) == 0 {
		return nil, fmt.Errorf("解析错误2")
	}

	if _, ok := json.Get("data").GetIndex(0).CheckGet("ep"); !ok {
		return nil, fmt.Errorf("解析错误3")
	}
	if _, ok := json.Get("data").GetIndex(0).CheckGet("ip"); !ok {
		return nil, fmt.Errorf("解析错误4")
	}
	//end 校验数据

	ep, _ := json.Get("data").GetIndex(0).Get("ep").String()
	ip, _ := json.Get("data").GetIndex(0).Get("ip").Int()

	totalSeconds := json.Get("data").GetIndex(0).Get("seconds").MustFloat64()

	streamSizes := make(map[VOD_STREAM_MODE]int)
	streamTypes := []VOD_STREAM_MODE{}
	streamSegs := make(map[VOD_STREAM_MODE][]VodSeg)

	segs, err := json.Get("data").GetIndex(0).Get("segs").Map()
	if err != nil {
		return nil, err
	}
	streamtypes := json.Get("data").GetIndex(0).Get("streamtypes").MustArray()

	for _, mode := range streamtypes {
		_m := mode.(string)
		for seg, _ := range segs {
			if seg == _m {
				//m3u8
				streamMode := y.vodM3U8StreamMode(seg)
				if streamMode == VOD_STREAM_MODE_UNDEFINED {
					continue
				}
				_msize := json.Get("data").GetIndex(0).Get("streamsizes").MustMap()[_m]
				modeSize, _ := _msize.(enjson.Number).Int64()
				streamSizes[streamMode] = int(modeSize) //strconv.Atoi(modeSize)
				streamTypes = append(streamTypes, streamMode)
				vss := []VodSeg{}

				newEp, token, sid := y.getEp(vid, ep)
				m3u8url := fmt.Sprintf(YOUKU_M3U8_FMT, utils.UrlEncode(newEp), ip, sid, token, seg, vid)
				vss = append(vss, VodSeg{
					No:      0,
					Seconds: int(totalSeconds),
					Size:    streamSizes[streamMode],
					Url:     m3u8url,
				})
				streamSegs[streamMode] = vss
				//download flvs
				downSizes, downSegs := y.getDownSegs(data, seg, vid)
				for __m, __size := range downSizes {
					streamSizes[__m] = __size
				}
				for __m, __vss := range downSegs {
					streamSegs[__m] = __vss
				}
			}
		}
	}
	vs := &VodStreams{
		TotalSeconds: totalSeconds,
		StreamSizes:  streamSizes,
		StreamTypes:  streamTypes,
		Segs:         streamSegs,
	}
	return vs, nil
}

func (y *YoukuReptileV2) getDownSegs(json string, seg string, vid string) (map[VOD_STREAM_MODE]int, map[VOD_STREAM_MODE][]VodSeg) {
	reptile_host := beego.AppConfig.String("reptile.vod.down.host")
	data := utils.UrlEncode(json)
	url := fmt.Sprintf("http://%s/getVideo.php", reptile_host)
	req := httplib.Post(url)
	req.Param("data", data)
	req.Param("src", "youku")
	resultJson, err := req.String()
	if err != nil {
		return nil, nil
	}
	j, err := utils.NewJson([]byte(resultJson))
	if err != nil {
		return nil, nil
	}
	maps, err := j.Map()
	if err != nil {
		return nil, nil
	}
	srcj, _ := utils.NewJson([]byte(json))
	supportModes := []string{"normal", "high", "super"}
	getSegs := func(mode string) string {
		switch mode {
		case "normal":
			return "flv"
		case "high":
			return "mp4"
		case "super":
			return "hd2"
		default:
			return "undefined"
		}
	}
	vsmap := make(map[VOD_STREAM_MODE][]VodSeg)
	szmap := make(map[VOD_STREAM_MODE]int)
	for _, m := range supportModes {
		if obj, ok := maps[m]; ok {
			vss := []VodSeg{}
			flvs := obj.([]interface{})
			i := 0
			for _, flvurl := range flvs {
				_flvurl := flvurl.(string)
				qx, has := srcj.Get("data").GetIndex(0).Get("segs").CheckGet(getSegs(m))
				if !has {
					break
				}
				seconds, _ := qx.GetIndex(i).Get("seconds").Int()
				sizestr, _ := qx.GetIndex(i).Get("size").String()
				size, _ := strconv.Atoi(sizestr)

				vss = append(vss, VodSeg{
					No:      i,
					Seconds: seconds,
					Size:    size,
					Url:     _flvurl,
				})

				i++
			}
			vsmap[y.vodDownStreamMode(getSegs(m))] = vss

			sz, has := srcj.Get("data").GetIndex(0).Get("streamsizes").CheckGet(getSegs(m))
			if has {
				_size, _ := sz.String()
				szmap[y.vodDownStreamMode(getSegs(m))], _ = strconv.Atoi(_size)
			}
		}
	}
	return szmap, vsmap
}

//func (y *YoukuReptileV2) getSid() string {
//	timestamps := time.Now().Unix()
//	thrids := 10000 + rand.Intn(9000)
//	return fmt.Sprintf("%d%d", timestamps, thrids)
//}

//func (y *YoukuReptileV2) getKey(key1 string, key2 string) string {
//	a, _ := strconv.ParseInt(key1, 16, 32)
//	b := a ^ 0xA55AA5A5
//	c := strconv.FormatInt(b, 16)
//	return key2 + c
//}

//func (y *YoukuReptileV2) getMixString(seed int) string {
//	mixed := ""
//	source := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ/\\:._-1234567890"
//	_len := len(source)
//	for i := 0; i < _len; i++ {
//		seed = (seed*211 + 30031) % 65536
//		index := int((float32(seed) / float32(65536)) * float32(len(source)))
//		c := source[index]
//		mixed = mixed + string(c)
//		source = strings.Replace(source, string(c), "", -1)
//	}
//	return mixed
//}

//func (y *YoukuReptileV2) getFileId(fileId string, seed int) string {
//	mixed := y.getMixString(seed)
//	ids := strings.Split(fileId, "*")
//	realId := ""
//	for _, id := range ids {
//		if len(id) == 0 {
//			continue
//		}
//		idx, _ := strconv.Atoi(id)
//		realId += string(mixed[idx])
//	}
//	return realId
//}

//func (y *YoukuReptileV2) getYkd(s string) string {
//	f := len(s)
//	b := 0
//	str := `ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/`
//	c := ""
//	for b < f {
//		b++
//		e := y.charCodeAt(s, b) & 255
//		if b == f {
//			c += y.charAt(str, e>>2)
//			c += y.charAt(str, (e&3)<<4)
//			c += "=="
//			break
//		}
//		b++
//		g := y.charCodeAt(s, b)
//		if b == f {
//			c += y.charAt(str, e>>2)
//			c += y.charAt(str, (e&3)<<4|(g&240)>>4)
//			c += y.charAt(str, (g&15)<<2)
//			break
//		}
//		b++
//		h := y.charCodeAt(s, b)
//		c += y.charAt(str, e>>2)
//		c += y.charAt(str, (e&3)<<4|(g&240)>>4)
//		c += y.charAt(str, (g&15)<<2|(h&192)>>6)
//		c += y.charAt(str, h&63)
//	}
//	return c
//}

//func (y *YoukuReptileV2) getYkna(s string) string {
//	sz := "-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,62,-1,-1,-1,63,52,53,54,55,56,57,58,59,60,61,-1,-1,-1,-1,-1,-1,-1,0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,-1,-1,-1,-1,-1,-1,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,-1,-1,-1,-1,-1"
//	h := strings.Split(sz, ",")
//	i := len(s)
//	f := 0
//	for f < i {
//		f++
//		for {
//			c := h[y.charCodeAt(s, f)&255]
//			if !(f < i && c == "-1") {
//				break
//			}
//		}
//		if c == "-1" {
//			break
//		}
//		f++
//		for {
//			b := h[y.charCodeAt(s, f)&255]
//			if !(f < i && b == "-1") {
//				break
//			}
//		}
//		if b == "-1" {
//			break
//		}
//		e +=
//	}
//}

//func (y *YoukuReptileV2) charCodeAt(str string, index int) int {
//	key := utils.Md5(str)
//	i := index + 1
//	if s, ok := y.charCode[key]; ok {
//		return int(s[i])
//	}
//	y.charCode[key] = str
//	return int(str[i])
//}

//func (y *YoukuReptileV2) charAt(str string, index int) string {
//	return str[index : index+1]
//}

func (y *YoukuReptileV2) vodM3U8StreamMode(mode string) VOD_STREAM_MODE {
	switch mode {
	case "flv":
		return VOD_STREAM_MODE_MSTD
	case "hd2":
		return VOD_STREAM_MODE_MSUPER
	case "mp4":
		return VOD_STREAM_MODE_MHIGH
	case "hd3":
		return VOD_STREAM_MODE_M1080P
	default:
		return VOD_STREAM_MODE_UNDEFINED
	}
}

func (y *YoukuReptileV2) vodDownStreamMode(mode string) VOD_STREAM_MODE {
	switch mode {
	case "flv":
		return VOD_STREAM_MODE_STANDARD_SP
	case "hd2":
		return VOD_STREAM_MODE_SUPER_SP
	case "mp4":
		return VOD_STREAM_MODE_HIGH_SP
	case "hd3":
		return VOD_STREAM_MODE_1080P_SP
	default:
		return VOD_STREAM_MODE_UNDEFINED
	}
}

func (y *YoukuReptileV2) M3u8ToSegs(m3u8url string) ([]VodSeg, error) {
	m3u8txt, err := httplib.Get(m3u8url).String()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBufferString(m3u8txt)
	playlist, ts, _ := m3u8.Decode(*buf, false)
	if ts != m3u8.MEDIA {
		return nil, fmt.Errorf("格式错误")
	}
	mpls := playlist.(*m3u8.MediaPlaylist)
	flvs := make(map[string]bool)
	vodsegs := []VodSeg{}
	i := 0
	for _, seg := range mpls.Segments {
		if seg == nil {
			continue
		}
		uri, err := url.Parse(seg.URI)
		if err != nil {
			continue
		}
		url := fmt.Sprintf("%s://%s%s", uri.Scheme, uri.Host, uri.Path)
		if _, ok := flvs[url]; !ok {
			flvs[url] = true
			vodsegs = append(vodsegs, VodSeg{
				No:      i,
				Seconds: 0,
				Size:    0,
				Url:     url,
			})
			i++
		}
	}
	return vodsegs, nil
}
