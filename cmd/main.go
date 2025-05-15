// cmd/main.go
package main

import (
	"fmt"
	"os"

	"github.com/johnsaigle/slack-thread-fetcher/pkg/slack"
	"github.com/spf13/cobra"
)

func main() {
	var token string
	var outputFile string
	
	var rootCmd = &cobra.Command{
		Use:   "slack-thread-fetcher [thread URLs...]",
		Short: "Fetch Slack thread comments",
		Long:  `A tool to fetch and save Slack thread comments from provided URLs.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fetcher := slack.NewSlackThreadFetcher(token)
			threadContents := make(map[string]string)
			
			for _, url := range args {
				channelID, threadTS, err := slack.ParseSlackURL(url)
				if err != nil {
					fmt.Printf("Error parsing URL %s: %v\n", url, err)
					continue
				}
				
				fmt.Printf("Fetching thread from channel %s, timestamp %s...\n", channelID, threadTS)
				messages, err := fetcher.GetThreadReplies(channelID, threadTS)
				if err != nil {
					fmt.Printf("Error fetching thread %s: %v\n", url, err)
					continue
				}
				
				threadContent := slack.FormatThreadContent(messages)
				threadContents[url] = threadContent
				fmt.Printf("Successfully fetched thread with %d messages\n", len(messages))
			}
			
			if len(threadContents) > 0 {
				err := slack.SaveToFile(outputFile, threadContents)
				if err != nil {
					fmt.Printf("Error saving to file: %v\n", err)
					os.Exit(1)
				}
				fmt.Printf("\nThread content saved to %s\n", outputFile)
			} else {
				fmt.Println("No threads were successfully fetched.")
			}
		},
	}
	
	rootCmd.Flags().StringVar(&token, "token", "", "Slack User Token (xoxp-). If not provided, uses SLACK_USER_TOKEN environment variable")
	rootCmd.Flags().StringVar(&outputFile, "output", "thread_content.txt", "Output file path")
	
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
