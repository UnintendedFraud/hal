package commands

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	tempest "github.com/Amatsagu/Tempest"
)

type MessageOfTheDay struct {
  Date time.Time
  Message tempest.Message
}

type PinnedData struct {
  Messages []tempest.Message
  Count int
  MessageOfTheDay *MessageOfTheDay
}

// not ideal but good enough for this
var data map[string]*PinnedData = make(map[string]*PinnedData)

func InitPinned(c tempest.Client, serverIDs []tempest.Snowflake) error {
  for _, sid := range serverIDs {
    channels, err := getChannels(c.Rest, sid.String())
    if err != nil {
      return err
    }

    for _, cid := range channels {
      messages, err := getPinnedMessages(c.Rest, cid.ID.String())
      if err != nil {
        fmt.Println("failed to get the pinned messages for channel", cid, err)
      }

      data[cid.ID.String()] = &PinnedData{
        Messages: messages,
        Count: len(messages),
        MessageOfTheDay: &MessageOfTheDay{},
      }
    }
  }

  return nil
}

var Pinned tempest.Command = tempest.Command{
  Name: "pinned",
  Description: "will show a random pinned message",
  Options: []tempest.Option{},
  SlashCommandHandler: func(itx tempest.CommandInteraction) {
    channelData, exist := data[itx.ChannelId.String()]
    if !exist {
      itx.SendLinearReply("no data for this channel", false)
    }

    if channelData.Count == 0 {
      itx.SendLinearReply("no pinned messages", false)
      return
    }

    today := time.Now()
    todayDate := time.Now().Format("2006/01/02")
    currentMessageDate := channelData.MessageOfTheDay.Date.Format("2006/01/02")

    if channelData.MessageOfTheDay == nil || todayDate != currentMessageDate {
      rand.Seed(time.Now().UnixNano())
      newMessage := channelData.Messages[rand.Intn(channelData.Count)] 
      
      channelData.MessageOfTheDay = &MessageOfTheDay{
        Date: today,
        Message: newMessage,
      }
    }

    fmt.Println("# message embeds #", len(channelData.MessageOfTheDay.Message.Embeds))

    embed := &tempest.Embed{
      Title: "Pin of the day",
      Author: &tempest.EmbedAuthor{
        Name: channelData.MessageOfTheDay.Message.Author.Username,
        IconUrl: channelData.MessageOfTheDay.Message.Author.FetchAvatarUrl(),
      },
      Timestamp: channelData.MessageOfTheDay.Message.Timestamp,
      Description: channelData.MessageOfTheDay.Message.Content,
    }

    if err := itx.SendReply(tempest.ResponseData{
      //Content: channelData.MessageOfTheDay.Message.Content,
      Embeds: []*tempest.Embed{embed},
      //Components: channelData.MessageOfTheDay.Message.Components,
    }, false); err != nil {
      itx.SendLinearReply(err.Error(), true)
    }
  },
}

func getChannels(rest tempest.Rest, serverID string) ([]Channel, error) {
  route := fmt.Sprintf("/guilds/%s/channels", serverID)

  bytes, err := rest.Request("GET", route, nil)
  if err != nil {
    return nil, err
  }

  channels := []Channel{}
  if err = json.Unmarshal(bytes, &channels); err != nil {
    return nil, err
  }

  return channels, nil
}

func getPinnedMessages(rest tempest.Rest, channelID string) ([]tempest.Message, error) {
  route := fmt.Sprintf("/channels/%s/pins", channelID)

  bytes, err := rest.Request("GET", route, nil)
  if err != nil {
    return nil, err
  }

  messages := []tempest.Message{}
  if err = json.Unmarshal(bytes, &messages); err != nil {
    return nil, err
  }

  return messages, nil
}

type Channel struct {
  ID tempest.Snowflake `json:"id"`
}
