package main

import (
	"testing"

	"github.com/vakrim/carcassonne-wave-collapse/tile"
)

func TestFindMatchingTile(t *testing.T) {
	pile := Pile{
		tile.CreateTile("FFFF"),
		tile.CreateTile("CCFF"),
	}

	query := "CCFF"
	if pile.FindMatchingTile(query).String() != "CCFF" {
		t.Errorf("Expected to find a tile matching %s, but found none or incorrect tile", query)
	}

	query = "C???"
	if pile.FindMatchingTile(query).String() != "CCFF" {
		t.Errorf("Expected to find a tile matching %s, but found none or incorrect tile", query)
	}

	query = "C??C"
	if pile.FindMatchingTile(query) != nil {
		t.Errorf("Expected to find no tile matching %s, but found one", query)
	}
}

func TestRemoveTile(t *testing.T) {
	pile := Pile{
		tile.CreateTile("FFFF"),
		tile.CreateTile("CCFF"),
		tile.CreateTile("FFCC"),
	}

	tileToRemove := &pile[1]
	pile.RemoveTile(tileToRemove)
	if pile.Size() != 2 {
		t.Errorf("Expected pile size 2 after removal, got %d", pile.Size())
	}
	for _, tile := range pile {
		if tile.String() == "CCFF" {
			t.Errorf("Tile 'CCFF' should have been removed from the pile")
		}
	}

	tileToRemove = &pile[0]
	pile.RemoveTile(tileToRemove)
	if pile.Size() != 1 {
		t.Errorf("Expected pile size 1 after removal, got %d", pile.Size())
	}
	if pile[0].String() != "FFCC" {
		t.Errorf("Expected remaining tile to be 'FFCC', got '%s'", pile[0].String())
	}

	tileToRemove = &pile[0]
	pile.RemoveTile(tileToRemove)
	if pile.Size() != 0 {
		t.Errorf("Expected pile size 0 after removal, got %d", pile.Size())
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when removing a tile not in the pile, but did not panic")
		}
	}()
	pile.RemoveTile(tileToRemove)
}
