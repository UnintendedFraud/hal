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
		fmt.Println("create ai client")
		aiclient := openai.NewClient(env.GetEnvVariables().OpenaiHalToken)

		fmt.Println("ai client created")

		res, err := aiclient.ChatCompletions(m.Content)
		if err != nil {
			log.Panicf("failed to query open ai with the following prompt [%s]. Error: %s", m.Content, err.Error())
		}

		fmt.Println("ai client replied")

		if len(res.Choices) > 0 {
			aiResponse := res.Choices[0].Message.Content
			if _, err = s.ChannelMessageSend(m.ChannelID, aiResponse); err != nil {
				log.Panicf("failed to send the response [%s] to the discord channel [%s]", aiResponse, err.Error())
			}
		}

		fmt.Println("### ai did not return any response ###", res.Usage, res.Choices)

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
