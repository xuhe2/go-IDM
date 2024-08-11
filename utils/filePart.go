package utils

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	BUFFER_SIZE = 1024 * 1024
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
	// download file
	var fileSize int64
	if fp.InMemory {
		// read response body
		fp.Data, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Print(ColorString(fmt.Sprintf("Error downloading part %d: %v", fp.Index, err), Red))
			return
		}
		// get file size
		fileSize = int64(len(fp.Data))
	} else {
		// write response body to file
		tmpFile, err := os.CreateTemp("", "go-IDM")
		if err != nil {
			log.Print(ColorString(fmt.Sprintf("Error create tmp file %d: %v", fp.Index, err), Red))
			return
		}
		defer tmpFile.Close()
		// store tmp file name
		fp.TmpFileName = tmpFile.Name()
		// craete a bufio.Writer
		writer := bufio.NewWriter(tmpFile)
		// write response body to file
		if fileSize, err = io.Copy(writer, resp.Body); err != nil {
			log.Print(ColorString(fmt.Sprintf("Error downloading part %d: %v", fp.Index, err), Red))
			return
		}
		// flush bufio.Writer
		if err := writer.Flush(); err != nil {
			log.Print(ColorString(fmt.Sprintf("Error flush writer %d: %v", fp.Index, err), Red))
			return
		}
	}
	// check response size
	if fileSize != fp.To-fp.From+1 {
		log.Print(ColorString(fmt.Sprintf("Error downloading part %d: expected %d bytes, got %d bytes", fp.Index, fp.To-fp.From+1, len(fp.Data)), Red))
		return
	}
	fp.Status = "Finished"
}

func (fp *FilePart) Close() {
	// clear tmp file
	if fp.TmpFileName != "" {
		os.Remove(fp.TmpFileName)
	}
}

func (fp *FilePart) GetProcess() int {
	return <-fp.processSignal
}
