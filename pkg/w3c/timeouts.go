package w3c

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// Timeouts represents  a timeouts that control
// the behavior of script evaluation, navigation, and element retrieval.
type Timeouts interface {

	// Get returns the timeouts implicit, pageLoad, script.
	Get(context.Context) (Timeout, error)

	// SetImplicit sets the amount of time the driver should wait when
	// searching for elements. The timeout will be rounded to nearest millisecond.
	SetImplicit(context.Context, time.Duration) error

	// SetPageLoad sets the amount of time the driver should wait when
	// loading a page. The timeout will be rounded to nearest millisecond.
	SetPageLoad(context.Context, time.Duration) error

	// SetScript sets the amount of time that asynchronous scripts
	// are permitted to run before they are aborted. The timeout will be rounded
	// to nearest millisecond.
	SetScript(context.Context, time.Duration) error
}

const (
	implicitTimeout = "implicit"
	pageLoadTimeout = "pageLoad"
	scriptTimeout   = "script"
)

type Timeout struct {
	Implicit uint `json:"implicit"`
	PageLoad uint `json:"pageLoad"`
	Script   uint `json:"script"`
}

func (t Timeout) GetImplicit() time.Duration {
	return time.Duration(t.Implicit) * time.Millisecond
}

func (t Timeout) GetPageLoad() time.Duration {
	return time.Duration(t.PageLoad) * time.Millisecond
}

func (t Timeout) GetScript() time.Duration {
	return time.Duration(t.Script) * time.Millisecond
}

type timeouts struct {
	id      string
	request Doer
}

// NewTimeouts creates a new instance of Timeouts.
func NewTimeouts(doer Doer, sessionID string) Timeouts {
	return &timeouts{
		id:      sessionID,
		request: doer,
	}
}

func (t *timeouts) Get(ctx context.Context) (info Timeout, err error) {
	resp, err := t.request.Do(ctx, http.MethodGet, "/session/"+t.id+"/timeouts", nil)
	if err != nil {
		return info, err
	}
	if err := json.Unmarshal(resp.Value, &info); err != nil {
		return info, err
	}
	return info, nil
}

func (t *timeouts) SetImplicit(ctx context.Context, d time.Duration) error {
	return t.set(ctx, implicitTimeout, d)
}

func (t *timeouts) SetPageLoad(ctx context.Context, d time.Duration) error {
	return t.set(ctx, pageLoadTimeout, d)
}

func (t *timeouts) SetScript(ctx context.Context, d time.Duration) error {
	return t.set(ctx, scriptTimeout, d)
}

func (t *timeouts) set(ctx context.Context, tm string, d time.Duration) error {
	p := Params{tm: d.Milliseconds()}
	resp, err := t.request.Do(ctx, http.MethodPost, "/session/"+t.id+"/timeouts", p)
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return ErrInvalidResponse
}
