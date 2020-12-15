package post

import (
	"errors"
	f "github.com/fauna/faunadb-go/v3/faunadb"
	"github.com/wlockiv/walkernews/graph/model"
	internalErrors "github.com/wlockiv/walkernews/internal/errors"
	"os"
)

func GetByRefV(refV f.RefV, postType model.PostType) (model.Post, error) {
	client := f.NewFaunaClient(os.Getenv("FDB_SERVER_KEY"))
	res, err := client.Query(f.Get(refV))
	if err != nil {
		return nil, internalErrors.NewDBError("(Comment) GetByRefV", err)
	}

	if postType.String() == model.PostTypeComment.String() {
		var comment *model.Comment
		if err := res.At(f.ObjKey("data")).Get(&comment); err != nil {
			return nil, internalErrors.NewUnmarshallError("comment", err)
		} else {
			return comment, nil
		}
	}

	if postType.String() == model.PostTypeLink.String() {
		var link *model.Link
		if err := res.At(f.ObjKey("data")).Get(&link); err != nil {
			return nil, internalErrors.NewUnmarshallError("link", err)
		} else {
			return link, nil
		}
	}

	err = errors.New("(Post) GetByRefV: must include a valid PostType")
	return nil, internalErrors.NewBaseError(err)
}
