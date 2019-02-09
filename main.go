package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
)

type config struct {
	Port          string `envconfig:"PORT" default:"3000"`
	SigningSecret string `envconfig:"SIGNING_SECRET" required:"true"`
	BotToken      string `envconfig:"BOT_TOKEN" required:"true"`
	BotID         string `envconfig:"BOT_ID"`
}

func main() {
	os.Exit(_main())
}

func _main() int {
	log.Print("starting up")
	var env config
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("error processing environment variables: %s", err)
		return 1
	}

	slackClient := slack.New(env.BotToken)
	slackHandler := &listener{
		client: slackClient,
		botID:  env.BotID,
	}

	go slackHandler.listen()

	r := chi.NewRouter()
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := withContext(r.Context(), env.SigningSecret, slackClient)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	r.Route("/", func(r chi.Router) {
		r.Post("/callback", callback)
		r.Post("/help", help)
	})

	log.Printf("server listening on :%s", env.Port)
	if err := http.ListenAndServe(":"+env.Port, r); err != nil {
		log.Printf("error: %s", err)
		return 1
	}

	return 0
}
