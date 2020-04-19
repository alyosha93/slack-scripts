package utils

import (
	"encoding/json"
	"net/http"

	"github.com/slack-go/slack"
)

// Msg is an intermediary struct used for posting messages
type Msg struct {
	Body        string
	Blocks      []slack.Block
	Attachments []slack.Attachment
	AsUser      bool
	UserID      string // Only set if you want to post ephemerally
	ResponseURL string // Only set if you want to delete an ephemeral message
}

// PostMsg sends the provided message to the channel designated by channelID
func PostMsg(client *slack.Client, msg Msg, channelID string) (string, string, error) {
	opts := []slack.MsgOption{
		slack.MsgOptionText(msg.Body, false),
		slack.MsgOptionBlocks(msg.Blocks...),
		slack.MsgOptionAttachments(msg.Attachments...),
		slack.MsgOptionAsUser(msg.AsUser),
		slack.MsgOptionEnableLinkUnfurl(),
	}

	if msg.UserID != "" {
		opts = append(opts, slack.MsgOptionPostEphemeral(msg.UserID))
	}

	channelID, ts, err := client.PostMessage(
		channelID,
		opts...,
	)

	if err != nil {
		return "", "", err
	}

	return channelID, ts, nil
}

// PostThreadMsg posts a message response into an existing thread
func PostThreadMsg(client *slack.Client, msg Msg, channelID string, threadTs string) error {
	_, _, err := client.PostMessage(
		channelID,
		slack.MsgOptionText(msg.Body, false),
		slack.MsgOptionAttachments(msg.Attachments...),
		slack.MsgOptionEnableLinkUnfurl(),
		slack.MsgOptionTS(threadTs),
	)

	return err
}

// UpdateMsg updates the provided message in the channel designated by channelID
func UpdateMsg(client *slack.Client, msg Msg, channelID, timestamp string) (string, string, string, error) {
	opts := []slack.MsgOption{
		slack.MsgOptionText(msg.Body, false),
		slack.MsgOptionBlocks(msg.Blocks...),
		slack.MsgOptionAttachments(msg.Attachments...),
		slack.MsgOptionAsUser(msg.AsUser),
		slack.MsgOptionEnableLinkUnfurl(),
	}

	if msg.ResponseURL != "" {
		opts = []slack.MsgOption{slack.MsgOptionDeleteOriginal(msg.ResponseURL)}
	}

	channelID, ts, text, err := client.UpdateMessage(
		channelID,
		timestamp,
		opts...,
	)

	if err != nil {
		return "", "", "", err
	}

	return channelID, ts, text, nil
}

// SendEmptyOK responds with status 200
func SendEmptyOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	return
}

// SendResp can be used to send simple callback responses
// NOTE: cannot be used in callback from block messages
func SendResp(w http.ResponseWriter, msg slack.Message) error {
	w.Header().Add("Content-type", "application/json")
	return json.NewEncoder(w).Encode(&msg)
}

// ReplaceOriginal replaces the original message with the newly encoded one
// NOTE: cannot be used in callback from block messages
func ReplaceOriginal(w http.ResponseWriter, msg slack.Message) error {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	msg.ReplaceOriginal = true
	return json.NewEncoder(w).Encode(&msg)
}

// SendOKAndDeleteOriginal responds with status 200 and deletes the original message
// NOTE: cannot be used in callback from block messages
func SendOKAndDeleteOriginal(w http.ResponseWriter) error {
	var msg slack.Message
	msg.DeleteOriginal = true
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(&msg)
}
