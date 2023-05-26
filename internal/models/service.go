package models

type Service struct {
	Name   string            `json:"name"`
	Order  int               `json:"order"`
	Status map[Region]Status `json:"status"`
}
