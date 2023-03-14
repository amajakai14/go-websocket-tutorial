package main

type MessageEnvelope struct {
	Action  int   `json:"action"`
	Message Menu  `json:"menu"`
	Target  *Room `json:"room"`
}

const (
	Join int = iota
	Leave
	Add
	Delete
	Reset
)
