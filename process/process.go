package process

import (
	"bytes"
	"context"
)

type Process interface {
	Run(ctx context.Context) error
	Stop(ctx context.Context) error
	Out() bytes.Buffer
}
