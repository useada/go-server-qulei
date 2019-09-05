package store

import (
	"context"

	"a.com/go-server/service/esearch/model"
	elastic "gopkg.in/olivere/elastic.v5"
)

type es struct {
	*elastic.Client
}

func NewElasticRepo(client *elastic.Client) Store {
	return &es{client}
}

func (e *es) UsersByName(name string, offset, limit int) ([]model.SearchInfo, error) {
	items := make([]model.SearchInfo, 0)
	if offset+limit > 10000 {
		return items, nil
	}

	queryer := elastic.NewBoolQuery().
		Should(elastic.NewQueryStringQuery(name).Field("name"))
	res, err := e.Client.Search().Index("user").Query(queryer).
		From(offset).Size(limit).Do(context.Background())
	if err != nil {
		return items, err
	}
	if res.Hits.TotalHits == 0 {
		return items, nil
	}

	for _, hit := range res.Hits.Hits {
		items = append(items, model.SearchInfo{Source: *hit.Source})
	}
	return items, err
}

func (e *es) UsersByNear(lat, lon float64, gender, offset, limit int) ([]model.SearchInfo, error) {
	items := make([]model.SearchInfo, 0)
	if offset+limit > 10000 {
		return items, nil
	}

	queryer := elastic.NewBoolQuery().Filter(elastic.NewTermQuery("state", 0),
		elastic.NewGeoDistanceQuery("geo").Point(lat, lon).Distance("10km"))
	sorter := elastic.NewGeoDistanceSort("geo").
		Point(lat, lon).Asc().Unit("km").SortMode("min").GeoDistance("plane")
	res, err := e.Client.Search().Index("user").Query(queryer).
		SortBy(sorter).From(offset).Size(limit).Do(context.Background())
	if err != nil {
		return items, err
	}
	if res.Hits.TotalHits == 0 {
		return items, nil
	}

	for _, hit := range res.Hits.Hits {
		item := model.SearchInfo{Source: *hit.Source}
		if len(hit.Sort) != 0 {
			item.Distance = hit.Sort[0].(float64)
		}
		items = append(items, item)
	}
	return items, nil
}

type Config struct {
	Hosts []string
	Auth  string
}

func NewElasticClient(conf Config) *elastic.Client {
	client, err := elastic.NewClient(elastic.SetURL(conf.Hosts...))
	if err != nil {
		panic(err)
	}
	return client
}
