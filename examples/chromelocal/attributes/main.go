package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mediabuyerbot/go-webdriver"
	"github.com/mediabuyerbot/go-webdriver/third_party/httpserver"
)

func main() {
	addr := httpserver.Addr()
	go func() {
		if err := httpserver.Run(addr); err != nil {
			exitWithError(err)
		}
	}()

	browser, err := webdriver.Chrome(nil)
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
		time.Sleep(time.Second)
	}()

	if err := browser.Url(addr); err != nil {
		exitWithError(err)
	}

	elem, err := browser.FindElementByID("promo-code-input")
	if err != nil {
		exitWithError(err)
	}

	placeholderAttr, _ := elem.Attr("placeholder")
	log.Printf("placeholderAttr=%s\n", placeholderAttr)

	elem, err = browser.FindElementByID("firstNameLabel")
	if err != nil {
		exitWithError(err)
	}

	firstNameLabel, _ := elem.Attr("for")
	log.Printf("firstNameLabel=%s\n", firstNameLabel)

	time.Sleep(5 * time.Second)
}

func exitWithError(err error) {
	fmt.Println(err)
	panic(err)
}
