package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"

	tempest "github.com/Amatsagu/Tempest"
)

var Pinned tempest.Command = tempest.Command{
  Name: "pinned",
  Description: "will show a random pinned message",
  Options: []tempest.Option{},
  SlashCommandHandler: func(itx tempest.CommandInteraction) {
    log.Print("inside pinned command")

    messages, err := getPinnedMessages(itx.Client.Rest, itx.ChannelId.String())
    if err != nil {
      log.Printf("failed to get messages: %s", err.Error())
      itx.SendLinearReply(err.Error(), false)
      return
    }

    log.Print("after get pinned messages")

    messagesCount := len(messages)

    if messagesCount == 0 {
      itx.SendLinearReply("no pinned messages", false)
      return
    }

    idx := rand.Intn(messagesCount) 

    log.Printf("just before send message, %d, %d", idx, messagesCount)
    log.Printf("# selected message: %+v", messages[idx])
    

    if err = itx.Client.CrosspostMessage(itx.ChannelId, messages[idx].Id); err != nil {
      log.Printf("failed to send message: %s", err.Error())
      itx.SendLinearReply(err.Error(), false)
    }

  },
}

func getPinnedMessages(rest tempest.Rest, channelID string) ([]tempest.Message, error) {
  route := fmt.Sprintf("/channels/%s/pins", channelID)
  log.Printf("Route: %s", route)
  bytes, err := rest.Request("GET", route, nil)
  if err != nil {
    return nil, err
  }

  messages := []tempest.Message{}
  if err = json.Unmarshal(bytes, &messages); err != nil {
    log.Printf("## error unmarshalling ## %s", err.Error())
    return nil, err
  }

  return messages, nil
}

