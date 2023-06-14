package models

type Status string

const (
	Operational Status = "operational"
	Disruption  Status = "disruption"
	Outage      Status = "outage"
	Unknown     Status = "unknown"
)
