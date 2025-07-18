package commands

import (
	"petpet/lib"
)

var PetpetImgCtx = lib.Command{
	Type:        4,
	Name:        "Petpet this image",
	Description: "",
	CommandHandler: func(interaction *lib.CommandInteraction) {
		interaction.SendSimpleReply("This is out? Oh damn. Please inform the developer about this.", true)
	},
}
