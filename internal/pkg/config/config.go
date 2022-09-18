package config

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type readerStruct struct {
	envs map[string]string
}

type list struct {
	SlackAuthToken string
}

func Read(envFile ...string) *list {
	envs, err := godotenv.Read(envFile...)

	if err != nil {
		log.Fatal().Err(err)
	}

	reader := &readerStruct{
		envs: envs,
	}

	return &list{
		SlackAuthToken: reader.getString("SLACK_BOT_AUTH_TOKEN"),
	}
}

func (r *readerStruct) getString(key string) string {
	result, ok := r.envs[key]
	if !ok {
		log.Fatal().Msgf("Undefined <%v> env variable", key)
	}

	return result
}
