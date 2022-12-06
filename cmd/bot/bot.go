package main

import (
	"context"
	"fmt"

	"github.com/vzxw/gobot/internal/pkg/emitter/slack"

	zeroLog "github.com/rs/zerolog/log"
	"github.com/vzxw/gobot/internal/pkg/config"
)

func main() {
	settings := config.Read(".env")

	// t := telegram.New(settings.TelegramAuthToken)
	s := slack.New(settings.SlackAppToken, slack.SlackEventsOps{
		Path: settings.SlackEventsPath,
		Port: settings.SlackEventsPort,
	})
	ctx := context.Background()

	msgChannel, err := s.Listen(ctx)
	if err != nil {
		zeroLog.Fatal().Err(err)
	}

	for msg := range msgChannel {
		if msg.Err != nil {
			zeroLog.Fatal().Err(msg.Err)
		}
		fmt.Println("Message:", msg.Text)
	}
}
