package reptile

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"utils"

	_ "code.google.com/p/go-charset/data"
	"code.google.com/p/mahonia"
	"github.com/m3u8"

	"github.com/astaxie/beego/httplib"
)

const (
	YOUKU_M3U8_FMT_FORV3   = "http://pl.youku.com/playlist/m3u8?ctype=12&ep=%s&ev=1&keyframe=1&oip=%d&sid=%s&token=%s&type=%s&vid=%s"
	VIDEOIDS_URL_FMT_FORV3 = "http://play.youku.com/play/get.json?vid=%s&ct=12"
	YOUKU_REP_USERAGENT    = "Mozilla/5.0 (iPhone; CPU iPhone OS 8_0_2 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile Safari/600.1.4"
	YOUKU_REP_COOKIE       = "r=\"8MaZyz7N UYvw0SPC D/S/7tik9Yye12tQKK2cPfLLL4J/AcsHURdYGBYIpRapcVRsgxjSz0UAQRjGlIAtiHvw==\";"
)

var last_reptile_time_v3 time.Time = time.Now()
var locker_v3 *sync.Mutex = new(sync.Mutex)

type YoukuReptileV3 struct {
}

func NewYoukuReptileV3() *YoukuReptileV3 {
	yr := &YoukuReptileV3{}
	return yr
}

func (y *YoukuReptileV3) getVid(url string) (string, error) {
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

func (y *YoukuReptileV3) encoder(a string, b []byte) []byte {
	//result := ""
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
		//result += string(_b)
	}
	//	if base64 {
	//		result = utils.ToWebBase64(byteR)
	//	}
	return byteR
}

func (y *YoukuReptileV3) generateEp(no int, streamFieldId string, sid string, token string) (fileid string, ep string) {
	number := strings.ToUpper(fmt.Sprintf("%x", no))
	if len(number) == 1 {
		number = "0" + number
	}
	fileid = streamFieldId[:8] + number + streamFieldId[10:]
	ep = y.getNewEp(sid, fileid, token)
	return
}

func (y *YoukuReptileV3) getNewEp(sid string, fileid string, token string) string {
	str := fmt.Sprintf("%s_%s_%s", sid, fileid, token)
	_ep := y.encoder("bf7e5f01", []byte(str))
	ep := utils.ToWebBase64(_ep)
	ep = strings.Replace(ep, "-", "%2B", -1)
	ep = strings.Replace(ep, "_", "%2F", -1)
	return ep
}

func (y *YoukuReptileV3) getEp(vid string, ep string) (newEp string, token string, sid string) {
	template1 := "becaf9be"
	template2 := "bf7e5f01"
	bytes, _ := utils.FromBase64(ep)
	ep = string(bytes)
	_temp := y.encoder(template1, bytes)
	parts := strings.Split(string(_temp), "_")
	sid = parts[0]
	token = parts[1]
	whole := fmt.Sprintf("%s_%s_%s", sid, vid, token)
	newbytes := []byte(whole)
	newEp = utils.ToBase64(y.encoder(template2, newbytes))
	return
}

func (y *YoukuReptileV3) Reptile(url string) (*VodStreams, error) {
	vid, err := y.getVid(url)
	if err != nil {
		return nil, err
	}
	videoids_url := fmt.Sprintf(VIDEOIDS_URL_FMT_FORV3, vid)
	req := httplib.Get(videoids_url)
	req.Header("Referer", "http://static.youku.com/")
	req.SetUserAgent(YOUKU_REP_USERAGENT)
	req.Header("Cookie", fmt.Sprintf("__ysuid=%d", time.Now().Unix()))
	data, err := req.String()
	if err != nil {
		return nil, err
	}
	return y.reptile(data, vid)
}

func (y *YoukuReptileV3) CallbackReptile(parameter string, url string, cmd string) (vs *VodStreams, clientReqUrl string, nextCmd string, err error) {
	lcmd := strings.ToLower(cmd)
	if lcmd == "client_get" {
		_vid, _err := y.getVid(url)
		if _err != nil {
			return nil, "", REP_CALLBACK_COMMAND_EXIT, err
		}
		videoids_url := fmt.Sprintf(VIDEOIDS_URL_FMT_FORV3, _vid)
		req := httplib.Get(videoids_url)
		req.Header("Referer", fmt.Sprintf("http://static.youku.com/", _vid))
		req.SetUserAgent(YOUKU_REP_USERAGENT)
		req.Header("Cookie", fmt.Sprintf("__ysuid=%d", time.Now().Unix()))
		data, err := req.String()
		if err != nil {
			return nil, "", "", err
		}
		_vs, _err := y.reptile(data, _vid)
		return _vs, "", REP_CALLBACK_COMMAND_COMPLETED, _err
	}
	vid, err := y.getVid(url)
	if err != nil {
		return nil, "", "", err
	}
	videoids_url := fmt.Sprintf(VIDEOIDS_URL_FMT_FORV3, vid)
	return nil, videoids_url, "client_get", nil
}

func (y *YoukuReptileV3) reptile(data string, vid string) (*VodStreams, error) {
	locker_v3.Lock()
	if time.Now().Sub(last_reptile_time_v3) < 5*time.Second { //1秒钟的间隔
		time.Sleep(5 * time.Second)
		last_reptile_time_v3 = time.Now()
	}
	locker_v3.Unlock()

	cachekey := fmt.Sprintf("youku_rep_vodstreams_%s", vid)
	cache := utils.GetCache()
	vs := &VodStreams{}
	err := cache.Get(cachekey, vs)
	if err == nil {
		return vs, nil
	}

	json, err := utils.NewJson([]byte(data))
	if err != nil {
		return nil, err
	}

	//start 校验数据
	if _, ok := json.CheckGet("data"); !ok {
		return nil, fmt.Errorf("解析错误1")
	}
	if _, ok := json.Get("data").Get("security").CheckGet("encrypt_string"); !ok {
		return nil, fmt.Errorf("解析错误3")
	}
	if _, ok := json.Get("data").Get("security").CheckGet("ip"); !ok {
		return nil, fmt.Errorf("解析错误4")
	}
	//end 校验数据

	ep, _ := json.Get("data").Get("security").Get("encrypt_string").String()
	ip, _ := json.Get("data").Get("security").Get("ip").Int()

	_tseconds, _ := json.Get("video").Get("seconds").String()
	totalSeconds, _ := strconv.Atoi(_tseconds)

	streamSizes := make(map[VOD_STREAM_MODE]int)
	streamTypes := []VOD_STREAM_MODE{}
	streamSegs := make(map[VOD_STREAM_MODE][]VodSeg)

	streams, _ := json.Get("data").Get("stream").Array()
	if streams == nil {
		return nil, fmt.Errorf("解析错误5")
	}

	//stream_ids := []string{"flv", "mp4", "hd2", "hd3", "hd"}

	//file_storage := libs.NewFileStorage()

	for _, stream := range streams {
		smaps := stream.(map[string]interface{})
		mode := smaps["stream_type"].(string)
		streamMode := y.vodM3U8StreamMode(mode)
		if streamMode == VOD_STREAM_MODE_UNDEFINED {
			continue
		}

		//new_ep, sid, token := y.getEp(vid, ep)
		//YOUKU_M3U8_FMT_FORV3   = "http://pl.youku.com/playlist/m3u8?ctype=12&ep=%s&ev=1&keyframe=1&oip=%d&sid=%s&token=%s&type=%s&vid=%s"
		//		m3u8url := fmt.Sprintf(YOUKU_M3U8_FMT_FORV3, utils.UrlEncode(new_ep), ip, token, sid, mode, vid)

		//		fmt.Println(m3u8url)

		//		req := httplib.Get(m3u8url)
		//		req.SetUserAgent(YOUKU_REP_USERAGENT)
		//		req.Header("Cookie", YOUKU_REP_COOKIE)

		//		m3u8_name := fmt.Sprintf("%s_%s.m3u8", vid, mode)
		//		m3u8data, _ := req.Bytes()

		//		node, err := file_storage.SaveFile(m3u8data, m3u8_name, 0)
		//		if err != nil {
		//			continue
		//		}
		streamField := smaps["stream_fileid"].(string)
		vss := []VodSeg{}
		segs := smaps["segs"].([]interface{})

		enc := mahonia.NewDecoder("ascii")
		_ep := enc.ConvertString(ep)
		_bep, _ := base64.StdEncoding.DecodeString(_ep)
		_bepEcd := y.encoder("becaf9be", _bep)

		_sts := strings.Split(string(_bepEcd), "_")
		sid := _sts[0]
		token := _sts[1]
		for i, seg := range segs {
			_seg := seg.(map[string]interface{})
			_key := _seg["key"].(string)
			//			s_sec16b := strconv.FormatInt(int64(i), 16) //hex.EncodeToString(utils.IntToBytes(int(s_sec)))
			//			s_sec16b = strings.ToUpper(s_sec16b)
			//			if len(s_sec16b) == 1 {
			//				s_sec16b = "0" + s_sec16b
			//			}
			fileid, _tep := y.generateEp(i, streamField, sid, token)
			//_gep, _ := utils.UrlDecode(_tep)
			_q := fmt.Sprintf("ctype=12&ev=1&K=%s&ep=%s&oip=%d&token=%s&yxon=1", _key, _tep, ip, token)
			playurl := fmt.Sprintf("http://k.youku.com/player/getFlvPath/sid/%s_00/st/%s/fileid/%s?%s", sid, y.getContainer(mode), fileid, _q)

			playData, _jerr := httplib.Get(playurl).Bytes()
			if _jerr != nil {
				break
			}
			playJson, _jerr := utils.NewJson(playData)
			if _jerr != nil {
				break
			}
			if arr, _jerr := playJson.Array(); _jerr == nil && len(arr) > 0 {
				playMaps, ok := arr[0].(map[string]interface{})
				if ok {
					if _tmpUrl, ok := playMaps["server"].(string); ok {
						playurl = _tmpUrl
					}
				}
			}
			vss = append(vss, VodSeg{
				No:      i,
				Seconds: int(totalSeconds),
				Size:    streamSizes[streamMode],
				Url:     playurl,
			})
		}
		streamSegs[streamMode] = vss

		//		streamSizes[streamMode] = 0
		//		streamTypes = append(streamTypes, streamMode)
		//		vss := []VodSeg{}
		//		vss = append(vss, VodSeg{
		//			No:      0,
		//			Seconds: int(totalSeconds),
		//			Size:    streamSizes[streamMode],
		//			Url:     file_storage.GetFileUrl(node.FileId), //m3u8url,
		//		})
		//		streamSegs[streamMode] = vss

		//Obsolete
		//		_m := mode.(string)
		//		for seg, _ := range segs {
		//			if seg == _m {
		//				//m3u8
		//				streamMode := y.vodM3U8StreamMode(seg)
		//				if streamMode == VOD_STREAM_MODE_UNDEFINED {
		//					continue
		//				}
		//				_msize := json.Get("data").GetIndex(0).Get("streamsizes").MustMap()[_m]
		//				modeSize, _ := _msize.(enjson.Number).Int64()
		//				streamSizes[streamMode] = int(modeSize) //strconv.Atoi(modeSize)
		//				streamTypes = append(streamTypes, streamMode)
		//				vss := []VodSeg{}

		//				newEp, token, sid := y.getEp(vid, ep)
		//				m3u8url := fmt.Sprintf(YOUKU_M3U8_FMT, utils.UrlEncode(newEp), ip, sid, token, seg, vid)
		//				vss = append(vss, VodSeg{
		//					No:      0,
		//					Seconds: int(totalSeconds),
		//					Size:    streamSizes[streamMode],
		//					Url:     m3u8url,
		//				})
		//				streamSegs[streamMode] = vss
		//				//download flvs
		//				downSizes, downSegs := y.getDownSegs(data, seg, vid)
		//				for __m, __size := range downSizes {
		//					streamSizes[__m] = __size
		//				}
		//				for __m, __vss := range downSegs {
		//					streamSegs[__m] = __vss
		//				}
		//			}
		//		}
	}
	vs = &VodStreams{
		TotalSeconds: float64(totalSeconds),
		StreamSizes:  streamSizes,
		StreamTypes:  streamTypes,
		Segs:         streamSegs,
	}
	cache.Set(cachekey, *vs, 30*time.Minute)
	return vs, nil
}

func (y *YoukuReptileV3) getDownSegs(json string, seg string, vid string) (map[VOD_STREAM_MODE]int, map[VOD_STREAM_MODE][]VodSeg) {
	return nil, nil
}

func (y *YoukuReptileV3) getContainer(mode string) string {
	switch mode {
	case "mp4", "mp4hd":
		return "mp4"
	case "mp4hd3":
		return "flv"
	default:
		return "flv"
	}
}

func (y *YoukuReptileV3) vodM3U8StreamMode(mode string) VOD_STREAM_MODE {
	switch mode {
	case "flv", "flvhd":
		return VOD_STREAM_MODE_MSTD
	case "hd2", "mp4hd2":
		return VOD_STREAM_MODE_MSUPER
	case "mp4", "mp4hd":
		return VOD_STREAM_MODE_MHIGH
	case "hd3", "mp4hd3":
		return VOD_STREAM_MODE_M1080P
	default:
		return VOD_STREAM_MODE_UNDEFINED
	}
}

func (y *YoukuReptileV3) vodDownStreamMode(mode string) VOD_STREAM_MODE {
	switch mode {
	case "flv", "flvhd":
		return VOD_STREAM_MODE_STANDARD_SP
	case "hd2", "mp4hd2":
		return VOD_STREAM_MODE_SUPER_SP
	case "mp4", "mp4hd":
		return VOD_STREAM_MODE_HIGH_SP
	case "hd3", "mp4hd3":
		return VOD_STREAM_MODE_1080P_SP
	default:
		return VOD_STREAM_MODE_UNDEFINED
	}
}

func (y *YoukuReptileV3) M3u8ToSegs(m3u8url string) ([]VodSeg, error) {
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
