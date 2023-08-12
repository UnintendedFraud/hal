package main

import (
	"fmt"
	"hal/commands"
	"hal/env"
	"hal/handlers"
	"log"
	"time"

	tempest "github.com/Amatsagu/Tempest"
	discordbot "github.com/bwmarrin/discordgo"
)

func main() {
	fmt.Println("HAL started")
	env := env.GetEnvVariables()

	dgclient, err := discordbot.New(env.Token)
	if err != nil {
		panic(err)
	}
	defer dgclient.Close()

	dgclient.AddHandler(handlers.OnMessageCreated)

	err = dgclient.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	client := tempest.CreateClient(tempest.ClientOptions{
		ApplicationId: env.AppID,
		PublicKey:     env.PublicKey,
		Token:         env.Token,
		PreCommandExecutionHandler: func(itx tempest.CommandInteraction) *tempest.ResponseData {
			log.Printf("running [%s] slash command", itx.Data.Name)
			return nil
		},
		Cooldowns: &tempest.ClientCooldownOptions{
			Duration:  5 * time.Second,
			Ephemeral: true,
			CooldownResponse: func(user tempest.User, timeLeft time.Duration) tempest.ResponseData {
				return tempest.ResponseData{
					Content: fmt.Sprintf("stop spamming, try again in %.2fs", timeLeft.Seconds()),
				}
			},
		},
	})

	if err = initialize(client, env.ServerIDs); err != nil {
		panic(err)
	}

	if err = client.RegisterCommand(commands.Pinned); err != nil {
		panic(err)
	}
	if err = client.RegisterCommand(commands.PsgRefreshFixtures); err != nil {
		panic(err)
	}
	if err = client.RegisterCommand(commands.PsgNextMatch); err != nil {
		panic(err)
	}

	if err = client.SyncCommands(env.ServerIDs, nil, false); err != nil {
		panic(err)
	}

	addr := fmt.Sprintf("%s:%s", env.Addr, env.Port)
	fmt.Println("starting server at", addr)

	if err := client.ListenAndServe(addr); err != nil {
		panic(err)
	}
}

func initialize(c tempest.Client, serverIDs []tempest.Snowflake) error {
	if err := commands.InitPinned(c, serverIDs); err != nil {
		return err
	}

	return nil
}
