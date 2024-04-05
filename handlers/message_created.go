package handlers

import (
	"fmt"
	"hal/openai"
	"log"
	"math/rand"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
)

const MAX_HISTORY = 50

const SPAM_PERIOD = 5 * time.Minute

var spams = []string{
	":GroGroDebile:",
	"Arrête de spam putain!",
	"Wesh...",
	"...",
	"Flemme.",
	"Laisse-moi tranquille!",
	"Fdr",
	":soxx:",
}

type userHistoryCount struct {
	date          time.Time
	bannedUntilAt time.Time
	count         int
}

type Handler struct {
	messagesHistory   []*openai.ChatMessage
	usersHistoryCount map[string]userHistoryCount

	client *openai.Client
}

func Init(token string) Handler {
	return Handler{
		client: openai.NewClient(token),

		messagesHistory:   []*openai.ChatMessage{},
		usersHistoryCount: map[string]userHistoryCount{},
	}
}

// Tramp: 161970441441902592

func (h Handler) OnMessageCreated(s *discordgo.Session, m *discordgo.MessageCreate) {
	isHal := m.Author.ID == s.State.User.ID

	h.addMessageToHistory(m.Message, isHal)

	if isHal || !containHal(m.Mentions, s.State.User.ID) {
		return
	}

	userSpamTooMuch := h.updateUserHistoryCount(m.Author.ID)

	if userSpamTooMuch {
		sendResponse(s, m.ChannelID, getRandomSpam())
	}

	res, err := h.client.Chat(h.messagesHistory)
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
	} else {
		sendResponse(s, m.ChannelID, "la fatigue")
	}
}

func (h Handler) addMessageToHistory(m *discordgo.Message, isHal bool) {
	var role string
	if isHal {
		role = "system"
	} else {
		role = "user"
	}

	h.messagesHistory = append(h.messagesHistory, &openai.ChatMessage{
		Role:    role,
		Content: cleanMessage(m.Content),
	})

	if len(h.messagesHistory) > MAX_HISTORY {
		h.messagesHistory = h.messagesHistory[1:]
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

func containHal(users []*discordgo.User, userID string) bool {
	for _, u := range users {
		if u.ID == userID {
			return true
		}
	}

	return false
}

func (h Handler) updateUserHistoryCount(userID string) bool {
	now := time.Now()

	u, ok := h.usersHistoryCount[userID]
	if !ok {
		h.usersHistoryCount[userID] = userHistoryCount{
			date:  now,
			count: 1,
		}

		return false
	}

	u.count++

	if now.Before(u.bannedUntilAt) {
		return true
	}

	if u.date.Add(SPAM_PERIOD).After(now) {
		u.count = 1
		u.date = now

		return false
	}

	if u.count > 5 {
		return true
	}

	return false
}

func getRandomSpam() string {
	return spams[rand.Intn(len(spams))]
}
