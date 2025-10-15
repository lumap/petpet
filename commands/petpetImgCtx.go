package commands

import (
	"petpet/lib"
)

var PetpetImgCtx = lib.Command{
	Type:        4,
	Name:        "Petpet this image",
	Description: "",
	IntegrationTypes: []int{lib.COMMAND_INTEGRATION_TYPE_GUILD, lib.COMMAND_INTEGRATION_TYPE_USER},
	Contexts:         []int{lib.COMMAND_CONTEXT_GUILD, lib.COMMAND_CONTEXT_BOT_DM, lib.COMMAND_CONTEXT_PRIVATE_CHANNEL},
	CommandHandler: func(interaction *lib.CommandInteraction) {
		interaction.SendSimpleReply("This is out? Oh damn. Please inform the developer about this.", true)
	},
}
