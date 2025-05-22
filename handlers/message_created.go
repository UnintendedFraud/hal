package handlers

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"time"

	"hal/gemini"
	"hal/openai"

	"github.com/bwmarrin/discordgo"
	"google.golang.org/genai"
)

const MAX_HISTORY = 50

const DISCORD_MAX_CHAR = 4000

const (
	SPAM_PERIOD         = 5 * time.Minute
	MESSAGES_PER_PERIOD = 10
)

var spams = []string{
	"Arrête de spam putain!",
	"Wesh...",
	"Flemme.",
	"Laisse-moi tranquille!",
	"Fdr",
}

// messagesHistory []*openai.ChatMessage = []*openai.ChatMessage{}
var geminiHistory []*genai.Content = []*genai.Content{}

var usersHistoryCount map[string]*userHistoryCount = map[string]*userHistoryCount{}

type userHistoryCount struct {
	date          time.Time
	bannedUntilAt time.Time
	count         int
}

type Handler struct {
	openaiClient *openai.Client
	geminiClient *gemini.Client
}

func Init(openaiToken string, geminiToken string) Handler {
	return Handler{
		openaiClient: openai.NewClient(openaiToken),
		geminiClient: gemini.NewClient(geminiToken),
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

	llmRes, err := h.geminiClient.GenerateContent(geminiHistory)
	if err != nil {
		sendResponse(
			s,
			m.ChannelID,
			fmt.Sprintf("X_X: %s", err.Error()),
		)

		log.Printf("\nfailed to query the llm with the following prompt [%s]. Error: %s", m.Content, err.Error())
		return
	}

	sendResponse(s, m.ChannelID, llmRes.Text())
}

func addMessageToHistory(m *discordgo.Message, isHal bool) []*genai.Content {
	message := cleanMessage(m.Content)
	if message == "" {
		return geminiHistory
	}

	var role string
	if isHal {
		role = "model"
	} else {
		role = "user"
	}

	geminiHistory = append(geminiHistory, &genai.Content{
		Role: role,
		Parts: []*genai.Part{
			genai.NewPartFromText(message),
		},
	})

	if len(geminiHistory) > MAX_HISTORY {
		geminiHistory = geminiHistory[1:]
	}

	return geminiHistory
}

func cleanMessage(p string) string {
	regex := regexp.MustCompile(`<@\d+>`)
	return regex.ReplaceAllString(p, "")
}

func sendResponse(s *discordgo.Session, channelID string, response string) {
	if _, err := s.ChannelMessageSend(channelID, truncateIfNeeded(response)); err != nil {
		log.Printf("\nfailed to send the response [%s] to the discord channel [%s]", response, err.Error())
	}
}

func truncateIfNeeded(response string) string {
	if len(response) < DISCORD_MAX_CHAR {
		return response
	}

	return fmt.Sprintf("%s %s", response[:DISCORD_MAX_CHAR-10], "[...]")
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
