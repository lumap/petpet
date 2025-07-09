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

	testServerID, err := lib.StringToSnowflake(os.Getenv("DISCORD_TEST_SERVER_ID")) // Register example commands only to this guild.
	if err != nil {
		log.Fatalln("failed to parse env variable to snowflake", err)
	}
	
	log.Println("Registering commands & static components...")
	bot.RegisterCommand(commands.Meow)
	bot.RegisterCommand(commands.Petpet)
	bot.RegisterSubCommand(subcommands.PetpetUser, "petpet")
	
	err = bot.SyncCommandsWithDiscord([]lib.Snowflake{testServerID})
	if err != nil {
		log.Fatalln("failed to sync local commands storage with Discord API", err)
	}
	
	http.HandleFunc("POST /", bot.DiscordRequestHandler)
	
	addr := os.Getenv("DISCORD_APP_ADDRESS")
	log.Printf("Serving application at: %s/\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalln("something went terribly wrong", err)
	}
}