package game

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"unicode"

	"github.com/sessionsdev/blue-octopus/internal/aiapi"
	"github.com/sessionsdev/blue-octopus/internal/redis"
)

func init() {
	gob.Register(aiapi.OpenAiMessage{})
}

func (g *Game) AppendToMessageHistory(message GameMessage) {
	g.GameMessageHistory = append(g.GameMessageHistory, message)
}

func (g *Game) UpdateGameState(userMessage GameMessage, stateUpdate GameUpdate, tokensUsed int, assistantMessage GameMessage) {
	g.handleLocationUpdate(stateUpdate)

	if stateUpdate.ProposedStateChanges.UpdatedPlayerInventory != nil {
		g.Player.Inventory = stateUpdate.ProposedStateChanges.UpdatedPlayerInventory
		log.Println("Updated Inventory: ", g.Player.Inventory)
	}

	if stateUpdate.ProposedStateChanges.NewStoryThreads != nil && len(stateUpdate.ProposedStateChanges.NewStoryThreads) > 0 {
		g.StoryThreads = append(g.StoryThreads, stateUpdate.ProposedStateChanges.NewStoryThreads...)
		log.Println("Updated Story Threads: ", g.StoryThreads)
	}

	g.TotalTokensUsed += tokensUsed
	log.Println("Total tokens used: ", g.TotalTokensUsed)

	g.GameMessageHistory = append(g.GameMessageHistory, userMessage)
	g.GameMessageHistory = append(g.GameMessageHistory, assistantMessage)
}

func (g *Game) handleLocationUpdate(stateUpdate GameUpdate) {
	proposedLocationName := stateUpdate.ProposedStateChanges.NewCurrentLocation
	normalized := strings.ReplaceAll(strings.ToLower(proposedLocationName), " ", "_")

	var locationToUpdate *Location

	if proposedLocationName != "" {
		locationToUpdate = g.World.Locations[normalized]

		if locationToUpdate == nil {
			locationToUpdate = &Location{
				LocationName:      proposedLocationName,
				EnemiesInLocation: []string{},
				RemovableItems:    []string{},
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
	if stateUpdate.ProposedStateChanges.UpdatedEnemiesInLocation != nil {
		locationToUpdate.EnemiesInLocation = stateUpdate.ProposedStateChanges.UpdatedEnemiesInLocation
	}

	if stateUpdate.ProposedStateChanges.UpdatedInteractiveObjectsInLocation != nil {
		locationToUpdate.InteractiveItems = stateUpdate.ProposedStateChanges.UpdatedInteractiveObjectsInLocation
	}

	if stateUpdate.ProposedStateChanges.UpdatedRemovableItemsInLocation != nil {
		locationToUpdate.RemovableItems = stateUpdate.ProposedStateChanges.UpdatedRemovableItemsInLocation
	}

	if stateUpdate.ProposedStateChanges.NewAdjacentLocations != nil {

		// for each new adjacent location, add the current location to the adjacent locations
		for _, potientialAdjacentLocation := range stateUpdate.ProposedStateChanges.NewAdjacentLocations {
			// normalize the location names
			normalized := strings.ReplaceAll(strings.ToLower(potientialAdjacentLocation), " ", "_")
			adjacentLocation := g.World.Locations[normalized]

			// if adjacent location does not exist, create it
			if adjacentLocation == nil {
				adjacentLocation = &Location{
					LocationName:      potientialAdjacentLocation,
					EnemiesInLocation: []string{},
					RemovableItems:    []string{},
					InteractiveItems:  []string{},
					AdjacentLocations: []string{locationToUpdate.LocationName},
				}

				// add the new location to the world
				g.World.Locations[normalized] = adjacentLocation
			}
		}
	}

	g.World.UpdateCurrentLocation(locationToUpdate)
	jsonOutput, err := json.MarshalIndent(g.World.CurrentLocation, "", "  ")
	if err != nil {
		log.Printf("JSON marshalling failed: %s", err)
	} else {
		log.Printf("NEW CURRENT LOCATION: %s", jsonOutput)
	}

}

type Player struct {
	Name      string   `json:"name"`
	Inventory []string `json:"inventory"`
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func InitializeNewGame(setupMessage GameMessage) *Game {
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
	newGame.SetupMessage = setupMessage
	newGame.CentralPlot = "The player starts outside a blue house in a forested area, with no specific instructions, but soon discovers that the main goal is to collect treasures and bring them back to the house. As the game progresses, the player navigates through a series of complex underground passages filled with puzzles, traps, and creatures such as grues (dangerous beings that inhabit dark places), all while gathering valuable items. The ultimate objective is to find all the treasures and secure them in the trophy case at the white house, achieving the rank of Master Adventurer. "
	newGame.StoryThreads = []string{
		"The player starts outside a blue house in a forested area with no memory",
		"There is a leaflet hidden in the mailbox explaining the general objective of the game.",
		"A crowbar is required to open the boarded door.  It is not in this location.",
		"There is a hidden and locked celler door in the back of the house.  It requires a key to open.",
	}

	return newGame
}

func buildInitialWorldLocations() map[string]*Location {
	blueHouse := &Location{
		LocationName:      "Blue House",
		EnemiesInLocation: []string{},
		RemovableItems:    []string{"leaflet"},
		InteractiveItems: []string{
			"Mailbox",
			"Borded Door",
		},
		AdjacentLocations: []string{},
	}

	riverLocation := &Location{
		LocationName:      "River",
		EnemiesInLocation: []string{"Mud Man"},
		RemovableItems:    []string{"leaf"},
		InteractiveItems: []string{
			"Boat",
			"Fish",
		},
		AdjacentLocations: []string{
			blueHouse.LocationName,
		},
	}

	theRoadLocation := &Location{
		LocationName:      "The Road",
		EnemiesInLocation: []string{"Bandit"},
		RemovableItems:    []string{"rock"},
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

type GameUpdate struct {
	Response             string                   `json:"response"`
	ProposedStateChanges GameStateResponseDetails `json:"proposed_state_changes"`
}

type GameStatePromptDetails struct {
	CurrentLocation       string   `json:"current_location"`
	AdjacentLocationNames []string `json:"adjacent_locations"`
	Inventory             []string `json:"player_inventory"`
	EnemiesInLocation     []string `json:"enemies_in_location"`
	InteractiveItems      []string `json:"interactive_objects_in_location"`
	RemovableItems        []string `json:"removable_items_in_location"`
	CentralPlot           string   `json:"central_plot"`
	StoryThreads          []string `json:"story_threads"`
}

func (g *Game) BuildGameStatePromptDetails() GameStatePromptDetails {
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

	return GameStatePromptDetails{
		CurrentLocation:       currentLocation.LocationName,
		AdjacentLocationNames: adjacentLocationNames,
		Inventory:             g.Player.Inventory,
		EnemiesInLocation:     currentLocation.EnemiesInLocation,
		InteractiveItems:      currentLocation.InteractiveItems,
		RemovableItems:        currentLocation.RemovableItems,
		CentralPlot:           g.CentralPlot,
		StoryThreads:          g.StoryThreads,
	}
}

func makeCamelCase(str string) string {
	lower := strings.ToLower(str)
	words := strings.Fields(lower)
	for i, word := range words {
		words[i] = string(unicode.ToUpper(rune(word[0]))) + word[1:]
	}

	return strings.Join(words, "")
}

func (gs *GameStatePromptDetails) GetJsonOrEmptyString() string {
	jsonString, err := json.Marshal(gs)
	if err != nil {
		return ""
	}

	return string(jsonString)
}

type GameStateResponseDetails struct {
	NewCurrentLocation                  string   `json:"new_current_location"`
	NewAdjacentLocations                []string `json:"new_adjacent_locations"`
	UpdatedEnemiesInLocation            []string `json:"updated_enemies_in_location"`
	UpdatedInteractiveObjectsInLocation []string `json:"updated_interactive_objects_in_location"`
	UpdatedRemovableItemsInLocation     []string `json:"updated_removable_items_in_location"`
	UpdatedPlayerInventory              []string `json:"updated_player_inventory"`
	NewStoryThreads                     []string `json:"new_story_threads"`
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
