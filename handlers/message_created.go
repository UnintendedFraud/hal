package handlers

import (
	"fmt"
	"hal/env"
	"hal/openai"
	"log"

	"github.com/bwmarrin/discordgo"
)

const MAX_HISTORY = 50

var messagesHistory = []*openai.ChatMessage{}

func OnMessageCreated(s *discordgo.Session, m *discordgo.MessageCreate) {
	isHal := m.Author.ID == s.State.User.ID

	fmt.Println("message received: ", m.Message.Content)

	addMessageToHistory(m.Message, isHal)

	fmt.Println("added message to history", len(messagesHistory))

	if isHal {
		return
	}

	fmt.Println("check contains user")

	if containUser(m.Mentions, s.State.User.ID) {
		fmt.Println("lets go")
		aiclient := openai.NewClient(env.GetEnvVariables().OpenaiHalToken)

		fmt.Println("client created")

		res, err := aiclient.Chat(messagesHistory)
		if err != nil {
			sendResponse(
				s,
				m.ChannelID,
				fmt.Sprintf("X_X: %s", err.Error()),
			)

			log.Panicf("failed to query open ai with the following prompt [%s]. Error: %s", m.Content, err.Error())
		}

		if len(res.Choices) > 0 {
			aiResponse := res.Choices[0].Message.Content
			sendResponse(s, m.ChannelID, aiResponse)
		}
	} else {
		fmt.Println("hal not mentioned")
	}
}

func addMessageToHistory(m *discordgo.Message, isHal bool) {
	var role string
	if isHal {
		role = "system"
	} else {
		role = "user"
	}

	messagesHistory = append(messagesHistory, &openai.ChatMessage{
		Role:    role,
		Content: m.Content,
	})

	if len(messagesHistory) > MAX_HISTORY {
		messagesHistory = messagesHistory[1:]
	}
}

func sendResponse(s *discordgo.Session, channelID string, response string) {
	if _, err := s.ChannelMessageSend(channelID, response); err != nil {
		log.Panicf("failed to send the response [%s] to the discord channel [%s]", response, err.Error())
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
