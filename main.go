package main

import (
	"fmt"
	"hal/commands"
	"log"
	"time"

	tempest "github.com/Amatsagu/Tempest"
)

const AppID tempest.Snowflake = 123
const PublicKey string = "public_key"
const Token string = "token"
const Addr string = "0.0.0.0:8080"

func main() {
  client := tempest.CreateClient(tempest.ClientOptions{
    ApplicationId: AppID,
    PublicKey: PublicKey,
    Token: Token,
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

  log.Printf("starting server at %s", Addr)

  client.RegisterCommand(commands.Hello)
  client.RegisterCommand(commands.Pinned)

  client.SyncCommands([]tempest.Snowflake{
    992760761812258868, // test server
  }, nil, false)

  if err := client.ListenAndServe(Addr); err != nil {
    panic(err)
  }
}
