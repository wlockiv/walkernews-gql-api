package model

import (
	f "github.com/fauna/faunadb-go/v3/faunadb"
	"time"
)

type Link struct {
	ID        string    `fauna:"id"`
	Title     string    `fauna:"title"`
	Address   string    `fauna:"address"`
	CreatedAt time.Time `fauna:"createdAt"`
	User      f.RefV    `fauna:"user"`
}
