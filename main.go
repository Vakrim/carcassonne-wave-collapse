package main

import (
	"fmt"
	"math"

	"github.com/vakrim/carcassonne-wave-collapse/tile"
)

func main() {
	pile := Pile{}

	const boardSize = 7

	const spareTiles = 10

	for range boardSize*boardSize + spareTiles {
		pile = append(pile, tile.CreateRandomTile())
	}

	board := Board{
		tiles: make([][]*tile.Tile, boardSize),
	}

	for i := range board.tiles {
		board.tiles[i] = make([]*tile.Tile, boardSize)
	}

	board.tiles[3][3] = pile.PopTop()

	fmt.Println("Initial board:")
	fmt.Println(board.String())
	fmt.Println()

	if err := solveWaveCollapse(&board, &pile, 0); err == nil {
		fmt.Println("Final board:")
		fmt.Println(board.String())
	} else {
		fmt.Println("Could not solve the puzzle with available tiles")
		fmt.Println("Error:", err)
		fmt.Println("Final board state:")
		fmt.Println(board.String())
	}
}

func solveWaveCollapse(board *Board, pile *Pile, recursiveCount int) error {
	minPos := findMinPossibilityPosition(board, pile)

	if minPos.row == -1 {
		if isBoardComplete(board) {
			return nil
		} else {
			return fmt.Errorf("no more positions to fill, but board is not complete")
		}
	}

	if minPos.possibilities == 0 {
		return fmt.Errorf("no possibilities left for position (%d, %d)", minPos.row, minPos.col)
	}

	// get the tile pattern for the position with the least possibilities
	pattern := board.GetTilePattern(minPos.row, minPos.col)

	matchingTiles := pile.Filter(pattern)
	if len(matchingTiles) == 0 {
		return fmt.Errorf("no matching tile found for position (%d, %d)", minPos.row, minPos.col)
	}

	for i := range matchingTiles {
		t := &matchingTiles[i]

		// place tile
		board.tiles[minPos.row][minPos.col] = t
		pile.RemoveTile(t)

		// recursively solve the rest of the board
		err := solveWaveCollapse(board, pile, recursiveCount+1)
		if err == nil {
			return nil // solved!
		}

		fmt.Printf("Backtracking from position (%d, %d) with tile: %s after %d recursions\n", minPos.row, minPos.col, t.String(), recursiveCount)

		// remove tile from board
		board.tiles[minPos.row][minPos.col] = nil
		// Add tile back to pile
		*pile = append(*pile, *t)
	}

	return fmt.Errorf("no solution found for position (%d, %d)", minPos.row, minPos.col)
}

type MinPossibilityPosition struct {
	row, col      int
	possibilities int
}

func findMinPossibilityPosition(board *Board, pile *Pile) MinPossibilityPosition {
	possibilities := board.CountPossibilities(pile)
	minPossibilities := math.MaxInt
	minRow, minCol := -1, -1

	for i := range possibilities {
		for j := range possibilities[i] {
			if !possibilities[i][j].alreadyPlaced &&
				possibilities[i][j].possibilities < minPossibilities {
				minPossibilities = possibilities[i][j].possibilities
				minRow, minCol = i, j
			}
		}
	}

	return MinPossibilityPosition{row: minRow, col: minCol, possibilities: minPossibilities}
}

func isBoardComplete(board *Board) bool {
	for i := range board.tiles {
		for j := range board.tiles[i] {
			if board.tiles[i][j] == nil {
				return false
			}
		}
	}
	return true
}
