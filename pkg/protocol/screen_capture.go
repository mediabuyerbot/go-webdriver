package protocol

import (
	"bytes"
	"context"
	"encoding/base64"
	"io/ioutil"
	"net/http"
)

type ScreenCapture interface {

	// Take take a screenshot of the current page.
	Take(ctx context.Context) ([]byte, error)

	// TakeElement take a screenshot of the element on the current page.
	TakeElement(ctx context.Context, elementID string) ([]byte, error)
}

type screenCapture struct {
	id      string
	request Doer
}

// NewScreenCapture creates a new instance of ScreenCapture.
func NewScreenCapture(doer Doer, sessionID string) ScreenCapture {
	return &screenCapture{
		id:      sessionID,
		request: doer,
	}
}

func (s *screenCapture) Take(ctx context.Context) ([]byte, error) {
	resp, err := s.request.Do(ctx, http.MethodGet, "/session/"+s.id+"/screenshot", nil)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(resp.Value)
	decoder := base64.NewDecoder(base64.StdEncoding, buf)
	return ioutil.ReadAll(decoder)
}

func (s *screenCapture) TakeElement(ctx context.Context, elementID string) ([]byte, error) {
	resp, err := s.request.Do(ctx, http.MethodGet, "/session/"+s.id+"/element/"+elementID+"/screenshot", nil)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(resp.Value)
	decoder := base64.NewDecoder(base64.StdEncoding, buf)
	return ioutil.ReadAll(decoder)
}
