package main

import (
	"github.com/vakrim/carcassonne-wave-collapse/tile"
)

type Pile []tile.Tile

func (p *Pile) Size() int {
	return len(*p)
}

func (p *Pile) hasMoreTiles() bool {
	return len(*p) > 0
}

func (p *Pile) PeekTop() *tile.Tile {
	if !p.hasMoreTiles() {
		panic("No more tiles in the pile")
	}
	return &(*p)[0]
}

func (p *Pile) PopTop() *tile.Tile {
	if !p.hasMoreTiles() {
		panic("No more tiles in the pile")
	}
	tile := (*p)[0]
	*p = (*p)[1:]
	return &tile
}

func (p *Pile) FindMatchingTile(query string) *tile.Tile {
	for _, t := range *p {
		if doesTileMatchQuery(t.Top(), string(query[0])) &&
			doesTileMatchQuery(t.Right(), string(query[1])) &&
			doesTileMatchQuery(t.Bottom(), string(query[2])) &&
			doesTileMatchQuery(t.Left(), string(query[3])) {
			return &t
		}
	}
	return nil
}

func (p *Pile) Filter(query string) Pile {
	var result Pile
	for _, t := range *p {
		if doesTileMatchQuery(t.Top(), string(query[0])) &&
			doesTileMatchQuery(t.Right(), string(query[1])) &&
			doesTileMatchQuery(t.Bottom(), string(query[2])) &&
			doesTileMatchQuery(t.Left(), string(query[3])) {
			result = append(result, t)
		}
	}
	return result
}

func (p *Pile) CountMatchingTiles(query string) int {
	count := 0
	for _, t := range *p {
		if doesTileMatchQuery(t.Top(), string(query[0])) &&
			doesTileMatchQuery(t.Right(), string(query[1])) &&
			doesTileMatchQuery(t.Bottom(), string(query[2])) &&
			doesTileMatchQuery(t.Left(), string(query[3])) {
			count++
		}
	}
	return count
}

func (p *Pile) RemoveTile(tileToRemove *tile.Tile) {
	for i, t := range *p {
		if t == *tileToRemove {
			*p = append((*p)[:i], (*p)[i+1:]...)
			return
		}
	}
	panic("Tile not found in the pile")
}

func doesTileMatchQuery(tileBorder string, queryBorder string) bool {
	return queryBorder == "?" || tileBorder == queryBorder
}
