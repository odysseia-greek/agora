package models

import "encoding/json"

func UnmarshalHealth(data []byte) (Health, error) {
	var r Health
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Health) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// swagger:model
type Health struct {
	// example: true
	Healthy bool `json:"healthy"`
	// example: 2023-06-07 15:02:11.678766777 +0000 UTC m=+5090.268683461
	Time     string         `json:"time"`
	Database DatabaseHealth `json:"databaseHealth,omitempty"`
	Memory   Memory         `json:"memory,omitempty"`
}

// swagger:model
type DatabaseHealth struct {
	// example: true
	Healthy bool `json:"healthy"`
	// example: aristoteles
	ClusterName string `json:"clusterName,omitempty"`
	// example: aristoteles-es-worker-0
	ServerName string `json:"serverName,omitempty"`
	// example: 8.8.0
	ServerVersion string `json:"serverVersion,omitempty"`
}

// swagger:model
type Memory struct {
	Free       uint64 `json:"free,omitempty"`
	Alloc      uint64 `json:"alloc,omitempty"`
	TotalAlloc uint64 `json:"totalAlloc,omitempty"`
	Sys        uint64 `json:"sys,omitempty"`
	Unit       string `json:"unit,omitempty"`
}
