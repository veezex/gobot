package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	zeroLog "github.com/rs/zerolog/log"

	"github.com/slack-go/slack/slackevents"

	slackContrib "github.com/slack-go/slack"

	"github.com/vzxw/gobot/internal/pkg/message"
)

type slack struct {
	secret string
	opts   SlackEventsOps
}

type SlackEventsOps struct {
	Port uint64
	Path string
}

func New(secret string, opts SlackEventsOps) *slack {
	return &slack{
		secret: secret,
		opts:   opts,
	}
}

func (s *slack) Listen(ctx context.Context) (<-chan message.Message, error) {
	result := make(chan message.Message, 10)

	http.HandleFunc(s.opts.Path, func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sv, err := slackContrib.NewSecretsVerifier(r.Header, s.secret)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if _, err := sv.Write(body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := sv.Ensure(); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		zeroLog.Info().Fields(eventsAPIEvent)
	})

	// run server
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.opts.Port), nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}
