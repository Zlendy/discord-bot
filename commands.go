package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type CommandHandler struct {
	Command discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var commands = map[string]*CommandHandler{
	"say": {
		Command: discordgo.ApplicationCommand{
			Name:        "say",
			Description: "Hacer que el bot diga algo",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "Mensaje",
					Required:    true,
				},
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			message := i.ApplicationCommandData().GetOption("message")

			_, err := s.ChannelMessageSend(i.ChannelID, message.StringValue())
			if err != nil {
				messageError(s, i, err)
				return
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: "Mensaje enviado",
				},
			})
		},
	},

	"rename": {
		Command: discordgo.ApplicationCommand{
			Name:        "rename",
			Description: "Cambiar el nombre de otro usuario",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "Usuario",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "Nombre",
					Required:    true,
				},
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			user := i.ApplicationCommandData().GetOption("user")
			name := i.ApplicationCommandData().GetOption("name")

			userId := user.UserValue(s).ID

			member, err := s.GuildMember(i.GuildID, userId)
			if err != nil {
				messageError(s, i, err)
				return
			}

			err = s.GuildMemberNickname(i.GuildID, userId, name.StringValue())
			if err != nil {
				messageError(s, i, err)
				return
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("El nombre de `%s` ha sido cambiado a `%s`", member.Nick, name.StringValue()),
				},
			})
		},
	},

	"join": {
		Command: discordgo.ApplicationCommand{
			Name:        "join",
			Description: "Unirse a tu canal de voz actual",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			voice, err := findUserVoiceState(s, i.Member.User.ID)
			if err != nil {
				messageError(s, i, err)
				return
			}

			_, err = s.ChannelVoiceJoin(i.GuildID, voice.ChannelID, false, false)
			if err != nil {
				messageError(s, i, err)
				return
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: "Me he unido a tu canal de voz",
				},
			})

			if err != nil {
				messageError(s, i, err)
				return
			}
		},
	},

	"leave": {
		Command: discordgo.ApplicationCommand{
			Name:        "leave",
			Description: "Salirse de un canal de voz",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			voice, ok := s.VoiceConnections[i.GuildID]
			if !ok {
				messageError(s, i, errors.New("Key not in map"))
				return
			}

			err := voice.Disconnect()
			if err != nil {
				messageError(s, i, err)
				return
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: "Me he salido del canal de voz",
				},
			})

			if err != nil {
				messageError(s, i, err)
				return
			}
		},
	},

	"activity": {
		Command: discordgo.ApplicationCommand{
			Name:        "activity",
			Description: "Comprobar que usuarios contienen un texto en su actividad",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "text",
					Description: "Texto",
					Required:    true,
				},
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			text := i.ApplicationCommandData().GetOption("text").StringValue()

			err := s.RequestGuildMembers(i.GuildID, "", 0, "", true)
			if err != nil {
				messageError(s, i, err)
				return
			}

			guild, err := s.State.Guild(i.GuildID)
			if err != nil {
				messageError(s, i, err)
				return
			}

			text_regexp := regexp.MustCompile(fmt.Sprintf("(?i)%s", regexp.QuoteMeta(text)))
			var message_users strings.Builder

			for _, presence := range guild.Presences {
				activities := presence.Activities

				var activity_text string
				if len(activities) > 0 {
					activity := activities[len(activities)-1]

					activity_text = activity.State
					if activity_text == "" {
						activity_text = activity.Name
					}
				}

				if text_regexp.MatchString(activity_text) {
					// Highlight substring contained in activity text. Example: "Golang goes brr" -> "**Go**lang **go**es brr"
					activity_text = text_regexp.ReplaceAllStringFunc(activity_text, func(s string) string {
						return fmt.Sprintf("**%s**", s)
					})

					// Append a formatted message that mentions the user and shows their activity
					fmt.Fprintf(&message_users, "- %s: %s\n", presence.User.Mention(), activity_text)
				}
			}

			var message strings.Builder
			if message_users.Len() == 0 {
				fmt.Fprintf(&message, "Ningún usuario contiene el texto `%s` en su actividad\n", text)
			} else {
				fmt.Fprintf(&message, "Estos usuarios contienen el texto `%s` en su actividad:\n", text)
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: message.String(),
				},
			})

			if err != nil {
				messageError(s, i, err)
				return
			}
		},
	},

	"russianroulette": {
		Command: discordgo.ApplicationCommand{
			Name:        "russianroulette",
			Description: "Si pierdes, te llevas un aislamiento temporal",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var response string
			var err error

			if rand.IntN(6) == 0 { // 1/6 possibilities to lose
				until := time.Now().Local().Add(time.Minute * time.Duration(5))
				err = s.GuildMemberTimeout(i.GuildID, i.Member.User.ID, &until)
				if err != nil {
					messageError(s, i, err)
					return
				}

				response = "Has perdido"
			} else {
				response = "No has perdido"
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: response,
				},
			})

			if err != nil {
				messageError(s, i, err)
				return
			}
		},
	},

	"chat": {
		Command: discordgo.ApplicationCommand{
			Name:        "chat",
			Description: "Habla con Ralsei",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "Mensaje",
					Required:    true,
				},
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			message := i.ApplicationCommandData().GetOption("message").StringValue()

			s.ChannelTyping(i.ChannelID)

			body := &Ollama{
				Model:  *OllamaModel,
				System: *OllamaSystemPrompt,
				Prompt: message,
				Stream: false,
				Think:  false,
			}

			marshalled, err := json.Marshal(body)
			if err != nil {
				messageError(s, i, err)
				return
			}

			req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/generate", *OllamaHost), bytes.NewReader(marshalled))
			if err != nil {
				messageError(s, i, err)
				return
			}

			client := &http.Client{Timeout: 30 * time.Second}

			res, err := client.Do(req)
			if err != nil {
				messageError(s, i, err)
				return
			}

			defer res.Body.Close()

			response_bytes, err := io.ReadAll(res.Body)
			if err != nil {
				messageError(s, i, err)
				return
			}

			var response_obj OllamaGenerateResponse
			err = json.Unmarshal(response_bytes, &response_obj)
			if err != nil {
				messageError(s, i, err)
				return
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: response_obj.Response,
				},
			})

			if err != nil {
				messageError(s, i, err)
				return
			}
		},
	},
}

func messageError(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprintf("Ha habido un error: %s", err),
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
	return nil, errors.New("No se ha podido encontrar el estado de voz del usuario")
}
