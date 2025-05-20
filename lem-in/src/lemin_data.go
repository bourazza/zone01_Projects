// src/filereader.go
package src

import (
	"fmt"
	"os"
	"sort"
)

// FileReader reads the content of a file
func FileReader(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}
	
	return string(data), nil
}

// FormatOutput formats the ant farm data for display
func FormatOutput(farm *AntFarm) string {
	output := fmt.Sprintf("%d\n", farm.NumAnts)
	output += "#rooms\n"
	
	// Print start room first
	output += "##start\n"
	startRoom := farm.Rooms[farm.StartRoom]
	output += fmt.Sprintf("%s %d %d\n", startRoom.Name, startRoom.X, startRoom.Y)
	
	// Print end room second
	output += "##end\n"
	endRoom := farm.Rooms[farm.EndRoom]
	output += fmt.Sprintf("%s %d %d\n", endRoom.Name, endRoom.X, endRoom.Y)
	
	// Print other rooms in alphabetical order
	var roomNames []string
	for name, room := range farm.Rooms {
		if !room.IsStart && !room.IsEnd {
			roomNames = append(roomNames, name)
		}
	}
	
	// Sort room names alphabetically
	sort.Strings(roomNames)
	
	// Print rooms
	for _, name := range roomNames {
		room := farm.Rooms[name]
		output += fmt.Sprintf("%s %d %d\n", room.Name, room.X, room.Y)
	}
	
	// Print links
	for _, link := range farm.Links {
		output += fmt.Sprintf("%s-%s\n", link.From, link.To)
	}
	
	return output
}