package game

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"log"

	utils "github.com/sessionsdev/blue-octopus/internal"
	"github.com/sessionsdev/blue-octopus/internal/redis"
)

func init() {
	gob.Register(GameMessage{})
	gob.Register(map[string]struct{}{})
}

type GameStateDetails struct {
	CurrentLocation       string   `json:"current_location"`
	AdjacentLocationNames []string `json:"adjacent_locations"`
	Inventory             []string `json:"player_inventory"`
	InteractiveItems      []string `json:"interactive_objects"`
	Enemies               []string `json:"enemies"`
	StoryThreads          []string `json:"story_threads"`
}

type GameStateUpdateResponse struct {
	CurrentLocation    string   `json:"current_location"`
	PotentialLocations []string `json:"potential_locations"`
	InventoryUpdates   struct {
		Added   []string `json:"added"`
		Removed []string `json:"removed"`
	} `json:"inventory_updates"`
	InteractiveObjectsUpdates struct {
		Added   []string `json:"added"`
		Removed []string `json:"removed"`
	} `json:"interactive_objects_updates"`
	EnemiesUpdates struct {
		Added    []string `json:"added"`
		Defeated []string `json:"defeated"`
	} `json:"enemies_updates"`
	StoryThreads []string `json:"story_threads"`
}

func (g *Game) BuildGameStateDetails() GameStateDetails {
	utils.GetStringsFromMap(g.World.CurrentLocation.PotentialLocations)
	return GameStateDetails{
		CurrentLocation:       g.World.CurrentLocation.LocationName,
		AdjacentLocationNames: utils.GetStringsFromMap(g.World.CurrentLocation.PotentialLocations),
		Inventory:             utils.GetStringsFromMap(g.Player.Inventory),
		InteractiveItems:      utils.GetStringsFromMap(g.World.CurrentLocation.InteractiveItems),
		Enemies:               utils.GetStringsFromMap(g.World.CurrentLocation.Enemies),
		StoryThreads:          g.StoryThreads,
	}
}

func (g *Game) UpdateGameHistory(userMessage GameMessage, assistantMessage GameMessage) {
	g.GameMessageHistory = append(g.GameMessageHistory, userMessage)
	g.GameMessageHistory = append(g.GameMessageHistory, assistantMessage)
}

func (g *Game) UpdateGameState(stateUpdate GameStateUpdateResponse) {
	g.handleLocationUpdate(stateUpdate)

	if stateUpdate.InventoryUpdates.Added != nil {
		log.Printf("Adding items to inventory: %v", stateUpdate.InventoryUpdates.Added)
		for _, item := range stateUpdate.InventoryUpdates.Added {
			g.Player.Inventory[item] = struct{}{}
		}
	}

	if stateUpdate.InventoryUpdates.Removed != nil {
		log.Printf("Removing items from inventory: %v", stateUpdate.InventoryUpdates.Removed)
		for _, item := range stateUpdate.InventoryUpdates.Removed {
			delete(g.Player.Inventory, item)
		}
	}

	if stateUpdate.StoryThreads != nil {
		g.StoryThreads = stateUpdate.StoryThreads
	}

	g.SaveGameToRedis()
}

func (g *Game) handleLocationUpdate(stateUpdate GameStateUpdateResponse) {
	proposedLocation := stateUpdate.CurrentLocation
	if proposedLocation != "" {
		location := g.World.NextLocation(proposedLocation)
		log.Println("Moving to new location: ", location.LocationName)
	}

	if stateUpdate.PotentialLocations != nil {
		log.Println("Updating potential locations: ", stateUpdate.PotentialLocations)
		for _, locationName := range stateUpdate.PotentialLocations {
			g.World.CurrentLocation.PotentialLocations[locationName] = struct{}{}
		}
	}

	if stateUpdate.InteractiveObjectsUpdates.Added != nil {
		log.Printf("Adding interactive objects: %v", stateUpdate.InteractiveObjectsUpdates.Added)
		for _, obj := range stateUpdate.InteractiveObjectsUpdates.Added {
			g.World.CurrentLocation.InteractiveItems[obj] = struct{}{}
		}
	}

	if stateUpdate.InteractiveObjectsUpdates.Removed != nil {
		log.Printf("Removing interactive objects: %v", stateUpdate.InteractiveObjectsUpdates.Removed)
		for _, obj := range stateUpdate.InteractiveObjectsUpdates.Removed {
			delete(g.World.CurrentLocation.InteractiveItems, obj)
		}
	}

	if stateUpdate.EnemiesUpdates.Added != nil {
		log.Printf("Adding enemies: %v", stateUpdate.EnemiesUpdates.Added)
		for _, enemy := range stateUpdate.EnemiesUpdates.Added {
			g.World.CurrentLocation.Enemies[enemy] = struct{}{}
		}
	}

	if stateUpdate.EnemiesUpdates.Defeated != nil {
		log.Printf("Defeating enemies: %v", stateUpdate.EnemiesUpdates.Defeated)
		for _, enemy := range stateUpdate.EnemiesUpdates.Defeated {
			delete(g.World.CurrentLocation.Enemies, enemy)
		}
	}

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
		Inventory: map[string]struct{}{},
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
		InteractiveItems: map[string]struct{}{
			"Borded Door": {},
		},
		PotentialLocations: map[string]struct{}{
			"River":    {},
			"The Road": {},
		},
		Enemies: map[string]struct{}{
			"Guard Dog": {},
		},
	}

	riverLocation := &Location{
		LocationName: "River",
		InteractiveItems: map[string]struct{}{
			"Boat": {},
		},
		PreviousLocation:   "Blue House",
		PotentialLocations: map[string]struct{}{},
	}

	theRoadLocation := &Location{
		LocationName: "The Road",
		InteractiveItems: map[string]struct{}{
			"Sign": {},
		},
		PreviousLocation:   "Blue House",
		PotentialLocations: map[string]struct{}{},
		Enemies: map[string]struct{}{
			"Bandit": {},
		},
	}

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
		log.Fatalf("Error loading game from redis: %s", gameId)
	}

	return decodeGame(encodedGame), nil
}

func safeRemoveString(slice []string, item string) []string {
	for i, s := range slice {
		if s == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func safeAddString(slice []string, item string) []string {
	for _, s := range slice {
		if s == item {
			return slice
		}
	}
	return append(slice, item)
}
