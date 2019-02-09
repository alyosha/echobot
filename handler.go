package main

import (
	"fmt"
	"net/http"

	"github.com/nlopes/slack"
)

type callbackHandler struct {
	signingSecret string
}

func (h callbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var respMsg slack.Message

	verifiedBody, err := verifyMsg(r, h.signingSecret)
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
			sendResp(w, fmtActionMsgResp(msg, messageText, postSelectUserAttachText, selectActions))
			return
		case "additional_user":
			messageText := fmt.Sprintf("Participants: %s, <@%s>", msg.Text[13:], action.SelectedOptions[0].Value)
			sendResp(w, fmtActionMsgResp(msg, messageText, "", selectActions))
			return
		case "cancel":
			respMsg.DeleteOriginal = true
			respMsg.Text = requestCancelledText
			sendResp(w, respMsg)
			return
		}
	}

	sendResp(w, respMsg)

}
