package model

import (
	"w2learn/pkg/def"
)

type HealthModel struct {
	ServerStatus   int `json:"serverStatus"`
	DatabaseStatus int `json:"databaseStatus"`
}

func GetDefaultHealthModel() *HealthModel {
	return &HealthModel{
		ServerStatus:   def.HealthStatusCheckOK,
		DatabaseStatus: def.HealthStatusNoCheck,
	}
}
