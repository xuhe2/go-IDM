package main

import (
	"flag"
	"log"

	"github.com/xuhe2/go-IDM/utils"
)

func main() {
	// parse command line arguments
	threads := flag.Int("t", 1, "number of threads")
	path := flag.String("p", "./", "path to save file")
	name := flag.String("n", "", "name of file")
	md5 := flag.String("md5", "", "calculate md5 hash")
	flag.Parse()
	args := flag.Args()
	// check if a url was provided
	if len(args) == 0 {
		panic("no url provided")
	}
	url := args[0]

	log.Printf("Starting download of %s with %d threads\n", url, *threads)
	fd := utils.NewFileDownloader(*name, url, *threads, *path, *md5)
	fd.Download()
}
