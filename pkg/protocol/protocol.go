package protocol

const (
	SuccessStatusCode                    = 0
	NoSuchDriverStatusCode               = 6
	NoSuchElementStatusCode              = 7
	NoSuchFrameStatusCode                = 8
	UnknownCommandStatusCode             = 9
	StaleElementReferenceStatusCode      = 10
	ElementNotVisibleStatusCode          = 11
	InvalidElementStateStatusCode        = 12
	UnknownErrorStatusCode               = 13
	ElementIsNotSelectableStatusCode     = 15
	JavaScriptErrorStatusCode            = 17
	XPathLookupErrorStatusCode           = 19
	TimeoutStatusCode                    = 21
	NoSuchWindowStatusCode               = 23
	InvalidCookieDomainStatusCode        = 24
	UnableToSetCookieStatusCode          = 25
	UnexpectedAlertOpenStatusCode        = 26
	NoAlertOpenErrorStatusCode           = 27
	ScriptTimeoutStatusCode              = 28
	InvalidElementCoordinatesStatusCode  = 29
	IMENotAvailableStatusCode            = 30
	IMEEngineActivationFailedStatusCode  = 31
	InvalidSelectorStatusCode            = 32
	SessionNotCreatedExceptionStatusCode = 33
	MoveTargetOutOfBoundsStatusCode      = 34
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

const (
	StatusMissingCommandParameters = "400: Missing Command Parameters"
	StatusUnknownCommand           = "404: Unknown command/Resource Not Found"
	StatusInvalidCommandMethod     = "405: Invalid Command Method"
	StatusFailedCommand            = "500: Failed Command"
	StatusUnimplementedCommand     = "501: Unimplemented Command"
)

var statusCode = map[int]string{
	SuccessStatusCode:                    "The command executed successfully.",
	NoSuchDriverStatusCode:               "A session is either terminated or not started.",
	NoSuchElementStatusCode:              "An element could not be located on the page using the given search parameters.",
	NoSuchFrameStatusCode:                "A request to switch to a frame could not be satisfied because the frame could not be found.",
	UnknownCommandStatusCode:             "The requested resource could not be found, or a request was received using an HTTP method that is not supported by the mapped resource.",
	StaleElementReferenceStatusCode:      "An element command failed because the referenced element is no longer attached to the DOM.",
	ElementNotVisibleStatusCode:          "An element command could not be completed because the element is not visible on the page.",
	InvalidElementStateStatusCode:        "An element command could not be completed because the element is in an invalid state (e.g. attempting to click a disabled element).",
	UnknownErrorStatusCode:               "An unknown server-side error occurred while processing the command.",
	ElementIsNotSelectableStatusCode:     "An attempt was made to select an element that cannot be selected.",
	JavaScriptErrorStatusCode:            "An error occurred while executing user supplied JavaScript.",
	XPathLookupErrorStatusCode:           "An error occurred while searching for an element by XPath.",
	TimeoutStatusCode:                    "An operation did not complete before its timeout expired.",
	NoSuchWindowStatusCode:               "A request to switch to a different window could not be satisfied because the window could not be found.",
	InvalidCookieDomainStatusCode:        "An illegal attempt was made to set a cookie under a different domain than the current page.",
	UnableToSetCookieStatusCode:          "A request to set a cookie's value could not be satisfied.",
	UnexpectedAlertOpenStatusCode:        "A modal dialog was open, blocking this operation.",
	NoAlertOpenErrorStatusCode:           "An attempt was made to operate on a modal dialog when one was not open.",
	ScriptTimeoutStatusCode:              "A script did not complete before its timeout expired.",
	InvalidElementCoordinatesStatusCode:  "The coordinates provided to an interactions operation are invalid.",
	IMENotAvailableStatusCode:            "IME was not available.",
	IMEEngineActivationFailedStatusCode:  "An IME engine could not be started.",
	InvalidSelectorStatusCode:            "Argument was an invalid selector (e.g. XPath/CSS).",
	SessionNotCreatedExceptionStatusCode: "A new session could not be created.",
	MoveTargetOutOfBoundsStatusCode:      "Target provided for a move action is out of bounds.",
}

type Capabilities map[string]interface{}

func copyCap(m map[string]interface{}) Capabilities {
	cp := make(Capabilities)
	for k, v := range m {
		vm, ok := v.(map[string]interface{})
		if ok {
			cp[k] = copyCap(vm)
		} else {
			cp[k] = v
		}
	}
	return cp
}
