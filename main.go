package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/olivere/elastic"
	"log"
	"os"
	"strings"
)

// FIXME pass config to yaml file and load configs there
var configs = struct {
	EsHost string
	File      string
	Separator string
	Index string
	Type string
}{
	EsHost: "http://localhost:9200",
	File:      "./log.txt",
	Separator: "-.-.-",
	Index: "app",
	Type: "log",
}

type EsStd struct {
	CreatedAt string `json:"created_at"`
	Level     string `json:"level"`
	Text      string `json:"text"`
}

var es = esClient()

func esClient() *elastic.Client {
	client, err := elastic.NewClient(
		elastic.SetURL(configs.EsHost),
		elastic.SetHealthcheck(false),
		elastic.SetSniff(false))

	if err != nil {
		log.Println("[Error] connecting to elastic - %v\n", err)
	} else {
		log.Println("Connected to elasticsearch OK.")
	}
	return client
}

func readFile(lines chan string, file string) {
	defer close(lines)

	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		lines <- scanner.Text()
	}
}

func parserToJson(lines chan string, parser chan EsStd) {
	defer close(parser)

	for line := range lines {
		var std EsStd
		stripped := strings.Split(line, fmt.Sprintf(" %s ", configs.Separator))
		if len(stripped) != 3 {
			log.Println("[Error] invalid log format ignoring line |", line)
			continue
		}
		std.CreatedAt, std.Level, std.Text = stripped[0], stripped[1], stripped[2]
		parser <- std
	}
}

func createBulk(parser chan EsStd) {
	bulkRequest := es.Bulk()
	for item := range parser {
		req := elastic.NewBulkIndexRequest().Index(configs.Index).Type(configs.Type).Doc(item)
		bulkRequest = bulkRequest.Add(req)
		// FIXME add timeout or buffer to send if the log is to big.
	}
	sendToEs(bulkRequest)
}

func sendToEs(bulk *elastic.BulkService) {
	response, err := bulk.Do(context.Background())
	if err != nil {
		fmt.Println("Error sending bulk to elasticsearch -> " + err.Error())
	}
	log.Println("Sending to elasticsearch errors: ", response.Errors)
}

func main() {
	lines := make(chan string)
	parser := make(chan EsStd)

	go readFile(lines, configs.File)
	go parserToJson(lines, parser)
	createBulk(parser)
}
