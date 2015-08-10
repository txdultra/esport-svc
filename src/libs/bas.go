package libs

import (
	"dbs"
	"errors"
	"time"
	"utils"
)

type Bas struct{}

func (b *Bas) gameCacheKey() string {
	return "mobile_base_games"
}

func (b *Bas) Games() []Game {
	cache := utils.GetCache()
	gs := &[]Game{}
	err := cache.Get(b.gameCacheKey(), gs)
	if err == nil {
		return *gs
	}
	var games []Game
	o := dbs.NewDefaultOrm()
	_, err = o.QueryTable(&Game{}).OrderBy("-display_order").All(&games)
	if err == nil {
		cache.Set(b.gameCacheKey(), games, 48*time.Hour)
	}
	return games
}

func (b *Bas) GetGame(id int) *Game {
	if id <= 0 {
		return nil
	}
	games := b.Games()
	for _, v := range games {
		if v.Id == id {
			return &v
		}
	}
	return nil
}

func (b *Bas) AddGame(game *Game) error {
	o := dbs.NewDefaultOrm()
	nums, err := o.QueryTable(&Game{}).Filter("en", game.En).Count()
	if err != nil {
		return errors.New(err.Error())
	}
	if nums > 0 {
		return errors.New("已存在相同的游戏")
	}
	id, err := o.Insert(game)
	if err == nil && id > 0 {
		cache := utils.GetCache()
		cache.Delete(b.gameCacheKey())
		return nil
	}
	return err
}

func (b *Bas) UpdateGame(game *Game) error {
	o := dbs.NewDefaultOrm()
	num, err := o.Update(game)
	if err != nil {
		return errors.New(err.Error())
	}
	if err == nil && num > 0 {
		cache := utils.GetCache()
		cache.Delete(b.gameCacheKey())
		return nil
	}
	return errors.New("不存在对应的游戏")
}

//matchs
func (b *Bas) matchCacheKey() string {
	return "mobile_base_matchs"
}

func (b *Bas) Matchs() []Match {
	cache := utils.GetCache()
	ms := &[]Match{}
	err := cache.Get(b.matchCacheKey(), ms)
	if err == nil {
		return *ms
	}
	var matchs []Match
	o := dbs.NewDefaultOrm()
	_, err = o.QueryTable(&Match{}).All(&matchs)
	if err == nil {
		cache.Set(b.matchCacheKey(), matchs, 12*time.Hour)
	}
	return matchs
}

func (b *Bas) Match(id int) *Match {
	matchs := b.Matchs()
	for _, v := range matchs {
		if v.Id == id {
			return &v
		}
	}
	return nil
}

func (b *Bas) AddMatch(match *Match) error {
	o := dbs.NewDefaultOrm()
	nums, err := o.QueryTable(&Match{}).Filter("en", match.En).Count()
	if err != nil {
		return errors.New(err.Error())
	}
	if nums > 0 {
		return errors.New("已存在相同的赛事")
	}
	id, err := o.Insert(match)
	if err == nil && id > 0 {
		cache := utils.GetCache()
		cache.Delete(b.matchCacheKey())
		return nil
	}
	return err
}

func (b *Bas) UpdateMatch(match *Match) error {
	o := dbs.NewDefaultOrm()
	num, err := o.Update(match)
	if err != nil {
		return errors.New(err.Error())
	}
	if num > 0 {
		cache := utils.GetCache()
		cache.Delete(b.matchCacheKey())
		return nil
	}
	return errors.New("不存在对应的赛事")
}

//smiley
func (b *Bas) GetSmileies(category string) []Smiley {
	o := dbs.NewDefaultOrm()
	smileys := []Smiley{}
	o.QueryTable(&Smiley{}).Filter("category", category).All(&smileys)
	return smileys
}

//home ads
func (b *Bas) lastHomeAdCacheKey() string {
	return "mobile_base_last_homead"
}
func (b *Bas) CreateHomeAd(ad *HomeAd) error {
	o := dbs.NewDefaultOrm()
	_, err := o.Insert(ad)
	if err == nil {
		cache := utils.GetCache()
		cache.Delete(b.lastHomeAdCacheKey())
	}
	return err
}

func (b *Bas) DeleteHomeAd(id int64) error {
	o := dbs.NewDefaultOrm()
	c, err := o.QueryTable(&HomeAd{}).Filter("id", id).Delete()
	if c > 0 {
		cache := utils.GetCache()
		cache.Delete(b.lastHomeAdCacheKey())
	}
	return err
}

func (b *Bas) LastNewHomeAd() *HomeAd {
	cache := utils.GetCache()
	ha := HomeAd{}
	err := cache.Get(b.lastHomeAdCacheKey(), &ha)
	if err == nil {
		return &ha
	}
	o := dbs.NewDefaultOrm()
	err = o.QueryTable(&ha).OrderBy("-end_time").Limit(1).One(&ha)
	if err == nil {
		cache.Add(b.lastHomeAdCacheKey(), ha, 24*time.Hour)
		return &ha
	}
	return nil
}
