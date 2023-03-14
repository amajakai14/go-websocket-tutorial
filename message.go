package main

type MessageEnvelope struct {
	Action  int   `json:"action"`
	Message Menu  `json:"menu"`
}

const (
	Join int = iota
	Add
	Delete
	Reset
)
