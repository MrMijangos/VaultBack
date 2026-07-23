package entities

import "time"

type Address struct {
	ID             string
	UserID         string
	Label          string
	Recipient      string
	Phone          string
	Street         string
	City           string
	State          string
	PostalCode     string
	ReferenceNotes string
	IsDefault      bool
	CreatedAt      time.Time
}
