package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
	cache "github.com/patrickmn/go-cache"
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
	fmt.Print("startup")
	var env config
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("Failed to process env var: %s", err)
		return 1
	}

	slackClient := slack.New(env.BotToken)
	slackHandler := &slackHandler{
		client: slackClient,
		botID:  env.BotID,
	}

	go slackHandler.listen()

	cache := cache.New(10*time.Minute, 30*time.Minute)

	http.Handle("/callback", callbackHandler{
		slackClient:   slackClient,
		signingSecret: env.SigningSecret,
		cache:         cache,
	})

	log.Printf("Server listening on :%s", env.Port)
	if err := http.ListenAndServe(":"+env.Port, nil); err != nil {
		log.Printf("Error: %s", err)
		return 1
	}

	return 0
}
