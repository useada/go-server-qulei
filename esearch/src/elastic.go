package main

import (
	"context"

	elastic "gopkg.in/olivere/elastic.v5"
)

type EsHandle struct {
	Client *elastic.Client
}

var ES *EsHandle

func InitElasticClient(conf ElasticConfigor) (err error) {
	ES = &EsHandle{}
	ES.Client, err = elastic.NewClient(elastic.SetURL(conf.Hosts...))
	return err
}

func (es *EsHandle) SearchByName(name string,
	offset, limit int) ([]SearchItem, error) {
	queryer := elastic.NewBoolQuery()
	queryer = queryer.Should(elastic.NewQueryStringQuery(name).Field("name"))

	items := make([]SearchItem, 0)
	res, err := es.Client.Search().Index("user").Query(queryer).
		From(offset).Size(limit).Do(context.Background())
	if err != nil {
		return items, err
	}

	if res.Hits.TotalHits == 0 {
		return items, nil
	}

	for _, hit := range res.Hits.Hits {
		items = append(items, SearchItem{Source: *hit.Source})
	}
	return items, err
}

func (es *EsHandle) SearchByNear(lat, lon float64, gender,
	offset, limit int) ([]SearchItem, error) {
	queryer := elastic.NewBoolQuery()
	queryer = queryer.Filter(elastic.NewTermQuery("state", 0),
		elastic.NewGeoDistanceQuery("geo").Point(lat, lon).Distance("10km"))

	sorter := elastic.NewGeoDistanceSort("geo").
		Point(lat, lon).Asc().Unit("km").SortMode("min").GeoDistance("plane")

	items := make([]SearchItem, 0)
	res, err := es.Client.Search().Index("user").Query(queryer).
		SortBy(sorter).From(offset).Size(limit).Do(context.Background())
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
