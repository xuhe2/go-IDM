package utils

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

type FileDownloader struct {
	FileDownloaderConfig
}

func NewFileDownloader(config FileDownloaderConfig) *FileDownloader {
	// if fileName is empty, extract it from the URL
	fileDownloader := &FileDownloader{FileDownloaderConfig: config}
	// get file info
	err := fileDownloader.GetInfo()
	if err != nil {
		log.Printf("Error: %s", err)
		return nil
	}
	// set file parts
	eachSize := fileDownloader.Size / int64(fileDownloader.Threads)
	for i := 0; i < fileDownloader.Threads; i++ {
		start := int64(i) * eachSize
		end := start + eachSize - 1
		// the last part
		if i == fileDownloader.Threads-1 {
			end = fileDownloader.Size - 1
		}
		// new file part
		fileDownloader.FileParts[i] = NewFilePart(fileDownloader.Config, i, start, end)
	}
	return fileDownloader
}

func (fd *FileDownloader) Download() {
	wg := sync.WaitGroup{}
	for _, part := range fd.FileParts {
		wg.Add(1)
		go func() {
			part.Download()
			wg.Done()
		}()
	}
	wg.Wait()
	// finish
	for _, part := range fd.FileParts {
		if part.Status != "Finished" {
			log.Printf("Error Status: %s", part.Status)
			return
		}
	}
	log.Print(ColorString("Download finished", Green))
	// merge and write into file
	fd.MergeAndWrite()
}

// merge the file parts
func (fd *FileDownloader) MergeAndWrite() {
	// open file
	file, err := os.Create(fd.Config.Path + "/" + fd.FileName)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		return
	}
	defer file.Close()
	// write file parts
	for _, part := range fd.FileParts {
		_, err := file.Write(part.Data)
		if err != nil {
			log.Printf("Error writing file: %v", err)
			return
		}
	}
}

func (fd *FileDownloader) GetInfo() error {
	// get file info
	log.Printf("Getting file info from %s\n", fd.Config.Url)
	// create a new request
	header := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.64",
	}
	req := NewHTTPRequest(fd.Config, header)
	// if req is nil, panic
	if req == nil {
		panic("Error creating request")
	}
	resp, err := http.DefaultClient.Do(req) // do the request
	// if resp is nil, panic
	if err != nil {
		panic(fmt.Sprintf("Error getting file info: %v", err))
	}
	// get info from response header
	// check status code
	if resp.StatusCode >= 300 {
		return errors.New("Status code error: " + resp.Status)
	}
	// check Accept-Ranges header
	if !fd.Config.Force && resp.Header.Get("Accept-Ranges") != "bytes" {
		return errors.New("Accept-Ranges error" + resp.Header.Get("Accept-Ranges"))
	}
	// get file name
	if fd.FileName == "" {
		fd.FileName = GetFileNameFromUrl(resp)
	}
	log.Printf("File name: %s\n", ColorString(fd.FileName, Green))
	// get file size
	fd.Size = resp.ContentLength
	log.Printf("File size: %v \n", ColorString(Bytes2Size(fd.Size), Green))
	return nil
}
