package game

import (
	"crypto/rand"
	"fmt"
	"log"

	utils "github.com/sessionsdev/blue-octopus/internal"
)

func InitializeNewGame() *Game {
	// generate a unique UUID for the game
	id, err := randomId()
	if err != nil {
		log.Fatal(err)
	}

	newGame := initializeEmptyGame()
	newGame.GameId = id
	newGame.TotalTokensUsed = 0

	world := &World{}
	buildInitialWorldLocations(world)
	world.CurrentLocation = world.Locations["blue_house"]

	player := &Player{
		Name:      "Adventurer",
		Inventory: []string{},
	}

	newGame.Player = player
	newGame.World = world
	newGame.StoryThreads = []string{
		"The player awoke at the blue house with no memory.",
		"The player must find a way to enter the locked blue house.",
		"The five trophies are scattered throughout the world.",
	}
	newGame.MainQuest = "Acquire the five trophies and place them in the trophy case inside the blue house."

	return newGame
}

func buildInitialWorldLocations(world *World) {
	blueHouse := world.SafeAddLocation("Blue House")
	blueHouse.InteractiveItems = append(blueHouse.InteractiveItems, "Locked Door")
	blueHouse.AdjacentLocations = append(blueHouse.AdjacentLocations, "River", "Eastern Road")

	river := world.SafeAddLocation("River")
	river.InteractiveItems = append(river.InteractiveItems, "Bridge")
	river.AdjacentLocations = append(river.AdjacentLocations, "Blue House")

	easternRoad := world.SafeAddLocation("Eastern Road")
	easternRoad.InteractiveItems = append(easternRoad.InteractiveItems, "Cave Entrance")
	easternRoad.Enemies = append(easternRoad.Enemies, "Goblin")
	easternRoad.AdjacentLocations = append(easternRoad.AdjacentLocations, "Blue House", "Cave")
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

func initializeEmptyGame() *Game {
	game := &Game{
		GameId: "",
		World: &World{
			Locations:        make(map[string]*Location),
			CurrentLocation:  &Location{},
			PreviousLocation: &Location{},
			VisitedLocations: utils.EmptyStringSet(),
		},
		Player: &Player{
			Name:      "",
			Inventory: []string{},
		},
		MainQuest:          "",
		StoryThreads:       []string{},
		GameMessageHistory: []GameMessage{},
		TotalTokensUsed:    0,
	}

	return game
}
