package dbs

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

func OpenFile(path string, o *opt.Options) (*leveldb.DB, error) {
	return leveldb.OpenFile(path, o)
}
