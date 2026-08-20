package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if m.Author.ID == s.State.User.ID {
		return
	}

	message := m.Content
	user_mention := fmt.Sprintf("<@%s>", s.State.User.ID)

	if !strings.Contains(message, user_mention) {
		return
	}

	message = strings.ReplaceAll(message, user_mention, "Ralsei")

	s.ChannelTyping(m.ChannelID)

	body := &Ollama{
		Model:  *OllamaModel,
		System: *OllamaSystemPrompt,
		Prompt: message,
		Stream: false,
		Think:  false,
	}

	marshalled, err := json.Marshal(body)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/generate", *OllamaHost), bytes.NewReader(marshalled))
	if err != nil {
		return
	}

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()

	response_bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	var response_obj OllamaGenerateResponse
	err = json.Unmarshal(response_bytes, &response_obj)
	if err != nil {
		return
	}

	s.ChannelMessageSendReply(m.ChannelID, response_obj.Response, m.Reference())
}
