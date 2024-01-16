package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type Response events.APIGatewayProxyResponse

func brower() error {
	fmt.Println("set default browser path")
	launcher.DefaultBrowserDir = "/tmp/rod/browser"

	// ヘッドレスブラウザを起動する
	u := launcher.New().Bin("/usr/bin/chromium").NoSandbox(true).MustLaunch()
	browser := rod.New().ControlURL(u)
	if err := browser.Connect(); err != nil {
		return fmt.Errorf("Failed to connect browser: %w", err)
	}
	fmt.Printf("Start browser: %s\n", u)

	// スクレイピング対象のページを指定する
	url := "https://www.google.com"
	p := browser.MustPage(url)
	if err := p.WaitLoad(); err != nil {
		return fmt.Errorf("Failed to wait load: %w", err)
	}

	// ページが完全にロードされるのを待つ
	if err := p.WaitLoad(); err != nil {
		return fmt.Errorf("page.WaitLoad: %w", err)
	}

	html, err := p.HTML()
	if err != nil {
		return fmt.Errorf("page.HTML: %w", err)
	}

	// 各要素のテキストを出力する
	fmt.Println(html)
	return nil
}

func Handler(ctx context.Context) (Response, error) {
	fmt.Println("come handler")
	if err := brower(); err != nil {
		return Response{}, err
	}

	resp := Response{
		StatusCode: 200,
		Body:       "hello world",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
