package models

type Service struct {
	Name   string `json:"name"`
	Status Status `json:"status"`
}
