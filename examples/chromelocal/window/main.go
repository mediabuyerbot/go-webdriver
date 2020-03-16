package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mediabuyerbot/go-webdriver"
	"github.com/mediabuyerbot/go-webdriver/pkg/w3c"
)

func main() {
	// run browser
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

	tabs := make([]w3c.WindowHandle, 10)

	// open windows
	if err := openTabs(browser, tabs); err != nil {
		exitWithError(err)
	}

	// switch tabs
	if err := switchTabs(browser, tabs); err != nil {
		exitWithError(err)
	}

	// switch sizes
	if err := switchSize(browser); err != nil {
		exitWithError(err)
	}

	// stats
	if err := stats(browser); err != nil {
		exitWithError(err)
	}

	// closeTabs
	if err := closeTabs(browser, tabs); err != nil {
		exitWithError(err)
	}
}

func wait() {
	time.Sleep(250 * time.Millisecond)
}

func sleep() {
	time.Sleep(time.Second)
}

func openTabs(browser *webdriver.Browser, tabs []w3c.WindowHandle) error {
	log.Println("1. openTabs")

	tab, err := browser.ActiveWindow()
	if err != nil {
		return err
	}
	tabs[0] = tab
	for i := 0; i < len(tabs)-1; i++ {
		win, err := browser.OpenTab()
		if err != nil {
			return err
		}
		tabs[i] = win
		wait()
	}
	return nil
}

func switchTabs(browser *webdriver.Browser, tabs []w3c.WindowHandle) error {
	log.Println("2. switchTabs")

	for i := 0; i < len(tabs)*3; i++ {
		tab := tabs[i%len(tabs)]
		if err := browser.SwitchTo(tab); !w3c.IsUnknownWindowHandler(err) {
			return err
		}
		wait()
	}
	return nil
}

func switchSize(browser *webdriver.Browser) error {
	log.Println("3. switchSize")

	var (
		width, height int
		x, y          int
	)
	width, height, err := browser.ScreenSize()
	if err != nil {
		return err
	}
	x, y, err = browser.ScreenPosition()
	if err != nil {
		return err
	}

	log.Println("3.1 -fullscreen")
	if err := browser.Fullscreen(); err != nil {
		return err
	}

	sleep()
	log.Println("3.2 -restore")
	if _, err := browser.ResizeWindow(w3c.Rect{
		Width:  width,
		Height: height,
		X:      x,
		Y:      y,
	}); err != nil {
		return err
	}

	sleep()
	log.Println("3.3 -minimize")
	if err := browser.Minimize(); err != nil {
		return err
	}

	sleep()
	log.Println("3.4 -maximize")
	if err := browser.Maximize(); err != nil {
		return err
	}
	return nil
}

func stats(browser *webdriver.Browser) error {
	log.Println("4. stats")

	tab, err := browser.ActiveWindow()
	if err != nil {
		return err
	}

	log.Println("4.1 Current window", tab)
	tabs, err := browser.Windows()
	if err != nil {
		return err
	}
	for _, tid := range tabs {
		if tab == tid {
			log.Println("- IsActive", tid)
			continue
		}
		log.Println("- ", tid)
	}
	return nil
}

func closeTabs(browser *webdriver.Browser, tabs []w3c.WindowHandle) error {
	log.Println("5. closeTabs")
	tab, err := browser.ActiveWindow()
	if err != nil {
		return err
	}
	for _, tid := range tabs {
		if len(tid) == 0 {
			continue
		}
		if tid == tab {
			log.Println(" - skip active tab", tid)
			continue
		}
		log.Println("- remove tab", tid)

		if err := browser.CloseWindow(tid); err != nil {
			return err
		}
	}

	return nil
}

func exitWithError(err error) {
	fmt.Println(err)
	panic(err)
}
