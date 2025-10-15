package commands

import (
	"petpet/lib"
	"petpet/pet_maker"
	"petpet/utils"
	"slices"
)

var PetpetMsgCtx = lib.Command{
	Type:        3,
	Name:        "Petpet the message's author",
	Description: "",
	CommandHandler: func(interaction *lib.CommandInteraction) {
		messageId := interaction.Data.TargetID
		lib.LogInfo(interaction.Data.Resolved.String())
		user := interaction.Data.Resolved.Messages[messageId].Author

		if slices.Contains(utils.BlacklistedUsers, user.ID.String()) {
			interaction.SendSimpleReply("This user is blacklisted, sorry.", true)
			return
		}

		avatar := user.AvatarURL()

		interaction.Defer(false)

		img := pet_maker.MakePetImage(avatar, 1, 128, 128)

		interaction.EditReply(lib.ResponseMessageData{
			Content: "<@" + interaction.GetUser().ID.String() + "> has pet <@" + user.ID.String() + "> :33",
		}, false, []lib.DiscordFile{
			{
				Filename: "petpet.gif",
				Reader:   img,
			},
		})
	},
}
