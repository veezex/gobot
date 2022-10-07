package main

import (
	"github.com/vzxw/gobot/internal/pkg/config"
	"github.com/vzxw/gobot/internal/pkg/receiver/telegram"
	"github.com/vzxw/gobot/internal/pkg/source/slack"
)

func main() {
	settings := config.Read(".env")

	t := telegram.New(settings.TelegramAuthToken)
	s := slack.New(settings.SlackAuthToken)

	t.ListenToSource(s)
}

/*
https://github.com/trestoa/slack-to-telegram-bot
*/
