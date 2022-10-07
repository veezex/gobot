package telegram

import (
	"github.com/vzxw/gobot/internal/pkg/source"
)

type telegram struct{}

func New(token string) *telegram {
	return &telegram{}
}

func (*telegram) ListenToSource(source source.MsgSource) {}
