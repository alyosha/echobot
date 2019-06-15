package main

import (
	"bytes"
	"fmt"

	utils "github.com/alyosha/slack-utils"
	"github.com/nlopes/slack"
)

const (
	selectBlockID     = "select_block"
	reqExpiredBlockID = "req_expired_block"
	selectActionID    = "select"
)

var (
	selectTxt       = slack.NewTextBlockObject(slack.PlainTextType, "Select", false, false)
	selectElem      = slack.NewOptionsSelectBlockElement(slack.OptTypeUser, selectTxt, selectActionID)
	startTxt        = slack.NewTextBlockObject(slack.PlainTextType, "Choose a member to get started", false, false)
	noUsersTxt      = slack.NewTextBlockObject(slack.MarkdownType, "*No users currently selected*", false, false)
	reqExpiredTxt   = slack.NewTextBlockObject(slack.PlainTextType, "It looks like your request has expired", false, false)
	reqCancelledTxt = slack.NewTextBlockObject(slack.PlainTextType, "Request cancelled", false, false)

	startSectionBlock        = slack.NewSectionBlock(startTxt, nil, nil)
	selectActionBlock        = slack.NewActionBlock(selectBlockID, selectElem, utils.CancelBtn)
	noUsersSectionBlock      = slack.NewSectionBlock(noUsersTxt, nil, nil)
	reqExpiredSectionBlock   = slack.NewSectionBlock(reqExpiredTxt, nil, nil)
	reqExpiredActionBlock    = slack.NewActionBlock(reqExpiredBlockID, utils.AckBtn)
	reqCancelledSectionBlock = slack.NewSectionBlock(reqCancelledTxt, nil, nil)

	startMsg = utils.Msg{
		Blocks: []slack.Block{startSectionBlock, selectActionBlock},
	}
	reqExpiredMsg = utils.Msg{
		Blocks: []slack.Block{reqExpiredSectionBlock},
	}
	reqCancelledMsg = utils.Msg{
		Blocks: []slack.Block{reqCancelledSectionBlock},
	}
)

func fmtRespMsg(req request) utils.Msg {
	userCount := len(req.users)

	if userCount == 0 {
		return utils.Msg{
			Blocks: []slack.Block{noUsersSectionBlock, selectActionBlock},
		}
	}

	var buf bytes.Buffer
	buf.WriteString("Selected users: ")
	for i, user := range req.users {
		buf.WriteString(fmt.Sprintf("<@%s>", user))
		if i+1 != userCount {
			buf.WriteString(",")
		}
	}

	respTxt := slack.NewTextBlockObject(
		slack.MarkdownType,
		fmt.Sprintf("Current user count: *%d*\n%s", userCount, buf.String()),
		false,
		false,
	)
	respSectionBlock := slack.NewSectionBlock(respTxt, nil, nil)

	return utils.Msg{
		Blocks: []slack.Block{respSectionBlock, selectActionBlock},
	}
}
