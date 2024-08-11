package utils

import "net/url"

type Config struct {
	Url   string
	Path  string
	Force bool
	Proxy *url.URL
}

type FileDownloaderConfig struct {
	Config
	FileName  string
	Size      int64
	Threads   int
	FileParts []*FilePart
	MD5       string
}

type FilePartConfig struct {
	Config
	Index  int
	From   int64
	To     int64
	Status string
	Data   []byte
	// signal channel for download finish
	// 0 is success, 1 is failed
	processSignal chan int
}
