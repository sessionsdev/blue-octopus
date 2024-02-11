package game

import (
	"crypto/rand"
	"fmt"
	"log"
)

func InitializeNewGame() *Game {
	newGameDetails := NewGameDetails{
		StartingLocation: "Blue House",
		PlayerName:       "Adventurer",
		StartingStoryThreads: []string{"The player awoke at the blue house with no memory.",
			"The player must find a way to enter the locked blue house.",
			"The five trophies are scattered throughout the world."},
		PlayerInventory:           []string{},
		StartingAdjacentLocations: []string{"River", "Eastern Road"},
		MainQuest:                 "Acquire the five trophies and place them in the trophy case inside the blue house.",
	}

	newGame := BuildNewGame(newGameDetails)
	newGame.GameId = randomId()
	newGame.TotalTokensUsed = 0

	return newGame
}

func randomId() (s string) {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	s = fmt.Sprintf("%x", b)
	return
}
