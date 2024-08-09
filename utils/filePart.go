package utils

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type FilePart struct {
	Url   string
	Index int
	From  int64
	To    int64
	Data  []byte
	// signal channel for download finish
	// 0 is success, 1 is failed
	processSignal chan int
}

func (fp *FilePart) Download(wg *sync.WaitGroup) {
	defer wg.Done()
	// craete header for http request
	header := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.64",
		"Range":      fmt.Sprintf("bytes=%d-%d", fp.From, fp.To),
	}
	// craete a new request
	req := NewHTTPRequest(fp.Url, header)
	if req == nil {
		return
	}
	// start download
	log.Printf("Start download part %d", fp.Index)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print(ColorString(fmt.Sprintf("Error downloading part %d: %v", fp.Index, err), Red))
		return
	}
	defer resp.Body.Close()
	// read response body
	size, err := resp.Body.Read(fp.Data)
	if err != nil {
		log.Print(ColorString(fmt.Sprintf("Error downloading part %d: %v", fp.Index, err), Red))
		return
	}
	// check response size
	if int64(size) != fp.To-fp.From+1 {
		log.Print(ColorString(fmt.Sprintf("Error downloading part %d: expected %d bytes, got %d bytes", fp.Index, fp.To-fp.From+1, size), Red))
		return
	}
}

func (fp *FilePart) GetProcess() int {
	return <-fp.processSignal
}
