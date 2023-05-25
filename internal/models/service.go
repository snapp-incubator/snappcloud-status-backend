package models

type Service struct {
	Name   string            `json:"name"`
	Status map[Region]Status `json:"status"`
}
