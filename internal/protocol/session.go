package protocol

import (
	"context"
	"encoding/json"
	"errors"
)

var (
	ErrSessionIDEmpty = errors.New("protocol: session id is empty")
)

type Session struct {
	id     string
	client Client
	cap    Capabilities
}

// NewSession creates a new instance of Session.
// The New Session command creates a new WebDriver session with the endpoint node. If the creation fails,
// a session not created error is returned.
// https://www.w3.org/TR/webdriver1/#new-session
func NewSession(client Client, desired, required Capabilities) (*Session, error) {
	if desired == nil {
		desired = map[string]interface{}{}
	}
	params := params{
		"desiredCapabilities":  desired,
		"requiredCapabilities": required,
	}
	resp, err := client.Post(context.Background(), "/session", params)
	if err != nil {
		return nil, err
	}
	var capabilities Capabilities
	if err := json.Unmarshal(resp.Value, &capabilities); err != nil {
		return nil, err
	}
	sessID := string(resp.SessionID)
	if len(sessID) == 0 {
		return nil, ErrSessionIDEmpty
	}
	return &Session{
		id:     sessID,
		cap:    capabilities,
		client: client,
	}, nil
}

// Capabilities returns the capabilities of the specified session.
func (s *Session) Capabilities() Capabilities {
	return s.cap
}

func (s *Session) GetTimeouts() {

}

// Delete delete the session.
func (s *Session) Delete(ctx context.Context) error {
	_, err := s.client.Delete(ctx, "/session/"+s.id)
	if err != nil {
		return err
	}
	return nil
}

// Status returns information about whether a remote end is in a state
// in which it can create new sessions, but may additionally include arbitrary
// meta information that is specific to the implementation.
// https://www.w3.org/TR/webdriver1/#status
func (s *Session) Status(ctx context.Context) (st Status, err error) {
	resp, err := s.client.Get(ctx, "/status")
	if err != nil {
		return st, err
	}
	if err := json.Unmarshal(resp.Value, &st); err != nil {
		return st, err
	}
	return st, nil
}

type Status struct {
	Build struct {
		Version  string `json:"version"`
		Revision string `json:"revision"`
		Time     string `json:"time"`
	}
	OS struct {
		Arch    string `json:"arch"`
		Name    string `json:"name"`
		Version string `json:"version"`
	}
}
