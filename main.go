package main

import (
	"fmt"
	"hal/commands"
	"hal/env"
	"hal/handlers"

	tempest "github.com/Amatsagu/Tempest"
	discordbot "github.com/bwmarrin/discordgo"
)

func main() {
	fmt.Println("HAL started")
	env := env.GetEnvVariables()

	var err error

	fmt.Println("env: ", env)

	discordToken := fmt.Sprintf("Bot %s", env.Token)
	dgclient, err := discordbot.New(discordToken)
	if err != nil {
		panic(fmt.Errorf("failed to create discord client: %s", err.Error()))
	}
	defer dgclient.Close()

	dgclient.AddHandler(handlers.OnMessageCreated)

	err = dgclient.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
}

func initialize(c tempest.Client, serverIDs []tempest.Snowflake) error {
	if err := commands.InitPinned(c, serverIDs); err != nil {
		return err
	}

	return nil
}
