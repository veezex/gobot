package emitter

import "github.com/vzxw/gobot/internal/pkg/message"

type MsgEmitter interface {
	Events() (<-chan message.Message, error)
}
