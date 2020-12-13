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

	l.ID = link.ID
	l.Title = link.Title
	l.Address = link.Address
	l.CreatedAt = link.CreatedAt
	l.UserID = link.UserID

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

func (l *Link) DeleteById(id, userKey string) error {
	client := f.NewFaunaClient(userKey)
	_, err := client.Query(
		f.Map(
			f.Paginate(f.MatchTerm("link_ref_by_id", id)),
			f.Lambda("LINK_REF", f.Delete(f.Var("LINK_REF"))),
		),
	)
	if err != nil {
		return err
	}

	return nil
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

//Obj{
//"data": Arr{
//Obj{"ref": RefV{ID: "284747993637193222", Collection: &RefV{ID: "Link", Collection: &RefV{ID: "collections"}}},
//"ts": 1607815698190000,
//"data": Obj{
//"createdAt": TimeV("2020-12-12T23:28:18.091775Z"),
//"id": "284747993637192198",
//"title": "Google",
//"userId": "284682068516930048",
//"address": "google.com"
//}}}}
