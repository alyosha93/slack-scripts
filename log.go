package utils

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/slack-go/slack"
)

func (c *Client) SendToLogChannel(msg Msg) error {
	_, err := c.PostMsg(msg, c.logChannel)
	if err != nil {
		return fmt.Errorf("c.PostMsg() > %w", err)
	}
	return nil
}

func (c *Client) SendToErrChannel(msgStr string, err error) error {
	var errMsgBody string

	if msgStr == "" {
		errMsgBody = fmt.Sprintf("`%v`", err)
	} else {
		errMsgBody = fmt.Sprintf("*%s*: `%v`", msgStr, err)
	}

	errMsg := Msg{
		Blocks: []slack.Block{
			NewTextBlock(errMsgBody, nil),
		},
	}

	threadMsg := Msg{
		Body: fmt.Sprintf("```\n%s\n```", string(debug.Stack())),
	}

	ts, err := c.PostMsg(errMsg, c.errChannel)
	if err != nil {
		return fmt.Errorf("c.PostMsg() > %w", err)
	}

	if err := c.PostThreadMsg(threadMsg, c.errChannel, ts); err != nil {
		return fmt.Errorf("c.PostThreadMsg() > %w", err)
	}

	return nil
}

func (c *Client) logRequest(cfg RequestLoggingConfig, endpoint, userID string) {
	if c.logChannel == "" || !cfg.Enabled || c.skipAdminLog(cfg.ExcludeAdmin, userID) {
		return
	}

	var logMsgBody string

	if cfg.MaskUserID {
		logMsgBody = fmt.Sprintf(
			"*endpoint:* `%s`\n*timestamp:* `%s`",
			endpoint,
			fmt.Sprintf("%d", time.Now().Unix()),
		)
	} else {
		logMsgBody = fmt.Sprintf(
			"*endpoint:* `%s`\n*user:* <@%s>\n*timestamp:* `%s`",
			endpoint,
			userID,
			fmt.Sprintf("%d", time.Now().Unix()),
		)
	}

	msg := Msg{
		Blocks: []slack.Block{
			NewTextBlock(logMsgBody, nil),
			DivBlock,
		},
	}

	if err := c.SendToLogChannel(msg); err != nil {
		_ = c.SendToErrChannel("failed to log request", err)
	}
}

func (c *Client) skipAdminLog(excludeAdmin bool, userID string) bool {
	if excludeAdmin && c.adminID == userID {
		return true
	}
	return false
}