package models

import "time"

type Task struct {
	ID           int64     `json:"id"`
	Header       string    `json:"header"`
	Description  string    `json:"description"`
	CreationTime time.Time `json:"creation_time"`
}
