package wordlebot

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"text/tabwriter"
	"unicode/utf8"

	log "github.com/sirupsen/logrus"

	"github.com/andrerfcsantos/wordle-discord-bot/db"
	"github.com/andrerfcsantos/wordle-discord-bot/wordle"
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
		b.config.AppID,
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

	fmt.Fprintf(tabw, "Name\tAvg. Score\tCount\n")
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
		fmt.Fprintf(tabw, "%s\t%.2f\t%d\n", entry.Username, entry.AvgScore, entry.Count)
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

func (b *WordleBot) MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	tracked, _ := b.repository.IsTrackedChannel(m.ChannelID)
	if !tracked {
		return
	}

	attempt, ok := wordle.ParseCopyPaste(m.Content)
	if !ok {
		return
	}

	err := b.saveWordleMessage(m.Message, attempt)
	if err != nil {
		return
	}

	var reaction string
	if attempt.Success {
		reaction = "✅"
	} else {
		reaction = "❌"
	}

	err = b.session.MessageReactionAdd(m.ChannelID, m.ID, reaction)
	if err != nil {
		log.Errorf("failed to add reaction: %v\n", err)
	}
}

type ProcessResult struct {
	WordleMessages int
	TotalMessages  int
}

func (b *WordleBot) ProcessChannelMessages(channelId string) (*ProcessResult, error) {
	messages, err := b.session.ChannelMessages(channelId, 100, "", "", "")
	if err != nil {
		return nil, fmt.Errorf("getting channel messages: %v", err)
	}

	var result ProcessResult

	for len(messages) > 0 {

		result.TotalMessages += len(messages)
		for _, m := range messages {
			attempt, ok := wordle.ParseCopyPaste(m.Content)
			if !ok {
				continue
			}

			err := b.saveWordleMessage(m, attempt)
			if err != nil {
				return nil, fmt.Errorf("saving wordle message: %v", err)
			}
			result.WordleMessages++
		}

		messages, err = b.session.ChannelMessages(channelId, 100, messages[len(messages)-1].ID, "", "")
		if err != nil {
			return nil, fmt.Errorf("getting channel messages: %v", err)
		}
	}

	return &result, nil
}

func (b *WordleBot) saveWordleMessage(m *discordgo.Message, attempt *wordle.Attempt) error {
	attemptsJson, err := json.Marshal(attempt.AttemptsDetail)
	if err != nil {
		return errors.New("failed to marshal attempts detail")
	}

	err = b.repository.SaveAttempt(db.Attempt{
		ChannelId:    m.ChannelID,
		UserId:       m.Author.ID,
		Day:          attempt.Day,
		UserName:     m.Author.Username,
		Attempts:     attempt.Attempts,
		MaxAttempts:  attempt.MaxAttempts,
		Success:      attempt.Success,
		AttemptsJson: string(attemptsJson),
		PostedAt:     m.Timestamp,
		Score:        attempt.Score,
	})

	if err != nil {
		return fmt.Errorf("saving attempt: %w", err)
	}

	return nil
}

func SanitizeMessage(message string) string {
	sanitized := strings.ReplaceAll(message, "\n", "")
	if utf8.RuneCountInString(sanitized) > 1024 {
		sanitized = string([]rune(sanitized)[:1024]) + "..."
	}

	return sanitized
}

func sliceHasString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
