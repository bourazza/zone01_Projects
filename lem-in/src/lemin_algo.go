// src/parser.go
package src

import (
	"fmt"
	"strconv"
	"strings"
)

// Room represents a room with a name and coordinates
type Room struct {
	Name string
	X    int
	Y    int
	IsStart bool
	IsEnd   bool
}

// Link represents a connection between two rooms
type Link struct {
	From string
	To   string
}

// AntFarm contains all data needed for the simulation
type AntFarm struct {
	NumAnts int
	Rooms   map[string]Room
	Links   []Link
	StartRoom string
	EndRoom   string
}

// ParseInt safely parses a string to int with error handling
func ParseInt(s string) (int, error) {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid integer: %s", s)
	}
	return val, nil
}

// ParseRoom parses a room line like "name x y"
func ParseRoom(line string) (*Room, error) {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid room format: %s", line)
	}
	
	name := parts[0]
	// Check if name starts with L or #
	if strings.HasPrefix(name, "L") || (strings.HasPrefix(name, "#") && name != "##start" && name != "##end") {
		return nil, fmt.Errorf("invalid room name: %s", name)
	}
	
	x, err := ParseInt(parts[1])
	if err != nil {
		return nil, err
	}
	
	y, err := ParseInt(parts[2])
	if err != nil {
		return nil, err
	}
	
	return &Room{
		Name: name,
		X:    x,
		Y:    y,
	}, nil
}

// ParseLink parses a link line like "room1-room2"
func ParseLink(line string) (*Link, error) {
	parts := strings.Split(line, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid link format: %s", line)
	}
	
	from := strings.TrimSpace(parts[0])
	to := strings.TrimSpace(parts[1])
	
	if from == "" || to == "" {
		return nil, fmt.Errorf("invalid link names: %s", line)
	}
	
	return &Link{
		From: from,
		To:   to,
	}, nil
}

// Parse processes the input data and builds an AntFarm structure
func Parse(data string) (*AntFarm, error) {
	lines := strings.Split(data, "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	
	farm := &AntFarm{
		Rooms: make(map[string]Room),
	}
	
	var i int
	var currentLine string
	var nextIsStart, nextIsEnd bool
	
	// Parse ant count from first non-comment line
	for i = 0; i < len(lines); i++ {
		currentLine = strings.TrimSpace(lines[i])
		if currentLine == "" || strings.HasPrefix(currentLine, "#") {
			// Skip empty lines and comments
			if currentLine == "##start" {
				nextIsStart = true
			} else if currentLine == "##end" {
				nextIsEnd = true
			}
			continue
		}
		
		// First non-comment line should be the number of ants
		numAnts, err := ParseInt(currentLine)
		if err != nil {
			return nil, fmt.Errorf("invalid number of ants: %s", currentLine)
		}
		if numAnts <= 0 {
			return nil, fmt.Errorf("invalid number of ants: %d", numAnts)
		}
		farm.NumAnts = numAnts
		i++
		break
	}
	
	// Parse rooms and links
	for ; i < len(lines); i++ {
		currentLine = strings.TrimSpace(lines[i])
		if currentLine == "" {
			continue
		}
		
		// Handle special commands
		if currentLine == "##start" {
			nextIsStart = true
			continue
		} else if currentLine == "##end" {
			nextIsEnd = true
			continue
		} else if strings.HasPrefix(currentLine, "#") {
			// Skip comments
			continue
		}
		
		// Check if it's a link (contains "-")
		if strings.Contains(currentLine, "-") {
			link, err := ParseLink(currentLine)
			if err != nil {
				return nil, err
			}
			
			// Verify that both rooms exist
			if _, exists := farm.Rooms[link.From]; !exists {
				return nil, fmt.Errorf("link references non-existent room: %s", link.From)
			}
			if _, exists := farm.Rooms[link.To]; !exists {
				return nil, fmt.Errorf("link references non-existent room: %s", link.To)
			}
			
			farm.Links = append(farm.Links, *link)
		} else {
			// It's a room
			room, err := ParseRoom(currentLine)
			if err != nil {
				return nil, err
			}
			
			// Check if room already exists
			if _, exists := farm.Rooms[room.Name]; exists {
				return nil, fmt.Errorf("duplicate room name: %s", room.Name)
			}
			
			// Set start/end flag if needed
			if nextIsStart {
				room.IsStart = true
				farm.StartRoom = room.Name
				nextIsStart = false
			} else if nextIsEnd {
				room.IsEnd = true
				farm.EndRoom = room.Name
				nextIsEnd = false
			}
			
			farm.Rooms[room.Name] = *room
		}
	}
	
	// Validate that we have start and end rooms
	if farm.StartRoom == "" {
		return nil, fmt.Errorf("no start room defined")
	}
	if farm.EndRoom == "" {
		return nil, fmt.Errorf("no end room defined")
	}
	
	return farm, nil
}
