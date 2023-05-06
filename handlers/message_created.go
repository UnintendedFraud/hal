package handlers

import (
	"fmt"
	"hal/env"
	"hal/openai"
	"log"

	"github.com/bwmarrin/discordgo"
)

func OnMessageCreated(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if containUser(m.Mentions, s.State.User.ID) {
		aiclient := openai.NewClient(env.GetEnvVariables().OpenaiHalToken)

		res, err := aiclient.ChatCompletions(m.Content)
		if err != nil {
			log.Panicf("failed to query open ai with the following prompt [%s]. Error: %s", m.Content, err.Error())
		}

		if len(res.Choices) > 0 {
			aiResponse := res.Choices[0].Message.Content
			if _, err = s.ChannelMessageSend(m.ChannelID, aiResponse); err != nil {
				log.Panicf("failed to send the response [%s] to the discord channel [%s]", aiResponse, err.Error())
			}
		}

		return
	}
}

func containUser(users []*discordgo.User, userID string) bool {
	for _, u := range users {
		if u.ID == userID {
			return true
		}
	}

	return false
}
