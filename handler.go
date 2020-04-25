package main

import (
	"fmt"
	"net/http"

	utils "github.com/alyosha/slack-utils"
	"github.com/patrickmn/go-cache"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type handler struct {
	client *slack.Client
	cache  *cache.Cache
	logger *zap.Logger
}

type request struct {
	users []string
}

func (h *handler) addUsers(w http.ResponseWriter, r *http.Request) {
	cmd, err := utils.VerifySlashCmd(r)
	if err != nil {
		h.logger.Error("failed to verify slash command", zap.Error(err))
		return
	}

	h.cache.Set(cmd.UserID, request{}, cache.DefaultExpiration)
	if _, err := utils.PostMsg(h.client, startMsg, cmd.ChannelID); err != nil {
		h.logger.Error("failed to handle message event", zap.Error(err))
	}
}

func (h *handler) callback(w http.ResponseWriter, r *http.Request) {
	callback, err := utils.VerifyCallbackMsg(r)
	if err != nil {
		h.logger.Error("failed to verify callback message", zap.Error(err))
		return
	}

	// We only have BlockAction callbacks in this application, so there is no risk of
	// index out of range here. In larger projects, implement a switch on callback type.
	action := callback.ActionCallback.BlockActions[0]

	switch action.ActionID {
	case utils.CancelActionID:
		h.cache.Delete(callback.User.ID)
		if err := utils.UpdateMsg(h.client, reqCancelledMsg, callback.Channel.ID, callback.Message.Timestamp); err != nil {
			h.logger.Error("failed to send request cancelled message", zap.Error(err))
		}
	case selectActionID:
		entry, found := h.cache.Get(callback.User.ID)
		if !found {
			if err := utils.UpdateMsg(h.client, reqExpiredMsg, callback.Channel.ID, callback.Message.Timestamp); err != nil {
				h.logger.Error("failed to send request expired message", zap.Error(err))
			}
			return
		}
		req := entry.(request)
		req.users = getUpdatedUsers(req.users, action.SelectedUser)
		if err := utils.UpdateMsg(h.client, fmtRespMsg(req), callback.Channel.ID, callback.Message.Timestamp); err != nil {
			h.logger.Error("failed to update message", zap.Error(err))
			return
		}
		h.cache.Set(callback.User.ID, req, cache.DefaultExpiration)
	}
}

func (h *handler) help(w http.ResponseWriter, r *http.Request) {
	var respMsg slack.Message

	_, err := utils.VerifySlashCmd(r)
	if err != nil {
		h.logger.Error("failed to verify slash command", zap.Error(err))
		return
	}

	respMsg.Text = "Set up a help message for users at this endpoint"

	utils.SendResp(w, respMsg)
}

func (h *handler) ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Bot is alive")
}

func getUpdatedUsers(existingUsers []string, newUser string) []string {
	var updatedUsers []string
	var alreadyPresent bool
	for _, user := range existingUsers {
		if user != newUser {
			updatedUsers = append(updatedUsers, user)
			continue
		}
		alreadyPresent = true
	}

	if !alreadyPresent {
		updatedUsers = append(updatedUsers, newUser)
	}

	return updatedUsers
}
