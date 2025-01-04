package types

type MacroType int

const (
	MacroTypeSequence MacroType = iota
	MacroTypeToggle
	MacroTypeHold
)

type MacroAction struct {
	Keys  []int `json:"keys"`
	Delay int   `json:"delay"`
}

type Macro struct {
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	ActivationKeys []int         `json:"activation_keys"`
	Type           MacroType     `json:"type"`
	Actions        []MacroAction `json:"actions"`
	IncludeTitles  string        `json:"include_titles"`
}
