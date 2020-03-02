package integration

import (
	"context"
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

//func TestChromedriverSessionCapabilities(t *testing.T) {
//	b, c := newChromeBrowser(t)
//	defer c()
//
//	ctx := context.Background()
//
//	// timeouts
//	b.Session().SetTimeouts(ctx, protocol.ImplicitTimeout, protocol.DefaultTimeoutMs)
//	b.Session().SetTimeouts(ctx, protocol.PageLoadTimeout, protocol.DefaultTimeoutMs)
//	b.Session().SetTimeouts(ctx, protocol.ScriptTimeout, protocol.DefaultTimeoutMs)
//	timeoutInfo, err := b.Session().GetTimeouts(ctx)
//	assert.Nil(t, err)
//	assert.Equal(t, timeoutInfo.Implicit, protocol.DefaultTimeoutMs, "timeoutInfo.Implicit")
//	assert.Equal(t, timeoutInfo.PageLoad, protocol.DefaultTimeoutMs, "timeoutInfo.PageLoad")
//	assert.Equal(t, timeoutInfo.Script, protocol.DefaultTimeoutMs, "timeoutInfo.Script")
//
//	// capabilities
//	cap := b.Session().Capabilities()
//	assert.Equal(t, cap["javascriptEnabled"].(bool), true, "javascriptEnabled")
//	assert.Equal(t, cap["databaseEnabled"].(bool), false, "databaseEnabled")
//	assert.Equal(t, cap["nativeEvents"].(bool), true, "nativeEvents")
//
//	// session id
//	assert.NotEmpty(t, b.Session().ID())
//
//	// session status
//	status, err := b.Session().Status(ctx)
//	assert.Nil(t, err, "session get status")
//	assert.NotEmpty(t, status["build"].(map[string]interface{})["version"])
//
//	// session delete
//	err = b.Session().Delete(ctx)
//	assert.Nil(t, err)
//}
