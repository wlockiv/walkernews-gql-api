package post

import (
	"errors"
	f "github.com/fauna/faunadb-go/v3/faunadb"
	"github.com/wlockiv/walkernews/graph/model"
	internalErrors "github.com/wlockiv/walkernews/internal/errors"
	"os"
)

func GetByRefV(refV f.RefV) (model.Post, error) {
	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_KEY"))
	res, err := client.Query(f.Get(refV))
	if err != nil {
		return nil, internalErrors.NewDBError("(Comment) GetByRefV", err)
	}

	if refV.Collection.ID == "comments" {
		var comment *model.Comment
		if err := res.At(f.ObjKey("data")).Get(&comment); err != nil {
			return nil, internalErrors.NewUnmarshallError("post/comment", err)
		}
		return comment, nil
	} else if refV.Collection.ID == "links" {
		var link *model.Link
		if err := res.At(f.ObjKey("data")).Get(&link); err != nil {
			return nil, internalErrors.NewUnmarshallError("post/link", err)
		}
		return link, nil
	}

	err = errors.New("(Post) GetByRefV: could not find collection '" + refV.Collection.ID + "'")
	return nil, internalErrors.NewBaseError(err)
}
