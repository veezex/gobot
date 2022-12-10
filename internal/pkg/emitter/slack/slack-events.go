package slack

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

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

func (s *slack) Listen() (<-chan message.Message, error) {
	result := make(chan message.Message)

	errNotify := func(err error) {
		result <- message.Message{
			Author: "",
			Text:   "",
			Err:    err,
		}
	}

	http.HandleFunc(s.opts.Path, func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errNotify(err)

			return
		}
		sv, err := slackContrib.NewSecretsVerifier(r.Header, s.secret)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errNotify(err)

			return
		}
		if _, err := sv.Write(body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			errNotify(err)

			return
		}
		if err := sv.Ensure(); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			errNotify(err)

			return
		}
		eventsAPIEvent, err := slackevents.ParseEvent(body, slackevents.OptionNoVerifyToken())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			errNotify(err)

			return
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal(body, &r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				errNotify(err)

				return
			}
			w.Header().Set("Content-Type", "text")

			_, err = w.Write([]byte(r.Challenge))
			if err != nil {
				errNotify(err)

				return
			}
		}

		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			if eventsAPIEvent.InnerEvent.Type == "message" {
				if err != nil {
					errNotify(err)
					return
				}

				messageEvent, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.MessageEvent)
				if !ok {
					errNotify(errors.New("Incoming event can't be converted to events.MessageEvent"))
				}

				// eventsAPIEvent.InnerEvent.Data["Text"]
				result <- message.Message{
					Author: messageEvent.User,
					Text:   messageEvent.Text,
					Err:    nil,
				}
			}
		}
	})

	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", s.opts.Port), nil)
		if err != nil {
			errNotify(err)
			close(result)
		}
	}()

	return result, nil
}
