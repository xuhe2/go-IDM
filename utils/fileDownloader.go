package utils

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
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
		filePartConfig := FilePartConfig{
			Config: fileDownloader.Config,
			Index:  i,
			From:   start,
			To:     end,
			Status: "Waiting",
			Data:   bytes.Buffer{},
		}
		fileDownloader.FileParts[i] = NewFilePart(filePartConfig)
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
		// close file part
		defer part.Close()
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
	// create a new writer
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	// write file parts
	if fd.Config.InMemory {
		// write file from memory
		for _, part := range fd.FileParts {
			_, err := writer.Write(part.Data.Bytes())
			if err != nil {
				log.Fatalf(ColorString("Error writing file: %v", Red), err)
			}
		}
	} else {
		// write file from disk
		for _, part := range fd.FileParts {
			// open tmp file
			tmpFile, err := os.Open(part.TmpFileName)
			if err != nil {
				log.Fatalf(ColorString("Error opening tmp file when merge: %v", Red), err)
			}
			defer tmpFile.Close()
			// write tmp file into file
			_, err = io.Copy(writer, tmpFile)
			if err != nil {
				log.Fatalf(ColorString("Error writing file when merge: %v", Red), err)
			}
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
	// do the request
	client := http.Client{}
	// if proxy is set, use it
	if fd.Config.Proxy != nil {
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(fd.Config.Proxy),
		}
	}
	resp, err := client.Do(req) // do the request
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
		return errors.New("Server does not support range requests, use `-f` to force download")
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

func (fd *FileDownloader) ShowProcess() {
	totalSize := fd.Size
	var downloadedSize int64
	for {
		// init downloaded size
		downloadedSize = 0
		// sleep 1 second
		time.Sleep(time.Second)
		// get downloaded size
		for _, part := range fd.FileParts {
			downloadedSize += part.GetSize()
		}
		// show process
		UpdateOutput(fmt.Sprintf("Downloaded: %v / %v", ColorString(Bytes2Size(downloadedSize), Green), Bytes2Size(totalSize)))
	}
}
