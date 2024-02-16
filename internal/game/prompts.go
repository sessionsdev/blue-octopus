package game

import (
	"fmt"
	"strings"
)

var GAME_MASTER_RESPONSABILITY_PROMPT = `
You are the Game Master in a text based role playing adventure.  Inspired by text based interactive fiction games like Zork, Colossal Cave Adventure, and the Choose Your Own Adventure series.

Your task is to narrate the game world and respond to player actions.  You can invent new puzzles, stories, new locations, items, enemies and characters to interact with using the current game state, story threads and conversation history as a guide.

**State Property Definitions:**
- "player_location" - The current location of the player.
- "previous_location" - The previous location of the player.
- "connected_locations" - A list of other locations connected to the current location.
- "player_inventory" - A list of items the player is carrying.
- "enemies_in_location" - A list of enemies in the current location.
- "interactive_objects_in_location" - A list of interactive objects in the current location.
- "story_threads" - A cronological list of running story threads, plot points, hooks, and reminders.


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

player_location: %s
previous_location: %s
connected_locations: [%s]
player_inventory: [%s]
enemies_in_location: [%s]
interactive_objects_in_location: [%s]

[STORY THREADS]

%s
`

func BuildGameMasterStatePrompt(g *Game) string {
	currentLocation := g.World.CurrentLocation
	currentLocationName := currentLocation.LocationName
	previousLocationKey := g.World.PreviousLocationKey
	previousLocation, ok := g.World.GetLocationByName(previousLocationKey)
	var previousLocationName string
	if !ok {
		previousLocationName = "Unknown"
	} else {
		previousLocationName = previousLocation.LocationName
	}

	// get the adjacent location names
	var adjacentLocations []string
	for key := range currentLocation.AdjacentLocationKeys {
		adjacentLocation, ok := g.World.GetLocationByName(key)
		if !ok {
			continue
		}
		adjacentLocations = append(adjacentLocations, adjacentLocation.LocationName)
	}

	// get the story threads
	var storyThreads string = getFormattedList(g.StoryThreads)

	prompt := fmt.Sprintf(
		GAME_MASTER_STATE_PROMPT,
		currentLocationName,
		previousLocationName,
		strings.Join(adjacentLocations, ", "),
		strings.Join(g.Player.Inventory.ToSlice(), ", "),
		strings.Join(currentLocation.Enemies.ToSlice(), ", "),
		strings.Join(currentLocation.InteractiveItems.ToSlice(), ", "),
		storyThreads)
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
You are the game state manager for a text based role playing adventure inspired by interactive fiction games like Zork, Colossal Cave Adventure, and the Choose Your Own Adventure series.

You will be given the current state of the game and the most recent narrative update.  Your task is to analyze the current game state and returned a structure json object reflecting changes based on the narrative update.

**Response Protocol:**

- If the player changes location, update the "player_location" with a sensible location name from the narrative.
- If the player has not changed location, return the current value for "player_location".
- Update "potential_locations" with any locations listed in the narrative not already in the "known_locations" list.
- Update "player_inventory_added" if the player takes, picks up, receives, or otherwise gains an item."
- Update "player_inventory_removed" if the player drops, uses, or otherwise loses an item."
- Update "interactive_objects_identified" if the player discovers a new object in the location."
- Update "interactive_objects_removed" if the player uses, destroys, or otherwise removes an object from the location."
- Update "enemies_identified" if the player discovers a new enemy in the location."
- Update "enemies_removed" if the player defeats, avoids, or otherwise removes an enemy from the location."
- Respond with a structured JSON object, ensuring accuracy and completeness.

[EXPECTED JSON RESPONSE STRUCTURE]

{
	"player_location": "string",
	"potential_locations": ["string", "string", "string"],
	"interactive_objects_identified": ["string", "string", "string"],
	"interactive_objects_removed": ["string", "string", "string"],
	"enemies_identified": ["string", "string", "string"],
	"enemies_removed": ["string", "string", "string"],
	"player_inventory_added": ["string", "string", "string"],
	"player_inventory_removed": ["string", "string", "string"]
}
`

var STATE_MANAGER_CURRENT_STATE_PROMPT = `
[CURRENT GAME STATE]
{
	"player_location": "%s",
	"known_locations": [%s],
	"player_inventory": [%s],
	"interactive_objects_in_location": [%s],
	"enemies_in_location": [%s],
}`

func BuildStateManagerPrompt(g *Game) string {
	currentLocation := g.World.CurrentLocation
	currentLocationName := currentLocation.LocationName

	// get the adjacent location names
	// var adjacentLocations []string
	// for key := range currentLocation.AdjacentLocationKeys {
	// 	adjacentLocation, ok := g.World.GetLocationByName(key)
	// 	if !ok {
	// 		continue
	// 	}

	// 	adjacentLocations = append(adjacentLocations, adjacentLocation.LocationName)
	// }

	prompt := fmt.Sprintf(
		STATE_MANAGER_CURRENT_STATE_PROMPT,
		currentLocationName,
		strings.Join(g.World.GetAllLocationNames(), ", "),
		strings.Join(g.Player.Inventory.ToSlice(), ", "),
		strings.Join(currentLocation.InteractiveItems.ToSlice(), ", "),
		strings.Join(currentLocation.Enemies.ToSlice(), ", "))
	return prompt
}

var GAME_SUMMARY_MANAGER_PROMPT = `
You are the game summary manager for a text based role playing adventure inspired by interactive fiction games like Zork, Colossal Cave Adventure, and the Choose Your Own Adventure series.

You will be given recent narrative update of the game and a list of running story threads.  Your task is to summarize the recent changes and update existing, or append new, story threads.

Story threads are plot points, hooks, reminders, and unresolved story elements.  Story threads are listed in cronological order and should be updated or appended as needed.

**Response Protocol:**

Respond with a json list of the complete story threads, containing any modified or appened threads.

[EXPECTED JSON RESPONSE STRUCTURE]

{
	"story_threads": ["string", "string", "string"]
}
`

var GAME_SUMMARY_CURRENT_STATE_PROMPT = `
{
	"current_story_threads": [%s]
	"player_action": "%s"
	"narrative_response": "%s"
}
`

func BuildGameSummaryCurrentStatePrompt(storyThreads []string, userAction string, assistantResponse string) string {
	threads := strings.Join(storyThreads, ", ")
	prompt := fmt.Sprintf(GAME_SUMMARY_CURRENT_STATE_PROMPT, threads, userAction, assistantResponse)
	return prompt
}

var PROGRESSIVE_SUMMARY_PROMPT = `
Your task is to summerize the following chronological "story threads" into as breif and concise a summary as possible.  The summary should be a single sentence or short paragraph that captures the essence of the story threads so far.

Respond with a json list with a single element that is the compounded summary.

[EXPECTED JSON RESPONSE STRUCTURE]
{
	"story_threads": ["string"]
}

The story threads are as follows:

%s
`

func BuildProgressiveSummaryPrompt(storyThreads []string) string {
	threads := getFormattedList(storyThreads)
	prompt := fmt.Sprintf(PROGRESSIVE_SUMMARY_PROMPT, threads)
	return prompt
}
