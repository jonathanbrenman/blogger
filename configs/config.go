package configs

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config interface {
	Load() *configs
}

type configs struct {
	Logs struct {
		Files     []struct {
			File string `yaml:"file"`
			Parser string `yaml:"parser"`
			Separator string `yaml:"separator"`
		} `yaml:"files"`
	} `yaml:"logs"`
	Elasticsearch struct {
		EsHost   string `yaml:"es_host"`
		EsIndex  string `yaml:"es_index"`
		EsType   string `yaml:"es_type"`
		Interval string `yaml:"interval"`
	} `yaml:"elasticsearch"`
}

func NewConfig() Config {
	return &configs{}
}

func (c *configs) Load() *configs {
	yamlFile, err := ioutil.ReadFile("./blogger.yaml")
    if err != nil {
        log.Fatal("[ Error ] - reading config yaml file. " + err.Error())
    }
    err = yaml.Unmarshal(yamlFile, &c)
    if err != nil || c.validate() != nil {
        log.Fatal("[ Error ] - format not valid for config file blogger.yaml. " + err.Error())
    }
    return c
}

func (c *configs) validate() error {
	if c.Elasticsearch.EsHost == "" || c.Elasticsearch.Interval == "" ||
		c.Elasticsearch.EsIndex == "" || c.Elasticsearch.EsType == "" || len(c.Logs.Files) == 0 {
		return errors.New("missing config fields, please checkout the documentation.")
	}
	return nil
}