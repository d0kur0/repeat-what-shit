package types

type MacroType int

const (
	MacroTypeSequence MacroType = iota
	MacroTypeToggle
	MacroTypeHold
)

type MacroAction struct {
	ID    string `json:"id"`
	Keys  []int  `json:"keys"`
	Delay int    `json:"delay"`
}

type Macro struct {
	ID             string        `json:"id"`
	Disabled       bool          `json:"disabled"`
	Name           string        `json:"name"`
	ActivationKeys []int         `json:"activation_keys"`
	Type           MacroType     `json:"type"`
	Actions        []MacroAction `json:"actions"`
	IncludeTitle   []string      `json:"include_title"`
}
