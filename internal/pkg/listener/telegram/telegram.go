package telegram

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	zeroLog "github.com/rs/zerolog/log"
	"github.com/vzxw/gobot/internal/pkg/emitter"
	"github.com/vzxw/gobot/internal/pkg/message"
)

type telegram struct {
	bot    *tgbotapi.BotAPI
	token  string
	chatId *int64
}

func New(token string) *telegram {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	return &telegram{
		token:  token,
		bot:    bot,
		chatId: nil,
	}
}

func (t *telegram) runUpdatesLoop(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	go func() {
		updates := t.bot.GetUpdatesChan(u)
		for {
			select {
			case <-ctx.Done():
				return
			case update := <-updates:
				if update.Message != nil { // If we got a message
					if update.Message.Text == t.token {
						// update channel id
						chatId := update.Message.Chat.ID
						t.chatId = &chatId

						msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Updating chat id: %d", update.Message.Chat.ID))
						msg.ReplyToMessageID = update.Message.MessageID

						_, err := t.bot.Send(msg)
						if err != nil {
							zeroLog.Error().Err(err)
						}
					}
				}
			}
		}
	}()
}

func (t *telegram) sendMessage(m message.Message) {
	if t.chatId == nil {
		zeroLog.Log().Msg("Undefined chat id")
		return
	}

	msg := tgbotapi.NewMessage(*t.chatId, fmt.Sprintf("%s: %s", m.Author, m.Text))
	_, err := t.bot.Send(msg)
	if err != nil {
		zeroLog.Error().Err(err)
	}
}

func (t *telegram) Listen(source emitter.MsgEmitter) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.runUpdatesLoop(ctx)

	msgChannel, err := source.Events()
	if err != nil {
		return errors.Wrap(err, "Listener error")
	}

	for msg := range msgChannel {
		if msg.Err != nil {
			zeroLog.Error().Err(msg.Err)
		} else {
			t.sendMessage(msg)
		}
	}

	return nil
}
