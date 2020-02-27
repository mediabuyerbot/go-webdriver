package test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mediabuyerbot/go-webdriver/pkg/protocol"

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
			"firefoxOptions": map[string]interface{}{
				"log": "{\"level\": \"trace\"}",
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

	ctx := context.Background()

	// timeouts
	b.Session().SetTimeouts(ctx, protocol.ImplicitTimeout, protocol.DefaultTimeoutMs)
	b.Session().SetTimeouts(ctx, protocol.PageLoadTimeout, protocol.DefaultTimeoutMs)
	b.Session().SetTimeouts(ctx, protocol.ScriptTimeout, protocol.DefaultTimeoutMs)

	timeoutInfo, err := b.Session().GetTimeouts(ctx)
	assert.Nil(t, err)
	assert.Equal(t, timeoutInfo.Implicit, protocol.DefaultTimeoutMs, "timeoutInfo.Implicit")
	assert.Equal(t, timeoutInfo.PageLoad, protocol.DefaultTimeoutMs, "timeoutInfo.PageLoad")
	assert.Equal(t, timeoutInfo.Script, protocol.DefaultTimeoutMs, "timeoutInfo.Script")

	// session id
	assert.NotEmpty(t, b.SessionID())

	// capabilities
	cap := b.Session().Capabilities()
	assert.Equal(t, cap["browserName"].(string), "firefox", "browserName")

	// session status
	status, err := b.Session().Status(ctx)
	assert.Nil(t, err, "session get status")
	assert.False(t, status["ready"].(bool))
}
