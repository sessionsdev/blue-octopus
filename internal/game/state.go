package game

import (
	"bytes"
	"encoding/gob"
	"log"

	utils "github.com/sessionsdev/blue-octopus/internal"
	"github.com/sessionsdev/blue-octopus/internal/redis"
)

func init() {
	gob.Register(GameMessage{})
	gob.Register(map[string]struct{}{})
	gob.Register(utils.StringSet{})
}

type GameStateDetails struct {
	CurrentLocation   string   `json:"current_location"`
	AdjacentLocations []string `json:"adjacent_locations"`
	Inventory         []string `json:"player_inventory"`
	InteractiveItems  []string `json:"interactive_objects"`
	Enemies           []string `json:"enemies"`
	StoryThreads      []string `json:"story_threads"`
}

type GameStateUpdateResponse struct {
	CurrentLocation          string   `json:"current_location"`
	UpdatedAdjacentLocations []string `json:"updated_adjacent_locations"`
	InventoryUpdates         struct {
		Added   []string `json:"added"`
		Removed []string `json:"removed"`
	} `json:"player_inventory_updates"`
	InteractiveObjectsUpdates struct {
		Added   []string `json:"added"`
		Removed []string `json:"removed"`
	} `json:"interactive_objects_in_location_updates"`
	EnemiesUpdates struct {
		Added    []string `json:"added"`
		Defeated []string `json:"defeated"`
	} `json:"enemies_in_location_updates"`
	StoryThreads []string `json:"story_threads"`
}

type PreparedStats struct {
	Location         string
	PreviousLocation string
	Inventory        []string
	Enemies          []string
	InteractiveItems []string
}

var PreparedStatsCache *PreparedStats

func (g *Game) populatePreparedStatsCache() {
	PreparedStatsCache = &PreparedStats{}

	location := g.World.CurrentLocation
	if location != nil {
		PreparedStatsCache.Location = location.LocationName
	} else {
		PreparedStatsCache.Location = "Unknown Location"
	}

	previousLocationName := g.World.SafePreviousLocation().LocationName
	PreparedStatsCache.PreviousLocation = previousLocationName

	PreparedStatsCache.Inventory = g.Player.Inventory.ToSlice()
	PreparedStatsCache.Enemies = g.World.CurrentLocation.Enemies.ToSlice()
	PreparedStatsCache.InteractiveItems = g.World.CurrentLocation.InteractiveItems.ToSlice()
}

func (g *Game) BuildGameStateDetails() GameStateDetails {
	return GameStateDetails{
		CurrentLocation:  g.World.CurrentLocation.LocationName,
		Inventory:        utils.GetStringsFromMap(g.Player.Inventory),
		InteractiveItems: utils.GetStringsFromMap(g.World.CurrentLocation.InteractiveItems),
		Enemies:          utils.GetStringsFromMap(g.World.CurrentLocation.Enemies),
		StoryThreads:     g.StoryThreads,
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
		g.Player.Inventory.AddAll(stateUpdate.InventoryUpdates.Added...)
	}

	if stateUpdate.InventoryUpdates.Removed != nil {
		log.Printf("Removing items from inventory: %v", stateUpdate.InventoryUpdates.Removed)
		g.Player.Inventory.RemoveAll(stateUpdate.InventoryUpdates.Removed...)
	}

	if stateUpdate.StoryThreads != nil {
		g.StoryThreads = stateUpdate.StoryThreads
	}
}

func (g *Game) handleLocationUpdate(stateUpdate GameStateUpdateResponse) {
	newCurrentLocation, ok := g.World.SafeAddLocation(stateUpdate.CurrentLocation)
	if !ok {
		log.Printf("Error adding location: %s", stateUpdate.CurrentLocation)
		return
	}

	if newCurrentLocation.getNormalizedName() != g.World.CurrentLocation.getNormalizedName() {
		location := g.World.NextLocation(newCurrentLocation)
		log.Println("Moving to new location: ", location.LocationName)
	}

	if len(stateUpdate.InteractiveObjectsUpdates.Added) != 0 {
		log.Printf("Adding interactive objects: %v", stateUpdate.InteractiveObjectsUpdates.Added)
		g.World.CurrentLocation.InteractiveItems.AddAll(stateUpdate.InteractiveObjectsUpdates.Added...)
	}

	if len(stateUpdate.InteractiveObjectsUpdates.Removed) != 0 {
		log.Printf("Removing interactive objects: %v", stateUpdate.InteractiveObjectsUpdates.Removed)
		g.World.CurrentLocation.InteractiveItems.RemoveAll(stateUpdate.InteractiveObjectsUpdates.Removed...)
	}

	if len(stateUpdate.EnemiesUpdates.Added) != 0 {
		log.Printf("Adding enemies: %v", stateUpdate.EnemiesUpdates.Added)
		g.World.CurrentLocation.Enemies.AddAll(stateUpdate.EnemiesUpdates.Added...)
	}

	if len(stateUpdate.EnemiesUpdates.Defeated) != 0 {
		log.Printf("Defeating enemies: %v", stateUpdate.EnemiesUpdates.Defeated)
		g.World.CurrentLocation.Enemies.RemoveAll(stateUpdate.EnemiesUpdates.Defeated...)
	}

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
