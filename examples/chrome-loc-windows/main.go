package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mediabuyerbot/go-webdriver/pkg/w3c"

	"github.com/mediabuyerbot/go-webdriver"
)

func main() {
	// run browser
	browser, err := webdriver.Chrome(nil)
	if err != nil {
		exitWithError(err)
	}
	defer browser.Close()

	wait := func() {
		time.Sleep(200 * time.Millisecond)
	}

	// open tabs
	tabs := make([]w3c.WindowHandle, 10)
	for i := 0; i < 10; i++ {
		tab, err := browser.OpenTab()
		if err != nil {
			exitWithError(err)
		}
		tabs[i] = tab
		wait()
	}

	for _, tab := range tabs {
		if err := browser.SwitchTo(tab); err != nil {
			exitWithError(err)
		}

		if err := browser.MoveTo(200, 400); err != nil {
			exitWithError(err)
		}

		if err := browser.ResizeTo(600, 500); err != nil {
			exitWithError(err)
		}

		wait()
	}

	time.Sleep(5 * time.Second)
}

func exitWithError(err error) {
	fmt.Println(err)
	os.Exit(1)
}
