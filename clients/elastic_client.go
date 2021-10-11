package clients

import (
	"BLogger/models"
	"context"
	"github.com/olivere/elastic"
	"log"
	"sync"
	"time"
)

type EsClient interface {
	CreateBulk(parser chan models.StandardLog)
}

type esClient struct {
	client *elastic.Client
	esHost string
	esIndex string
	esType string
	interval int
}

var singletonEsClient EsClient

func NewElasticClient(esHost, esIndex, esType string, interval int) EsClient {
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
		interval,
	}
}

func (e *esClient) CreateBulk(chParser chan models.StandardLog) {
	bulkRequest := e.client.Bulk()
	mu := sync.Mutex{}
	for item := range chParser {
		req := elastic.NewBulkIndexRequest().Index(e.esIndex).Type(e.esType).Doc(item)
		bulkRequest = bulkRequest.Add(req)
		time.AfterFunc(time.Duration(e.interval) * time.Second, func() {
			if bulkRequest != nil {
				mu.Lock()
				e.save(bulkRequest)
				bulkRequest = nil
				mu.Unlock()
			}
		})
	}
	mu.Lock()
	defer mu.Unlock()
	e.save(bulkRequest)
}

func (e *esClient) save(bulk *elastic.BulkService) {
	defer func() {
        if err := recover(); err != nil {
            log.Println("connection refused to elasticsearch")
        }
    }()
	response, err := bulk.Do(context.Background())
	if err != nil || response.Errors {
		log.Println("Error sending bulk to elasticsearch -> " + err.Error())
		return
	}
	log.Println("Sending to elasticsearch errors: ", response.Errors)
}