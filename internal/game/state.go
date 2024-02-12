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

type GameStateUpdateResponse struct {
	PlayerLocation              string   `json:"player_location"`
	PotentialLocations          []string `json:"potential_locations"`
	InteactiveObjectsIdentified []string `json:"interactive_objects_identified"`
	InteractiveObjectsRemoved   []string `json:"interactive_objects_removed"`
	EnemiesIdentified           []string `json:"enemies_identified"`
	EnemiesRemoved              []string `json:"enemies_removed"`
	PlayerInventoryAdded        []string `json:"player_inventory_added"`
	PlayerInventoryRemoved      []string `json:"player_inventory_removed"`
	StoryThreads                []string `json:"current_story_threads"`
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

	previousLocation, _ := g.World.GetLocationByName(g.World.PreviousLocationKey)
	if previousLocation != nil {
		PreparedStatsCache.PreviousLocation = previousLocation.LocationName
	} else {
		PreparedStatsCache.PreviousLocation = "Unknown Location"
	}

	PreparedStatsCache.Inventory = g.Player.Inventory.ToSlice()
	PreparedStatsCache.Enemies = g.World.CurrentLocation.Enemies.ToSlice()
	PreparedStatsCache.InteractiveItems = g.World.CurrentLocation.InteractiveItems.ToSlice()
}

func (g *Game) UpdateGameHistory(userMessage GameMessage, assistantMessage GameMessage) {
	g.GameMessageHistory = append(g.GameMessageHistory, userMessage)
	g.GameMessageHistory = append(g.GameMessageHistory, assistantMessage)
}

func (g *Game) UpdateGameState(stateUpdate GameStateUpdateResponse) {
	g.handleLocationUpdate(stateUpdate)

	if len(stateUpdate.PlayerInventoryAdded) > 0 {
		g.Player.Inventory.AddAll(stateUpdate.PlayerInventoryAdded...)
	}

	if len(stateUpdate.PlayerInventoryRemoved) > 0 {
		g.Player.Inventory.RemoveAll(stateUpdate.PlayerInventoryRemoved...)
	}
}

func (g *Game) handleLocationUpdate(stateUpdate GameStateUpdateResponse) {
	potentialLocationName := stateUpdate.PlayerLocation
	newOrExistingLocation := g.World.SafeAddLocation(potentialLocationName)
	currentLocation := g.World.NextLocation(newOrExistingLocation)

	if len(stateUpdate.PotentialLocations) > 0 {
		for _, adjacentLocation := range stateUpdate.PotentialLocations {
			adjacentLocation := g.World.SafeAddLocation(adjacentLocation)
			currentLocation.SafeAddAdjacentLocation(adjacentLocation)
		}

	}

	if len(stateUpdate.InteactiveObjectsIdentified) > 0 {
		currentLocation.InteractiveItems.AddAll(stateUpdate.InteactiveObjectsIdentified...)
	}

	if len(stateUpdate.InteractiveObjectsRemoved) > 0 {
		currentLocation.InteractiveItems.RemoveAll(stateUpdate.InteractiveObjectsRemoved...)
	}

	if len(stateUpdate.EnemiesIdentified) > 0 {
		currentLocation.Enemies.AddAll(stateUpdate.EnemiesIdentified...)
	}

	if len(stateUpdate.EnemiesRemoved) > 0 {
		currentLocation.Enemies.RemoveAll(stateUpdate.EnemiesRemoved...)
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
		return InitializeNewGame(), err
	}

	return decodeGame(encodedGame), nil
}
