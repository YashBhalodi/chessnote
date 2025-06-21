package chessnote_test

import (
	"strings"
	"testing"

	"github.com/YashBhalodi/chessnote"
)

func TestSplitMultiGame(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		pgn     string
		want    []string
		wantLen int
	}{
		{
			name:    "no games",
			pgn:     "",
			want:    []string{},
			wantLen: 0,
		},
		{
			name: "single game",
			pgn:  `[Event "Test"]` + "\n" + `1. e4 e5 *`,
			want: []string{
				`[Event "Test"]` + "\n" + `1. e4 e5 *`,
			},
			wantLen: 1,
		},
		{
			name: "two games with unix newlines",
			pgn:  "[Event \"1\"]\n1. e4 *\n\n[Event \"2\"]\n1. d4 *",
			want: []string{
				"[Event \"1\"]\n1. e4 *",
				"[Event \"2\"]\n1. d4 *",
			},
			wantLen: 2,
		},
		{
			name: "two games with windows newlines",
			pgn:  "[Event \"1\"]\r\n1. e4 *\r\n\r\n[Event \"2\"]\r\n1. d4 *",
			want: []string{
				"[Event \"1\"]\n1. e4 *",
				"[Event \"2\"]\n1. d4 *",
			},
			wantLen: 2,
		},
		{
			name:    "empty with only whitespace",
			pgn:     "  \n \r\n  ",
			want:    []string{},
			wantLen: 0,
		},
		{
			name: "multiple games with extra spacing",
			pgn:  "\n\n[Event \"1\"]\n1. e4 e5 *\n\n\n[Event \"2\"]\n1. d4 d5 *\n\n",
			want: []string{
				"[Event \"1\"]\n1. e4 e5 *",
				"[Event \"2\"]\n1. d4 d5 *",
			},
			wantLen: 2,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := chessnote.SplitMultiGame(tt.pgn)
			if len(got) != tt.wantLen {
				t.Fatalf("SplitMultiGame() got %d games, want %d", len(got), tt.wantLen)
			}
			for i := range got {
				// We trim space here because the splitter can leave trailing/leading space
				// which is not significant.
				if strings.TrimSpace(got[i]) != strings.TrimSpace(tt.want[i]) {
					t.Errorf("game %d mismatch:\ngot:\n%s\nwant:\n%s", i, got[i], tt.want[i])
				}
			}
		})
	}
}
