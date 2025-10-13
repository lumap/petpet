package petpetsubcommands

import (
	"fmt"
	"petpet/lib"
	"petpet/pet_maker"
	"petpet/utils"
	"slices"
	"strings"
)

func processUser(interaction *lib.CommandInteraction, files *[]lib.DiscordFile, untypedUserId string) {
	userId, err := lib.StringToSnowflake(untypedUserId)
	if err != nil {
		interaction.SendSimpleReply("Couldn't parse user ID. What have you done? How? :/", true)
		return
	}

	if slices.Contains(utils.BlacklistedUsers, userId) {
		interaction.SendSimpleReply("This user is blacklisted, sorry.", true)
		return
	}

	member := interaction.Data.Resolved.Members[userId]
	user := interaction.Data.Resolved.Users[userId]

	useServerAvatar, err := interaction.GetBoolOptionValue("use_server_avatar", true)
	if err != nil {
		interaction.SendSimpleReply("Couldn't parse use_server_avatar option.", true)
		return
	}

	avatar := user.AvatarURL()
	if useServerAvatar && member != nil {
		avatar = member.GuildAvatarURL()
	}

	speed, err := interaction.GetFloatOptionValue("speed", 1.0)
	if err != nil {
		interaction.SendSimpleReply("Couldn't parse speed option.", true)
		return
	}
	width, err := interaction.GetIntOptionValue("width", 128)
	if err != nil {
		interaction.SendSimpleReply("Couldn't parse width option.", true)
		return
	}
	height, err := interaction.GetIntOptionValue("height", 128)
	if err != nil {
		interaction.SendSimpleReply("Couldn't parse height option.", true)
		return
	}

	img := pet_maker.MakePetImage(avatar, speed, width, height)

	*files = append(*files, lib.DiscordFile{
		Filename: "petpet.gif",
		Reader:   img,
	})
}

var PetpetUser = lib.Command{
	Name:        "user",
	Description: "Petpet someone's pfp",
	Options:     append(utils.PetpetCommandUserOptions, utils.PetpetCommandOptions...),
	CommandHandler: func(interaction *lib.CommandInteraction) {
		ephemeral, err := interaction.GetBoolOptionValue("ephemeral", false)
		if err != nil {
			interaction.SendSimpleReply("Couldn't parse ephemeral option.", true)
			return
		}

		interaction.Defer(ephemeral)

		files := []lib.DiscordFile{}
		userIDs := []string{}

		// it can support up to 10 but i'm not gonna bother implementing more than 4 until i figure out monetization
		for i := 1; i <= 4; i++ {
			optionName := "user"
			if i > 1 {
				optionName = optionName + fmt.Sprint(i)
			}
			user, err := interaction.GetStringOptionValue(optionName, "")
			if err == nil && user != "" {
				userIDs = append(userIDs, user)
				processUser(interaction, &files, user)
			}
		}

		// format user IDs to be mentions
		for i, id := range userIDs {
			userIDs[i] = fmt.Sprintf("<@%s>", id)
		}

		mentionedUsers := strings.Join(userIDs, ", ")
		// replace the last ", " with " and " if there are multiple users
		if len(userIDs) > 1 {
			lastComma := strings.LastIndex(mentionedUsers, ", ")
			mentionedUsers = mentionedUsers[:lastComma] + " and" + mentionedUsers[lastComma+1:]
		}

		interaction.EditReply(lib.ResponseMessageData{
			Content:         "<@" + interaction.GetUser().ID.String() + "> has pet " + mentionedUsers + " :3",
			AllowedMentions: &lib.AllowedMentions{},
		}, ephemeral, files)
	},
}
