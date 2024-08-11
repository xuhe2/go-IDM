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
		log.Printf(ColorString("Error downloading part %d: %v", Red), fp.Index, err)
		return
	}
	defer resp.Body.Close()
	// download file
	var fileSize int64
	if fp.InMemory {
		fileSize = fp.save2memory(resp)
	} else {
		fileSize = fp.save2disk(resp)
	}
	// check response size
	if fileSize != fp.To-fp.From+1 {
		log.Printf(ColorString("Error downloading part %d: expected %d bytes, got %d bytes", Red), fp.Index, fp.To-fp.From+1, fileSize)
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

// save file part to memory
// return file size
func (fp *FilePart) save2memory(resp *http.Response) int64 {
	var err error
	// read response body
	fileSize, err := io.Copy(&fp.Data, resp.Body)
	if err != nil {
		log.Printf(ColorString("Error downloading part %d: %v", Red), fp.Index, err)
		return 0
	}
	// get file size
	return fileSize
}

// save file part to tmp dir in disk
// return file size
func (fp *FilePart) save2disk(resp *http.Response) int64 {
	// write response body to file
	tmpFile, err := os.CreateTemp("", "go-IDM")
	if err != nil {
		log.Printf(ColorString("Error create tmp file %d: %v", Red), fp.Index, err)
		return 0
	}
	defer tmpFile.Close()
	// store tmp file name
	fp.TmpFileName = tmpFile.Name()
	// craete a bufio.Writer
	writer := bufio.NewWriter(tmpFile)
	// write response body to file
	fileSize, err := io.Copy(writer, resp.Body)
	if err != nil {
		log.Printf(ColorString("Error downloading part %d: %v", Red), fp.Index, err)
		return 0
	}
	// flush bufio.Writer
	if err := writer.Flush(); err != nil {
		log.Printf(ColorString("Error downloading part %d: %v", Red), fp.Index, err)
		return 0
	}
	return fileSize
}

func (fp *FilePart) GetSize() int64 {
	if fp.InMemory {
		// get file size from memory
		return int64(fp.Data.Len())
	} else {
		// get file size from disk
		fileInfo, err := os.Stat(fp.TmpFileName)
		if err != nil {
			log.Printf("Error getting file size: %v", err)
			return 0
		}
		return fileInfo.Size()
	}
}
