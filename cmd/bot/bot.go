package main

import (
	zeroLog "github.com/rs/zerolog/log"
	"github.com/vzxw/gobot/internal/pkg/listener/telegram"

	"github.com/vzxw/gobot/internal/pkg/emitter/slack"

	"github.com/vzxw/gobot/internal/pkg/config"
)

func main() {
	settings := config.Read(".env")

	emitter := slack.New(settings.SlackSigningSecret, slack.EventOpts{
		Path: settings.SlackEventsPath,
		Port: settings.SlackEventsPort,
	})

	listener := telegram.New(settings.TelegramAuthToken)

	err := listener.Listen(emitter)
	if err != nil {
		zeroLog.Fatal().Err(err)
	}
}
