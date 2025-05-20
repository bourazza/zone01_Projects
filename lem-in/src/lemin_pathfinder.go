// src/pathfinder.go
package src

import (
	"container/heap"
	"sort"
)

// Path represents a path from start to end
type Path struct {
	Rooms []string
	Length int
}

// PriorityQueue for Dijkstra's algorithm
type PriorityQueue []*PQItem

type PQItem struct {
	RoomName  string
	Distance  int
	Index     int
}

// Required methods for heap.Interface
func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Distance < pq[j].Distance
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*PQItem)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// Find paths from start to end
func FindShortestPath(farm *AntFarm, excludedRooms map[string]bool) (*Path, error) {
	// Create adjacency list from links
	graph := make(map[string][]string)
	for _, link := range farm.Links {
		// Skip links that contain excluded rooms
		if excludedRooms[link.From] || excludedRooms[link.To] {
			continue
		}
		
		// Links are bidirectional
		graph[link.From] = append(graph[link.From], link.To)
		graph[link.To] = append(graph[link.To], link.From)
	}
	
	// Initialize distances and previous map
	dist := make(map[string]int)
	prev := make(map[string]string)
	
	// Initialize all distances to infinity
	for roomName := range farm.Rooms {
		dist[roomName] = int(^uint(0) >> 1) // Max int value
	}
	
	// Distance to start is 0
	dist[farm.StartRoom] = 0
	
	// Priority queue for Dijkstra's algorithm
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	
	// Push start room with distance 0
	heap.Push(&pq, &PQItem{
		RoomName: farm.StartRoom,
		Distance: 0,
	})
	
	// Dijkstra's algorithm
	for pq.Len() > 0 {
		// Get room with minimum distance
		current := heap.Pop(&pq).(*PQItem)
		
		// If we reached the end room, we're done
		if current.RoomName == farm.EndRoom {
			break
		}
		
		// Skip if we've already found a shorter path
		if current.Distance > dist[current.RoomName] {
			continue
		}
		
		// Explore neighbors
		for _, neighbor := range graph[current.RoomName] {
			// Skip excluded rooms
			if excludedRooms[neighbor] {
				continue
			}
			
			// Distance to neighbor through current room
			distance := dist[current.RoomName] + 1
			
			// If we found a shorter path to neighbor
			if distance < dist[neighbor] {
				dist[neighbor] = distance
				prev[neighbor] = current.RoomName
				
				heap.Push(&pq, &PQItem{
					RoomName: neighbor,
					Distance: distance,
				})
			}
		}
	}
	
	// Check if end room is reachable
	if _, exists := prev[farm.EndRoom]; !exists && farm.EndRoom != farm.StartRoom {
		return nil, nil // No path exists
	}
	
	// Reconstruct path from end to start
	path := &Path{}
	current := farm.EndRoom
	
	// Add rooms to path in reverse order
	for current != farm.StartRoom {
		path.Rooms = append([]string{current}, path.Rooms...)
		current = prev[current]
	}
	
	// Add start room
	path.Rooms = append([]string{farm.StartRoom}, path.Rooms...)
	path.Length = len(path.Rooms) - 1 // Length is number of edges, which is rooms - 1
	
	return path, nil
}

// FindMultiplePaths finds multiple non-overlapping paths
func FindMultiplePaths(farm *AntFarm) ([]*Path, error) {
	// Start with no excluded rooms except start and end (they can be shared)
	excludedRooms := make(map[string]bool)
	
	// We'll store all found paths here
	var allPaths []*Path
	
	// Keep finding paths until we can't find any more
	for {
		path, err := FindShortestPath(farm, excludedRooms)
		if err != nil {
			return nil, err
		}
		
		// If no path exists, we're done
		if path == nil {
			break
		}
		
		// Add path to results
		allPaths = append(allPaths, path)
		
		// Mark intermediate rooms as excluded for future paths
		// (don't exclude start and end)
		for i := 1; i < len(path.Rooms)-1; i++ {
			excludedRooms[path.Rooms[i]] = true
		}
	}
	
	// Sort paths by length (shortest first)
	sort.Slice(allPaths, func(i, j int) bool {
		return allPaths[i].Length < allPaths[j].Length
	})
	
	return allPaths, nil
}