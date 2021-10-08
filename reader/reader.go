package reader

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
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
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}()
	f, err := os.Open(r.file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	offset, err := f.Seek(0, io.SeekEnd)
	buffer := make([]byte, 1024, 1024)
	for {
		readBytes, err := f.ReadAt(buffer, offset)
		if err != nil  && err != io.EOF{
			log.Println(err)
			break
		}
		offset += int64(readBytes)
		if readBytes != 0 {
			fmt.Println(string(buffer[:readBytes]))
			lines <- string(buffer[:readBytes])
		}
		time.Sleep(time.Second)
	}
}