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

	// returns error
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/"+testWebElementID+"/attribute/class", nil).Times(1).Return(
		nil, elementWebErr)
	val, err = webElem.GetAttribute(ctx, "class")
	assert.Error(t, err)
	assert.Equal(t, elementWebErr, err)
	assert.Empty(t, val)
}
