package game

import (
	"log"
	"strings"

	utils "github.com/sessionsdev/blue-octopus/internal"
)

type Location struct {
	LocationName      string
	AdjacentLocations []string
	InteractiveItems  []string
	Enemies           []string
}

func (l *Location) SafeAddAdjacentLocation(newAdjacentLocation string) []string {
	if l.AdjacentLocations == nil {
		l.AdjacentLocations = []string{}
	}

	currentAdjacentLocations := l.AdjacentLocations
	if !utils.Contains(currentAdjacentLocations, newAdjacentLocation) {
		l.AdjacentLocations = append(currentAdjacentLocations, newAdjacentLocation)
	}

	return l.AdjacentLocations
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

	log.Printf("Current location: %s", w.CurrentLocation.LocationName)
	log.Printf("Next location: %s", nextLocation.LocationName)

	// determine which direction we moved by looking at the current location directions and the next location name

	if w.VisitedLocations == nil {
		w.VisitedLocations = utils.EmptyStringSet()
	}

	// update visited locations
	w.VisitedLocations.AddAll(w.CurrentLocation.LocationName, nextLocation.LocationName)

	// update previous, current locations and adjacent locations
	w.CurrentLocation.AdjacentLocations = append(w.CurrentLocation.AdjacentLocations, nextLocation.LocationName)
	nextLocation.AdjacentLocations = append(nextLocation.AdjacentLocations, w.CurrentLocation.LocationName)

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
	log.Println("Location already exists: ", locationName)

	// if it doesn't exist, create it
	if !ok {
		location = &Location{
			LocationName:      locationName,
			AdjacentLocations: []string{},
			InteractiveItems:  []string{},
			Enemies:           []string{},
		}

		// add the location to the world
		log.Println("Adding location: ", locationName)
		w.Locations[normalizedName] = location
	}

	// return the location and whether it was added
	return location
}

func (w *World) SafePreviousLocation() *Location {
	if w.PreviousLocation == nil {
		return &Location{
			LocationName: "Unknown",
		}
	}

	return w.PreviousLocation
}
