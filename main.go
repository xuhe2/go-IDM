package main

import (
	"flag"
	"log"
)

func main() {
	// parse command line arguments
	threads := flag.Int("t", 1, "number of threads")
	flag.Parse()
	args := flag.Args()
	// check if a url was provided
	if len(args) == 0 {
		panic("no url provided")
	}
	url := args[0]

	log.Printf("URL: %s", url)
	log.Printf("Threads: %d", *threads)
}
