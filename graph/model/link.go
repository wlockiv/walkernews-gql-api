package model

import (
	"time"
)

type Link struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
	UserID    string    `json:"user"`
}
