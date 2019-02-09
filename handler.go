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
	signingSecret string
}

func (h callbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var respMsg slack.Message

	verifiedBody, err := verifyMessage(r, h.signingSecret)
	if err != nil {
		return
	}

	msg, callbackID := verifiedBody.OriginalMessage, verifiedBody.CallbackID

	switch callbackID {
	case userInputID:
		action := verifiedBody.Actions[0]
		switch action.Name {
		case "select":
			messageText := fmt.Sprintf("Participants: <@%s>", action.SelectedOptions[0].Value)
			sendResponse(w, formatActionMessageResponse(msg, messageText, postSelectUserAttachText, selectActions))
			return
		case "additional_user":
			messageText := fmt.Sprintf("Participants: %s, <@%s>", msg.Text[13:], action.SelectedOptions[0].Value)
			sendResponse(w, formatActionMessageResponse(msg, messageText, "", selectActions))
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

func sendResponse(w http.ResponseWriter, msg slack.Message) {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&msg)
	return
}
