package game

import (
	"crypto/rand"
	"fmt"
	"log"
)

func InitializeNewGame() *Game {
	newGameDetails := NewGameDetails{
		StartingLocation:          "Blue House",
		PlayerName:                "Adventurer",
		PlayerInventory:           []string{},
		StartingAdjacentLocations: []string{"River", "Eastern Road"},
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
