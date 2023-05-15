package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
)

func createImageUrlHandler(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
	if limit == 0 {
		response.Reply("Sorry, I'm out of credits, please contact my admin to refill the credits.")
		return
	}
	prompt := request.Param("prompt")
	response.Reply("Creating image...")
	url, err := createImageUrl(prompt)
	limit = limit - 1
	if err != nil {
		log.Printf("Image creation error: %v\n", err)
		response.Reply("Sorry, something went wrong.")
		return
	}
	log.Println("Image created:", url)
	response.Reply(fmt.Sprintf("Here's your link: %s", url))
}

func createImageUploadHandler(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
	if limit == 0 {
		response.Reply("Sorry, I'm out of credits, please contact my admin to refill the credits.")
		return
	}

	channel := botCtx.Event().ChannelID
	if channel == "" {
		log.Println("Could not get channel ID")
		response.Reply("Sorry, something went wrong.")
		return
	}

	prompt := request.Param("prompt")
	response.Reply("Creating image...")
	imgBytes, err := createImageBytes(prompt)
	limit = limit - 1
	if err != nil {
		log.Printf("Image creation error: %v\n", err)
		response.Reply("Sorry, something went wrong.")
		return
	}

	response.Reply("Uploading image...")
	filename := fmt.Sprintf("slack-image-%d.png", time.Now().UnixNano())
	err = saveImage(imgBytes, filename)
	if err != nil {
		log.Printf("Image save error: %v\n", err)
		response.Reply("Sorry, something went wrong.")
		return
	}

	fileUploadParams := slack.FileUploadParameters{
		Filetype: "image/png",
		Filename: filename,
		File:     filename,
		Channels: []string{channel},
	}
	_, err = botCtx.APIClient().UploadFile(fileUploadParams)
	if err != nil {
		log.Printf("Image upload error: %v\n", err)
		response.Reply("Sorry, something went wrong.")
		return
	}

	if filename != "" {
		if err := os.Remove(filename); err != nil {
			log.Printf("Image file deletion error: %v\n", err)
		}
	}
}
