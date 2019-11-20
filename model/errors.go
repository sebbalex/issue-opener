package model

// Error validation definition
type Error struct {
	Key    string `json:"key"`
	Reason string `json:"reason"`
}
