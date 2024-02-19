package game

func InitializeNewGame() *Game {
	newGameDetails := NewGameDetails{
		StartingLocation:          "Blue House",
		PlayerName:                "Adventurer",
		PlayerInventory:           []string{},
		StartingAdjacentLocations: []string{"River", "Eastern Road"},
	}

	newGame := BuildNewGame(newGameDetails)
	newGame.TotalTokensUsed = 0

	return newGame
}
