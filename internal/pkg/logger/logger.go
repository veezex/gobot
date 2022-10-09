package logger

import zeroLog "github.com/rs/zerolog/log"

type loggerInfo struct {
	prefix string
}

func NewInfo(prefix string) *loggerInfo {
	return &loggerInfo{
		prefix: prefix,
	}
}

func (li *loggerInfo) Output(_ int, msg string) error {
	zeroLog.Info().Msgf("%s: ", li.prefix, msg)
	return nil
}
