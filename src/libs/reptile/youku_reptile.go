package reptile

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/httplib"
)

type XzJsonObj struct {
	Quality string        `json:"quality"`
	Site    string        `json:"site"`
	Title   string        `json:"title"`
	Files   []*XzJsonFile `json:"files"`
}

type XzJsonFile struct {
	Bytes  int    `json:"bytes"`
	FType  string `json:"ftype"`
	FUrl   string `json:"furl"`
	Second int    `json:"seconds"`
	Size   string `json:"size"`
	Time   string `json:"time"`
}

const (
	YOUKU_FLV_PATH_GET_URL = "http://f.youku.com/player/getFlvPath/sid/%s_%s/st/%s/fileid/%s?K=%s&ts=%d"
	YOUKU_GETLIST_VIDEOIDS = "http://v.youku.com/player/getPlayList/VideoIDS/"
	//token/b27592bb86ed6b26e78367530427f49d
	XZ_REPTILE_URL = "http://api.flvxz.com/token/8f36d905e6ebaf2feae395baddbf12e9/site/%s/vid/%s/jsonp/purejson"
)

var last_reptile_time time.Time = time.Now()
var locker *sync.Mutex = new(sync.Mutex)

type YoukuReptile struct{}

func NewYoukuReptile() *YoukuReptile {
	return &YoukuReptile{}
}

func (r *YoukuReptile) Reptile(url string) (*VodStreams, error) {
	locker.Lock()
	if time.Now().Sub(last_reptile_time) < 1*time.Second { //1秒钟的间隔
		time.Sleep(1 * time.Second)
		last_reptile_time = time.Now()
	}
	locker.Unlock()

	re := regexp.MustCompile("http://v.youku.com/v_show/id_(\\w+)\\.html")
	matchs := re.FindSubmatch([]byte(url))
	if len(matchs) != 2 {
		return nil, errors.New("youku url format fail")
	}
	yk_id := string(matchs[1])
	__idx := strings.Index(yk_id, "_")
	if __idx > 0 {
		yk_id = yk_id[0:__idx]
	}
	data_url := fmt.Sprintf(XZ_REPTILE_URL, "youku", yk_id)
	req := httplib.Get(data_url)
	rsp, err := req.Response()
	defer func() {
		if rsp != nil && !rsp.Close {
			rsp.Body.Close()
		}
	}()
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}
		xzs := []*XzJsonObj{}
		err = json.Unmarshal(body, &xzs)
		if err != nil {
			return nil, err
		}
		total_seconds := 0
		stream_sizes := make(map[VOD_STREAM_MODE]int)
		stream_types := []VOD_STREAM_MODE{}
		segs := map[VOD_STREAM_MODE][]VodSeg{}
		for _, v := range xzs {
			stream_mode := r.vodStreamModeV2(v.Quality)
			if stream_mode == VOD_STREAM_MODE_UNDEFINED {
				continue //未知清晰度略过
			}
			_no := 0
			_bytes := 0
			_seconds := 0
			vsegs := []VodSeg{}
			for _, f := range v.Files {
				//_b, _ := strconv.Atoi(f.Bytes)
				_b := f.Bytes
				_bytes += _b
				_seconds += f.Second
				_seg := VodSeg{
					No:      _no,
					Seconds: f.Second,
					Size:    _b,
					Url:     f.FUrl,
				}
				vsegs = append(vsegs, _seg)
				_no++
			}
			stream_sizes[stream_mode] = _bytes
			stream_types = append(stream_types, stream_mode)
			total_seconds += _seconds
			segs[stream_mode] = vsegs
		}
		vs := &VodStreams{
			TotalSeconds: float64(total_seconds),
			StreamSizes:  stream_sizes,
			StreamTypes:  stream_types,
			Segs:         segs,
		}
		return vs, nil
	}
	return nil, errors.New("http response status not ok")
}

func (r *YoukuReptile) vodStreamModeV2(mode string) VOD_STREAM_MODE {
	switch mode {
	case "标清", "FLV标清":
		return VOD_STREAM_MODE_STANDARD
	case "分段_标清_FLV":
		return VOD_STREAM_MODE_STANDARD_SP
	case "超清":
		return VOD_STREAM_MODE_SUPER
	case "分段_超清_FLV":
		return VOD_STREAM_MODE_SUPER_SP
	case "高清":
		return VOD_STREAM_MODE_HIGH
	case "分段_高清_MP4":
		return VOD_STREAM_MODE_HIGH_SP
	case "1080P":
		return VOD_STREAM_MODE_1080P
	case "分段_1080P_FLV":
		return VOD_STREAM_MODE_1080P_SP
	case "手机标清", "单段_标清_M3U8":
		return VOD_STREAM_MODE_MSTD
	case "高清M3U8", "单段_高清_M3U8":
		return VOD_STREAM_MODE_MHIGH
	case "单段_超清_M3U8":
		return VOD_STREAM_MODE_MSUPER
	case "单段_标清_MP4":
		return VOD_STREAM_MODE_MSTD_MP4
	default:
		return VOD_STREAM_MODE_UNDEFINED
	}
}

//youku url: http://v.youku.com/v_show/id_XNTgxMTk1NjI0.html
//func (r *YoukuReptile) Reptile(url string) (*VodStreams, error) {
//	re := regexp.MustCompile("http://v.youku.com/v_show/id_(\\w+)\\.html")
//	matchs := re.FindSubmatch([]byte(url))
//	if len(matchs) != 2 {
//		return nil, errors.New("youku url format fail")
//	}
//	yk_id := string(matchs[1])
//	__idx := strings.Index(yk_id, "_")
//	if __idx > 0 {
//		yk_id = yk_id[0:__idx]
//	}
//	data_url := YOUKU_GETLIST_VIDEOIDS + yk_id
//	req := httplib.Get(data_url)
//	rsp, err := req.Response()
//	defer func() {
//		if !rsp.Close {
//			rsp.Body.Close()
//		}
//	}()
//	if err != nil {
//		return nil, err
//	}
//	if rsp.StatusCode == http.StatusOK {
//		body, err := ioutil.ReadAll(rsp.Body)
//		if err != nil {
//			return nil, err
//		}
//		json_data, err := utils.NewJson(body)
//		if err != nil {
//			return nil, err
//		}
//		//解析youku数据
//		vs := new(VodStreams)
//		j_data := json_data.Get("data").GetIndex(0)
//		total_seconds_str, _ := j_data.Get("seconds").String()
//		total_seconds, _ := strconv.ParseFloat(total_seconds_str, 16)
//		stream_types, _ := j_data.Get("streamtypes").Array()
//		stream_sizes, _ := j_data.Get("streamsizes").Map()
//		stream_segs, _ := j_data.Get("segs").Map()
//		stream_fileids, _ := j_data.Get("streamfileids").Map()
//		stream_seed, _ := j_data.Get("seed").Int64()

//		vs.TotalSeconds = total_seconds
//		vs.StreamTypes = make([]VOD_STREAM_MODE, 0, len(stream_types))
//		for _, v := range stream_types {
//			if val, ok := v.(string); ok {
//				vs.StreamTypes = append(vs.StreamTypes, r.vodStreamMode(val))
//			}
//		}
//		vs.StreamSizes = make(map[VOD_STREAM_MODE]int)
//		for k, v := range stream_sizes {
//			if val, ok := v.(string); ok {
//				i, _ := strconv.Atoi(val)
//				vs.StreamSizes[r.vodStreamMode(k)] = int(i)
//			}
//		}

//		vs.Segs = make(map[VOD_STREAM_MODE][]VodSeg)
//		for k, v := range stream_segs {
//			if val, ok := stream_fileids[k].(string); ok {
//				mode := r.vodStreamMode(k)
//				vs.Segs[mode] = r.decodeVod(int(stream_seed), v.([]interface{}), val, mode)
//			}
//		}
//		return vs, nil
//	}

//	return nil, errors.New("http response status not ok")
//}

func (r *YoukuReptile) vodStreamMode(mode string) VOD_STREAM_MODE {
	switch mode {
	case "flv":
		return VOD_STREAM_MODE_STANDARD
	case "hd2":
		return VOD_STREAM_MODE_SUPER
	case "mp4":
		return VOD_STREAM_MODE_HIGH
	case "hd3":
		return VOD_STREAM_MODE_1080P
	default:
		return VOD_STREAM_MODE_UNDEFINED
	}
}

func (r *YoukuReptile) vodStreamModeAddr(mode VOD_STREAM_MODE) string {
	switch mode {
	case VOD_STREAM_MODE_HIGH:
		return "mp4"
	default:
		return "flv"
	}
}

func (r *YoukuReptile) decodeVod(seed int, sources []interface{}, fileid string, mode VOD_STREAM_MODE) []VodSeg {
	estr := r.getFileId(fileid, seed)
	vss := make([]VodSeg, 0, len(sources))
	for _, source := range sources {
		kvs := source.(map[string]interface{})
		s_key := kvs["k"].(string)
		var s_no int = 0
		if v, ok := kvs["no"].(string); ok {
			s_no, _ = strconv.Atoi(v)
		}
		//fmt.Println(reflect.TypeOf(kvs["seconds"]))
		//json无双引号数字解析成json.Number 冏
		var s_sec int64 = 0
		if v, ok := kvs["seconds"].(json.Number); ok {
			s_sec, _ = v.Int64()
		}
		s_size := kvs["size"].(string)
		s_sec16b := strconv.FormatInt(int64(s_no), 16) //hex.EncodeToString(utils.IntToBytes(int(s_sec)))
		s_sec16b = strings.ToUpper(s_sec16b)
		if len(s_sec16b) == 1 {
			s_sec16b = "0" + s_sec16b
		}

		s1 := estr[0:8]
		s2 := estr[10:]
		new_id := s1 + s_sec16b + s2
		sid := r.getSid()
		vs := VodSeg{}
		vs.No = s_no
		vs.Seconds = int(s_sec)
		vs.Size, _ = strconv.Atoi(s_size)
		vs.Url = fmt.Sprintf(YOUKU_FLV_PATH_GET_URL, sid, s_sec16b, r.vodStreamModeAddr(mode), new_id, s_key, s_sec)
		vss = append(vss, vs)
	}
	return vss
}

func (r *YoukuReptile) getMixString(seed int) string {
	mixed := ""
	source := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ/\\:._-1234567890"
	_len := len(source)
	for i := 0; i < _len; i++ {
		seed = (seed*211 + 30031) % 65536
		index := int((float32(seed) / float32(65536)) * float32(len(source)))
		c := source[index]
		mixed = mixed + string(c)
		source = strings.Replace(source, string(c), "", -1)
	}
	return mixed
}

func (r *YoukuReptile) getFileId(fileId string, seed int) string {
	mixed := r.getMixString(seed)
	ids := strings.Split(fileId, "*")
	realId := ""
	for _, id := range ids {
		if len(id) == 0 {
			continue
		}
		idx, _ := strconv.Atoi(id)
		realId += string(mixed[idx])
	}
	return realId
}

func (r *YoukuReptile) getSid() string {
	timestamps := time.Now().Unix()
	seconds := 1000 + rand.Intn(999)
	thrids := 1000 + rand.Intn(9000)
	return fmt.Sprintf("%d%d%d", timestamps, seconds, thrids)
}

//////////////////////////////////////////////////////////////////////////////////////////
//优酷用户专题视频抓取
//////////////////////////////////////////////////////////////////////////////////////////

type YoukuUserReptile struct{}

func (r *YoukuUserReptile) analysisHtmlToVodDatas(html string) []*RVodData {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(html)
	doc, err := goquery.NewDocumentFromReader(buf)
	arr := []*RVodData{}
	if err == nil {
		doc.Find(".yk-row").Each(func(i int, s *goquery.Selection) {
			s.Find(".yk-col4").Each(func(j int, ss *goquery.Selection) {
				cTime, _ := ss.Attr("c_time")
				img, _ := ss.Find(".v-thumb img").Attr("src")
				title, _ := ss.Find(".v-meta-title a").Attr("title")
				purl, _ := ss.Find(".v-link a").Attr("href")
				ctime, _ := utils.StrToTime(cTime)
				_url, err := url.Parse(purl)
				if err != nil {
					return
				}
				purl = fmt.Sprintf("%s://%s%s", _url.Scheme, _url.Host, _url.Path)
				rvd := &RVodData{title, img, "jpg", purl, ctime}
				arr = append(arr, rvd)
			})
		})
	}
	return arr
}

func (r *YoukuUserReptile) analysisPageUrlToParameters(url string) map[string]string {
	//simple:/u/UMzkzOTQ4ODcy/videos/order_1_view_1_page_2_spg_1_stt_2169_sid_183301486_sst_1403998645
	params := make(map[string]string)
	if len(url) == 0 {
		return params
	}
	last_idx := strings.LastIndex(url, "/")
	last_str := url[last_idx+1:]
	arrs := strings.Split(last_str, "_")
	for i := 0; i < len(arrs); i += 2 {
		pname := arrs[i]
		pvalue := arrs[i+1]
		params[pname] = pvalue
	}
	return params
}

func (r *YoukuUserReptile) analysisUrlPages(html string) (maps map[int]string, maxPage int) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(html)
	doc, err := goquery.NewDocumentFromReader(buf)
	maps = make(map[int]string)
	maxPage = 0
	if err == nil {
		doc.Find(".YK-pages a").Each(func(i int, s *goquery.Selection) {
			if _, exist := s.Attr("_h"); exist {
				pageNum := s.Text()
				href, _ := s.Attr("href")
				page, err := strconv.Atoi(pageNum)
				if err == nil {
					maps[page] = href
					if page > maxPage {
						maxPage = page
					}
				}
			}
		})
	}
	return maps, maxPage
}

func (r *YoukuUserReptile) reptileSinglePageData(page int, ukey string, last_str string) []*RVodData {
	url_format := "http://i.youku.com/u/%s/videos/fun_ajaxload/?__rt=1&__ro=&v_page=%d&page_num=%d&page_order=1&q=&last_str=%s"
	rvss := []*RVodData{}
	for i := 1; i <= 3; i++ {
		page_num := ((page - 1) * 3) + i
		post_url := fmt.Sprintf(url_format, ukey, page, page_num, utils.UrlEncode(last_str))
		html, err := httplib.Get(post_url).String()
		if err == nil {
			rvs := r.analysisHtmlToVodDatas(html)
			for _, v := range rvs {
				rvss = append(rvss, v)
			}
		}
	}
	return rvss
}

func (r *YoukuUserReptile) Reptile(url string, saveVodFunc func([]*RVodData) bool) error {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	__index := strings.LastIndex(url, "/")
	uKey := url[__index+1:]
	ucpage_url := "http://i.youku.com/u/%s/videos/view_1_order_1"
	get_url := fmt.Sprintf(ucpage_url, uKey)
	html, err := httplib.Get(get_url).String()
	if err != nil {
		return err
	}
	mpage := 1
	for {
		pages, _ := r.analysisUrlPages(html)
		pup := make(map[string]string)
		if len(pages) > 0 {
			for k, v := range pages { //后一页
				if k > mpage {
					pup = r.analysisPageUrlToParameters(v)
					break
				}
			}
			if len(pup) == 0 { //如果是最后一页,则选择前一页
				for k, v := range pages {
					if k == mpage-1 {
						pup = r.analysisPageUrlToParameters(v)
						break
					}
				}
			}
		}
		if len(pup) == 0 && len(pages) == 0 { //只有一页的特殊情况
			_rvd := r.reptileSinglePageData(mpage, uKey, "")
			saveVodFunc(_rvd)
			break
		}
		var j_page, j_vid, j_sortv, j_total string
		if v, ok := pup["spg"]; ok {
			j_page = v
		}
		if v, ok := pup["sid"]; ok {
			j_vid = v
		}
		if v, ok := pup["sst"]; ok {
			j_sortv = v
		}
		if v, ok := pup["stt"]; ok {
			j_total = v
		}
		last_str := fmt.Sprintf(`{"page":%s,"vid":%s,"sort_value":%s,"total":%s}`, j_page, j_vid, j_sortv, j_total)
		if len(pup) == 0 {
			last_str = ""
		}
		rvd := r.reptileSinglePageData(mpage, uKey, last_str)
		if len(rvd) == 0 {
			break
		}
		if !saveVodFunc(rvd) {
			return nil
		}

		mpage++
		ucpage_url = "http://i.youku.com/u/%s/videos/order_1_view_1_page_%d_spg_%s_stt_%s_sid_%s_sst_%s"
		post_url := fmt.Sprintf(ucpage_url, uKey, mpage, j_page, j_total, j_vid, j_sortv)
		_html, err := httplib.Get(post_url).String()
		html = _html
		if err != nil {
			return err
		}
		if mpage > 50 {
			return errors.New("页数不能超过50页")
		}
	}
	return errors.New("抓取失败")
}

func (r *YoukuUserReptile) ValidateUrl(url string) error {
	return errors.New("未实现此功能")
}
