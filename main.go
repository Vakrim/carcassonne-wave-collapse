package main

import (
	"github.com/vakrim/carcassonne-wave-collapse/tile"
)

func main() {
	pile := Pile{}

	for range 50 {
		pile = append(pile, tile.CreateRandomTile())
	}

	const boardSize = 7

	board := Board{
		tiles: make([][]*tile.Tile, boardSize),
	}

	for i := range board.tiles {
		board.tiles[i] = make([]*tile.Tile, boardSize)
	}

	board.tiles[3][3] = pile.PopTop()

	print(board.String())
}
