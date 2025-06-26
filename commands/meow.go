package commands

import (
	tempest "github.com/amatsagu/tempest"
)

var Meow tempest.Command = tempest.Command{
	Name:        "Meow meow!",
	Description: "",
	Type: 2,
	SlashCommandHandler: func(itx *tempest.CommandInteraction) {
		itx.SendLinearReply("Meow~ 🐱", true)
	},
}