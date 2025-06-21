package util_test

import (
	"testing"

	"github.com/YashBhalodi/chessnote/internal/util"
)

func TestIsFile(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		r    rune
		want bool
	}{
		{"lowercase a", 'a', true},
		{"lowercase h", 'h', true},
		{"lowercase g", 'g', true},
		{"outside range low", '`', false},
		{"outside range high", 'i', false},
		{"uppercase A", 'A', false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := util.IsFile(tt.r); got != tt.want {
				t.Errorf("IsFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsRank(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		r    rune
		want bool
	}{
		{"rank 1", '1', true},
		{"rank 8", '8', true},
		{"rank 5", '5', true},
		{"outside range low", '0', false},
		{"outside range high", '9', false},
		{"letter a", 'a', false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := util.IsRank(tt.r); got != tt.want {
				t.Errorf("IsRank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsWhitespace(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		r    rune
		want bool
	}{
		{"space", ' ', true},
		{"tab", '\t', true},
		{"newline", '\n', true},
		{"letter", 'a', false},
		{"digit", '1', false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := util.IsWhitespace(tt.r); got != tt.want {
				t.Errorf("IsWhitespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsLetter(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		r    rune
		want bool
	}{
		{"lowercase", 'a', true},
		{"uppercase", 'Z', true},
		{"digit", '5', false},
		{"symbol", '*', false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := util.IsLetter(tt.r); got != tt.want {
				t.Errorf("IsLetter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsDigit(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		r    rune
		want bool
	}{
		{"digit 0", '0', true},
		{"digit 9", '9', true},
		{"lowercase", 'a', false},
		{"uppercase", 'Z', false},
		{"symbol", '*', false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := util.IsDigit(tt.r); got != tt.want {
				t.Errorf("IsDigit() = %v, want %v", got, tt.want)
			}
		})
	}
}
