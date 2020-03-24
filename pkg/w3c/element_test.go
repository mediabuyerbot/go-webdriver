package w3c

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

const testWebElementID = "73101597-492f-4ffe-8f75-bd7bd0acb691"

var elementWebErr = &Error{
	Code:    "code",
	Message: "msg",
}

func newWebElement(t *testing.T, sessID string) (WebElement, *MockDoer, func()) {
	ctrl := gomock.NewController(t)
	cli := NewMockDoer(ctrl)
	cx := webElement{
		sid:     sessID,
		wid:     testWebElementID,
		request: cli,
	}
	return cx, cli, func() {
		ctrl.Finish()
	}
}

func TestWebElement_ID(t *testing.T) {
	webElem, _, done := newWebElement(t, "123")
	defer done()
	assert.Equal(t, testWebElementID, webElem.ID())
}

func TestWebElement_Click(t *testing.T) {
	webElem, cli, done := newWebElement(t, "123")
	defer done()

	ctx := context.Background()

	// returns success
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/element/"+testWebElementID+"/click", nil).Times(1).Return(
		&Response{
			Value: []byte(`null`),
		}, nil)
	err := webElem.Click(ctx)
	assert.Nil(t, err)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/element/"+testWebElementID+"/click", nil).Times(1).Return(nil, elementWebErr)
	err = webElem.Click(ctx)
	assert.Error(t, err)
	assert.Equal(t, elementWebErr, err)

	// returns error (custom JSON format)
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/element/"+testWebElementID+"/click", nil).Times(1).Return(
		&Response{
			Value: []byte(`{}`),
		}, nil)
	err = webElem.Click(ctx)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidResponse, err)
}

func TestWebElement_Find(t *testing.T) {
	webElem, cli, done := newWebElement(t, "123")
	defer done()

	ctx := context.Background()
	want := "73101597-492f-4ffe-8f75-bd7bd0acb691"

	// returns success
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/element/"+testWebElementID+"/elements", gomock.Any()).Times(1).Return(&Response{
		Value: []byte(`[{"element-6066-11e4-a52e-4f735466cecf":"73101597-492f-4ffe-8f75-bd7bd0acb691"}]`),
	}, nil).Do(func(_ context.Context, path string, method string, p Params) {
		assert.Equal(t, ByCSSSelector, p["using"])
		assert.Equal(t, ".className", p["value"])
	})
	elems, err := webElem.Find(ctx, ByCSSSelector, ".className")
	assert.Nil(t, err)
	assert.Len(t, elems, 1)
	assert.Equal(t, want, elems[0].ID())

	// returns error
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/element/"+testWebElementID+"/elements", gomock.Any()).Times(1).Return(&Response{}, elementsErr)
	elems, err = webElem.Find(ctx, ByCSSSelector, ".className")
	assert.Error(t, err)
	assert.Equal(t, elementsErr, err)
	assert.Nil(t, elems)

	// returns error no such element
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/element/"+testWebElementID+"/elements", gomock.Any()).Times(1).Return(&Response{
		Value: []byte(`[{"element":"73101597-492f-4ffe-8f75-bd7bd0acb691"}]`),
	}, nil)
	elems, err = webElem.Find(ctx, ByCSSSelector, ".className")
	assert.Error(t, err)
	assert.Equal(t, ErrNoSuchElement, err)
	assert.Nil(t, elems)

	// returns error invalid arguments
	elems, err = webElem.Find(ctx, ByCSSSelector, "")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidArguments, err)
	assert.Nil(t, elems)

	// returns error unmarshal JSON
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/element/"+testWebElementID+"/elements", gomock.Any()).Times(1).Return(&Response{
		Value: []byte(`[{element":"73101597-492f-4ffe-8f75-bd7bd0acb691"}]`),
	}, nil)
	elems, err = webElem.Find(ctx, ByCSSSelector, "#id")
	assert.Error(t, err)
	assert.Nil(t, elems)
}

func TestWebElement_FindOne(t *testing.T) {
	webElem, cli, done := newWebElement(t, "123")
	defer done()

	ctx := context.Background()
	want := "73101597-492f-4ffe-8f75-bd7bd0acb691"

	// returns success
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/element/"+testWebElementID+"/element", gomock.Any()).Times(1).Return(&Response{
		Value: []byte(`{"element-6066-11e4-a52e-4f735466cecf":"73101597-492f-4ffe-8f75-bd7bd0acb691"}`),
	}, nil).Do(func(_ context.Context, path string, method string, p Params) {
		assert.Equal(t, ByCSSSelector, p["using"])
		assert.Equal(t, "#id", p["value"])
	})
	elem, err := webElem.FindOne(ctx, ByCSSSelector, "#id")
	assert.Nil(t, err)
	assert.Equal(t, want, webElem.ID())

	// returns error
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/element/"+testWebElementID+"/element", gomock.Any()).Times(1).Return(&Response{}, elementsErr)
	elem, err = webElem.FindOne(ctx, ByCSSSelector, "#id")
	assert.Error(t, err)
	assert.Equal(t, elementsErr, err)
	assert.Nil(t, elem)

	// returns error no such element
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/element/"+testWebElementID+"/element", gomock.Any()).Times(1).Return(&Response{
		Value: []byte(`{"element":"73101597-492f-4ffe-8f75-bd7bd0acb691"}`),
	}, nil)
	elem, err = webElem.FindOne(ctx, ByCSSSelector, "#id")
	assert.Error(t, err)
	assert.Equal(t, ErrNoSuchElement, err)
	assert.Nil(t, elem)

	// returns error invalid arguments
	elem, err = webElem.FindOne(ctx, ByCSSSelector, "")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidArguments, err)
	assert.Nil(t, elem)

	// returns error unmarshal JSON
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/element/"+testWebElementID+"/element", gomock.Any()).Times(1).Return(&Response{
		Value: []byte(`{element":"73101597-492f-4ffe-8f75-bd7bd0acb691"}`),
	}, nil)
	elem, err = webElem.FindOne(ctx, ByCSSSelector, "#id")
	assert.Error(t, err)
	assert.Nil(t, elem)
}

func TestWebElement_GetAttribute(t *testing.T) {
	webElem, cli, done := newWebElement(t, "123")
	defer done()

	ctx := context.Background()

	// returns success
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/attribute/class", nil).Times(1).Return(
		&Response{Value: []byte(`"className"`)}, nil)
	val, err := webElem.GetAttribute(ctx, "class")
	assert.Nil(t, err)
	assert.Equal(t, "className", val)

	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/attribute/class", nil).Times(1).Return(
		&Response{Value: []byte(`null`)}, nil)
	val, err = webElem.GetAttribute(ctx, "class")
	assert.Nil(t, err)
	assert.Equal(t, "", val)

	// returns error invalid args
	val, err = webElem.GetAttribute(ctx, "")
	assert.Error(t, err)
	assert.Empty(t, val)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/attribute/class", nil).Times(1).Return(
		nil, elementWebErr)
	val, err = webElem.GetAttribute(ctx, "class")
	assert.Error(t, err)
	assert.Equal(t, elementWebErr, err)
	assert.Empty(t, val)

	// returns error JSON unmarshal
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/attribute/class", nil).Times(1).Return(
		&Response{Value: []byte(`{`)}, nil)
	val, err = webElem.GetAttribute(ctx, "class")
	assert.Error(t, err)
	assert.Empty(t, val)
}

func TestWebElement_GetCSSValue(t *testing.T) {
	webElem, cli, done := newWebElement(t, "123")
	defer done()

	ctx := context.Background()

	// return success
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/css/border", nil).Times(1).Return(
		&Response{
			Value: []byte(`"rgba(0, 0, 0, 0)"`),
		}, nil)
	val, err := webElem.GetCSSValue(ctx, "border")
	assert.Nil(t, err)
	assert.Equal(t, "rgba(0, 0, 0, 0)", val)

	// returns error invalid args
	val, err = webElem.GetCSSValue(ctx, "")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidArguments, err)
	assert.Empty(t, val)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/css/border", nil).Times(1).Return(
		nil, elementWebErr)
	val, err = webElem.GetCSSValue(ctx, "border")
	assert.Error(t, err)
	assert.Equal(t, elementWebErr, err)
	assert.Empty(t, val)

	// returns error JSON unmarshal
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/css/border", nil).Times(1).Return(
		&Response{
			Value: []byte(`{`),
		}, nil)
	val, err = webElem.GetCSSValue(ctx, "border")
	assert.Error(t, err)
	assert.Empty(t, val)
}

func TestWebElement_GetProperty(t *testing.T) {
	webElem, cli, done := newWebElement(t, "123")
	defer done()

	ctx := context.Background()

	// returns success
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/property/id", nil).Times(1).Return(
		&Response{
			Value: []byte(`"id"`),
		}, nil)
	val, err := webElem.GetProperty(ctx, "id")
	assert.Nil(t, err)
	assert.Equal(t, "id", val)

	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/property/id", nil).Times(1).Return(
		&Response{
			Value: []byte(`null`),
		}, nil)
	val, err = webElem.GetProperty(ctx, "id")
	assert.Nil(t, err)
	assert.Equal(t, "", val)

	// returns error invalid args
	val, err = webElem.GetProperty(ctx, "")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidArguments, err)
	assert.Empty(t, val)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/property/id", nil).Times(1).Return(
		nil, elementWebErr)
	val, err = webElem.GetProperty(ctx, "id")
	assert.Error(t, err)
	assert.Equal(t, elementWebErr, err)
	assert.Empty(t, val)

	// returns error JSON unmarshal
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/property/id", nil).Times(1).Return(
		&Response{
			Value: []byte(`{`),
		}, nil)
	val, err = webElem.GetProperty(ctx, "id")
	assert.Error(t, err)
	assert.Empty(t, val)
}

func TestWebElement_Text(t *testing.T) {
	webElem, cli, done := newWebElement(t, "123")
	defer done()

	ctx := context.Background()

	// returns success
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/text", nil).Times(1).Return(
		&Response{
			Value: []byte(`"hoho"`),
		}, nil)
	val, err := webElem.Text(ctx)
	assert.Nil(t, err)
	assert.Equal(t, "hoho", val)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/text", nil).Times(1).Return(
		nil, elementWebErr)
	val, err = webElem.Text(ctx)
	assert.Error(t, err)
	assert.Empty(t, val)

	// returns empty
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/text", nil).Times(1).Return(
		&Response{
			Value: []byte(`null`),
		}, nil)
	val, err = webElem.Text(ctx)
	assert.Nil(t, err)
	assert.Empty(t, val)

	// returns errors JSON unmarshal
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/text", nil).Times(1).Return(
		&Response{
			Value: []byte(`"`),
		}, nil)
	val, err = webElem.Text(ctx)
	assert.Error(t, err)
	assert.Empty(t, val)
}

func TestWebElement_IsEnabled(t *testing.T) {
	webElem, cli, done := newWebElement(t, "123")
	defer done()

	ctx := context.Background()

	// returns success
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/enabled", nil).Times(1).Return(
		&Response{
			Value: []byte(`true`),
		}, nil)
	flag, err := webElem.IsEnabled(ctx)
	assert.Nil(t, err)
	assert.True(t, flag)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/enabled", nil).Times(1).Return(
		nil, elementWebErr)
	flag, err = webElem.IsEnabled(ctx)
	assert.Error(t, err)
	assert.False(t, flag)

	// returns error JSON unmarshal
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/enabled", nil).Times(1).Return(
		&Response{
			Value: []byte(`"`),
		}, nil)
	flag, err = webElem.IsEnabled(ctx)
	assert.Error(t, err)
	assert.False(t, flag)
}

func TestWebElement_IsSelected(t *testing.T) {
	webElem, cli, done := newWebElement(t, "123")
	defer done()

	ctx := context.Background()

	// returns success
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/selected", nil).Times(1).Return(
		&Response{
			Value: []byte(`true`),
		}, nil)
	flag, err := webElem.IsSelected(ctx)
	assert.Nil(t, err)
	assert.True(t, flag)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/selected", nil).Times(1).Return(
		nil, elementWebErr)
	flag, err = webElem.IsSelected(ctx)
	assert.Error(t, err)
	assert.False(t, flag)

	// returns error JSON unmarshal
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/selected", nil).Times(1).Return(
		&Response{
			Value: []byte(`{`),
		}, nil)
	flag, err = webElem.IsSelected(ctx)
	assert.Error(t, err)
	assert.False(t, flag)
}

func TestWebElement_TagName(t *testing.T) {
	webElem, cli, done := newWebElement(t, "123")
	defer done()

	ctx := context.Background()

	// returns success
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/name", nil).Times(1).Return(
		&Response{
			Value: []byte(`style`),
		}, nil)
	val, err := webElem.TagName(ctx)
	assert.Nil(t, err)
	assert.Equal(t, "style", val)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/name", nil).Times(1).Return(
		nil, elementWebErr)
	val, err = webElem.TagName(ctx)
	assert.Error(t, err)
	assert.Empty(t, val)

	// returns null
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/name", nil).Times(1).Return(
		&Response{
			Value: []byte(`null`),
		}, nil)
	val, err = webElem.TagName(ctx)
	assert.Empty(t, val)
	assert.Nil(t, err)
}

func TestWebElement_Rect(t *testing.T) {
	webElem, cli, done := newWebElement(t, "123")
	defer done()

	ctx := context.Background()

	// returns success
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/rect", nil).Times(1).Return(
		&Response{
			Value: []byte(`{"height":300,"width":645,"x":0,"y":0}`),
		}, nil)
	val, err := webElem.Rect(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 300, val.Height)
	assert.Equal(t, 645, val.Width)
	assert.Equal(t, 0, val.X)
	assert.Equal(t, 0, val.Y)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/rect", nil).Times(1).Return(
		nil, elementWebErr)
	val, err = webElem.Rect(ctx)
	assert.Error(t, err)
	assert.Equal(t, 0, val.Height)
	assert.Equal(t, 0, val.Width)
	assert.Equal(t, 0, val.X)
	assert.Equal(t, 0, val.Y)

	// returns error JSON unmarshal
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/rect", nil).Times(1).Return(
		&Response{
			Value: []byte(`"`),
		}, nil)
	val, err = webElem.Rect(ctx)
	assert.Error(t, err)
}

func TestWebElement_Clear(t *testing.T) {
	webElem, cli, done := newWebElement(t, "123")
	defer done()

	ctx := context.Background()

	// returns success
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/element/"+testWebElementID+"/clear", nil).Times(1).Return(
		&Response{
			Value: []byte(`null`),
		}, nil)
	err := webElem.Clear(ctx)
	assert.Nil(t, err)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/element/"+testWebElementID+"/clear", nil).Times(1).Return(
		nil, elementWebErr)
	err = webElem.Clear(ctx)
	assert.Error(t, err)

	// returns error invalid response
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/element/"+testWebElementID+"/clear", nil).Times(1).Return(
		&Response{
			Value: []byte(`{}`),
		}, nil)
	err = webElem.Clear(ctx)
	assert.Equal(t, ErrInvalidResponse, err)
}

func TestWebElement_SendKeys(t *testing.T) {
	webElem, cli, done := newWebElement(t, "123")
	defer done()

	ctx := context.Background()

	// returns success
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/element/"+testWebElementID+"/value", gomock.Any()).Times(1).Return(
		&Response{
			Value: []byte(`null`),
		}, nil).Do(func(_ context.Context, method, path string, p Params) {
		assert.Equal(t, string(PageDownKey)+" "+string(PageUpKey), p["text"])
	})
	err := webElem.SendKeys(ctx, PageDownKey, PageUpKey)
	assert.Nil(t, err)

	// returns error invalid args
	err = webElem.SendKeys(ctx)
	assert.Equal(t, ErrInvalidArguments, err)

	// returns error
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/element/"+testWebElementID+"/value", gomock.Any()).Times(1).Return(
		nil, elementWebErr)
	err = webElem.SendKeys(ctx, PageDownKey, PageUpKey)
	assert.Error(t, err)

	// returns error invalid response
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/element/"+testWebElementID+"/value", gomock.Any()).Times(1).Return(
		&Response{
			Value: []byte(`{}`),
		}, nil)
	err = webElem.SendKeys(ctx, PageDownKey, PageUpKey)
	assert.Error(t, err)
}
