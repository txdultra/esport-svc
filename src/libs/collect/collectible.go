package collect

import (
	"labix.org/v2/mgo/bson"
	"time"
)

type Collectible struct {
	Id             bson.ObjectId `bson:"_id"`
	Uid            int64         `bson:"uid"`
	RelId          string        `bson:"rel_id"`
	RelType        string        `bson:"rel_type"`
	PreviewContent string        `bson:"preview_content"`
	PreviewImg     int64         `bson:"preview_img"`
	CreateTime     time.Time     `bson:"create_time"`
}

type CollectCompler interface {
	CompleCollectible(c *Collectible) error
}

var complers map[string]CollectCompler = make(map[string]CollectCompler)

func RegisterCompler(relType string, assemble CollectCompler) {
	if _, ok := complers[relType]; ok {
		return
	}
	complers[relType] = assemble
}

func getCompler(relType string) CollectCompler {
	if m, ok := complers[relType]; ok {
		return m
	}
	return nil
}
