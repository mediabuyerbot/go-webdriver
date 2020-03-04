package integration

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mediabuyerbot/go-webdriver"
)

const chromeRemoteAddrKey = "WEBGO_CHROMEDRIVER_ADDR"

func chromeBrowser(t *testing.T) (webdriver.Session, func()) {
	addr := os.Getenv(chromeRemoteAddrKey)
	if len(addr) == 0 {
		t.Skip()
	}
	browser, err := webdriver.NewSession(
		addr,
		webdriver.DesiredCapabilities{
			"platform": "linux",
			"chromeOptions": map[string][]string{
				"args": []string{
					// "--headless",
					// "--disable-gpu",
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

func TestChrome_Sessions(t *testing.T) {
	browser, closeSession := chromeBrowser(t)
	defer closeSession()

	ctx := context.Background()
	capabilities := browser.Session().Capabilities()
	assert.Equal(t, capabilities.BrowserName(), "chrome")
	assert.NotEmpty(t, browser.SessionID())

	status, err := browser.Session().Status(ctx)
	assert.Nil(t, err)
	assert.True(t, status.Ready)
}

func TestChrome_ScreenCapture(t *testing.T) {
	browser, closeSession := chromeBrowser(t)
	defer closeSession()

	ctx := context.Background()
	err := browser.Navigation().NavigateTo(ctx, "https://google.com")
	assert.Nil(t, err)

	reader, err := browser.ScreenCapture().Take(ctx)
	assert.Nil(t, err)
	b, err := ioutil.ReadAll(reader)
	assert.Nil(t, err)
	assert.NotEmpty(t, b)
}
