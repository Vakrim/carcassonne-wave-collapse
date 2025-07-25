//go:build visual

package main

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// VisualizationSolver wraps the solving logic with visual updates
type VisualizationSolver struct {
	board   *Board
	pile    *Pile
	game    *VisualizationGame
	delay   time.Duration
	solving bool
}

func NewVisualizationSolver(board *Board, pile *Pile) *VisualizationSolver {
	game := NewVisualizationGame(board, pile)
	return &VisualizationSolver{
		board:   board,
		pile:    pile,
		game:    game,
		delay:   time.Millisecond * 500, // 500ms delay between steps
		solving: false,
	}
}

func (vs *VisualizationSolver) StartSolving() {
	if vs.solving {
		return
	}
	vs.solving = true
	go func() {
		defer func() {
			vs.solving = false
		}()
		
		fmt.Println("Starting visualization solve...")
		err := vs.solveWaveCollapseVisualized(0)
		if err == nil {
			fmt.Println("Success! All tiles have been placed.")
		} else {
			fmt.Printf("Could not place all tiles: %v\n", err)
		}
	}()
}

func (vs *VisualizationSolver) Update() error {
	return vs.game.Update()
}

func (vs *VisualizationSolver) Draw(screen *ebiten.Image) {
	vs.game.Draw(screen)
}

func (vs *VisualizationSolver) Layout(outsideWidth, outsideHeight int) (int, int) {
	return vs.game.Layout(outsideWidth, outsideHeight)
}

func (vs *VisualizationSolver) solveWaveCollapseVisualized(recursiveCount int) error {
	// Check if all tiles are used
	if len(*vs.pile) == 0 {
		return nil // Success - all tiles used
	}

	minPositions := findMinPossibilityPosition(vs.board, vs.pile)

	if len(minPositions) == 0 {
		return fmt.Errorf("no more valid positions to place remaining %d tiles", len(*vs.pile))
	}

	// Try each position with minimum possibilities
	for _, minPos := range minPositions {
		if minPos.possibilities == 0 {
			continue
		}

		// get the tile pattern for the position with the least possibilities
		pattern := vs.board.GetTilePattern(minPos.row, minPos.col)

		matchingTiles := vs.pile.Filter(pattern)
		if len(matchingTiles) == 0 {
			continue
		}

		// Find the tile with the least placement options across the board
		bestTile := findBestTile(matchingTiles, vs.board, vs.pile)
		if bestTile == nil {
			continue
		}

		// place tile
		vs.board.tiles[minPos.row][minPos.col] = bestTile
		vs.pile.RemoveTile(bestTile)

		// Wait to show the placement
		time.Sleep(vs.delay)

		// recursively solve the rest of the board
		err := vs.solveWaveCollapseVisualized(recursiveCount + 1)
		if err == nil {
			return nil // solved!
		}

		fmt.Printf("Backtracking from position (%d, %d) with tile: %s after %d recursions\n", minPos.row, minPos.col, bestTile.String(), recursiveCount)

		// remove tile from board
		vs.board.tiles[minPos.row][minPos.col] = nil
		// Add tile back to pile
		*vs.pile = append(*vs.pile, *bestTile)

		// Wait to show the backtrack
		time.Sleep(vs.delay)
	}

	return fmt.Errorf("no solution found for any of the %d minimum possibility positions", len(minPositions))
}