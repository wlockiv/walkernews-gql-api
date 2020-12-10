package model

import (
	"time"
)

type Link struct {
	ID        string    `json:"id" dynamo:"id"`
	Title     string    `json:"title" dynamo:"title"`
	Address   string    `json:"address" dynamo:"address"`
	CreatedAt time.Time `json:"createdAt" dynamo:"createdAt"`
	UserID    string    `json:"user" dynamo:"user"`
}
