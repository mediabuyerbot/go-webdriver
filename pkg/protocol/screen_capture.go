package protocol

import "context"

type ScreenCapture interface {

	// TakeScreenshot take a screenshot of the current page.
	TakeScreenshot(ctx context.Context) ([]byte, error)

	// TakeElementScreenshot take a screenshot of the element on the current page.
	TakeElementScreenshot(ctx context.Context, elementID string) ([]byte, error)
}
