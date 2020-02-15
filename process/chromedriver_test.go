package process

import (
	"context"
	"testing"
	"time"
)

func TestDefaultChromeDriverWithRunAndStop(t *testing.T) {
	if t.Skipped() {
		return
	}

	ctx := context.Background()
	ch := NewChromeDriver()
	go func() {
		if err := ch.Run(ctx); err != nil {
			t.Log(ch.Out())
			t.Fatal(err)
		}
	}()

	<-time.After(time.Second * 5)

	if err := ch.Stop(ctx); err != nil {
		t.Log(ch.Out())
		t.Fatal(err)
	}
}
