package main

import (
	"fmt"
	"os"
)

// This is a demo function to show what the visualization looks like in text form
func demonstrateVisualization() {
	fmt.Println("=== Carcassonne Wave Collapse Visualization Demo ===")
	fmt.Println()
	fmt.Println("The visualization shows a real-time graphical representation of the wave collapse algorithm.")
	fmt.Println()
	fmt.Println("Features:")
	fmt.Println("• Color-coded tile borders:")
	fmt.Println("  - Green: Field (F)")
	fmt.Println("  - Brown: City (C)")
	fmt.Println("  - Blue: Stream (S)")
	fmt.Println("  - Gray: Road (R)")
	fmt.Println()
	fmt.Println("• Real-time tile placement with visual delays")
	fmt.Println("• Backtracking visualization when algorithm needs to retry")
	fmt.Println("• Live tile count display")
	fmt.Println()
	fmt.Println("Example tile appearance:")
	fmt.Println("┌──────┐")
	fmt.Println("│ CITY │ <- Brown border (top)")
	fmt.Println("│      │")
	fmt.Println("│FIELD │ <- Green border (left)")
	fmt.Println("│      │")
	fmt.Println("│FIELD │ <- Green border (bottom)")
	fmt.Println("└──────┘")
	fmt.Println("   ^")
	fmt.Println("   Green border (right)")
	fmt.Println()
	fmt.Println("To see the actual visualization, run:")
	fmt.Println("  go build -tags visual")
	fmt.Println("  ./carcassonne-wave-collapse --visual")
	fmt.Println()
	fmt.Println("Note: Requires a display environment to run the graphical version.")
}

func init() {
	// If --demo flag is provided, show the demo
	if len(os.Args) > 1 && os.Args[1] == "--demo" {
		demonstrateVisualization()
		os.Exit(0)
	}
}