package protocol

import (
	"context"
	"encoding/json"
	"errors"
)

const (
	// ImplicitTimeout  gives the timeout of when to abort locating an element.
	ImplicitTimeout Timeout = "implicit"

	// PageLoadTimeout provides the timeout limit used to interrupt navigation of the browsing context.
	PageLoadTimeout Timeout = "pageLoad"

	// ScriptTimeout determines when to interrupt a script that is being evaluated.
	ScriptTimeout Timeout = "script"

	DefaultTimeoutMs = Ms(10000)
)

// ErrUnknownSession returns when a new session is created with empty sessionId.
var ErrUnknownSession = errors.New("protocol: unknown session")

// Session represents the connection between a local end and a specific remote end.
// Session is equivalent to a single instantiation of a particular user agent, including all its child browsers.
type Session interface {
	// ID returns the unique session id.
	ID() SessionID

	// Capabilities returns the capabilities of the specified session.
	Capabilities() Capabilities

	// GetTimeouts returns the timeouts implicit, pageLoad, script.
	GetTimeouts(context.Context) (TimeoutInfo, error)

	// SetTimeouts configure the amount of time that a particular type of operation can execute for before
	// they are aborted and a |Timeout| error is returned to the client.  Valid values are: "script" for script timeouts,
	// "implicit" for modifying the implicit wait timeout and "pageLoad" for setting a page load timeout.
	SetTimeouts(context.Context, Timeout, Ms) error

	// Status returns information about whether a remote end is in a state
	// in which it can create new sessions, but may additionally include arbitrary
	// meta information that is specific to the implementation.
	Status(context.Context) (Status, error)

	// Delete delete the session.
	Delete(context.Context) error

	// SetImplicitTimeout set the session implicit wait timeout.
	SetImplicitTimeout(context.Context, Ms) error

	// SetPageLoadTimeout set the session page load timeout.
	SetPageLoadTimeout(context.Context, Ms) error

	// SetScriptTimeout set the session script timeout.
	SetScriptTimeout(context.Context, Ms) error
}

type (
	session struct {
		id     string
		client Client
		cap    Capabilities
	}

	Timeout   string
	SessionID string
	Ms        int

	Status map[string]interface{}

	TimeoutInfo struct {
		Implicit Ms `json:"implicit"`
		PageLoad Ms `json:"pageLoad"`
		Script   Ms `json:"script"`
	}
)

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
func (s *session) ID() SessionID {
	return SessionID(s.id)
}

// Capabilities returns the capabilities of the specified session.
func (s *session) Capabilities() Capabilities {
	return s.cap
}

// GetTimeouts returns the timeouts implicit, pageLoad, script.
func (s *session) GetTimeouts(ctx context.Context) (info TimeoutInfo, err error) {
	resp, err := s.client.Get(ctx, "/session/"+s.id+"/timeouts")
	if err != nil {
		return info, err
	}
	if err := json.Unmarshal(resp.Value, &info); err != nil {
		return info, err
	}
	return info, nil
}

// SetTimeouts configure the amount of time that a particular type of operation can execute for before
// they are aborted and a |Timeout| error is returned to the client.  Valid values are: "script" for script timeouts,
// "implicit" for modifying the implicit wait timeout and "pageLoad" for setting a page load timeout.
func (s *session) SetTimeouts(ctx context.Context, t Timeout, ms Ms) error {
	params := params{string(t): ms}
	if _, err := s.client.Post(ctx, "/session/"+s.id+"/timeouts", params); err != nil {
		return err
	}
	return nil
}

// SetImplicitTimeout set the session implicit wait timeout.
func (s *session) SetImplicitTimeout(ctx context.Context, ms Ms) error {
	return s.SetTimeouts(ctx, ImplicitTimeout, ms)
}

// SetPageLoadTimeout set the session page load timeout.
func (s *session) SetPageLoadTimeout(ctx context.Context, ms Ms) error {
	return s.SetTimeouts(ctx, PageLoadTimeout, ms)
}

// SetScriptTimeout set the session script timeout.
func (s *session) SetScriptTimeout(ctx context.Context, ms Ms) error {
	return s.SetTimeouts(ctx, ScriptTimeout, ms)
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
