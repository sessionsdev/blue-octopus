package game

import "encoding/json"

var SETUP_PROMPT = `
**Game Master Role in Text-Based Adventure:**

As the Game Master, you orchestrate a text-based adventure, drawing inspiration from classics like Zork and Colossal Cave Adventure. You are tasked with guiding players through a dynamically evolving world, creating locations, characters, and storylines in response to their journey. Your narrative should adapt to player actions, enriching the game with new challenges and discoveries.

**Responsibilities:**

- **Creative World-Building:** Continuously introduce new locations, characters, and items, enriching the game world.
- **Engaging Narration:** Provide vivid descriptions of scenes, characters, and challenges, enhancing the immersive experience.
- **Challenge Simulation:** Design encounters requiring strategy, making gameplay rewarding.
- **Dynamic Interaction:** React to player inputs by weaving new storylines and challenges, fitting seamlessly into the narrative.
- **Proactive State Expansion:** Leverage player actions to suggest new locations, enemies, and plot developments.
- **Adaptive Storytelling:** Craft a narrative that evolves with player actions, steering the game towards new intrigues.

**Response Protocol:**

Combine detailed narrative descriptions with a structured JSON object to outline game state expansions. Encourage exploration and progression by aligning new elements with player actions and storylines.

**JSON Template:**

{
  "response": "Narrative detailing new encounters or items.",
  "proposed_state_changes": {
    "new_current_location": null or "Location Name",
    "new_adjacent_locations": ["New Locations"],
    "updated_enemies_in_location": ["Enemies List"],
    "updated_interactive_objects_in_location": ["Objects List"],
    "updated_removable_items_in_location": ["Items List"],
    "updated_player_inventory": ["Inventory Items"],
    "story_threads": ["The list of new and modified story threads, in their entirety, that should be active after this response."]
  }
}

**Examples:**

1. **Player Action:** "I examine the mailbox."
   - **Narrative Response:** "Approaching the mailbox, you find it shimmering oddly. Inside lies a mysterious, glowing key."
   - **State Change:**
  {
    "response": "Opening the mailbox reveals a glowing key.",
    "proposed_state_changes": { 
      "updated_player_inventory": ["Glowing Key"],
      "story_threads": [
        "Previous Story Thread 1",
        "Previous Story Thread 2","
        "The Mystery of the Glowing Key is unfolding."]
    }
  }

2. **Player Action:** "I head east towards the forest."
   - **Narrative Response:** "You find yourself on the dark forest's edge, filled with whispers and the scent of adventure."
   - **State Change:**
  {
    "response": "Entering the whispering forest, adventure calls.",
    "proposed_state_changes": {
      "new_current_location": "Whispering Forest",
      "new_adjacent_locations": ["Ancient Ruins", "Mystic River"],
      "story_threads": [
        "The Forest's Whisper"
        ]
    }
  }

Ensure each response and state update reflects the evolving game world, offering players new opportunities for exploration and interaction. Your creativity shapes a unique and memorable adventure for each player.

**Initial Prompt for Players:** "You are standing in an open field west of a white house, with a boarded front door. There is a small mailbox here."
`

type GameStatePromptDetails struct {
	CurrentLocation       string   `json:"current_location"`
	AdjacentLocationNames []string `json:"adjacent_locations"`
	Inventory             []string `json:"player_inventory"`
	EnemiesInLocation     []string `json:"enemies_in_location"`
	InteractiveItems      []string `json:"interactive_objects_in_location"`
	RemovableItems        []string `json:"removable_items_in_location"`
	CentralPlot           string   `json:"central_plot"`
	StoryThreads          []string `json:"story_threads"`
}

func GetJsonFieldDescriptionsForPromptDetails() string {
	return `{
		"current_location": "The name of the current location.",
		"adjacent_locations": ["A list of the names of the adjacent locations."],
		"player_inventory": ["A list of the items in the player's inventory."],
		"enemies_in_location": ["A list of the enemies in the current location."],
		"interactive_objects_in_location": ["A list of the interactive objects in the current location.",
		"removable_items_in_location": ["A list of the removable items in the current location."],
		"central_plot": "The central plot of the game.",
		"story_threads": ["A list of the story threads."]
	}`
}

func (g *Game) BuildGameStatePromptDetails() GameStatePromptDetails {
	currentLocation := g.World.CurrentLocation
	player := g.Player

	if currentLocation == nil {
		currentLocation = &Location{}
	}

	if player == nil {
		player = &Player{}
	}

	adjacentLocationNames := make([]string, len(currentLocation.AdjacentLocations))
	for i, location := range currentLocation.AdjacentLocations {
		adjacentLocationNames[i] = location
	}

	return GameStatePromptDetails{
		CurrentLocation:       currentLocation.LocationName,
		AdjacentLocationNames: adjacentLocationNames,
		Inventory:             g.Player.Inventory,
		EnemiesInLocation:     currentLocation.EnemiesInLocation,
		InteractiveItems:      currentLocation.InteractiveItems,
		RemovableItems:        currentLocation.RemovableItems,
		CentralPlot:           g.CentralPlot,
		StoryThreads:          g.StoryThreads,
	}
}

func (gs *GameStatePromptDetails) GetJsonOrEmptyString() string {
	jsonString, err := json.Marshal(gs)
	if err != nil {
		return ""
	}

	return string(jsonString)
}
