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
	StoryThreads      []string `json:"current_story_threads"`
}

type GameStateUpdateResponse struct {
	CurrentLocation              string   `json:"current_location"`
	AdjacentLocations            []string `json:"adjacent_locations"`
	PlayerInventory              []string `json:"player_inventory"`
	InteractiveObjectsInLocation []string `json:"interactive"`
	EnemiesInLocation            []string `json:"enemies"`
	StoryThreads                 []string `json:"story_threads"`
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

	PreparedStatsCache.Inventory = g.Player.Inventory
	PreparedStatsCache.Enemies = g.World.CurrentLocation.Enemies.ToSlice()
	PreparedStatsCache.InteractiveItems = g.World.CurrentLocation.InteractiveItems.ToSlice()
}

func (g *Game) BuildGameStateDetails() GameStateDetails {
	return GameStateDetails{
		CurrentLocation:  g.World.CurrentLocation.LocationName,
		Inventory:        g.Player.Inventory,
		InteractiveItems: g.World.CurrentLocation.InteractiveItems.ToSlice(),
		Enemies:          g.World.CurrentLocation.Enemies.ToSlice(),
		StoryThreads:     g.World.CurrentLocation.StoryThreads,
	}
}

func (g *Game) UpdateGameHistory(userMessage GameMessage, assistantMessage GameMessage) {
	g.GameMessageHistory = append(g.GameMessageHistory, userMessage)
	g.GameMessageHistory = append(g.GameMessageHistory, assistantMessage)
}

func (g *Game) UpdateGameState(stateUpdate GameStateUpdateResponse) {
	g.handleLocationUpdate(stateUpdate)

	if stateUpdate.PlayerInventory != nil {
		log.Printf("Updating player inventory: %v", stateUpdate.PlayerInventory)
		g.Player.Inventory = stateUpdate.PlayerInventory
	}
}

func (g *Game) handleLocationUpdate(stateUpdate GameStateUpdateResponse) {
	updatedLocationName := stateUpdate.CurrentLocation

	g.World.SafeAddLocation(updatedLocationName)
	location, _ := g.World.GetLocationByName(updatedLocationName)
	if location == nil {
		return
	}

	if location.getNormalizedName() != g.World.CurrentLocation.getNormalizedName() {
		location := g.World.NextLocation(location)
		log.Println("Moving to new location: ", location.LocationName)
	}

	g.World.CurrentLocation.InteractiveItems.AddAll(stateUpdate.InteractiveObjectsInLocation...)
	g.World.CurrentLocation.Enemies.AddAll(stateUpdate.EnemiesInLocation...)
	g.World.CurrentLocation.StoryThreads = stateUpdate.StoryThreads

	if len(stateUpdate.AdjacentLocations) > 0 {
		for _, newLocation := range stateUpdate.AdjacentLocations {
			g.World.CurrentLocation.SafeAddAdjacentLocation(newLocation)
		}
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
