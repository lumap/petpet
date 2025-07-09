package commands

import (
	tempest "github.com/amatsagu/tempest"
	// "petpet/pet_maker"
	// "petpet/utils"
)

var Petpet tempest.Command = tempest.Command{
	Name:             "petpet",
	Description:      "Petpet someone. Easy.",
	// Type:             1,
	// Options: []tempest.CommandOption{
	// 	{
	// 		Name: "user",
	// 		Description: "Petpet someone's pfp",
	// 		Type: tempest.SUB_OPTION_TYPE,
	// 		Options: append(utils.PetpetCommandUserOptions, utils.PetpetCommandOptions...),
	// 	},
	// 	{
	// 		Name:        "image_via_upload",
	// 		Description: "Petpet an image via upload",
	// 		Type: tempest.SUB_OPTION_TYPE,
	// 		Options: append([]tempest.CommandOption{{
	// 			Type:        11,
	// 			Name:        "image",
	// 			Description: "The image to petpet",
	// 			Required:    true,
	// 		},}, utils.PetpetCommandOptions...),
	// 	},
	// 	{
	// 		Name:        "image_via_url",
	// 		Description: "Petpet an image via URL",
	// 		Type: tempest.SUB_OPTION_TYPE,
	// 		Options: append([]tempest.CommandOption{{
	// 			Type:        3,
	// 			Name:        "image_url",
	// 			Description: "The URL of the image to petpet",
	// 			Required:    true,
	// 		},}, utils.PetpetCommandOptions...),
	// 	},
	// },
}