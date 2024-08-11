package utils

import (
	"bytes"
	"net/url"
)

type Config struct {
	Url      string
	Path     string
	Force    bool
	Proxy    *url.URL
	InMemory bool
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
	Data   bytes.Buffer
	// tmp file name
	TmpFileName string
}
