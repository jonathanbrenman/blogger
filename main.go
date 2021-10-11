package main

import (
	"BLogger/clients"
	"BLogger/configs"
	"BLogger/models"
	"BLogger/parser"
	"BLogger/reader"
	"strconv"
	"sync"
)

func main() {
	config := configs.NewConfig().Load()
	interval, _ := strconv.Atoi(config.Elasticsearch.Interval)

	// Create elastic search client singleton
	esClient := clients.NewElasticClient(
		config.Elasticsearch.EsHost,
		config.Elasticsearch.EsIndex,
		config.Elasticsearch.EsType,
		interval,
	)

	wg := sync.WaitGroup{}
	for _, file := range config.Logs.Files {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			// Initialize channels
			chLines := make(chan string)
			chParser := make(chan models.StandardLog)

			// Read logs
			fReader := reader.NewReader(file.File)

			// Thread process
			go fReader.ReadFile(chLines)
			go parser.New(file.Parser, file.Separator).ToJson(chLines, chParser)
			esClient.CreateBulk(chParser)
			wg.Done()
		}(&wg)
	}
	wg.Wait()
}
