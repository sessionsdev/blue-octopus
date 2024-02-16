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
	StartingStoryThreads      []string `json:"starting_story_threads"`
	PlayerName                string   `json:"player_name"`
	PlayerInventory           []string `json:"player_inventory"`
	StartingAdjacentLocations []string `json:"starting_adjacent_locations"`
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
			Locations:           make(map[string]*Location),
			CurrentLocation:     nil,
			PreviousLocationKey: "",
		},
		Player: &Player{
			Name:      details.PlayerName,
			Inventory: utils.NewStringSet(details.PlayerInventory...),
		},
		GameMessageHistory: []GameMessage{},
		TotalTokensUsed:    0,
	}

	// Add the starting location to the world
	game.World.SafeAddLocation(details.StartingLocation)
	game.World.CurrentLocation, _ = game.World.GetLocationByName(details.StartingLocation)
	game.World.CurrentLocation.AdjacentLocationKeys = utils.EmptyStringSet()

	for _, locationName := range details.StartingAdjacentLocations {
		game.World.SafeAddLocation(locationName)
		location, _ := game.World.GetLocationByName(locationName)
		location.AdjacentLocationKeys.AddAll(game.World.CurrentLocation.getNormalizedName())
		game.World.CurrentLocation.AdjacentLocationKeys.AddAll(location.getNormalizedName())
	}

	return game
}
