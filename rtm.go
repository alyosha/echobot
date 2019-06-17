package main

import (
	"fmt"
	"strings"

	utils "github.com/alyosha/slack-utils"
	"github.com/nlopes/slack"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

type listener struct {
	client *slack.Client
	cache  *cache.Cache
	logger *zap.Logger
	botID  string
}

// listen waits for new message events in any channels the bot has been invited to
func (l *listener) listen() {
	rtm := l.client.NewRTM()

	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch event := msg.Data.(type) {
		case *slack.MessageEvent:
			if l.isBotImperative(event) {
				l.cache.Set(event.User, request{}, cache.DefaultExpiration)
				if _, _, err := utils.PostMsg(l.client, startMsg, event.Channel); err != nil {
					l.logger.Error("failed to handle message event", zap.Error(err))
				}
			}
		}
	}
}

// Only respond if the message begins with the bot's name
func (l *listener) isBotImperative(event *slack.MessageEvent) bool {
	if l.botID == "" {
		l.logger.Info("received message but BotID not set", zap.String("content", event.Msg.Text))
		return false
	}

	msg := splitMsg(event.Msg.Text)
	if len(msg) == 0 || msg[0] != fmt.Sprintf("<@%s>", l.botID) {
		return false
	}

	return true
}

func splitMsg(msg string) []string {
	return strings.Split(strings.TrimSpace(msg), " ")[0:]
}
