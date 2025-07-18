package petpetsubcommands

import (
	"petpet/lib"
	"petpet/pet_maker"
	"petpet/utils"
)

var PetpetImageURL = lib.Command{
	Name:        "image_url",
	Description: "Petpet an image (via external URL)",
	Options:     append(utils.PetpetCommandImageURLOptions, utils.PetpetCommandOptions...),
	CommandHandler: func(interaction *lib.CommandInteraction) {

		imageURL, err := interaction.GetStringOptionValue("image_url", "")
		if err != nil {
			interaction.SendSimpleReply("Invalid image URL provided.", true)
			return
		}

		isImage, err := utils.IsLinkAnImageURL(imageURL)
		if err != nil {
			interaction.SendSimpleReply("Couldn't check if the URL is an image. You shouldn't see this.", true)
			return
		}
		if !isImage {
			interaction.SendSimpleReply("The provided URL is not an image.", true)
			return
		}

		ephemeral, err := interaction.GetBoolOptionValue("ephemeral", false)
		if err != nil {
			interaction.SendSimpleReply("Couldn't parse ephemeral option. You shouldn't see this.", true)
			return
		}

		interaction.Defer(ephemeral)

		speed, err := interaction.GetFloatOptionValue("speed", 1.0)
		if err != nil {
			interaction.SendSimpleReply("Couldn't parse speed option. You shouldn't see this.", true)
			return
		}
		width, err := interaction.GetIntOptionValue("width", 128)
		if err != nil {
			interaction.SendSimpleReply("Couldn't parse width option. You shouldn't see this.", true)
			return
		}
		height, err := interaction.GetIntOptionValue("height", 128)
		if err != nil {
			interaction.SendSimpleReply("Couldn't parse height option. You shouldn't see this.", true)
			return
		}

		img := pet_maker.MakePetImage(imageURL, speed, width, height)

		interaction.EditReply(lib.ResponseMessageData{}, ephemeral, []lib.DiscordFile{
			{
				Filename: "petpet.gif",
				Reader:   img,
			},
		})
	},
}
