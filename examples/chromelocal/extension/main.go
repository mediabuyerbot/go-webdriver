package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mediabuyerbot/go-webdriver"
)

func main() {
	// load unpacked extension
	ext, err := webdriver.LoadChromeExtension("./extension.crx")
	if err != nil {
		exitWithError(err)
	}

	opts := webdriver.ChromeOptions()
	// add extension
	_ = opts.AddExtension(ext)

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
	time.Sleep(15 * time.Second)
}

func exitWithError(err error) {
	fmt.Println(err)
	panic(err)
}
