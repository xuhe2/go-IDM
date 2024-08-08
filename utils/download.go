package utils

import (
	"log"
)

func Download(name string, url string, threads int, path string, md5 string) {
	fileDownloader := NewFileDownloader(name, url, threads, path, md5)
	err := fileDownloader.GetInfo()
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}
}
