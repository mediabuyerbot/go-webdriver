package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mediabuyerbot/go-webdriver"
)

func main() {
	// load extension from local path
	ext, err := webdriver.LoadChromeExtension("./extension.crx")
	if err != nil {
		exitWithError(err)
	}

	opts := webdriver.ChromeOptions()
	// add extension
	_ = opts.AddExtension(ext)

	// creates a new instance of the remote browser
	browser, err := webdriver.Remote("http://localhost:9515", opts.Build())
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
	time.Sleep(15 * time.Second)
}

func exitWithError(err error) {
	fmt.Println(err)
	panic(err)
}
