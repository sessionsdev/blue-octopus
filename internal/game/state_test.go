package game

import (
	"testing"

	utils "github.com/sessionsdev/blue-octopus/internal"
)

func TestUpdateGameState(t *testing.T) {
	testCurrentLocation := &Location{
		LocationName:      "Test Current Location",
		AdjacentLocations: utils.EmptyStringSet(),
		InteractiveItems:  utils.EmptyStringSet(),
		Enemies:           utils.EmptyStringSet(),
	}

	testGame := &Game{
		Player: &Player{
			Inventory: utils.EmptyStringSet(),
		},
		World: &World{
			Locations: map[string]*Location{
				"test_current_location": testCurrentLocation,
			},
			CurrentLocation:  testCurrentLocation,
			PreviousLocation: nil,
			VisitedLocations: utils.EmptyStringSet(),
		},
		StoryThreads: []string{
			"Test Story Thread 1",
		},
	}

	stateUpdate := GameStateUpdateResponse{}

	// Test with empty stateUpdate
	testGame.UpdateGameState(stateUpdate)
	if testGame.World.CurrentLocation.LocationName != "Test Current Location" {
		t.Errorf("Expected current location to be 'Test Current Location', but got %s", testGame.World.CurrentLocation.LocationName)
	}

	// Test with some inventory updates
	stateUpdate.InventoryUpdates.Added = []string{"item1", "item2"}
	stateUpdate.InventoryUpdates.Removed = []string{"item1"}
	testGame.UpdateGameState(stateUpdate)
	if !testGame.Player.Inventory.Contains("item2") {
		t.Errorf("Expected inventory to contain 'item2', but it didn't")
	}

	// Test with some story threads
	stateUpdate.StoryThreads = []string{"thread2", "thread3"}
	testGame.UpdateGameState(stateUpdate)
	if len(testGame.StoryThreads) != 3 {
		t.Errorf("Expected story threads to have 3 items, but got %d", len(testGame.StoryThreads))
	}

	// Test with some location updates
	stateUpdate.CurrentLocation = "Test New Location"
	stateUpdate.UpdatedAdjacentLocations = []string{"Test Even Newer Location"}
	stateUpdate.InteractiveObjectsUpdates.Added = []string{"object1", "object2"}
	stateUpdate.InteractiveObjectsUpdates.Removed = []string{"object3", "object4"}
	stateUpdate.EnemiesUpdates.Added = []string{"enemy1", "enemy2"}
	stateUpdate.EnemiesUpdates.Defeated = []string{"enemy1"}
	stateUpdate.StoryThreads = []string{"thread4", "thread5"}
	testGame.UpdateGameState(stateUpdate)

	if testGame.World.CurrentLocation.LocationName != "Test New Location" {
		t.Errorf("Expected current location to be 'Test New Location', but got %s", testGame.World.CurrentLocation.LocationName)
	}

	if !testGame.World.CurrentLocation.AdjacentLocations.Contains("Test Even Newer Location") {
		t.Errorf("Expected current location to have adjacent location 'Test Even Newer Location', but it didn't")
	}

	if !testGame.World.CurrentLocation.AdjacentLocations.Contains("Test Current Location") {
		t.Errorf("Expected current location to have adjacent location 'Test Current Location', but it didn't")
	}

	if !testGame.World.CurrentLocation.InteractiveItems.Contains("object1") {
		t.Errorf("Expected current location to have interactive item 'object1', but it didn't")
	}

	if testGame.World.CurrentLocation.InteractiveItems.Contains("object3") {
		t.Errorf("Expected current location to not have interactive item 'object3', but it did")
	}

	if !testGame.World.CurrentLocation.Enemies.Contains("enemy2") {
		t.Errorf("Expected current location to have enemy 'enemy2', but it didn't")
	}

	if testGame.World.CurrentLocation.Enemies.Contains("enemy1") {
		t.Errorf("Expected current location to not have enemy 'enemy1', but it did")
	}

	if len(testGame.StoryThreads) != 5 {
		t.Errorf("Expected story threads to have 5 items, but got %d", len(testGame.StoryThreads))
	}

}
