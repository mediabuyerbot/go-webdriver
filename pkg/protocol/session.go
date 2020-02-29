package protocol

import (
	"context"
	"encoding/json"
	"errors"
)

// ErrUnknownSession returns when a new session is created with empty sessionId.
var ErrUnknownSession = errors.New("protocol: unknown session")

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

type (
	session struct {
		id     string
		client Client
		cap    Capabilities
	}
)

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
func NewSession(client Client, desired, required interface{}) (Session, error) {
	if desired == nil {
		desired = map[string]interface{}{}
	}
	if required == nil {
		required = map[string]interface{}{}
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
	sessID := resp.SessionID
	// firefox
	if len(sessID) == 0 {
		sess, ok := capabilities["sessionId"]
		if !ok {
			return nil, ErrUnknownSession
		}
		sessID = sess.(string)
		if len(sessID) == 0 {
			return nil, ErrUnknownSession
		}
		caps, ok := capabilities["capabilities"]
		if ok {
			capabilities = copyCap(caps.(map[string]interface{}))
		}
	}
	return &session{
		id:     sessID,
		cap:    capabilities,
		client: client,
	}, nil
}

// ID returns the unique session id.
func (s *session) ID() string {
	return s.id
}

// Capabilities returns the capabilities of the specified session.
func (s *session) Capabilities() Capabilities {
	return s.cap
}

// Delete delete the session.
func (s *session) Delete(ctx context.Context) error {
	_, err := s.client.Delete(ctx, "/session/"+s.id)
	if err != nil {
		return err
	}
	return nil
}

// Status returns information about whether a remote end is in a state
// in which it can create new sessions, but may additionally include arbitrary
// meta information that is specific to the implementation.
func (s *session) Status(ctx context.Context) (st Status, err error) {
	resp, err := s.client.Get(ctx, "/status")
	if err != nil {
		return st, err
	}
	if err := json.Unmarshal(resp.Value, &st); err != nil {
		return st, err
	}
	return st, nil
}
