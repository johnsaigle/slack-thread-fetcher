// pkg/slack/utils.go
package slack

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ParseSlackURL extracts channel ID and timestamp from a Slack thread URL
func ParseSlackURL(url string) (string, string, error) {
	// Decode URL to handle any percent encoding
	url = strings.Replace(url, "%3A", ":", -1)
	url = strings.Replace(url, "%2F", "/", -1)

	pattern := `archives/([A-Z0-9]+)/p(\d{10})(\d+)`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(url)

	if matches == nil || len(matches) != 4 {
		return "", "", fmt.Errorf("invalid Slack thread URL: %s", url)
	}

	channelID := matches[1]
	threadTS := fmt.Sprintf("%s.%s", matches[2], matches[3])

	return channelID, threadTS, nil
}

// FormatThreadContent formats the thread messages into a readable string
func FormatThreadContent(messages []Message) string {
	var content strings.Builder
	for _, msg := range messages {
		content.WriteString(fmt.Sprintf("%s (@%s) - %s\n%s\n\n", 
			msg.RealName, 
			msg.UserName, 
			msg.Timestamp, 
			msg.Text))
	}
	return content.String()
}

// SaveToFile writes thread contents to the specified file
func SaveToFile(outputFile string, contents map[string]string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// For each thread, write its content
	threadNum := 1
	for url, content := range contents {
		if threadNum > 1 {
			file.WriteString("\n\n" + strings.Repeat("=", 50) + "\n\n")
		}
		file.WriteString(fmt.Sprintf("=== Thread %d: %s ===\n\n", threadNum, url))
		file.WriteString(content)
		threadNum++
	}

	return nil
}
