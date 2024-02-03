package game

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type Game struct {
	GameId          string `json:"game_id"`
	World           World  `json:"world"`
	Player          Player `json:"player"`
	TotalTokensUsed int    `json:"total_tokens_used"`
}

func (g *Game) GetJsonRepresentation() string {
	json, err := json.Marshal(g)
	if err != nil {
		log.Fatal(err)
	}

	return string(json)
}

func (g *Game) GetPartialJsonRepresentation() string {
	partialGame := PartialGame{
		Player: g.Player,
		World:  g.World,
	}

	json, err := json.Marshal(partialGame)
	if err != nil {
		log.Fatal(err)
	}

	return string(json)
}

func (g *Game) UpdateGameState(proposedStateChanges ProposedStateChanges, tokensUsed int) {
	proposedLocation := proposedStateChanges.NewCurrentLocation

	if proposedLocation != "" && proposedLocation != g.World.CurrentLocation.LocationName {
		for _, location := range g.World.Locations {
			if location.LocationName == proposedLocation {
				g.World.UpdateLocation(&location)
				log.Println("Location updated: ", location)
				break
			} else {
				newLocation := Location{
					LocationName:      proposedLocation,
					EnemiesInLocation: proposedStateChanges.UpdatedEnemiesInLocation,
				}
				g.World.addNewLocation(newLocation)
				g.World.UpdateLocation(&newLocation)
				log.Println("New location added: ", newLocation)
				break
			}
		}
	}

	updatedInventory := proposedStateChanges.UpdatedInventory
	if len(updatedInventory) > 0 {
		g.Player.Inventory = updatedInventory
		log.Println("Updated inventory: ", g.Player.Inventory)
	}

	g.TotalTokensUsed += tokensUsed
}

type ProposedStateChanges struct {
	NewCurrentLocation       string   `json:"current_location_name"`
	InteractiveObjects       []string `json:"interactive_objects_in_location"`
	RemovableItems           []string `json:"removable_items"`
	UpdatedEnemiesInLocation []string `json:"updated_enemies_in_location"`
	UpdatedInventory         []string `json:"player_inventory"`
}

type PartialGame struct {
	Player Player `json:"player"`
	World  World  `json:"world"`
}

type Player struct {
	Name      string   `json:"name"`
	Inventory []string `json:"inventory"`
}

type Location struct {
	LocationName          string   `json:"location_name"`
	EnemiesInLocation     []string `json:"enemies_in_location"`
	RemovableItems        []string `json:"removable_items"`
	InteractiveItems      []string `json:"interactive_items"`
	AdjacentLocationNames []string `json:"adjacent_location_names"`
}

func (l *Location) getLowerCaseLocationName() string {
	return strings.ToLower(l.LocationName)
}

type World struct {
	Locations       map[string]Location `json:"locations"`
	CurrentLocation *Location           `json:"current_location"`
}

func (w *World) addNewLocation(newLocation Location) {
	w.Locations[newLocation.LocationName] = newLocation
}

func (w *World) UpdateLocation(newLocation *Location) {
	w.CurrentLocation = newLocation
}

func InitializeNewGame() Game {

	forestLocation := Location{
		LocationName: "Forest",
		EnemiesInLocation: []string{
			"Orc",
			"Goblin",
		},
		RemovableItems: []string{
			"Key",
		},
		InteractiveItems: []string{
			"Tree",
		},
		AdjacentLocationNames: []string{
			"Cave",
		},
	}

	caveLocation := Location{
		LocationName: "Cave",
		EnemiesInLocation: []string{
			"Dragon",
		},
		RemovableItems: []string{
			"Treasure Chest",
		},
		InteractiveItems: []string{
			"Stalactite",
		},
		AdjacentLocationNames: []string{
			"Forest",
		},
	}

	// generate a unique UUID for the game
	id, err := randomId()
	if err != nil {
		log.Fatal(err)
	}

	world := World{
		Locations:       map[string]Location{},
		CurrentLocation: &forestLocation,
	}

	world.addNewLocation(forestLocation)
	world.addNewLocation(caveLocation)

	player := Player{
		Name: "Adventurer",
		Inventory: []string{
			"Sword",
			"Shield",
			"Health Potion",
		},
	}

	game := new(Game)
	game.Player = player

	game.GameId = id
	game.World = world

	return *game
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
