package models

type Service struct {
	Name   string  `json:"name"`
	States []State `json:"states"`
}
