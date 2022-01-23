package wordle_test

import (
	"testing"

	"github.com/andrerfcsantos/wordle-discord-bot/wordle"
)

func TestParseCopyPaste(t *testing.T) {
	tests := []struct {
		name    string
		paste   string
		want    *wordle.Attempt
		wantErr bool
	}{
		{
			name: "success",
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
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := wordle.ParseCopyPaste(tt.paste)
			if ok == tt.wantErr {
				t.Errorf("ok = %v, wantErr = %v", ok, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil {
				if got.Day != tt.want.Day {
					t.Errorf("ParseCopyPaste() got Day = %v, want %v", got.Day, tt.want.Day)
				}
				if got.MaxAttempts != tt.want.MaxAttempts {
					t.Errorf("ParseCopyPaste() got MaxAttempts = %v, want %v", got.MaxAttempts, tt.want.MaxAttempts)
				}
				if got.Attempts != tt.want.Attempts {
					t.Errorf("ParseCopyPaste() got Attempts = %v, want %v", got.Attempts, tt.want.Attempts)
				}
				if got.Success != tt.want.Success {
					t.Errorf("ParseCopyPaste() got Success = %v, want %v", got.Success, tt.want.Success)
				}
				if len(got.AttemptsDetail) != len(tt.want.AttemptsDetail) {
					t.Errorf("ParseCopyPaste() got AttemptsDetail = %v, want %v", got.AttemptsDetail, tt.want.AttemptsDetail)
				}
				for i := range got.AttemptsDetail {
					if got.AttemptsDetail[i] != tt.want.AttemptsDetail[i] {
						t.Errorf("%v != %v", got.AttemptsDetail[i], tt.want.AttemptsDetail[i])
					}
				}
			}
		})
	}
}
