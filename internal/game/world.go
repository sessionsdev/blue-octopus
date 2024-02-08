package game

import (
	"strings"
)

type Location struct {
	LocationName       string
	PreviousLocation   string
	InteractiveItems   []string
	potentialLocations []string
	Obstacles          []string
	Enemies            []string
}

func (l *Location) getNormalizedName() string {
	name := strings.ToLower(l.LocationName)
	return strings.ReplaceAll(name, " ", "_")
}

type World struct {
	Locations       map[string]*Location
	CurrentLocation *Location
}

func (w *World) NextLocation(newLocation *Location) {
	// new location doesn't exist in the world, add it
	w.SafeAddLocation(newLocation)
	newLocation.PreviousLocation = w.CurrentLocation.LocationName
	w.CurrentLocation = newLocation
}

func (w *World) GetLocationByName(locationName string) (*Location, bool) {
	normalized := strings.ReplaceAll(strings.ToLower(locationName), " ", "_")
	location, ok := w.Locations[normalized]
	return location, ok
}

func (w *World) SafeAddLocation(location *Location) {
	normalized := location.getNormalizedName()
	_, ok := w.Locations[normalized]
	if !ok {
		w.Locations[normalized] = location
		return
	}
}
