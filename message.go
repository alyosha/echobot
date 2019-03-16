package main

import (
	"github.com/alyosha/slack-utils"
	"github.com/nlopes/slack"
)

const (
	postSelectUserAttachText = "Okay, need at least one more participant!"
	requestCancelledText     = "Request cancelled!"
	userInputID              = "user_input"
	selectAction             = "select"
	additionalUserAction     = "additional_user"
	cancelAction             = "cancel"
)

// Sample response message featuring two interactive actions
var startMsg = utils.Msg{
	Attachments: []slack.Attachment{
		{
			Text:       "Who will be participating?",
			Color:      "#f9a41b",
			CallbackID: userInputID,
			Actions: []slack.AttachmentAction{
				{
					Name:       selectAction,
					Type:       "select",
					DataSource: "users",
				},
				{
					Name:  cancelAction,
					Text:  "Cancel",
					Type:  "button",
					Style: "danger",
				},
			},
		},
	},
	Body: "Let's get started!",
}

var selectActions = []slack.AttachmentAction{
	{
		Name:       additionalUserAction,
		Type:       "select",
		DataSource: "users",
	},
	{
		Name:  cancelAction,
		Text:  "Cancel",
		Type:  "button",
		Style: "danger",
	},
}

func fmtActionMsgResp(originalMsg slack.Message, msgText string, attachText string, actions []slack.AttachmentAction) slack.Message {
	originalMsg.ReplaceOriginal = true
	originalMsg.Text = msgText
	originalMsg.Attachments[0].Text = attachText
	originalMsg.Attachments[0].Actions = actions

	return originalMsg
}
