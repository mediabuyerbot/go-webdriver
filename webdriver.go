package webdriver

const (
	Success                    = 0
	NoSuchDriver               = 6
	NoSuchElement              = 7
	NoSuchFrame                = 8
	UnknownCommand             = 9
	StaleElementReference      = 10
	ElementNotVisible          = 11
	InvalidElementState        = 12
	UnknownError               = 13
	ElementIsNotSelectable     = 15
	JavaScriptError            = 17
	XPathLookupError           = 19
	Timeout                    = 21
	NoSuchWindow               = 23
	InvalidCookieDomain        = 24
	UnableToSetCookie          = 25
	UnexpectedAlertOpen        = 26
	NoAlertOpenError           = 27
	ScriptTimeout              = 28
	InvalidElementCoordinates  = 29
	IMENotAvailable            = 30
	IMEEngineActivationFailed  = 31
	InvalidSelector            = 32
	SessionNotCreatedException = 33
	MoveTargetOutOfBounds      = 34
)

const (
	CapabilityBrowserName             = "browserName"
	CapabilityBrowserVersion          = "browserVersion"
	CapabilityPlatformName            = "platformName"
	CapabilityAcceptInsecureCerts     = "acceptInsecureCerts"
	CapabilityPageLoadStrategy        = "pageLoadStrategy"
	CapabilityProxy                   = "proxy"
	CapabilitySetWindowRect           = "setWindowRect"
	CapabilityTimeouts                = "timeouts"
	CapabilityUnhandledPromptBehavior = "unhandledPromptBehavior"
)

var statusCode = map[int]string{
	Success:                    "The command executed successfully.",
	NoSuchDriver:               "A session is either terminated or not started.",
	NoSuchElement:              "An element could not be located on the page using the given search parameters.",
	NoSuchFrame:                "A request to switch to a frame could not be satisfied because the frame could not be found.",
	UnknownCommand:             "The requested resource could not be found, or a request was received using an HTTP method that is not supported by the mapped resource.",
	StaleElementReference:      "An element command failed because the referenced element is no longer attached to the DOM.",
	ElementNotVisible:          "An element command could not be completed because the element is not visible on the page.",
	InvalidElementState:        "An element command could not be completed because the element is in an invalid state (e.g. attempting to click a disabled element).",
	UnknownError:               "An unknown server-side error occurred while processing the command.",
	ElementIsNotSelectable:     "An attempt was made to select an element that cannot be selected.",
	JavaScriptError:            "An error occurred while executing user supplied JavaScript.",
	XPathLookupError:           "An error occurred while searching for an element by XPath.",
	Timeout:                    "An operation did not complete before its timeout expired.",
	NoSuchWindow:               "A request to switch to a different window could not be satisfied because the window could not be found.",
	InvalidCookieDomain:        "An illegal attempt was made to set a cookie under a different domain than the current page.",
	UnableToSetCookie:          "A request to set a cookie's value could not be satisfied.",
	UnexpectedAlertOpen:        "A modal dialog was open, blocking this operation.",
	NoAlertOpenError:           "An attempt was made to operate on a modal dialog when one was not open.",
	ScriptTimeout:              "A script did not complete before its timeout expired.",
	InvalidElementCoordinates:  "The coordinates provided to an interactions operation are invalid.",
	IMENotAvailable:            "IME was not available.",
	IMEEngineActivationFailed:  "An IME engine could not be started.",
	InvalidSelector:            "Argument was an invalid selector (e.g. XPath/CSS).",
	SessionNotCreatedException: "A new session could not be created.",
	MoveTargetOutOfBounds:      "Target provided for a move action is out of bounds.",
}

type Capabilities map[string]interface{}

type Webdriver struct {
	client Client
}

func (w *Webdriver) NewBrowser() (*Browser, error) {
	return &Browser{
		cookies: &Cookies{sessionID: ""},
		window:  &Window{sessionID: ""},
		element: &Element{sessionID: ""},
	}, nil
}

func (w *Webdriver) Close() error {
	return nil
}
