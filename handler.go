package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alyosha/slack-utils"
	"github.com/nlopes/slack"
)

func callback(w http.ResponseWriter, r *http.Request) {
	var respMsg slack.Message

	verifiedBody, err := utils.VerifyCallbackMsg(r)
	if err != nil {
		log.Fatalf("failed to verify callback message: %s", err)
	}

	msg, callbackID := verifiedBody.OriginalMessage, verifiedBody.CallbackID

	switch callbackID {
	case userInputID:
		action := verifiedBody.Actions[0]
		switch action.Name {
		case selectAction:
			messageText := fmt.Sprintf("Participants: <@%s>", action.SelectedOptions[0].Value)
			utils.SendResp(w, fmtActionMsgResp(msg, messageText, postSelectUserAttachText, selectActions))
			return
		case additionalUserAction:
			messageText := fmt.Sprintf("Participants: %s, <@%s>", msg.Text[13:], action.SelectedOptions[0].Value)
			utils.SendResp(w, fmtActionMsgResp(msg, messageText, "", selectActions))
			return
		case cancelAction:
			respMsg.DeleteOriginal = true
			respMsg.Text = requestCancelledText
			utils.SendResp(w, respMsg)
			return
		}
	}

	utils.SendResp(w, respMsg)
}

func help(w http.ResponseWriter, r *http.Request) {
	var respMsg slack.Message

	_, err := utils.VerifySlashCmd(r)
	if err != nil {
		log.Fatalf("failed to verify slash command: %s", err)
	}

	respMsg.Text = "Set up a help message for users at this endpoint"

	utils.SendResp(w, respMsg)
}
