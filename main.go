package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
	"github.com/shomali11/slacker"
)

var c *openai.Client
var limit = 10

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Failed to load .env file:", err)
		return
	}

	c = openai.NewClient(os.Getenv("OPENAI_API_KEY"))
}

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	go printCommandEvents(bot.CommandEvents())

	bot.Command("dalle-hello", helloCmd)
	bot.Command("dalle-url <prompt>", urlCmd)
	bot.Command("dalle-upload <prompt>", uploadCmd)
	bot.Command("dalle-credits", creditsCmd)
	bot.Command("dalle-help", helpCmd)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func printCommandEvents(commandEvents <-chan *slacker.CommandEvent) {
	for event := range commandEvents {
		fmt.Println("Command Event. Timestamp:", event.Timestamp, "Command:", event.Command, "Parameters:", event.Parameters)
	}
}
