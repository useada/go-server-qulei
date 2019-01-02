package main

import (
	"context"

	elastic "gopkg.in/olivere/elastic.v5"
)

type EsHandle struct {
	Client *elastic.Client
}

var ES = &EsHandle{}

func InitElasticClient(conf ElasticConfigor) error {
	var err error
	ES.Client, err = elastic.NewClient(elastic.SetURL(conf.Hosts...))
	return err
}

func (es *EsHandle) UsersByName(name string,
	offset, limit int) ([]SearchModel, error) {
	items := make([]SearchModel, 0)
	if offset+limit > 10000 {
		return items, nil
	}

	queryer := elastic.NewBoolQuery().
		Should(elastic.NewQueryStringQuery(name).Field("name"))
	res, err := es.Client.Search().Index("user").Query(queryer).
		From(offset).Size(limit).Do(context.Background())
	if err != nil {
		return items, err
	}
	if res.Hits.TotalHits == 0 {
		return items, nil
	}

	for _, hit := range res.Hits.Hits {
		items = append(items, SearchModel{Source: *hit.Source})
	}
	return items, err
}

func (es *EsHandle) UsersByNear(lat, lon float64,
	gender, offset, limit int) ([]SearchModel, error) {
	items := make([]SearchModel, 0)
	if offset+limit > 10000 {
		return items, nil
	}

	queryer := elastic.NewBoolQuery().Filter(elastic.NewTermQuery("state", 0),
		elastic.NewGeoDistanceQuery("geo").Point(lat, lon).Distance("10km"))
	sorter := elastic.NewGeoDistanceSort("geo").
		Point(lat, lon).Asc().Unit("km").SortMode("min").GeoDistance("plane")
	res, err := es.Client.Search().Index("user").Query(queryer).
		SortBy(sorter).From(offset).Size(limit).Do(context.Background())
	if err != nil {
		return items, err
	}
	if res.Hits.TotalHits == 0 {
		return items, nil
	}

	for _, hit := range res.Hits.Hits {
		item := SearchModel{Source: *hit.Source}
		if len(hit.Sort) != 0 {
			item.Distance = hit.Sort[0].(float64)
		}
		items = append(items, item)
	}
	return items, nil
}
