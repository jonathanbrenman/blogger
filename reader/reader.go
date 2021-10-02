package reader

import (
	"bufio"
	"os"
)

type Reader interface {
	ReadFile(lines chan string)
}

type readerImpl struct {
	file string
}

func NewReader(file string) Reader {
	return &readerImpl{file}
}

func (r *readerImpl) ReadFile(lines chan string) {
	defer close(lines)

	f, err := os.Open(r.file)
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