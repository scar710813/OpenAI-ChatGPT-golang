package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	gpt3 "github.com/csuzhang/gpt-go"
	"github.com/spf13/cobra"
)

func main() {
	log.SetOutput(new(NullWriter))
	apiKey := getAPIKeyFromEnv()
	ctx := context.Background()
	client := gpt3.NewClient(apiKey)
	rootCmd := &cobra.Command{
		Use:   "chatgpt",
		Short: "Chat with ChatGPT in console.",
		Run: func(cmd *cobra.Command, args []string) {
			scanner := bufio.NewScanner(os.Stdin)
			exit := false
			for !exit {
				fmt.Print("Please input question(q for exit): ")
				if !scanner.Scan() {
					break
				}
				questionParam := validateQuestion(scanner.Text())
				switch questionParam {
				case "q":
					exit = true
				case "":
					continue
				default:
					GetResponse(client, ctx, questionParam)
				}
			}
		},
	}
	log.Fatal(rootCmd.Execute())
}

func getAPIKeyFromEnv() string {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		panic("Missing API KEY")
	}
	return apiKey
}

func validateQuestion(question string) string {
	quest := strings.Trim(question, " ")
	keywords := []string{"", "loop", "break", "continue", "cls", "exit", "block"}
	for _, x := range keywords {
		if quest == x {
			return ""
		}
	}
	return quest
}

// GetResponse get response from gpt3
func GetResponse(client gpt3.Client, ctx context.Context, question string) {
	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			question,
		},
		MaxTokens:   gpt3.IntPtr(3000),
		Temperature: gpt3.Float32Ptr(0),
	}, func(resp *gpt3.CompletionResponse) {
		fmt.Print(resp.Choices[0].Text)
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(13)
	}
	fmt.Printf("\n")
}

// NullWriter is a writer on which all Write calls succeed
type NullWriter int

// Write implements io.Writer
func (NullWriter) Write([]byte) (int, error) { return 0, nil }
