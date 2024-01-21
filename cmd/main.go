package main

import (
	"fmt"
	"log"

	"github.com/shoet/go-lambda-docker/internal"
)

func main() {
	fmt.Println("run handler")
	if _, err := internal.CopyBrowser(); err != nil {
		log.Fatalf("could not copy browser: %v", err)
	}
	p := internal.NewPlaywrightClient()
	if err := p.RunScrape(); err != nil {
		log.Fatalf("could not run scrape: %v", err)
	}
}
