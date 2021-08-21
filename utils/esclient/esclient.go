package esclient

import (
	"context"
	"site/config"

	"github.com/olivere/elastic/v7"
)

func NewClient(c context.Context) (*elastic.Client, error) {

	esEndpoint, err := config.GetConfig(c, config.ElasticSearchEndpoint)
	if err != nil {
		return nil, err
	}

	esUsername, err := config.GetConfig(c, config.ElasticSearchUsername)
	if err != nil {
		return nil, err
	}

	esPassword, err := config.GetConfig(c, config.ElasticSearchPassword)
	if err != nil {
		return nil, err
	}

	return elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(esEndpoint.Value), elastic.SetBasicAuth(esUsername.Value, esPassword.Value))
}
