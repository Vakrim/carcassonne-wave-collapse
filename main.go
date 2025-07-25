package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/vakrim/carcassonne-wave-collapse/tile"
)

package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
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

	// Check if we can run the visualization
	if os.Getenv("DISPLAY") == "" {
		fmt.Println("No display detected, visualization requires a display environment")
		fmt.Println("Please run in an environment with a display server")
		os.Exit(1)
	}

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

type MinPossibilityPosition struct {
	row, col      int
	possibilities int
}

func findMinPossibilityPosition(board *Board, pile *Pile) []MinPossibilityPosition {
	possibilities := board.CountPossibilities(pile)
	minPossibilities := math.MaxInt
	var positions []MinPossibilityPosition

	// Find minimum number of possibilities
	for i := range possibilities {
		for j := range possibilities[i] {
			if !possibilities[i][j].alreadyPlaced &&
				possibilities[i][j].possibilities > 0 &&
				hasAdjacentTile(board, i, j) &&
				possibilities[i][j].possibilities < minPossibilities {
				minPossibilities = possibilities[i][j].possibilities
			}
		}
	}

	// Collect all positions with minimum possibilities
	for i := range possibilities {
		for j := range possibilities[i] {
			if !possibilities[i][j].alreadyPlaced &&
				possibilities[i][j].possibilities == minPossibilities &&
				hasAdjacentTile(board, i, j) {
				positions = append(positions, MinPossibilityPosition{
					row:           i,
					col:           j,
					possibilities: minPossibilities,
				})
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

// finds the tile from matchingTiles that has the fewest alternative placement options
func findBestTile(matchingTiles Pile, board *Board, pile *Pile) *tile.Tile {
	if len(matchingTiles) == 0 {
		return nil
	}

	bestTile := &matchingTiles[0]
	minAlternatives := math.MaxInt

	for i := range matchingTiles {
		currentTile := &matchingTiles[i]
		alternatives := countTilePlacementOptions(currentTile, board, pile)

		if alternatives < minAlternatives {
			minAlternatives = alternatives
			bestTile = currentTile
		}
	}

	return bestTile
}

// counts how many positions on the board this tile could be placed
func countTilePlacementOptions(targetTile *tile.Tile, board *Board, pile *Pile) int {
	count := 0
	possibilities := board.CountPossibilities(pile)

	for i := range possibilities {
		for j := range possibilities[i] {
			if !possibilities[i][j].alreadyPlaced && possibilities[i][j].possibilities > 0 &&
				hasAdjacentTile(board, i, j) {
				pattern := board.GetTilePattern(i, j)
				if targetTile.MatchesQuery(pattern) {
					count++
				}
			}
		}
	}

	return count
}
