package test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mediabuyerbot/go-webdriver"
)

const firefoxRemoteAddrKey = "WEBGO_FIREFOXDRIVER_ADDR"

func firefoxBrowser(t *testing.T) (webdriver.Session, func()) {
	addr := os.Getenv(firefoxRemoteAddrKey)
	if len(addr) == 0 {
		t.Skip()
	}
	browser, err := webdriver.NewSession(
		addr,
		webdriver.DesiredCapabilities{
			"platform": "linux",
			"firefoxOptions": map[string]interface{}{
				"log": "{\"level\": \"trace\"}",
				"args": []string{
					"-headless",
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

func TestFirefox_Sessions(t *testing.T) {
	browser, closeSession := firefoxBrowser(t)
	defer closeSession()

	ctx := context.Background()
	capabilities := browser.Session().Capabilities()
	assert.Equal(t, capabilities.BrowserName(), "firefox")
	assert.NotEmpty(t, browser.SessionID())

	status, err := browser.Session().Status(ctx)
	assert.Nil(t, err)
	assert.True(t, status.Ready)
}

//func TestSessionFirefoxdriverSessionCapabilities(t *testing.T) {
//	b, c := newFirefoxBrowser(t)
//	defer c()
//
//	ctx := context.Background()
//
//	// timeouts
//	b.Session().SetTimeouts(ctx, protocol.ImplicitTimeout, protocol.DefaultTimeoutMs)
//	b.Session().SetTimeouts(ctx, protocol.PageLoadTimeout, protocol.DefaultTimeoutMs)
//	b.Session().SetTimeouts(ctx, protocol.ScriptTimeout, protocol.DefaultTimeoutMs)
//
//	timeoutInfo, err := b.Session().GetTimeouts(ctx)
//	assert.Nil(t, err)
//	assert.Equal(t, timeoutInfo.Implicit, protocol.DefaultTimeoutMs, "timeoutInfo.Implicit")
//	assert.Equal(t, timeoutInfo.PageLoad, protocol.DefaultTimeoutMs, "timeoutInfo.PageLoad")
//	assert.Equal(t, timeoutInfo.Script, protocol.DefaultTimeoutMs, "timeoutInfo.Script")
//
//	// session id
//	assert.NotEmpty(t, b.SessionID())
//
//	// capabilities
//	cap := b.Session().Capabilities()
//	assert.Equal(t, cap["browserName"].(string), "firefox", "browserName")
//
//	// session status
//	status, err := b.Session().Status(ctx)
//	assert.Nil(t, err, "session get status")
//	assert.False(t, status["ready"].(bool))
//}
