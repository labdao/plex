package models

import "time"

type FileEvent struct {
	ID        int       `json:"id"`
	FileID    int       `json:"file_id"`
	UserID    int       `json:"user_id"`
	EventTime time.Time `json:"event_time"`
	EventType string    `json:"event_type"`
}
