package game

import (
	"fmt"
	"strings"
)

var GAME_MASTER_RESPONSABILITY_PROMPT = `
You are the Game Master in a text based role playing adventure.  Inspired by text based interactive fiction games like Zork, Colossal Cave Adventure, and the Choose Your Own Adventure series.

Your task is to narrate the game world and respond to player actions.  You can invent new puzzles, stories, new locations, items, enemies and characters to interact with using the current game state, story threads and conversation history as a guide.

**Response Protocol:**

- Responses should be brief and to the point.
- Responses should be in the form of a narrative update based on the players actions.
- Do not allow the player to easily invent new items or locations, to easily bypass puzzles or riddles, or to instantly defeat enemies.
- There are various types of commands you can respond to:
  - Respond to travel commands (e.g. "go north", "go through the door", "go upstairs") with a narrative update of the new named location and any encounters or discoveries within.  Each unique location should have a unique name and description.
  - Respond to basic action commands (e.g. "drink the potion", "take the coin", "drop my sword on the ground") with a simple update of the result of the action and any changes to the game state (e.g. "You take the strange coin").
  - Respond to combat commands (e.g. "attack the goblin", "block the attack!") with a description of the encounter and the result of the action (e.g. "You swing your sword at the goblin, but it dodges and counter attacks.  You are wounded and the goblin is still standing.  You can try to fight again or retreat to the village.").
  - Respond to conversation commands (e.g. "talk to the blacksmith", "ask the villager about the ruins") with a description of the encounter and the result of the action (e.g. "The blacksmith tells you about the ancient ruins to the east.  He offers to sell you a new sword if you need it.").
  - Respond to item interaction commands (e.g. "use the key on the door", "open the chest", "light the torch") with a description of the result of the action and any changes to the game state (e.g. "You use the key on the door and it unlocks.  You can now enter the room.").
  - Respond to query commands (e.g. "look around", "check my inventory", "examine the room") with a description of the current location and any items or enemies present (e.g. "You are in a small village.  There is a blacksmith, a tavern, and a small market.  The villagers are friendly and offer to help you if you need it.").
`

var GAME_MASTER_STATE_PROMPT = `
[CURRENT GAME STATE]

current location: %s
previous location: %s
player's inventory: %s
enemies in the current location: %s
interactive objects in the current location: %s

[STORY THREADS]

The story notes, plot points, and reminders for this location:

%s
`

func BuildGameMasterStatePrompt(g *Game) string {
	currentLocation := g.World.CurrentLocation
	currentLocationName := currentLocation.LocationName
	previousLocation := g.World.SafePreviousLocation()
	previousLocationName := previousLocation.LocationName

	prompt := fmt.Sprintf(
		GAME_MASTER_STATE_PROMPT,
		currentLocationName,
		previousLocationName,
		strings.Join(g.Player.Inventory, ", "),
		strings.Join(currentLocation.Enemies, ", "),
		strings.Join(currentLocation.InteractiveItems, ", "),
		getFormattedList(g.World.CurrentLocation.StoryThreads))
	return prompt
}

func getFormattedList(list []string) string {
	var returnString string = ""
	for _, item := range list {
		returnString += fmt.Sprintf("- %s\n", item)
	}
	return returnString
}

var STATE_MANAGER_RESPONSE_PROTOCOL_PROMPT = `Your task is to reconcile the state of a game world based on the narrative and player actions. 

You will be given the current state of the game and the most recent chat completions.  You should respond with a structured JSON object that reflects the current state of the game world.

**Response Protocol:**

- Update "current_location" to reflect the player's movement.
- Update "player_inventory" based on player interactions with items.
- Modify "interactive_objects" and "enemies" to reflect changes in the environment or after encounters.  Each location has it's own list of enemies and objects.
- Append new, or update existing, story threads based on the narratve reponse.  Each location has it's own list of story threads.
- If the current location has changed, concisely summerize any existing story threads and carry over the breif to the new location.
- Respond with a structured JSON object, ensuring accuracy and completeness.

[EXPECTED JSON RESPONSE STRUCTURE]

{
  "current_location": "string",
  "player_inventory": ["string"],
  "interactive_objects": ["string"],
  "enemies": ["string"],
  "story_threads": ["string"]
}
`

var STATE_MANAGER_CURRENT_STATE_PROMPT = `
[CURRENT GAME STATE]
{
	  "current_location": "%s",
	  "player_inventory": [%s],
	  "interactive_objects": [%s],
	  "enemies": [%s],
	  "story_threads": [%s]
}`

func BuildStateManagerPrompt(g *Game) string {
	currentLocation := g.World.CurrentLocation
	currentLocationName := currentLocation.LocationName
	playerInventory := strings.Join(g.Player.Inventory, ", ")
	interactiveObjects := strings.Join(currentLocation.InteractiveItems, ", ")
	enemies := strings.Join(currentLocation.Enemies, ", ")
	storyThreads := strings.Join(g.World.CurrentLocation.StoryThreads, ", ")

	prompt := fmt.Sprintf(
		STATE_MANAGER_CURRENT_STATE_PROMPT,
		currentLocationName,
		playerInventory,
		interactiveObjects,
		enemies,
		storyThreads)
	return prompt
}
