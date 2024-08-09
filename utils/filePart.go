package utils

import (
	"fmt"
	"log"
	"net/http"
)

type FilePart struct {
	Url   string
	Index int
	From  int64
	To    int64
	Data  []byte
	// signal channel for download finish
	// 0 is success, 1 is failed
	finishSignal chan int
}

func (fp *FilePart) Download() {
	// craete header for http request
	header := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.64",
		"Range":      fmt.Sprintf("bytes=%d-%d", fp.From, fp.To),
	}
	// craete a new request
	req := NewHTTPRequest(fp.Url, header)
	if req == nil {
		fp.finishSignal <- 1 // download failed
		return
	}
	// start download
	log.Printf("Start download part %d", fp.Index)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error downloading part %d: %v", fp.Index, err)
		fp.finishSignal <- 1 // download failed
		return
	}
	defer resp.Body.Close()
	// read response body
	size, err := resp.Body.Read(fp.Data)
	if err != nil {
		log.Printf("Error reading part %d: %v", fp.Index, err)
		fp.finishSignal <- 1 // download failed
		return
	}
	// check response size
	if int64(size) != fp.To-fp.From+1 {
		log.Printf("Error downloading part %d: expected %d bytes, got %d bytes", fp.Index, fp.To-fp.From+1, size)
		fp.finishSignal <- 1 // download failed
		return
	}
}
