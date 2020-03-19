package w3c

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

var elementsErr = &Error{
	Code:    "code",
	Message: "msg",
}

func newElement(t *testing.T, sessID string) (Elements, *MockDoer, func()) {
	ctrl := gomock.NewController(t)
	cli := NewMockDoer(ctrl)
	cx := NewElements(cli, sessID)
	return cx, cli, func() {
		ctrl.Finish()
	}
}

func TestElements_Active(t *testing.T) {
	elem, cli, done := newElement(t, "123")
	defer done()

	ctx := context.TODO()
	want := "73101597-492f-4ffe-8f75-bd7bd0acb691"

	// returns success
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/active", gomock.Any()).Times(1).Return(&Response{
		Value: []byte(`{"element-6066-11e4-a52e-4f735466cecf":"73101597-492f-4ffe-8f75-bd7bd0acb691"}`),
	}, nil).Do(func(_ context.Context, path string, method string, p Params) {
		assert.Nil(t, p)
	})
	webElem, err := elem.Active(ctx)
	assert.Nil(t, err)
	assert.Equal(t, want, webElem.ID())

	// returns error
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/active", gomock.Any()).Times(1).Return(nil, elementsErr)
	webElem, err = elem.Active(ctx)
	assert.Error(t, err)
	assert.Equal(t, elementsErr, err)
	assert.Nil(t, webElem)

	// returns error unmarshal
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/active", gomock.Any()).Times(1).Return(&Response{
		Value: []byte(`{element-6066-11e4-a52e-4f735466cecf":"73101597-492f-4ffe-8f75-bd7bd0acb691"}`),
	}, nil).Do(func(_ context.Context, path string, method string, p Params) {
		assert.Nil(t, p)
	})
	webElem, err = elem.Active(ctx)
	assert.Error(t, err)
	assert.Nil(t, webElem)

	// returns error no such element
	cli.EXPECT().Do(ctx, http.MethodGet, "/session/123/element/active", gomock.Any()).Times(1).Return(&Response{
		Value: []byte(`{"element":"73101597-492f-4ffe-8f75-bd7bd0acb691"}`),
	}, nil).Do(func(_ context.Context, path string, method string, p Params) {
		assert.Nil(t, p)
	})
	webElem, err = elem.Active(ctx)
	assert.Error(t, err)
	assert.Equal(t, ErrNoSuchElement, err)
	assert.Nil(t, webElem)
}

func TestElements_Find(t *testing.T) {
	elem, cli, done := newElement(t, "123")
	defer done()

	ctx := context.TODO()
	want := "73101597-492f-4ffe-8f75-bd7bd0acb691"

	// returns success
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/elements", gomock.Any()).Times(1).Return(&Response{
		Value: []byte(`[{"element-6066-11e4-a52e-4f735466cecf":"73101597-492f-4ffe-8f75-bd7bd0acb691"}]`),
	}, nil).Do(func(_ context.Context, path string, method string, p Params) {
		assert.Equal(t, ByCSSSelector, p["using"])
		assert.Equal(t, ".className", p["value"])
	})
	webElems, err := elem.Find(ctx, ByCSSSelector, ".className")
	assert.Nil(t, err)
	assert.Len(t, webElems, 1)
	assert.Equal(t, want, webElems[0].ID())

	// returns error
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/elements", gomock.Any()).Times(1).Return(&Response{}, elementsErr)
	webElems, err = elem.Find(ctx, ByCSSSelector, ".className")
	assert.Error(t, err)
	assert.Equal(t, elementsErr, err)
	assert.Nil(t, webElems)

	// returns error no such element
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/elements", gomock.Any()).Times(1).Return(&Response{
		Value: []byte(`[{"element":"73101597-492f-4ffe-8f75-bd7bd0acb691"}]`),
	}, nil)
	webElems, err = elem.Find(ctx, ByCSSSelector, ".className")
	assert.Error(t, err)
	assert.Equal(t, ErrNoSuchElement, err)
	assert.Nil(t, webElems)

	// returns error invalid arguments
	webElems, err = elem.Find(ctx, ByCSSSelector, "")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidArguments, err)
	assert.Nil(t, webElems)

	// returns error unmarshal JSON
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/elements", gomock.Any()).Times(1).Return(&Response{
		Value: []byte(`[{element":"73101597-492f-4ffe-8f75-bd7bd0acb691"}]`),
	}, nil)
	webElems, err = elem.Find(ctx, ByCSSSelector, "#id")
	assert.Error(t, err)
	assert.Nil(t, webElems)
}

func TestElements_FindOne(t *testing.T) {
	elem, cli, done := newElement(t, "123")
	defer done()

	ctx := context.TODO()
	want := "73101597-492f-4ffe-8f75-bd7bd0acb691"

	// returns success
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/element", gomock.Any()).Times(1).Return(&Response{
		Value: []byte(`{"element-6066-11e4-a52e-4f735466cecf":"73101597-492f-4ffe-8f75-bd7bd0acb691"}`),
	}, nil).Do(func(_ context.Context, path string, method string, p Params) {
		assert.Equal(t, ByCSSSelector, p["using"])
		assert.Equal(t, "#id", p["value"])
	})
	webElem, err := elem.FindOne(ctx, ByCSSSelector, "#id")
	assert.Nil(t, err)
	assert.Equal(t, want, webElem.ID())

	// returns error
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/element", gomock.Any()).Times(1).Return(&Response{}, elementsErr)
	webElem, err = elem.FindOne(ctx, ByCSSSelector, "#id")
	assert.Error(t, err)
	assert.Equal(t, elementsErr, err)
	assert.Nil(t, webElem)

	// returns error no such element
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/element", gomock.Any()).Times(1).Return(&Response{
		Value: []byte(`{"element":"73101597-492f-4ffe-8f75-bd7bd0acb691"}`),
	}, nil)
	webElem, err = elem.FindOne(ctx, ByCSSSelector, "#id")
	assert.Error(t, err)
	assert.Equal(t, ErrNoSuchElement, err)
	assert.Nil(t, webElem)

	// returns error invalid arguments
	webElem, err = elem.FindOne(ctx, ByCSSSelector, "")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidArguments, err)
	assert.Nil(t, webElem)

	// returns error unmarshal JSON
	cli.EXPECT().Do(ctx, http.MethodPost, "/session/123/element", gomock.Any()).Times(1).Return(&Response{
		Value: []byte(`{element":"73101597-492f-4ffe-8f75-bd7bd0acb691"}`),
	}, nil)
	webElem, err = elem.FindOne(ctx, ByCSSSelector, "#id")
	assert.Error(t, err)
	assert.Nil(t, webElem)
}
