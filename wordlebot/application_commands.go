package wordlebot

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func (b *WordleBot) setupApplicationCommands() error {
	command := &discordgo.ApplicationCommand{
		Name:        "wordle",
		Type:        discordgo.ChatApplicationCommand,
		Description: "Wordle stats",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "track",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Description: "Marks a channel to be tracked for past and future wordle messages.",
			},
			{
				Name:        "leaderboard",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Description: "Displays the leaderboard for the current channel.",
			},
		},
	}

	for _, guild := range b.config.InteractionGuilds {
		cmd, err := b.session.ApplicationCommandCreate(
			b.config.AppID,
			guild,
			command,
		)
		if err != nil {
			log.Errorf("creating application command for guild %q: %v\n", guild, err)
		}

		if guild == "" {
			b.appCmd = cmd
		}
	}

	b.session.AddHandler(b.ApplicationCommandHandler)
	b.session.AddHandler(b.DeleteMessageHandler)
	b.session.AddHandler(b.UpdateMessageHandler)

	return nil
}

func (b *WordleBot) ApplicationCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()

	var err error
	switch data.Options[0].Name {
	case "track":
		err = b.HandleTrackInteraction(s, i)
	case "leaderboard":
		err = b.HandleLeaderboardInteraction(s, i)
	}

	if err != nil {
		log.Errorf("responding to interaction: %v\n", err)
	}
}
