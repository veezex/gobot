package main

import (
	"context"
	"fmt"

	zeroLog "github.com/rs/zerolog/log"
	"github.com/vzxw/gobot/internal/pkg/config"
	"github.com/vzxw/gobot/internal/pkg/source/slack"
)

func main() {
	settings := config.Read(".env")

	// t := telegram.New(settings.TelegramAuthToken)
	s := slack.New(settings.SlackAppToken, settings.SlackBotToken)
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
