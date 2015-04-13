package admincp

import (
	"utils"
)

// 缓存管理 API
type CacheCPController struct {
	AdminController
}

func (c *CacheCPController) Prepare() {
	c.AdminController.Prepare()
}

// @Title 清空缓存
// @Description 清空缓存
// @Param   ck   path string true  "key"
// @Success 200
// @router /clean [get]
func (c *CacheCPController) CleanCache() {
	key := c.GetString("ck")
	if len(key) > 0 {
		cache := utils.GetCache()
		localCache := utils.GetLocalCache()

		if key == "all" {
			cache.Flush()
			localCache.Flush()
		} else {
			cache.Delete(key)
			localCache.Delete(key)
		}
		c.WriteString("success")
		return
	}
	c.WriteString("not exist")
}
