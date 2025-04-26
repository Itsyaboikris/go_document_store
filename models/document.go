package models

import "time"

type Document struct {
	ID        string                 `json:"_id"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}
