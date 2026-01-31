package main

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

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

	"rename": {
		Command: discordgo.ApplicationCommand{
			Name: "rename",
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.SpanishES: "renombrar",
			},
			Description: "Change another user's name",
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.SpanishES: "Cambiar el nombre de otro usuario",
			},
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type: discordgo.ApplicationCommandOptionUser,
					Name: "user",
					NameLocalizations: map[discordgo.Locale]string{
						discordgo.SpanishES: "usuario",
					},
					Description: "User",
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.SpanishES: "Usuario",
					},
					Required: true,
				},
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "name",
					NameLocalizations: map[discordgo.Locale]string{
						discordgo.SpanishES: "nombre",
					},
					Description: "Name",
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.SpanishES: "Nombre",
					},
					Required: true,
				},
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			user := i.ApplicationCommandData().GetOption("user")
			name := i.ApplicationCommandData().GetOption("name")

			userId := user.UserValue(s).ID

			member, err := s.GuildMember(i.GuildID, userId)
			if err != nil {
				messageError(s, i)
				return
			}

			err = s.GuildMemberNickname(i.GuildID, userId, name.StringValue())
			if err != nil {
				messageError(s, i)
				return
			}

			var response string
			switch i.Locale {
			case discordgo.SpanishES:
				response = fmt.Sprintf("El nombre de `%s` ha sido cambiado a `%s`", member.Nick, name.StringValue())
			default:
				response = fmt.Sprintf("`%s`'s name has been changed to `%s`", member.Nick, name.StringValue())
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: response,
				},
			})
		},
	},

	"join": {
		Command: discordgo.ApplicationCommand{
			Name: "join",
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.SpanishES: "unirse",
			},
			Description: "Join your current voice channel",
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.SpanishES: "Unirse a tu canal de voz actual",
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			voice, err := findUserVoiceState(s, i.Member.User.ID)
			if err != nil {
				messageError(s, i)
				return
			}

			_, err = s.ChannelVoiceJoin(i.GuildID, voice.ChannelID, false, false)
			if err != nil {
				messageError(s, i)
				return
			}

			var response string
			switch i.Locale {
			case discordgo.SpanishES:
				response = "Me he unido a tu canal de voz"
			default:
				response = "I joined your voice channel"
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: response,
				},
			})

			if err != nil {
				messageError(s, i)
				return
			}
		},
	},

	"leave": {
		Command: discordgo.ApplicationCommand{
			Name: "leave",
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.SpanishES: "salirse",
			},
			Description: "Leave a voice channel",
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.SpanishES: "Salirse de un canal de voz",
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			voice, ok := s.VoiceConnections[i.GuildID]
			if !ok {
				messageError(s, i)
				return
			}

			err := voice.Disconnect()
			if err != nil {
				messageError(s, i)
				return
			}

			var response string
			switch i.Locale {
			case discordgo.SpanishES:
				response = "Me he salido del canal de voz"
			default:
				response = "I left the voice channel"
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: response,
				},
			})

			if err != nil {
				messageError(s, i)
				return
			}
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

func findUserVoiceState(s *discordgo.Session, userId string) (*discordgo.VoiceState, error) {
	for _, guild := range s.State.Guilds {
		for _, voice := range guild.VoiceStates {
			if voice.UserID == userId {
				return voice, nil
			}
		}
	}
	return nil, errors.New("Could not find user's voice state")
}
