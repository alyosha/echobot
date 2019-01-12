package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/nlopes/slack"
)

type callbackHandler struct {
	slackClient   *slack.Client
	signingSecret string
}

func (h callbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var respMsg slack.Message

	msg, err := verifyMessage(r, h.signingSecret)
	if err != nil {
		return
	}

	message, callbackID := msg.OriginalMessage, msg.CallbackID

	switch callbackID {
	case userInputID:
		action := msg.Actions[0]
		switch action.Name {
		case "select":
			messageText := fmt.Sprintf("Participants: <@%s>", action.SelectedOptions[0].Value)
			sendResponse(w, formatActionMessageResponse(message, messageText, postSelectUserAttachText, selectActions))
			return
		case "additional_user":
			messageText := fmt.Sprintf("Participants: %s, <@%s>", message.Text[13:], action.SelectedOptions[0].Value)
			sendResponse(w, formatActionMessageResponse(message, messageText, "", selectActions))
			return
		case "cancel":
			respMsg.DeleteOriginal = true
			respMsg.Text = requestCancelledText
			sendResponse(w, respMsg)
			return
		}
	}

	sendResponse(w, respMsg)

}
func verifyMessage(req *http.Request, signingSecret string) (verifiedBody *slack.AttachmentActionCallback, err error) {
	if req.Method != http.MethodPost {
		log.Printf("Invalid method: %s, want POST", req.Method)
		return nil, err
	}

	sv, err := slack.NewSecretsVerifier(req.Header, signingSecret)
	if err != nil {
		log.Printf("Error initializing new SecretsVerifier: %s", err)
		return nil, err
	}

	var buf bytes.Buffer
	dest := io.MultiWriter(&buf, &sv)
	if _, err := io.Copy(dest, req.Body); err != nil {
		log.Printf("Error writing body to SecretsVerifier: %s", err)
		return nil, err
	}

	if err := sv.Ensure(); err != nil {
		log.Printf("Invalid signing secret: %s", err)
		return nil, err
	}

	jsonBody, err := url.QueryUnescape(buf.String()[8:])
	if err != nil {
		log.Printf("Error unescaping request body: %s", err)
		return
	}

	var message *slack.AttachmentActionCallback
	if err := json.Unmarshal([]byte(jsonBody), &message); err != nil {
		log.Printf("Error decoding JSON message from Slack: %s", err)
		return nil, err
	}

	return message, nil
}

func sendResponse(w http.ResponseWriter, message slack.Message) {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&message)
	return
}
