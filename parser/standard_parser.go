package parser

import (
	"BLogger/models"
	"fmt"
	"log"
	"strings"
)

type stdParser struct {
	separator string
}

func newStdParser(separator string) Parser {
	return &stdParser{separator}
}

func (s *stdParser) ToJson(chLines chan string, chParser chan models.StandardLog) {
	defer close(chParser)

	for line := range chLines {
		var std models.StandardLog
		stripped := strings.Split(line, fmt.Sprintf(" %s ", s.separator))
		if len(stripped) != 3 {
			log.Println("[Error] invalid log format ignoring line |", line)
			continue
		}
		std.CreatedAt, std.Level, std.Text = stripped[0], stripped[1], stripped[2]
		chParser <- std
	}
}