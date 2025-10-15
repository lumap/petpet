package commands

import (
	"petpet/lib"
)

var Petpet lib.Command = lib.Command{
	Name:        "petpet",
	Description: "Petpet someone. Easy.",
	IntegrationTypes: []int{lib.COMMAND_INTEGRATION_TYPE_GUILD, lib.COMMAND_INTEGRATION_TYPE_USER},
	Contexts:         []int{lib.COMMAND_CONTEXT_GUILD, lib.COMMAND_CONTEXT_BOT_DM, lib.COMMAND_CONTEXT_PRIVATE_CHANNEL},
}
