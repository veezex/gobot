package config

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type readerStruct struct {
	envs map[string]string
}

type List struct {
	SlackAuthToken    string
	TelegramAuthToken string
}

func Read(envFile ...string) List {
	envs, err := godotenv.Read(envFile...)
	if err != nil {
		log.Fatal().Err(err)
	}

	reader := &readerStruct{
		envs: envs,
	}

	return List{
		SlackAuthToken:    reader.getString("SLACK_BOT_AUTH_TOKEN"),
		TelegramAuthToken: reader.getString("TELEGRAM_BOT_AUTH_TOKEN"),
	}
}

func (r *readerStruct) getString(key string) string {
	result, ok := r.envs[key]
	if !ok {
		log.Fatal().Msgf("Undefined <%v> env variable", key)
	}

	return result
}
