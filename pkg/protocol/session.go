package protocol

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

var ErrUnknownSession = errors.New("protocol: unknown session id")

// Session represents the connection between a local end and a specific remote end.
// Session is equivalent to a single instantiation of a particular user agent, including all its child browsers.
type Session interface {
	// ID returns the unique session id.
	ID() string

	// Capabilities returns the capabilities of the specified session.
	Capabilities() Capabilities

	// Status returns information about whether a remote end is in a state
	// in which it can create new sessions, but may additionally include arbitrary
	// meta information that is specific to the implementation.
	Status(context.Context) (Status, error)

	// Delete delete the session.
	Delete(context.Context) error
}

type session struct {
	id      string
	request Doer
	cap     Capabilities
}

type Status struct {
	// protocol specification
	Ready   bool   `json:"ready"`
	Message string `json:"message"`

	// extension
	Build struct {
		Version  string `json:"version"`
		Revision string `json:"revision"`
		Time     string `json:"time"`
	} `json:"build"`
	OS struct {
		Arch    string `json:"arch"`
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"os"`
}

func (s Status) HasExtensionInfo() bool {
	return len(s.Build.Version) > 0
}

// NewSession creates a new instance of Session.
// The new session command creates a new WebDriver session with the endpoint node. If the creation fails,
// a session not created error is returned.
func NewSession(request Doer, opts Options) (Session, error) {
	browserOptions := Params{}

	if opts.AlwaysMatch() != nil {
		browserOptions["alwaysMatch"] = opts.AlwaysMatch()
	}
	if len(opts.FirstMatch()) > 0 {
		browserOptions["firstMatch"] = opts.FirstMatch()
	}
	if opts.Proxy() != nil {
		browserOptions["proxy"] = opts.Proxy()
	}
	browserConfiguration := Params{
		"capabilities": browserOptions,
	}
	resp, err := request.Do(context.Background(), http.MethodPost, "/session", browserConfiguration)
	if err != nil {
		return nil, err
	}

	var sessResp sessionResponse
	if err := json.Unmarshal(resp.Value, &sessResp); err != nil {
		return nil, err
	}
	if err := sessResp.Validate(); err != nil {
		return nil, err
	}
	return &session{
		id:      sessResp.SessionID,
		cap:     sessResp.Capabilities,
		request: request,
	}, nil
}

func (s *session) ID() string {
	return s.id
}

func (s *session) Capabilities() Capabilities {
	return s.cap
}

func (s *session) Delete(ctx context.Context) error {
	resp, err := s.request.Do(ctx, http.MethodDelete, "/session/"+s.id, nil)
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return ErrInvalidResponse
}

func (s *session) Status(ctx context.Context) (st Status, err error) {
	resp, err := s.request.Do(ctx, http.MethodGet, "/status", nil)
	if err != nil {
		return st, err
	}
	if err := json.Unmarshal(resp.Value, &st); err != nil {
		return st, err
	}
	return st, nil
}

type sessionResponse struct {
	SessionID    string       `json:"sessionId"`
	Capabilities Capabilities `json:"capabilities"`
}

func (sr sessionResponse) Validate() error {
	if len(sr.SessionID) == 0 {
		return ErrUnknownSession
	}
	return nil
}
