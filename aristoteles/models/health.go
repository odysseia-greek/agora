package models

import "encoding/json"

func UnmarshalDatabaseHealth(data []byte) (DatabaseHealth, error) {
	var r DatabaseHealth
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *DatabaseHealth) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type DatabaseHealth struct {
	Healthy       bool   `json:"healthy"`
	ClusterName   string `json:"clusterName,omitempty"`
	ServerName    string `json:"serverName,omitempty"`
	ServerVersion string `json:"serverVersion,omitempty"`
}
