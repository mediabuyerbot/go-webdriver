package w3c

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"net/http"
)

type ScreenCapture interface {

	// Take takes a screenshot of the current page.
	Take(ctx context.Context) (io.Reader, error)

	// TakeElement takes a screenshot of the element on the current page.
	TakeElement(ctx context.Context, elementID string) (io.Reader, error)
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

func (s *screenCapture) Take(ctx context.Context) (io.Reader, error) {
	resp, err := s.request.Do(ctx, http.MethodGet, "/session/"+s.id+"/screenshot", nil)
	if err != nil {
		return nil, err
	}
	if len(resp.Value) == 0 {
		return nil, ErrInvalidResponse
	}
	buf := bytes.NewBuffer(resp.Value[1 : len(resp.Value)-1])
	return base64.NewDecoder(base64.StdEncoding, buf), nil
}

func (s *screenCapture) TakeElement(ctx context.Context, elementID string) (io.Reader, error) {
	resp, err := s.request.Do(ctx, http.MethodGet, "/session/"+s.id+"/element/"+elementID+"/screenshot", nil)
	if err != nil {
		return nil, err
	}
	if len(resp.Value) == 0 {
		return nil, ErrInvalidResponse
	}
	buf := bytes.NewBuffer(resp.Value[1 : len(resp.Value)-1])
	return base64.NewDecoder(base64.StdEncoding, buf), nil
}
