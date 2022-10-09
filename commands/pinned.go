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
    messages, err := getPinnedMessages(itx.Client.Rest, itx.ChannelId.String())
    if err != nil {
      log.Printf("failed to get messages: %s", err.Error())
      itx.SendLinearReply(err.Error(), false)
      return
    }

    messagesCount := len(messages)

    if messagesCount == 0 {
      itx.SendLinearReply("no pinned messages", false)
      return
    }

    idx := rand.Intn(messagesCount) 

    itx.Client.SendMessage(itx.ChannelId, messages[idx])
  },
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

