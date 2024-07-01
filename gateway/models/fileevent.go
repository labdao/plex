package models

import "time"

type FileEvent struct {
	ID            int       `json:"id"`
	FileID        int       `json:"file_id"`
	WalletAddress string    `json:"wallet_address"`
	EventTime     time.Time `json:"event_time"`
	EventType     string    `json:"event_type"`
}
