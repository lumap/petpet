package commands

import (
	"petpet/lib"
)

var Meow lib.Command = lib.Command{
	Name:        "meow",
	Description: "Meow, meow!",
	Type:        lib.COMMAND_TYPE_CHAT_INPUT,
	IntegrationTypes: []int{lib.COMMAND_INTEGRATION_TYPE_GUILD, lib.COMMAND_INTEGRATION_TYPE_USER},
	Contexts:         []int{lib.COMMAND_CONTEXT_GUILD, lib.COMMAND_CONTEXT_BOT_DM, lib.COMMAND_CONTEXT_PRIVATE_CHANNEL},
	CommandHandler: func(interaction *lib.CommandInteraction) {
		interaction.SendSimpleReply("Mrrowww~", false)
	},
}
