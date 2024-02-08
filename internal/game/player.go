package game

type Player struct {
	Name      string              `json:"name"`
	Inventory map[string]struct{} `json:"inventory"`
}
