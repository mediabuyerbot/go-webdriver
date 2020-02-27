package protocol

import (
	"context"
	"encoding/json"
	"errors"
)

const (
	ImplicitTimeout Timeout = "implicit"
	PageLoadTimeout Timeout = "page load"
	ScriptTimeout   Timeout = "script"

	DefaultTimeoutMs = Ms(10000)
)

var (
	ErrSessionIDEmpty = errors.New("protocol: session id is empty")
)

type (
	// Session represents the connection between a local end and a specific remote end.
	// https://www.w3.org/TR/webdriver1/#sessions
	Session struct {
		id     string
		client Client
		cap    Capabilities
	}

	Timeout string
	Ms      int
)

// NewSession creates a new instance of Session.
// The New Session command creates a new WebDriver session with the endpoint node. If the creation fails,
// a session not created error is returned.
// https://www.w3.org/TR/webdriver1/#new-session
func NewSession(client Client, desired, required interface{}) (*Session, error) {
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

// ID returns the session id.
func (s *Session) ID() string {
	return s.id
}

// Capabilities returns the capabilities of the specified session.
func (s *Session) Capabilities() Capabilities {
	return s.cap
}

// GetTimeouts returns the timeouts (implicit, pageLoad, script).
// https://www.w3.org/TR/webdriver1/#get-timeouts
func (s *Session) GetTimeouts(ctx context.Context) (info TimeoutInfo, err error) {
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
// "implicit" for modifying the implicit wait timeout and "page load" for setting a page load timeout.
// https://www.w3.org/TR/webdriver1/#set-timeouts
func (s *Session) SetTimeouts(ctx context.Context, t Timeout, ms Ms) error {
	params := params{"type": t, "ms": ms}
	if _, err := s.client.Post(ctx, "/session/"+s.id+"/timeouts", params); err != nil {
		return err
	}
	return nil
}

// SetImplicitTimeout set the session implicit wait timeout.
func (s *Session) SetImplicitTimeout(ctx context.Context, ms Ms) error {
	return s.SetTimeouts(ctx, ImplicitTimeout, ms)
}

// SetPageLoadTimeout set the session page load timeout.
func (s *Session) SetPageLoadTimeout(ctx context.Context, ms Ms) error {
	return s.SetTimeouts(ctx, PageLoadTimeout, ms)
}

// SetScriptTimeout set the session script timeout.
func (s *Session) SetScriptTimeout(ctx context.Context, ms Ms) error {
	return s.SetTimeouts(ctx, ScriptTimeout, ms)
}

// SetTimeoutsAsyncScrypt set the amount of time, in milliseconds, that asynchronous scripts
// executed by ExecuteScriptAsync() are permitted to run before they are aborted
// and a |Timeout| error is returned to the client.
func (s Session) SetTimeoutsAsyncScript(ctx context.Context, ms Ms) error {
	params := params{"ms": ms}
	if _, err := s.client.Post(ctx, "/session/"+s.id+"/timeouts/async_script", params); err != nil {
		return err
	}
	return nil
}

// SetTimeoutsImplicitWait set the amount of time the driver should wait when searching for elements.
// When searching for a single element, the driver should poll the page until an element is found or the timeout expires,
// whichever occurs first. When searching for multiple elements, the driver should poll the page until at least one element
// is found or the timeout expires, at which point it should return an empty list.
// If this command is never sent, the driver should default to an implicit wait of 0ms.
func (s Session) SetTimeoutsImplicitWait(ctx context.Context, ms Ms) error {
	params := params{"ms": ms}
	if _, err := s.client.Post(ctx, "/session/"+s.id+"/timeouts/implicit_wait", params); err != nil {
		return err
	}
	return nil
}

// Delete delete the session.
// https://www.w3.org/TR/webdriver1/#delete-session
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

type Status map[string]interface{}

type TimeoutInfo struct {
	Implicit Ms `json:"implicit"`
	PageLoad Ms `json:"pageLoad"`
	Script   Ms `json:"script"`
}
