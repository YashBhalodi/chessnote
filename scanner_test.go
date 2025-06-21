package chessnote

import (
	"strings"
	"testing"
)

func TestScanner(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name  string
		input string
		want  []Token
	}{
		{
			name:  "tags and result",
			input: `[Event "Test"] *`,
			want: []Token{
				{Type: LBRACKET, Literal: "["},
				{Type: IDENT, Literal: "Event"},
				{Type: STRING, Literal: "Test"},
				{Type: RBRACKET, Literal: "]"},
				{Type: ASTERISK, Literal: "*"},
				{Type: EOF},
			},
		},
		{
			name:  "simple move",
			input: `1. e4`,
			want: []Token{
				{Type: NUMBER, Literal: "1"},
				{Type: DOT, Literal: "."},
				{Type: IDENT, Literal: "e4"},
				{Type: EOF},
			},
		},
		{
			name:  "disambiguated move",
			input: `Rdf8`,
			want: []Token{
				{Type: IDENT, Literal: "Rdf8"},
				{Type: EOF},
			},
		},
		{
			name:  "pawn promotion",
			input: `e8=Q`,
			want: []Token{
				{Type: IDENT, Literal: "e8=Q"},
				{Type: EOF},
			},
		},
		{
			name:  "pawn promotion with check",
			input: `e8=Q+`,
			want: []Token{
				{Type: IDENT, Literal: "e8=Q+"},
				{Type: EOF},
			},
		},
		{
			name:  "capture with disambiguation and promotion",
			input: `exd8=R#`,
			want: []Token{
				{Type: IDENT, Literal: "exd8=R#"},
				{Type: EOF},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewScanner(strings.NewReader(tc.input))
			for i, wantToken := range tc.want {
				gotToken := s.Scan()
				if gotToken.Type != wantToken.Type {
					t.Fatalf("test %d: token type wrong. got=%v, want=%v", i, gotToken.Type, wantToken.Type)
				}
				if gotToken.Literal != wantToken.Literal {
					t.Fatalf("test %d: token literal wrong. got=%q, want=%q", i, gotToken.Literal, wantToken.Literal)
				}
			}
		})
	}
}
