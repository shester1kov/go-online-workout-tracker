package models

import "time"

type HealthStatus struct {
	Status    string
	Timestamp time.Time
	Details   map[string]string
}

func (hs *HealthStatus) IsHealthy() bool {
	return hs.Status == "up"
}
