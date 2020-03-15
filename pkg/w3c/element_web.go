package w3c

import (
	"context"
	"encoding/json"
	"net/http"
)

type WebElement interface {

	// Click clicks on the element.
	Click(ctx context.Context) error

	// SendKeys types into the element.
	SendKeys(ctx context.Context, keys ...Key) error

	// Clear clears the element.
	Clear(ctx context.Context) error

	// FindOne finds a child element.
	FindOne(ctx context.Context, by FindElementStrategy, value string) (WebElement, error)

	// Find finds multiple children elements.
	Find(ctx context.Context, by FindElementStrategy, value string) ([]WebElement, error)

	// TagName returns the element's name.
	TagName(ctx context.Context) (string, error)

	// Text returns the text of the element.
	Text(ctx context.Context) (string, error)

	// IsSelected returns true if element is selected.
	IsSelected(ctx context.Context) (bool, error)

	// IsEnabled returns true if the element is enabled.
	IsEnabled(ctx context.Context) (bool, error)

	// GetAttribute returns the named attribute of the element.
	GetAttribute(ctx context.Context, name string) (string, error)

	// Rect returns the element's size.
	Rect(ctx context.Context) (Rect, error)

	// GetProperty returns the value of the specified property of the element.
	GetProperty(ctx context.Context, name string) (string, error)

	// GetCSSValue returns the value of the specified CSS property of the element.
	GetCSSValue(ctx context.Context, name string) (string, error)
}

// Point is a 2D point.
type Point struct {
	X, Y int
}

type webElement struct {
	wid     string
	sid     string
	request Doer
}

// Click clicks on the element.
func (w webElement) Click(ctx context.Context) error {
	resp, err := w.request.Do(ctx, http.MethodPost, "/session/"+w.sid+"/element/"+w.wid+"/click", nil)
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return ErrInvalidResponse
}

// SendKeys types into the element.
func (w webElement) SendKeys(ctx context.Context, keys ...Key) error {
	if len(keys) == 0 {
		return ErrInvalidArguments
	}
	p := Params{
		"value": keys,
	}
	resp, err := w.request.Do(ctx, http.MethodPost, "/session/"+w.sid+"/element/"+w.wid+"/value", p)
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return ErrInvalidResponse
}

// Clear clears the element.
func (w webElement) Clear(ctx context.Context) error {
	resp, err := w.request.Do(ctx, http.MethodPost, "/session/"+w.sid+"/element/"+w.wid+"/clear", nil)
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return ErrInvalidResponse
}

// FindOne finds a child element.
func (w webElement) FindOne(ctx context.Context, by FindElementStrategy, value string) (WebElement, error) {
	if len(value) == 0 {
		return nil, ErrInvalidArguments
	}
	p := Params{
		"using": by,
		"value": value,
	}
	resp, err := w.request.Do(ctx, http.MethodPost, "/session/"+w.sid+"/element/"+w.wid+"/element", p)
	if err != nil {
		return nil, err
	}
	var er elemResp
	if err := json.Unmarshal(resp.Value, &er); err != nil {
		return nil, err
	}
	return webElement{
		wid:     er.ID,
		sid:     w.sid,
		request: w.request,
	}, nil
}

// Find finds multiple children elements.
func (w webElement) Find(ctx context.Context, by FindElementStrategy, value string) ([]WebElement, error) {
	if len(value) == 0 {
		return nil, ErrInvalidArguments
	}
	p := Params{
		"using": by,
		"value": value,
	}
	resp, err := w.request.Do(ctx, http.MethodPost, "/session/"+w.sid+"/element/"+w.wid+"/elements", p)
	if err != nil {
		return nil, err
	}
	var elms []elemResp
	if err := json.Unmarshal(resp.Value, &elms); err != nil {
		return nil, err
	}
	webElements := make([]WebElement, len(elms))
	for i, wid := range elms {
		webElements[i] = webElement{
			wid:     wid.ID,
			sid:     w.sid,
			request: w.request,
		}
	}
	return webElements, nil
}

// TagName returns the element's name.
func (w webElement) TagName(ctx context.Context) (string, error) {
	resp, err := w.request.Do(ctx, http.MethodGet, "/session/"+w.sid+"/element/"+w.wid+"/name", nil)
	if err != nil {
		return "", err
	}
	return string(resp.Value), nil
}

// Text returns the text of the element.
func (w webElement) Text(ctx context.Context) (string, error) {
	resp, err := w.request.Do(ctx, http.MethodGet, "/session/"+w.sid+"/element/"+w.wid+"/text", nil)
	if err != nil {
		return "", err
	}
	return string(resp.Value), nil
}

// IsSelected returns true if element is selected.
func (w webElement) IsSelected(ctx context.Context) (bool, error) {
	resp, err := w.request.Do(ctx, http.MethodGet, "/session/"+w.sid+"/element/"+w.wid+"/selected", nil)
	if err != nil {
		return false, err
	}
	return string(resp.Value) == "true", nil
}

// IsEnabled returns true if the element is enabled.
func (w webElement) IsEnabled(ctx context.Context) (bool, error) {
	resp, err := w.request.Do(ctx, http.MethodGet, "/session/"+w.sid+"/element/"+w.wid+"/enabled", nil)
	if err != nil {
		return false, err
	}
	return string(resp.Value) == "true", nil
}

// GetAttribute returns the named attribute of the element.
func (w webElement) GetAttribute(ctx context.Context, name string) (string, error) {
	resp, err := w.request.Do(ctx, http.MethodGet, "/session/"+w.sid+"/element/"+w.wid+"/attribute/"+name, nil)
	if err != nil {
		return "", err
	}
	return string(resp.Value), nil
}

// GetProperty returns the value of the specified property of the element.
func (w webElement) GetProperty(ctx context.Context, name string) (string, error) {
	resp, err := w.request.Do(ctx, http.MethodGet, "/session/"+w.sid+"/element/"+w.wid+"/property/"+name, nil)
	if err != nil {
		return "", err
	}
	return string(resp.Value), nil
}

// GetCSSValue returns the value of the specified CSS property of the element.
func (w webElement) GetCSSValue(ctx context.Context, name string) (string, error) {
	resp, err := w.request.Do(ctx, http.MethodGet, "/session/"+w.sid+"/element/"+w.wid+"/css/"+name, nil)
	if err != nil {
		return "", err
	}
	return string(resp.Value), nil
}

// Rect returns the element's size.
func (w webElement) Rect(ctx context.Context) (r Rect, err error) {
	resp, err := w.request.Do(ctx, http.MethodGet, "/session/"+w.sid+"/element/"+w.wid+"/rect", nil)
	if err != nil {
		return r, err
	}
	if err := json.Unmarshal(resp.Value, &r); err != nil {
		return r, err
	}
	return r, nil
}
