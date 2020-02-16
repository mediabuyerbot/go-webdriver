package chromedriver

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
	ch := New()
	go func() {
		if err := ch.Run(ctx); err != nil {
			t.Log(ch.Out())
			t.Fatal(err)
		}
	}()

	<-time.After(time.Second * 2)

	if err := ch.Stop(ctx); err != nil {
		t.Log(ch.Out())
		t.Fatal(err)
	}
}
