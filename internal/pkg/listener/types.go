package listener

import (
	"context"

	"github.com/vzxw/gobot/internal/pkg/emitter"
)

type MsgListener interface {
	ListenToSource(ctx context.Context, source emitter.MsgEmitter)
}
