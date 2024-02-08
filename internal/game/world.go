package game

import (
	"strings"
)

type Location struct {
	LocationName       string
	PreviousLocation   string
	InteractiveItems   map[string]struct{}
	PotentialLocations map[string]struct{}
	Enemies            map[string]struct{}
}

func (l *Location) getNormalizedName() string {
	name := strings.ToLower(l.LocationName)
	return strings.ReplaceAll(name, " ", "_")
}

type World struct {
	Locations       map[string]*Location
	CurrentLocation *Location
}

func (w *World) NextLocation(newLocation string) *Location {
	// new location doesn't exist in the world, add it
	location := w.SafeAddLocation(newLocation)
	location.PreviousLocation = w.CurrentLocation.LocationName
	w.CurrentLocation = location

	return w.CurrentLocation
}

func (w *World) GetLocationByName(locationName string) (*Location, bool) {
	normalized := strings.ReplaceAll(strings.ToLower(locationName), " ", "_")
	location, ok := w.Locations[normalized]
	return location, ok
}

func (w *World) SafeAddLocation(locationName string) *Location {
	location := &Location{
		LocationName:       locationName,
		InteractiveItems:   make(map[string]struct{}),
		PotentialLocations: make(map[string]struct{}),
		Enemies:            make(map[string]struct{}),
	}

	normalizedName := location.getNormalizedName()
	location, ok := w.Locations[normalizedName]
	if !ok {
		w.Locations[normalizedName] = location
		return location
	}
	return location
}
