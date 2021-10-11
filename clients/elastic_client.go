package clients

import (
	"BLogger/models"
	"context"
	"github.com/olivere/elastic"
	"log"
)

type EsClient interface {
	CreateBulk(chParser chan models.StandardLog, chBulk chan *elastic.BulkService)
	Save(bulk *elastic.BulkService)
}

type esClient struct {
	client *elastic.Client
	esHost string
	esIndex string
	esType string
}

var singletonEsClient EsClient

func NewElasticClient(esHost, esIndex, esType string) EsClient {
	if singletonEsClient != nil {
		return singletonEsClient
	}

	client, err := elastic.NewClient(
		elastic.SetURL(esHost),
		elastic.SetHealthcheck(false),
		elastic.SetSniff(false))

	if err != nil {
		log.Printf("[Error] connecting to elastic - %v\n", err)
	}

	return &esClient{
		client,
		esHost,
		esIndex,
		esType,
	}
}

func (e *esClient) CreateBulk(chParser chan models.StandardLog, chBulk chan *elastic.BulkService) {
	bulkRequest := e.client.Bulk()
	for item := range chParser {
		req := elastic.NewBulkIndexRequest().Index(e.esIndex).Type(e.esType).Doc(item)
		bulkRequest = bulkRequest.Add(req)
		chBulk <- bulkRequest
	}
}

func (e *esClient) Save(bulk *elastic.BulkService) {
	defer func() {
        if err := recover(); err != nil {
            log.Println("connection refused to elasticsearch")
        }
    }()
	if bulk == nil {
		return
	}
	response, err := bulk.Do(context.Background())
	if err != nil || response.Errors {
		log.Println("Error sending bulk to elasticsearch -> " + err.Error())
		return
	}
	log.Println("Sending to elasticsearch errors: ", response.Errors)
}