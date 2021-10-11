package main

import (
	"BLogger/clients"
	"BLogger/configs"
	"BLogger/models"
	"BLogger/parser"
	"BLogger/reader"
	"fmt"
	"github.com/olivere/elastic"
	"strconv"
	"sync"
	"time"
)

func main() {
	config := configs.NewConfig().Load()
	interval, _ := strconv.Atoi(config.Elasticsearch.Interval)

	// Create elastic search client singleton
	esClient := clients.NewElasticClient(
		config.Elasticsearch.EsHost,
		config.Elasticsearch.EsIndex,
		config.Elasticsearch.EsType,
	)

	wg := sync.WaitGroup{}
	for _, file := range config.Logs.Files {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			// Initialize channels
			chLines := make(chan string)
			chParser := make(chan models.StandardLog)
			chBulk := make(chan *elastic.BulkService)

			// Read logs
			fReader := reader.NewReader(file.File)

			// Thread process
			go fReader.ReadFile(chLines)
			go parser.New(file.Parser, file.Separator).ToJson(chLines, chParser)
			go func(ch chan *elastic.BulkService) {
				for now := range time.Tick(time.Duration(interval) * time.Second) {
					fmt.Println("Saving bulk in elastic search", now)
					esClient.Save(<-ch)
				}
			}(chBulk)
			esClient.CreateBulk(chParser, chBulk)
			wg.Done()
		}(&wg)
	}
	wg.Wait()
}
