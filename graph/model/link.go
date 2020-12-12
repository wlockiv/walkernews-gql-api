package model

import (
	f "github.com/fauna/faunadb-go/v3/faunadb"
	"os"
	"time"
)

type Link struct {
	ID        string    `fauna:"id"`
	Title     string    `fauna:"title"`
	Address   string    `fauna:"address"`
	CreatedAt time.Time `fauna:"createdAt"`
	UserID    string    `fauna:"userId"`
}

func (l *Link) incorporate(link Link) {
	l.ID = link.ID
	l.Title = link.Title
	l.Address = link.Address
	l.CreatedAt = link.CreatedAt
	l.UserID = link.UserID
}

func (l *Link) Save(userKey string) error {
	client := f.NewFaunaClient(userKey)
	res, err := client.Query(f.Create(
		f.Collection("Link"), f.Obj{
			"data": f.Obj{
				"id":        f.NewId(),
				"title":     l.Title,
				"address":   l.Address,
				"userId":    l.UserID,
				"createdAt": f.Now(),
			},
		},
	))
	if err != nil {
		return err
	}

	var link Link
	if err := res.At(f.ObjKey("data")).Get(&link); err != nil {
		return err
	} else {
		// Update the model
		l.ID = link.ID
	}

	return nil
}

func (l *Link) GetById(id string) error {
	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_KEY"))
	res, err := client.Query(
		f.Get(
			f.MatchTerm(f.Index("link_by_id"), id),
		),
	)
	if err != nil {
		return err
	}

	var link Link
	if err := res.At(f.ObjKey("data")).Get(&link); err != nil {
		return err
	}

	l.incorporate(link)

	return nil
}

func (l *Link) GetAll() ([]*Link, error) {
	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_KEY"))
	res, err := client.Query(
		f.Map(
			f.Paginate(f.Match(f.Index("links_sorted_by_createdAt_desc"))),
			f.Lambda([]string{"ts", "ref"}, f.Select("data", f.Get(f.Var("ref")))),
		),
	)
	if err != nil {
		panic(err)
	}

	var links []*Link
	if err := res.At(f.ObjKey("data")).Get(&links); err != nil {
		panic(err)
	}

	return links, nil
}

func NewLinkModel(title, address, userId string) *Link {
	link := Link{
		ID:        f.NewId().String(),
		Title:     title,
		Address:   address,
		CreatedAt: time.Now(),
		UserID:    userId,
	}

	return &link
}
