package main

import "github.com/bwmarrin/discordgo"

type CommandHandler struct {
	Command discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var commands = map[string]*CommandHandler{
	"say": {
		Command: discordgo.ApplicationCommand{
			Name: "say",
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.SpanishES: "decir",
			},
			Description: "Make the bot say something",
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.SpanishES: "Hacer que el bot diga algo",
			},
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "message",
					NameLocalizations: map[discordgo.Locale]string{
						discordgo.SpanishES: "mensaje",
					},
					Description: "Message",
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.SpanishES: "Mensaje",
					},
					Required: true,
				},
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			message := i.ApplicationCommandData().GetOption("message")

			_, err := s.ChannelMessageSend(i.ChannelID, message.StringValue())
			if err != nil {
				messageError(s, i)
				return
			}

			var response string
			switch i.Locale {
			case discordgo.SpanishES:
				response = "Mensaje enviado"
			default:
				response = "Message sent"
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: response,
				},
			})
		},
	},
}

func messageError(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var response string
	switch i.Locale {
	case discordgo.SpanishES:
		response = "Ha habido un error"
	default:
		response = "There was an error"
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: response,
		},
	})
}
