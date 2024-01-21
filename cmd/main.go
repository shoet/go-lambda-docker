package main

import (
	"fmt"
	"log"

	"github.com/shoet/go-lambda-docker/internal"
)

func main() {
	fmt.Println("run handler")
	browserPath, err := internal.CopyBrowser()
	if err != nil {
		log.Fatalf("could not copy browser: %v", err)
	}
	p := internal.NewPlaywrightClient()
	if err := p.RunScrape(browserPath); err != nil {
		log.Fatalf("could not run scrape: %v", err)
	}
}
