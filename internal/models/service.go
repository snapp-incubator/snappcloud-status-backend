package models

type Service struct {
	Name   string            `json:"name"`
	Query  string            `json:"-"`
	Status map[Region]Status `json:"status"`
}
