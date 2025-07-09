package petpetsubcommands

import (
	"petpet/lib"
	"petpet/pet_maker"
	"petpet/utils"
	"slices"
)

var PetpetUser = lib.Command{
	Name: "user",
	Description: "Petpet someone's pfp",
	Options: append(utils.PetpetCommandUserOptions, utils.PetpetCommandOptions...),
	CommandHandler: func(itx *lib.CommandInteraction) {

		untypedUser, err := itx.GetStringOptionValue("user_to_petpet", "")
		if err != nil {
			itx.SendSimpleReply("Invalid user ID provided.", true)
			return
		}
		userId, err := lib.StringToSnowflake(untypedUser)
		if err != nil {
			itx.SendSimpleReply("Couldn't parse user ID. You shouldn't see this.", true)
			return
		}

		if slices.Contains(utils.BlacklistedUsers, userId) {
			itx.SendSimpleReply("This user is blacklisted, sorry.", true)
			return
		}

		member := itx.Data.Resolved.Members[userId]
		user := itx.Data.Resolved.Users[userId]

		useServerAvatar, err := itx.GetBoolOptionValue("use_server_avatar", true)
		if err != nil {
			itx.SendSimpleReply("Couldn't parse use_server_avatar option. You shouldn't see this.", true)
			return
		}

		avatar := member.GuildAvatarURL()
		if !useServerAvatar || avatar == "" {
			avatar = user.AvatarURL()
		}

		ephemeral, err := itx.GetBoolOptionValue("ephemeral", false)
		if err != nil {
			itx.SendSimpleReply("Couldn't parse ephemeral option. You shouldn't see this.", true)
			return
		}

		itx.Defer(ephemeral)

		speed, err := itx.GetFloatOptionValue("speed", 1.0)
		if err != nil {
			itx.SendSimpleReply("Couldn't parse speed option. You shouldn't see this.", true)
			return
		}
		width, err := itx.GetIntOptionValue("width", 128)
		if err != nil {
			itx.SendSimpleReply("Couldn't parse width option. You shouldn't see this.", true)
			return
		}
		height, err := itx.GetIntOptionValue("height", 128)
		if err != nil {
			itx.SendSimpleReply("Couldn't parse height option. You shouldn't see this.", true)
			return
		}

		img := pet_maker.MakePetImage(avatar, speed, width, height)

		itx.EditReply(lib.ResponseMessageData{}, ephemeral, []lib.DiscordFile{
			{
				Filename: "petpet.gif",
				Reader:   img,
			},
		})
	},
}