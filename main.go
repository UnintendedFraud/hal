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

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dgclient.Close()

	fmt.Println("after .Close")
	// addr := fmt.Sprintf(":%s", env.Port)
	// fmt.Println("starting server at", addr)
	//
	// if err := http.ListenAndServe(addr); err != nil {
	// 	panic(err)
	// }
}
