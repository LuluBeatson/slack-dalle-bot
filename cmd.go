package main

import (
	"fmt"

	"github.com/shomali11/slacker"
)

var helloCmd = &slacker.CommandDefinition{
	Description: "Say hello",
	Examples:    []string{"hello"},
	Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
		response.Reply("Hello!")
	},
}

var uploadCmd = &slacker.CommandDefinition{
	Description: "Create an image from a prompt and upload it to the channel",
	Examples:    []string{"Portrait of a humanoid parrot in a classic costume, high detail, realistic light, unreal engine"},
	Handler:     createImageUploadHandler,
}

var urlCmd = &slacker.CommandDefinition{
	Description: "Create an image from a prompt and return the URL",
	Examples:    []string{"Portrait of a humanoid parrot in a classic costume, high detail, realistic light, unreal engine"},
	Handler:     createImageUrlHandler,
}

var creditsCmd = &slacker.CommandDefinition{
	Description: "Show remaining credits",
	Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
		response.Reply(fmt.Sprintf("Remaining image credits: %d", limit))
	},
}

var helpCmd = &slacker.CommandDefinition{
	Description: "Show help",
	Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
		response.Reply("Available commands:\n" +
			"• dalle-hello: Say hello to dalle-bot\n" +
			"• dalle-url <prompt>: Create an image from a prompt and return the URL\n" +
			"• dalle-upload <prompt>: Create an image from a prompt and upload it to the channel\n" +
			"• dalle-credits: Show remaining credits\n" +
			"• dalle-help: Show help")
	},
}
