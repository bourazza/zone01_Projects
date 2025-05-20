// main.go
package main

import (
	"fmt"
	"lem-in/src"
	"os"
)

func main() {
	// Check arguments
	if len(os.Args) != 2 {
		fmt.Println("ERROR: invalid data format, usage: go run . [filename]")
		return
	}
	
	filename := os.Args[1]
	
	// Read file content
	data, err := src.FileReader(filename)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	
	// Parse the data
	farm, err := src.Parse(data)
	if err != nil {
		fmt.Printf("ERROR: invalid data format, %v\n", err)
		return
	}
	
	// Print the formatted input
	  // Add an empty line as shown in the examples
	
	// Run the simulation
	moves, err := src.RunSimulation(farm)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	fmt.Print(src.FormatOutput(farm))
	fmt.Println()
	
	// Print the moves
	for _, move := range moves {
		fmt.Println(move)
	}
}