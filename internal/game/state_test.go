package game

import (
	"testing"
)

func TestUpdateGameState(t *testing.T) {
	newGameDetails := NewGameDetails{
		StartingLocation:          "Test Current Location",
		PlayerName:                "Test Player",
		PlayerInventory:           []string{"item1", "item2"},
		StartingAdjacentLocations: []string{"Test Adjacent Location"},
		MainQuest:                 "Test Main Quest",
	}

	testGame := BuildNewGame(newGameDetails)

	stateUpdate := GameStateUpdateResponse{
		CurrentLocation:   "Test Current Location",
		AdjacentLocations: []string{},
		PlayerInventory: []string{
			"item1",
			"item2",
		},
		InteractiveObjectsInLocation: []string{},
		EnemiesInLocation:            []string{},
		StoryThreads:                 []string{"Test Story Thread 2"},
	}

	// Test with current values stateUpdate
	testGame.UpdateGameState(stateUpdate)
	if testGame.World.CurrentLocation.LocationName != "Test Current Location" {
		t.Errorf("Expected current location to be 'Test Current Location', but got %s", testGame.World.CurrentLocation.LocationName)
	}

	// Test with some inventory updates
	stateUpdate.PlayerInventory = []string{"item1", "item2"}
	testGame.UpdateGameState(stateUpdate)
	if len(testGame.Player.Inventory) != 2 {
		t.Errorf("Expected inventory to have 2 items, but got %d", len(testGame.Player.Inventory))
	}

	// Test with some story threads
	stateUpdate.StoryThreads = []string{"Test Story Thread 3"}
	testGame.UpdateGameState(stateUpdate)

	// Test with some location updates
	stateUpdate.CurrentLocation = "Test New Location"
	stateUpdate.AdjacentLocations = []string{"Test Even Newer Location"}
	stateUpdate.PlayerInventory = []string{"item3", "item4"}
	stateUpdate.InteractiveObjectsInLocation = []string{"object1", "object2"}
	stateUpdate.EnemiesInLocation = []string{"enemy1", "enemy2"}
	stateUpdate.StoryThreads = []string{"Test Story Thread 4"}
	testGame.UpdateGameState(stateUpdate)
}
