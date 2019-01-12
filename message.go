package main

import "github.com/nlopes/slack"

const (
	postSelectUserAttachText = "Okay, need at least one more participant!"
	requestCancelledText     = "Request cancelled!"
	userInputID              = "user_input"
)

type message struct {
	attachments []slack.Attachment
	messageBody string
}

// Sample response message featuring two interactive actions
var startMessage = message{
	attachments: []slack.Attachment{
		{
			Text:       "Who will be participating?",
			Color:      "#f9a41b",
			CallbackID: userInputID,
			Actions: []slack.AttachmentAction{
				{
					Name:       "select",
					Type:       "select",
					DataSource: "users",
				},
				{
					Name:  "cancel",
					Text:  "Cancel",
					Type:  "button",
					Style: "danger",
				},
			},
		},
	},
	messageBody: "Let's get started!",
}

var selectActions = []slack.AttachmentAction{
	{
		Name:       "additional_user",
		Type:       "select",
		DataSource: "users",
	},
	{
		Name:  "cancel",
		Text:  "Cancel",
		Type:  "button",
		Style: "danger",
	},
}

func formatActionMessageResponse(originalMessage slack.Message, messageText string, attachText string, actions []slack.AttachmentAction) slack.Message {
	originalMessage.ReplaceOriginal = true
	originalMessage.Text = messageText
	originalMessage.Attachments[0].Text = attachText
	originalMessage.Attachments[0].Actions = actions

	return originalMessage
}
