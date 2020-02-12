package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	"github.com/vvakame/sdlog/aelog"
	"go.opencensus.io/plugin/ochttp"
)

func main() {

	cfg := Environments()
	api := slack.New(cfg.SlackBotOAuthAccessToken)

	http.HandleFunc("/slack/bot", func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			aelog.Warningf(r.Context(), "err: %s", err)
			return
		}

		body := buf.String()

		aelog.Debugf(r.Context(), "body: %s", body)

		event, err := slackevents.ParseEvent(
			json.RawMessage(body),
			slackevents.OptionVerifyToken(
				&slackevents.TokenComparator{
					VerificationToken: cfg.SlackBotVerificationToken,
				},
			),
		)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			aelog.Warningf(r.Context(), "err: %s", err)
			w.Write([]byte(err.Error()))
			return
		}

		aelog.Debugf(r.Context(), "event.Type: %s", event.Type)

		switch event.Type {
		case slackevents.URLVerification:
			// for Event subscription setup
			urlVerificationEvent, ok := event.Data.(*slackevents.EventsAPIURLVerificationEvent)
			if !ok {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			b, err := json.Marshal(map[string]string{"challenge": urlVerificationEvent.Challenge})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				aelog.Warningf(r.Context(), "err: %s", err)
				w.Write([]byte(err.Error()))
				return
			}
			w.Write(b)
			return

		case slackevents.CallbackEvent:
			switch ev := event.InnerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				aelog.Debugf(r.Context(), "received: %s", ev.Text)
				_, _, err = api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
				if err != nil {
					w.WriteHeader(http.StatusForbidden)
					aelog.Warningf(r.Context(), "err: %s", err)
					w.Write([]byte(err.Error()))
					return
				}
			}
		}

		aelog.Debugf(r.Context(), "dump: %v", event)
	})

	log.Printf("Listening on port %s", cfg.Port)
	err := http.ListenAndServe(":"+cfg.Port, &ochttp.Handler{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := aelog.WithHTTPRequest(r.Context(), r)
			r = r.WithContext(ctx)
			http.DefaultServeMux.ServeHTTP(w, r)
		}),
		Propagation: &propagation.HTTPFormat{},
	})
	if err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	Port string `envconfig:"PORT" default:"8080"`

	SlackBotClientID          string `envconfig:"SLACK_BOT_CLIENT_ID" required:"true"`
	SlackBotClientSecret      string `envconfig:"SLACK_BOT_CLIENT_SECRET" required:"true"`
	SlackBotSigningSecret     string `envconfig:"SLACK_BOT_SIGNING_SECRET" required:"true"`
	SlackBotVerificationToken string `envconfig:"SLACK_BOT_VERIFICATION_TOKEN" required:"true"`
	SlackBotOAuthAccessToken  string `envconfig:"SLACK_BOT_OAUTH_ACCESS_TOKEN" required:"true"`
}

func Environments() *Config {
	cfg := &Config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
