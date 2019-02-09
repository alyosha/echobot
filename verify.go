package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/nlopes/slack"
)

func verifyCallbackMsg(r *http.Request) (verifiedBody *slack.AttachmentActionCallback, err error) {
	if r.Method != http.MethodPost {
		log.Printf("invalid method: %s, want POST", r.Method)
		return nil, err
	}

	buf, err := checkSecretAndWriteBody(r)
	if err != nil {
		return nil, err
	}

	jsonBody, err := url.QueryUnescape(buf.String()[8:])
	if err != nil {
		log.Printf("error unescaping request body: %s", err)
		return
	}

	var msg *slack.AttachmentActionCallback
	if err := json.Unmarshal([]byte(jsonBody), &msg); err != nil {
		log.Printf("error decoding JSON message from Slack: %s", err)
		return nil, err
	}

	return msg, nil
}

func verifySlashCommand(r *http.Request) (verifiedBody *slack.SlashCommand, err error) {
	if r.Method != http.MethodPost {
		log.Printf("invalid method: %s, want POST", r.Method)
		return nil, err
	}

	buf, err := checkSecretAndWriteBody(r)
	if err != nil {
		return nil, err
	}

	body, err := url.ParseQuery(string(buf.String()))
	if err != nil {
		log.Printf("error parsing query body: %s", err)
		return nil, err
	}

	msg := parseCommand(body)

	return &msg, nil
}

func checkSecretAndWriteBody(r *http.Request) (bytes.Buffer, error) {
	var buf bytes.Buffer
	signingSecret, err := getSigningSecret(r.Context())
	if err != nil {
		log.Printf("failed to extract signing secret from context %s", err)
		return buf, err
	}

	sv, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		log.Printf("error initializing new SecretsVerifier: %s", err)
		return buf, err
	}

	dest := io.MultiWriter(&buf, &sv)
	if _, err := io.Copy(dest, r.Body); err != nil {
		log.Printf("error writing body to SecretsVerifier: %s", err)
		return buf, err
	}

	if err := sv.Ensure(); err != nil {
		log.Printf("invalid signing secret: %s", err)
		return buf, err
	}

	return buf, nil
}

func parseCommand(body url.Values) (s slack.SlashCommand) {
	s.Token = body.Get("token")
	s.TeamID = body.Get("team_id")
	s.TeamDomain = body.Get("team_domain")
	s.EnterpriseID = body.Get("enterprise_id")
	s.EnterpriseName = body.Get("enterprise_name")
	s.ChannelID = body.Get("channel_id")
	s.ChannelName = body.Get("channel_name")
	s.UserID = body.Get("user_id")
	s.UserName = body.Get("user_name")
	s.Command = body.Get("command")
	s.Text = body.Get("text")
	s.ResponseURL = body.Get("response_url")
	s.TriggerID = body.Get("trigger_id")

	return s
}
