package ast

import (
	"testing"

	"github.com/adamstimb/rmbasicx64yar/internal/app/rmbasicx64yar/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{
					TokenType: token.LET,
					Literal:   "LET",
				},
				Name: &Identifier{
					Token: token.Token{
						TokenType: token.IdentifierLiteral,
						Literal:   "Myvar"},
					Value: "Myvar",
				},
				Value: &Identifier{
					Token: token.Token{
						TokenType: token.IdentifierLiteral,
						Literal:   "Anothervar",
					},
					Value: "Anothervar",
				},
			},
		},
	}
	if program.String() != "LET Myvar := Anothervar;" {
		t.Errorf("program.String() wrong, got %q", program.String())
	}
}
