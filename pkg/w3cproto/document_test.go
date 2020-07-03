package w3cproto

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

var documentErr = &Error{
	Code:    "code",
	Message: "msg",
}

func newDocument(t *testing.T, sessID string) (Document, *MockDoer, func()) {
	ctrl := gomock.NewController(t)
	cli := NewMockDoer(ctrl)
	cx := NewDocument(cli, "123")
	return cx, cli, func() {
		ctrl.Finish()
	}
}

func TestDocument_GetPageSource(t *testing.T) {
	doc, cli, done := newDocument(t, "123")
	defer done()

	ctx := context.TODO()
	wantSource := `<html><body></body></html>`

	// returns success
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/source", nil).Times(1).Return(
		&Response{
			Value: []byte(wantSource),
		}, nil)
	haveSource, err := doc.GetPageSource(ctx)
	assert.Nil(t, err)
	assert.Equal(t, wantSource, haveSource)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/source", nil).Times(1).Return(nil, documentErr)
	haveSource, err = doc.GetPageSource(ctx)
	assert.Error(t, err)
	assert.Empty(t, haveSource)
}

func TestDocument_ExecuteScript(t *testing.T) {
	doc, cli, done := newDocument(t, "123")
	defer done()

	ctx := context.TODO()
	wantScript := `return {'title': 'webdriver'}`

	// returns success
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/execute/sync",
		Params{
			"script": wantScript, "args": []interface{}{}}).Times(1).Return(
		&Response{
			Value: []byte(`{"title": "webdriver"}`),
		}, nil)
	b, err := doc.ExecuteScript(ctx, wantScript, nil)
	assert.Nil(t, err)
	var data map[string]interface{}
	err = json.Unmarshal(b, &data)
	assert.Nil(t, err)
	assert.Equal(t, data["title"], "webdriver")

	// returns error
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/execute/sync", gomock.Any()).Times(1).Return(
		nil, documentErr).Do(func(_ context.Context, method string, path string, p Params) {
		assert.Equal(t, p["script"], wantScript)
	})
	b, err = doc.ExecuteScript(ctx, wantScript, nil)
	assert.Equal(t, documentErr, err)
	assert.Empty(t, b)
}

func TestDocument_ExecuteAsyncScript(t *testing.T) {
	doc, cli, done := newDocument(t, "123")
	defer done()

	ctx := context.TODO()
	wantScript := `var done = arguments[0]; done({"title": "webdriver"});`

	// returns success
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/execute/async",
		gomock.Any()).Times(1).Return(
		&Response{
			Value: []byte(`{"title": "webdriver"}`),
		}, nil).Do(func(_ context.Context, method string, path string, p Params) {
		assert.Equal(t, wantScript, p["script"])
	})
	b, err := doc.ExecuteAsyncScript(ctx, wantScript, nil)
	assert.Nil(t, err)
	var data map[string]interface{}
	err = json.Unmarshal(b, &data)
	assert.Nil(t, err)
	assert.Equal(t, data["title"], "webdriver")

	// returns error
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/execute/async", gomock.Any()).Times(1).Return(
		nil, documentErr)
	b, err = doc.ExecuteAsyncScript(ctx, wantScript, nil)
	assert.Equal(t, documentErr, err)
	assert.Empty(t, b)
}
