package game

import (
	"strings"
)

type Location struct {
	LocationName      string   `json:"location_name"`
	EnemiesInLocation []string `json:"enemies_in_location"`
	RemovableItems    []string `json:"removable_items"`
	InteractiveItems  []string `json:"interactive_items"`
	AdjacentLocations []string `json:"adjacent_locations"`
}

func (l *Location) getNormalizedName() string {
	name := strings.ToLower(l.LocationName)
	return strings.ReplaceAll(name, " ", "_")
}

type World struct {
	Locations       map[string]*Location `json:"locations"`
	CurrentLocation *Location            `json:"current_location"`
}

func (w *World) UpdateCurrentLocation(newLocation *Location) {
	w.CurrentLocation = newLocation
}

func (w *World) VisualizeLocationTree() string {
	visited := make(map[*Location]bool)
	return w.dfs(w.CurrentLocation, visited, "")
}

func (w *World) dfs(location *Location, visited map[*Location]bool, indent string) string {
	if visited[location] {
		return ""
	}

	visited[location] = true
	result := indent + location.LocationName + "\n"

	for _, adjacentLocation := range location.AdjacentLocations {
		realLocation := w.Locations[adjacentLocation]
		result += w.dfs(realLocation, visited, indent+"\t")
	}

	return result
}
