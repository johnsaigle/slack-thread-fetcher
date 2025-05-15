// pkg/slack/fetcher.go
package slack

import (
	"fmt"
	"os"
	"strconv"
	"time"

	slackapi "github.com/slack-go/slack"
)

// SlackThreadFetcher handles communication with the Slack API
type SlackThreadFetcher struct {
	client *slackapi.Client
}

// Message represents a message in a Slack thread
type Message struct {
	UserID    string
	UserName  string
	RealName  string
	Email     string
	Text      string
	Timestamp string
	IsParent  bool
}

// NewSlackThreadFetcher creates a new instance of SlackThreadFetcher
func NewSlackThreadFetcher(token string) *SlackThreadFetcher {
	if token == "" {
		token = os.Getenv("SLACK_USER_TOKEN")
		if token == "" {
			fmt.Println("Error: Slack token is required. Set SLACK_USER_TOKEN environment variable or pass token using --token flag.")
			os.Exit(1)
		}
	}
	return &SlackThreadFetcher{
		client: slackapi.New(token),
	}
}

// GetThreadReplies retrieves all messages in a thread
func (f *SlackThreadFetcher) GetThreadReplies(channelID string, threadTS string) ([]Message, error) {
	params := slackapi.GetConversationRepliesParameters{
		ChannelID: channelID,
		Timestamp: threadTS,
	}

	conv, _, _, err := f.client.GetConversationReplies(&params)
	if err != nil {
		return nil, fmt.Errorf("error fetching thread: %v", err)
	}

	var messages []Message
	for _, msg := range conv {
		userInfo, err := f.getUserInfo(msg.User)
		if err != nil {
			userInfo = map[string]string{
				"name":      "Unknown User",
				"real_name": "Unknown User",
				"email":     "",
			}
		}

		messages = append(messages, Message{
			UserID:    msg.User,
			UserName:  userInfo["name"],
			RealName:  userInfo["real_name"],
			Email:     userInfo["email"],
			Text:      msg.Text,
			Timestamp: formatTimestamp(msg.Timestamp),
			IsParent:  msg.Timestamp == threadTS,
		})
	}

	return messages, nil
}

// getUserInfo retrieves information about a Slack user
func (f *SlackThreadFetcher) getUserInfo(userID string) (map[string]string, error) {
	user, err := f.client.GetUserInfo(userID)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"name":      user.Name,
		"real_name": user.RealName,
		"email":     user.Profile.Email,
	}, nil
}

// formatTimestamp converts a Slack timestamp to a human-readable format
func formatTimestamp(ts string) string {
	tsFloat, err := strconv.ParseFloat(ts, 64)
	if err != nil {
		return ts
	}
	t := time.Unix(int64(tsFloat), 0)
	return t.Format("2006-01-02 15:04:05")
}
