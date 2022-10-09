package main

import (
	"fmt"
	"hal/commands"
	"log"
	"os"
	"time"

	tempest "github.com/Amatsagu/Tempest"
)

func main() {
  client := tempest.CreateClient(tempest.ClientOptions{
    ApplicationId: tempest.StringToSnowflake(os.Getenv("APP_ID")),
    PublicKey: os.Getenv("PUBLIC_KEY"),
    Token: os.Getenv("TOKEN"),
    PreCommandExecutionHandler: func(itx tempest.CommandInteraction) *tempest.ResponseData {
      log.Printf("running [%s] slash command", itx.Data.Name)
      return nil
    },
    Cooldowns: &tempest.ClientCooldownOptions{
      Duration: 5 * time.Second,
      Ephemeral: true,
      CooldownResponse: func(user tempest.User, timeLeft time.Duration) tempest.ResponseData {
        return tempest.ResponseData{
          Content: fmt.Sprintf("stop spamming, try again in %.2fs", timeLeft.Seconds()),
        }  
      },
    },
  }) 


  client.RegisterCommand(commands.Hello)
  client.RegisterCommand(commands.Pinned)

  client.SyncCommands([]tempest.Snowflake{
    992760761812258868, // test server
  }, nil, false)


  addr := os.Getenv("RAILWAY_STATIC_URL")
  if err := client.ListenAndServe(addr); err != nil {
    panic(err)
  }

  fmt.Println("starting server at", addr)
}
