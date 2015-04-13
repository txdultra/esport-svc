package libs

import (
	"dbs"
	"errors"
	"fmt"
	"time"
	"utils"

	"github.com/astaxie/beego/orm"
)

func NewRecommendService() *RecommendService {
	return &RecommendService{}
}

type RecommendService struct{}

func (c *RecommendService) cacheKey(id int64) string {
	return fmt.Sprintf("mobile_recommend_id:%d", id)
}

func (c *RecommendService) categoryCacheKey(category string) string {
	return fmt.Sprintf("mobile_recommeds_category:%s", category)
}

func (c *RecommendService) Create(recommend Recommend) (int64, error) {
	//if len(recommend.Title) == 0 {
	//	return 0, errors.New("Title属性不能为空")
	//}
	if len(recommend.Category) == 0 {
		return 0, errors.New("Category不能为空")
	}
	o := dbs.NewDefaultOrm()
	recommend.PostTime = time.Now()
	id, err := o.Insert(&recommend)
	if err != nil {
		return 0, errors.New("保存数据失败:" + err.Error())
	}
	recommend.Id = id
	cache := utils.GetCache()
	cache.Set(c.cacheKey(id), recommend, 5*time.Hour)
	cache.Delete(c.categoryCacheKey(recommend.Category))
	return id, nil
}

func (c *RecommendService) Delete(id int64) error {
	recommend := c.Get(id)
	if recommend == nil {
		return errors.New("不存在指定的推荐项目")
	}
	o := dbs.NewDefaultOrm()
	o.Delete(recommend)
	recommend.Id = id
	cache := utils.GetCache()
	cache.Delete(c.cacheKey(recommend.Id))
	cache.Delete(c.categoryCacheKey(recommend.Category))
	return nil
}

func (c *RecommendService) Update(recommend Recommend) error {
	if len(recommend.Title) == 0 {
		return errors.New("Title属性不能为空")
	}
	if len(recommend.Category) == 0 {
		return errors.New("Category不能为空")
	}
	o := dbs.NewDefaultOrm()
	recommend.PostTime = time.Now()
	_, err := o.Update(&recommend)
	if err == nil {
		cache := utils.GetCache()
		cache.Replace(c.cacheKey(recommend.Id), recommend, 5*time.Hour)
		cache.Delete(c.categoryCacheKey(recommend.Category))
		return nil
	}
	return err
}

func (c *RecommendService) Get(id int64) *Recommend {
	cache := utils.GetCache()
	recommend := Recommend{}
	err := cache.Get(c.cacheKey(id), &recommend)
	if err != nil {
		o := dbs.NewDefaultOrm()
		recommend.Id = id
		err = o.Read(&recommend)
		if err == orm.ErrNoRows || err == orm.ErrMissPK || err != nil {
			return nil
		}
		cache.Set(c.cacheKey(id), recommend, 5*time.Hour)
	}
	return &recommend
}

func (c *RecommendService) Gets(category string) []*Recommend {
	ckey := c.categoryCacheKey(category)
	cache := utils.GetCache()
	ids := []int64{}
	err := cache.Get(ckey, &ids)
	if err != nil {
		o := dbs.NewDefaultOrm()
		qs := o.QueryTable(&Recommend{}).Filter("category", category).OrderBy("-display_order", "-enabled")
		var lst []*Recommend
		_, err := qs.All(&lst)
		if err == nil {
			for _, m := range lst {
				ids = append(ids, m.Id)
			}
		}
	}
	cache.Set(ckey, ids, 5*time.Hour)
	result := []*Recommend{}
	for _, id := range ids {
		recommend := c.Get(id)
		if recommend != nil {
			result = append(result, recommend)
		}
	}
	return result
}
