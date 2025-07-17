package tile

import "math/rand"

type Border int

const (
	Field Border = iota
	City
	Stream
	Road
	BorderLength
)

func (b Border) String() string {
	switch b {
	case Field:
		return "F"
	case City:
		return "C"
	case Stream:
		return "S"
	case Road:
		return "R"
	default:
		panic("Unknown border type")
	}
}

type Tile struct {
	top    Border
	right  Border
	bottom Border
	left   Border
}

func (t *Tile) Top() string {
	return t.top.String()
}
func (t *Tile) Right() string {
	return t.right.String()
}
func (t *Tile) Bottom() string {
	return t.bottom.String()
}
func (t *Tile) Left() string {
	return t.left.String()
}

func (t *Tile) String() string {
	return t.Top() + t.Right() + t.Bottom() + t.Left()
}

func CreateRandomTile() Tile {
	return Tile{
		top:    getRandomBorder(),
		left:   getRandomBorder(),
		right:  getRandomBorder(),
		bottom: getRandomBorder(),
	}
}

func CreateTile(borders string) Tile {
	if len(borders) != 4 {
		panic("Invalid borders string length")
	}
	return Tile{
		top:    parseBorder(borders[0]),
		right:  parseBorder(borders[1]),
		bottom: parseBorder(borders[2]),
		left:   parseBorder(borders[3]),
	}
}

func parseBorder(b byte) Border {
	switch b {
	case 'F':
		return Field
	case 'C':
		return City
	case 'S':
		return Stream
	case 'R':
		return Road
	default:
		panic("Unknown border type")
	}
}

func getRandomBorder() Border {
	return Border(rand.Intn(int(BorderLength)))
}
