# Blogger
Collect, parse to json, create bulk and send to elastic search

## Requeriments
* GO 1.17 (required to build)
* blogger.yaml (required to run)

## Settings
* Create blogger.yaml on same location that the binary file <br>

Example:
```yaml
logs:
  separator: "-.-.-" # indicates how you separate the log lines
  files:
    - "./log.txt"

elasticsearch:
  es_host: "http://localhost:9200"
  es_index: "app"
  es_type: "logs"
  interval: "5" # (in seconds)
```

## The logs should be with this format
* created_at
* [log_level]
* text

For example:
``` bash
2021/10/02 13:38:22 -.-.- [info] -.-.- Pepe
2021/10/02 13:38:22 -.-.- [error] -.-.- this is and error code 123
2021/10/02 13:38:22 -.-.- [debug] -.-.- this message is for debug issue
```

## Kibana + ElasticSearch on docker-compose to test it.
``` bash
$ docker-compose up -d
$ go build .
$ ./BLogger # Binary Size: 8.6MB

# Kibana: http://localhost:5601
# ElasticSearch: http://localhost:9200
```