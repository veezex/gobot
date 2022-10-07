package receiver

import "github.com/vzxw/gobot/internal/pkg/source"

type MsgReceiver interface {
	ListenToSource(source source.MsgSource)
}
