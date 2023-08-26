package main

import (
	"fmt"
	"hal/env"
	"hal/handlers"

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

	dgclient.AddHandler(handlers.OnMessageCreated)

	err = dgclient.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("before .Close")
	dgclient.Close()

	fmt.Println("after .Close")
	// addr := fmt.Sprintf(":%s", env.Port)
	// fmt.Println("starting server at", addr)
	//
	// if err := http.ListenAndServe(addr); err != nil {
	// 	panic(err)
	// }
}
