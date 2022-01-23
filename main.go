package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/andrerfcsantos/wordle-discord-bot/wordlebot"
	log "github.com/sirupsen/logrus"
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

	file, err := os.OpenFile("logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	// Cleanly close down the Discord session.
	defer bot.Close()

	// Wait here until CTRL-C or other term signal is received.
	log.Info("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
