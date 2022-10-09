package commands

import tempest "github.com/Amatsagu/Tempest"

var Pinned tempest.Command = tempest.Command{
  Name: "pinned",
  Description: "will show a random pinned message",
  Options: []tempest.Option{},
  SlashCommandHandler: func(itx tempest.CommandInteraction) {

  },
}
