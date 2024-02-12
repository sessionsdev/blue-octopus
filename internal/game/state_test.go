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
	}

	testGame := BuildNewGame(newGameDetails)

	stateUpdate := GameStateUpdateResponse{
		PlayerLocation:              "Test New Location",
		PotentialLocations:          []string{"Test Even Newer Location"},
		InteactiveObjectsIdentified: []string{"object1", "object2"},
		InteractiveObjectsRemoved:   []string{"object3", "object4"},
		EnemiesIdentified:           []string{"enemy1", "enemy2"},
		EnemiesRemoved:              []string{"enemy3", "enemy4"},
		PlayerInventoryAdded:        []string{"item3", "item4"},
		PlayerInventoryRemoved:      []string{"item5", "item6"},
		StoryThreads:                []string{"Test Story Thread 1", "Test Story Thread 2"},
	}

	// Test with current values stateUpdate
	testGame.UpdateGameState(stateUpdate)
	if testGame.World.CurrentLocation.LocationName != "Test Current Location" {
		t.Errorf("Expected current location to be 'Test Current Location', but got %s", testGame.World.CurrentLocation.LocationName)
	}

	// Test with some inventory updates
	stateUpdate.PlayerInventoryAdded = []string{"item3", "item4"}
	testGame.UpdateGameState(stateUpdate)
	if len(testGame.Player.Inventory) != 2 {
		t.Errorf("Expected inventory to have 2 items, but got %d", len(testGame.Player.Inventory))
	}

	// Test with some story threads
	stateUpdate.StoryThreads = []string{"Test Story Thread 3"}
	testGame.UpdateGameState(stateUpdate)

	// Test with some location updates
	stateUpdate.PlayerLocation = "Test New Location"
	stateUpdate.PotentialLocations = []string{"Test Even Newer Location"}

	testGame.UpdateGameState(stateUpdate)
}
