package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/vakrim/carcassonne-wave-collapse/tile"
)

const (
	screenWidth   = 800
	screenHeight  = 600
	tileSize      = 40
	boardOffsetX  = 50
	boardOffsetY  = 50
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
	board *Board
	pile  *Pile
}

func NewVisualizationGame(board *Board, pile *Pile) *VisualizationGame {
	return &VisualizationGame{
		board: board,
		pile:  pile,
	}
}

func (g *VisualizationGame) Update() error {
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
				g.drawEmptyTile(screen, x, y)
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

func (g *VisualizationGame) drawEmptyTile(screen *ebiten.Image, x, y int) {
	ebitenutil.DrawRect(screen, float64(x), float64(y), tileSize, tileSize, emptyColor)
	// Draw border
	ebitenutil.DrawRect(screen, float64(x), float64(y), tileSize, 1, color.RGBA{200, 200, 200, 255})
	ebitenutil.DrawRect(screen, float64(x), float64(y), 1, tileSize, color.RGBA{200, 200, 200, 255})
	ebitenutil.DrawRect(screen, float64(x+tileSize-1), float64(y), 1, tileSize, color.RGBA{200, 200, 200, 255})
	ebitenutil.DrawRect(screen, float64(x), float64(y+tileSize-1), tileSize, 1, color.RGBA{200, 200, 200, 255})
}

func (g *VisualizationGame) drawInfo(screen *ebiten.Image) {
	// Draw pile count and other info
	infoY := 10
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Tiles remaining: %d", len(*g.pile)), 10, infoY)
	
	// Draw instructions
	ebitenutil.DebugPrintAt(screen, "Close window to exit", 10, infoY+20)
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