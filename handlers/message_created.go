package handlers

import (
	"fmt"
	"hal/env"
	"hal/openai"
	"log"

	"github.com/bwmarrin/discordgo"
)

func OnMessageCreated(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Printf("new message posted: [%s], par [%s], authorID [%s], state user id [%s]", m.Content, m.Author.Username, m.Author.ID, s.State.User.ID)

	log.Printf("author [%+v], state user [%+v]", m.Author, s.State.User)
	if m.Author.ID == s.State.User.ID {
		return
	}

	if containUser(m.Mentions, s.State.User.ID) {
		fmt.Println("going through openai stuff")
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
		fmt.Println("## mentions ##", u.ID, u.Username, "-----", userID)
		if u.ID == userID {
			return true
		}
	}

	return false
}
