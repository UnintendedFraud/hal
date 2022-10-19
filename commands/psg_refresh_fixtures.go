package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	tempest "github.com/Amatsagu/Tempest"
)

type Matches []Match

func (m Matches) Len() int {
  return len(m)
}

func (m Matches) Less(i int, j int) bool {
  return m[i].Fixture.Timestamp < m[j].Fixture.Timestamp
}

func (m Matches) Swap(i int, j int) {
  m[i], m[j] = m[j], m[i]
}

const PSG_FILE_PATH = "../data/football/psg_fixtures.json"

const endpoint = "https://v3.football.api-sports.io/fixtures"

const authorizedUser = "160073741076922369"

var PsgRefreshFixtures tempest.Command = tempest.Command{
  Name: "psg-refresh-fixtures",
  Description: "met à jour les matchs du psg si ça change ou des matchs sont ajoutés (e.g. CL)",
  Options: []tempest.Option{},
  SlashCommandHandler: func(itx tempest.CommandInteraction) {
    if itx.User.Id.String() != authorizedUser {
      itx.SendLinearReply("command cannot be used by you :(", true)
      return
    }

    client := &http.Client{}

    url := fmt.Sprintf("%s?season=2022&team=85&timezone=Europe/Paris", endpoint)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
      itx.SendLinearReply(fmt.Sprintf("failed to create the request: %s", err.Error()), true)
      return
    }

    req.Header.Add("x-rapidapi-key", "a22d9050a0d04cfeb988bffe89854617")
    req.Header.Add("x-rapidapi-host", "v3.football.api-sports.io")

    res, err := client.Do(req)
    if err != nil {
      itx.SendLinearReply(fmt.Sprintf("request failed: %s", err.Error()), true)
      return
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
      itx.SendLinearReply(fmt.Sprintf("failed to read the response body: %s", err.Error()), true)
      return
    }

    r := &Response{}
    if err = json.Unmarshal(body, r); err != nil {
      itx.SendLinearReply(fmt.Sprintf("failed to unmarshal the response: %s", err.Error()), true)
      return
    }

    matches := Matches(r.Response)
    sort.Sort(matches)

    sortedBytes, err := json.Marshal(matches)
    if err != nil {
      itx.SendLinearReply(fmt.Sprintf("failed to marshal the sorted list: %s", err.Error()), true)
      return
    }

    f, err := os.Create(PSG_FILE_PATH)
    if err != nil {
      itx.SendLinearReply(fmt.Sprintf("failed to create the json file: %s", err.Error()), true)
      return
    }

    if _, err = f.Write(sortedBytes); err != nil {
      itx.SendLinearReply(fmt.Sprintf("failed to write the json: %s", err.Error()), true)
      return
    }

    itx.SendLinearReply("matches refreshed", true)
  },
}

var PsgNextMatch tempest.Command = tempest.Command{
  Name: "psg-next-match",
  Description: "montre le prochain match du PSG",
  Options: []tempest.Option{},
  SlashCommandHandler: func(itx tempest.CommandInteraction) {
    file, err := ioutil.ReadFile(PSG_FILE_PATH)
    if err != nil {
      itx.SendLinearReply(fmt.Sprintf("failed to open the json file: %s", err.Error()), true)
      return
    }

    var list Matches
    if err = json.Unmarshal(file, &list); err != nil {
      itx.SendLinearReply(fmt.Sprintf("failed to unmarshal the list: %s", err.Error()), true)
      return
    }

    now := time.Now()

    var nextMatch *Match

    for _, m := range list {
      if m.Fixture.Date.Before(now) {
        continue
      }

      nextMatch = &m
    }

    if err := itx.SendReply(tempest.ResponseData{
      Content: formatNextPSGMatchContent(nextMatch),
    }, false); err != nil {
      log.Printf("failed to send reply with the pinned message: %s",err.Error())
      itx.SendLinearReply(err.Error(), true)
    }
  },
}

func formatNextPSGMatchContent(m *Match) string {
  if m == nil {
    return "no next match available"
  }

  return fmt.Sprintf(
    "%s - %s, %s (%s, %s)",
    m.Teams.Home.Name,
    m.Teams.Away.Name,
    m.Fixture.Date.Format("02 Jan 2006, 15:04"),
    m.Fixture.Venue.Name,
    m.Fixture.Venue.City,
  )
}


type Match struct {
  Fixture Fixture `json:"fixture"`
  League LeagueStruct `json:"league"`
  Teams struct {
    Home TeamStruct `json:"home"`
    Away TeamStruct `json:"away"`
  }`json:"teams"`
}

type Fixture struct {
  ID int64 `json:"id"`
  Referee string `json:"referee"`
  Timezone string `json:"timezone"`
  Date time.Time `json:"date"`
  Timestamp int64 `json:"timestamp"`
  Venue VenueStruct `json:"venue"`
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
