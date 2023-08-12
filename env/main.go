package env

import (
	"os"
	"strings"

	tempest "github.com/Amatsagu/Tempest"
)

func GetEnvVariables() Env {
	ids := strings.Split(os.Getenv("DISCORD_SERVER_IDS"), ",")

	serverIDs := []tempest.Snowflake{}
	for _, id := range ids {
		serverIDs = append(serverIDs, tempest.StringToSnowflake(id))
	}

	return Env{
		AppID:          tempest.StringToSnowflake(os.Getenv("HAL_APP_ID")),
		PublicKey:      os.Getenv("HAL_PUBLIC_KEY"),
		Token:          os.Getenv("HAL_TOKEN"),
		Port:           os.Getenv("PORT"),
		Addr:           os.Getenv("ADDR"),
		ServerIDs:      serverIDs,
		OpenaiHalToken: os.Getenv("OPENAI_HAL"),
	}
}

type Env struct {
	AppID          tempest.Snowflake
	PublicKey      string
	Token          string
	Port           string
	Addr           string
	ServerIDs      []tempest.Snowflake
	OpenaiHalToken string
}
