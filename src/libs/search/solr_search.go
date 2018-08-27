package search

import (
	"fmt"
	"strings"

	"github.com/solr"
)

type SolrSearcher struct {
	options *SearchOptions
}

func NewSolrSearcher(opts *SearchOptions) *SolrSearcher {
	searcher := &SolrSearcher{}
	searcher.options = opts
	return searcher
}

func (s *SolrSearcher) getClient(core string) (*solr.Connection, error) {
	return solr.Init(s.options.Host, s.options.Port, core)
}

func (s *SolrSearcher) fqParamsSplice(vals []uint64, exclude bool) string {
	if len(vals) == 0 {
		return ""
	}
	if len(vals) == 1 {
		return fmt.Sprintf("%d", vals[0])
	}
	str := "("
	for _, val := range vals {
		str += fmt.Sprintf("%d OR ", val)
	}
	str = strings.TrimRight(str, " OR ")
	str += ")"
	return str
}

func (s *SolrSearcher) VideoQuery(words string) ([]int64, int, error) {
	solrClient, _ := s.getClient(s.options.IndexName)
	q := solr.Query{
		Rows:  s.options.Limit,
		Start: s.options.Offset,
		Sort:  "post_time desc",
	}
	params := solr.URLParamMap{}
	params["q"] = []string{"title:" + strings.ToLower(words)}
	params["df"] = []string{"title"}
	params["fl"] = []string{"id"}

	fqItems := []string{}
	if len(s.options.Filters) > 0 {
		for _, filter := range s.options.Filters {
			fqHeader := "%s:%s"
			if filter.Exclude {
				fqHeader = "!%s:%s"
			}
			fqItem := fmt.Sprintf(fqHeader, filter.Attr, s.fqParamsSplice(filter.Values, filter.Exclude))
			fqItems = append(fqItems, fqItem)
		}
	}
	if len(s.options.FilterRangeInt) > 0 {
		for _, filter := range s.options.FilterRangeInt {
			if filter.Exclude {
				fqItems = append(fqItems, fmt.Sprintf("!%s:[%d TO %d]", filter.Attr, filter.Min, filter.Max))
			} else {
				fqItems = append(fqItems, fmt.Sprintf("%s:[%d TO %d]", filter.Attr, filter.Min, filter.Max))
			}
		}
	}
	params["fq"] = fqItems
	q.Params = params

	res, err := solrClient.Select(&q)
	ids := []int64{}
	if err != nil {
		return ids, 0, err
	}
	result := res.Results
	total := result.NumFound
	for i := 0; i < result.Len(); i++ {
		sId := result.Get(i).Field("id")
		id := sId.(float64)
		ids = append(ids, int64(id))
	}
	return ids, total, nil
}

func (s *SolrSearcher) MemberQuery(words string) ([]int64, int, error) {
	solrClient, _ := s.getClient(s.options.IndexName)
	q := solr.Query{
		Rows:  s.options.Limit,
		Start: s.options.Offset,
		Sort:  "fans desc",
	}
	params := solr.URLParamMap{}
	params["q"] = []string{"nick_name:" + strings.ToLower(words)}
	params["df"] = []string{"nick_name"}
	params["fl"] = []string{"uid"}

	fqItems := []string{}
	if len(s.options.Filters) > 0 {
		for _, filter := range s.options.Filters {
			fqHeader := "%s:%s"
			if filter.Exclude {
				fqHeader = "!%s:%s"
			}
			fqItem := fmt.Sprintf(fqHeader, filter.Attr, s.fqParamsSplice(filter.Values, filter.Exclude))
			fqItems = append(fqItems, fqItem)
		}
	}
	if len(s.options.FilterRangeInt) > 0 {
		for _, filter := range s.options.FilterRangeInt {
			if filter.Exclude {
				fqItems = append(fqItems, fmt.Sprintf("!%s:[%d TO %d]", filter.Attr, filter.Min, filter.Max))
			} else {
				fqItems = append(fqItems, fmt.Sprintf("%s:[%d TO %d]", filter.Attr, filter.Min, filter.Max))
			}
		}
	}
	params["fq"] = fqItems
	q.Params = params

	res, err := solrClient.Select(&q)
	ids := []int64{}
	if err != nil {
		return ids, 0, err
	}
	result := res.Results
	total := result.NumFound
	for i := 0; i < result.Len(); i++ {
		sId := result.Get(i).Field("uid")
		id := sId.(float64)
		ids = append(ids, int64(id))
	}
	return ids, total, nil
}

func (s *SolrSearcher) LiveQuery(words string) ([]int64, int, error) {
	solrClient, _ := s.getClient(s.options.IndexName)
	q := solr.Query{
		Rows:  s.options.Limit,
		Start: s.options.Offset,
		Sort:  "status asc,onlines desc",
	}
	params := solr.URLParamMap{}
	params["q"] = []string{"keywords:" + strings.ToLower(words)}
	params["df"] = []string{"keywords"}
	params["fl"] = []string{"id"}

	fqItems := []string{}
	if len(s.options.Filters) > 0 {
		for _, filter := range s.options.Filters {
			fqHeader := "%s:%s"
			if filter.Exclude {
				fqHeader = "!%s:%s"
			}
			fqItem := fmt.Sprintf(fqHeader, filter.Attr, s.fqParamsSplice(filter.Values, filter.Exclude))
			fqItems = append(fqItems, fqItem)
		}
	}
	if len(s.options.FilterRangeInt) > 0 {
		for _, filter := range s.options.FilterRangeInt {
			if filter.Exclude {
				fqItems = append(fqItems, fmt.Sprintf("!%s:[%d TO %d]", filter.Attr, filter.Min, filter.Max))
			} else {
				fqItems = append(fqItems, fmt.Sprintf("%s:[%d TO %d]", filter.Attr, filter.Min, filter.Max))
			}
		}
	}
	params["fq"] = fqItems
	q.Params = params

	res, err := solrClient.Select(&q)
	ids := []int64{}
	if err != nil {
		return ids, 0, err
	}
	result := res.Results
	total := result.NumFound
	for i := 0; i < result.Len(); i++ {
		sId := result.Get(i).Field("id")
		id := sId.(float64)
		ids = append(ids, int64(id))
	}
	return ids, total, nil
}

func (s *SolrSearcher) ProgramQuery(words string) ([]int64, int, error) {
	solrClient, _ := s.getClient(s.options.IndexName)
	q := solr.Query{
		Rows:  s.options.Limit,
		Start: s.options.Offset,
		Sort:  "stime ASC,onlines DESC",
	}
	params := solr.URLParamMap{}
	params["q"] = []string{"keywords:" + strings.ToLower(words)}
	params["df"] = []string{"keywords"}
	params["fl"] = []string{"id"}

	fqItems := []string{}
	if len(s.options.Filters) > 0 {
		for _, filter := range s.options.Filters {
			fqHeader := "%s:%s"
			if filter.Exclude {
				fqHeader = "!%s:%s"
			}
			fqItem := fmt.Sprintf(fqHeader, filter.Attr, s.fqParamsSplice(filter.Values, filter.Exclude))
			fqItems = append(fqItems, fqItem)
		}
	}
	if len(s.options.FilterRangeInt) > 0 {
		for _, filter := range s.options.FilterRangeInt {
			if filter.Exclude {
				fqItems = append(fqItems, fmt.Sprintf("!%s:[%d TO %d]", filter.Attr, filter.Min, filter.Max))
			} else {
				fqItems = append(fqItems, fmt.Sprintf("%s:[%d TO %d]", filter.Attr, filter.Min, filter.Max))
			}
		}
	}
	params["fq"] = fqItems
	q.Params = params

	res, err := solrClient.Select(&q)
	ids := []int64{}
	if err != nil {
		return ids, 0, err
	}
	result := res.Results
	total := result.NumFound
	for i := 0; i < result.Len(); i++ {
		sId := result.Get(i).Field("id")
		id := sId.(float64)
		ids = append(ids, int64(id))
	}
	return ids, total, nil
}

//map数组
//[]interface{}{
//			map[string]interface{}{"id": 22, "title": "abc"},
//			map[string]interface{}{"id": 23, "title": "def"},
//			map[string]interface{}{"id": 24, "title": "def"},
//		}
func (s *SolrSearcher) Update(kvs []interface{}) error {
	solrClient, _ := s.getClient(s.options.IndexName)
	f := map[string]interface{}{
		"add": kvs,
	}
	_, err := solrClient.Update(f, true)

	if err != nil {
		return err
	}
	return nil
}

//条件字段+参数
//[]interface{}{
//			map[string]interface{}{"id": 22, "title": "abc"},
//			map[string]interface{}{"id": 23, "title": "def"},
//			map[string]interface{}{"id": 24, "title": "def"},
//		}
func (s *SolrSearcher) Delete(kvs []interface{}) error {
	solrClient, _ := s.getClient(s.options.IndexName)
	f := map[string]interface{}{
		"delete": kvs,
	}
	_, err := solrClient.Update(f, true)

	if err != nil {
		return err
	}
	return nil
}
