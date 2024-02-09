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
		Inventory: utils.EmptyStringSet(),
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
	blueHouse, _ := world.SafeAddLocation("Blue House")
	blueHouse.InteractiveItems.AddAll("Locked Door")
	blueHouse.AdjacentLocations.AddAll("Eastern Road", "River")

	river, _ := world.SafeAddLocation("River")
	river.InteractiveItems.AddAll("Boat", "Risky Bridge")
	river.AdjacentLocations.AddAll("Blue House")

	easternRoad, _ := world.SafeAddLocation("Eastern Road")
	easternRoad.InteractiveItems.AddAll("Sign")
	easternRoad.Enemies.AddAll("Goblin")
	easternRoad.AdjacentLocations.AddAll("Blue House")
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
			Inventory: utils.EmptyStringSet(),
		},
		MainQuest:          "",
		StoryThreads:       []string{},
		GameMessageHistory: []GameMessage{},
		TotalTokensUsed:    0,
	}

	return game
}
