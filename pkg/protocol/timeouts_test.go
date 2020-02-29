package protocol

import (
	"testing"
)

func TestTimeouts_GetTimeouts(t *testing.T) {
	//ctrl := gomock.NewController(t)
	//defer ctrl.Finish()
	//
	//cli := NewMockClient(ctrl)
	//sess := testNewSession(cli, t)
	//ctx := context.Background()
	//cli.EXPECT().Get(ctx, "/session/123/timeouts").Times(1).Return(&Response{
	//	Value: []byte(`{"script": 10000, "pageLoad": 10000, "implicit": 10000}`),
	//}, nil)
	//
	//ti, err := sess.GetTimeouts(ctx)
	//assert.Nil(t, err)
	//assert.Equal(t, ti.Implicit, DefaultTimeoutMs)
	//assert.Equal(t, ti.Script, DefaultTimeoutMs)
	//assert.Equal(t, ti.PageLoad, DefaultTimeoutMs)
}
