package search

import (
	"errors"
	"fmt"
	"strings"

	"github.com/huichen/sego"
	"github.com/yunge/sphinx"
)

const (
	vod_idx_name     = "vods_idx"
	member_idx_name  = "member_idx"
	live_idx_name    = "live_idx"
	program_idx_name = "program_idx"
)

type SearchOptions struct {
	Host             string
	Port             int
	Timeout          int
	Offset           int            // how many records to seek from result-set start
	Limit            int            // how many records to return from result-set starting at offset (default is 20)
	MaxMatches       int            // max matches to retrieve
	MatchMode        int            // query matching mode (default is SPH_MATCH_ALL)
	RankMode         int            // SPH_RANK_PROXIMITY_BM25
	SortMode         int            // SPH_SORT_EXTENDED
	SortBy           string         // @relevance DESC,u1 ASC,@id DESC(排序模式：@relevance和@id是内置变量，@relevance代表相关度权值，@id等于search_id，u1为字段名）
	FieldWeights     map[string]int //
	Filters          []SearchFilter
	FilterRangeInt   []FilterRangeInt //u1:0,100,false;u2:50,90,true（整数范围过滤器：字段u1 >= 0并且u1 <= 100，字段u2 < 50并且u2 > 90）
	FilterRangeFloat FilterRangeFloat //u1:1.23,99.645,false;u2:1034.3,7834.56,true（浮点数范围过滤器：字段u1 >= 1.23并且u1 <= 99.645，字段u2 < 1034.3并且u2 > 7834.56）
	Excerpt          SearchExcerpt    //高亮设置
}

type SearchExcerpt struct {
	Excerpts       int      //（是否开启高亮显示与文本摘要，1开启 或 0关闭）
	ExcerptsBefore string   //<font color=red>  （高亮显示与文本摘要，如果为空值则不进行高亮显示与文本摘要。在匹配的关键字前面插入的字符串。）
	ExcerptsAfter  string   //</font>  （高亮显示与文本摘要，如果为空值则不进行高亮显示与文本摘要。在匹配的关键字之后插入的字符串。）
	ExcerptsLimit  int      //（高亮显示与文本摘要，如果为空值则不进行高亮显示与文本摘要。摘要最多包含的符号（码点）数。）
	ExcerptsFields []string //（仅对指定的字段进行高亮显示，其余字段不进行高亮显示，如果此参数为空，则默认所有的字符型字段都进行高亮显示）
}

type FilterRangeInt struct {
	Attr    string
	Min     uint64
	Max     uint64
	Exclude bool
}

type FilterRangeFloat struct {
	Attr string
	Min  float64
	Max  float64
}

type SearchFilter struct {
	Attr    string
	Values  []uint64
	Exclude bool
}

type Searcher struct {
	options *SearchOptions
}

func NewSearcher(opts *SearchOptions) *Searcher {
	searcher := &Searcher{}
	searcher.options = opts
	return searcher
}

func (f *Searcher) Segment(text string) string {
	// 分词
	byts := []byte(text)
	segments := segmenter.Segment(byts)

	// 处理分词结果
	return sego.SegmentsToStringV2(segments, false)
}

func (f *Searcher) sphinxOptions() *sphinx.Options {
	return &sphinx.Options{
		Host:       f.options.Host,
		Port:       f.options.Port,
		Timeout:    f.options.Timeout,
		Offset:     f.options.Offset,     // how many records to seek from result-set start
		Limit:      f.options.Limit,      // how many records to return from result-set starting at offset (default is 20)
		MaxMatches: f.options.MaxMatches, // max matches to retrieve
		MatchMode:  f.options.MatchMode,  // query matching mode (default is SPH_MATCH_ALL)
		RankMode:   f.options.RankMode,
	}
}

func (f *Searcher) sph_match_mode(match_mode string) int {
	//all,any,phrase,boolean,extended,fullscan,extended2
	switch strings.ToLower(match_mode) {
	case "any":
		return sphinx.SPH_MATCH_ANY
	case "phrase":
		return sphinx.SPH_MATCH_PHRASE
	case "boolean":
		return sphinx.SPH_MATCH_BOOLEAN
	case "extended":
		return sphinx.SPH_MATCH_EXTENDED
	case "fullscan":
		return sphinx.SPH_MATCH_FULLSCAN
	case "extended2":
		return sphinx.SPH_MATCH_EXTENDED2
	default:
		return sphinx.SPH_MATCH_ALL
	}
}

func (f *Searcher) VideoQuery(words string, match_mode string) ([]int64, int, error) {
	sc := sphinx.NewClient(f.sphinxOptions())
	if err := sc.Error(); err != nil {
		return nil, 0, errors.New("connect error:" + err.Error())
	}
	defer sc.Close()
	sph_match := f.sph_match_mode(match_mode)
	sc.SetMatchMode(sph_match)
	sc.SetRankingMode(sphinx.SPH_RANK_SPH04)
	sc.SetSortMode(sphinx.SPH_SORT_EXTENDED, "@rank DESC,@weight DESC,post_time DESC")
	if len(f.options.Filters) > 0 {
		for _, filter := range f.options.Filters {
			sc.SetFilter(filter.Attr, filter.Values, filter.Exclude)
		}
	}
	if len(f.options.FilterRangeInt) > 0 {
		for _, filter := range f.options.FilterRangeInt {
			sc.SetFilterRange(filter.Attr, filter.Min, filter.Max, filter.Exclude)
		}
	}
	words = strings.ToUpper(words) //统一转成大写
	sego_words := words            //f.Segment(words)
	res, err := sc.Query(sego_words, vod_idx_name, "")
	if err != nil {
		fmt.Println(err)
		return nil, 0, errors.New("query fail:" + err.Error())
	}
	var ids []int64
	for _, r := range res.Matches {
		ids = append(ids, int64(r.DocId))
	}
	return ids, res.TotalFound, nil
}

func (f *Searcher) MemberQuery(words string, match_mode string, sortBy string) ([]int64, int, error) {
	sc := sphinx.NewClient(f.sphinxOptions())
	if err := sc.Error(); err != nil {
		return nil, 0, errors.New("connect error:" + err.Error())
	}
	sph_match := f.sph_match_mode(match_mode)
	sc.SetMatchMode(sph_match)
	sc.SetRankingMode(sphinx.SPH_RANK_SPH04)
	if len(sortBy) == 0 {
		sc.SetSortMode(sphinx.SPH_SORT_EXTENDED, "@rank DESC,@weight DESC,fans DESC")
	} else {
		sc.SetSortMode(sphinx.SPH_SORT_EXTENDED, sortBy)
	}
	if len(f.options.Filters) > 0 {
		for _, filter := range f.options.Filters {
			sc.SetFilter(filter.Attr, filter.Values, filter.Exclude)
		}
	}
	if len(f.options.FilterRangeInt) > 0 {
		for _, filter := range f.options.FilterRangeInt {
			sc.SetFilterRange(filter.Attr, filter.Min, filter.Max, filter.Exclude)
		}
	}
	words = strings.ToUpper(words) //统一转成大写
	sego_words := words            //f.Segment(words)
	res, err := sc.Query(sego_words, member_idx_name, "")
	if err != nil {
		return nil, 0, errors.New("query fail:" + err.Error())
	}
	var ids []int64
	for _, r := range res.Matches {
		ids = append(ids, int64(r.DocId))
	}
	return ids, res.TotalFound, nil
}

func (f *Searcher) LiveQuery(words string, match_mode string) ([]int64, int, error) {
	sc := sphinx.NewClient(f.sphinxOptions())
	if err := sc.Error(); err != nil {
		return nil, 0, errors.New("connect error:" + err.Error())
	}
	sph_match := f.sph_match_mode(match_mode)
	sc.SetMatchMode(sph_match)
	sc.SetRankingMode(sphinx.SPH_RANK_SPH04)
	sc.SetSortMode(sphinx.SPH_SORT_EXTENDED, "@rank DESC,@weight DESC,status ASC,onlines DESC")
	if len(f.options.Filters) > 0 {
		for _, filter := range f.options.Filters {
			sc.SetFilter(filter.Attr, filter.Values, filter.Exclude)
		}
	}
	if len(f.options.FilterRangeInt) > 0 {
		for _, filter := range f.options.FilterRangeInt {
			sc.SetFilterRange(filter.Attr, filter.Min, filter.Max, filter.Exclude)
		}
	}
	words = strings.ToUpper(words) //统一转成大写
	sego_words := words            //f.Segment(words)
	res, err := sc.Query(sego_words, live_idx_name, "")
	if err != nil {
		return nil, 0, errors.New("query fail:" + err.Error())
	}
	var ids []int64
	for _, r := range res.Matches {
		ids = append(ids, int64(r.DocId))
	}
	return ids, res.TotalFound, nil
}

func (f *Searcher) ProgramQuery(words string, match_mode string) ([]int64, int, error) {
	sc := sphinx.NewClient(f.sphinxOptions())
	if err := sc.Error(); err != nil {
		return nil, 0, errors.New("connect error:" + err.Error())
	}
	sph_match := f.sph_match_mode(match_mode)
	sc.SetMatchMode(sph_match)
	sc.SetRankingMode(sphinx.SPH_RANK_SPH04)
	sc.SetSortMode(sphinx.SPH_SORT_EXTENDED, "@weight DESC,stime ASC,onlines DESC")
	if len(f.options.Filters) > 0 {
		for _, filter := range f.options.Filters {
			sc.SetFilter(filter.Attr, filter.Values, filter.Exclude)
		}
	}
	if len(f.options.FilterRangeInt) > 0 {
		for _, filter := range f.options.FilterRangeInt {
			sc.SetFilterRange(filter.Attr, filter.Min, filter.Max, filter.Exclude)
		}
	}
	words = strings.ToUpper(words) //统一转成大写
	sego_words := words            //f.Segment(words)
	res, err := sc.Query(sego_words, program_idx_name, "")
	if err != nil {
		return nil, 0, errors.New("query fail:" + err.Error())
	}
	var ids []int64
	for _, r := range res.Matches {
		ids = append(ids, int64(r.DocId))
	}
	return ids, res.TotalFound, nil
}

func (f *Searcher) UpdateAttributes(index string, attrs []string, values [][]interface{}) (ndocs int, err error) {
	sc := sphinx.NewClient(f.sphinxOptions())
	if err := sc.Error(); err != nil {
		return 0, errors.New("connect error:" + err.Error())
	}
	ndocs, err = sc.UpdateAttributes(index, attrs, values, false)
	return
}

func (f *Searcher) Query(words string, sorts []string, index string, match_mode string) ([]int64, int, error) {
	sc := sphinx.NewClient(f.sphinxOptions())
	if err := sc.Error(); err != nil {
		return nil, 0, errors.New("connect error:" + err.Error())
	}
	sph_match := f.sph_match_mode(match_mode)
	sc.SetMatchMode(sph_match)
	sc.SetRankingMode(sphinx.SPH_RANK_SPH04)
	if len(sorts) > 0 {
		sc.SetSortMode(sphinx.SPH_SORT_EXTENDED, strings.Join(sorts, ","))
	} else {
		sc.SetSortMode(sphinx.SPH_SORT_EXTENDED, "@weight DESC")
	}
	if len(f.options.Filters) > 0 {
		for _, filter := range f.options.Filters {
			sc.SetFilter(filter.Attr, filter.Values, filter.Exclude)
		}
	}
	if len(f.options.FilterRangeInt) > 0 {
		for _, filter := range f.options.FilterRangeInt {
			sc.SetFilterRange(filter.Attr, filter.Min, filter.Max, filter.Exclude)
		}
	}
	words = strings.ToUpper(words) //统一转成大写
	sego_words := words            //f.Segment(words)
	res, err := sc.Query(sego_words, index, "")
	if err != nil {
		return nil, 0, errors.New("query fail:" + err.Error())
	}
	var ids []int64
	for _, r := range res.Matches {
		ids = append(ids, int64(r.DocId))
	}
	return ids, res.TotalFound, nil
}
