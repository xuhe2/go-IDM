package utils

import (
	"errors"
	"log"
	"mime"
	"net/http"
	"path/filepath"
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
	req := NewHTTPRequest(fd.Url)
	// if req is nil, panic
	if req == nil {
		panic("Error creating request")
	}
	resp, err := http.DefaultClient.Do(req) // do the request
	// if resp is nil, panic
	if err != nil {
		log.Fatalf("Error getting file info: %v", err)
		panic("Error getting file info")
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
	return nil
}

func NewHTTPRequest(url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil) // create a new request
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
		return nil
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.64")
	return req
}

// GetFileNameFromUrl extracts the file name from a URL
func GetFileNameFromUrl(response *http.Response) string {
	// if `Content-Disposition` exist, extract the file name from it
	if contentDisposition := response.Header.Get("Content-Disposition"); contentDisposition != "" {
		_, params, err := mime.ParseMediaType(contentDisposition)
		if err != nil {
			log.Fatalf("Error parsing content disposition: %v", err)
			return ""
		}
		if fileName, ok := params["filename"]; ok {
			return fileName
		}
	}
	// if `Content-Type` exist, extract the file name from it
	fileName := filepath.Base(response.Request.URL.Path) // extract the file name from the URL
	return fileName
}
