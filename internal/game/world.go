package game

import (
	"log"
	"strings"

	"github.com/sessionsdev/blue-octopus/internal/util"
)

type Location struct {
	LocationName         string
	AdjacentLocationKeys util.StringSet
	InteractiveItems     util.StringSet
	Enemies              util.StringSet
}

func (l *Location) SafeAddAdjacentLocation(adjLocation *Location) {
	// if the location has no adjacent locations, initialize it
	if l.AdjacentLocationKeys == nil {
		l.AdjacentLocationKeys = make(map[string]struct{})
	}

	// if the new adjacent location is empty, return
	if adjLocation == nil {
		return
	}

	// if the new adjacent location is the same as the current location, return
	if l.getNormalizedName() == adjLocation.getNormalizedName() {
		return
	}

	// if the new adjacent location is not already in the list of adjacent locations, add it
	if !l.AdjacentLocationKeys.Contains(adjLocation.getNormalizedName()) {
		l.AdjacentLocationKeys.AddAll(adjLocation.getNormalizedName())
	}
}

func (l *Location) getNormalizedName() string {
	return normalizedLocationName(l.LocationName)
}

type World struct {
	Locations           map[string]*Location
	CurrentLocation     *Location
	PreviousLocationKey string
}

func (w *World) NextLocation(nextLocation *Location) *Location {
	if nextLocation == nil {
		return w.CurrentLocation
	}

	if w.CurrentLocation == nil {
		w.CurrentLocation = nextLocation
		return w.CurrentLocation
	}

	if w.CurrentLocation.getNormalizedName() == nextLocation.getNormalizedName() {
		return w.CurrentLocation
	}

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
			AdjacentLocationKeys: util.EmptyStringSet(),
			InteractiveItems:     util.EmptyStringSet(),
			Enemies:              util.EmptyStringSet(),
		}

		// add the location to the world
		log.Println("Adding location: ", locationName)
		w.Locations[normalizedName] = location
	}

	return location
}

func (w *World) GetAllLocationNames() []string {
	var locationNames []string
	for _, value := range w.Locations {
		locationNames = append(locationNames, value.LocationName)
	}
	return locationNames
}

func normalizedLocationName(locationName string) string {
	return strings.ReplaceAll(strings.ToLower(locationName), " ", "_")
}
