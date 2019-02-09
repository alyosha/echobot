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

func verifyMsg(req *http.Request, signingSecret string) (verifiedBody *slack.AttachmentActionCallback, err error) {
	if req.Method != http.MethodPost {
		log.Printf("invalid method: %s, want POST", req.Method)
		return nil, err
	}

	sv, err := slack.NewSecretsVerifier(req.Header, signingSecret)
	if err != nil {
		log.Printf("error initializing new SecretsVerifier: %s", err)
		return nil, err
	}

	var buf bytes.Buffer
	dest := io.MultiWriter(&buf, &sv)
	if _, err := io.Copy(dest, req.Body); err != nil {
		log.Printf("error writing body to SecretsVerifier: %s", err)
		return nil, err
	}

	if err := sv.Ensure(); err != nil {
		log.Printf("invalid signing secret: %s", err)
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
