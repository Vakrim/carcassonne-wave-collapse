//go:build !visual

package main

import "fmt"

func runWithVisualization() {
	fmt.Println("Visual mode not available. Rebuild with: go build -tags visual")
	fmt.Println("Or run without the --visual flag for console mode.")
}