package share

import (
	"fmt"
	"libs/vod"
	"regexp"
	"strconv"
)

///前端数据输出时进行特定数据装配

const (
	share_tag_vod_regex = `\[vod:(\d+)\]` //不带缩略图
	share_tag_pic_regex = `\[pic:(\d+)\]` //

	share_tag_vod_fmt = `[vod:%s]`
	share_tag_pic_fmt = `[pic:%s]`
)

//前端封装输出对象的通用代理对象
type ResOutputProxyObject struct {
	Id           string
	Title        string
	Content      string
	ThumbnailPic int64
	Uid          int64
}

var rts map[SHARE_KIND]ResTransformFunc = make(map[SHARE_KIND]ResTransformFunc)
var trs map[SHARE_KIND]TransformResFunc = make(map[SHARE_KIND]TransformResFunc)
var pfs map[SHARE_KIND]ResPicFileIdFunc = make(map[SHARE_KIND]ResPicFileIdFunc)
var outps map[SHARE_KIND]ResToOutputProxyObjectFunc = make(map[SHARE_KIND]ResToOutputProxyObjectFunc)

type ResTransformFunc func(resource string) (SHARE_KIND, string, error)
type TransformResFunc func(id string, args ...string) string
type ResPicFileIdFunc func(res *ShareResource) int64
type ResToOutputProxyObjectFunc func(res *ShareResource) *ResOutputProxyObject

////////////////////////////////////////////////////////////////////////////////////////////////////
func RegisterShareResourceTransformFuncs(shareKind SHARE_KIND, Func ResTransformFunc) {
	if _, ok := rts[shareKind]; ok {
		return
	}
	rts[shareKind] = Func
}

func ShareResourceTransformFuncs() map[SHARE_KIND]ResTransformFunc {
	return rts
}

////////////////////////////////////////////////////////////////////////////////////////////////////
func RegisterShareTransformResourceFuncs(shareKind SHARE_KIND, Func TransformResFunc) {
	if _, ok := trs[shareKind]; ok {
		return
	}
	trs[shareKind] = Func
}

func ShareTransformResourceFunc(shareKind SHARE_KIND) TransformResFunc {
	if f, ok := trs[shareKind]; ok {
		return f
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////
func RegisterResPicFileIdFuncs(shareKind SHARE_KIND, Func ResPicFileIdFunc) {
	if _, ok := pfs[shareKind]; ok {
		return
	}
	pfs[shareKind] = Func
}

func ShareResPicFileIdFunc(shareKind SHARE_KIND) ResPicFileIdFunc {
	if f, ok := pfs[shareKind]; ok {
		return f
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////
func RegisterResToOutputProxyObjectFuncs(shareKind SHARE_KIND, Func ResToOutputProxyObjectFunc) {
	if _, ok := outps[shareKind]; ok {
		return
	}
	outps[shareKind] = Func
}

func ShareResToOutputProxyObjectFunc(shareKind SHARE_KIND) ResToOutputProxyObjectFunc {
	if f, ok := outps[shareKind]; ok {
		return f
	}
	return nil
}

func initShareKindFuncs() {
	//注册资源转换方法
	RegisterShareTransformResourceFuncs(SHARE_KIND_VOD, func(id string, args ...string) string {
		return fmt.Sprintf(share_tag_vod_fmt, id)
	})
	RegisterShareTransformResourceFuncs(SHARE_KIND_PIC, func(id string, args ...string) string {
		return fmt.Sprintf(share_tag_pic_fmt, id)
	})

	//注册资源处理方法
	RegisterShareResourceTransformFuncs(SHARE_KIND_VOD, func(resource string) (SHARE_KIND, string, error) {
		if ok, _ := regexp.MatchString(share_tag_vod_regex, resource); ok {
			rep := regexp.MustCompile(share_tag_vod_regex)
			arr := rep.FindStringSubmatch(resource)
			return SHARE_KIND_VOD, arr[1], nil
		}
		return SHARE_KIND_EMPTY, "", fmt.Errorf("未匹配VOD_KIND")
	})
	RegisterShareResourceTransformFuncs(SHARE_KIND_PIC, func(resource string) (SHARE_KIND, string, error) {
		if ok, _ := regexp.MatchString(share_tag_pic_regex, resource); ok {
			rep := regexp.MustCompile(share_tag_pic_regex)
			arr := rep.FindStringSubmatch(resource)
			return SHARE_KIND_PIC, arr[1], nil
		}
		return SHARE_KIND_EMPTY, "", fmt.Errorf("未匹配PIC_KIND")
	})

	//注册资源图片文件id
	RegisterResPicFileIdFuncs(SHARE_KIND_VOD, func(res *ShareResource) int64 {
		_id, _ := strconv.ParseInt(res.Id, 10, 64)
		vods := &vod.Vods{}
		v := vods.Get(_id, false)
		if v == nil {
			return 0
		}
		return v.Img
	})
	RegisterResPicFileIdFuncs(SHARE_KIND_PIC, func(res *ShareResource) int64 {
		_id, _ := strconv.ParseInt(res.Id, 10, 64)
		return _id
	})
}
