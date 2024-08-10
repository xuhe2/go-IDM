package utils

type Config struct {
	Url   string
	Path  string
	Force bool
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
