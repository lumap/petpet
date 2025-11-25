package commands

import (
	"petpet/lib"
)

var Donate lib.Command = lib.Command{
	Name:        "donate",
	Description: "Want to support PetPet? This is how!",
	Type:        lib.COMMAND_TYPE_CHAT_INPUT,
	IntegrationTypes: []int{lib.COMMAND_INTEGRATION_TYPE_GUILD, lib.COMMAND_INTEGRATION_TYPE_USER},
	Contexts:         []int{lib.COMMAND_CONTEXT_GUILD, lib.COMMAND_CONTEXT_BOT_DM, lib.COMMAND_CONTEXT_PRIVATE_CHANNEL},
	CommandHandler: func(interaction *lib.CommandInteraction) {
		interaction.SendSimpleReply("Hey there! If you're interested in supporting PetPet, you can do do through this link: https://ko-fi.com/lumap.\nPetPet is and will remain free to use for everyone forever. You are not required to help, but it would be greatly appreciated!", true)
	},
}
