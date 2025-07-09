package commands

import (
	tempest "github.com/amatsagu/tempest"
)

var PetpetImgCtx = tempest.Command{
	Type:             4,
	Name:             "Petpet this image",
	Description:      "",
	SlashCommandHandler: func(itx *tempest.CommandInteraction) {
		itx.SendLinearReply("This is out? Oh damn. Please inform the developer about this.", true)
	},
}
