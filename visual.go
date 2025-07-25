//go:build visual

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/vakrim/carcassonne-wave-collapse/tile"
)

func runWithVisualization() {
	pile, err := loadTilesFromFile("tiles.txt")
	if err != nil {
		log.Fatalf("Error loading tiles: %v", err)
	}

	const boardSize = 12

	board := Board{
		tiles: make([][]*tile.Tile, boardSize),
	}

	for i := range board.tiles {
		board.tiles[i] = make([]*tile.Tile, boardSize)
	}

	board.tiles[6][6] = pile.PopTop()

	fmt.Printf("Loaded %d tiles from file\n", len(pile))
	fmt.Println("Starting visualization...")

	solver := NewVisualizationSolver(&board, &pile)

	// Start solving in background after a brief delay
	go func() {
		time.Sleep(time.Second * 2) // Wait 2 seconds before starting
		solver.StartSolving()
	}()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Carcassonne Wave Collapse Visualization")

	if err := ebiten.RunGame(solver); err != nil {
		log.Fatal(err)
	}
}