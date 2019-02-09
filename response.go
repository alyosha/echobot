package main

import (
	"encoding/json"
	"net/http"

	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

// postMsg sends the message provided in the method params to the channel designated
func (s *listener) postMsg(msg message, channel string) (timestamp string, err error) {
	params := slack.PostMessageParameters{
		Attachments: msg.attachments,
	}
	_, ts, err := s.client.PostMessage(channel, msg.body, params)
	if err != nil {
		return "", errors.Wrap(err, "failed to post message to Slack")
	}

	return ts, nil
}

func sendResp(w http.ResponseWriter, msg slack.Message) {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&msg)
	return
}
