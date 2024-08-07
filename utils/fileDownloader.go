package utils

import "strings"

type FileDownloader struct {
	FileName string
	Size     int64
	Url      string
	Threads  int
	Path     string
}

func NewFileDownloader(fileName string, size int64, url string, threads int, path string) *FileDownloader {
	if fileName == "" {
		fileName = GetFileNameFromUrl(url)

	}
	return &FileDownloader{
		FileName: fileName,
		Size:     size,
		Url:      url,
		Threads:  threads,
		Path:     path,
	}
}

// GetFileNameFromUrl extracts the file name from a URL
func GetFileNameFromUrl(url string) string {
	fileName := strings.Split(url, "/")[len(strings.Split(url, "/"))-1]
	return fileName
}
