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
	opts   EventOpts
}

type EventOpts struct {
	Port uint64
	Path string
}

func New(secret string, opts EventOpts) *slack {
	return &slack{
		secret: secret,
		opts:   opts,
	}
}

func (s *slack) Listen(ctx context.Context) (<-chan message.Message, error) {
	result := make(chan message.Message, 10)

	finish := func(err error) {
		result <- message.Message{
			Author: "",
			Text:   "",
			Err:    err,
		}
		close(result)
	}

	http.HandleFunc(s.opts.Path, func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			finish(err)
			return
		}
		sv, err := slackContrib.NewSecretsVerifier(r.Header, s.secret)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			finish(err)
			return
		}
		if _, err := sv.Write(body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			finish(err)
			return
		}
		if err := sv.Ensure(); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			finish(err)
			return
		}
		eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			finish(err)
			return
		}

		zeroLog.Info().Fields(eventsAPIEvent)
	})

	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", s.opts.Port), nil)
		if err != nil {
			finish(err)
		}
	}()

	return result, nil
}
