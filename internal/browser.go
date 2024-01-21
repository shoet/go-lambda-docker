package internal

import (
	"fmt"
	"os"
	"path/filepath"

	cp "github.com/otiai10/copy"
	"github.com/playwright-community/playwright-go"
)

type Crawler interface {
	FetchContents(url string) (string, string, error)
}

type PlaywrightClient struct {
	browser playwright.Browser
}

func NewPlaywrightClient() (*PlaywrightClient, func() error, error) {
	browserBaseDir := "/tmp/playwright/browser"
	runOption := &playwright.RunOptions{
		SkipInstallBrowsers: true,
		DriverDirectory:     browserBaseDir,
		Browsers:            []string{"chromium"},
		Verbose:             true,
	}
	if err := playwright.Install(runOption); err != nil {
		return nil, nil, fmt.Errorf("could not install playwright: %v", err)
	}
	pw, err := playwright.Run(runOption)
	if err != nil {
		return nil, nil, fmt.Errorf("could not run playwright: %v", err)
	}

	matches, err := filepath.Glob(filepath.Join(browserBaseDir, "chromium-*"))
	if err != nil {
		return nil, nil, fmt.Errorf("could not find browser: %v", err)
	}

	if len(matches) == 0 {
		return nil, nil, fmt.Errorf("could not find browser")
	}
	browserPath := filepath.Join(matches[0], "chrome-linux", "chrome")

	chromiumOptions := playwright.BrowserTypeLaunchOptions{
		Headless:        playwright.Bool(true),
		ExecutablePath:  playwright.String(browserPath),
		Timeout:         playwright.Float(1200000),
		ChromiumSandbox: playwright.Bool(false),
		Args: []string{
			"--no-sandbox",
			"--single-process",
			"--disable-gpu-sandbox",
		},
	}
	browser, err := pw.Chromium.Launch(chromiumOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("could not launch browser: %v", err)
	}
	closer := func() error {
		if err := browser.Close(); err != nil {
			return fmt.Errorf("could not close browser: %v", err)
		}
		if err := pw.Stop(); err != nil {
			return fmt.Errorf("could not stop playwright: %v", err)
		}
		return nil
	}
	return &PlaywrightClient{
		browser: browser,
	}, closer, nil
}

func CopyBrowser() (string, error) {
	src := "/var/playwright/browser/chromium-1091"
	dst := "/tmp/playwright/browser/chromium-1091"

	if _, err := os.Stat(dst); os.IsNotExist(err) {
		if err := cp.Copy(src, dst); err != nil {
			return "", fmt.Errorf("could not copy browser: %v", err)
		}
	}
	return dst, nil
}

func (p *PlaywrightClient) FetchPage(url string) (playwright.Page, error) {
	page, err := p.browser.NewPage()
	if err != nil {
		return nil, fmt.Errorf("could not create page: %v", err)
	}
	_, err = page.Goto(url)
	if err != nil {
		return nil, fmt.Errorf("could not goto page: %v", err)
	}
	return page, nil
}

func Run() error {
	b, closer, err := NewPlaywrightClient()
	if err != nil {
		return fmt.Errorf("could not create playwright client: %v", err)
	}

	url := "https://google.com"
	page, err := b.FetchPage(url)
	if err != nil {
		return fmt.Errorf("could not run scrape: %v", err)
	}

	txt, err := page.Content()
	if err != nil {
		return fmt.Errorf("could not get content: %v", err)
	}
	fmt.Println(txt[:10])

	if err := closer(); err != nil {
		return fmt.Errorf("could not close playwright client: %v", err)
	}
	return nil
}
