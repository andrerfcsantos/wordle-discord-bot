package wordlebot

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/andrerfcsantos/wordle-discord-bot/db"
	"github.com/andrerfcsantos/wordle-discord-bot/wordle"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"strings"
	"unicode/utf8"
)

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
		log.Errorf("failed to save wordle message: %v\n", err)
		return
	}

	err = b.session.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	if err != nil {
		log.Errorf("failed to add reaction: %v\n", err)
	}
}

func (b *WordleBot) DeleteMessageHandler(s *discordgo.Session, m *discordgo.MessageDelete) {
	tracked, _ := b.repository.IsTrackedChannel(m.ChannelID)
	if !tracked {
		return
	}

	_, err := b.repository.AttemptForMessage(m.ChannelID, m.ID)
	if err != nil {
		log.Errorf("failed to get attempt for message: %v\n", err)
		return
	}

	_, err = b.repository.DeleteAttemptForMessage(m.ChannelID, m.ID)
	if err != nil {
		log.Errorf("failed to delete attempt: %v\n", err)
	}
}

func (b *WordleBot) UpdateMessageHandler(s *discordgo.Session, m *discordgo.MessageUpdate) {
	tracked, _ := b.repository.IsTrackedChannel(m.ChannelID)
	if !tracked {
		return
	}

	a, err := b.repository.AttemptForMessage(m.ChannelID, m.ID)
	if err != nil {
		log.Errorf("failed to get attempt for message: %v\n", err)
		return
	}

	attempt, ok := wordle.ParseCopyPaste(m.Content)

	if a.MessageId == m.ID && !ok {
		_, err := b.repository.DeleteAttemptForMessage(m.ChannelID, m.ID)
		if err != nil {
			log.Errorf("failed to delete attempt: %v\n", err)
		}
	} else {
		if ok {
			err := b.saveWordleMessage(m.Message, attempt)
			if err != nil {
				log.Errorf("failed to save wordle message: %v\n", err)
			}
		}
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

			err = b.session.MessageReactionAdd(m.ChannelID, m.ID, "✅")
			if err != nil {
				log.Errorf("failed to add reaction: %v\n", err)
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
		MessageId:    m.ID,
		ChannelId:    m.ChannelID,
		UserId:       m.Author.ID,
		Day:          attempt.Day,
		UserName:     m.Author.Username,
		Attempts:     attempt.Attempts,
		MaxAttempts:  attempt.MaxAttempts,
		Success:      attempt.Success,
		AttemptsJson: string(attemptsJson),
		PostedAt:     m.Timestamp,
		HardMode:     attempt.HardMode,
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
