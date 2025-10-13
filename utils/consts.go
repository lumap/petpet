package utils

import "petpet/lib"

var BlacklistedUsers = []lib.Snowflake{
	1118171067463241868,
}

var PetpetCommandOptions = []lib.CommandOption{
	{
		Type:        4,
		Name:        "width",
		Description: "The width of the gif (default is 128)",
		Required:    false,
		MinValue:    8,
		MaxValue:    1024,
	},
	{
		Type:        4,
		Name:        "height",
		Description: "The height of the gif (default is 128)",
		Required:    false,
		MinValue:    8,
		MaxValue:    1024,
	},
	{
		Type:        10,
		Name:        "speed",
		Description: "How fast the petting is (default is 1, min 0.125, max 3)",
		Required:    false,
		MinValue:    0.125,
		MaxValue:    3,
	},
	{
		Type:        5,
		Name:        "ephemeral",
		Description: "Whether or not to make the message ephemeral (default is false)",
		Required:    false,
	},
}

var PetpetCommandUserOptions = []lib.CommandOption{
	{
		Type:        6,
		Name:        "user",
		Description: "The user to petpet",
		Required:    true,
	},
	{
		Type:        6,
		Name:        "user2",
		Description: "The second user to petpet",
		Required:    false,
	},
	{
		Type:        6,
		Name:        "user3",
		Description: "The third user to petpet",
		Required:    false,
	},
	{
		Type:        6,
		Name:        "user4",
		Description: "The fourth user to petpet",
		Required:    false,
	},
	// {
	// 	Type:        6,
	// 	Name:        "user5",
	// 	Description: "The fifth user to petpet",
	// 	Required:    false,
	// },
	// {
	// 	Type:        6,
	// 	Name:        "user6",
	// 	Description: "The sixth user to petpet",
	// 	Required:    false,
	// },
	// {
	// 	Type:        6,
	// 	Name:        "user7",
	// 	Description: "The seventh user to petpet",
	// 	Required:    false,
	// },
	// {
	// 	Type:        6,
	// 	Name:        "user8",
	// 	Description: "The eighth user to petpet",
	// 	Required:    false,
	// },
	// {
	// 	Type:        6,
	// 	Name:        "user9",
	// 	Description: "The ninth user to petpet",
	// 	Required:    false,
	// },
	// {
	// 	Type:        6,
	// 	Name:        "user10",
	// 	Description: "The tenth user to petpet",
	// 	Required:    false,
	// },
	{
		Type:        5,
		Name:        "use_server_avatar",
		Description: "Whether to use the server avatars (default is true; applies to all users)",
		Required:    false,
	},
}

var PetpetCommandImageURLOptions = []lib.CommandOption{
	{
		Type:        3,
		Name:        "image_url",
		Description: "The image's URL",
		Required:    true,
	},
}

var PetpetCommandImageUploadOptions = []lib.CommandOption{
	{
		Type:        11,
		Name:        "image_upload",
		Description: "The image to petpet",
		Required:    true,
	},
}
