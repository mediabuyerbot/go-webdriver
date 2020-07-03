package w3cproto

import (
	"context"
	"encoding/json"
	"net/http"
)

type FindElementStrategy string

// Methods by which to find elements.
const (
	// Returns an element whose ID attribute matches the search value.
	ByID = FindElementStrategy("id")

	// Returns an element matching an XPath expression.
	ByXPATH = FindElementStrategy("xpath")

	// Returns an anchor element whose visible text matches the search value.
	ByLinkText = FindElementStrategy("link text")

	// Returns an anchor element whose visible text partially matches the search value.
	ByPartialLinkText = FindElementStrategy("partial link text")

	// Returns an element whose NAME attribute matches the search value.
	ByName = FindElementStrategy("name")

	// Returns an element whose tag name matches the search value.
	ByTagName = FindElementStrategy("tag name")

	// Returns an element whose class name contains the search value.
	// Compound class names are not permitted.
	ByClassName = FindElementStrategy("class name")

	// Returns an element matching a CSS selector.
	ByCSSSelector = FindElementStrategy("css selector")
)

type Key string

// Special keyboard keys, for SendKeys.
const (
	NullKey       = Key('\ue000')
	CancelKey     = Key('\ue001')
	HelpKey       = Key('\ue002')
	BackspaceKey  = Key('\ue003')
	TabKey        = Key('\ue004')
	ClearKey      = Key('\ue005')
	ReturnKey     = Key('\ue006')
	EnterKey      = Key('\ue007')
	ShiftKey      = Key('\ue008')
	ControlKey    = Key('\ue009')
	AltKey        = Key('\ue00a')
	PauseKey      = Key('\ue00b')
	EscapeKey     = Key('\ue00c')
	SpaceKey      = Key('\ue00d')
	PageUpKey     = Key('\ue00e')
	PageDownKey   = Key('\ue00f')
	EndKey        = Key('\ue010')
	HomeKey       = Key('\ue011')
	LeftArrowKey  = Key('\ue012')
	UpArrowKey    = Key('\ue013')
	RightArrowKey = Key('\ue014')
	DownArrowKey  = Key('\ue015')
	InsertKey     = Key('\ue016')
	DeleteKey     = Key('\ue017')
	SemicolonKey  = Key('\ue018')
	EqualsKey     = Key('\ue019')
	Numpad0Key    = Key('\ue01a')
	Numpad1Key    = Key('\ue01b')
	Numpad2Key    = Key('\ue01c')
	Numpad3Key    = Key('\ue01d')
	Numpad4Key    = Key('\ue01e')
	Numpad5Key    = Key('\ue01f')
	Numpad6Key    = Key('\ue020')
	Numpad7Key    = Key('\ue021')
	Numpad8Key    = Key('\ue022')
	Numpad9Key    = Key('\ue023')
	MultiplyKey   = Key('\ue024')
	AddKey        = Key('\ue025')
	SeparatorKey  = Key('\ue026')
	SubstractKey  = Key('\ue027')
	DecimalKey    = Key('\ue028')
	DivideKey     = Key('\ue029')
	F1Key         = Key('\ue031')
	F2Key         = Key('\ue032')
	F3Key         = Key('\ue033')
	F4Key         = Key('\ue034')
	F5Key         = Key('\ue035')
	F6Key         = Key('\ue036')
	F7Key         = Key('\ue037')
	F8Key         = Key('\ue038')
	F9Key         = Key('\ue039')
	F10Key        = Key('\ue03a')
	F11Key        = Key('\ue03b')
	F12Key        = Key('\ue03c')
	MetaKey       = Key('\ue03d')
)

type Elements interface {

	// FindOne finds an element on the page, starting from the document root.
	FindOne(ctx context.Context, by FindElementStrategy, value string) (WebElement, error)

	// Find finds multiple elements on the page, starting from the document root.
	Find(ctx context.Context, by FindElementStrategy, value string) ([]WebElement, error)

	// Active returns the currently active element on the page.
	Active(ctx context.Context) (WebElement, error)
}

type elemResp map[string]string

func (er elemResp) ID() (s string, ok bool) {
	s, ok = er[WebElementIdentifier]
	return
}

type elements struct {
	id      string
	request Doer
}

func NewElements(doer Doer, sessID string) Elements {
	return &elements{
		id:      sessID,
		request: doer,
	}
}

func (e *elements) FindOne(ctx context.Context, by FindElementStrategy, value string) (WebElement, error) {
	if len(value) == 0 {
		return nil, ErrInvalidArguments
	}
	p := Params{
		"value": value,
		"using": by,
	}
	resp, err := e.request.Do(ctx, http.MethodPost, "/session/"+e.id+"/element", p)
	if err != nil {
		return nil, err
	}
	var er elemResp
	if err := json.Unmarshal(resp.Value, &er); err != nil {
		return nil, err
	}
	id, ok := er.ID()
	if !ok {
		return nil, ErrNoSuchElement
	}
	return webElement{
		wid:     id,
		sid:     e.id,
		request: e.request,
	}, nil
}

func (e *elements) Find(ctx context.Context, by FindElementStrategy, value string) ([]WebElement, error) {
	if len(value) == 0 {
		return nil, ErrInvalidArguments
	}
	p := Params{
		"value": value,
		"using": by,
	}
	resp, err := e.request.Do(ctx, http.MethodPost, "/session/"+e.id+"/elements", p)
	if err != nil {
		return nil, err
	}
	var elms []elemResp
	if err := json.Unmarshal(resp.Value, &elms); err != nil {
		return nil, err
	}
	webElements := make([]WebElement, len(elms))
	for i, wid := range elms {
		id, ok := wid.ID()
		if !ok {
			return nil, ErrNoSuchElement
		}
		webElements[i] = webElement{
			wid:     id,
			sid:     e.id,
			request: e.request,
		}
	}
	return webElements, nil
}

func (e *elements) Active(ctx context.Context) (WebElement, error) {
	resp, err := e.request.Do(ctx, http.MethodGet, "/session/"+e.id+"/element/active", nil)
	if err != nil {
		return nil, err
	}
	var er elemResp
	if err := json.Unmarshal(resp.Value, &er); err != nil {
		return nil, err
	}
	id, ok := er.ID()
	if !ok {
		return nil, ErrNoSuchElement
	}
	return webElement{
		wid:     id,
		sid:     e.id,
		request: e.request,
	}, nil
}
