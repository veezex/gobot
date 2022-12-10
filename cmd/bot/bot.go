package main

import (
	"fmt"

	"github.com/vzxw/gobot/internal/pkg/emitter/slack"

	zeroLog "github.com/rs/zerolog/log"
	"github.com/vzxw/gobot/internal/pkg/config"
)

func main() {
	settings := config.Read(".env")

	// t := telegram.New(settings.TelegramAuthToken)
	s := slack.New(settings.SlackSigningSecret, slack.EventOpts{
		Path: settings.SlackEventsPath,
		Port: settings.SlackEventsPort,
	})

	msgChannel, err := s.Listen()
	if err != nil {
		zeroLog.Fatal().Err(err)
	}

	for msg := range msgChannel {
		if msg.Err != nil {
			zeroLog.Fatal().Err(msg.Err)
		} else {
			fmt.Println("Message:", msg.Text)
		}
	}
}
