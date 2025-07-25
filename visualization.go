package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/vakrim/carcassonne-wave-collapse/tile"
)

const (
	screenWidth  = 800
	screenHeight = 600
	tileSize     = 40
	boardOffsetX = 50
	boardOffsetY = 50
)

var (
	// Colors for different border types
	fieldColor  = color.RGBA{34, 139, 34, 255}   // Forest green
	cityColor   = color.RGBA{139, 69, 19, 255}   // Brown
	streamColor = color.RGBA{30, 144, 255, 255}  // Dodger blue
	roadColor   = color.RGBA{128, 128, 128, 255} // Gray
	emptyColor  = color.RGBA{245, 245, 245, 255} // White smoke
)

type VisualizationGame struct {
	board         *Board
	pile          *Pile
	possibilities [][]PossibilitiesCount
	solver        *VisualizationSolver
}

func NewVisualizationGame(board *Board, pile *Pile) *VisualizationGame {
	return &VisualizationGame{
		board:         board,
		pile:          pile,
		possibilities: board.CountPossibilities(pile),
	}
}

func (g *VisualizationGame) UpdatePossibilities() {
	g.possibilities = g.board.CountPossibilities(g.pile)
}

func (g *VisualizationGame) SetSolver(solver *VisualizationSolver) {
	g.solver = solver
}

func (g *VisualizationGame) Update() error {
	// Check if space key is pressed to speed up visualization
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if g.solver != nil {
			g.solver.delay = 0
		}
	}
	return nil
}

func (g *VisualizationGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{240, 240, 240, 255}) // Light gray background

	// Draw the board
	g.drawBoard(screen)

	// Draw info text
	g.drawInfo(screen)
}

func (g *VisualizationGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *VisualizationGame) drawBoard(screen *ebiten.Image) {
	for row := range g.board.tiles {
		for col := range g.board.tiles[row] {
			x := boardOffsetX + col*tileSize
			y := boardOffsetY + row*tileSize

			if g.board.tiles[row][col] != nil {
				g.drawTile(screen, g.board.tiles[row][col], x, y)
			} else {
				g.drawEmptyTile(screen, x, y, row, col)
			}
		}
	}
}

func (g *VisualizationGame) drawTile(screen *ebiten.Image, t *tile.Tile, x, y int) {
	// Draw tile background
	ebitenutil.DrawRect(screen, float64(x), float64(y), tileSize, tileSize, color.RGBA{255, 255, 255, 255})

	// Draw borders
	borderSize := 8.0

	// Top border
	topColor := getBorderColor(t.Top())
	ebitenutil.DrawRect(screen, float64(x), float64(y), tileSize, borderSize, topColor)

	// Right border
	rightColor := getBorderColor(t.Right())
	ebitenutil.DrawRect(screen, float64(x+tileSize-int(borderSize)), float64(y), borderSize, tileSize, rightColor)

	// Bottom border
	bottomColor := getBorderColor(t.Bottom())
	ebitenutil.DrawRect(screen, float64(x), float64(y+tileSize-int(borderSize)), tileSize, borderSize, bottomColor)

	// Left border
	leftColor := getBorderColor(t.Left())
	ebitenutil.DrawRect(screen, float64(x), float64(y), borderSize, tileSize, leftColor)

	// Draw tile border
	ebitenutil.DrawRect(screen, float64(x), float64(y), tileSize, 1, color.Black)
	ebitenutil.DrawRect(screen, float64(x), float64(y), 1, tileSize, color.Black)
	ebitenutil.DrawRect(screen, float64(x+tileSize-1), float64(y), 1, tileSize, color.Black)
	ebitenutil.DrawRect(screen, float64(x), float64(y+tileSize-1), tileSize, 1, color.Black)
}

func (g *VisualizationGame) drawEmptyTile(screen *ebiten.Image, x, y, row, col int) {
	ebitenutil.DrawRect(screen, float64(x), float64(y), tileSize, tileSize, emptyColor)

	// Draw border
	ebitenutil.DrawRect(screen, float64(x), float64(y), tileSize, 1, color.RGBA{200, 200, 200, 255})
	ebitenutil.DrawRect(screen, float64(x), float64(y), 1, tileSize, color.RGBA{200, 200, 200, 255})
	ebitenutil.DrawRect(screen, float64(x+tileSize-1), float64(y), 1, tileSize, color.RGBA{200, 200, 200, 255})
	ebitenutil.DrawRect(screen, float64(x), float64(y+tileSize-1), tileSize, 1, color.RGBA{200, 200, 200, 255})

	// Draw possibilities count if available
	if g.possibilities != nil && row < len(g.possibilities) && col < len(g.possibilities[row]) {
		possCount := g.possibilities[row][col]
		if !possCount.alreadyPlaced && possCount.possibilities > 0 {
			// Show possibilities count in the center of the tile
			text := fmt.Sprintf("%d", possCount.possibilities)
			textX := x + tileSize/2 - 4 // Center text (rough approximation)
			textY := y + tileSize/2 - 4
			ebitenutil.DebugPrintAt(screen, text, textX, textY)
		}
	}
}

func (g *VisualizationGame) drawInfo(screen *ebiten.Image) {
	// Draw pile count and other info
	infoY := 10
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Tiles remaining: %d", len(*g.pile)), 10, infoY)
	
	// Show current delay
	delayText := "Normal speed"
	if g.solver != nil && g.solver.delay == 0 {
		delayText = "Fast mode (0ms delay)"
	} else if g.solver != nil {
		delayText = fmt.Sprintf("Delay: %dms", g.solver.delay.Milliseconds())
	}
	ebitenutil.DebugPrintAt(screen, delayText, 10, infoY+20)

	// Draw instructions
	ebitenutil.DebugPrintAt(screen, "Press SPACE to speed up", 10, infoY+40)
	ebitenutil.DebugPrintAt(screen, "Close window to exit", 10, infoY+60)
}

func getBorderColor(border string) color.Color {
	switch border {
	case "F":
		return fieldColor
	case "C":
		return cityColor
	case "S":
		return streamColor
	case "R":
		return roadColor
	default:
		return color.Black
	}
}
