package game

import (
	"testing"

	utils "github.com/sessionsdev/blue-octopus/internal"
)

func TestUpdateGameState(t *testing.T) {
	testCurrentLocation := &Location{
		LocationName:      "Test Current Location",
		AdjacentLocations: []string{},
		InteractiveItems:  []string{},
		Enemies:           []string{},
	}

	testGame := &Game{
		Player: &Player{
			Inventory: []string{},
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
	stateUpdate.PlayerInventory = []string{"item1", "item2"}
	testGame.UpdateGameState(stateUpdate)

	// Test with some story threads
	stateUpdate.StoryThreads = []string{"thread2", "thread3"}
	testGame.UpdateGameState(stateUpdate)
	if len(testGame.StoryThreads) != 3 {
		t.Errorf("Expected story threads to have 3 items, but got %d", len(testGame.StoryThreads))
	}

	// Test with some location updates
	stateUpdate.UpdatedLocationName = "Test New Location"
	stateUpdate.UpdatedAdjacentLocations = []string{"Test Even Newer Location"}
	stateUpdate.PlayerInventory = []string{"item3", "item4"}
	stateUpdate.InteractiveObjectsInLocation = []string{"object1", "object2"}
	stateUpdate.EnemiesInLocation = []string{"enemy1", "enemy2"}
	stateUpdate.StoryThreads = []string{"thread4", "thread5"}
	testGame.UpdateGameState(stateUpdate)

}
