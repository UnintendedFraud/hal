package handlers

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"time"

	"hal/openai"

	"github.com/bwmarrin/discordgo"
)

const MAX_HISTORY = 50

const (
	SPAM_PERIOD         = 5 * time.Minute
	MESSAGES_PER_PERIOD = 10
)

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

var messagesHistory []*openai.ChatMessage = []*openai.ChatMessage{}

var usersHistoryCount map[string]*userHistoryCount = map[string]*userHistoryCount{}

type userHistoryCount struct {
	date          time.Time
	bannedUntilAt time.Time
	count         int
}

type Handler struct {
	client *openai.Client
}

func Init(token string) Handler {
	return Handler{
		client: openai.NewClient(token),
	}
}

// Tramp: 161970441441902592

func (h Handler) OnMessageCreated(s *discordgo.Session, m *discordgo.MessageCreate) {
	isHal := m.Author.ID == s.State.User.ID

	addMessageToHistory(m.Message, isHal)

	if isHal || !containHal(m.Mentions, s.State.User.ID) {
		return
	}

	userSpamTooMuch := updateUserHistoryCount(m.Author.ID)

	if userSpamTooMuch {
		sendResponse(s, m.ChannelID, getRandomSpam())
		return
	}

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
	} else {
		sendResponse(s, m.ChannelID, "la fatigue")
	}
}

func addMessageToHistory(m *discordgo.Message, isHal bool) []*openai.ChatMessage {
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

	return messagesHistory
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

func updateUserHistoryCount(userID string) bool {
	now := time.Now()

	u, ok := usersHistoryCount[userID]
	if !ok {
		usersHistoryCount[userID] = &userHistoryCount{
			date:  now,
			count: 1,
		}

		return false
	}

	u.count++

	if now.Before(u.bannedUntilAt) {
		return true
	}

	if !u.bannedUntilAt.IsZero() {
		u.bannedUntilAt = time.Time{}
		u.count = 1
	}

	if u.date.Add(SPAM_PERIOD).Before(now) {
		u.count = 1
		u.date = now
		return false
	}

	if u.count > MESSAGES_PER_PERIOD {
		u.bannedUntilAt = now.Add(30 * time.Minute)
		return true
	}

	return false
}

func getRandomSpam() string {
	return spams[rand.Intn(len(spams))]
}
