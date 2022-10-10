package main

import (
	"fmt"
	"hal/commands"
	"log"
	"os"
	"time"

	tempest "github.com/Amatsagu/Tempest"
)

var channelIDs []tempest.Snowflake = []tempest.Snowflake{
  992760761812258868, // test server
}

func main() {
  env := getEnvVariables()

  client := tempest.CreateClient(tempest.ClientOptions{
    ApplicationId: env.AppID,
    PublicKey: env.PublicKey,
    Token: env.Token,
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

  if err := initialize(client); err != nil {
    panic(err)
  }

  client.RegisterCommand(commands.Pinned)

  client.SyncCommands(channelIDs, nil, false)


  addr := fmt.Sprintf("%s:%s", env.Addr, env.Port)
  fmt.Println("starting server at", addr)

  if err := client.ListenAndServe(addr); err != nil {
    panic(err)
  }

}

func initialize(c tempest.Client) error {
  if err := commands.InitPinned(c, channelIDs); err != nil {
    return err
  }

  return nil
}

func getEnvVariables() Env {
  if os.Getenv("RAILWAY_ENVIRONMENT") == "production" {
    return Env{
      AppID: tempest.StringToSnowflake(os.Getenv("APP_ID")),
      PublicKey: os.Getenv("PUBLIC_KEY"),
      Token: os.Getenv("TOKEN"),
      Port: "8080",
      Addr: os.Getenv("ADDR"),
    }
  }

  return Env{}
}


type Env struct {
  AppID tempest.Snowflake
  PublicKey string
  Token string
  Port string
  Addr string
}
