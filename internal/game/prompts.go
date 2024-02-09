package game

import (
	"fmt"
	"log"
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
[Current Game State]
{
  main_quest: %s,
  current_location: %s,
  previous_location: %s,
  adjacent_locations: [%s],
  players_inventory: [%s],
  enemies": [%s],
  interactive_objects: [%s],
  story_threads: [%s]
}
`

func BuildGameStatePrompt(g *Game) string {
	mainQuest := g.MainQuest
	currentLocation := g.World.CurrentLocation
	currentLocationName := currentLocation.LocationName
	previousLocation := g.World.SafePreviousLocation()
	previousLocationName := previousLocation.LocationName

	prompt := fmt.Sprintf(
		GAME_MASTER_STATE_PROMPT,
		mainQuest,
		currentLocationName,
		previousLocationName,
		strings.Join(currentLocation.AdjacentLocations, ", "),
		strings.Join(g.Player.Inventory, ", "),
		strings.Join(currentLocation.Enemies, ", "),
		strings.Join(currentLocation.InteractiveItems, ", "),
		getFormattedList(g.StoryThreads))

	log.Println("GAME_STATE_PROMPT: ", prompt)
	return prompt
}

func getFormattedList(list []string) string {
	var returnString string = ""
	for _, item := range list {
		returnString += fmt.Sprintf("- %s\n", item)
	}
	return returnString
}

var STATE_MANAGER_RESPONSE_PROTOCOL_PROMPT = `
Your task is to update the game state based on a narrative context and current state of the game.

**JSON Response Structure:**

State updates will always apply to the current location in the response.

Respond with a JSON object containing the following fields:

{
  "current_location": "Location Name",
  "adjacent_locations": ["existing", "and", "added", "adjacent", "locations"],
  "player_inventory": ["complete", "list", "of", "player", "inventory"],
  "interactive_objects": ["complete", "list", "of", "interactive", "objects"],
  "enemies": ["complete", "list", "of", "enemies"],
  "new_story_threads": ["new", "story", "threads", "to", "append"]
}

**Examples:**

*Initial Provided State:*
{
  "main_quest": "Find the lost treasure",
  "current_location": "Mountain Pass",
  "previous_location": "Village",
  "adjacent_locations": ["Forest", "Cave"],
  "player_inventory": ["sword", "shield", "potion"],
  "enemies": ["goblin", "skeleton", "dragon"],
  "interactive_objects": ["chest", "door", "key"],
  "story_threads": ["meet the blacksmith", "explore the forest"]
}

*Player Action:*
Player: "Go into the cave"
Assistant: "You find yourself in a dark cave.  The air is damp and the sound of dripping water echoes through the chamber.  You can see a faint light to the north."
`
