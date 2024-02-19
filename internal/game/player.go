package game

import "github.com/sessionsdev/blue-octopus/internal/util"

type Player struct {
	Name      string
	Inventory util.StringSet
}
