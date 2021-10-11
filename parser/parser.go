package parser

import (
	"BLogger/models"
	"log"
)

type Parser interface {
	ToJson(lines chan string, parser chan models.StandardLog)
}

func New(pType, separator string) Parser {
	// Here we add new parsers strategies.
	pMap := make(map[string]Parser)
	pMap["std"] = newStdParser(separator)
	pMap["syslog"] = newStdParser(separator)

	if _, ok := pMap[pType]; !ok {
		 log.Fatal("invalid parser type: " + pType)
	}
	return pMap[pType]
}