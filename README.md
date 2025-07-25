# Carcassonne Wave Collapse

A Go implementation of the Wave Function Collapse algorithm applied to Carcassonne-style tile placement with real-time graphical visualization.

## Overview

This project implements a tile-based puzzle solver inspired by the board game Carcassonne. It uses concepts from the Wave Function Collapse algorithm to generate valid tile placements on a board where tiles must match their neighbors according to specific rules. The application provides real-time visualization of the solving process using ebitengine.

## Features

- **Tile System**: Tiles with four borders (top, right, bottom, left) that can be:

  - `F` - Field
  - `C` - City
  - `S` - Stream
  - `R` - Road

- **Board Management**: 2D grid system for placing tiles with constraint checking

- **Pile System**: Manages available tiles for placement

- **Pattern Matching**: Validates tile placement based on neighboring tiles

- **Possibility Counting**: Calculates how many tiles from the pile can fit in each empty position

- **Real-time Visualization**: Graphical display showing the wave collapse algorithm in action with color-coded tile borders and backtracking visualization

## Installation

1. Clone the repository:

```bash
git clone https://github.com/vakrim/carcassonne-wave-collapse.git
cd carcassonne-wave-collapse
```

2. Build the application:

```bash
go build
```

## Usage

Run the visualization:

```bash
./carcassonne-wave-collapse
```

Or run directly:

```bash
go run .
```

**Note**: Requires a display environment to run. The application uses ebitengine for graphics and needs a display server (X11 on Linux, etc.).

### Running Tests

```bash
go test ./...
```

## Visualization

The visualization shows:
- **Real-time tile placement**: Watch as the wave collapse algorithm places tiles one by one
- **Color-coded borders**: Different colors for each border type (Field=Green, City=Brown, Stream=Blue, Road=Gray)
- **Backtracking**: Visual feedback when the algorithm needs to backtrack and try different placements
- **Tile count**: Shows remaining tiles in the pile

### Controls
- Close the window to exit the visualization

### Example Tile Patterns

Tiles are represented as 4-character strings representing `[Top][Right][Bottom][Left]` borders:

- `"FFFF"` - All field borders
- `"CCFF"` - City on top and right, field on bottom and left
- `"RCRC"` - Road on top and bottom, city on left and right

## API Reference

### Tile Package

- `CreateRandomTile()` - Creates a tile with random borders
- `CreateTile(borders string)` - Creates a tile from a 4-character pattern
- `tile.String()` - Returns the tile's border pattern

### Board

- `GetTilePattern(row, col int)` - Gets the required pattern for a position
- `CountPossibilities(pile *Pile)` - Counts valid tiles for each empty position
- `BoardFromString(s string)` - Creates a board from string representation

### Pile

- `PopTop()` - Removes and returns the top tile
- `PeekTop()` - Returns the top tile without removing it
- `CountMatchingTiles(pattern string)` - Counts tiles that match a pattern

## License

This project is open source and available under the [MIT License](LICENSE).

## Acknowledgments

- Inspired by the Wave Function Collapse algorithm
- Based on the tile-matching mechanics of Carcassonne board game
- Visualization powered by [Ebitengine](https://ebitengine.org/)
