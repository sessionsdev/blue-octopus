package game

import (
	utils "github.com/sessionsdev/blue-octopus/internal"
)

type Player struct {
	Name      string          `json:"name"`
	Inventory utils.StringSet `json:"inventory"`
}
