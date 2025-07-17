package main

import (
	"fmt"
	"math"

	"github.com/vakrim/carcassonne-wave-collapse/tile"
)

func main() {
	pile := Pile{}

	for range 100 {
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

	fmt.Println("Initial board:")
	fmt.Println(board.String())
	fmt.Println()

	if err := solveWaveCollapse(&board, &pile); err == nil {
		fmt.Println("Final board:")
		fmt.Println(board.String())
	} else {
		fmt.Println("Could not solve the puzzle with available tiles")
		fmt.Println("Error:", err)
		fmt.Println("Final board state:")
		fmt.Println(board.String())
	}
}

func solveWaveCollapse(board *Board, pile *Pile) error {
	for {
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

		pattern := board.GetTilePattern(minPos.row, minPos.col)

		matchingTile := pile.FindMatchingTile(pattern)
		if matchingTile == nil {
			return fmt.Errorf("no matching tile found for position (%d, %d)", minPos.row, minPos.col)
		}

		board.tiles[minPos.row][minPos.col] = matchingTile

		pile.RemoveTile(matchingTile)
	}
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
