package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nlopes/slack"
)

func callback(w http.ResponseWriter, r *http.Request) {
	var respMsg slack.Message

	verifiedBody, err := verifyCallbackMsg(r)
	if err != nil {
		log.Printf("failed to verify callback message: %s", err)
		return
	}

	msg, callbackID := verifiedBody.OriginalMessage, verifiedBody.CallbackID

	switch callbackID {
	case userInputID:
		action := verifiedBody.Actions[0]
		switch action.Name {
		case selectAction:
			messageText := fmt.Sprintf("Participants: <@%s>", action.SelectedOptions[0].Value)
			sendResp(w, fmtActionMsgResp(msg, messageText, postSelectUserAttachText, selectActions))
			return
		case additionalUserAction:
			messageText := fmt.Sprintf("Participants: %s, <@%s>", msg.Text[13:], action.SelectedOptions[0].Value)
			sendResp(w, fmtActionMsgResp(msg, messageText, "", selectActions))
			return
		case cancelAction:
			respMsg.DeleteOriginal = true
			respMsg.Text = requestCancelledText
			sendResp(w, respMsg)
			return
		}
	}

	sendResp(w, respMsg)
}

func help(w http.ResponseWriter, r *http.Request) {
	var respMsg slack.Message

	_, err := verifySlashCommand(r)
	if err != nil {
		log.Printf("failed to verify slash command: %s", err)
		return
	}

	respMsg.Text = "Set up a help message for users at this endpoint"

	sendResp(w, respMsg)
}
