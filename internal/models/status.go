package models

type Status uint8

const (
	Operational Status = 1
	Warning     Status = 2
	Outage      Status = 3
)
