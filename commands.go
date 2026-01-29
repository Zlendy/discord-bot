package main

import "github.com/bwmarrin/discordgo"

type CommandHandler struct {
	Command discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var commands = map[string]*CommandHandler{
	"hello-world": {
		Command: discordgo.ApplicationCommand{
			Name:        "hello-world",
			Description: "Hello, world!",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hello, world!",
				},
			})
		},
	},

	"goodbye": {
		Command: discordgo.ApplicationCommand{
			Name:        "goodbye",
			Description: "Sayonara",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Sayonara",
				},
			})
		},
	},
}
