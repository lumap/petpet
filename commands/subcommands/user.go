package petpetsubcommands

import (
	"github.com/amatsagu/tempest"
	"petpet/utils"
	"slices"
)

var PetpetUser = tempest.Command{
	Name: "user",
	Description: "Petpet someone's pfp",
	Options: append(utils.PetpetCommandUserOptions, utils.PetpetCommandOptions...),
	SlashCommandHandler: func(itx *tempest.CommandInteraction) {

		// get user and its avatar
		untypedUser, _ := itx.GetOptionValue("user_to_petpet")
		user := untypedUser.(tempest.User)

		if slices.Contains(utils.BlacklistedUsers, user.ID) {
			itx.SendLinearReply("This user is blacklisted, sorry.", true)
			return
		}

		member := itx.ResolveMember(user.ID)
		var avatar = ""
		if member == nil {
			avatar = user.AvatarURL()
		} else {
			avatar = member.GuildAvatarURL()
			if avatar == "" {
				avatar = member.User.AvatarURL()
			}
		}

		// defer
		ephemeral, present := itx.GetOptionValue("ephemeral")
		if !present {
			ephemeral = false
		}
		itx.Defer(ephemeral.(bool))

		itx.SendFollowUp(tempest.ResponseMessageData{
			Attachments: []tempest.Attachment{

			},
		}, ephemeral.(bool))
	},
}