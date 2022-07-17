package wordlebot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"strings"
	"text/tabwriter"
)

func (b *WordleBot) HandleTrackInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ok, _ := b.repository.IsTrackedChannel(i.ChannelID)
	var err error
	if ok {
		err = b.session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This channel is already being tracked.",
			},
		})
		return err
	}

	message := "The wordle bot is now tracking this channel for wordle messages. " +
		"Past messages will also be scanned for wordle copy/pastes."
	err = b.session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
	if err != nil {
		return fmt.Errorf("responding to initial interaction: %w", err)
	}

	err = b.repository.TrackChannel(i.ChannelID)
	if err != nil {
		return fmt.Errorf("tracking channel: %w", err)
	}

	procResult, err := b.ProcessChannelMessages(i.ChannelID)
	if err != nil {
		return fmt.Errorf("processing channel messages: %w", err)
	}

	_, err = b.session.InteractionResponseEdit(
		i.Interaction,
		&discordgo.WebhookEdit{
			Content: message + "\n" +
				fmt.Sprintf("📩 Import of old messages is now completed! %d messages were processed "+
					"of which %d were wordle copy/pastes", procResult.TotalMessages, procResult.WordleMessages),
		},
	)
	if err != nil {
		return fmt.Errorf("updating interaction: %w", err)
	}
	return nil
}

func (b *WordleBot) HandleLeaderboardInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	var builder strings.Builder

	tabw := tabwriter.NewWriter(&builder, 2, 2, 2, ' ', tabwriter.TabIndent)

	fmt.Fprintf(tabw, "Name\tScore\tAvg. Attempts\tGames Played\n")
	entries, err := b.repository.Leaderboard(i.ChannelID)
	if err != nil {
		log.Errorf("Error handling leaderboard interaction: %v", err)

		b.session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "There was a problem processing this request, sorry :(",
			},
		})

		return err
	}

	for _, entry := range entries {
		fmt.Fprintf(tabw, "%s\t%.2f\t%.2f\t%d\n",
			entry.Username, entry.TotalScore, entry.AvgAttempts, entry.Played)
	}

	tabw.Flush()
	table := builder.String()

	err = b.session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Here's the current leaderboard:\n" +
				"```\n" + table + "\n```",
		},
	})

	return err
}

func (b *WordleBot) HandleDayInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	err := b.session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponsePong,
	})
	return err
}
