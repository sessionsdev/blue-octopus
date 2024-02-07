package game

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/sessionsdev/blue-octopus/internal/aiapi"
	"github.com/sessionsdev/blue-octopus/internal/redis"
)

func init() {
	gob.Register(aiapi.OpenAiMessage{})
}

type GameStateDetails struct {
	CurrentLocation       string   `json:"current_location"`
	AdjacentLocationNames []string `json:"adjacent_locations"`
	Inventory             []string `json:"player_inventory"`
	InteractiveItems      []string `json:"interactive_objects_in_location"`
	Obstacles             []string `json:"obstacles_in_location"`
	StoryThreads          []string `json:"story_threads"`
}

func (g *Game) BuildGameStateDetails() GameStateDetails {
	currentLocation := g.World.CurrentLocation
	player := g.Player

	if currentLocation == nil {
		currentLocation = &Location{}
	}

	if player == nil {
		player = &Player{}
	}

	adjacentLocationNames := make([]string, len(currentLocation.AdjacentLocations))
	for i, location := range currentLocation.AdjacentLocations {
		adjacentLocationNames[i] = location
	}

	return GameStateDetails{
		CurrentLocation:       currentLocation.LocationName,
		AdjacentLocationNames: adjacentLocationNames,
		Inventory:             g.Player.Inventory,
		InteractiveItems:      currentLocation.InteractiveItems,
		StoryThreads:          g.StoryThreads,
		Obstacles:             currentLocation.Obstacles,
	}
}

func (gs *GameStateDetails) GetJsonOrEmptyString() string {
	jsonString, err := json.Marshal(gs)
	if err != nil {
		return ""
	}

	return string(jsonString)
}

func (g *Game) UpdateGameHistory(userMessage GameMessage, assistantMessage GameMessage) {
	g.GameMessageHistory = append(g.GameMessageHistory, userMessage)
	g.GameMessageHistory = append(g.GameMessageHistory, assistantMessage)
}

func (g *Game) UpdateGameState(stateUpdate GameStateDetails) {
	g.handleLocationUpdate(stateUpdate)

	if stateUpdate.Inventory != nil {
		g.Player.Inventory = stateUpdate.Inventory
		log.Println("Updated Inventory: ", g.Player.Inventory)
	}

	if stateUpdate.StoryThreads != nil {
		g.StoryThreads = stateUpdate.StoryThreads
	}

	g.SaveGameToRedis()
}

func (g *Game) handleLocationUpdate(stateUpdate GameStateDetails) {
	proposedLocationName := stateUpdate.CurrentLocation
	normalized := strings.ReplaceAll(strings.ToLower(proposedLocationName), " ", "_")

	var locationToUpdate *Location

	if proposedLocationName != "" {
		locationToUpdate = g.World.Locations[normalized]

		if locationToUpdate == nil {
			locationToUpdate = &Location{
				LocationName:      proposedLocationName,
				Obstacles:         []string{},
				InteractiveItems:  []string{},
				AdjacentLocations: []string{g.World.CurrentLocation.LocationName},
			}

			g.World.Locations[normalized] = locationToUpdate
		}
	} else {
		// if the location is not provided, no update is needed
		return
	}

	// update the location with the new state
	if stateUpdate.Obstacles != nil {
		locationToUpdate.Obstacles = stateUpdate.Obstacles
	}

	if stateUpdate.InteractiveItems != nil {
		locationToUpdate.InteractiveItems = stateUpdate.InteractiveItems
	}

	if stateUpdate.AdjacentLocationNames != nil {

		// for each new adjacent location, add the current location to the adjacent locations
		for _, potientialAdjacentLocation := range stateUpdate.AdjacentLocationNames {
			// normalize the location names
			normalized := strings.ReplaceAll(strings.ToLower(potientialAdjacentLocation), " ", "_")
			adjacentLocation := g.World.Locations[normalized]

			// if adjacent location does not exist, create it
			if adjacentLocation == nil {
				adjacentLocation = &Location{
					LocationName:      potientialAdjacentLocation,
					Obstacles:         []string{},
					InteractiveItems:  []string{},
					AdjacentLocations: []string{locationToUpdate.LocationName},
				}

				// add the new location to the world
				g.World.Locations[normalized] = adjacentLocation
			}
		}
	}

	g.World.UpdateCurrentLocation(locationToUpdate)
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func InitializeNewGame() *Game {
	// generate a unique UUID for the game
	id, err := randomId()
	if err != nil {
		log.Fatal(err)
	}

	newGame := new(Game)
	newGame.GameId = id
	newGame.TotalTokensUsed = 0

	world := World{
		Locations: buildInitialWorldLocations(),
	}
	world.CurrentLocation = world.Locations["blue_house"]

	player := Player{
		Name:      "Adventurer",
		Inventory: []string{},
	}

	newGame.Player = &player
	newGame.World = &world
	newGame.StoryThreads = []string{
		"The player awoke at the blue house with no memory.",
		"The player must find a way to enter the blue house.",
		"The blue house contains a trophy case that must be filled with various treasures to complete the game.",
	}

	return newGame
}

func buildInitialWorldLocations() map[string]*Location {
	blueHouse := &Location{
		LocationName: "Blue House",
		Obstacles:    []string{"Borded Door", "Locked Cellar Door"},
		InteractiveItems: []string{
			"Mailbox",
			"Borded Door",
		},
		AdjacentLocations: []string{},
	}

	riverLocation := &Location{
		LocationName: "River",
		Obstacles:    []string{"Cross the river"},
		InteractiveItems: []string{
			"Hidden Boat",
			"Fish",
		},
		AdjacentLocations: []string{
			blueHouse.LocationName,
		},
	}

	theRoadLocation := &Location{
		LocationName: "The Road",
		Obstacles:    []string{"Bandin on the road"},
		InteractiveItems: []string{
			"Sign",
			"Tree",
		},
		AdjacentLocations: []string{
			blueHouse.LocationName,
		},
	}

	blueHouse.AdjacentLocations = append(blueHouse.AdjacentLocations, riverLocation.LocationName, theRoadLocation.LocationName)

	locations := make(map[string]*Location)
	locations[theRoadLocation.getNormalizedName()] = theRoadLocation
	locations[riverLocation.getNormalizedName()] = riverLocation
	locations[blueHouse.getNormalizedName()] = blueHouse

	return locations
}

func randomId() (s string, err error) {
	b := make([]byte, 8)
	_, err = rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	s = fmt.Sprintf("%x", b)
	return
}

func (g *Game) encodeGame() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(g)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}

func decodeGame(encodedGame []byte) *Game {
	dec := gob.NewDecoder(bytes.NewReader(encodedGame))
	var game Game
	err := dec.Decode(&game)
	if err != nil {
		log.Fatal(err)
	}

	return &game
}

func (g *Game) SaveGameToRedis() {
	log.Printf("Saving game to redis: %s", g.GameId)

	encodedGame := g.encodeGame()
	err := redis.SetGob(g.GameId, encodedGame, 0)
	if err != nil {
		log.Printf("Error saving game to redis: %s", g.GameId)
	}
}

func LoadGameFromRedis(gameId string) (*Game, error) {
	log.Printf("Loading game from redis: %s", gameId)

	encodedGame, err := redis.GetGob(gameId)
	if err != nil {
		if _, ok := err.(*redis.NotFoundError); ok {
			return nil, err
		} else if err != nil {
			log.Printf("Error loading game from redis: %s", gameId)
			log.Fatal(err)
		}
	}

	return decodeGame(encodedGame), nil
}
