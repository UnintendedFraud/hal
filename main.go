package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"hal/env"
	"hal/handlers"
	"hal/llm"

	discordbot "github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	fmt.Println("HAL started")
	env := env.GetEnvVariables()

	handler := handlers.Init(&env)

	discordToken := fmt.Sprintf("Bot %s", env.Token)
	dgclient, err := discordbot.New(discordToken)
	if err != nil {
		panic(fmt.Errorf("failed to create discord client: %s", err.Error()))
	}

	dgclient.AddHandler(handler.OnMessageCreated)

	err = dgclient.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// ---
	llmClient := llm.NewClient(&env)
	res, err := llmClient.GenerateContent("Yo poto")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	// ---

	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dgclient.Close()
}
