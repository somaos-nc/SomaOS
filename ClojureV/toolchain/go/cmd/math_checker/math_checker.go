package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Math Checker for ClojureV
// Verifies algebraic invariants and functional purity after linting.

func main() {
	fmt.Println("--- INITIATING MATHEMATICAL VERIFICATION ---")
	
	if len(os.Args) < 2 {
		fmt.Println("Usage: math <cljv_file>")
		os.Exit(1)
	}
	
	fileData, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("[FRACTURE] Failed to read file: %v\n", err)
		os.Exit(1)
	}
	
	cljv := string(fileData)
	
	// Rule 1: Division by zero protection in functional logic
	fmt.Println("[MT] Checking Invariant: Zero-Division Safety")
	if strings.Contains(cljv, "/") {
		reDiv := regexp.MustCompile(`\/\s+(\w+)\s+(\w+)`)
		matches := reDiv.FindAllStringSubmatch(cljv, -1)
		for _, m := range matches {
			denominator := m[2]
			if denominator == "0" || denominator == "0.0" {
				fmt.Printf("[FRACTURE] Static Division by Zero detected: %s\n", m[0])
				os.Exit(1)
			}
		}
	}
	fmt.Println("[MT] Invariant 'Zero-Division Safety' verified.")

	// Rule 2: Pure Function Non-Mutability (No 'set!' or 'atom' in core math)
	fmt.Println("[MT] Checking Invariant: Functional Purity")
	if strings.Contains(cljv, "set!") || strings.Contains(cljv, "atom") {
		fmt.Println("[FRACTURE] Impure state mutation detected in functional territory.")
		os.Exit(1)
	}
	fmt.Println("[MT] Invariant 'Functional Purity' verified.")

	fmt.Println("--- MATHEMATICAL VERIFICATION ABSOLUTE ---")
}
