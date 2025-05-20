// src/simulation.go
package src

import (
	"fmt"
	"strings"
)

// Ant represents an ant in the simulation
type Ant struct {
	ID       int
	Position string   // Current room name
	Path     []string // Path this ant is following
	PathIdx  int      // Current position in the path
	Done     bool     // Whether the ant has reached the end
}

// Simulation manages the ant movement simulation
type Simulation struct {
	Farm  *AntFarm
	Ants  []*Ant
	Paths []*Path
}

// CalculatePathCosts calculates the cost of sending an ant through each path
// The cost is the path length + the number of ants already assigned to the path
func calculatePathCosts(paths []*Path, antCounts []int) []int {
	costs := make([]int, len(paths))
	for i, path := range paths {
		costs[i] = path.Length + antCounts[i]
	}
	return costs
}

// NewSimulation creates a new simulation with the given farm and paths
func NewSimulation(farm *AntFarm, paths []*Path) *Simulation {
	sim := &Simulation{
		Farm:  farm,
		Paths: paths,
		Ants:  make([]*Ant, farm.NumAnts),
	}

	// We'll keep track of how many ants are assigned to each path
	antCounts := make([]int, len(paths))

	// Initialize ants with optimal path distribution
	for i := 0; i < farm.NumAnts; i++ {
		ant := &Ant{
			ID:       i + 1, // Ant IDs start at 1
			Position: farm.StartRoom,
			PathIdx:  0,
			Done:     false,
		}

		// Choose the path with the lowest current cost
		if len(paths) > 0 {
			costs := calculatePathCosts(paths, antCounts)
			bestPathIdx := 0
			lowestCost := costs[0]

			for j := 1; j < len(costs); j++ {
				if costs[j] < lowestCost {
					lowestCost = costs[j]
					bestPathIdx = j
				}
			}

			ant.Path = paths[bestPathIdx].Rooms
			antCounts[bestPathIdx]++
		}

		sim.Ants[i] = ant
	}

	return sim
}

// RunSimulation runs the simulation until all ants reach the end room
func RunSimulation(farm *AntFarm) ([]string, error) {
	// Find paths from start to end
	paths, err := FindMultiplePaths(farm)
	if err != nil {
		return nil, err
	}

	if len(paths) == 0 {
		return nil, fmt.Errorf("no path from start to end")
	}

	// Create simulation
	sim := NewSimulation(farm, paths)

	// Run simulation until all ants reach the end
	moves := []string{}

	for {
		// Check if all ants have reached the end
		allDone := true
		for _, ant := range sim.Ants {
			if !ant.Done {
				allDone = false
				break
			}
		}

		if allDone {
			break
		}

		// Perform one turn of movement
		turnMoves := sim.SimulateTurn()
		if len(turnMoves) > 0 {
			moves = append(moves, strings.Join(turnMoves, " "))
		} else {
			// If no moves were made this turn but we're not done, something's wrong
			if !allDone {
				return nil, fmt.Errorf("simulation deadlock: no moves possible but not all ants are done")
			}
		}
	}

	return moves, nil
}

// SimulateTurn simulates one turn of ant movement and returns the moves made
func (sim *Simulation) SimulateTurn() []string {
	// This keeps track of which rooms will be occupied after this turn
	occupied := make(map[string]bool)

	// Start by marking rooms that have ants that won't move
	for _, ant := range sim.Ants {
		if !ant.Done && (ant.PathIdx >= len(ant.Path)-1 || ant.Position == sim.Farm.EndRoom) {
			occupied[ant.Position] = true
		}
	}

	// Track moves that will be made this turn
	moves := []string{}

	// For each ant not at the end yet
	for _, ant := range sim.Ants {
		// Skip ants that are done
		if ant.Done {
			continue
		}

		// If ant is at the end, mark it done
		if ant.Position == sim.Farm.EndRoom {
			ant.Done = true
			continue
		}

		// Skip ants that are at the last position in their path
		if ant.PathIdx >= len(ant.Path)-1 {
			continue
		}

		// Next room in the path
		nextRoom := ant.Path[ant.PathIdx+1]

		// Check if next room is available (not occupied by another ant)
		// Start and end rooms can hold multiple ants
		if nextRoom != sim.Farm.EndRoom && nextRoom != sim.Farm.StartRoom && occupied[nextRoom] {
			continue // Room is occupied, cannot move
		}

		// Move ant to next room
		ant.Position = nextRoom
		ant.PathIdx++

		// Mark next room as occupied
		if nextRoom != sim.Farm.EndRoom && nextRoom != sim.Farm.StartRoom {
			occupied[nextRoom] = true
		}

		// If ant reached the end, mark it done
		if nextRoom == sim.Farm.EndRoom {
			ant.Done = true
		}

		// Record the move
		move := fmt.Sprintf("L%d-%s", ant.ID, nextRoom)
		moves = append(moves, move)
	}

	return moves
}
