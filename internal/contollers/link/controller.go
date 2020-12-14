package link

import (
	"errors"
	f "github.com/fauna/faunadb-go/v3/faunadb"
	"github.com/wlockiv/walkernews/graph/model"
	"os"
)

func Create(newLink model.NewLink, userKey string) (*model.Link, error) {
	client := f.NewFaunaClient(userKey)
	res, err := client.Query(f.Create(
		f.Collection("links"), f.Obj{
			"data": f.Obj{
				"id":        f.NewId(),
				"title":     newLink.Title,
				"address":   newLink.Address,
				"user":      f.Identity(),
				"createdAt": f.Now(),
			},
		},
	))
	if err != nil {
		return nil, err
	}

	var link *model.Link
	if err := res.At(f.ObjKey("data")).Get(&link); err != nil {
		return nil, err
	}

	return link, nil
}

func GetById(id string) (*model.Link, error) {
	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_KEY"))
	res, err := client.Query(
		f.Get(
			f.MatchTerm(f.Index("link_by_id"), id),
		),
	)
	if err != nil {
		return nil, err
	}

	var link *model.Link
	if err := res.At(f.ObjKey("data")).Get(&link); err != nil {
		return nil, err
	}

	return link, nil
}

func GetAll() ([]*model.Link, error) {
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

	var links []*model.Link
	if err := res.At(f.ObjKey("data")).Get(&links); err != nil {
		panic(err)
	}

	return links, nil
}

func DeleteById(id, userKey string) error {
	client := f.NewFaunaClient(userKey)
	res, err := client.Query(
		f.Map(
			f.Paginate(f.MatchTerm("link_ref_by_id", id)),
			f.Lambda("LINK_REF", f.Select("data", f.Delete(f.Var("LINK_REF")))),
		),
	)
	if err != nil {
		return err
	}

	var links []*model.Link
	if err := res.At(f.ObjKey("data")).Get(&links); err != nil {
		return err
	} else if len(links) == 0 {
		err := errors.New("a link with the provided id could not be found")
		return err
	}

	return nil
}
