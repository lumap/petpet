package commands

import (
	"petpet/lib"
)

var Meow lib.Command = lib.Command{
	Name:        "meow",
	Description: "Meow, meow!",
	Type: lib.COMMAND_TYPE_CHAT_INPUT,
	CommandHandler: func(itx *lib.CommandInteraction) {
		itx.SendSimpleReply("Mrrowww~", false)
	},
}