package game

import (
	"strings"

	utils "github.com/sessionsdev/blue-octopus/internal"
)

type Location struct {
	LocationName      string
	AdjacentLocations utils.StringSet
	InteractiveItems  utils.StringSet
	Enemies           utils.StringSet
}

func (l *Location) getNormalizedName() string {
	name := strings.ToLower(l.LocationName)
	return strings.ReplaceAll(name, " ", "_")
}

type World struct {
	Locations        map[string]*Location
	CurrentLocation  *Location
	PreviousLocation *Location
	VisitedLocations utils.StringSet
}

func (w *World) NextLocation(nextLocation *Location) *Location {
	if nextLocation == nil {
		return w.CurrentLocation
	}

	// determine which direction we moved by looking at the current location directions and the next location name

	if w.VisitedLocations == nil {
		w.VisitedLocations = utils.EmptyStringSet()
	}

	// update visited locations
	w.VisitedLocations.AddAll(w.CurrentLocation.LocationName, nextLocation.LocationName)

	// update previous, current locations and adjacent locations
	w.CurrentLocation.AdjacentLocations.AddAll(nextLocation.LocationName)
	nextLocation.AdjacentLocations.AddAll(w.CurrentLocation.LocationName)
	w.PreviousLocation = w.CurrentLocation
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

func (w *World) SafeAddLocation(locationName string) (*Location, bool) {
	if w.Locations == nil {
		w.Locations = make(map[string]*Location)
	}

	if locationName == "" {
		return nil, false
	}

	// normalize the location name
	name := strings.ToLower(locationName)
	normalizedName := strings.ReplaceAll(name, " ", "_")

	// check if the location already exists
	location, ok := w.Locations[normalizedName]

	// if it doesn't exist, create it
	if !ok {
		location = &Location{
			LocationName:     locationName,
			InteractiveItems: utils.EmptyStringSet(),
			Enemies:          utils.EmptyStringSet(),
		}

		// add the location to the world
		w.Locations[location.getNormalizedName()] = location
	}

	// return the location and whether it was added
	return location, true
}

func (w *World) SafePreviousLocation() *Location {
	if w.PreviousLocation == nil {
		return &Location{
			LocationName: "Unknown",
		}
	}

	return w.PreviousLocation
}
