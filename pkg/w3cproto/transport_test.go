package w3cproto

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func TestTransport_DoEmptyRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	client, transport := newHttpClient(t, ctrl)
	request := WithClient(client)
	ctx := context.Background()
	transport.EXPECT().Do(gomock.Any()).Times(1).Return(&http.Response{
		Body:       makeBody(`{"value": null}`),
		StatusCode: http.StatusOK,
	}, nil).Do(func(req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		assert.Nil(t, err)
		assert.Equal(t, "{}", string(body))
		assert.Equal(t, "/session/123/window/new", req.URL.Path)
		assert.Equal(t, http.MethodPost, req.Method)
	})
	resp, err := request.Do(ctx, http.MethodPost, "/session/123/window/new", nil)
	assert.Nil(t, err)
	assert.True(t, resp.Success())
}
