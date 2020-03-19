package webdriver

import (
	"context"

	"github.com/mediabuyerbot/go-webdriver/pkg/w3c"
)

type WebElement struct {
	elem w3c.WebElement
	ctx  context.Context
}
