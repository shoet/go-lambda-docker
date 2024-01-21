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
	if err := internal.Run(); err != nil {
		log.Fatalf("could not run scrape: %v", err)
	}
}
