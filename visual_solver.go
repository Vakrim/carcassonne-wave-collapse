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
	solver := &VisualizationSolver{
		board:   board,
		pile:    pile,
		game:    game,
		delay:   time.Millisecond * 500, // delay between steps
		solving: false,
	}
	game.SetSolver(solver) // Set the solver reference for keyboard handling
	return solver
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
	// Update possibilities display
	vs.game.UpdatePossibilities()

	// Check if all tiles are used
	if len(*vs.pile) == 0 {
		return nil // Success - all tiles used
	}

	sortedPositions := getSortedAvailablePositions(vs.board, vs.pile)

	if len(sortedPositions) == 0 {
		return fmt.Errorf("no more valid positions to place remaining %d tiles", len(*vs.pile))
	}

	currentTile := vs.pile.PeekTop()

	for _, pos := range sortedPositions {
		pattern := vs.board.GetTilePattern(pos.row, pos.col)

		// Check if current tile matches this position
		if currentTile.MatchesQuery(pattern) {
			// Place the tile
			placedTile := vs.pile.PopTop()
			vs.board.tiles[pos.row][pos.col] = placedTile

			// Update possibilities display after placement
			vs.game.UpdatePossibilities()

			// Wait to show the placement
			vs.sleep()

			err := vs.solveWaveCollapseVisualized(recursiveCount + 1)
			if err == nil {
				return nil // solved!
			}

			vs.board.tiles[pos.row][pos.col] = nil
			vs.pile.PushTop(placedTile)

			// Update possibilities display after backtrack
			vs.game.UpdatePossibilities()

			// Wait to show the backtrack
			vs.sleepFaster()
		}
	}

	return fmt.Errorf("current tile %s cannot be placed in any available position", currentTile.String())
}

func (vs *VisualizationSolver) sleep() {
	if vs.delay <= 0 {
		return // No delay needed
	}

	time.Sleep(vs.delay)
}

func (vs *VisualizationSolver) sleepFaster() {
	vs.sleep()

	if vs.delay > 10*time.Millisecond {
		vs.delay -= 10 * time.Millisecond // Decrease delay to speed up visualization
	} else {
		vs.delay = 0 // Minimum delay to keep visualization responsive
	}
}
