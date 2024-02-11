package game

import (
	"log"
	"strings"

	utils "github.com/sessionsdev/blue-octopus/internal"
)

type Location struct {
	LocationName         string
	AdjacentLocationKeys utils.StringSet
	InteractiveItems     utils.StringSet
	Enemies              utils.StringSet
	StoryThreads         []string
}

func (l *Location) SafeAddAdjacentLocation(newAdjacentLocation string) {
	// if the location has no adjacent locations, initialize it
	if l.AdjacentLocationKeys == nil {
		l.AdjacentLocationKeys = make(map[string]struct{})
	}

	// if the new adjacent location is empty, return
	if newAdjacentLocation == "" {
		return
	}

	// if the new adjacent location is the same as the current location, return
	normalizedName := normalizedLocationName(newAdjacentLocation)
	if l.getNormalizedName() == normalizedName {
		return
	}

	// if the new adjacent location is not already in the list of adjacent locations, add it
	currentAdjacentLocations := l.AdjacentLocationKeys
	if !currentAdjacentLocations.Contains(normalizedName) {
		l.AdjacentLocationKeys.AddAll(normalizedName)
	}
}

func (l *Location) getNormalizedName() string {
	return normalizedLocationName(l.LocationName)
}

type World struct {
	Locations           map[string]*Location
	CurrentLocation     *Location
	PreviousLocationKey string
	VisitedLocations    utils.StringSet
}

func (w *World) NextLocation(nextLocation *Location) *Location {
	if nextLocation == nil {
		return w.CurrentLocation
	}

	log.Printf("Current location: %s", w.CurrentLocation.LocationName)
	log.Printf("Next location: %s", nextLocation.LocationName)

	// determine which direction we moved by looking at the current location directions and the next location name

	if w.VisitedLocations == nil {
		w.VisitedLocations = utils.EmptyStringSet()
	}

	// update visited locations
	w.VisitedLocations.AddAll(w.CurrentLocation.LocationName, nextLocation.LocationName)

	// update previous, current locations and adjacent locations
	w.CurrentLocation.AdjacentLocationKeys.AddAll(nextLocation.getNormalizedName())
	nextLocation.AdjacentLocationKeys.AddAll(w.CurrentLocation.getNormalizedName())

	w.PreviousLocationKey = w.CurrentLocation.getNormalizedName()
	w.CurrentLocation = nextLocation
	return w.CurrentLocation
}

func (w *World) GetLocationByName(locationName string) (*Location, bool) {
	if locationName == "" {
		return nil, false
	}

	normalized := strings.ReplaceAll(strings.ToLower(locationName), " ", "_")
	location, ok := w.Locations[normalized]
	return location, ok
}

func (w *World) SafeAddLocation(locationName string) *Location {
	if w.Locations == nil {
		w.Locations = make(map[string]*Location)
	}

	if locationName == "" {
		return nil
	}

	// normalize the location name
	name := strings.ToLower(locationName)
	normalizedName := strings.ReplaceAll(name, " ", "_")

	// check if the location already exists
	location, ok := w.Locations[normalizedName]

	// if it doesn't exist, create it
	if !ok {
		location = &Location{
			LocationName:         locationName,
			AdjacentLocationKeys: utils.EmptyStringSet(),
			InteractiveItems:     utils.EmptyStringSet(),
			Enemies:              utils.EmptyStringSet(),
		}

		// add the location to the world
		log.Println("Adding location: ", locationName)
		w.Locations[normalizedName] = location
	}

	return location
}

func normalizedLocationName(locationName string) string {
	return strings.ReplaceAll(strings.ToLower(locationName), " ", "_")
}
