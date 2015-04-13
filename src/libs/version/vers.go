package version

import (
	"dbs"
	"fmt"
	"strings"
	"time"
	"utils"
)

type VCS struct{}

func (v *VCS) ConvertPlatform(plat string) (error, MOBILE_PLATFORM) {
	switch plat {
	case string(MOBILE_PLATFORM_ANDROID):
		return nil, MOBILE_PLATFORM_ANDROID
	case string(MOBILE_PLATFORM_APPLE):
		return nil, MOBILE_PLATFORM_APPLE
	case string(MOBILE_PLATFORM_WPHONE):
		return nil, MOBILE_PLATFORM_WPHONE
	default:
		return fmt.Errorf("不存在平台标识"), ""
	}
}

func (v *VCS) Create(ver *ClientVersion) (int64, error) {
	if len(ver.DownloadUrl) == 0 {
		return 0, fmt.Errorf("下载地址不能为空")
	}
	if len(ver.Platform) == 0 {
		return 0, fmt.Errorf("未选择客户端平台")
	}
	if ver.Version <= 0 {
		return 0, fmt.Errorf("版本号必须大于0")
	}
	if len(ver.Ver) == 0 {
		return 0, fmt.Errorf("必须要有版本名称,类似:V1.000")
	}
	ver.Ver = strings.ToLower(ver.Ver)
	ver.PostTime = time.Now()
	o := dbs.NewDefaultOrm()
	id, err := o.Insert(ver)
	if err != nil {
		return 0, err
	}
	ver.Id = id
	cache := utils.GetCache()
	cache.Delete(v.cacheKey(ver.Platform))
	return id, nil
}

func (v *VCS) Update(ver *ClientVersion) error {
	ver.Ver = strings.ToLower(ver.Ver)
	o := dbs.NewDefaultOrm()
	_, err := o.Update(ver)
	if err != nil {
		return err
	}
	cache := utils.GetCache()
	cache.Delete(v.cacheKey(ver.Platform))
	return nil
}

func (v *VCS) Del(id int64, platform MOBILE_PLATFORM) error {
	ver := v.GetClientVersionById(platform, id)
	if ver == nil {
		return fmt.Errorf("版本不存在")
	}

	o := dbs.NewDefaultOrm()
	o.QueryTable(&ClientVersion{}).Filter("id", id).Delete()
	cache := utils.GetCache()
	cache.Delete(v.cacheKey(platform))
	return nil
}

func (v *VCS) cacheKey(platform MOBILE_PLATFORM) string {
	return fmt.Sprintf("mobile_version_plt:%s", platform)
}

func (v *VCS) GetClientVersions(platform MOBILE_PLATFORM) []*ClientVersion {
	cache := utils.GetCache()
	vers := []*ClientVersion{}
	err := cache.Get(v.cacheKey(platform), &vers)
	if err != nil {
		o := dbs.NewDefaultOrm()
		qs := o.QueryTable(&ClientVersion{})
		_, err = qs.Filter("platform", platform).OrderBy("-post_time").All(&vers)
		if err != nil {
			return vers
		}
		cache.Set(v.cacheKey(platform), vers, utils.StrToDuration("6h"))
	}
	return vers
}

func (v *VCS) GetLastClientVersion(platform MOBILE_PLATFORM) *ClientVersion {
	pvers := v.GetClientVersions(platform)
	var ver *ClientVersion
	var version float64 = 0
	for _, vr := range pvers {
		if vr.Version > version {
			ver = vr
			version = vr.Version
		}
	}
	return ver
}

func (v *VCS) GetClientVersion(platform MOBILE_PLATFORM, ver string) *ClientVersion {
	pvers := v.GetClientVersions(platform)
	vname := strings.ToLower(ver)
	for _, vr := range pvers {
		if vr.Ver == vname {
			return vr
		}
	}
	return nil
}

func (v *VCS) GetClientVersionById(platform MOBILE_PLATFORM, id int64) *ClientVersion {
	pvers := v.GetClientVersions(platform)
	for _, vr := range pvers {
		if vr.Id == id {
			return vr
		}
	}
	return nil
}
