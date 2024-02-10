package game

import utils "github.com/sessionsdev/blue-octopus/internal"

type Message interface {
	NewMessage(string, string) Message
}

type GameMessage struct {
	Provider string `json:"provider"`
	Message  string `json:"message"`
}

type NewGameDetails struct {
	StartingLocation          string   `json:"starting_location"`
	PlayerName                string   `json:"player_name"`
	PlayerInventory           []string `json:"player_inventory"`
	StartingAdjacentLocations []string `json:"starting_adjacent_locations"`
	MainQuest                 string   `json:"main_quest"`
	StoryThreads              []string `json:"story_threads"`
}

func (m *GameMessage) NewMessage(provider string, message string) Message {
	return &GameMessage{Provider: provider, Message: message}
}

type Game struct {
	GameId             string        `json:"game_id"`
	World              *World        `json:"world"`
	Player             *Player       `json:"player"`
	MainQuest          string        `json:"main_quest"`
	StoryThreads       []string      `json:"story_threads"`
	GameMessageHistory []GameMessage `json:"game_message_history"`
	TotalTokensUsed    int           `json:"total_tokens_used"`
}

func (g *Game) GetRecentHistory(numItems int) []GameMessage {
	currentHistory := g.GameMessageHistory
	if len(currentHistory) > 5 {
		// Take the 5 most recent history items
		return currentHistory[len(currentHistory)-5:]
	} else {
		return currentHistory
	}
}

func BuildNewGame(details NewGameDetails) *Game {
	game := &Game{
		World: &World{
			Locations:        make(map[string]*Location),
			CurrentLocation:  nil,
			PreviousLocation: nil,
			VisitedLocations: utils.EmptyStringSet(),
		},
		Player: &Player{
			Name:      details.PlayerName,
			Inventory: details.PlayerInventory,
		},
		MainQuest:          details.MainQuest,
		StoryThreads:       details.StoryThreads,
		GameMessageHistory: []GameMessage{},
		TotalTokensUsed:    0,
	}

	// Create the starting location
	startingLocation := &Location{
		LocationName:      details.StartingLocation,
		AdjacentLocations: details.StartingAdjacentLocations,
		InteractiveItems:  []string{},
		Enemies:           []string{},
	}

	// Add the starting location to the world
	game.World.SafeAddLocation(startingLocation.LocationName)
	game.World.CurrentLocation = startingLocation

	return game
}
