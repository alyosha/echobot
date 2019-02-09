package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/nlopes/slack"
)

type listener struct {
	client *slack.Client
	botID  string
}

// listen waits for message events
func (s *listener) listen() {
	rtm := s.client.NewRTM()

	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch event := msg.Data.(type) {
		case *slack.MessageEvent:
			if s.isValidMsg(event) {
				if _, err := s.postMsg(startMsg, event.Channel); err != nil {
					log.Printf("problem handling message event: %s", err)
				}
			}
		}
	}
}

func (s *listener) isValidMsg(event *slack.MessageEvent) bool {
	if s.botID == "" {
		log.Printf("received the following message: %s", event.Msg.Text)
		return false
	}
	msg := splitMsg(event.Msg.Text)
	if len(msg) == 0 || msg[0] != fmt.Sprintf("<@%s>", s.botID) {
		return false
	}

	return true
}

func splitMsg(msg string) []string {
	return strings.Split(strings.TrimSpace(msg), " ")[0:]
}
