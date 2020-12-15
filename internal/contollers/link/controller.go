package link

import (
	f "github.com/fauna/faunadb-go/v3/faunadb"
	"github.com/wlockiv/walkernews/graph/model"
	internalErr "github.com/wlockiv/walkernews/internal/errors"
	"os"
)

func Create(newLink model.NewLink, userKey string) (*model.Link, error) {
	client := f.NewFaunaClient(userKey)
	res, err := client.Query(
		f.Call("create_link", f.Arr{newLink.Title, newLink.Address}),
	)
	if err != nil {
		return nil, internalErr.NewDBError("(Link) Create", err)
	}

	var link *model.Link
	if err := res.At(f.ObjKey("data")).Get(&link); err != nil {
		return nil, internalErr.NewUnmarshallError("link", err)
	}

	return link, nil
}

func GetById(id string) (*model.Link, error) {
	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_KEY"))
	res, err := client.Query(
		f.Get(f.MatchTerm(f.Index("link_by_id"), id)))
	if err != nil {
		return nil, internalErr.NewDBError("(Link) GetById", err)
	}

	var link *model.Link
	if err := res.At(f.ObjKey("data")).Get(&link); err != nil {
		return nil, internalErr.NewUnmarshallError("link", err)
	}

	return link, nil
}

func GetByRefV(linkRef f.RefV) (*model.Link, error) {
	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_CLIENT_KEY"))
	res, err := client.Query(f.Get(linkRef))
	if err != nil {
		return nil, internalErr.NewDBError("(Link) GetByRefV", err)
	}

	var link *model.Link
	if err := res.At(f.ObjKey("data")).Get(&link); err != nil {
		return nil, internalErr.NewUnmarshallError("link", err)
	}

	return link, nil
}

func GetAll() ([]*model.Link, error) {
	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_CLIENT_KEY"))
	res, err := client.Query(
		f.Map(
			f.Paginate(f.Match(f.Index("links_sorted_by_createdAt_desc"))),
			f.Lambda([]string{"ts", "ref"}, f.Select("data", f.Get(f.Var("ref")))),
		),
	)
	if err != nil {
		return nil, internalErr.NewDBError("(Link) GetAll", err)
	}

	var links []*model.Link
	if err := res.At(f.ObjKey("data")).Get(&links); err != nil {
		return nil, internalErr.NewUnmarshallError("[]link", err)
	}

	return links, nil
}

func DeleteById(id, userKey string) (*model.Link, error) {
	client := f.NewFaunaClient(userKey)
	res, err := client.Query(
		f.Let().Bind(
			"linkRef", f.Select("ref", f.Get(f.MatchTerm(f.Index("link_ref_by_id"), id))),
		).In(
			f.Select("data", f.Delete(f.Var("linkRef"))),
		),
	)
	if err != nil {
		return nil, internalErr.NewDBError("DeleteById", err)
	}

	var link *model.Link
	if err := res.Get(&link); err != nil {
		return nil, internalErr.NewUnmarshallError("link", err)
	}

	return link, nil
}
