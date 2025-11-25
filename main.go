package main

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"petpet/commands"
	subcommands "petpet/commands/subcommands"

	godotenv "github.com/joho/godotenv"

	"petpet/lib"
)

func main() {
	lib.LogInfo("Loading environmental variables...")
	if err := godotenv.Load(".env"); err != nil {
		lib.LogError("failed to load env variables", err)
	}

	lib.LogInfo("Creating new client...")
	bot := lib.CreateBot(os.Getenv("DISCORD_BOT_TOKEN"), os.Getenv("DISCORD_PUBLIC_KEY"))

	lib.LogInfo("Registering commands & static components...")
	bot.RegisterCommand(commands.Meow)
	bot.RegisterCommand(commands.Donate)

	bot.RegisterCommand(commands.Petpet)
	bot.RegisterSubCommand(subcommands.PetpetUser, "petpet")
	bot.RegisterSubCommand(subcommands.PetpetImageURL, "petpet")
	bot.RegisterSubCommand(subcommands.PetpetImageUpload, "petpet")

	bot.RegisterCommand(commands.PetpetMsgCtx)
	bot.RegisterCommand(commands.PetpetUserCtx)
	bot.RegisterCommand(commands.PetpetImgCtx)

	if os.Getenv("SYNC_COMMANDS") == "1" {
		lib.LogInfo("Syncing commands with Discord API...")
		if os.Getenv("TEST_SERVER_ID") != "" {
			testServerID, err := lib.StringToSnowflake(os.Getenv("TEST_SERVER_ID"))
			if err != nil {
				lib.LogError("failed to parse TEST_SERVER_ID", err)
			}
			if err = bot.SyncCommandsWithDiscord([]lib.Snowflake{testServerID}); err != nil {
				lib.LogError("failed to sync commands with Discord API", err)
			}
		} else {
			lib.LogInfo("No test server ID provided, syncing commands globally...")
			if err := bot.SyncCommandsWithDiscord([]lib.Snowflake{}); err != nil {
				lib.LogError("failed to sync commands with Discord API", err)
			}
		}
	}

	http.HandleFunc("POST /", bot.DiscordRequestHandler)

	addr := os.Getenv("DISCORD_APP_ADDRESS")
	lib.LogInfo("Serving application at: %s/\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		lib.LogError("something went terribly wrong", err)
	}
}
