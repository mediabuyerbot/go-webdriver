package w3c

import (
	"context"
	"net/http"
)

type Document interface {

	// GetPageSource returns a string serialization of the DOM of the current browsing context active document.
	GetPageSource(ctx context.Context) (source string, err error)

	// ExecuteScript inject a snippet of JavaScript into the page for execution in the context of the currently
	// selected frame. The executed script is assumed to be synchronous and the result of evaluating the script
	// is returned to the client. The script argument defines the script to execute in the form of a function body.
	// The value returned by that function will be returned to the client. The function will be invoked with
	// the provided args array and the values may be accessed via the arguments object in the order specified.
	// Arguments may be any JSON-primitive, array, or JSON object. JSON objects that define a WebElement
	// reference will be converted to the corresponding DOM element. Likewise, any WebElements in the script result
	// will be returned to the client as WebElement JSON objects.
	ExecuteScript(ctx context.Context, script string, args []interface{}) ([]byte, error)

	// ExecuteAsyncScript  inject a snippet of JavaScript into the page for execution in the context of the
	// currently selected frame. The executed script is assumed to be asynchronous and must signal that
	// is done by invoking the provided callback, which is always provided as the final argument to the function.
	// The value to this callback will be returned to the client. Asynchronous script commands may not span page loads.
	// If an unload event is fired while waiting for a script result, an error should be returned to the client.
	// The script argument defines the script to execute in teh form of a function body. The function will be invoked
	// with the provided args array and the values may be accessed via the arguments object in the order specified.
	// The final argument will always be a callback function that must be invoked to signal that the script has finished.
	// Arguments may be any JSON-primitive, array, or JSON object. JSON objects that define a WebElement reference
	// will be converted to the corresponding DOM element. Likewise, any WebElements in the script result
	// will be returned to the client as WebElement JSON objects.
	ExecuteAsyncScript(ctx context.Context, script string, args []interface{}) ([]byte, error)
}

type document struct {
	request Doer
	id      string
}

// NewDocument creates a new instance of Document.
func NewDocument(doer Doer, sessID string) Document {
	return &document{
		request: doer,
		id:      sessID,
	}
}

func (d *document) GetPageSource(ctx context.Context) (source string, err error) {
	resp, err := d.request.Do(ctx, http.MethodGet, "/session/"+d.id+"/source", nil)
	if err != nil {
		return source, err
	}
	return string(resp.Value), nil
}

func (d *document) ExecuteScript(ctx context.Context, script string, args []interface{}) ([]byte, error) {
	if args == nil {
		args = []interface{}{}
	}
	p := Params{"script": script, "args": args}
	resp, err := d.request.Do(ctx, http.MethodPost, "/session/"+d.id+"/execute/sync", p)
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

func (d *document) ExecuteAsyncScript(ctx context.Context, script string, args []interface{}) ([]byte, error) {
	p := Params{"script": script, "args": args}
	resp, err := d.request.Do(ctx, http.MethodPost, "/session/"+d.id+"/execute/async", p)
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}
