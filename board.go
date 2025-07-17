package main

import (
	"strings"

	"github.com/vakrim/carcassonne-wave-collapse/tile"
)

type Board struct {
	tiles [][]*tile.Tile
}

type PossibilitiesCount struct {
	possibilities int
	alreadyPlaced bool
}

func (b *Board) CountPossibilities(pile *Pile) [][]PossibilitiesCount {
	possibilities := make([][]PossibilitiesCount, len(b.tiles))
	for i := range b.tiles {
		possibilities[i] = make([]PossibilitiesCount, len(b.tiles[i]))
		for j := range b.tiles[i] {
			if b.tiles[i][j] == nil {
				possibilities[i][j] = PossibilitiesCount{
					possibilities: pile.CountMatchingTiles(b.GetTilePattern(i, j)),
					alreadyPlaced: false,
				}
			} else {
				possibilities[i][j] = PossibilitiesCount{
					possibilities: 0,
					alreadyPlaced: true,
				}
			}
		}
	}
	return possibilities
}

func (b *Board) GetTilePattern(row, col int) string {
	pattern := ""
	if row > 0 && b.tiles[row-1][col] != nil {
		pattern += b.tiles[row-1][col].Bottom()
	} else {
		pattern += "?"
	}
	if col < len(b.tiles[row])-1 && b.tiles[row][col+1] != nil {
		pattern += b.tiles[row][col+1].Left()
	} else {
		pattern += "?"
	}
	if row < len(b.tiles)-1 && b.tiles[row+1][col] != nil {
		pattern += b.tiles[row+1][col].Top()
	} else {
		pattern += "?"
	}
	if col > 0 && b.tiles[row][col-1] != nil {
		pattern += b.tiles[row][col-1].Right()
	} else {
		pattern += "?"
	}
	return pattern
}

func (b *Board) String() string {
	var sb strings.Builder
	for rowNumber, row := range b.tiles {
		for _, t := range row {
			sb.WriteString("[")
			if t != nil {
				sb.WriteString(t.String())
			} else {
				sb.WriteString("    ")
			}
			sb.WriteString("]")
		}
		if rowNumber < len(b.tiles)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

const tileStringRepresentationLength = 6

func BoardFromString(s string) Board {
	lines := strings.Split(s, "\n")
	board := Board{
		tiles: make([][]*tile.Tile, len(lines)),
	}
	for i, line := range lines {
		board.tiles[i] = make([]*tile.Tile, len(line)/tileStringRepresentationLength)
		for j := 0; j < len(line); j += tileStringRepresentationLength {
			char := line[j : j+tileStringRepresentationLength]
			if char != "[    ]" {
				t := tile.CreateTile(string(char[1:5]))
				board.tiles[i][j/tileStringRepresentationLength] = &t
			}
		}
	}
	return board
}
