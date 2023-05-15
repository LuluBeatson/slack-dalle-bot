package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image/png"
	"log"
	"os"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
	"github.com/shomali11/slacker"
)

var c *openai.Client

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Failed to load .env file:", err)
		return
	}

	c = openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	// createImageFile("A robot artist in a cute simplified style")
}

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	go printCommandEvents(bot.CommandEvents())

	bot.Command("hello", &slacker.CommandDefinition{
		Description: "Say hello",
		Examples:    []string{"hello"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("Hello!")
		},
	})

	bot.Command("<prompt>", &slacker.CommandDefinition{
		Description: "Create an image from a prompt and return the URL",
		Examples:    []string{"Portrait of a humanoid parrot in a classic costume, high detail, realistic light, unreal engine"},
		Handler:     createImageUrlHandler,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func createImageUrlHandler(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
	log.Println("Creating image...")
	prompt := request.Param("prompt")
	ctx := context.Background()

	response.Reply("Creating image...")
	reqUrl := openai.ImageRequest{
		Prompt:         prompt,
		Size:           openai.CreateImageSize256x256,
		ResponseFormat: openai.CreateImageResponseFormatURL,
		N:              1,
	}

	respUrl, err := c.CreateImage(ctx, reqUrl)
	if err != nil {
		response.Reply(fmt.Sprintf("Image creation error: %v\n", err))
		return
	}
	response.Reply(respUrl.Data[0].URL)
	log.Println("Image created:", respUrl.Data[0].URL)
}

func createImageFile(prompt string) {
	ctx := context.Background()
	reqBase64 := openai.ImageRequest{
		Prompt:         prompt,
		Size:           openai.CreateImageSize256x256,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
		N:              1,
	}

	respBase64, err := c.CreateImage(ctx, reqBase64)
	if err != nil {
		fmt.Printf("Image creation error: %v\n", err)
		return
	}

	imgBytes, err := base64.StdEncoding.DecodeString(respBase64.Data[0].B64JSON)
	if err != nil {
		fmt.Printf("Base64 decode error: %v\n", err)
		return
	}

	r := bytes.NewReader(imgBytes)
	imgData, err := png.Decode(r)
	if err != nil {
		fmt.Printf("PNG decode error: %v\n", err)
		return
	}

	file, err := os.Create("example.png")
	if err != nil {
		fmt.Printf("File creation error: %v\n", err)
		return
	}
	defer file.Close()

	if err := png.Encode(file, imgData); err != nil {
		fmt.Printf("PNG encode error: %v\n", err)
		return
	}

	fmt.Println("The image was saved as example.png")
}

func printCommandEvents(commandEvents <-chan *slacker.CommandEvent) {
	for event := range commandEvents {
		fmt.Println("Command Event. Timestamp:", event.Timestamp, "Command:", event.Command, "Parameters:", event.Parameters)
	}
}
