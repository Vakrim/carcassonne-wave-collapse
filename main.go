package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/vakrim/carcassonne-wave-collapse/tile"
)

func main() {
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

func loadTilesFromFile(filename string) (Pile, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var pile Pile
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && len(line) == 4 {
			t := tile.CreateTile(line)
			pile = append(pile, t)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return pile, nil
}

type PositionWithPossibilities struct {
	row, col      int
	possibilities int
}

func getSortedAvailablePositions(board *Board, pile *Pile) []PositionWithPossibilities {
	possibilities := board.CountPossibilities(pile)
	var positions []PositionWithPossibilities

	// Collect all valid positions
	for i := range possibilities {
		for j := range possibilities[i] {
			if !possibilities[i][j].alreadyPlaced &&
				possibilities[i][j].possibilities > 0 &&
				hasAdjacentTile(board, i, j) {
				positions = append(positions, PositionWithPossibilities{
					row:           i,
					col:           j,
					possibilities: possibilities[i][j].possibilities,
				})
			}
		}
	}

	// Sort by possibilities (ascending - least possibilities first)
	for i := 0; i < len(positions)-1; i++ {
		for j := i + 1; j < len(positions); j++ {
			if positions[i].possibilities > positions[j].possibilities {
				positions[i], positions[j] = positions[j], positions[i]
			}
		}
	}

	return positions
}

func hasAdjacentTile(board *Board, row, col int) bool {
	directions := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} // up, down, left, right

	for _, dir := range directions {
		newRow, newCol := row+dir[0], col+dir[1]
		if newRow >= 0 && newRow < len(board.tiles) &&
			newCol >= 0 && newCol < len(board.tiles[0]) &&
			board.tiles[newRow][newCol] != nil {
			return true
		}
	}
	return false
}
