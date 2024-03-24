package main

import (
	"fmt"
	"hal/env"
	"hal/handlers"
	"os"
	"os/signal"
	"syscall"

	discordbot "github.com/bwmarrin/discordgo"
)

func main() {
	fmt.Println("HAL started")
	env := env.GetEnvVariables()

	var err error

	handler := handlers.Init(env.OpenaiHalToken)

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

	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dgclient.Close()
}
