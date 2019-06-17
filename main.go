package main

import (
	"log"
	"net/http"
	"os"
	"time"

	utils "github.com/alyosha/slack-utils"
	"github.com/go-chi/chi"
	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

const (
	exitOK = iota
	exitError
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
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initalize zap logger: %s", err)
	}
	defer logger.Sync()

	logger.Info("starting up")

	var env config
	if err := envconfig.Process("", &env); err != nil {
		logger.Error("error processing environment variables", zap.Error(err))
		return exitError
	}

	client := slack.New(env.BotToken)
	cache := cache.New(15*time.Minute, 30*time.Minute)

	l := listener{
		client: client,
		cache:  cache,
		logger: logger,
		botID:  env.BotID,
	}
	h := handler{
		client: client,
		cache:  cache,
		logger: logger,
	}

	go l.listen()

	r := chi.NewRouter()
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := utils.WithSigningSecret(r.Context(), env.SigningSecret)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	r.Route("/", func(r chi.Router) {
		r.Post("/callback", h.callback)
		r.Post("/help", h.help)
	})

	logger.Info("server listening", zap.String("port", env.Port))
	if err := http.ListenAndServe(":"+env.Port, r); err != nil {
		logger.Error("failed to start http server", zap.Error(err))
		return exitError
	}

	return exitOK
}
