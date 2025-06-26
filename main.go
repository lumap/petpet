package main

import (
	"petpet/commands"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	tempest "github.com/amatsagu/tempest"
	godotenv "github.com/joho/godotenv"
)

func main() {
	log.Println("Loading environmental variables...")
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln("failed to load env variables", err)
	}

	log.Println("Creating new Tempest client...")
	client := tempest.NewClient(tempest.ClientOptions{
		Token:     os.Getenv("DISCORD_BOT_TOKEN"),
		PublicKey: os.Getenv("DISCORD_PUBLIC_KEY"),
	})

	testServerID, err := tempest.StringToSnowflake(os.Getenv("DISCORD_TEST_SERVER_ID")) // Register example commands only to this guild.
	if err != nil {
		log.Fatalln("failed to parse env variable to snowflake", err)
	}
	
	log.Println("Registering commands & static components...")
	client.RegisterCommand(commands.Meow)
	
	err = client.SyncCommandsWithDiscord([]tempest.Snowflake{testServerID}, nil, false)
	if err != nil {
		log.Fatalln("failed to sync local commands storage with Discord API", err)
	}
	
	http.HandleFunc("POST /", client.DiscordRequestHandler)
	
	addr := os.Getenv("DISCORD_APP_ADDRESS")
	log.Printf("Serving application at: %s/\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalln("something went terribly wrong", err)
	}
}