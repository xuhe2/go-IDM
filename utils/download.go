package utils

// this file contains the pure function

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"path/filepath"
)

func Download(name string, url string, threads int, path string, md5 string) {
	fileDownloader := NewFileDownloader(name, url, threads, path, md5)
	err := fileDownloader.GetInfo()
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}
}

func NewHTTPRequest(url string, header map[string]string) *http.Request {
	req, err := http.NewRequest("GET", url, nil) // create a new request
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil
	}
	for key, value := range header {
		req.Header.Set(key, value)
	}
	return req
}

// GetFileNameFromUrl extracts the file name from a URL
func GetFileNameFromUrl(response *http.Response) string {
	// if `Content-Disposition` exist, extract the file name from it
	if contentDisposition := response.Header.Get("Content-Disposition"); contentDisposition != "" {
		_, params, err := mime.ParseMediaType(contentDisposition)
		if err != nil {
			log.Printf("Error parsing content disposition: %v", err)
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

// convert bytes to human readable format
func Bytes2Size(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
)

// string with green color
func ColorString(s string, color string) string {
	return color + s + Reset
}
