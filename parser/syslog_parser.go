package parser

import (
	"BLogger/models"
	"fmt"
)

type syslogParser struct {}

func newSyslogParser() Parser {
	return &syslogParser{}
}

func (s *syslogParser) ToJson(chLines chan string, chParser chan models.StandardLog) {
	defer close(chParser)

	for line := range chLines {
		var std models.StandardLog
		fmt.Println(line)
		// std.CreatedAt, std.Level, std.Text = stripped[0], stripped[1], stripped[2]
		chParser <- std
	}
}