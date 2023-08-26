package handlers

import (
	"fmt"
	"hal/env"
	"hal/openai"
	"log"

	"github.com/bwmarrin/discordgo"
)

func OnMessageCreated(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Printf("\nmessage received 1: [%s]", m.Message.Content)

	if m.Author.ID == s.State.User.ID {
		return
	}

	if containUser(m.Mentions, s.State.User.ID) {
		fmt.Printf("\nmessage received 2: [%s]", m.Message.Content)

		aiclient := openai.NewClient(env.GetEnvVariables().OpenaiHalToken)

		fmt.Printf("\nAfter openai client: %+v", aiclient)

		res, err := aiclient.Completions(m.Content)
		if err != nil {
			log.Panicf("failed to query open ai with the following prompt [%s]. Error: %s", m.Content, err.Error())
		}

		fmt.Printf("\nopen ai responses: [%d]", len(res.Choices))

		if len(res.Choices) > 0 {
			aiResponse := res.Choices[0].Text
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
