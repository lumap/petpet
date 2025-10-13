package petpetsubcommands

import (
	"petpet/lib"
	"petpet/pet_maker"
	"petpet/utils"
)

var PetpetImageUpload = lib.Command{
	Name:        "image_upload",
	Description: "Petpet an uploaded image",
	Options:     append(utils.PetpetCommandImageUploadOptions, utils.PetpetCommandOptions...),
	CommandHandler: func(interaction *lib.CommandInteraction) {

		untypedImage, err := interaction.GetAttachmentOptionId("image_upload", "")
		if err != nil {
			interaction.SendSimpleReply("Invalid image URL provided.", true)
			return
		}

		imageId, err := lib.StringToSnowflake(untypedImage)
		if err != nil {
			interaction.SendSimpleReply("Couldn't parse image ID.", true)
			return
		}

		image := interaction.Data.Resolved.Attachments[imageId]

		isImage, err := utils.IsLinkAnImageURL(image.URL)
		if err != nil {
			interaction.SendSimpleReply("Couldn't check if the URL is an image.", true)
			return
		}
		if !isImage {
			interaction.SendSimpleReply("The provided URL is not an image.", true)
			return
		}

		ephemeral, err := interaction.GetBoolOptionValue("ephemeral", false)
		if err != nil {
			interaction.SendSimpleReply("Couldn't parse ephemeral option.", true)
			return
		}

		interaction.Defer(ephemeral)

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

		img := pet_maker.MakePetImage(image.URL, speed, width, height)

		interaction.EditReply(lib.ResponseMessageData{
			Content:         "<@" + interaction.GetUser().ID.String() + "> has pet an uploaded image :3",
			AllowedMentions: &lib.AllowedMentions{},
		}, ephemeral, []lib.DiscordFile{
			{
				Filename: "petpet.gif",
				Reader:   img,
			},
		})
	},
}
