package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/andrerfcsantos/wordle-discord-bot/wordlebot"
)

func main() {
	bot, err := wordlebot.New(&wordlebot.Config{
		Token:             os.Getenv("WORDLE_DISCORD_BOT_TOKEN"),
		AppID:             os.Getenv("WORDLE_DISCORD_BOT_APP_ID"),
		InteractionGuilds: strings.Split(os.Getenv("WORDLE_DISCORD_BOT_INTERACTION_GUILDS"), ","),
	})
	if err != nil {
		panic("error creating new bot: " + err.Error())
	}

	// Cleanly close down the Discord session.
	defer bot.Close()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
