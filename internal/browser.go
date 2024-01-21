package internal

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	cp "github.com/otiai10/copy"
	"github.com/playwright-community/playwright-go"
)

type PlaywrightClient struct {
	browserBaseDir string
}

func NewPlaywrightClient() *PlaywrightClient {
	return &PlaywrightClient{
		browserBaseDir: "/tmp/playwright/browser",
	}
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

func (p *PlaywrightClient) RunScrape() error {
	runOption := &playwright.RunOptions{
		SkipInstallBrowsers: true,
		DriverDirectory:     p.browserBaseDir,
		Browsers:            []string{"chromium"},
		Verbose:             true,
	}
	fmt.Println("set default browser path")
	fmt.Println(runOption.DriverDirectory)
	if err := playwright.Install(runOption); err != nil {
		log.Fatalf("could not install playwright: %v", err)
	}
	ent, err := os.ReadDir(p.browserBaseDir)
	if err != nil {
		log.Fatalf("could not read dir: %v", err)
	}
	fmt.Println("### entries")
	for _, e := range ent {
		fmt.Println(e.Name())
	}
	pw, err := playwright.Run(runOption)
	if err != nil {
		log.Fatalf("could not launch playwright: %v", err)
	}
	defer pw.Stop()

	matches, err := filepath.Glob(filepath.Join(p.browserBaseDir, "chromium-*"))
	if err != nil {
		log.Fatalf("could not glob: %v", err)
	}
	fmt.Println("### matches")
	fmt.Println(matches)

	browserPath := matches[0]
	browserPath = filepath.Join(browserPath, "chrome-linux", "chrome")

	chromiumOptions := playwright.BrowserTypeLaunchOptions{
		Headless:        playwright.Bool(true),
		ExecutablePath:  playwright.String(browserPath),
		Timeout:         playwright.Float(1200000),
		ChromiumSandbox: playwright.Bool(false),
		Args: []string{
			"--autoplay-policy=user-gesture-required",
			"--disable-background-networking",
			"--disable-background-timer-throttling",
			"--disable-backgrounding-occluded-windows",
			"--disable-breakpad",
			"--disable-client-side-phishing-detection",
			"--disable-component-update",
			"--disable-default-apps",
			"--disable-dev-shm-usage",
			"--disable-domain-reliability",
			"--disable-extensions",
			"--disable-features=AudioServiceOutOfProcess",
			"--disable-hang-monitor",
			"--disable-ipc-flooding-protection",
			"--disable-notifications",
			"--disable-offer-store-unmasked-wallet-cards",
			"--disable-popup-blocking",
			"--disable-print-preview",
			"--disable-prompt-on-repost",
			"--disable-renderer-backgrounding",
			"--disable-setuid-sandbox",
			"--disable-speech-api",
			"--disable-sync",
			"--disk-cache-size=33554432",
			"--hide-scrollbars",
			"--ignore-gpu-blacklist",
			"--metrics-recording-only",
			"--mute-audio",
			"--no-default-browser-check",
			"--no-first-run",
			"--no-pings",
			"--no-sandbox",
			"--no-zygote",
			"--password-store=basic",
			"--use-gl=swiftshader",
			"--use-mock-keychain",
			"--single-process",
			"--disable-gpu-sandbox",
		},
	}
	browser, err := pw.Chromium.Launch(chromiumOptions)
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	defer browser.Close()
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	// if _, err = page.Goto("https://google.com"); err != nil {
	if _, err = page.Goto("https://news.yahoo.co.jp/articles/89b9e71a181813e422f7183d3194adb0a80ddb5f"); err != nil {
		log.Fatalf("could not goto: %v", err)
	}
	txt, err := page.Content()
	if err != nil {
		log.Fatalf("could not get content: %v", err)
	}
	fmt.Println(txt[:10])
	return nil
}
