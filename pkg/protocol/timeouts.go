package protocol

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
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

// ErrTimeoutConfiguration returns when an invalid timeout type is specified.
var ErrTimeoutConfiguration = errors.New("protocol: timeouts configuration")

// Timeouts represents  a timeouts that control
// the behavior of script evaluation, navigation, and element retrieval.
type Timeouts interface {

	// GetTimeouts returns the timeouts implicit, pageLoad, script.
	GetTimeouts(context.Context) (TimeoutInfo, error)

	// SetTimeouts configure the amount of time that a particular type of operation can execute for before
	// they are aborted and a |Timeout| error is returned to the client.  Valid values are: "script" for script timeouts,
	// "implicit" for modifying the implicit wait timeout and "pageLoad" for setting a page load timeout.
	SetTimeouts(context.Context, Timeout, Ms) error

	// SetImplicitTimeout set the session implicit wait timeout.
	SetImplicitTimeout(context.Context, Ms) error

	// SetPageLoadTimeout set the session page load timeout.
	SetPageLoadTimeout(context.Context, Ms) error

	// SetScriptTimeout set the session script timeout.
	SetScriptTimeout(context.Context, Ms) error
}

type (
	Timeout     string
	Ms          int
	TimeoutInfo struct {
		Implicit Ms `json:"implicit"`
		PageLoad Ms `json:"pageLoad"`
		Script   Ms `json:"script"`
	}
)

type timeouts struct {
	id      string
	request Doer
}

func (t Timeout) String() string {
	return string(t)
}

// NewTimeouts creates a new instance of Timeouts.
func NewTimeouts(doer Doer, sessionID string) Timeouts {
	return &timeouts{
		id:      sessionID,
		request: doer,
	}
}

func (t *timeouts) validate(tm Timeout) bool {
	switch tm {
	case ImplicitTimeout, PageLoadTimeout, ScriptTimeout:
		return true
	default:
		return false
	}
}

// GetTimeouts returns the timeouts implicit, pageLoad, script.
func (t *timeouts) GetTimeouts(ctx context.Context) (info TimeoutInfo, err error) {
	resp, err := t.request.Do(ctx, http.MethodGet, "/session/"+t.id+"/timeouts", nil)
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
func (t *timeouts) SetTimeouts(ctx context.Context, tm Timeout, ms Ms) error {
	if ok := t.validate(tm); !ok {
		return ErrTimeoutConfiguration
	}
	p := Params{tm.String(): ms}
	resp, err := t.request.Do(ctx, http.MethodPost, "/session/"+t.id+"/timeouts", p)
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return ErrInvalidResponse
}

// SetImplicitTimeout set the session implicit wait timeout.
func (t *timeouts) SetImplicitTimeout(ctx context.Context, ms Ms) error {
	return t.SetTimeouts(ctx, ImplicitTimeout, ms)
}

// SetPageLoadTimeout set the session page load timeout.
func (t *timeouts) SetPageLoadTimeout(ctx context.Context, ms Ms) error {
	return t.SetTimeouts(ctx, PageLoadTimeout, ms)
}

// SetScriptTimeout set the session script timeout.
func (t *timeouts) SetScriptTimeout(ctx context.Context, ms Ms) error {
	return t.SetTimeouts(ctx, ScriptTimeout, ms)
}
