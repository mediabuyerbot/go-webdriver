package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mediabuyerbot/go-webdriver/pkg/w3c"

	"github.com/mediabuyerbot/go-webdriver"
)

func main() {
	opts := webdriver.ChromeOptions()
	opts.AddArgument(
		// "--headless",
		"--allow-http-background-page",
		"--no-sandbox")

	opts.
		SetPref("browser.show_fullscreen_toolbar", 1).
		SetPref("hide_web_store_icon", 1)

	opts.
		SetLocalState("key", "value")

	_ = opts.SetProxy(&w3c.Proxy{
		Type:     w3c.ProxyDirectType,
		HTTPPort: 8080,
	})

	// run chrome
	browser, err := webdriver.Chrome(opts)
	if err != nil {
		exitWithError(err)
	}
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic:", err)
		}
		if err := browser.Close(); err != nil {
			log.Println(err)
		}
	}()
	time.Sleep(5 * time.Second)
}

func exitWithError(err error) {
	fmt.Println(err)
	panic(err)
}
