package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type FilePart struct {
	FilePartConfig
}

func NewFilePart(config FilePartConfig) *FilePart {
	return &FilePart{
		FilePartConfig: config,
	}
}

func (fp *FilePart) Download() {
	// craete header for http request
	header := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.64",
		"Range":      fmt.Sprintf("bytes=%d-%d", fp.From, fp.To),
	}
	// craete a new request
	req := NewHTTPRequest(fp.Config, header)
	if req == nil {
		return
	}
	// start download
	fp.Status = "Downloading"
	log.Printf("Start download part %d", fp.Index)
	// create http client
	client := &http.Client{}
	// set proxy if exist
	if fp.Config.Proxy != nil {
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(fp.Config.Proxy),
		}
	}
	// send request
	resp, err := client.Do(req)
	if err != nil {
		log.Print(ColorString(fmt.Sprintf("Error downloading part %d: %v", fp.Index, err), Red))
		return
	}
	defer resp.Body.Close()
	// read response body
	fp.Data, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Print(ColorString(fmt.Sprintf("Error downloading part %d: %v", fp.Index, err), Red))
		return
	}
	// check response size
	if int64(len(fp.Data)) != fp.To-fp.From+1 {
		log.Print(ColorString(fmt.Sprintf("Error downloading part %d: expected %d bytes, got %d bytes", fp.Index, fp.To-fp.From+1, len(fp.Data)), Red))
		return
	}
	fp.Status = "Finished"
}

func (fp *FilePart) GetProcess() int {
	return <-fp.processSignal
}
