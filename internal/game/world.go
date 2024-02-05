package game

import (
	"fmt"
	"strings"
)

type Location struct {
	LocationName      string      `json:"location_name"`
	EnemiesInLocation []string    `json:"enemies_in_location"`
	RemovableItems    []string    `json:"removable_items"`
	InteractiveItems  []string    `json:"interactive_items"`
	AdjacentLocations []*Location `json:"adjacent_locations"`
}

func (l *Location) getNormalizedName() string {
	name := strings.ToLower(l.LocationName)
	return strings.ReplaceAll(name, " ", "_")
}

type LocationJson struct {
	LocationName      string   `json:"location_name"`
	EnemiesInLocation []string `json:"enemies_in_location"`
	RemovableItems    []string `json:"removable_items"`
	InteractiveItems  []string `json:"interactive_items"`
	AdjacentLocations []string `json:"adjacent_locations"`
}

func (l *Location) getLocationJson() LocationJson {
	adjacentLocationNames := make([]string, len(l.AdjacentLocations))
	for i, location := range l.AdjacentLocations {
		adjacentLocationNames[i] = location.LocationName
	}

	return LocationJson{
		LocationName:      l.LocationName,
		EnemiesInLocation: l.EnemiesInLocation,
		RemovableItems:    l.RemovableItems,
		InteractiveItems:  l.InteractiveItems,
		AdjacentLocations: adjacentLocationNames,
	}
}

type World struct {
	Locations       map[string]*Location `json:"locations"`
	CurrentLocation *Location            `json:"current_location"`
}

type WorldJson struct {
	Locations       map[string]LocationJson `json:"locations"`
	CurrentLocation LocationJson            `json:"current_location"`
}

func (w *World) GetWorldJson() WorldJson {
	locations := make(map[string]LocationJson)
	for name, location := range w.Locations {
		locations[name] = location.getLocationJson()
	}

	return WorldJson{
		Locations:       locations,
		CurrentLocation: w.CurrentLocation.getLocationJson(),
	}
}

func (w *World) UpdateCurrentLocation(newLocation *Location) {
	w.CurrentLocation = newLocation
}

func (w *World) VisualizeLocationTree() {
	visited := make(map[*Location]bool)
	w.dfs(w.CurrentLocation, visited, "")
}

func (w *World) dfs(location *Location, visited map[*Location]bool, indent string) {
	if visited[location] {
		return
	}

	visited[location] = true
	fmt.Println(indent + location.LocationName)

	for _, adjacentLocation := range location.AdjacentLocations {
		w.dfs(adjacentLocation, visited, indent+"\t")
	}
}
