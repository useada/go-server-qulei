package main

import (
	"context"

	elastic "gopkg.in/olivere/elastic.v5"
)

type ESHandle struct {
	Client *elastic.Client
}

func (e *ESHandle) SearchByName(name string, offset, limit int) ([]SearchItem, error) {
	items := make([]SearchItem, 0)

	queryer := elastic.NewBoolQuery()
	queryer = queryer.Should(elastic.NewQueryStringQuery(name).Field("name"))
	res, err := ES.Client.Search().Index("user").Query(queryer).
		From(offset).Size(limit).Do(context.Background())
	if err != nil {
		return items, err
	}

	if res.Hits.TotalHits == 0 {
		return items, nil
	}

	for _, hit := range res.Hits.Hits {
		item := SearchItem{Source: *hit.Source}
		items = append(items, item)
	}
	return items, err
}

func (e *ESHandle) SearchByNear(lat, lon float64, gender, offset, limit int) ([]SearchItem, error) {
	items := make([]SearchItem, 0)

	queryer := elastic.NewBoolQuery()
	queryer = queryer.Filter(elastic.NewTermQuery("state", 0),
		elastic.NewGeoDistanceQuery("position").Point(lat, lon).Distance("10km"))
	// TODO 这里加上面一条会计算两遍geo吗?
	dist := elastic.NewGeoDistanceSort("position").
		Point(lat, lon).Asc().Unit("km").GeoDistance("plane")

	res, err := ES.Client.Search().Index("user").Query(queryer).
		SortBy(dist).From(offset).Size(limit).Do(context.Background())
	if err != nil {
		return items, err
	}

	if res.Hits.TotalHits == 0 {
		return items, nil
	}

	for _, hit := range res.Hits.Hits {
		item := SearchItem{Source: *hit.Source}
		if len(hit.Sort) != 0 {
			item.Distance = hit.Sort[0].(float64)
		}
		items = append(items, item)
	}
	return items, nil
}

type SearchItem struct {
	Source   []byte
	Distance float64
}

var ES *ESHandle

func InitElasticClient(conf ElasticConfigor) (err error) {
	ES = &ESHandle{}
	ES.Client, err = elastic.NewClient(elastic.SetURL(conf.Hosts...))
	return err
}
