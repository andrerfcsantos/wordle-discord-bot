package wordle_test

import (
	"testing"

	"github.com/andrerfcsantos/wordle-discord-bot/wordle"
)

// https://go.dev/play/p/8xrll8rJxcH

func TestParseCopyPaste(t *testing.T) {
	tests := []struct {
		name    string
		paste   string
		want    *wordle.Attempt
		wantErr bool
	}{
		{
			name: "paste with game won",
			paste: `Wordle 217 3/6

⬛⬛⬛🟩🟩
⬛⬛🟨🟩🟩
🟩🟩🟩🟩🟩`,
			want: &wordle.Attempt{
				Day:         217,
				MaxAttempts: 6,
				Attempts:    3,
				Success:     true,
				AttemptsDetail: []string{
					"⬛⬛⬛🟩🟩",
					"⬛⬛🟨🟩🟩",
					"🟩🟩🟩🟩🟩",
				},
			},
		},
		{
			name: "paste with game won with message after",
			paste: `Wordle 219 4/6

⬛⬛⬛⬛⬛
⬛🟨⬛🟨⬛
🟩🟩🟩⬛⬛
🟩🟩🟩🟩🟩
Bot, be good`,
			want: &wordle.Attempt{
				Day:         219,
				MaxAttempts: 6,
				Attempts:    4,
				Success:     true,
				AttemptsDetail: []string{
					"⬛⬛⬛⬛⬛",
					"⬛🟨⬛🟨⬛",
					"🟩🟩🟩⬛⬛",
					"🟩🟩🟩🟩🟩",
				},
			},
			wantErr: false,
		},
		{
			name: "paste with game won with message after and before",
			paste: `Bot, be good
Wordle 219 4/6

⬛⬛⬛⬛⬛
⬛🟨⬛🟨⬛
🟩🟩🟩⬛⬛
🟩🟩🟩🟩🟩
Bot, be good`,
			want: &wordle.Attempt{
				Day:         219,
				MaxAttempts: 6,
				Attempts:    4,
				Success:     true,
				AttemptsDetail: []string{
					"⬛⬛⬛⬛⬛",
					"⬛🟨⬛🟨⬛",
					"🟩🟩🟩⬛⬛",
					"🟩🟩🟩🟩🟩",
				},
			},
			wantErr: false,
		},
		{
			name: "paste with game lost with message after and before",
			paste: `Bot, be good
Wordle 219 X/6

⬛⬛⬛⬛⬛
⬛🟨⬛🟨⬛
🟩🟩🟩⬛⬛
🟩🟩🟩⬛⬛
🟩🟩🟩⬛⬛
🟩🟩🟩⬛⬛
Bot, be good`,
			want: &wordle.Attempt{
				Day:         219,
				MaxAttempts: 6,
				Attempts:    6,
				Success:     false,
				AttemptsDetail: []string{
					"⬛⬛⬛⬛⬛",
					"⬛🟨⬛🟨⬛",
					"🟩🟩🟩⬛⬛",
					"🟩🟩🟩⬛⬛",
					"🟩🟩🟩⬛⬛",
					"🟩🟩🟩⬛⬛",
				},
			},
			wantErr: false,
		},
		{
			name: "paste with invalid game due to incomplete row",
			paste: `Wordle 219 X/6

⬛⬛⬛⬛⬛
⬛🟨⬛🟨⬛
🟩🟩🟩⬛
🟩🟩🟩⬛⬛
🟩🟩🟩⬛⬛
🟩🟩🟩⬛⬛`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := wordle.ParseCopyPaste(tt.paste)

			if ok && tt.wantErr {
				t.Fatal("ParseCopyPaste() returned OK, but error wanted")
			}

			if !ok && !tt.wantErr {
				t.Fatal("ParseCopyPaste() returned error, but none was expected")
			}

			if !ok && tt.wantErr {
				return
			}

			if (got == nil || tt.want == nil) && got != tt.want {
				t.Fatalf("ParseCopyPaste() = %v, want %v", got, tt.want)
			}

			if !attemptsEqual(got, tt.want) {
				t.Fatalf("ParseCopyPaste() = %v, want %v", got, tt.want)
			}
		})
	}
}

func attemptsEqual(a, b *wordle.Attempt) bool {
	if a.Day != b.Day {
		return false
	}
	if a.MaxAttempts != b.MaxAttempts {
		return false
	}
	if a.Attempts != b.Attempts {
		return false
	}
	if a.Success != b.Success {
		return false
	}
	if len(a.AttemptsDetail) != len(b.AttemptsDetail) {
		return false
	}
	for i := range a.AttemptsDetail {
		if a.AttemptsDetail[i] != b.AttemptsDetail[i] {
			return false
		}
	}
	return true
}
