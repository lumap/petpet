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

	member := interaction.Data.Resolved.Members[userId]
	user := interaction.Data.Resolved.Users[userId]

	useServerAvatar, err := interaction.GetBoolOptionValue("use_server_avatar", true)
	if err != nil {
		interaction.SendSimpleReply("Couldn't parse use_server_avatar option.", true)
		return
	}

	avatar := user.AvatarURL()
	if useServerAvatar && member != nil && member.GuildAvatarHash != "" {
		avatar = member.GuildAvatarURL(interaction.GuildID.String(), userId.String())
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
		mentions := []string{}
		if notify, err := interaction.GetBoolOptionValue("notify_users", false); err != nil {
			interaction.SendSimpleReply("Couldn't parse notify_users option.", true)
			return
		} else if notify {
			mentions = append(mentions, "users")
		}

		interaction.Defer(ephemeral)

		files := []lib.DiscordFile{}
		userIDs := []string{}
		blacklistDetected := false

		for i := 1; i <= 10; i++ {
			optionName := "user"
			if i > 1 {
				optionName = optionName + fmt.Sprint(i)
			}
			user, err := interaction.GetStringOptionValue(optionName, "")
			if err == nil && user != "" {
				if slices.Contains(utils.BlacklistedUsers, user) {
					blacklistDetected = true
					continue
				}
				userIDs = append(userIDs, user)
				processUser(interaction, &files, user)
			}
		}

		content := "<@" + interaction.GetUser().ID.String() + "> has pet "
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
		content += mentionedUsers + " :3"
		if blacklistDetected {
			content += "\n-# A user you tried to petpet has been blacklisted and was ignored."
		}

		interaction.EditReply(lib.ResponseMessageData{
			Content:         content,
			AllowedMentions: &lib.AllowedMentions{Parse: mentions},
		}, ephemeral, files)
	},
}
