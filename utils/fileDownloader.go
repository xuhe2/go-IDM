package utils

import "strings"

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
	if fileName == "" {
		fileName = GetFileNameFromUrl(url)
	}
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

// GetFileNameFromUrl extracts the file name from a URL
func GetFileNameFromUrl(url string) string {
	fileName := strings.Split(url, "/")[len(strings.Split(url, "/"))-1]
	return fileName
}
