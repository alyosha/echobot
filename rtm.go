package main

import (
	"fmt"
	"log"
	"strings"

	utils "github.com/alyosha/slack-utils"
	"github.com/nlopes/slack"
)

type listener struct {
	client *slack.Client
	botID  string
}

// listen waits for message events
func (l *listener) listen() {
	rtm := l.client.NewRTM()

	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch event := msg.Data.(type) {
		case *slack.MessageEvent:
			if l.isBotImperative(event) {
				if _, err := utils.PostMsg(l.client, startMsg, event.Channel); err != nil {
					log.Printf("problem handling message event: %s", err)
				}
			}
		}
	}
}

func (l *listener) isBotImperative(event *slack.MessageEvent) bool {
	if l.botID == "" {
		log.Printf("received the following message: %s", event.Msg.Text)
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
