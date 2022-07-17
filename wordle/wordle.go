package wordle

import (
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Attempt struct {
	Day            int
	MaxAttempts    int
	Attempts       int
	Success        bool
	AttemptsDetail []string
	HardMode       bool
}

var wordleRegex *regexp.Regexp

func init() {
	wordleRegex = regexp.MustCompile(`(?si)Wordle (\d+) (\d|X)\/(\d)(\*?)\n\n((⬜|🟩|🟨|⬛){5}\n?)+`)
}

func ParseCopyPaste(paste string) (*Attempt, bool) {
	if !wordleRegex.MatchString(paste) {
		return nil, false
	}

	matches := wordleRegex.FindStringSubmatch(paste)
	day, _ := strconv.Atoi(matches[1])
	maxAttempts, _ := strconv.Atoi(matches[3])

	success, nAttempts := true, 6
	if matches[2] == "X" {
		success = false
	} else {
		nAttempts, _ = strconv.Atoi(matches[2])
	}

	var hardMode bool
	if matches[4] == "*" {
		hardMode = true
	}

	triesStr := strings.Split(matches[0], "\n\n")[1]
	triesStr = strings.ReplaceAll(triesStr, "⬜", "⬛")
	attemptLines := strings.Split(triesStr, "\n")
	attempts := make([]string, 0)
	for _, attemptLine := range attemptLines {
		l := strings.TrimSpace(attemptLine)

		if l == "" {
			continue
		}

		if utf8.RuneCountInString(l) != 5 {
			return nil, false
		}
		attempts = append(attempts, l)
	}

	if len(attempts) != nAttempts {
		return nil, false
	}

	return &Attempt{
		Day:            day,
		MaxAttempts:    maxAttempts,
		Attempts:       nAttempts,
		Success:        success,
		AttemptsDetail: attempts,
		HardMode:       hardMode,
	}, true
}
