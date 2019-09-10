package store

import (
	elastic "gopkg.in/olivere/elastic.v5"
)

type es struct {
	*elastic.Client
}

func NewElasticRepo(client *elastic.Client) Store {
	return &es{client}
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
