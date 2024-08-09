package utils

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type FileDownloader struct {
	FileName  string
	Size      int64
	Url       string
	Threads   int
	Path      string
	FileParts []FilePart
	MD5       string
}

func NewFileDownloader(fileName string, url string, threads int, path string, md5 string) *FileDownloader {
	// if fileName is empty, extract it from the URL
	return &FileDownloader{
		FileName:  fileName,
		Size:      0,
		Url:       url,
		Threads:   threads,
		Path:      path,
		FileParts: make([]FilePart, threads),
		MD5:       md5,
	}
}

func (fd *FileDownloader) GetInfo() error {
	// get file info
	log.Printf("Getting file info from %s\n", fd.Url)
	// create a new request
	header := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.64",
	}
	req := NewHTTPRequest(fd.Url, header)
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
	if resp.Header.Get("Accept-Ranges") != "bytes" {
		return errors.New("Accept-Ranges error" + resp.Header.Get("Accept-Ranges"))
	}
	// get file name
	if fd.FileName == "" {
		fd.FileName = GetFileNameFromUrl(resp)
	}
	log.Printf("File name: %s\n", fd.FileName)
	// get file size
	fd.Size = resp.ContentLength
	log.Printf("File size: %v \n", Bytes2Size(fd.Size))
	return nil
}
