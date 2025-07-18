package commands

import (
	"petpet/lib"
)

var Meow lib.Command = lib.Command{
	Name:        "meow",
	Description: "Meow, meow!",
	Type:        lib.COMMAND_TYPE_CHAT_INPUT,
	CommandHandler: func(interaction *lib.CommandInteraction) {
		interaction.SendSimpleReply("Mrrowww~", false)
	},
}
