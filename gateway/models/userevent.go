package models

import "time"

type UserEvent struct {
	ID            int       `json:"id"`
	WalletAddress string    `json:"wallet_address"`
	ApiKeyID      int       `json:"api_key_id"`
	EventTime     time.Time `json:"event_time"`
	EventType     string    `json:"event_type"`
}
