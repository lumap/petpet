package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"petpet/commands"
	subcommands "petpet/commands/subcommands"

	godotenv "github.com/joho/godotenv"

	"petpet/lib"
)

func main() {
	log.Println("Loading environmental variables...")
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln("failed to load env variables", err)
	}

	log.Println("Creating new client...")
	bot := lib.CreateBot(os.Getenv("DISCORD_BOT_TOKEN"), os.Getenv("DISCORD_PUBLIC_KEY"))
	
	log.Println("Registering commands & static components...")
	bot.RegisterCommand(commands.Meow)
	
	bot.RegisterCommand(commands.Petpet)
	bot.RegisterSubCommand(subcommands.PetpetUser, "petpet")
	bot.RegisterSubCommand(subcommands.PetpetImageURL, "petpet")
	bot.RegisterSubCommand(subcommands.PetpetImageUpload, "petpet")

	bot.RegisterCommand(commands.PetpetMsgCtx)
	bot.RegisterCommand(commands.PetpetUserCtx)
	bot.RegisterCommand(commands.PetpetImgCtx)
	

	if os.Getenv("SYNC_COMMANDS") == "1" {
		log.Println("Syncing commands with Discord API...")
		if os.Getenv("TEST_SERVER_ID") != "" {
			testServerID, err := lib.StringToSnowflake(os.Getenv("DISCORD_TEST_SERVER_ID"))
			if err != nil {
				log.Fatalln("failed to parse TEST_SERVER_ID", err)
			}
			err = bot.SyncCommandsWithDiscord([]lib.Snowflake{testServerID})
			if err != nil {
				log.Fatalln("failed to sync commands with Discord API", err)
			}
		} else {
			log.Println("No test server ID provided, syncing commands globally...")
			err := bot.SyncCommandsWithDiscord([]lib.Snowflake{})
			if err != nil {
				log.Fatalln("failed to sync commands with Discord API", err)
			}
		}
	}
	
	http.HandleFunc("POST /", bot.DiscordRequestHandler)
	
	addr := os.Getenv("DISCORD_APP_ADDRESS")
	log.Printf("Serving application at: %s/\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalln("something went terribly wrong", err)
	}
}