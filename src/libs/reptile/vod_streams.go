package reptile

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type VOD_STREAM_MODE string

const (
	VOD_STREAM_MODE_STANDARD    VOD_STREAM_MODE = "std"
	VOD_STREAM_MODE_STANDARD_SP VOD_STREAM_MODE = "std_sp"
	VOD_STREAM_MODE_HIGH        VOD_STREAM_MODE = "high"
	VOD_STREAM_MODE_HIGH_SP     VOD_STREAM_MODE = "high_sp"
	VOD_STREAM_MODE_SUPER       VOD_STREAM_MODE = "super"
	VOD_STREAM_MODE_SUPER_SP    VOD_STREAM_MODE = "super_sp"
	VOD_STREAM_MODE_1080P       VOD_STREAM_MODE = "1080p"
	VOD_STREAM_MODE_1080P_SP    VOD_STREAM_MODE = "1080p_sp"
	VOD_STREAM_MODE_MSTD        VOD_STREAM_MODE = "m_std"
	VOD_STREAM_MODE_MSTD_MP4    VOD_STREAM_MODE = "m_std_mp4"
	VOD_STREAM_MODE_MHIGH       VOD_STREAM_MODE = "m_high"
	VOD_STREAM_MODE_MSUPER      VOD_STREAM_MODE = "m_super"
	VOD_STREAM_MODE_M1080P      VOD_STREAM_MODE = "m_1080p"
	VOD_STREAM_MODE_UNDEFINED   VOD_STREAM_MODE = "undefined"
)

var ALL_VOD_STREAM_MODES []VOD_STREAM_MODE = []VOD_STREAM_MODE{
	VOD_STREAM_MODE_STANDARD,
	VOD_STREAM_MODE_HIGH,
	VOD_STREAM_MODE_SUPER,
	VOD_STREAM_MODE_1080P,
	VOD_STREAM_MODE_MSTD,
	VOD_STREAM_MODE_MHIGH,
	VOD_STREAM_MODE_MSTD_MP4,
	VOD_STREAM_MODE_MSUPER,
	VOD_STREAM_MODE_M1080P,
	VOD_STREAM_MODE_UNDEFINED,
}

type VOD_SOURCE string

const (
	VOD_SOURCE_YOUKU VOD_SOURCE = "youku"
	VOD_SOURCE_NONE  VOD_SOURCE = "none"
)

func GetVodSource(url string) VOD_SOURCE {
	var lowerUrl = strings.ToLower(url)
	if strings.Contains(lowerUrl, "http://v.youku.com/v_show/id_") {
		return VOD_SOURCE_YOUKU
	}
	return VOD_SOURCE_NONE
}

func GetUcSource(url string) VOD_SOURCE {
	var lowerUrl = strings.ToLower(url)
	matched, _ := regexp.MatchString("http://i.youku.com/u/(\\w+)", lowerUrl)
	if matched {
		return VOD_SOURCE_YOUKU
	}
	return VOD_SOURCE_NONE
}

type VodStreams struct {
	TotalSeconds float64                      `json:"total_seconds"`
	StreamSizes  map[VOD_STREAM_MODE]int      `json:"stream_sizes"`
	StreamTypes  []VOD_STREAM_MODE            `json:"stream_types"`
	Segs         map[VOD_STREAM_MODE][]VodSeg `json:"segs"`
}

type VodSeg struct {
	No      int    `json:"no"`
	Seconds int    `json:"seconds"`
	Size    int    `json:"size"`
	Url     string `json:"url"`
}

type RVodData struct {
	Title    string    `json:"title"`
	Image    string    `json:"img"`
	ImgExt   string    `json:"img_ext"`
	PlayUrl  string    `json:"play_url"`
	PostTime time.Time `json:"time"`
}

func ReptileService(source VOD_SOURCE) IReptile {
	switch source {
	case VOD_SOURCE_YOUKU:
		return IReptile(NewYoukuReptileV3())
	default:
		return nil
	}
}

func CallbackReptileService(source VOD_SOURCE) ICallbackReptile {
	switch source {
	case VOD_SOURCE_YOUKU:
		return ICallbackReptile(NewYoukuReptileV3())
	default:
		return nil
	}
}

func ReptileUcService(source VOD_SOURCE) IUserReptile {
	switch source {
	case VOD_SOURCE_YOUKU:
		return IUserReptile(&YoukuUserReptile{})
	default:
		return nil
	}
}

func ConvertVodStreamMode(m string) VOD_STREAM_MODE {
	mm := strings.ToLower(m)
	switch mm {
	case "std":
		return VOD_STREAM_MODE_STANDARD
	case "std_sp":
		return VOD_STREAM_MODE_STANDARD_SP
	case "high":
		return VOD_STREAM_MODE_HIGH
	case "high_sp":
		return VOD_STREAM_MODE_HIGH_SP
	case "super":
		return VOD_STREAM_MODE_SUPER
	case "super_sp":
		return VOD_STREAM_MODE_SUPER_SP
	case "1080p":
		return VOD_STREAM_MODE_1080P
	case "1080p_sp":
		return VOD_STREAM_MODE_1080P_SP
	case "m_std":
		return VOD_STREAM_MODE_MSTD
	case "m_std_mp4":
		return VOD_STREAM_MODE_MSTD_MP4
	case "m_high":
		return VOD_STREAM_MODE_MHIGH
	case "m_super":
		return VOD_STREAM_MODE_MSUPER
	case "m_1080p":
		return VOD_STREAM_MODE_M1080P
	default:
		return VOD_STREAM_MODE_UNDEFINED
	}
}

func ConvertVodModeName(mode VOD_STREAM_MODE) string {
	switch mode {
	case VOD_STREAM_MODE_STANDARD, VOD_STREAM_MODE_STANDARD_SP:
		return "标清"
	case VOD_STREAM_MODE_HIGH, VOD_STREAM_MODE_HIGH_SP:
		return "高清"
	case VOD_STREAM_MODE_SUPER, VOD_STREAM_MODE_SUPER_SP:
		return "超清"
	case VOD_STREAM_MODE_1080P, VOD_STREAM_MODE_1080P_SP:
		return "1080P"
	case VOD_STREAM_MODE_MSTD:
		return "手机标清"
	case VOD_STREAM_MODE_MSTD_MP4:
		return "手机标清mp4"
	case VOD_STREAM_MODE_MHIGH:
		return "手机高清"
	case VOD_STREAM_MODE_MSUPER:
		return "手机超清"
	case VOD_STREAM_MODE_UNDEFINED:
		return "未知"
	default:
		return string(mode)
	}
}

func BuildM3u8(vods *VodStreams, mode VOD_STREAM_MODE) string {
	maxSeconds := 0
	lst, ok := vods.Segs[mode]
	if ok {
		for _, v := range lst {
			if v.Seconds > maxSeconds {
				maxSeconds = v.Seconds
			}
		}
	} else {
		return "not exist " + string(mode)
	}
	buffer := bytes.NewBufferString("")
	buffer.WriteString("#EXTM3U\n")
	buffer.WriteString("#EXT-X-TARGETDURATION:" + strconv.Itoa(maxSeconds) + "\n")
	buffer.WriteString("#EXT-X-VERSION:2\n")
	for _, v := range lst {
		buffer.WriteString("#EXTINF:" + strconv.Itoa(v.Seconds) + ",\n")
		buffer.WriteString(v.Url + "\n")
	}
	buffer.WriteString("#EXT-X-ENDLIST\n")
	return buffer.String()
}
