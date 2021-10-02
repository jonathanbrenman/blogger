package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var configs = struct {
	File      string
	Separator string
}{
	File:      "./log.txt",
	Separator: "-.-.-",
}

type EsStd struct {
	CreatedAt string
	Level     string
	Text      string
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

func parserToJson(lines chan string) {
	for line := range lines {
		var std EsStd
		stripped := strings.Split(line, fmt.Sprintf(" %s ", configs.Separator))
		if len(stripped) != 3 {
			log.Println("[Error] invalid log format ignoring line |", line)
			continue
		}
		std.CreatedAt, std.Level, std.Text = stripped[0], stripped[1], stripped[2]
		fmt.Println(std)
		//fmt.Println(line)
	}
}

func main() {
	lines := make(chan string)
	go readFile(lines, configs.File)
	parserToJson(lines)
}
