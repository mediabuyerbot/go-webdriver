package protocol

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func TestNewSessionWithoutW3CCompatibility(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli := NewMockClient(ctrl)
	cli.EXPECT().Post(context.Background(), "/session", gomock.Any()).
		Times(1).
		Return(&Response{
			SessionID: "123",
			Status:    0,
			Value:     json.RawMessage(`{}`),
		}, nil)

	desired := make(map[string]string)
	desired["platform"] = "linux"
	required := make(map[string]string)
	sess, err := NewSession(cli, desired, required)

	assert.Nil(t, err)
	assert.Equal(t, sess.ID(), SessionID("123"))
}
