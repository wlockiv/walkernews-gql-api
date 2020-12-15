package model

import f "github.com/fauna/faunadb-go/v3/faunadb"

type Comment struct {
	ID      string `fauna:"id"`
	Content string `fauna:"content"`
	Parent  f.RefV `fauna:"parent"` // Comment
	User    f.RefV `fauna:"user"`
}

func (c *Comment) IsPost() {}
