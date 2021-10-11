package parser

import (
	"BLogger/helpers"
	"BLogger/models"
	"log"
	"strings"
)

type syslogParser struct {}

func newSyslogParser() Parser {
	return &syslogParser{}
}

func (s *syslogParser) ToJson(chLines chan string, chParser chan models.StandardLog) {
	defer close(chParser)

	for line := range chLines {
		lines := strings.Split(line, "\n")
		for _, l := range lines {
			var std models.StandardLog
			stripped := strings.Split(l, " ")
			if len(stripped) < 5 {
				log.Println("skipping", stripped)
				continue
			}
			std.Level = "info"
			std.CreatedAt = helpers.SyslogDateToString(stripped[0], stripped[1], stripped[2])
			std.Host = stripped[3]
			std.Process = stripped[4]
			std.Text = strings.Join(stripped[5:], " ")
			chParser <- std
		}
	}
}