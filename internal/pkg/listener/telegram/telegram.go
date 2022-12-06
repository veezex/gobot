package telegram

import "github.com/vzxw/gobot/internal/pkg/emitter"

type telegram struct{}

func New(token string) *telegram {
	return &telegram{}
}

func (*telegram) ListenToSource(source emitter.MsgEmitter) {}
