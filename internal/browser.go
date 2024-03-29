package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cp "github.com/otiai10/copy"
	"github.com/playwright-community/playwright-go"
)

type Crawler interface {
	FetchContents(url string) (string, string, error)
}

type PlaywrightClient struct {
	browser playwright.Browser
}

type PlaywrightClientConfig struct {
	BrowserLaunchTimeoutSec int
}

var _ Crawler = (*PlaywrightClient)(nil)

func NewPlaywrightClient(config *PlaywrightClientConfig) (*PlaywrightClient, func() error, error) {
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

	var browserLaunchTimeoutSec float64
	if config.BrowserLaunchTimeoutSec != 0 {
		browserLaunchTimeoutSec = float64(config.BrowserLaunchTimeoutSec) * 1000
	}

	chromiumOptions := playwright.BrowserTypeLaunchOptions{
		Headless:        playwright.Bool(true),
		ExecutablePath:  playwright.String(browserPath),
		Timeout:         playwright.Float(float64(browserLaunchTimeoutSec)),
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

func (p *PlaywrightClient) FetchContents(url string) (string, string, error) {
	page, err := p.FetchPage(url)
	if err != nil {
		return "", "", fmt.Errorf("could not fetch page: %v", err)
	}

	titles, err := page.Locator("h1").All()
	if err != nil {
		return "", "", fmt.Errorf("could not get h1: %v", err)
	}
	titleBuilder := strings.Builder{}
	for _, t := range titles {
		text, err := t.TextContent()
		if err != nil {
			return "", "", fmt.Errorf("could not get text content: %v", err)
		}
		if _, err := titleBuilder.WriteString(text); err != nil {
			return "", "", fmt.Errorf("could not write string: %v", err)
		}
	}

	contents, err := page.Locator("p").All()
	if err != nil {
		return "", "", fmt.Errorf("could not get p: %v", err)
	}
	contentBuilder := strings.Builder{}
	for _, c := range contents {
		text, err := c.TextContent()
		if err != nil {
			return "", "", fmt.Errorf("could not get text content: %v", err)
		}
		if _, err := contentBuilder.WriteString(text); err != nil {
			return "", "", fmt.Errorf("could not write string: %v", err)
		}
	}
	return titleBuilder.String(), contentBuilder.String(), nil
}

func Run() error {
	config := &PlaywrightClientConfig{
		BrowserLaunchTimeoutSec: 120,
	}
	b, closer, err := NewPlaywrightClient(config)
	if err != nil {
		return fmt.Errorf("could not create playwright client: %v", err)
	}

	url := "https://news.yahoo.co.jp/articles/0376263c7ac6dfd8b2bcebcec07802f254af2e55"
	title, content, err := b.FetchContents(url)
	if err != nil {
		return fmt.Errorf("could not fetch contents: %v", err)
	}

	fmt.Println("### title")
	fmt.Println(title)

	fmt.Println("### content")
	fmt.Println(content)

	if err := closer(); err != nil {
		return fmt.Errorf("could not close playwright client: %v", err)
	}
	return nil
}
