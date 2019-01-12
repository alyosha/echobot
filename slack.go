package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

type slackHandler struct {
	client *slack.Client
	botID  string
}

// listen waits for message events
func (s *slackHandler) listen() {
	rtm := s.client.NewRTM()

	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch event := msg.Data.(type) {
		case *slack.MessageEvent:
			if s.botID == "" {
				log.Printf("Received the following message: %s", event.Msg.Text)
			}
			m := strings.Split(strings.TrimSpace(event.Msg.Text), " ")[0:]
			if len(m) != 0 && m[0] == fmt.Sprintf("<@%s>", s.botID) {
				if _, err := s.postMessage(startMessage, event.Channel); err != nil {
					log.Printf("Problem handling message event: %s", err)
				}
			}
		}
	}
}

// postMessage sends the message provided in the method params to the channel designated
func (s *slackHandler) postMessage(msg message, channel string) (timestamp string, err error) {
	params := slack.PostMessageParameters{
		Attachments: msg.attachments,
	}
	_, ts, err := s.client.PostMessage(channel, msg.messageBody, params)
	if err != nil {
		return "", errors.Wrap(err, "Failed to post message to Slack")
	}

	return ts, nil
}
