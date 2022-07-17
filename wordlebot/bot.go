package wordlebot

import (
	"fmt"
	"github.com/andrerfcsantos/wordle-discord-bot/db"
	"github.com/bwmarrin/discordgo"
)

type Config struct {
	Token             string
	AppID             string
	InteractionGuilds []string
}

func New(config *Config) (*WordleBot, error) {
	var bot WordleBot
	var err error

	bot.config = *config
	bot.config.InteractionGuilds = make([]string, len(config.InteractionGuilds))
	copy(bot.config.InteractionGuilds, config.InteractionGuilds)

	if !sliceHasString(bot.config.InteractionGuilds, "") {
		bot.config.InteractionGuilds = append([]string{""}, bot.config.InteractionGuilds...)
	}

	bot.session, err = discordgo.New("Bot " + bot.config.Token)
	if err != nil {
		return nil, fmt.Errorf("creating new bot: %w", err)
	}

	// bot.session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	bot.session.AddHandler(bot.MessageCreateHandler)

	err = bot.session.Open()
	if err != nil {
		return nil, fmt.Errorf("opening connection to discord: %w", err)
	}

	err = bot.setupApplicationCommands()
	if err != nil {
		return nil, fmt.Errorf("setting up application commands: %w", err)
	}

	bot.repository, err = db.NewRepository()
	if err != nil {
		return nil, fmt.Errorf("creating bot DB repository: %w", err)
	}

	err = bot.repository.RunMigrations()
	if err != nil {
		return nil, fmt.Errorf("migrating database: %w", err)
	}

	return &bot, nil
}

type WordleBot struct {
	session    *discordgo.Session
	repository *db.Repository
	appCmd     *discordgo.ApplicationCommand
	config     Config
}

func (b *WordleBot) Close() error {
	return b.session.Close()
}
