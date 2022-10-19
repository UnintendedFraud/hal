package main

import (
	"fmt"
	"hal/commands"
	"log"
	"os"
	"strings"
	"time"

	tempest "github.com/Amatsagu/Tempest"
)

func main_old() {
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

  if err := initialize(client, env.ServerIDs); err != nil {
    panic(err)
  }

  client.RegisterCommand(commands.Pinned)
  client.RegisterCommand(commands.PsgRefreshFixtures)
  client.RegisterCommand(commands.PsgNextMatch)

  client.SyncCommands(env.ServerIDs, nil, false)


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

func getEnvVariables() Env {
  if os.Getenv("RAILWAY_ENVIRONMENT") == "production" {
    ids := strings.Split(os.Getenv("SERVER_IDS"), ",")

    serverIDs := []tempest.Snowflake{}
    for _, id := range ids {
      serverIDs = append(serverIDs, tempest.StringToSnowflake(id)) 
    }

    return Env{
      AppID: tempest.StringToSnowflake(os.Getenv("APP_ID")),
      PublicKey: os.Getenv("PUBLIC_KEY"),
      Token: os.Getenv("TOKEN"),
      Port: "8080",
      Addr: os.Getenv("ADDR"),
      ServerIDs: serverIDs,
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
  ServerIDs []tempest.Snowflake
}
