package handlers

import (
	"fmt"
	"hal/openai"
	"log"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

const MAX_HISTORY = 50

var messagesHistory = []*openai.ChatMessage{}

type Handler struct {
	client *openai.Client
}

func Init(token string) Handler {
	return Handler{
		client: openai.NewClient(token),
	}
}

func (h Handler) OnMessageCreated(s *discordgo.Session, m *discordgo.MessageCreate) {
	isHal := m.Author.ID == s.State.User.ID

	addMessageToHistory(m.Message, isHal)

	if isHal {
		return
	}

	if containUser(m.Mentions, s.State.User.ID) {
		res, err := h.client.Chat(messagesHistory)
		if err != nil {
			sendResponse(
				s,
				m.ChannelID,
				fmt.Sprintf("X_X: %s", err.Error()),
			)

			log.Printf("\nfailed to query open ai with the following prompt [%s]. Error: %s", m.Content, err.Error())
			return
		}

		if len(res.Choices) > 0 {
			aiResponse := res.Choices[0].Message.Content
			sendResponse(s, m.ChannelID, aiResponse)
		}
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
		Content: cleanMessage(m.Content),
	})

	if len(messagesHistory) > MAX_HISTORY {
		messagesHistory = messagesHistory[1:]
	}
}

func cleanMessage(p string) string {
	regex := regexp.MustCompile(`<@\d+>`)
	return regex.ReplaceAllString(p, "")
}

func sendResponse(s *discordgo.Session, channelID string, response string) {
	if _, err := s.ChannelMessageSend(channelID, response); err != nil {
		log.Printf("\nfailed to send the response [%s] to the discord channel [%s]", response, err.Error())
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
