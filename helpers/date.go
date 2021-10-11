package helpers

import (
	"fmt"
	"time"
)

var (
	months = map[string]string{
		"Jan":  "01",
		"Feb":  "02",
		"Mar":  "03",
		"Apr":  "04",
		"May":  "05",
		"June": "06",
		"July": "07",
		"Aug":  "08",
		"Sept": "09",
		"Oct":  "10",
		"Nov":  "11",
		"Dec":  "12",
	}
)

func SyslogDateToString(month, day, times string) string {
	year, _, _ := time.Now().Date()
	return fmt.Sprintf("%d-%s-%s %s", year, months[month], day, times)
}
