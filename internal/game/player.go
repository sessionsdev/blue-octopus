package game

type Player struct {
	Name      string   `json:"name"`
	Inventory []string `json:"inventory"`
}
