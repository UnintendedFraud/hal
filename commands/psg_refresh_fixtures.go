package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	tempest "github.com/Amatsagu/Tempest"
)

const endpoint = "https://v3.football.api-sports.io/fixtures"

var PsgRefreshFixtures tempest.Command = tempest.Command{
  Name: "psg-refresh-fixtures",
  Description: "met à jour les matchs du psg si ça change ou des matchs sont ajoutés (e.g. CL)",
  Options: []tempest.Option{},
  SlashCommandHandler: func(itx tempest.CommandInteraction) {
    if itx.User.Id != 160073741076922369 {
      itx.SendLinearReply("command cannot be used by you :(", true)
      return
    }

    client := &http.Client{}

    payload := &Payload{
      Season: "2022",
      Team: "85",
    }

    payloadBytes, err := json.Marshal(payload)
    if err != nil {
      itx.SendLinearReply(fmt.Sprintf("failed to marshal the payload: %s", err.Error()), true)
      return
    }

    req, err := http.NewRequest("GET", endpoint, bytes.NewReader(payloadBytes))
    if err != nil {
      itx.SendLinearReply(fmt.Sprintf("failed to create the request: %s", err.Error()), true)
      return
    }

    req.Header.Add("x-rapidapi-key", os.Getenv("FOOTBALL_API_KEY"))
    req.Header.Add("x-rapidapi-host", "v3.football.api-sports.io")

    res, err := client.Do(req)
    if err != nil {
      itx.SendLinearReply(fmt.Sprintf("failed to execute the request: %s", err.Error()), true)
      return
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
      itx.SendLinearReply(fmt.Sprintf("failed to read the response: %s", err.Error()), true)
      return
    }

    var r *Response
    if err = json.Unmarshal(body, r); err != nil {
      itx.SendLinearReply(fmt.Sprintf("failed to unmarshal the response: %s", err.Error()), true)
      return
    }

    itx.SendLinearReply(fmt.Sprintf("it worked, matches count: %d", r.Results), true)
  },
}

type Match struct {
  Fixture 
}

type Fixture struct {
  ID int64 `json:"id"`
  Referee string `json:"referee"`
  Timezone string `json:"timezone"`
  Date string `json:"date"`
  Timestamp string `json:"timestamp"`
  Venue VenueStruct `json:"venue"`
  League LeagueStruct `json:"league"`
  Teams struct {
    Home TeamStruct `json:"home"`
    Away TeamStruct `json:"away"`
  }`json:"teams"`
}

type TeamStruct struct {
  ID int64 `json:"id"`
  Name string `json:"name"`
}

type LeagueStruct struct {
  ID int64 `json:"id"`
  Name string `json:"name"`
  Country string `json:"country"`
  Round string `json:"round"`
}

type VenueStruct struct {
  ID int64 `json:"id"`
  Name string `json:"name"`
  City string `json:"city"`
}

type Response struct {
  Get string `json:"get"`
  Parameters map[string]string `json:"parameters"`
  Errors []interface{} `json:"errors"`
  Results int64 `json:"results"`
  Paging Paging `json:"paging"`
  Response []Match `json:"response"`
}

type Paging struct {
  Current int64 `json:"current"`
  Total int64 `json:"total"`
}

type Payload struct {
  Season string `json:"season"`
  Team string `json:"team"`
}
