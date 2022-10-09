package commands

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	tempest "github.com/Amatsagu/Tempest"
)

const timeout time.Duration = 10000000000 

var Pinned tempest.Command = tempest.Command{
  Name: "pinned",
  Description: "will show a random pinned message",
  Options: []tempest.Option{},
  SlashCommandHandler: func(itx tempest.CommandInteraction) {

    _, close := itx.Client.AwaitComponent([]string{itx.Data.CustomId}, timeout)

    messages, err := getPinnedMessages(itx.Client.Rest, itx.ChannelId.String())
    if err != nil {
      itx.SendLinearReply(err.Error(), false)
      return
    }


    messagesCount := len(messages)

    if messagesCount == 0 {
      itx.SendLinearReply("no pinned messages", false)
      return
    }

    idx := rand.Intn(messagesCount) 

    if _, err = itx.Client.SendMessage(itx.ChannelId, messages[idx]); err != nil {
      itx.SendLinearReply(err.Error(), false)
    }

    close()
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

