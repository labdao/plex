package models

import "time"

type UserEvent struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ApiKeyID  int       `json:"api_key_id"`
	EventTime time.Time `json:"event_time"`
	EventType string    `json:"event_type"`
}
