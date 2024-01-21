package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/shoet/go-lambda-docker/internal"
)

func Handler(ctx context.Context) (string, error) {
	if _, err := internal.CopyBrowser(); err != nil {
		log.Fatalf("could not copy browser: %v", err)
	}

	fmt.Println("come handler")
	p := internal.NewPlaywrightClient()
	if err := p.RunScrape(); err != nil {
		return "failed", fmt.Errorf("could not run scrape: %v", err)
	}

	return "success", nil
}

var browserPath string

func main() {
	lambda.Start(Handler)
}
