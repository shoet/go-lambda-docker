package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/playwright-community/playwright-go"
)

type Response events.APIGatewayProxyResponse

func brower() error {

	// ヘッドレスブラウザを起動する
	// u := launcher.New().Bin("/usr/bin/chromium").NoSandbox(true).MustLaunch()

	// // ヘッドレスブラウザを起動する
	// // u := launcher.New().Bin("/usr/bin/chromium").WorkingDir("/tmp").NoSandbox(true).MustLaunch()
	// u := launcher.New().WorkingDir("/tmp").NoSandbox(true).MustLaunch()

	// browser := rod.New().ControlURL(u)
	// if err := browser.Connect(); err != nil {
	// 	return fmt.Errorf("Failed to connect browser: %w", err)
	// }
	// fmt.Printf("Start browser: %s\n", u)

	fmt.Println("set default browser path")
	launcher.DefaultBrowserDir = "/tmp/rod/browser"
	// u := launcher.New().NoSandbox(true).MustLaunch()
	// u := launcher.New().Bin("/usr/bin/chromium").NoSandbox(true).MustLaunch()
	u := launcher.New().Bin("/usr/bin/chromium-browser").NoSandbox(true).MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

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

func browser2() error {
	runOption := &playwright.RunOptions{
		SkipInstallBrowsers: true,
	}
	if err := playwright.Install(runOption); err != nil {
		log.Fatalf("could not install playwright: %v", err)
	}
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not launch playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch()
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	if _, err = page.Goto("https://google.com"); err != nil {
		log.Fatalf("could not goto: %v", err)
	}
	return nil
}

func main() {
	fmt.Println("run handler")
	if err := browser2(); err != nil {
		fmt.Println("#### error")
		fmt.Println(err)
	}
}
