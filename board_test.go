package main

import (
	"reflect"
	"testing"

	"github.com/vakrim/carcassonne-wave-collapse/tile"
)

func TestBoardFromString(t *testing.T) {
	input := `[    ][    ][    ]
[    ][RCCC][    ]
[    ][    ][    ]`

	board := BoardFromString(input)

	if board.String() != input {
		t.Errorf("Expected:\n%s\nGot:\n%s", input, board.String())
	}
}

func TestGetTilePattern(t *testing.T) {
	board := BoardFromString(`[    ][    ][    ]
[    ][RCCC][    ]
[    ][    ][CCCC]`)

	expected := [][]string{
		{"????", "??R?", "????"},
		{"?C??", "????", "??CC"},
		{"????", "CC??", "????"},
	}

	result := [][]string{}

	for i := 0; i < len(board.tiles); i++ {
		for j := 0; j < len(board.tiles[i]); j++ {
			pattern := board.GetTilePattern(i, j)
			if len(result) <= i {
				result = append(result, []string{})
			}
			result[i] = append(result[i], pattern)
		}
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, result)
	}
}

func TestCountPossibilities(t *testing.T) {
	input := `[    ][    ][    ]
[    ][RCCC][    ]
[    ][    ][CCCC]`

	t.Run("Empty Tiles", func(t *testing.T) {

		board := BoardFromString(input)

		possibilities := board.CountPossibilities(&Pile{})
		expected := [][]PossibilitiesCount{
			{{0, false}, {0, false}, {0, false}},
			{{0, false}, {0, true}, {0, false}},
			{{0, false}, {0, false}, {0, true}},
		}

		if !reflect.DeepEqual(possibilities, expected) {
			t.Errorf("Expected:\n%v\nGot:\n%v", expected, possibilities)
		}
	})

	t.Run("With Tiles", func(t *testing.T) {
		pile := Pile{
			tile.CreateTile("FFFF"),
			tile.CreateTile("CCFF"),
			tile.CreateTile("RCRC"),
		}

		board := BoardFromString(input)

		possibilities := board.CountPossibilities(&pile)
		expected := [][]PossibilitiesCount{
			{{3, false}, {1, false}, {3, false}},
			{{2, false}, {0, true}, {0, false}},
			{{3, false}, {1, false}, {0, true}},
		}

		if !reflect.DeepEqual(possibilities, expected) {
			t.Errorf("Expected:\n%v\nGot:\n%v", expected, possibilities)
		}
	})
}
