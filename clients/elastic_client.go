package clients

import (
	"BLogger/models"
	"context"
	"fmt"
	"github.com/olivere/elastic"
	"log"
	"strings"
	"sync"
	"time"
)

type EsClient interface {
	CreateBulk(parser chan models.StandardLog)
	ParserToJson(lines chan string, parser chan models.StandardLog)
}

type esClient struct {
	client *elastic.Client
	esHost string
	esIndex string
	esType string
	separator string
	interval int
}

var singletonEsClient EsClient

func NewElasticClient(esHost, esIndex, esType, separator string, interval int) EsClient {
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
		separator,
		interval,
	}
}

func (e *esClient) CreateBulk(parser chan models.StandardLog) {
	bulkRequest := e.client.Bulk()
	mu := sync.Mutex{}
	for item := range parser {
		req := elastic.NewBulkIndexRequest().Index(e.esIndex).Type(e.esType).Doc(item)
		bulkRequest = bulkRequest.Add(req)
		// If the log is to big after 5 seconds will send a bulk and reset it
		time.AfterFunc(time.Duration(e.interval) * time.Second, func() {
			if bulkRequest != nil {
				mu.Lock()
				e.sendToEs(bulkRequest)
				bulkRequest = nil
				mu.Unlock()
			}
		})
	}
	mu.Lock()
	defer mu.Unlock()
	e.sendToEs(bulkRequest)
}

func (e *esClient) ParserToJson(lines chan string, parser chan models.StandardLog) {
	defer close(parser)

	for line := range lines {
		var std models.StandardLog
		stripped := strings.Split(line, fmt.Sprintf(" %s ", e.separator))
		if len(stripped) != 3 {
			log.Println("[Error] invalid log format ignoring line |", line)
			continue
		}
		std.CreatedAt, std.Level, std.Text = stripped[0], stripped[1], stripped[2]
		parser <- std
	}
}

func (e *esClient) sendToEs(bulk *elastic.BulkService) {
	defer func() {
        if err := recover(); err != nil {
            log.Println("connection refused to elasticsearch")
        }
    }()
	response, err := bulk.Do(context.Background())
	if err != nil || response.Errors {
		log.Println("Error sending bulk to elasticsearch -> " + err.Error())
	}
	log.Println("Sending to elasticsearch errors: ", response.Errors)
}