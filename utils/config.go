package utils

import "net/url"

type Config struct {
	Url   string
	Path  string
	Force bool
	Proxy *url.URL
}

type FileDownloaderConfig struct {
	Config    Config
	FileName  string
	Size      int64
	Threads   int
	FileParts []*FilePart
	MD5       string
}

type FilePartConfig struct {
	Config Config
}
