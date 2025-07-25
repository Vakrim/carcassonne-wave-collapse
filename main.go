package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/vakrim/carcassonne-wave-collapse/tile"
)

func main() {
	// Check for command line flag to run with visualization
	if len(os.Args) > 1 && os.Args[1] == "--visual" {
		runWithVisualization()
		return
	}

	// Original console-based implementation
	runConsoleVersion()
}

func runConsoleVersion() {
	pile, err := loadTilesFromFile("tiles.txt")
	if err != nil {
		fmt.Printf("Error loading tiles: %v\n", err)
		return
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
	fmt.Println("Initial board:")
	fmt.Println(board.String())
	fmt.Println()

	if err := solveWaveCollapse(&board, &pile, 0); err == nil {
		fmt.Println("Success! All tiles have been placed:")
		fmt.Println(board.String())
		fmt.Printf("Tiles remaining in pile: %d\n", len(pile))
	} else {
		fmt.Println("Could not place all tiles")
		fmt.Println("Error:", err)
		fmt.Printf("Tiles remaining in pile: %d\n", len(pile))
		fmt.Println("Final board state:")
		fmt.Println(board.String())
	}
}

func solveWaveCollapse(board *Board, pile *Pile, recursiveCount int) error {
	// Check if all tiles are used
	if len(*pile) == 0 {
		return nil // Success - all tiles used
	}

	minPositions := findMinPossibilityPosition(board, pile)

	if len(minPositions) == 0 {
		return fmt.Errorf("no more valid positions to place remaining %d tiles", len(*pile))
	}

	// Try each position with minimum possibilities
	for _, minPos := range minPositions {
		if minPos.possibilities == 0 {
			continue
		}

		// get the tile pattern for the position with the least possibilities
		pattern := board.GetTilePattern(minPos.row, minPos.col)

		matchingTiles := pile.Filter(pattern)
		if len(matchingTiles) == 0 {
			continue
		}

		// Find the tile with the least placement options across the board
		bestTile := findBestTile(matchingTiles, board, pile)
		if bestTile == nil {
			continue
		}

		// place tile
		board.tiles[minPos.row][minPos.col] = bestTile
		pile.RemoveTile(bestTile)

		// recursively solve the rest of the board
		err := solveWaveCollapse(board, pile, recursiveCount+1)
		if err == nil {
			return nil // solved!
		}

		fmt.Printf("Backtracking from position (%d, %d) with tile: %s after %d recursions\n", minPos.row, minPos.col, bestTile.String(), recursiveCount)

		// remove tile from board
		board.tiles[minPos.row][minPos.col] = nil
		// Add tile back to pile
		*pile = append(*pile, *bestTile)
	}

	return fmt.Errorf("no solution found for any of the %d minimum possibility positions", len(minPositions))
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
