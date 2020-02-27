package test

import (
	"context"
	"os"
	"testing"

	"github.com/mediabuyerbot/go-webdriver"
)

func newFirefoxBrowser(t *testing.T) (*webdriver.Session, func()) {
	addr := os.Getenv("WEBGO_FIREFOXDRIVER_ADDR")
	if len(addr) == 0 {
		t.Fatal("WEBGO_FIREFOXDRIVER_ADDR env is not assigned")
	}
	browser, err := webdriver.NewSession(
		addr,
		webdriver.DesiredCapabilities{
			"platform": "linux",
			"firefoxOptions": map[string][]string{
				"args": []string{
					"--headless",
					"--disable-gpu",
					"--no-sandbox",
				},
			},
		},
		webdriver.RequiredCapabilities{},
	)
	if err != nil {
		t.Fatal(err)
	}
	return browser, func() {
		if err := browser.Close(context.TODO()); err != nil {
			t.Fatal(err)
		}
	}
}

func TestSessionFirefoxdriverSessionCapabilities(t *testing.T) {
	b, c := newFirefoxBrowser(t)
	defer c()

	b = b

}
