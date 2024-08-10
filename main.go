package main

import (
	"flag"
	"log"
	"time"

	"github.com/xuhe2/go-IDM/utils"
)

func main() {
	nowTime := time.Now()

	// parse command line arguments
	threads := flag.Int("t", 1, "number of threads")
	path := flag.String("p", ".", "path to save file")
	name := flag.String("n", "", "name of file")
	md5 := flag.String("md5", "", "calculate md5 hash")
	force := flag.Bool("f", false, "force download")
	flag.Parse()
	args := flag.Args()
	// check if a url was provided
	if len(args) == 0 {
		panic("no url provided")
	}
	url := args[0]

	log.Printf("Starting download of %s with %d threads\n", url, *threads)

	// create a new file downloader config
	fdConfig := utils.FileDownloaderConfig{
		Config: utils.Config{
			Url:   url,
			Path:  *path,
			Force: *force,
		},
		FileName:  *name,
		Size:      0,
		Threads:   *threads,
		FileParts: make([]*utils.FilePart, *threads),
		MD5:       *md5,
	}
	// create a new file downloader
	fd := utils.NewFileDownloader(fdConfig)
	if fd == nil {
		panic("failed to create file downloader")
	}

	// show a loading animation
	loadingStr := "Loading..."
	go func() {
		for {
			time.Sleep(1 * time.Second)
			utils.UpdateOutput(loadingStr)
			loadingStr += "."
		}
	}()

	// download the file
	fd.Download()

	// show the time taken
	log.Printf("Downloaded in %s", utils.ColorString(time.Since(nowTime).String(), utils.Green))
}
