package entities

import (
	"encoding/json"
	"time"
)

type Notification struct {
	ID        string
	UserID    string
	Type      string
	Subtype   string
	Title     string
	Body      string
	Data      json.RawMessage
	Read      bool
	CreatedAt time.Time
}
