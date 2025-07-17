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
