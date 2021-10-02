package main

import (
	"BLogger/clients"
	"BLogger/configs"
	"BLogger/models"
	"BLogger/reader"
	"strconv"
	"sync"
)

func main() {
	// TODO daemon true / false
	// Load configs from blogger.yaml (same dir)
	config := configs.NewConfig().Load()
	interval, _ := strconv.Atoi(config.Elasticsearch.Interval)

	// Create elastic search client singleton
	esClient := clients.NewElasticClient(
		config.Elasticsearch.EsHost,
		config.Elasticsearch.EsIndex,
		config.Elasticsearch.EsType,
		config.Logs.Separator,
		interval,
	)

	wg := sync.WaitGroup{}
	for _, file := range config.Logs.Files {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			// Initialize channels
			lines := make(chan string)
			parser := make(chan models.StandardLog)

			// Read logs
			fReader := reader.NewReader(file)

			// Thread process
			go fReader.ReadFile(lines)
			go esClient.ParserToJson(lines, parser)
			esClient.CreateBulk(parser)
			wg.Done()
		}(&wg)
		wg.Wait()
	}
}
